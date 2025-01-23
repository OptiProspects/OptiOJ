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
		// 如果题目不是公开的，只有管理员和题目创建者可以访问
		if userID == 0 {
			return nil, errors.New("无权访问该题目")
		}
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

	// 获取用户提交状态
	if userID > 0 {
		type Result struct {
			Status *string `gorm:"column:status"`
		}
		var result Result
		err := config.DB.Raw(`
			SELECT 
				CASE 
					WHEN EXISTS (SELECT 1 FROM submissions WHERE problem_id = ? AND user_id = ? AND status = 'accepted') THEN 'accepted'
					WHEN EXISTS (SELECT 1 FROM submissions WHERE problem_id = ? AND user_id = ?) THEN 'attempted'
					ELSE NULL 
				END as status
		`, problemID, userID, problemID, userID).Scan(&result).Error
		if err != nil {
			return nil, fmt.Errorf("获取用户提交状态失败: %v", err)
		}
		detail.UserStatus = result.Status
	}

	return &detail, nil
}

// GetProblemList 获取题目列表
func GetProblemList(req *models.ProblemListRequest, userID uint) (*models.ProblemListResponse, error) {
	var problems []struct {
		models.Problem
		AcceptCount     int64   `gorm:"column:accept_count"`
		SubmissionCount int64   `gorm:"column:submission_count"`
		AcceptRate      float64 `gorm:"column:accept_rate"`
		UserStatus      *string `gorm:"column:user_status"`
	}

	// 构建基础查询
	query := config.DB.Table("problems").
		Select(`
			problems.*,
			COUNT(CASE WHEN submissions.status = 'accepted' THEN 1 END) as accept_count,
			COUNT(submissions.id) as submission_count,
			IFNULL(COUNT(CASE WHEN submissions.status = 'accepted' THEN 1 END) * 100.0 / 
				NULLIF(COUNT(submissions.id), 0), 0) as accept_rate,
			(SELECT 
				CASE 
					WHEN ? = 0 THEN NULL
					WHEN EXISTS (SELECT 1 FROM submissions s2 WHERE s2.problem_id = problems.id AND s2.user_id = ? AND s2.status = 'accepted') THEN 'accepted'
					WHEN EXISTS (SELECT 1 FROM submissions s2 WHERE s2.problem_id = problems.id AND s2.user_id = ?) THEN 'attempted'
					ELSE NULL 
				END
			) as user_status
		`, userID, userID, userID).
		Joins("LEFT JOIN submissions ON submissions.problem_id = problems.id").
		Group("problems.id")

	// 应用筛选条件
	if req.Title != "" {
		query = query.Where("problems.title LIKE ?", "%"+req.Title+"%")
	}
	if req.Difficulty != "" {
		query = query.Where("problems.difficulty = ?", req.Difficulty)
	}
	if len(req.Tags) > 0 {
		query = query.Joins("JOIN problem_tag_relations ptr ON ptr.problem_id = problems.id").
			Where("ptr.tag_id IN ?", req.Tags)
	}
	if len(req.Categories) > 0 {
		query = query.Joins("JOIN problem_category_relations pcr ON pcr.problem_id = problems.id").
			Where("pcr.category_id IN ?", req.Categories)
	}
	if req.IsPublic != nil {
		query = query.Where("problems.is_public = ?", *req.IsPublic)
	} else if userID == 0 {
		// 未登录用户只能看到公开题目
		query = query.Where("problems.is_public = ?", true)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("统计题目总数失败: %v", err)
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&problems).Error; err != nil {
		return nil, fmt.Errorf("查询题目列表失败: %v", err)
	}

	// 转换为响应格式
	var problemList []models.ProblemListItem
	for _, p := range problems {
		// 获取标签
		var tags []models.ProblemTag
		if err := config.DB.Table("problem_tags").
			Joins("JOIN problem_tag_relations ptr ON ptr.tag_id = problem_tags.id").
			Where("ptr.problem_id = ?", p.ID).
			Find(&tags).Error; err != nil {
			return nil, fmt.Errorf("获取题目标签失败: %v", err)
		}

		// 获取分类
		var categories []models.ProblemCategory
		if err := config.DB.Table("problem_categories").
			Joins("JOIN problem_category_relations pcr ON pcr.category_id = problem_categories.id").
			Where("pcr.problem_id = ?", p.ID).
			Find(&categories).Error; err != nil {
			return nil, fmt.Errorf("获取题目分类失败: %v", err)
		}

		item := models.ProblemListItem{
			ID:              p.ID,
			Title:           p.Title,
			Difficulty:      p.Difficulty,
			Tags:            tags,
			Categories:      categories,
			AcceptCount:     p.AcceptCount,
			SubmissionCount: p.SubmissionCount,
			AcceptRate:      p.AcceptRate,
			UserStatus:      p.UserStatus,
		}
		problemList = append(problemList, item)
	}

	return &models.ProblemListResponse{
		Problems:    problemList,
		TotalCount:  total,
		PageSize:    req.PageSize,
		CurrentPage: req.Page,
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

	// 如果指定了分类ID，验证是否为二级分类
	if req.CategoryID != nil {
		var category models.TagCategory
		if err := config.DB.First(&category, *req.CategoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return 0, fmt.Errorf("指定的分类不存在")
			}
			return 0, fmt.Errorf("查询分类信息失败: %v", err)
		}

		// 检查是否为二级分类
		if category.ParentID == nil {
			return 0, fmt.Errorf("标签只能添加到二级分类下")
		}

		// 检查父分类是否为一级分类
		var parentCategory models.TagCategory
		if err := config.DB.First(&parentCategory, *category.ParentID).Error; err != nil {
			return 0, fmt.Errorf("查询父分类信息失败: %v", err)
		}
		if parentCategory.ParentID != nil {
			return 0, fmt.Errorf("标签只能添加到二级分类下")
		}
	}

	// 创建标签
	tag := &models.ProblemTag{
		Name:       req.Name,
		Color:      req.Color,
		CategoryID: req.CategoryID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := config.DB.Create(tag).Error; err != nil {
		return 0, fmt.Errorf("创建标签失败: %v", err)
	}

	return tag.ID, nil
}

// UpdateTag 更新标签
func UpdateTag(tagID uint64, req *models.UpdateTagRequest) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 检查标签是否存在
		var tag models.ProblemTag
		if err := tx.First(&tag, tagID).Error; err != nil {
			return fmt.Errorf("标签不存在: %v", err)
		}

		updates := make(map[string]interface{})

		// 更新标签名
		if req.Name != nil {
			// 检查新名称是否与其他标签重复
			var count int64
			if err := tx.Model(&models.ProblemTag{}).
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

		// 更新分类
		if req.CategoryID != nil {
			// 验证是否为二级分类
			var category models.TagCategory
			if err := tx.First(&category, *req.CategoryID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("指定的分类不存在")
				}
				return fmt.Errorf("查询分类信息失败: %v", err)
			}

			// 检查是否为二级分类
			if category.ParentID == nil {
				return fmt.Errorf("标签只能添加到二级分类下")
			}

			// 检查父分类是否为一级分类
			var parentCategory models.TagCategory
			if err := tx.First(&parentCategory, *category.ParentID).Error; err != nil {
				return fmt.Errorf("查询父分类信息失败: %v", err)
			}
			if parentCategory.ParentID != nil {
				return fmt.Errorf("标签只能添加到二级分类下")
			}

			updates["category_id"] = *req.CategoryID
		}

		if len(updates) > 0 {
			updates["updated_at"] = time.Now()
			if err := tx.Model(&tag).Updates(updates).Error; err != nil {
				return fmt.Errorf("更新标签失败: %v", err)
			}
		}

		return nil
	})
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
	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取标签总数失败: %v", err)
	}

	// 获取分页数据
	var tags []models.TagWithCategory
	if err := query.Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("获取标签列表失败: %v", err)
	}

	// 获取标签的分类信息
	for i := range tags {
		if tags[i].CategoryID != nil {
			var category models.TagCategory
			if err := config.DB.First(&category, *tags[i].CategoryID).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return nil, fmt.Errorf("获取标签分类信息失败: %v", err)
				}
			} else {
				tags[i].Category = &category
				// 获取分类路径
				path, err := getCategoryPath(config.DB, category)
				if err != nil {
					return nil, fmt.Errorf("获取分类路径失败: %v", err)
				}
				tags[i].CategoryPath = path
			}
		}
	}

	// 获取一级分类列表
	var categories []models.TagCategory
	if err := config.DB.Where("parent_id IS NULL").Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %v", err)
	}

	return &models.TagListResponse{
		Tags:       tags,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Categories: categories,
	}, nil
}

// getCategoryPath 获取分类的完整路径
func getCategoryPath(db *gorm.DB, category models.TagCategory) ([]string, error) {
	path := []string{category.Name}
	currentID := category.ParentID

	for currentID != nil {
		var parent models.TagCategory
		if err := db.First(&parent, *currentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				break
			}
			return nil, err
		}
		path = append([]string{parent.Name}, path...)
		currentID = parent.ParentID
	}

	return path, nil
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

// CreateTagCategory 创建标签分类
func CreateTagCategory(req *models.CreateTagCategoryRequest) (uint64, error) {
	// 如果指定了父分类，检查父分类是否存在
	if req.ParentID != nil {
		var count int64
		if err := config.DB.Model(&models.TagCategory{}).Where("id = ?", *req.ParentID).Count(&count).Error; err != nil {
			return 0, fmt.Errorf("检查父分类是否存在失败: %v", err)
		}
		if count == 0 {
			return 0, fmt.Errorf("父分类不存在")
		}
	}

	// 检查同级分类下是否有重名
	var count int64
	query := config.DB.Model(&models.TagCategory{}).Where("name = ?", req.Name)
	if req.ParentID != nil {
		query = query.Where("parent_id = ?", *req.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("检查分类名是否存在失败: %v", err)
	}
	if count > 0 {
		return 0, fmt.Errorf("同级分类下已存在名为 %s 的分类", req.Name)
	}

	// 创建分类
	category := &models.TagCategory{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(category).Error; err != nil {
		return 0, fmt.Errorf("创建分类失败: %v", err)
	}

	return category.ID, nil
}

// UpdateTagCategory 更新标签分类
func UpdateTagCategory(categoryID uint64, req *models.UpdateTagCategoryRequest) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 检查分类是否存在
		var category models.TagCategory
		if err := tx.First(&category, categoryID).Error; err != nil {
			return fmt.Errorf("分类不存在: %v", err)
		}

		updates := make(map[string]interface{})

		// 更新分类名
		if req.Name != nil {
			// 检查新名称是否与同级其他分类重复
			var count int64
			query := tx.Model(&models.TagCategory{}).
				Where("name = ? AND id != ?", *req.Name, categoryID)
			if req.ParentID != nil {
				query = query.Where("parent_id = ?", *req.ParentID)
			} else {
				query = query.Where("parent_id IS NULL")
			}
			if err := query.Count(&count).Error; err != nil {
				return fmt.Errorf("检查分类名是否存在失败: %v", err)
			}
			if count > 0 {
				return fmt.Errorf("同级分类下已存在名为 %s 的分类", *req.Name)
			}
			updates["name"] = *req.Name
		}

		// 更新描述
		if req.Description != nil {
			updates["description"] = *req.Description
		}

		// 更新父分类
		if req.ParentID != nil {
			// 检查新父分类是否存在
			var count int64
			if err := tx.Model(&models.TagCategory{}).Where("id = ?", *req.ParentID).Count(&count).Error; err != nil {
				return fmt.Errorf("检查父分类是否存在失败: %v", err)
			}
			if count == 0 {
				return fmt.Errorf("父分类不存在")
			}

			// 检查是否会形成循环依赖
			if *req.ParentID == categoryID {
				return fmt.Errorf("不能将分类设置为自己的子分类")
			}

			// 检查新父分类是否是当前分类的子分类
			var childIDs []uint64
			if err := getChildCategoryIDs(tx, categoryID, &childIDs); err != nil {
				return fmt.Errorf("检查子分类失败: %v", err)
			}
			for _, childID := range childIDs {
				if childID == *req.ParentID {
					return fmt.Errorf("不能将分类设置为其子分类的子分类")
				}
			}

			updates["parent_id"] = req.ParentID
		}

		if len(updates) > 0 {
			updates["updated_at"] = time.Now()
			if err := tx.Model(&category).Updates(updates).Error; err != nil {
				return fmt.Errorf("更新分类失败: %v", err)
			}
		}

		return nil
	})
}

// DeleteTagCategory 删除标签分类
func DeleteTagCategory(categoryID uint64) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 检查分类是否存在
		var category models.TagCategory
		if err := tx.First(&category, categoryID).Error; err != nil {
			return fmt.Errorf("分类不存在: %v", err)
		}

		// 检查是否有子分类
		var childCount int64
		if err := tx.Model(&models.TagCategory{}).Where("parent_id = ?", categoryID).Count(&childCount).Error; err != nil {
			return fmt.Errorf("检查子分类失败: %v", err)
		}
		if childCount > 0 {
			return fmt.Errorf("该分类下还有子分类，无法删除")
		}

		// 检查是否有标签使用该分类
		var tagCount int64
		if err := tx.Model(&models.ProblemTag{}).Where("category_id = ?", categoryID).Count(&tagCount).Error; err != nil {
			return fmt.Errorf("检查标签使用情况失败: %v", err)
		}
		if tagCount > 0 {
			return fmt.Errorf("该分类下还有标签，无法删除")
		}

		// 删除分类
		if err := tx.Delete(&category).Error; err != nil {
			return fmt.Errorf("删除分类失败: %v", err)
		}

		return nil
	})
}

// GetTagCategoryList 获取标签分类列表
func GetTagCategoryList(req *models.GetTagCategoryListRequest) (*models.GetTagCategoryListResponse, error) {
	var categories []models.TagCategory

	// 构建查询
	query := config.DB.Model(&models.TagCategory{})
	if req.ParentID != nil {
		query = query.Where("parent_id = ?", *req.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	// 执行查询
	if err := query.Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %v", err)
	}

	// 转换为响应格式
	var categoryDetails []models.TagCategoryDetail
	for _, category := range categories {
		detail := models.TagCategoryDetail{
			TagCategory: category,
		}

		// 获取子分类
		if err := getChildCategories(config.DB, category.ID, &detail.Children); err != nil {
			return nil, fmt.Errorf("获取子分类失败: %v", err)
		}

		categoryDetails = append(categoryDetails, detail)
	}

	return &models.GetTagCategoryListResponse{
		Categories: categoryDetails,
	}, nil
}

// getChildCategories 递归获取子分类
func getChildCategories(db *gorm.DB, parentID uint64, children *[]models.TagCategoryDetail) error {
	var categories []models.TagCategory
	if err := db.Where("parent_id = ?", parentID).Order("created_at DESC").Find(&categories).Error; err != nil {
		return err
	}

	for _, category := range categories {
		detail := models.TagCategoryDetail{
			TagCategory: category,
		}
		if err := getChildCategories(db, category.ID, &detail.Children); err != nil {
			return err
		}
		*children = append(*children, detail)
	}

	return nil
}

// getChildCategoryIDs 递归获取所有子分类ID
func getChildCategoryIDs(db *gorm.DB, parentID uint64, childIDs *[]uint64) error {
	var categories []models.TagCategory
	if err := db.Where("parent_id = ?", parentID).Find(&categories).Error; err != nil {
		return err
	}

	for _, category := range categories {
		*childIDs = append(*childIDs, category.ID)
		if err := getChildCategoryIDs(db, category.ID, childIDs); err != nil {
			return err
		}
	}

	return nil
}

// GetTagCategoryTree 获取标签分类树形结构
func GetTagCategoryTree() (*models.GetTagCategoryListResponse, error) {
	// 获取所有一级分类
	req := &models.GetTagCategoryListRequest{
		ParentID: nil,
	}
	return GetTagCategoryList(req)
}
