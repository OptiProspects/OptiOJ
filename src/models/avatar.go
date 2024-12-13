package models

import "time"

type Avatar struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Filename   string    `json:"filename"`
	UploadTime time.Time `json:"upload_time"`
}

type AvatarUploadResponse struct {
	Success  bool   `json:"success"`
	Filename string `json:"filename,omitempty"`
	Error    string `json:"error,omitempty"`
}
