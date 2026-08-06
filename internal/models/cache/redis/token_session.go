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

package rediscache

import (
	"time"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
)

// TokenSession 登录会话数据，序列化后存储于 Redis
type TokenSession struct {
	UserID    xSnowflake.SnowflakeID `json:"user_id"`
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	LoginIP   string                 `json:"login_ip"`
	UserAgent string                 `json:"user_agent"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiredAt time.Time              `json:"expired_at"`
}
