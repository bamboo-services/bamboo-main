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

package constants

import xCtx "github.com/bamboo-services/bamboo-base-go/defined/context"

const (
	// ContextKeyUser 用户信息上下文键
	ContextKeyUser = "user"
	// ContextKeyUserID 用户 ID 上下文键
	ContextKeyUserID = "user_id"
	// ContextKeyToken 认证令牌上下文键
	ContextKeyToken = "token"

	// HeaderAuthorization 认证请求头名称
	HeaderAuthorization = "Authorization"
	// HeaderContentType 内容类型请求头名称
	HeaderContentType = "Content-Type"

	// TokenPrefix Bearer 认证令牌前缀
	TokenPrefix = "Bearer "

	// StatusActive 启用状态值
	StatusActive = 1
	// StatusInactive 停用状态值
	StatusInactive = 0

	// LinkStatusPending 链接待审核状态
	LinkStatusPending = 0
	// LinkStatusApproved 链接已通过状态
	LinkStatusApproved = 1
	// LinkStatusRejected 链接已拒绝状态
	LinkStatusRejected = 2

	// LinkFailNormal 链接正常状态
	LinkFailNormal = 0
	// LinkFailBroken 链接失效状态
	LinkFailBroken = 1

	// RoleAdmin 管理员角色标识
	RoleAdmin = "admin"
	// RoleModerator 审核员角色标识
	RoleModerator = "moderator"
	// RoleUser 普通用户角色标识
	RoleUser = "user"

	// EmailTypeApply 申请邮件类型
	EmailTypeApply = "apply"
	// EmailTypeApproved 审核通过邮件类型
	EmailTypeApproved = "approved"
	// EmailTypeRejected 审核拒绝邮件类型
	EmailTypeRejected = "rejected"
	// EmailTypePasswordReset 密码重置邮件类型
	EmailTypePasswordReset = "password_reset"
)

const (
	// ContextCustomConfig 自定义配置上下文键
	ContextCustomConfig xCtx.ContextKey = "context_custom_config"
)
