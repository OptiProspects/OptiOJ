package models

import "time"

// Team 团队
type Team struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Avatar      string    `json:"avatar,omitempty"`
	CreatedBy   uint64    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TeamMember 团队成员
type TeamMember struct {
	TeamID   uint64    `json:"team_id"`
	UserID   uint64    `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

// TeamAssignment 团队作业
type TeamAssignment struct {
	ID          uint64    `json:"id"`
	TeamID      uint64    `json:"team_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedBy   uint64    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TeamAssignmentProblem 团队作业题目
type TeamAssignmentProblem struct {
	AssignmentID  uint64  `json:"assignment_id"`
	ProblemID     uint64  `json:"problem_id"`
	ProblemType   string  `json:"problem_type"`    // 题目类型：global-全局题目，team-团队题目
	TeamProblemID *uint64 `json:"team_problem_id"` // 团队题目ID，当problem_type为team时有效
	OrderIndex    int     `json:"order_index"`
	Score         int     `json:"score"`
}

// TeamProblemList 团队题单
type TeamProblemList struct {
	ID          uint64    `json:"id"`
	TeamID      uint64    `json:"team_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedBy   uint64    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TeamProblemListItem 团队题单项目
type TeamProblemListItem struct {
	ListID     uint64 `json:"list_id"`
	ProblemID  uint64 `json:"problem_id"`
	OrderIndex int    `json:"order_index"`
	Note       string `json:"note"`
}

// TeamInvitation 团队邀请
type TeamInvitation struct {
	ID        uint64    `json:"id"`
	TeamID    uint64    `json:"team_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedBy uint64    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateTeamRequest 创建团队请求
type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateTeamRequest 更新团队请求
type UpdateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateAssignmentRequest 创建作业请求
type CreateAssignmentRequest struct {
	TeamID      uint64    `json:"team_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Problems    []struct {
		ProblemID     uint64 `json:"problem_id" binding:"required"`
		ProblemType   string `json:"problem_type" binding:"required,oneof=global team"` // 题目类型：global-全局题目，team-团队题目
		TeamProblemID uint64 `json:"team_problem_id,omitempty"`                         // 团队题目ID，当problem_type为team时必填
		OrderIndex    int    `json:"order_index"`
		Score         int    `json:"score"`
	} `json:"problems" binding:"required"`
}

// UpdateAssignmentRequest 更新作业请求
type UpdateAssignmentRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Problems    []struct {
		ProblemID     uint64 `json:"problem_id"`
		ProblemType   string `json:"problem_type" binding:"omitempty,oneof=global team"` // 题目类型：global-全局题目，team-团队题目
		TeamProblemID uint64 `json:"team_problem_id,omitempty"`                          // 团队题目ID，当problem_type为team时必填
		OrderIndex    int    `json:"order_index"`
		Score         int    `json:"score"`
	} `json:"problems"`
}

// CreateProblemListRequest 创建题单请求
type CreateProblemListRequest struct {
	TeamID      uint64 `json:"team_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	Problems    []struct {
		ProblemID  uint64 `json:"problem_id" binding:"required"`
		OrderIndex int    `json:"order_index"`
		Note       string `json:"note"`
	} `json:"problems" binding:"required"`
}

// UpdateProblemListRequest 更新题单请求
type UpdateProblemListRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	Problems    []struct {
		ProblemID  uint64 `json:"problem_id"`
		OrderIndex int    `json:"order_index"`
		Note       string `json:"note"`
	} `json:"problems"`
}

// TeamListRequest 团队列表请求
type TeamListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Keyword  string `form:"keyword"`
	Scope    string `form:"scope" binding:"omitempty,oneof=all joined"` // 查询范围：all-所有团队，joined-已加入的团队
}

// TeamListResponse 团队列表响应
type TeamListResponse struct {
	Teams    []TeamDetail `json:"teams"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

// TeamOwnerInfo 团队创建者信息
type TeamOwnerInfo struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname,omitempty"` // 团队内名称
}

// TeamDetail 团队详情
type TeamDetail struct {
	Team
	MemberCount int            `json:"member_count"`
	UserRole    string         `json:"user_role,omitempty"`
	IsJoined    bool           `json:"is_joined"`
	Owner       *TeamOwnerInfo `json:"owner,omitempty"`
}

// TeamAvatar 团队头像
type TeamAvatar struct {
	ID         uint64    `json:"id"`
	TeamID     uint64    `json:"team_id"`
	Filename   string    `json:"filename"`
	UploadTime time.Time `json:"upload_time"`
}

// AvatarUploadResponse 头像上传响应
type TeamAvatarUploadResponse struct {
	Success  bool   `json:"success"`
	Filename string `json:"filename,omitempty"`
	Error    string `json:"error,omitempty"`
}

// TeamMemberListRequest 团队成员列表请求
type TeamMemberListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Role     string `form:"role"`    // 可选，按角色筛选
	Keyword  string `form:"keyword"` // 可选，搜索用户名
}

// TeamMemberInfo 团队成员信息
type TeamMemberInfo struct {
	UserID   uint64    `json:"user_id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	Nickname string    `json:"nickname,omitempty"` // 团队内名称
}

// TeamMemberListResponse 团队成员列表响应
type TeamMemberListResponse struct {
	Members  []TeamMemberInfo `json:"members"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// TeamNickname 团队内名称
type TeamNickname struct {
	TeamID    uint64    `json:"team_id"`
	UserID    uint64    `json:"user_id"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateTeamNicknameRequest 更新团队内名称请求
type UpdateTeamNicknameRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
}

// AvailableProblemListRequest 获取可用题目列表请求
type AvailableProblemListRequest struct {
	TeamID   uint64 `form:"team_id" binding:"required"`
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Keyword  string `form:"keyword"`                                        // 搜索关键字
	Type     string `form:"type" binding:"omitempty,oneof=all global team"` // 题目类型：all-所有题目，global-全局题目，team-团队题目
}

// AvailableProblemInfo 可用题目信息
type AvailableProblemInfo struct {
	ID            uint64    `json:"id"`              // 题目ID
	Type          string    `json:"type"`            // 题目类型：global-全局题目，team-团队题目
	TeamProblemID uint64    `json:"team_problem_id"` // 团队题目ID，当type为team时有效
	Title         string    `json:"title"`           // 题目标题
	TimeLimit     int       `json:"time_limit"`      // 时间限制
	MemoryLimit   int       `json:"memory_limit"`    // 内存限制
	Tags          string    `json:"tags"`            // 标签列表，JSON字符串
	Difficulty    string    `json:"difficulty"`      // 难度等级：beginner-入门，easy-简单，medium-中等，hard-困难，expert-专家
	CreatedBy     uint64    `json:"created_by"`      // 创建者ID
	CreatedAt     time.Time `json:"created_at"`      // 创建时间
}

// AvailableProblemTag 可用题目标签
type AvailableProblemTag struct {
	ID    uint64 `json:"id"`    // 标签ID
	Name  string `json:"name"`  // 标签名称
	Color string `json:"color"` // 标签颜色
}

// AvailableProblemListResponse 获取可用题目列表响应
type AvailableProblemListResponse struct {
	Problems []AvailableProblemInfo `json:"problems"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

// AssignmentProblemDetail 作业题目详情
type AssignmentProblemDetail struct {
	ID              uint64          `json:"id"`               // 题目ID
	Type            string          `json:"type"`             // 题目类型：global-全局题目，team-团队题目
	TeamProblemID   *uint64         `json:"team_problem_id"`  // 团队题目ID，当type为team时有效
	Title           string          `json:"title"`            // 题目标题
	TimeLimit       int             `json:"time_limit"`       // 时间限制
	MemoryLimit     int             `json:"memory_limit"`     // 内存限制
	Tags            string          `json:"tags"`             // 标签列表，JSON字符串
	Difficulty      string          `json:"difficulty"`       // 难度等级
	Score           int             `json:"score"`            // 分值
	OrderIndex      int             `json:"order_index"`      // 题目顺序
	SubmissionStats SubmissionStats `json:"submission_stats"` // 提交统计
}

// SubmissionStats 提交统计
type SubmissionStats struct {
	TotalCount    int            `json:"total_count"`    // 总提交数
	AcceptedCount int            `json:"accepted_count"` // 通过数
	AcceptedRate  float64        `json:"accepted_rate"`  // 通过率
	StatusCounts  map[string]int `json:"status_counts"`  // 各状态数量
}

// GetAssignmentProblemsRequest 获取作业题目列表请求
type GetAssignmentProblemsRequest struct {
	AssignmentID uint64 `form:"assignment_id" binding:"required"`
	TeamID       uint64 `form:"team_id" binding:"required"`
}

// GetAssignmentProblemsResponse 获取作业题目列表响应
type GetAssignmentProblemsResponse struct {
	Problems []AssignmentProblemDetail `json:"problems"`
}

// GetAssignmentProblemDetailRequest 获取作业题目详情请求
type GetAssignmentProblemDetailRequest struct {
	AssignmentID uint64 `form:"assignment_id" binding:"required"`
	ProblemID    uint64 `form:"problem_id" binding:"required"`
	ProblemType  string `form:"problem_type" binding:"required,oneof=global team"`
	TeamID       uint64 `form:"team_id" binding:"required"`
}

// AssignmentProblemFullDetail 作业题目完整详情
type AssignmentProblemFullDetail struct {
	ID                uint64            `json:"id"`                 // 题目ID
	Type              string            `json:"type"`               // 题目类型：global-全局题目，team-团队题目
	TeamProblemID     *uint64           `json:"team_problem_id"`    // 团队题目ID，当type为team时有效
	Title             string            `json:"title"`              // 题目标题
	Description       string            `json:"description"`        // 题目描述
	InputDescription  string            `json:"input_description"`  // 输入说明
	OutputDescription string            `json:"output_description"` // 输出说明
	TimeLimit         int               `json:"time_limit"`         // 时间限制
	MemoryLimit       int               `json:"memory_limit"`       // 内存限制
	Tags              []ProblemTag      `json:"tags"`               // 标签列表
	Difficulty        string            `json:"difficulty"`         // 难度等级
	Score             int               `json:"score"`              // 分值
	SampleCases       string            `json:"sample_cases"`       // 样例
	Hint              string            `json:"hint"`               // 提示
	Source            string            `json:"source"`             // 来源
	DifficultySystem  DifficultySystem  `json:"difficulty_system"`  // 难度等级系统
	CreatedBy         uint64            `json:"created_by"`         // 创建者ID
	CreatedAt         time.Time         `json:"created_at"`         // 创建时间
	UpdatedAt         time.Time         `json:"updated_at"`         // 更新时间
	Categories        []ProblemCategory `json:"categories"`         // 分类
	UserStatus        *string           `json:"user_status"`        // 用户状态：null-未提交, accepted-已通过, attempted-尝试过
	SubmissionStats   SubmissionStats   `json:"submission_stats"`   // 提交统计
}

// SubmitAssignmentCodeRequest 提交作业代码请求
type SubmitAssignmentCodeRequest struct {
	AssignmentID uint64 `json:"assignment_id" binding:"required"`
	ProblemID    uint64 `json:"problem_id" binding:"required"`
	ProblemType  string `json:"problem_type" binding:"required,oneof=global team"`
	TeamID       uint64 `json:"team_id" binding:"required"`
	Language     string `json:"language" binding:"required"`
	Code         string `json:"code" binding:"required"`
}

// GetAssignmentSubmissionsRequest 获取作业提交记录请求
type GetAssignmentSubmissionsRequest struct {
	AssignmentID uint64 `form:"assignment_id" binding:"required"`
	TeamID       uint64 `form:"team_id" binding:"required"`
	UserID       uint64 `form:"user_id"`                                          // 可选，按用户筛选
	ProblemID    uint64 `form:"problem_id"`                                       // 可选，按题目筛选
	Status       string `form:"status"`                                           // 可选，按状态筛选
	OrderBy      string `form:"order_by" binding:"omitempty,oneof=id created_at"` // 排序字段：id-按提交ID，created_at-按提交时间
	OrderType    string `form:"order_type" binding:"omitempty,oneof=asc desc"`    // 排序方式：asc-升序，desc-降序
	Page         int    `form:"page" binding:"required,min=1"`
	PageSize     int    `form:"page_size" binding:"required,min=1,max=100"`
}

// AssignmentSubmissionInfo 作业提交记录信息
type AssignmentSubmissionInfo struct {
	ID           uint64    `json:"id"`            // 提交ID
	ProblemID    uint64    `json:"problem_id"`    // 题目ID
	ProblemType  string    `json:"problem_type"`  // 题目类型：global-全局题目，team-团队题目
	ProblemTitle string    `json:"problem_title"` // 题目标题
	UserID       uint64    `json:"user_id"`       // 用户ID
	Username     string    `json:"username"`      // 用户名
	Nickname     string    `json:"nickname"`      // 团队内名称
	Language     string    `json:"language"`      // 编程语言
	Status       string    `json:"status"`        // 判题状态
	TimeUsed     int       `json:"time_used"`     // 运行时间（毫秒）
	MemoryUsed   int       `json:"memory_used"`   // 内存使用（KB）
	Score        int       `json:"score"`         // 题目分值
	CreatedAt    time.Time `json:"created_at"`    // 提交时间
}

// GetAssignmentSubmissionsResponse 获取作业提交记录响应
type GetAssignmentSubmissionsResponse struct {
	Submissions []AssignmentSubmissionInfo `json:"submissions"`
	Total       int64                      `json:"total"`
	Page        int                        `json:"page"`
	PageSize    int                        `json:"page_size"`
}
