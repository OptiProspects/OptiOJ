package models

import "time"

type Admin struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Role      string    `json:"role"`       // 角色类型：super_admin, admin
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// 设置表名
func (Admin) TableName() string {
	return "admins"
}

// AdminDetailItem 管理员详细信息
type AdminDetailItem struct {
	ID            uint      `json:"id"`
	UserID        uint64    `json:"user_id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	LastLoginTime time.Time `json:"last_login_time"`
	LastLoginIP   string    `json:"last_login_ip"`
}
