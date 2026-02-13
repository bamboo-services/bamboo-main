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

package dto

import (
	"time"
)

// SystemUserSimpleDTO 系统用户简单DTO - 用于列表查询
type SystemUserSimpleDTO struct {
	ID       int64  `json:"id"`       // 用户主键
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 昵称
}

// SystemUserNormalDTO 系统用户标准DTO - 用于分页查询
type SystemUserNormalDTO struct {
	ID          int64      `json:"id"`            // 用户主键
	Username    string     `json:"username"`      // 用户名
	Email       string     `json:"email"`         // 邮箱
	Nickname    string     `json:"nickname"`      // 昵称
	Role        string     `json:"role"`          // 角色
	Status      int        `json:"status"`        // 状态
	EmailVerify bool       `json:"email_verify"`  // 邮箱验证状态
	LastLoginAt *time.Time `json:"last_login_at"` // 最后登录时间
	CreatedAt   time.Time  `json:"created_at"`    // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`    // 更新时间
}

// SystemUserDetailDTO 系统用户详细DTO - 用于详情查询
type SystemUserDetailDTO struct {
	ID          int64      `json:"id"`            // 用户主键
	Username    string     `json:"username"`      // 用户名
	Email       string     `json:"email"`         // 邮箱
	Nickname    *string    `json:"nickname"`      // 昵称
	Avatar      *string    `json:"avatar"`        // 头像
	Role        string     `json:"role"`          // 角色
	Status      int        `json:"status"`        // 状态
	EmailVerify bool       `json:"email_verify"`  // 邮箱验证状态
	LastLoginAt *time.Time `json:"last_login_at"` // 最后登录时间
	CreatedAt   time.Time  `json:"created_at"`    // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`    // 更新时间
}
