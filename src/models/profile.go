package models

import (
	"encoding/json"
	"strings"
	"time"
)

type Profile struct {
	ID       int        `json:"id" gorm:"primaryKey"`
	UserID   int        `json:"user_id"`
	Bio      string     `json:"bio"`               // 个人签名
	Gender   string     `json:"gender"`            // 性别
	School   string     `json:"school"`            // 学校
	Birthday *time.Time `json:"birthday"`          // 生日（带时区）
	Location string     `json:"-"`                 // 现居地(内部存储用)
	Province string     `json:"province" gorm:"-"` // 省份（仅用于JSON）
	City     string     `json:"city" gorm:"-"`     // 城市（仅用于JSON）
	RealName string     `json:"real_name"`         // 真实姓名
	CreateAt time.Time  `json:"create_at"`         // 创建时间
	UpdateAt time.Time  `json:"update_at"`         // 更新时间
}

// UnmarshalJSON 实现自定义的 JSON 解析
func (p *Profile) UnmarshalJSON(data []byte) error {
	// 创建一个临时结构体来解析基本字段
	type Alias Profile
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 如果 Location 不为空，解析出省份和城市
	if p.Location != "" {
		parts := strings.Split(p.Location, "-")
		if len(parts) == 2 {
			p.Province = parts[0]
			p.City = parts[1]
		}
	}
	return nil
}

// MarshalJSON 添加自定义的 JSON 序列化
func (p Profile) MarshalJSON() ([]byte, error) {
	type Alias Profile

	// 解析 Location 字段
	province, city := "", ""
	if p.Location != "" {
		parts := strings.Split(p.Location, "-")
		if len(parts) == 2 {
			province = parts[0]
			city = parts[1]
		}
	}

	return json.Marshal(&struct {
		Alias
		Birthday string `json:"birthday"`
		Province string `json:"province"`
		City     string `json:"city"`
	}{
		Alias: (Alias)(p),
		Birthday: func() string {
			if p.Birthday == nil {
				return ""
			}
			return p.Birthday.Format(time.RFC3339)
		}(),
		Province: province,
		City:     city,
	})
}

type UpdateProfileRequest struct {
	Bio      string `json:"bio"`
	Gender   string `json:"gender"`
	School   string `json:"school"`
	Birthday string `json:"birthday"`
	Province string `json:"province"`
	City     string `json:"city"`
	RealName string `json:"real_name"`
}

// ActivityLevel 活跃度等级
type ActivityLevel int

const (
	ActivityLevelNone   ActivityLevel = 0 // 无提交
	ActivityLevelLow    ActivityLevel = 1 // 低活跃度 (1-3次)
	ActivityLevelMedium ActivityLevel = 2 // 中等活跃度 (4-6次)
	ActivityLevelHigh   ActivityLevel = 3 // 高活跃度 (7-9次)
	ActivityLevelSuper  ActivityLevel = 4 // 超高活跃度 (10次及以上)
)

// DailyActivity 每日活跃度
type DailyActivity struct {
	Date  string        `json:"date"`  // 日期，格式：YYYY-MM-DD
	Count int           `json:"count"` // 提交次数
	Level ActivityLevel `json:"level"` // 活跃度等级
}

// ActivityResponse 活跃度响应
type ActivityResponse struct {
	Activities []DailyActivity `json:"activities"`  // 活跃度数据
	TotalDays  int             `json:"total_days"`  // 统计天数
	MaxCount   int             `json:"max_count"`   // 单日最大提交次数
	TotalCount int             `json:"total_count"` // 总提交次数
	AcceptRate float64         `json:"accept_rate"` // 90天内通过率
}
