package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// CreateProblem 创建题目
func CreateProblem(req *models.CreateProblemRequest, createdBy uint64) (uint64, error) {
	problem := &models.Problem{
		Title:             req.Title,
		Description:       req.Description,
		InputDescription:  req.InputDescription,
		OutputDescription: req.OutputDescription,
		Samples:           req.Samples,
		Hint:              req.Hint,
		Source:            req.Source,
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
			updates["samples"] = *req.Samples
		}
		if req.Hint != nil {
			updates["hint"] = *req.Hint
		}
		if req.Source != nil {
			updates["source"] = *req.Source
		}
		if req.Difficulty != nil {
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

		if err := tx.Model(&models.Problem{}).Where("id = ?", problemID).Updates(updates).Error; err != nil {
			return err
		}

		// 更新分类关系
		if req.CategoryIDs != nil {
			if err := tx.Where("problem_id = ?", problemID).Delete(&models.ProblemCategory{}).Error; err != nil {
				return err
			}
			if len(req.CategoryIDs) > 0 {
				var values []map[string]interface{}
				for _, categoryID := range req.CategoryIDs {
					values = append(values, map[string]interface{}{
						"problem_id":  problemID,
						"category_id": categoryID,
					})
				}
				if err := tx.Table("problem_category_relations").Create(values).Error; err != nil {
					return err
				}
			}
		}

		// 更新标签关系
		if req.TagIDs != nil {
			if err := tx.Where("problem_id = ?", problemID).Delete(&models.ProblemTag{}).Error; err != nil {
				return err
			}
			if len(req.TagIDs) > 0 {
				var values []map[string]interface{}
				for _, tagID := range req.TagIDs {
					values = append(values, map[string]interface{}{
						"problem_id": problemID,
						"tag_id":     tagID,
					})
				}
				if err := tx.Table("problem_tag_relations").Create(values).Error; err != nil {
					return err
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
func UploadTestCase(problemID uint64, inputFile, outputFile *os.File, isSample bool) error {
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
		IsSample:   isSample,
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
func GetTestCases(problemID uint64) ([]models.TestCase, error) {
	var testCases []models.TestCase

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
