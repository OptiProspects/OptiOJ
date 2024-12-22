package models

import "time"

// Problem 题目模型
type Problem struct {
	ID                uint64    `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	InputDescription  string    `json:"input_description"`
	OutputDescription string    `json:"output_description"`
	Samples           string    `json:"samples"` // JSON格式存储样例输入输出
	Hint              string    `json:"hint"`
	Source            string    `json:"source"`
	Difficulty        string    `json:"difficulty"`
	TimeLimit         int       `json:"time_limit"`
	MemoryLimit       int       `json:"memory_limit"`
	IsPublic          bool      `json:"is_public"`
	CreatedBy         uint64    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ProblemCategory 题目分类
type ProblemCategory struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *uint64   `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// ProblemTag 题目标签
type ProblemTag struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

// TestCase 测试用例
type TestCase struct {
	ID         uint64    `json:"id"`
	ProblemID  uint64    `json:"problem_id"`
	LocalID    int       `json:"local_id"`    // 测试用例在当前题目中的序号
	InputFile  string    `json:"input_file"`  // 输入文件路径
	OutputFile string    `json:"output_file"` // 输出文件路径
	IsSample   bool      `json:"is_sample"`
	CreatedAt  time.Time `json:"created_at"`
}

// ProblemDetail 题目详细信息（包含分类和标签）
type ProblemDetail struct {
	Problem
	Categories []ProblemCategory `json:"categories"`
	Tags       []ProblemTag      `json:"tags"`
}

// CreateProblemRequest 创建题目请求
type CreateProblemRequest struct {
	Title             string   `json:"title" binding:"required"`
	Description       string   `json:"description" binding:"required"`
	InputDescription  string   `json:"input_description"`
	OutputDescription string   `json:"output_description"`
	Samples           string   `json:"samples"`
	Hint              string   `json:"hint"`
	Source            string   `json:"source"`
	Difficulty        string   `json:"difficulty" binding:"required,oneof=easy medium hard"`
	TimeLimit         int      `json:"time_limit" binding:"required,min=100,max=10000"`
	MemoryLimit       int      `json:"memory_limit" binding:"required,min=16,max=1024"`
	IsPublic          bool     `json:"is_public"`
	CategoryIDs       []uint64 `json:"category_ids"`
	TagIDs            []uint64 `json:"tag_ids"`
}

// UpdateProblemRequest 更新题目请求
type UpdateProblemRequest struct {
	Title             *string  `json:"title"`
	Description       *string  `json:"description"`
	InputDescription  *string  `json:"input_description"`
	OutputDescription *string  `json:"output_description"`
	Samples           *string  `json:"samples"`
	Hint              *string  `json:"hint"`
	Source            *string  `json:"source"`
	Difficulty        *string  `json:"difficulty" binding:"omitempty,oneof=easy medium hard"`
	TimeLimit         *int     `json:"time_limit" binding:"omitempty,min=100,max=10000"`
	MemoryLimit       *int     `json:"memory_limit" binding:"omitempty,min=16,max=1024"`
	IsPublic          *bool    `json:"is_public"`
	CategoryIDs       []uint64 `json:"category_ids"`
	TagIDs            []uint64 `json:"tag_ids"`
}

// ProblemListRequest 题目列表请求
type ProblemListRequest struct {
	Page       int      `form:"page" binding:"required,min=1"`
	PageSize   int      `form:"page_size" binding:"required,min=1,max=100"`
	Title      string   `form:"title"`
	Difficulty string   `form:"difficulty" binding:"omitempty,oneof=easy medium hard"`
	CategoryID *uint64  `form:"category_id"`
	TagIDs     []uint64 `form:"tag_ids"`
	IsPublic   *bool    `form:"is_public"`
}

// ProblemListResponse 题目列表响应
type ProblemListResponse struct {
	Problems []ProblemDetail `json:"problems"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// TestCaseUploadRequest 测试用例上传请求
type TestCaseUploadRequest struct {
	ProblemID uint64 `form:"problem_id" binding:"required"`
	IsSample  bool   `form:"is_sample"`
}
