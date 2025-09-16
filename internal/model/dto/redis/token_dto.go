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

package dtoRedis

import "time"

// TokenDTO Token令牌数据传输对象
type TokenDTO struct {
	UserUUID  string    `json:"user_uuid"`  // 用户UUID
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 用户邮箱
	Role      string    `json:"role"`       // 用户角色
	LoginIP   string    `json:"login_ip"`   // 登录IP地址
	UserAgent string    `json:"user_agent"` // 用户代理
	CreatedAt time.Time `json:"created_at"` // Token创建时间
	ExpiredAt time.Time `json:"expired_at"` // Token过期时间
}
