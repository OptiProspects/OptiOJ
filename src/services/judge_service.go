package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateSubmission 创建提交记录并执行判题
func CreateSubmission(req *models.SubmissionRequest, userID uint64) (uint64, error) {
	// 获取题目信息
	var problem models.Problem
	if err := config.DB.First(&problem, req.ProblemID).Error; err != nil {
		return 0, fmt.Errorf("获取题目信息失败: %v", err)
	}

	// 创建提交记录
	submission := &models.Submission{
		ProblemID: req.ProblemID,
		UserID:    userID,
		Language:  req.Language,
		Code:      req.Code,
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	if err := config.DB.Create(submission).Error; err != nil {
		return 0, fmt.Errorf("创建提交记录失败: %v", err)
	}

	// 异步执行判题
	go func() {
		// 更新状态为判题中
		if err := config.DB.Model(submission).Updates(map[string]interface{}{
			"status":     models.StatusJudging,
			"updated_at": time.Now(),
		}).Error; err != nil {
			return
		}

		// 判题配置
		judgeConfig := &models.JudgeConfig{
			TimeLimit:   problem.TimeLimit,
			MemoryLimit: problem.MemoryLimit,
			Language:    req.Language,
			Code:        req.Code,
		}

		// 获取测试用例
		var testCases []models.TestCase
		if err := config.DB.Where("problem_id = ?", problem.ID).Find(&testCases).Error; err != nil {
			return
		}

		// 调用 gRPC 判题服务
		result, err := GetJudgeClient().Submit(judgeConfig, testCases)
		if err != nil {
			config.DB.Model(submission).Updates(map[string]interface{}{
				"status":        models.StatusSystemError,
				"error_message": err.Error(),
				"updated_at":    time.Now(),
			})
			return
		}

		// 开启事务保存结果
		err = config.DB.Transaction(func(tx *gorm.DB) error {
			// 保存每个测试点的结果
			for i, testResult := range result.TestCaseResults {
				judgeResult := &models.JudgeResult{
					SubmissionID: submission.ID,
					TestCaseID:   testCases[i].ID,
					Status:       testResult.Status,
					TimeUsed:     int(testResult.TimeUsed),
					MemoryUsed:   int(testResult.MemoryUsed),
					CreatedAt:    time.Now(),
				}
				if err := tx.Create(judgeResult).Error; err != nil {
					return err
				}
			}

			// 更新提交记录的最终状态
			updates := map[string]interface{}{
				"status":      result.Status,
				"time_used":   result.TimeUsed,
				"memory_used": result.MemoryUsed,
				"updated_at":  time.Now(),
			}
			if result.ErrorMessage != "" {
				updates["error_message"] = result.ErrorMessage
			}

			return tx.Model(submission).Updates(updates).Error
		})

		if err != nil {
			logrus.Errorf("保存判题结果失败: %v", err)
			config.DB.Model(submission).Updates(map[string]interface{}{
				"status":        models.StatusSystemError,
				"error_message": fmt.Sprintf("保存判题结果失败: %v", err),
				"updated_at":    time.Now(),
			})
		}
	}()

	return submission.ID, nil
}

// convertStatus 将 gRPC 状态码转换为系统状态
func convertStatus(status int32) string {
	switch status {
	case 0:
		return models.StatusAccepted
	case 1:
		return models.StatusWrongAnswer
	case 2:
		return models.StatusTimeLimitExceed
	case 3:
		return models.StatusMemoryLimitExceed
	case 4:
		return models.StatusRuntimeError
	default:
		return models.StatusSystemError
	}
}

// GetSubmissionList 获取提交记录列表
func GetSubmissionList(req *models.SubmissionListRequest) (*models.SubmissionListResponse, error) {
	var submissions []models.SubmissionDetail
	var total int64

	// 构建查询
	query := config.DB.Model(&models.Submission{})

	// 应用��选条件
	if req.ProblemID != nil {
		query = query.Where("problem_id = ?", *req.ProblemID)
	}
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Language != "" {
		query = query.Where("language = ?", req.Language)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取提交记录总数失败: %v", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).
		Preload("Problem").
		Preload("User").
		Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("获取提交记录列表失败: %v", err)
	}

	return &models.SubmissionListResponse{
		Submissions: submissions,
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}, nil
}

// GetSubmissionDetail 获取提交记录详情
func GetSubmissionDetail(submissionID uint64) (*models.SubmissionDetail, error) {
	var detail models.SubmissionDetail

	if err := config.DB.Model(&models.Submission{}).
		Preload("Problem").
		Preload("User").
		Preload("Results").
		First(&detail, submissionID).Error; err != nil {
		return nil, fmt.Errorf("获��提交记录详情失败: %v", err)
	}

	return &detail, nil
}

func judge(submission *models.Submission) error {
	// 更新状态为判题中
	if err := config.DB.Model(submission).Updates(map[string]interface{}{
		"status":     models.StatusJudging,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("更新判题状态失败: %v", err)
	}

	// 获取题目信息
	var problem models.Problem
	if err := config.DB.First(&problem, submission.ProblemID).Error; err != nil {
		return fmt.Errorf("获取题目信息失败: %v", err)
	}

	// 获取测试用例
	var testCases []models.TestCase
	if err := config.DB.Where("problem_id = ?", problem.ID).Find(&testCases).Error; err != nil {
		return fmt.Errorf("获取测试用例失败: %v", err)
	}

	// 判题配置
	judgeConfig := &models.JudgeConfig{
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		Language:    submission.Language,
		Code:        submission.Code,
	}

	// 调用 gRPC 判题服务
	result, err := GetJudgeClient().Submit(judgeConfig, testCases)
	if err != nil {
		return fmt.Errorf("判题失败: %v", err)
	}

	// 更新提交记录的最终状态
	updates := map[string]interface{}{
		"status":      result.Status,
		"time_used":   result.TimeUsed,
		"memory_used": result.MemoryUsed,
		"updated_at":  time.Now(),
	}
	if result.ErrorMessage != "" {
		updates["error_message"] = result.ErrorMessage
	}

	if err := config.DB.Model(submission).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新提交记录失败: %v", err)
	}

	return nil
}

// Debug 在线调试代码
func Debug(req *models.DebugRequest) (*models.DebugResponse, error) {
	// 创建临时文件存储输入数据
	inputFile, err := os.CreateTemp("", "debug_input_*.txt")
	if err != nil {
		return nil, fmt.Errorf("创建输入文件失败: %v", err)
	}
	defer os.Remove(inputFile.Name())
	defer inputFile.Close()

	// 写入输入数据
	if _, err := inputFile.WriteString(req.Input); err != nil {
		return nil, fmt.Errorf("写入输入数据失败: %v", err)
	}

	// 创建临时文件存储期望输出
	outputFile, err := os.CreateTemp("", "debug_output_*.txt")
	if err != nil {
		return nil, fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer os.Remove(outputFile.Name())
	defer outputFile.Close()

	// 写入期望输出数据
	if _, err := outputFile.WriteString(req.ExpectedOutput); err != nil {
		return nil, fmt.Errorf("写入期望输出数据失败: %v", err)
	}

	// 构造测试用例
	testCase := models.TestCase{
		InputFile:  inputFile.Name(),
		OutputFile: outputFile.Name(),
	}

	// 构造判题配置
	judgeConfig := &models.JudgeConfig{
		TimeLimit:   req.TimeLimit,
		MemoryLimit: req.MemoryLimit,
		Language:    req.Language,
		Code:        req.Code,
	}

	// 调用判题服务
	result, err := GetJudgeClient().Submit(judgeConfig, []models.TestCase{testCase})
	if err != nil {
		return nil, fmt.Errorf("调用判题服务失败: %v", err)
	}

	// 构造响应
	response := &models.DebugResponse{
		Status:         result.Status,
		TimeUsed:       float64(result.TimeUsed),
		MemoryUsed:     float64(result.MemoryUsed),
		ErrorMessage:   result.ErrorMessage,
		ExpectedOutput: req.ExpectedOutput,
	}

	// 如果运行成功，获取输出并比对结果
	if len(result.TestCaseResults) > 0 {
		response.Output = result.TestCaseResults[0].ActualOutput
		// 如果有预期输出，则进行比对
		if req.ExpectedOutput != "" {
			response.IsCorrect = (result.Status == models.StatusAccepted)
		}
	}

	return response, nil
}
