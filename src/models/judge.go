package models

import "time"

// 提交状态常量
const (
	StatusPending           = "pending"             // 等待判题
	StatusJudging           = "judging"             // 判题中
	StatusAccepted          = "accepted"            // 通过
	StatusWrongAnswer       = "wrong_answer"        // 答案错误
	StatusTimeLimitExceed   = "time_limit_exceed"   // 超时
	StatusMemoryLimitExceed = "memory_limit_exceed" // 内存超限
	StatusRuntimeError      = "runtime_error"       // 运行时错误
	StatusCompileError      = "compile_error"       // 编译错误
	StatusSystemError       = "system_error"        // 系统错误
)

// 支持的编程语言
const (
	LangC      = "c"
	LangCPP    = "cpp"
	LangJava   = "java"
	LangPython = "python"
	LangGo     = "go"
)

// Submission 提交记录
type Submission struct {
	ID           uint64    `json:"id"`
	ProblemID    uint64    `json:"problem_id"`
	UserID       uint64    `json:"user_id"`
	Language     string    `json:"language"`
	Code         string    `json:"code"`
	Status       string    `json:"status"`
	TimeUsed     *int      `json:"time_used"`
	MemoryUsed   *int      `json:"memory_used"`
	ErrorMessage *string   `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// JudgeResult 判题结果
type JudgeResult struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	SubmissionID uint64    `json:"submission_id"`
	TestCaseID   uint64    `json:"test_case_id"`
	Status       string    `json:"status"`
	TimeUsed     int       `json:"time_used"`
	MemoryUsed   int       `json:"memory_used"`
	ErrorMessage *string   `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

// SubmissionRequest 提交代码请求
type SubmissionRequest struct {
	ProblemID uint64 `json:"problem_id" binding:"required"`
	Language  string `json:"language" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

// SubmissionListRequest 提交记录列表请求
type SubmissionListRequest struct {
	Page      int     `form:"page" binding:"required,min=1"`
	PageSize  int     `form:"page_size" binding:"required,min=1,max=100"`
	ProblemID *uint64 `form:"problem_id"`
	UserID    *uint64 `form:"user_id"`
	Language  string  `form:"language"`
	Status    string  `form:"status"`
}

// SubmissionListResponse 提交记录列表响应
type SubmissionListResponse struct {
	Submissions []SubmissionDetail `json:"submissions"`
	Total       int64              `json:"total"`
	Page        int                `json:"page"`
	PageSize    int                `json:"page_size"`
}

// SubmissionDetail 提交记录详情
type SubmissionDetail struct {
	Submission
	Problem *Problem      `json:"problem"`
	User    *User         `json:"user"`
	Results []JudgeResult `json:"results,omitempty" gorm:"foreignKey:SubmissionID"`
}

// JudgeConfig 判题配置
type JudgeConfig struct {
	TimeLimit   int    `json:"time_limit"`   // 时间限制（毫秒）
	MemoryLimit int    `json:"memory_limit"` // 内存限制（MB）
	Language    string `json:"language"`     // 编程语言
	Code        string `json:"code"`         // 源代码
	TestCase    struct {
		Input  string `json:"input"`  // 输入文件路径
		Output string `json:"output"` // 输出文件路径
	} `json:"test_case"`
}

// CompileResult 编译结果
type CompileResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// TestCaseResult 单个测试点的运行结果
type TestCaseResult struct {
	Status       string  `json:"status"`
	TimeUsed     float64 `json:"time_used"`   // 单位：毫秒
	MemoryUsed   float64 `json:"memory_used"` // 单位：KB
	ActualOutput string  `json:"actual_output"`
	TestCaseID   int     `json:"test_case_id"`
}

// RunResult 运行结果
type RunResult struct {
	Status          string           `json:"status"`
	TimeUsed        int              `json:"time_used"`
	MemoryUsed      int              `json:"memory_used"`
	ErrorMessage    string           `json:"error_message,omitempty"`
	Output          string           `json:"output,omitempty"`
	TestCaseResults []TestCaseResult `json:"test_case_results,omitempty"`
}
