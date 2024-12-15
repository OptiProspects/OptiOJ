package models

import "time"

type UserLogin struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	LoginTime   time.Time `json:"login_time"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	LoginStatus string    `json:"login_status"`
	FailReason  string    `json:"fail_reason,omitempty"`
	Location    string    `json:"location,omitempty"`
}

// TableName 设置表名
func (UserLogin) TableName() string {
	return "user_logins"
}

// LoginHistoryRequest 登录历史查询请求
type LoginHistoryRequest struct {
	UserID    uint      `form:"user_id"`
	StartTime time.Time `form:"start_time"`
	EndTime   time.Time `form:"end_time"`
	Status    string    `form:"status"`
	Page      int       `form:"page" binding:"required,min=1"`
	PageSize  int       `form:"page_size" binding:"required,min=1,max=100"`
}
