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
	dto.SystemUserDetailDTO
}

// MessageResponse 通用消息响应
type MessageResponse struct {
	Message string `json:"message"` // 响应消息
}

// AuthLogoutResponse 用户登出响应
type AuthLogoutResponse = MessageResponse

// AuthPasswordChangeResponse 密码修改响应
type AuthPasswordChangeResponse = MessageResponse

// AuthPasswordResetResponse 密码重置响应
type AuthPasswordResetResponse = MessageResponse