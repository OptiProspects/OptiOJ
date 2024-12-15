package models

import "time"

type UserBan struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Reason    string    `json:"reason"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uint64    `json:"created_by"`
	IsActive  bool      `json:"is_active"`
}

// TableName 设置表名
func (UserBan) TableName() string {
	return "user_bans"
}
