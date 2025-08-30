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

package constants

const (
	// 上下文键
	ContextKeyUser   = "user"
	ContextKeyUserID = "user_id"
	ContextKeyToken  = "token"

	// HTTP 头部
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"

	// 认证相关
	TokenPrefix = "Bearer "

	// 系统状态
	StatusActive   = 1
	StatusInactive = 0

	// 链接状态
	LinkStatusPending  = 0 // 待审核
	LinkStatusApproved = 1 // 已通过
	LinkStatusRejected = 2 // 已拒绝

	// 链接失效状态
	LinkFailNormal = 0 // 正常
	LinkFailBroken = 1 // 失效

	// 邮件类型
	EmailTypeApply         = "apply"
	EmailTypeApproved      = "approved"
	EmailTypeRejected      = "rejected"
	EmailTypePasswordReset = "password_reset"
)
