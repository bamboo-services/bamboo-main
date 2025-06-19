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

package base

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// UserTokenDTO
//
// 用户令牌数据传输对象；
// 包含用户令牌、刷新令牌、用户代理、用户IP地址、创建时间、过期时间和刷新时间等信息；
// 该结构体用于在用户认证和授权过程中传递令牌相关信息。
type UserTokenDTO struct {
	Token        uuid.UUID   `json:"token_uuid" dc:"用户令牌"`
	RefreshToken uuid.UUID   `json:"refresh_token_uuid" dc:"用户刷新令牌"`
	UserAgent    string      `json:"user_agent" dc:"用户代理"`
	UserIP       string      `json:"user_ip" dc:"用户IP地址"`
	CreatedAt    *gtime.Time `json:"created_at" dc:"创建时间"`
	ExpiresAt    *gtime.Time `json:"expires_at" dc:"过期时间"`
	RefreshAt    *gtime.Time `json:"refresh_at" dc:"刷新时间"`
}
