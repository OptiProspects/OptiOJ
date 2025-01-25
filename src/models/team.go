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
	AssignmentID uint64 `json:"assignment_id"`
	ProblemID    uint64 `json:"problem_id"`
	OrderIndex   int    `json:"order_index"`
	Score        int    `json:"score"`
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
		ProblemID  uint64 `json:"problem_id" binding:"required"`
		OrderIndex int    `json:"order_index"`
		Score      int    `json:"score"`
	} `json:"problems" binding:"required"`
}

// UpdateAssignmentRequest 更新作业请求
type UpdateAssignmentRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Problems    []struct {
		ProblemID  uint64 `json:"problem_id"`
		OrderIndex int    `json:"order_index"`
		Score      int    `json:"score"`
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
