package models

import "time"

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
}

type RegisterRequest struct {
	UserName         string `json:"userName"`
	PassWord         string `json:"passWord"`
	RequestEmail     string `json:"requestEmail"`
	RequestPhone     string `json:"requestPhone"`
	VerificationCode string `json:"verificationCode"`
	VerificationType string `json:"verificationType"` // "email" 或 "phone"
}

type LoginRequest struct {
	AccountInfo string `json:"accountInfo"`
	PassWord    string `json:"passWord"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Username string `form:"username"`
	Email    string `form:"email"`
	Phone    string `form:"phone"`
	Status   string `form:"status"` // normal, banned
}

// UserUpdateRequest 更新用户信息请求
type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// UserBanRequest 封禁用户请求
type UserBanRequest struct {
	UserID     uint      `json:"user_id" binding:"required"`
	BanReason  string    `json:"ban_reason" binding:"required"`
	BanExpires time.Time `json:"ban_expires"` // 可选，为空表示永久封禁
}

// 在 User 结构体下面添加新的结构体
type UserListItem struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Status        string    `json:"status"`
	BanReason     string    `json:"ban_reason,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastLoginTime time.Time `json:"last_login_time" gorm:"column:last_login_time"`
	LastLoginIP   string    `json:"last_login_ip" gorm:"column:last_login_ip"`
	Role          string    `json:"role"`
}

// GenerateUsersRequest 批量生成用户请求
type GenerateUsersRequest struct {
	Prefix string `json:"prefix" binding:"required"`               // 用户名前缀
	Suffix string `json:"suffix"`                                  // 用户名后缀（可选）
	Count  int    `json:"count" binding:"required,min=1,max=1000"` // 生成数量
	Domain string `json:"domain" binding:"required"`               // 邮箱域名
}

// GeneratedUserInfo 生成的用户信息
type GeneratedUserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// GenerateUsersResponse 批量生成用户响应
type GenerateUsersResponse struct {
	Users []GeneratedUserInfo `json:"users"`
	Total int                 `json:"total"`
}
