/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package response

import (
	"bamboo-main/internal/model/dto"
	"time"
)

// AuthLoginResponse 用户登录响应
type AuthLoginResponse struct {
	User      dto.SystemUserDetailDTO `json:"user"`       // 用户信息
	Token     string                  `json:"token"`      // 访问令牌
	CreatedAt time.Time               `json:"created_at"` // Token创建时间
	ExpiredAt time.Time               `json:"expired_at"` // Token过期时间
}

// AuthRegisterResponse 用户注册响应
type AuthRegisterResponse struct {
	User      dto.SystemUserDetailDTO `json:"user"`       // 用户信息
	Token     string                  `json:"token"`      // 访问令牌
	CreatedAt time.Time               `json:"created_at"` // Token创建时间
	ExpiredAt time.Time               `json:"expired_at"` // Token过期时间
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
