package response

import (
	"bamboo-main/internal/model/dto"
)

// AuthLoginResponse 用户登录响应
type AuthLoginResponse struct {
	User  dto.SystemUserDetailDTO `json:"user"`  // 用户信息
	Token string                  `json:"token"` // 访问令牌
}

// AuthUserInfoResponse 用户信息响应
type AuthUserInfoResponse struct {
	Data dto.SystemUserDetailDTO `json:"data"` // 用户详细信息
}

// AuthLogoutResponse 用户登出响应
type AuthLogoutResponse struct {
	Message string `json:"message"` // 响应消息
}

// AuthPasswordChangeResponse 密码修改响应
type AuthPasswordChangeResponse struct {
	Message string `json:"message"` // 响应消息
}

// AuthPasswordResetResponse 密码重置响应
type AuthPasswordResetResponse struct {
	Message string `json:"message"` // 响应消息
}