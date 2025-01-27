package models

import "time"

// TeamProblem 团队私有题目
type TeamProblem struct {
	ID                uint64    `json:"id"`
	TeamID            uint64    `json:"team_id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	InputDescription  string    `json:"input_description"`
	OutputDescription string    `json:"output_description"`
	SampleCases       string    `json:"sample_cases"`
	Hint              string    `json:"hint"`
	TimeLimit         int       `json:"time_limit"`
	MemoryLimit       int       `json:"memory_limit"`
	CreatedBy         uint64    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TeamProblemTestCase 团队私有题目测试用例
type TeamProblemTestCase struct {
	ID         uint64    `json:"id"`
	ProblemID  uint64    `json:"problem_id"`
	InputData  string    `json:"input_data"`
	OutputData string    `json:"output_data"`
	IsSample   bool      `json:"is_sample"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateTeamProblemRequest 创建团队私有题目请求
type CreateTeamProblemRequest struct {
	TeamID            uint64 `json:"team_id" binding:"required"`
	Title             string `json:"title" binding:"required"`
	Description       string `json:"description" binding:"required"`
	InputDescription  string `json:"input_description"`
	OutputDescription string `json:"output_description"`
	SampleCases       string `json:"sample_cases"`
	Hint              string `json:"hint"`
	TimeLimit         int    `json:"time_limit"`
	MemoryLimit       int    `json:"memory_limit"`
	TestCases         []struct {
		InputData  string `json:"input_data" binding:"required"`
		OutputData string `json:"output_data" binding:"required"`
		IsSample   bool   `json:"is_sample"`
	} `json:"test_cases" binding:"required,min=1"`
}

// UpdateTeamProblemRequest 更新团队私有题目请求
type UpdateTeamProblemRequest struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	InputDescription  string `json:"input_description"`
	OutputDescription string `json:"output_description"`
	SampleCases       string `json:"sample_cases"`
	Hint              string `json:"hint"`
	TimeLimit         int    `json:"time_limit"`
	MemoryLimit       int    `json:"memory_limit"`
	TestCases         []struct {
		InputData  string `json:"input_data"`
		OutputData string `json:"output_data"`
		IsSample   bool   `json:"is_sample"`
	} `json:"test_cases"`
}

// TeamProblemDetail 团队私有题目详情
type TeamProblemDetail struct {
	TeamProblem
	TestCases []TeamProblemTestCase `json:"test_cases,omitempty"`
}

// TeamProblemListRequest 团队私有题目列表请求
type TeamProblemListRequest struct {
	TeamID   uint64 `form:"team_id" binding:"required"`
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// TeamProblemListResponse 团队私有题目列表响应
type TeamProblemListResponse struct {
	Problems []TeamProblem `json:"problems"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}
