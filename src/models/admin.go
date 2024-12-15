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
