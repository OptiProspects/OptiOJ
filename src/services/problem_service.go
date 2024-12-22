package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// isValidDifficulty 验证难度等级是否有效
func isValidDifficulty(system models.DifficultySystem, difficulty string) bool {
	return models.IsValidDifficulty(system, difficulty)
}

// CreateProblem 创建题目
func CreateProblem(req *models.CreateProblemRequest, createdBy uint64) (uint64, error) {
	// 验证难度等级
	if !isValidDifficulty(req.DifficultySystem, req.Difficulty) {
		return 0, fmt.Errorf("无效的难度等级: %s", models.GetDifficultyDisplay(req.DifficultySystem, req.Difficulty))
	}

	problem := &models.Problem{
		Title:             req.Title,
		Description:       req.Description,
		InputDescription:  req.InputDescription,
		OutputDescription: req.OutputDescription,
		Samples:           req.Samples,
		Hint:              req.Hint,
		Source:            req.Source,
		DifficultySystem:  req.DifficultySystem,
		Difficulty:        req.Difficulty,
		TimeLimit:         req.TimeLimit,
		MemoryLimit:       req.MemoryLimit,
		IsPublic:          req.IsPublic,
		CreatedBy:         createdBy,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return problem.ID, config.DB.Transaction(func(tx *gorm.DB) error {
		// 创建题目
		if err := tx.Create(problem).Error; err != nil {
			return err
		}

		// 添加分类关系
		if len(req.CategoryIDs) > 0 {
			var values []map[string]interface{}
			for _, categoryID := range req.CategoryIDs {
				values = append(values, map[string]interface{}{
					"problem_id":  problem.ID,
					"category_id": categoryID,
				})
			}
			if err := tx.Table("problem_category_relations").Create(values).Error; err != nil {
				return err
			}
		}

		// 添加标签关系
		if len(req.TagIDs) > 0 {
			var values []map[string]interface{}
			for _, tagID := range req.TagIDs {
				values = append(values, map[string]interface{}{
					"problem_id": problem.ID,
					"tag_id":     tagID,
				})
			}
			if err := tx.Table("problem_tag_relations").Create(values).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateProblem 更新题目
func UpdateProblem(problemID uint64, req *models.UpdateProblemRequest) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 获取原题目信息
		var problem models.Problem
		if err := tx.First(&problem, problemID).Error; err != nil {
			return fmt.Errorf("题目不存在: %v", err)
		}

		// 更新题目基本信息
		updates := make(map[string]interface{})
		if req.Title != nil {
			updates["title"] = *req.Title
		}
		if req.Description != nil {
			updates["description"] = *req.Description
		}
		if req.InputDescription != nil {
			updates["input_description"] = *req.InputDescription
		}
		if req.OutputDescription != nil {
			updates["output_description"] = *req.OutputDescription
		}
		if req.Samples != nil {
			// 验证样例数据格式
			if !json.Valid([]byte(*req.Samples)) {
				return fmt.Errorf("样例数据必须是有效的 JSON 格式")
			}
			updates["samples"] = *req.Samples
		}
		if req.Hint != nil {
			updates["hint"] = *req.Hint
		}
		if req.Source != nil {
			updates["source"] = *req.Source
		}

		// 更新难度等级系统
		if req.DifficultySystem != nil {
			updates["difficulty_system"] = *req.DifficultySystem
			// 如果只更新难度等级系统，需要验证现有难度等级是否符合新系统
			if req.Difficulty == nil && !isValidDifficulty(*req.DifficultySystem, problem.Difficulty) {
				updates["difficulty"] = models.DifficultyNormalUnrated // 重置为暂无评级
			}
		}

		// 更新难度等级
		if req.Difficulty != nil {
			system := problem.DifficultySystem
			if req.DifficultySystem != nil {
				system = *req.DifficultySystem
			}
			if !isValidDifficulty(system, *req.Difficulty) {
				return fmt.Errorf("无效的难度等级: %s", models.GetDifficultyDisplay(system, *req.Difficulty))
			}
			updates["difficulty"] = *req.Difficulty
		}

		if req.TimeLimit != nil {
			updates["time_limit"] = *req.TimeLimit
		}
		if req.MemoryLimit != nil {
			updates["memory_limit"] = *req.MemoryLimit
		}
		if req.IsPublic != nil {
			updates["is_public"] = *req.IsPublic
		}
		updates["updated_at"] = time.Now()

		if err := tx.Model(&problem).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新题目基本信息失败: %v", err)
		}

		// 更新分类关系
		if req.CategoryIDs != nil {
			// 删除旧的分类关系
			if err := tx.Table("problem_category_relations").Where("problem_id = ?", problemID).Delete(nil).Error; err != nil {
				return fmt.Errorf("删除旧的分类关系失败: %v", err)
			}

			// 添加新的分类关系
			if len(req.CategoryIDs) > 0 {
				var values []map[string]interface{}
				for _, categoryID := range req.CategoryIDs {
					// 检查分类是否存在
					var count int64
					if err := tx.Model(&models.ProblemCategory{}).Where("id = ?", categoryID).Count(&count).Error; err != nil {
						return fmt.Errorf("检查分类是否存在失败: %v", err)
					}
					if count == 0 {
						return fmt.Errorf("分类 ID %d 不存在", categoryID)
					}

					values = append(values, map[string]interface{}{
						"problem_id":  problemID,
						"category_id": categoryID,
					})
				}
				if err := tx.Table("problem_category_relations").Create(values).Error; err != nil {
					return fmt.Errorf("添加新的分类关系失败: %v", err)
				}
			}
		}

		// 更新标签关系
		if req.TagIDs != nil {
			// 删除旧的标签关系
			if err := tx.Table("problem_tag_relations").Where("problem_id = ?", problemID).Delete(nil).Error; err != nil {
				return fmt.Errorf("删除旧的标签关系失败: %v", err)
			}

			// 添加新的标签关系
			if len(req.TagIDs) > 0 {
				var values []map[string]interface{}
				for _, tagID := range req.TagIDs {
					// 检查标签是否存在
					var count int64
					if err := tx.Model(&models.ProblemTag{}).Where("id = ?", tagID).Count(&count).Error; err != nil {
						return fmt.Errorf("检查标签是否存在失败: %v", err)
					}
					if count == 0 {
						return fmt.Errorf("标签 ID %d 不存在", tagID)
					}

					values = append(values, map[string]interface{}{
						"problem_id": problemID,
						"tag_id":     tagID,
					})
				}
				if err := tx.Table("problem_tag_relations").Create(values).Error; err != nil {
					return fmt.Errorf("添加新的标签关系失败: %v", err)
				}
			}
		}

		return nil
	})
}

// DeleteProblem 删除题目
func DeleteProblem(problemID uint64) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 获取该题目所有的提交记录ID
		var submissionIDs []uint64
		if err := tx.Model(&models.Submission{}).
			Where("problem_id = ?", problemID).
			Pluck("id", &submissionIDs).Error; err != nil {
			return err
		}

		// 获取该题目所有的测试用例ID
		var testCaseIDs []uint64
		if err := tx.Model(&models.TestCase{}).
			Where("problem_id = ?", problemID).
			Pluck("id", &testCaseIDs).Error; err != nil {
			return err
		}

		// 删除所有相关的判题结果（包括提交记录和测试用例相关的）
		if err := tx.Where("submission_id IN ? OR test_case_id IN ?", submissionIDs, testCaseIDs).
			Delete(&models.JudgeResult{}).Error; err != nil {
			return err
		}

		// 删除提交记录
		if err := tx.Where("problem_id = ?", problemID).
			Delete(&models.Submission{}).Error; err != nil {
			return err
		}

		// 获取并删除测试用例相关数据
		var testCases []models.TestCase
		if err := tx.Where("problem_id = ?", problemID).
			Find(&testCases).Error; err != nil {
			return err
		}

		// 删除测试用例文件
		for _, tc := range testCases {
			os.Remove(tc.InputFile)
			os.Remove(tc.OutputFile)
		}

		// 删除测试用例记录
		if err := tx.Where("problem_id = ?", problemID).
			Delete(&models.TestCase{}).Error; err != nil {
			return err
		}

		// 删除题目相关的其他数据
		if err := tx.Where("problem_id = ?", problemID).
			Delete(&models.ProblemCategory{}).Error; err != nil {
			return err
		}
		if err := tx.Where("problem_id = ?", problemID).
			Delete(&models.ProblemTag{}).Error; err != nil {
			return err
		}

		// 最后删除题目
		if err := tx.Delete(&models.Problem{}, problemID).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetProblemDetail 获取题目详情
func GetProblemDetail(problemID uint64, userID uint64) (*models.ProblemDetail, error) {
	var problem models.Problem
	if err := config.DB.First(&problem, problemID).Error; err != nil {
		return nil, err
	}

	// 检查访问权限
	if !problem.IsPublic {
		isAdmin, _ := IsAdmin(uint(userID))
		if !isAdmin && problem.CreatedBy != userID {
			return nil, errors.New("无权访问该题目")
		}
	}

	var detail models.ProblemDetail
	detail.Problem = problem

	// 获取分类信息
	if err := config.DB.Model(&problem).
		Select("problem_categories.*").
		Joins("JOIN problem_category_relations ON problem_category_relations.problem_id = problems.id").
		Joins("JOIN problem_categories ON problem_categories.id = problem_category_relations.category_id").
		Find(&detail.Categories).Error; err != nil {
		return nil, err
	}

	// 获取标签信息
	if err := config.DB.Model(&problem).
		Select("problem_tags.*").
		Joins("JOIN problem_tag_relations ON problem_tag_relations.problem_id = problems.id").
		Joins("JOIN problem_tags ON problem_tags.id = problem_tag_relations.tag_id").
		Find(&detail.Tags).Error; err != nil {
		return nil, err
	}

	return &detail, nil
}

// GetProblemList 获取题目列表
func GetProblemList(req *models.ProblemListRequest, userID uint64) (*models.ProblemListResponse, error) {
	query := config.DB.Model(&models.Problem{})

	// 处理查询条件
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.DifficultySystem != nil {
		query = query.Where("difficulty_system = ?", *req.DifficultySystem)
	}
	if req.Difficulty != "" {
		query = query.Where("difficulty = ?", req.Difficulty)
	}
	if req.CategoryID != nil {
		query = query.Joins("JOIN problem_category_relations ON problem_category_relations.problem_id = problems.id").
			Where("problem_category_relations.category_id = ?", *req.CategoryID)
	}
	if len(req.TagIDs) > 0 {
		query = query.Joins("JOIN problem_tag_relations ON problem_tag_relations.problem_id = problems.id").
			Where("problem_tag_relations.tag_id IN ?", req.TagIDs)
	}
	if req.IsPublic != nil {
		query = query.Where("is_public = ?", *req.IsPublic)
	}

	// 检查用户权限
	isAdmin, _ := IsAdmin(uint(userID))
	if !isAdmin {
		query = query.Where("is_public = ? OR created_by = ?", true, userID)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	var problems []models.Problem
	if err := query.Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&problems).Error; err != nil {
		return nil, err
	}

	// 获取每个题目的分类和标签信息
	var problemDetails []models.ProblemDetail
	for _, problem := range problems {
		detail := models.ProblemDetail{Problem: problem}

		// 获取分类信息
		if err := config.DB.Model(&problem).
			Select("problem_categories.*").
			Joins("JOIN problem_category_relations ON problem_category_relations.problem_id = problems.id").
			Joins("JOIN problem_categories ON problem_categories.id = problem_category_relations.category_id").
			Find(&detail.Categories).Error; err != nil {
			return nil, err
		}

		// 获取标签信息
		if err := config.DB.Model(&problem).
			Select("problem_tags.*").
			Joins("JOIN problem_tag_relations ON problem_tag_relations.problem_id = problems.id").
			Joins("JOIN problem_tags ON problem_tags.id = problem_tag_relations.tag_id").
			Find(&detail.Tags).Error; err != nil {
			return nil, err
		}

		problemDetails = append(problemDetails, detail)
	}

	return &models.ProblemListResponse{
		Problems: problemDetails,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UploadTestCase 上传测试用例
func UploadTestCase(problemID uint64, inputFile, outputFile *os.File) error {
	// 创建测试用例存储目录
	testCaseDir := filepath.Join("data", "testcases", fmt.Sprintf("problem_%d", problemID))
	if err := os.MkdirAll(testCaseDir, 0755); err != nil {
		return err
	}

	// 生成唯一的文件名
	timestamp := time.Now().UnixNano()
	InputFile := filepath.Join(testCaseDir, fmt.Sprintf("input_%d.txt", timestamp))
	OutputFile := filepath.Join(testCaseDir, fmt.Sprintf("output_%d.txt", timestamp))

	// 关闭原文件
	inputFile.Close()
	outputFile.Close()

	// 复制文件而不是移动
	if err := copyFile(inputFile.Name(), InputFile); err != nil {
		return fmt.Errorf("复制输入文件失败: %v", err)
	}
	if err := copyFile(outputFile.Name(), OutputFile); err != nil {
		// 如果输出文件复制失败，清理已复制的输入文件
		os.Remove(InputFile)
		return fmt.Errorf("复制输出文件失败: %v", err)
	}

	// 删除临时文件
	os.Remove(inputFile.Name())
	os.Remove(outputFile.Name())

	// 创建测试用例记录
	testCase := &models.TestCase{
		ProblemID:  problemID,
		InputFile:  InputFile,
		OutputFile: OutputFile,
		CreatedAt:  time.Now(),
	}

	return config.DB.Create(testCase).Error
}

// DeleteTestCase 删除测试用例
func DeleteTestCase(testCaseID uint64) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 先删除与该测试用例相关的判题结果
		if err := tx.Where("test_case_id = ?", testCaseID).
			Delete(&models.JudgeResult{}).Error; err != nil {
			return err
		}

		// 获取测试用例信息
		var testCase models.TestCase
		if err := tx.First(&testCase, testCaseID).Error; err != nil {
			return err
		}

		// 删除测试用例文件
		os.Remove(testCase.InputFile)
		os.Remove(testCase.OutputFile)

		// 删除测试用例记录
		return tx.Delete(&testCase).Error
	})
}

// GetTestCases 获取题目的测试用例列表
func GetTestCases(problemID uint64) ([]models.TestCaseWithLocalID, error) {
	var testCases []models.TestCaseWithLocalID

	// 使用 ROW_NUMBER() 窗口函数计算局部 ID
	err := config.DB.Raw(`
		SELECT 
			t.*,
			ROW_NUMBER() OVER (
				PARTITION BY t.problem_id 
				ORDER BY t.created_at ASC, t.id ASC
			) as local_id
		FROM test_cases t
		WHERE t.problem_id = ?
		ORDER BY t.created_at ASC, t.id ASC
	`, problemID).Scan(&testCases).Error

	return testCases, err
}

// copyFile 复制文件从源路径到目标路径
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// GetTestCaseContent 获取测试用例内容
func GetTestCaseContent(testCaseID uint64) (*models.TestCaseContentResponse, error) {
	// 获取测试用例基本信息
	var testCase models.TestCaseWithLocalID
	err := config.DB.Raw(`
		SELECT 
			t.*,
			ROW_NUMBER() OVER (
				PARTITION BY t.problem_id 
				ORDER BY t.created_at ASC, t.id ASC
			) as local_id
		FROM test_cases t
		WHERE t.id = ?
	`, testCaseID).Scan(&testCase).Error
	if err != nil {
		return nil, fmt.Errorf("获取测试用例信息失败: %v", err)
	}

	// 读取输入文件内容
	input, err := os.ReadFile(testCase.InputFile)
	if err != nil {
		return nil, fmt.Errorf("读取输入文件失败: %v", err)
	}

	// 读取输出文件内容
	output, err := os.ReadFile(testCase.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("读取输出文件失败: %v", err)
	}

	// 构造响应
	response := &models.TestCaseContentResponse{
		ID:        testCase.ID,
		LocalID:   testCase.LocalID,
		ProblemID: testCase.ProblemID,
		Input:     string(input),
		Output:    string(output),
	}

	return response, nil
}

// CreateTag 创建标签
func CreateTag(req *models.CreateTagRequest) (uint64, error) {
	// 检查标签名是否已存在
	var count int64
	if err := config.DB.Model(&models.ProblemTag{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("检查标签名是否存在失败: %v", err)
	}
	if count > 0 {
		return 0, fmt.Errorf("标签名 %s 已存在", req.Name)
	}

	// 验证颜色格式
	if !isValidHexColor(req.Color) {
		return 0, fmt.Errorf("无效的颜色格式，应为十六进制颜色值，如 #FF0000")
	}

	// 创建标签
	tag := &models.ProblemTag{
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(tag).Error; err != nil {
		return 0, fmt.Errorf("创建标签失败: %v", err)
	}

	return tag.ID, nil
}

// UpdateTag 更新标签
func UpdateTag(tagID uint64, req *models.UpdateTagRequest) error {
	// 检查标签是否存在
	var tag models.ProblemTag
	if err := config.DB.First(&tag, tagID).Error; err != nil {
		return fmt.Errorf("标签不存在: %v", err)
	}

	updates := make(map[string]interface{})

	// 更新标签名
	if req.Name != nil {
		// 检查新名称是否与其他标签重复
		var count int64
		if err := config.DB.Model(&models.ProblemTag{}).
			Where("name = ? AND id != ?", *req.Name, tagID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("检查标签名是否存在失败: %v", err)
		}
		if count > 0 {
			return fmt.Errorf("标签名 %s 已存在", *req.Name)
		}
		updates["name"] = *req.Name
	}

	// 更新标签颜色
	if req.Color != nil {
		if !isValidHexColor(*req.Color) {
			return fmt.Errorf("无效的颜色格式，应为十六进制颜色值，如 #FF0000")
		}
		updates["color"] = *req.Color
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&tag).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新标签失败: %v", err)
		}
	}

	return nil
}

// DeleteTag 删除标签
func DeleteTag(tagID uint64) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 检查标签是否存在
		var tag models.ProblemTag
		if err := tx.First(&tag, tagID).Error; err != nil {
			return fmt.Errorf("标签不存在: %v", err)
		}

		// 检查标签是否被使用
		var count int64
		if err := tx.Model(&models.ProblemTag{}).
			Joins("JOIN problem_tag_relations ON problem_tag_relations.tag_id = problem_tags.id").
			Where("problem_tags.id = ?", tagID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("检查标签使用情况失败: %v", err)
		}
		if count > 0 {
			return fmt.Errorf("标签正在被使用，无法删除")
		}

		// 删除标签
		if err := tx.Delete(&tag).Error; err != nil {
			return fmt.Errorf("删除标签失败: %v", err)
		}

		return nil
	})
}

// GetTagList 获取标签列表
func GetTagList(req *models.TagListRequest) (*models.TagListResponse, error) {
	query := config.DB.Model(&models.ProblemTag{})

	// 应用搜索条件
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取标签总数失败: %v", err)
	}

	// 获取分页数据
	var tags []models.ProblemTag
	if err := query.Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("获取标签列表失败: %v", err)
	}

	return &models.TagListResponse{
		Tags:     tags,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// isValidHexColor 验证十六进制颜色值格式
func isValidHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		if !((color[i] >= '0' && color[i] <= '9') ||
			(color[i] >= 'a' && color[i] <= 'f') ||
			(color[i] >= 'A' && color[i] <= 'F')) {
			return false
		}
	}
	return true
}

// SwitchDifficultySystem 切换难度等级系统
func SwitchDifficultySystem(req *models.SwitchDifficultySystemRequest) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 获取所有题目的当前难度等级系统和难度
		var problems []struct {
			ID               uint64                  `gorm:"column:id"`
			DifficultySystem models.DifficultySystem `gorm:"column:difficulty_system"`
			Difficulty       string                  `gorm:"column:difficulty"`
		}
		if err := tx.Model(&models.Problem{}).Select("id, difficulty_system, difficulty").Find(&problems).Error; err != nil {
			return fmt.Errorf("获取题目列表失败: %v", err)
		}

		// 批量更新题目的难度等级系统和难度
		for _, problem := range problems {
			// 获取映射后的难度等级
			newDifficulty := models.GetDifficultyMapping(
				problem.DifficultySystem,
				req.DifficultySystem,
				problem.Difficulty,
			)

			// 更新题目
			if err := tx.Model(&models.Problem{}).Where("id = ?", problem.ID).Updates(map[string]interface{}{
				"difficulty_system": req.DifficultySystem,
				"difficulty":        newDifficulty,
				"updated_at":        time.Now(),
			}).Error; err != nil {
				return fmt.Errorf("更新题目 %d 失败: %v", problem.ID, err)
			}
		}

		return nil
	})
}

// GetDifficultySystem 获取难度等级系统
func GetDifficultySystem() (*models.GetDifficultySystemResponse, error) {
	// 获取最新的一个题目的难度系统作为当前系统
	var currentSystem models.DifficultySystem
	if err := config.DB.Model(&models.Problem{}).
		Select("difficulty_system").
		Order("updated_at DESC").
		Limit(1).
		Scan(&currentSystem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有题目，默认使用普通难度系统
			currentSystem = models.DifficultySystemNormal
		} else {
			return nil, fmt.Errorf("获取当前难度系统失败: %v", err)
		}
	}

	// 如果没有获取到有效的难度系统，使用默认值
	if currentSystem == "" {
		currentSystem = models.DifficultySystemNormal
	}

	// 构造响应
	response := &models.GetDifficultySystemResponse{
		CurrentSystem: currentSystem,
		Systems:       models.OrderedDifficultySystems,
	}

	return response, nil
}
