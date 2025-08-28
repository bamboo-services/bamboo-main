package helper

import "time"

// UserSession 用户会话结构
type UserSession struct {
	UserUUID string    `json:"user_uuid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	LoginAt  time.Time `json:"login_at"`
	ExpireAt time.Time `json:"expire_at"`
}