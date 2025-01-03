package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CreateMessage 创建站内信
func CreateMessage(senderID *uint64, receiverID uint64, msgType string, title string, content string) error {
	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Type:       msgType,
		Title:      title,
		Content:    content,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}
	return config.DB.Create(message).Error
}

// GetMessageList 获取站内信列表
func GetMessageList(userID uint64, req *models.MessageListRequest) (*models.MessageListResponse, error) {
	var messages []models.Message
	var total int64
	var unreadCount int64

	query := config.DB.Model(&models.Message{}).Where("receiver_id = ?", userID)

	// 添加消息类型筛选
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 添加已读状态筛选（不包括团队申请消息）
	if req.IsRead != nil {
		query = query.Where("(type != ? OR (type = ? AND is_read = ?))",
			models.MessageTypeTeamApplication,
			models.MessageTypeTeamApplication,
			false)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取未读数（不包括团队申请消息）
	if err := config.DB.Model(&models.Message{}).
		Where("receiver_id = ? AND is_read = ? AND type != ?",
			userID, false, models.MessageTypeTeamApplication).
		Count(&unreadCount).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(req.PageSize).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	// 为每条消息添加可执行的操作，并处理团队申请消息的特殊显示
	for i := range messages {
		// 检查是否是管理员（对于团队申请消息）
		isAdmin := false
		if messages[i].Type == models.MessageTypeTeamApplication && !messages[i].IsProcessed {
			var application models.TeamApplication
			if err := config.DB.Where("id = ?", messages[i].ApplicationID).
				First(&application).Error; err == nil {
				role, err := GetTeamUserRole(application.TeamID, userID)
				if err == nil && (role == "owner" || role == "admin") {
					isAdmin = true
				}
			}
		}
		messages[i].Actions = models.GetMessageActions(&messages[i], isAdmin)

		// 团队申请消息不显示已读状态
		if messages[i].Type == models.MessageTypeTeamApplication {
			messages[i].IsRead = false
			messages[i].ReadAt = nil
		}
	}

	return &models.MessageListResponse{
		Messages:    messages,
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
		UnreadCount: unreadCount,
	}, nil
}

// MarkMessageAsRead 标记消息为已读
func MarkMessageAsRead(messageID uint64, userID uint64) error {
	// 检查消息类型，团队申请消息不能标记已读
	var message models.Message
	if err := config.DB.First(&message, messageID).Error; err != nil {
		return err
	}
	if message.Type == models.MessageTypeTeamApplication {
		return fmt.Errorf("团队申请消息不能标记已读")
	}

	now := time.Now()
	result := config.DB.Model(&models.Message{}).
		Where("id = ? AND receiver_id = ? AND is_read = ?", messageID, userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("消息不存在、无权限或已读")
	}
	return nil
}

// BatchMarkMessagesAsRead 批量标记消息为已读
func BatchMarkMessagesAsRead(messageIDs []uint64, userID uint64) error {
	if len(messageIDs) == 0 {
		return nil
	}

	// 检查是否包含团队申请消息
	var count int64
	if err := config.DB.Model(&models.Message{}).
		Where("id IN ? AND type = ?", messageIDs, models.MessageTypeTeamApplication).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("团队申请消息不能标记已读")
	}

	now := time.Now()
	result := config.DB.Model(&models.Message{}).
		Where("id IN ? AND receiver_id = ? AND is_read = ?", messageIDs, userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})

	if result.Error != nil {
		return result.Error
	}
	return nil
}

// MarkAllMessagesAsRead 标记所有消息为已读
func MarkAllMessagesAsRead(userID uint64) error {
	now := time.Now()
	return config.DB.Model(&models.Message{}).
		Where("receiver_id = ? AND is_read = ? AND type != ?", userID, false, models.MessageTypeTeamApplication).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// DeleteMessage 删除消息
func DeleteMessage(messageID uint64, userID uint64) error {
	// 检查消息类型，团队申请消息不能删除
	var message models.Message
	if err := config.DB.First(&message, messageID).Error; err != nil {
		return err
	}
	if message.Type == models.MessageTypeTeamApplication {
		return fmt.Errorf("团队申请消息不能删除")
	}

	result := config.DB.Where("id = ? AND receiver_id = ?", messageID, userID).
		Delete(&models.Message{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("消息不存在或无权限")
	}
	return nil
}

// GetUnreadMessageCount 获取未读消息数量
func GetUnreadMessageCount(userID uint64) (*models.UnreadCountResponse, error) {
	var total int64
	// 获取总未读数（排除团队申请消息）
	err := config.DB.Model(&models.Message{}).
		Where("receiver_id = ? AND is_read = ? AND type != ?",
			userID, false, models.MessageTypeTeamApplication).
		Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 按类型统计未读数
	var results []struct {
		Type  string
		Count int64
	}
	err = config.DB.Model(&models.Message{}).
		Select("type, count(*) as count").
		Where("receiver_id = ? AND ((type != ? AND is_read = ?) OR (type = ? AND is_processed = ?))",
			userID,
			models.MessageTypeTeamApplication, false,
			models.MessageTypeTeamApplication, false).
		Group("type").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	// 构建响应
	byTypes := make(map[string]models.TypeCountInfo)
	for _, result := range results {
		// 对于团队申请消息，使用未处理数量而不是未读数量
		if result.Type == models.MessageTypeTeamApplication {
			byTypes[result.Type] = models.TypeCountInfo{
				Count:       result.Count,
				Description: "待处理的" + models.GetMessageTypeDescription(result.Type),
			}
		} else {
			byTypes[result.Type] = models.TypeCountInfo{
				Count:       result.Count,
				Description: models.GetMessageTypeDescription(result.Type),
			}
		}
	}

	return &models.UnreadCountResponse{
		Total:   total,
		ByTypes: byTypes,
	}, nil
}

// CreateTeamApplication 创建团队申请
func CreateTeamApplication(userID uint64, req *models.TeamApplicationRequest) error {
	// 检查用户是否已经是团队成员
	isMember, err := IsTeamMember(req.TeamID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return fmt.Errorf("您已经是团队成员")
	}

	// 检查是否有待处理的申请
	var count int64
	if err := config.DB.Model(&models.TeamApplication{}).
		Where("team_id = ? AND user_id = ? AND status = ?", req.TeamID, userID, "pending").
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("您已有待处理的申请")
	}

	// 开启事务
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 创建申请记录
		application := &models.TeamApplication{
			TeamID:    req.TeamID,
			UserID:    userID,
			Message:   req.Message,
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := tx.Create(application).Error; err != nil {
			return err
		}

		// 获取团队信息
		var team models.Team
		if err := tx.First(&team, req.TeamID).Error; err != nil {
			return err
		}

		// 获取申请者信息
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		// 向团队所有者和管理员发送通知
		var admins []models.TeamMember
		if err := tx.Where("team_id = ? AND role IN ('owner', 'admin')", req.TeamID).
			Find(&admins).Error; err != nil {
			return err
		}

		// 创建关联的消息记录
		for _, admin := range admins {
			message := &models.Message{
				SenderID:      &userID,
				ReceiverID:    admin.UserID,
				Type:          models.MessageTypeTeamApplication,
				Title:         fmt.Sprintf("新的团队申请 - %s", team.Name),
				Content:       fmt.Sprintf("用户 %s 申请加入团队 %s\n申请信息：%s", user.Username, team.Name, req.Message),
				IsRead:        false,
				IsProcessed:   false,
				ApplicationID: application.ID,
				CreatedAt:     time.Now(),
			}
			if err := tx.Create(message).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetTeamApplicationList 获取团队申请列表
func GetTeamApplicationList(req *models.TeamApplicationListRequest, userID uint64) (*models.TeamApplicationListResponse, error) {
	var applications []models.TeamApplication
	var total int64

	// 构建基础查询
	query := config.DB.Model(&models.TeamApplication{})

	// 如果指定了团队ID，检查用户权限
	if req.TeamID > 0 {
		role, err := GetTeamUserRole(req.TeamID, userID)
		if err != nil {
			return nil, err
		}
		if role != "owner" && role != "admin" {
			return nil, fmt.Errorf("权限不足")
		}
		query = query.Where("team_id = ?", req.TeamID)
	} else {
		// 如果没有指定团队ID，只能查看自己的申请
		query = query.Where("user_id = ?", userID)
	}

	// 添加状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(req.PageSize).
		Find(&applications).Error; err != nil {
		return nil, err
	}

	// 构建响应
	responses := make([]models.TeamApplicationResponse, len(applications))
	for i, app := range applications {
		responses[i].Application = app

		// 获取团队信息
		if err := config.DB.First(&responses[i].Team, app.TeamID).Error; err != nil {
			return nil, err
		}

		// 获取用户信息
		var user models.User
		if err := config.DB.Select("id, username, email").
			First(&user, app.UserID).Error; err != nil {
			return nil, err
		}
		responses[i].User.ID = user.ID
		responses[i].User.Username = user.Username
		responses[i].User.Email = user.Email
	}

	return &models.TeamApplicationListResponse{
		Applications: responses,
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}, nil
}

// HandleTeamApplication 处理团队申请
func HandleTeamApplication(req *models.TeamApplicationHandleRequest, operatorID uint64) error {
	var application models.TeamApplication
	if err := config.DB.First(&application, req.ApplicationID).Error; err != nil {
		return err
	}

	// 检查操作者权限
	role, err := GetTeamUserRole(application.TeamID, operatorID)
	if err != nil {
		return err
	}
	if role != "owner" && role != "admin" {
		return fmt.Errorf("权限不足")
	}

	// 检查申请状态
	if application.Status != "pending" {
		return fmt.Errorf("该申请已被处理")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 更新申请状态
		application.Status = req.Status
		application.UpdatedAt = time.Now()
		if err := tx.Save(&application).Error; err != nil {
			return err
		}

		// 如果同意申请，添加为团队成员
		if req.Status == "approved" {
			member := &models.TeamMember{
				TeamID:   application.TeamID,
				UserID:   application.UserID,
				Role:     "member",
				JoinedAt: time.Now(),
			}
			if err := tx.Create(member).Error; err != nil {
				return err
			}
		}

		// 获取团队信息
		var team models.Team
		if err := tx.First(&team, application.TeamID).Error; err != nil {
			return err
		}

		// 获取操作者信息
		var operator models.User
		if err := tx.First(&operator, operatorID).Error; err != nil {
			return err
		}

		// 向申请者发送通知
		resultMessage := &models.Message{
			SenderID:   &operatorID,
			ReceiverID: application.UserID,
			Type:       "team_application_result",
			Title:      fmt.Sprintf("团队申请结果 - %s", team.Name),
			Content: fmt.Sprintf("您申请加入团队 %s 的请求已%s\n处理人：%s\n处理意见：%s",
				team.Name,
				getStatusText(req.Status),
				operator.Username,
				req.Message),
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(resultMessage).Error; err != nil {
			return err
		}

		// 更新原申请消息的状态
		if err := tx.Model(&models.Message{}).
			Where("type = ? AND application_id = ?", models.MessageTypeTeamApplication, application.ID).
			Updates(map[string]interface{}{
				"is_processed": true,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// getStatusText 获取状态文本
func getStatusText(status string) string {
	if status == "approved" {
		return "通过"
	}
	return "被拒绝"
}

// IsTeamMember 检查用户是否是团队成员
func IsTeamMember(teamID uint64, userID uint64) (bool, error) {
	var count int64
	err := config.DB.Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count).Error
	return count > 0, err
}
