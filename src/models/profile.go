package models

import (
	"encoding/json"
	"strings"
	"time"
)

type Profile struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Bio      string    `json:"bio"`       // 个人签名
	Gender   string    `json:"gender"`    // 性别
	School   string    `json:"school"`    // 学校
	Birthday time.Time `json:"birthday"`  // 生日（带时区）
	Location string    `json:"-"`         // 现居地(内部存储用)
	Province string    `json:"province"`  // 省份
	City     string    `json:"city"`      // 城市
	RealName string    `json:"real_name"` // 真实姓名
	CreateAt time.Time `json:"create_at"` // 创建时间
	UpdateAt time.Time `json:"update_at"` // 更新时间
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
		Alias:    (Alias)(p),
		Birthday: p.Birthday.Format(time.RFC3339), // 使用 ISO 8601 格式输出带时区的时间
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
