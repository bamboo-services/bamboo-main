/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package apiAuth

import (
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
)

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string  `json:"username" binding:"required,min=1,max=50" example:"admin"`
	Email    string  `json:"email" binding:"required,email" example:"admin@example.com"`
	Nickname *string `json:"nickname" binding:"omitempty,min=1,max=50" example:"筱锋"`
	Password string  `json:"password" binding:"required,min=6,max=100" example:"password123"`
	Code     string  `json:"code" binding:"required,len=6" example:"123456"`
}

// RegisterCodeRequest 发送注册邮箱验证码请求
type RegisterCodeRequest struct {
	Email string `json:"email" binding:"required,email" example:"admin@example.com"`
}

// RegisterResponse 用户注册响应
type RegisterResponse struct {
	User      entity.SystemUser `json:"user"`
	Token     string            `json:"token"`
	CreatedAt time.Time         `json:"created_at"`
	ExpiredAt time.Time         `json:"expired_at"`
}
