package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CreateTeamProblem 创建团队私有题目
func CreateTeamProblem(req *models.CreateTeamProblemRequest, userID uint64) (uint64, error) {
	// 检查用户权限
	role, err := GetTeamUserRole(req.TeamID, userID)
	if err != nil {
		return 0, err
	}
	if role != "owner" && role != "admin" {
		return 0, errors.New("权限不足")
	}

	problem := &models.TeamProblem{
		TeamID:            req.TeamID,
		Title:             req.Title,
		Description:       req.Description,
		InputDescription:  req.InputDescription,
		OutputDescription: req.OutputDescription,
		SampleCases:       req.SampleCases,
		Hint:              req.Hint,
		TimeLimit:         req.TimeLimit,
		MemoryLimit:       req.MemoryLimit,
		CreatedBy:         userID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return problem.ID, config.DB.Transaction(func(tx *gorm.DB) error {
		// 创建题目
		if err := tx.Create(problem).Error; err != nil {
			return err
		}

		// 添加测试用例
		for _, tc := range req.TestCases {
			testCase := &models.TeamProblemTestCase{
				ProblemID:  problem.ID,
				InputData:  tc.InputData,
				OutputData: tc.OutputData,
				IsSample:   tc.IsSample,
				CreatedAt:  time.Now(),
			}
			if err := tx.Create(testCase).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateTeamProblem 更新团队私有题目
func UpdateTeamProblem(problemID uint64, req *models.UpdateTeamProblemRequest, userID uint64) error {
	var problem models.TeamProblem
	if err := config.DB.First(&problem, problemID).Error; err != nil {
		return err
	}

	// 检查用户权限
	role, err := GetTeamUserRole(problem.TeamID, userID)
	if err != nil {
		return err
	}
	if role != "owner" && role != "admin" && problem.CreatedBy != userID {
		return errors.New("权限不足")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 更新题目信息
		updates := make(map[string]interface{})
		if req.Title != "" {
			updates["title"] = req.Title
		}
		if req.Description != "" {
			updates["description"] = req.Description
		}
		if req.InputDescription != "" {
			updates["input_description"] = req.InputDescription
		}
		if req.OutputDescription != "" {
			updates["output_description"] = req.OutputDescription
		}
		if req.SampleCases != "" {
			updates["sample_cases"] = req.SampleCases
		}
		if req.Hint != "" {
			updates["hint"] = req.Hint
		}
		if req.TimeLimit > 0 {
			updates["time_limit"] = req.TimeLimit
		}
		if req.MemoryLimit > 0 {
			updates["memory_limit"] = req.MemoryLimit
		}
		updates["updated_at"] = time.Now()

		if err := tx.Model(&problem).Updates(updates).Error; err != nil {
			return err
		}

		// 如果提供了新的测试用例，更新测试用例
		if req.TestCases != nil {
			// 删除旧的测试用例
			if err := tx.Where("problem_id = ?", problemID).Delete(&models.TeamProblemTestCase{}).Error; err != nil {
				return err
			}

			// 添加新的测试用例
			for _, tc := range req.TestCases {
				testCase := &models.TeamProblemTestCase{
					ProblemID:  problemID,
					InputData:  tc.InputData,
					OutputData: tc.OutputData,
					IsSample:   tc.IsSample,
					CreatedAt:  time.Now(),
				}
				if err := tx.Create(testCase).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// DeleteTeamProblem 删除团队私有题目
func DeleteTeamProblem(problemID uint64, userID uint64) error {
	var problem models.TeamProblem
	if err := config.DB.First(&problem, problemID).Error; err != nil {
		return err
	}

	// 检查用户权限
	role, err := GetTeamUserRole(problem.TeamID, userID)
	if err != nil {
		return err
	}
	if role != "owner" && role != "admin" && problem.CreatedBy != userID {
		return errors.New("权限不足")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 删除测试用例
		if err := tx.Where("problem_id = ?", problemID).Delete(&models.TeamProblemTestCase{}).Error; err != nil {
			return err
		}

		// 删除题目
		return tx.Delete(&problem).Error
	})
}

// GetTeamProblemDetail 获取团队私有题目详情
func GetTeamProblemDetail(problemID uint64, userID uint64) (*models.TeamProblemDetail, error) {
	var problem models.TeamProblem
	if err := config.DB.First(&problem, problemID).Error; err != nil {
		return nil, err
	}

	// 检查用户权限
	role, err := GetTeamUserRole(problem.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, errors.New("您不是团队成员")
	}

	detail := &models.TeamProblemDetail{
		TeamProblem: problem,
	}

	// 获取测试用例
	if role == "owner" || role == "admin" || problem.CreatedBy == userID {
		var testCases []models.TeamProblemTestCase
		if err := config.DB.Where("problem_id = ?", problemID).Find(&testCases).Error; err != nil {
			return nil, err
		}
		detail.TestCases = testCases
	}

	return detail, nil
}

// GetTeamProblemList 获取团队私有题目列表
func GetTeamProblemList(req *models.TeamProblemListRequest, userID uint64) (*models.TeamProblemListResponse, error) {
	// 检查用户权限
	role, err := GetTeamUserRole(req.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, errors.New("您不是团队成员")
	}

	query := config.DB.Model(&models.TeamProblem{}).Where("team_id = ?", req.TeamID)

	// 添加关键字搜索
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	var problems []models.TeamProblem
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).
		Order("created_at DESC").
		Find(&problems).Error; err != nil {
		return nil, err
	}

	return &models.TeamProblemListResponse{
		Problems: problems,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
