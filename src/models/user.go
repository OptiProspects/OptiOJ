package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
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
