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

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=1,max=50" example:"admin"`
	Password string `json:"password" binding:"required,min=6,max=100" example:"password123"`
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	User      entity.SystemUser `json:"user"`
	Token     string            `json:"token"`
	CreatedAt time.Time         `json:"created_at"`
	ExpiredAt time.Time         `json:"expired_at"`
}
