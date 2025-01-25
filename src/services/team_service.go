package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTeam 创建团队
func CreateTeam(req *models.CreateTeamRequest, creatorID uint64) (uint64, error) {
	team := &models.Team{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   creatorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return team.ID, config.DB.Transaction(func(tx *gorm.DB) error {
		// 创建团队
		if err := tx.Create(team).Error; err != nil {
			return err
		}

		// 添加创建者为团队所有者
		member := &models.TeamMember{
			TeamID:   team.ID,
			UserID:   creatorID,
			Role:     "owner",
			JoinedAt: time.Now(),
		}
		return tx.Create(member).Error
	})
}

// UpdateTeam 更新团队信息
func UpdateTeam(teamID uint64, req *models.UpdateTeamRequest, userID uint64) error {
	// 检查用户权限
	role, err := GetTeamUserRole(teamID, userID)
	if err != nil {
		return err
	}
	if role != "owner" && role != "admin" {
		return errors.New("权限不足")
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	updates["updated_at"] = time.Now()

	return config.DB.Model(&models.Team{}).Where("id = ?", teamID).Updates(updates).Error
}

// DeleteTeam 删除团队
func DeleteTeam(teamID uint64, userID uint64) error {
	// 检查用户权限
	role, err := GetTeamUserRole(teamID, userID)
	if err != nil {
		return err
	}
	if role != "owner" {
		return errors.New("只有团队所有者可以删除团队")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 删除团队成员
		if err := tx.Where("team_id = ?", teamID).Delete(&models.TeamMember{}).Error; err != nil {
			return err
		}

		// 删除团队作业和相关题目
		var assignments []models.TeamAssignment
		if err := tx.Where("team_id = ?", teamID).Find(&assignments).Error; err != nil {
			return err
		}
		for _, assignment := range assignments {
			if err := tx.Where("assignment_id = ?", assignment.ID).Delete(&models.TeamAssignmentProblem{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("team_id = ?", teamID).Delete(&models.TeamAssignment{}).Error; err != nil {
			return err
		}

		// 删除团队题单和相关题目
		var lists []models.TeamProblemList
		if err := tx.Where("team_id = ?", teamID).Find(&lists).Error; err != nil {
			return err
		}
		for _, list := range lists {
			if err := tx.Where("list_id = ?", list.ID).Delete(&models.TeamProblemListItem{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("team_id = ?", teamID).Delete(&models.TeamProblemList{}).Error; err != nil {
			return err
		}

		// 删除团队邀请
		if err := tx.Where("team_id = ?", teamID).Delete(&models.TeamInvitation{}).Error; err != nil {
			return err
		}

		// 删除团队
		return tx.Delete(&models.Team{}, teamID).Error
	})
}

// GetTeamDetail 获取团队详情
func GetTeamDetail(teamID uint64, userID uint64) (*models.TeamDetail, error) {
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		return nil, err
	}

	detail := &models.TeamDetail{
		Team: team,
	}

	// 获取成员数量
	var count int64
	if err := config.DB.Model(&models.TeamMember{}).Where("team_id = ?", teamID).Count(&count).Error; err != nil {
		return nil, err
	}
	detail.MemberCount = int(count)

	// 获取用户角色和是否是团队成员
	if userID > 0 {
		role, _ := GetTeamUserRole(teamID, userID)
		detail.UserRole = role
		detail.IsJoined = role != ""
	}

	// 获取创建者基本信息
	owner, err := GetTeamOwnerInfo(teamID, team.CreatedBy)
	if err != nil {
		return nil, err
	}
	detail.Owner = owner

	return detail, nil
}

// GetTeamOwnerInfo 获取团队创建者信息
func GetTeamOwnerInfo(teamID uint64, ownerID uint64) (*models.TeamOwnerInfo, error) {
	// 查询用户基本信息
	var user struct {
		ID       uint64 `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := config.DB.Table("users").
		Select("id, username, email").
		Where("id = ?", ownerID).
		First(&user).Error; err != nil {
		return nil, err
	}

	// 查询团队内名称
	var nickname string
	err := config.DB.Table("team_nicknames").
		Select("nickname").
		Where("team_id = ? AND user_id = ?", teamID, ownerID).
		Take(&nickname).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &models.TeamOwnerInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: nickname,
	}, nil
}

// GetTeamList 获取团队列表
func GetTeamList(req *models.TeamListRequest, userID uint64) (*models.TeamListResponse, error) {
	var teams []models.Team
	var total int64

	// 构建基础查询
	query := config.DB.Model(&models.Team{})

	// 如果是查询已加入的团队
	if req.Scope == "joined" && userID > 0 {
		query = query.Joins("INNER JOIN team_members ON teams.id = team_members.team_id").
			Where("team_members.user_id = ?", userID)
	}

	// 如果提供了关键字，搜索团队名称和描述
	if req.Keyword != "" {
		query = query.Where("teams.name LIKE ? OR teams.description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&teams).Error; err != nil {
		return nil, err
	}

	// 转换为 TeamDetail 并填充额外信息
	teamDetails := make([]models.TeamDetail, len(teams))
	for i, team := range teams {
		teamDetails[i].Team = team

		// 获取成员数量
		var count int64
		if err := config.DB.Model(&models.TeamMember{}).Where("team_id = ?", team.ID).Count(&count).Error; err != nil {
			return nil, err
		}
		teamDetails[i].MemberCount = int(count)

		// 获取用户角色
		if userID > 0 {
			role, _ := GetTeamUserRole(team.ID, userID)
			teamDetails[i].UserRole = role
			teamDetails[i].IsJoined = role != ""
		}

		// 获取创建者信息
		owner, err := GetTeamOwnerInfo(team.ID, team.CreatedBy)
		if err != nil {
			return nil, err
		}
		teamDetails[i].Owner = owner
	}

	return &models.TeamListResponse{
		Teams:    teamDetails,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetTeamUserRole 获取用户在团队中的角色
func GetTeamUserRole(teamID uint64, userID uint64) (string, error) {
	var member models.TeamMember
	err := config.DB.Where("team_id = ? AND user_id = ?", teamID, userID).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

// CreateTeamInvitation 创建团队邀请
func CreateTeamInvitation(teamID uint64, userID uint64) (*models.TeamInvitation, error) {
	// 检查用户权限
	role, err := GetTeamUserRole(teamID, userID)
	if err != nil {
		return nil, err
	}
	if role != "owner" && role != "admin" {
		return nil, errors.New("权限不足")
	}

	// 生成邀请码
	code := generateInviteCode()

	invitation := &models.TeamInvitation{
		TeamID:    teamID,
		Code:      code,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7天有效期
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(invitation).Error; err != nil {
		return nil, err
	}

	return invitation, nil
}

// JoinTeamByInvitation 通过邀请码加入团队
func JoinTeamByInvitation(code string, userID uint64) error {
	var invitation models.TeamInvitation
	if err := config.DB.Where("code = ? AND expires_at > ?", code, time.Now()).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("无效或已过期的邀请码")
		}
		return err
	}

	// 检查用户是否已经是团队成员
	var count int64
	if err := config.DB.Model(&models.TeamMember{}).Where("team_id = ? AND user_id = ?", invitation.TeamID, userID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("您已经是团队成员")
	}

	// 添加用户为团队成员
	member := &models.TeamMember{
		TeamID:   invitation.TeamID,
		UserID:   userID,
		Role:     "member",
		JoinedAt: time.Now(),
	}

	return config.DB.Create(member).Error
}

// generateInviteCode 生成邀请码
func generateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	rand.Seed(time.Now().UnixNano())
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// CreateAssignment 创建团队作业
func CreateAssignment(req *models.CreateAssignmentRequest, userID uint64) (uint64, error) {
	// 检查用户权限
	role, err := GetTeamUserRole(req.TeamID, userID)
	if err != nil {
		return 0, err
	}
	if role != "owner" && role != "admin" {
		return 0, errors.New("权限不足")
	}

	// 验证时间
	if req.StartTime.After(req.EndTime) {
		return 0, errors.New("开始时间不能晚于结束时间")
	}

	assignment := &models.TeamAssignment{
		TeamID:      req.TeamID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return assignment.ID, config.DB.Transaction(func(tx *gorm.DB) error {
		// 创建作业
		if err := tx.Create(assignment).Error; err != nil {
			return err
		}

		// 添加题目
		for _, problem := range req.Problems {
			ap := &models.TeamAssignmentProblem{
				AssignmentID: assignment.ID,
				ProblemID:    problem.ProblemID,
				OrderIndex:   problem.OrderIndex,
				Score:        problem.Score,
			}
			if err := tx.Create(ap).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateAssignment 更新团队作业
func UpdateAssignment(assignmentID uint64, req *models.UpdateAssignmentRequest, userID uint64) error {
	var assignment models.TeamAssignment
	if err := config.DB.First(&assignment, assignmentID).Error; err != nil {
		return err
	}

	// 检查用户权限
	role, err := GetTeamUserRole(assignment.TeamID, userID)
	if err != nil {
		return err
	}
	if role != "owner" && role != "admin" {
		return errors.New("权限不足")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 更新作业信息
		updates := make(map[string]interface{})
		if req.Title != "" {
			updates["title"] = req.Title
		}
		if req.Description != "" {
			updates["description"] = req.Description
		}
		if !req.StartTime.IsZero() {
			updates["start_time"] = req.StartTime
		}
		if !req.EndTime.IsZero() {
			updates["end_time"] = req.EndTime
		}
		updates["updated_at"] = time.Now()

		if err := tx.Model(&assignment).Updates(updates).Error; err != nil {
			return err
		}

		// 如果提供了新的题目列表，更新题目
		if req.Problems != nil {
			// 删除旧的题目
			if err := tx.Where("assignment_id = ?", assignmentID).Delete(&models.TeamAssignmentProblem{}).Error; err != nil {
				return err
			}

			// 添加新的题目
			for _, problem := range req.Problems {
				ap := &models.TeamAssignmentProblem{
					AssignmentID: assignmentID,
					ProblemID:    problem.ProblemID,
					OrderIndex:   problem.OrderIndex,
					Score:        problem.Score,
				}
				if err := tx.Create(ap).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// CreateProblemList 创建团队题单
func CreateProblemList(req *models.CreateProblemListRequest, userID uint64) (uint64, error) {
	// 检查用户权限
	role, err := GetTeamUserRole(req.TeamID, userID)
	if err != nil {
		return 0, err
	}
	if role == "" {
		return 0, errors.New("您不是团队成员")
	}

	list := &models.TeamProblemList{
		TeamID:      req.TeamID,
		Title:       req.Title,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return list.ID, config.DB.Transaction(func(tx *gorm.DB) error {
		// 创建题单
		if err := tx.Create(list).Error; err != nil {
			return err
		}

		// 添加题目
		for _, problem := range req.Problems {
			item := &models.TeamProblemListItem{
				ListID:     list.ID,
				ProblemID:  problem.ProblemID,
				OrderIndex: problem.OrderIndex,
				Note:       problem.Note,
			}
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateProblemList 更新团队题单
func UpdateProblemList(listID uint64, req *models.UpdateProblemListRequest, userID uint64) error {
	var list models.TeamProblemList
	if err := config.DB.First(&list, listID).Error; err != nil {
		return err
	}

	// 检查用户权限
	if list.CreatedBy != userID {
		role, err := GetTeamUserRole(list.TeamID, userID)
		if err != nil {
			return err
		}
		if role != "owner" && role != "admin" {
			return errors.New("权限不足")
		}
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 更新题单信息
		updates := make(map[string]interface{})
		if req.Title != "" {
			updates["title"] = req.Title
		}
		if req.Description != "" {
			updates["description"] = req.Description
		}
		updates["is_public"] = req.IsPublic
		updates["updated_at"] = time.Now()

		if err := tx.Model(&list).Updates(updates).Error; err != nil {
			return err
		}

		// 如果提供了新的题目列表，更新题目
		if req.Problems != nil {
			// 删除旧的题目
			if err := tx.Where("list_id = ?", listID).Delete(&models.TeamProblemListItem{}).Error; err != nil {
				return err
			}

			// 添加新的题目
			for _, problem := range req.Problems {
				item := &models.TeamProblemListItem{
					ListID:     listID,
					ProblemID:  problem.ProblemID,
					OrderIndex: problem.OrderIndex,
					Note:       problem.Note,
				}
				if err := tx.Create(item).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// UpdateTeamMemberRole 更新团队成员角色
func UpdateTeamMemberRole(teamID uint64, targetUserID uint64, newRole string, operatorID uint64) error {
	// 检查操作者权限
	operatorRole, err := GetTeamUserRole(teamID, operatorID)
	if err != nil {
		return err
	}
	if operatorRole != "owner" {
		return errors.New("只有团队所有者可以修改成员角色")
	}

	// 不能修改自己的角色
	if targetUserID == operatorID {
		return errors.New("不能修改自己的角色")
	}

	// 验证新角色是否有效
	validRoles := map[string]bool{"admin": true, "member": true}
	if !validRoles[newRole] {
		return errors.New("无效的角色")
	}

	// 更新角色
	result := config.DB.Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamID, targetUserID).
		Update("role", newRole)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不是团队成员")
	}

	return nil
}

// RemoveTeamMember 移除团队成员
func RemoveTeamMember(teamID uint64, targetUserID uint64, operatorID uint64) error {
	// 检查操作者权限
	operatorRole, err := GetTeamUserRole(teamID, operatorID)
	if err != nil {
		return err
	}

	// 获取目标用户角色
	targetRole, err := GetTeamUserRole(teamID, targetUserID)
	if err != nil {
		return err
	}

	// 权限检查
	if operatorRole != "owner" {
		if operatorRole != "admin" || targetRole == "owner" || targetRole == "admin" {
			return errors.New("权限不足")
		}
	}

	// 不能移除自己
	if targetUserID == operatorID {
		return errors.New("不能移除自己")
	}

	// 不能移除团队所有者
	if targetRole == "owner" {
		return errors.New("不能移除团队所有者")
	}

	// 移除成员
	result := config.DB.Where("team_id = ? AND user_id = ?", teamID, targetUserID).Delete(&models.TeamMember{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不是团队成员")
	}

	return nil
}

// GetAssignmentDetail 获取作业详情
func GetAssignmentDetail(assignmentID uint64, userID uint64) (*models.TeamAssignment, error) {
	var assignment models.TeamAssignment
	if err := config.DB.First(&assignment, assignmentID).Error; err != nil {
		return nil, err
	}

	// 检查用户权限
	role, err := GetTeamUserRole(assignment.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, errors.New("您不是团队成员")
	}

	// 获取作业题目
	var problems []models.TeamAssignmentProblem
	if err := config.DB.Where("assignment_id = ?", assignmentID).
		Order("order_index").
		Find(&problems).Error; err != nil {
		return nil, err
	}

	// TODO: 获取每个题目的提交状态

	return &assignment, nil
}

// GetAssignmentList 获取作业列表
func GetAssignmentList(teamID uint64, userID uint64) ([]models.TeamAssignment, error) {
	// 检查用户权限
	role, err := GetTeamUserRole(teamID, userID)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, errors.New("您不是团队成员")
	}

	var assignments []models.TeamAssignment
	if err := config.DB.Where("team_id = ?", teamID).
		Order("created_at DESC").
		Find(&assignments).Error; err != nil {
		return nil, err
	}

	return assignments, nil
}

// GetProblemListDetail 获取题单详情
func GetProblemListDetail(listID uint64, userID uint64) (*models.TeamProblemList, error) {
	var list models.TeamProblemList
	if err := config.DB.First(&list, listID).Error; err != nil {
		return nil, err
	}

	// 检查访问权限
	if !list.IsPublic {
		role, err := GetTeamUserRole(list.TeamID, userID)
		if err != nil {
			return nil, err
		}
		if role == "" {
			return nil, errors.New("无权访问题单")
		}
	}

	// 获取题单题目
	var problems []models.TeamProblemListItem
	if err := config.DB.Where("list_id = ?", listID).
		Order("order_index").
		Find(&problems).Error; err != nil {
		return nil, err
	}

	// TODO: 获取每个题目的提交状态

	return &list, nil
}

// GetProblemListList 获取题单列表
func GetProblemListList(teamID uint64, userID uint64) ([]models.TeamProblemList, error) {
	query := config.DB.Where("team_id = ?", teamID)

	// 如果不是团队成员，只能看到公开题单
	role, _ := GetTeamUserRole(teamID, userID)
	if role == "" {
		query = query.Where("is_public = ?", true)
	}

	var lists []models.TeamProblemList
	if err := query.Order("created_at DESC").Find(&lists).Error; err != nil {
		return nil, err
	}

	return lists, nil
}

// GetTeamAvatarPath 获取团队头像存储路径
func GetTeamAvatarPath() string {
	baseDir := filepath.Join("data", "team", "avatar")
	os.MkdirAll(baseDir, 0755)
	return baseDir
}

// SaveTeamAvatar 保存团队头像
func SaveTeamAvatar(teamID uint64, file *multipart.FileHeader) (*models.TeamAvatar, error) {
	var resultAvatar *models.TeamAvatar

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return fmt.Errorf("不支持的文件格式，仅支持 jpg、jpeg 和 png")
		}

		// 生成唯一文件名
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		filePath := filepath.Join(GetTeamAvatarPath(), filename)

		// 保存文件
		if err := os.MkdirAll(GetTeamAvatarPath(), 0755); err != nil {
			return fmt.Errorf("创建目录失败: %v", err)
		}

		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("打开文件失败: %v", err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("创建目标文件失败: %v", err)
		}
		defer dst.Close()

		// 读取文件内容并写入新文件
		buffer := make([]byte, 1024*1024) // 1MB buffer
		for {
			n, err := src.Read(buffer)
			if err != nil {
				break
			}
			if _, err := dst.Write(buffer[:n]); err != nil {
				return fmt.Errorf("写入文件失败: %v", err)
			}
		}

		// 更新数据库
		avatar := &models.TeamAvatar{
			TeamID:     teamID,
			Filename:   filename,
			UploadTime: time.Now(),
		}

		// 删除旧头像
		var oldAvatar models.TeamAvatar
		if err := tx.Where("team_id = ?", teamID).First(&oldAvatar).Error; err == nil {
			// 删除旧文件
			oldPath := filepath.Join(GetTeamAvatarPath(), oldAvatar.Filename)
			os.Remove(oldPath)
			// 更新记录
			if err := tx.Model(&oldAvatar).Updates(avatar).Error; err != nil {
				return fmt.Errorf("更新头像记录失败: %v", err)
			}
		} else {
			// 创建新记录
			if err := tx.Create(avatar).Error; err != nil {
				return fmt.Errorf("保存头像记录失败: %v", err)
			}
		}

		// 更新团队的avatar字段和更新时间
		now := time.Now()
		updates := map[string]interface{}{
			"avatar":     filename,
			"updated_at": now,
		}
		if err := tx.Model(&models.Team{}).Where("id = ?", teamID).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新团队头像字段失败: %v", err)
		}

		resultAvatar = avatar
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resultAvatar, nil
}

// GetTeamAvatar 获取团队头像
func GetTeamAvatar(filename string) (*models.TeamAvatar, error) {
	var avatar models.TeamAvatar
	if err := config.DB.Where("filename = ?", filename).First(&avatar).Error; err != nil {
		return nil, err
	}
	return &avatar, nil
}

// RemoveTeamAvatar 删除团队头像
func RemoveTeamAvatar(teamID uint64) error {
	var avatar models.TeamAvatar
	// 查找团队的头像记录
	if err := config.DB.Where("team_id = ?", teamID).First(&avatar).Error; err != nil {
		return fmt.Errorf("未找到头像记录")
	}

	// 删除文件
	filePath := filepath.Join(GetTeamAvatarPath(), avatar.Filename)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除头像文件失败: %v", err)
	}

	// 删除数据库记录
	if err := config.DB.Delete(&avatar).Error; err != nil {
		return fmt.Errorf("删除头像记录失败: %v", err)
	}

	// 清空团队的avatar字段
	if err := config.DB.Model(&models.Team{}).Where("id = ?", teamID).Update("avatar", "").Error; err != nil {
		return fmt.Errorf("更新团队头像字段失败: %v", err)
	}

	return nil
}

// UpdateTeamNickname 更新团队内名称
func UpdateTeamNickname(teamID uint64, userID uint64, nickname string) error {
	// 检查用户是否是团队成员
	role, err := GetTeamUserRole(teamID, userID)
	if err != nil {
		return err
	}
	if role == "" {
		return errors.New("您不是团队成员")
	}

	if nickname == "" {
		// 如果昵称为空，则删除记录
		return config.DB.Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&models.TeamNickname{}).Error
	}

	// 更新或创建团队内名称
	return config.DB.Exec(`
		INSERT INTO team_nicknames (team_id, user_id, nickname)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
		nickname = VALUES(nickname),
		updated_at = CURRENT_TIMESTAMP
	`, teamID, userID, nickname).Error
}

// GetTeamNickname 获取团队内名称
func GetTeamNickname(teamID uint64, userID uint64) (string, error) {
	var nickname string
	err := config.DB.Table("team_nicknames").
		Select("nickname").
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Take(&nickname).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	return nickname, err
}

// GetTeamMemberList 获取团队成员列表
func GetTeamMemberList(teamID uint64, req *models.TeamMemberListRequest, userID uint64) (*models.TeamMemberListResponse, error) {
	// 检查用户权限
	role, err := GetTeamUserRole(teamID, userID)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, errors.New("您不是团队成员")
	}

	// 构建基础查询
	query := config.DB.Table("team_members").
		Select(`
			team_members.user_id,
			users.username,
			avatars.filename as avatar,
			team_members.role,
			team_members.joined_at,
			team_nicknames.nickname
		`).
		Joins("LEFT JOIN users ON team_members.user_id = users.id").
		Joins("LEFT JOIN avatars ON team_members.user_id = avatars.user_id").
		Joins("LEFT JOIN team_nicknames ON team_members.team_id = team_nicknames.team_id AND team_members.user_id = team_nicknames.user_id").
		Where("team_members.team_id = ?", teamID)

	// 添加角色筛选
	if req.Role != "" {
		query = query.Where("team_members.role = ?", req.Role)
	}

	// 添加关键字搜索
	if req.Keyword != "" {
		query = query.Where("(users.username LIKE ? OR team_nicknames.nickname LIKE ?)",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	var members []models.TeamMemberInfo
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).
		Order("FIELD(team_members.role, 'owner', 'admin', 'member') ASC, team_members.joined_at DESC").
		Scan(&members).Error; err != nil {
		return nil, err
	}

	return &models.TeamMemberListResponse{
		Members:  members,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
