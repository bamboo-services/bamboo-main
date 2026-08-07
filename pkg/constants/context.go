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

	// LinkStatusPending 链接申请阶段待审核状态（新友链申请，等待管理员审核）
	LinkStatusPending = 0
	// LinkStatusApproved 链接已通过状态
	LinkStatusApproved = 1
	// LinkStatusRejected 链接已拒绝状态
	LinkStatusRejected = 2
	// LinkStatusTakedownPending 链接下架待审核状态（用户申请下架，等待管理员审核）
	LinkStatusTakedownPending = 3
	// LinkStatusTakenDown 链接已下架状态（管理员批准下架，不再公开展示）
	LinkStatusTakenDown = 4
	// LinkStatusEditPending 链接修改阶段待审核状态（用户申请修改展示位置/颜色，等待管理员审核）
	LinkStatusEditPending = 5

	// SponsorStatusPending 赞助待审核状态（公开申请后默认，前台不可见）
	SponsorStatusPending = 0
	// SponsorStatusApproved 赞助已通过状态（通过审核后前台展示）
	SponsorStatusApproved = 1
	// SponsorStatusRejected 赞助已拒绝状态（拒绝后前台不可见）
	SponsorStatusRejected = 2

	// LinkFailNormal 链接正常状态
	LinkFailNormal = 0
	// LinkFailBroken 链接失效状态
	LinkFailBroken = 1

	// LinkLevelRegular 一般友链（1×1 紧凑卡片）
	LinkLevelRegular = 0
	// LinkLevelClose 好友友链（1×1 富式卡片）
	LinkLevelClose = 1
	// LinkLevelPremium 高级友链（2×2 特写卡片，hover 展开站点截图）
	LinkLevelPremium = 2
	// LinkLevelAd 广告友链（1×1 居中卡片，带推广标识）
	LinkLevelAd = 3

	// EmailTypeApply 申请邮件类型
	EmailTypeApply = "apply"
	// EmailTypeApproved 审核通过邮件类型
	EmailTypeApproved = "approved"
	// EmailTypeRejected 审核拒绝邮件类型
	EmailTypeRejected = "rejected"
	// EmailTypePasswordReset 密码重置邮件类型
	EmailTypePasswordReset = "password_reset"

	// SponsorEmailTypeApply 赞助申请通知邮件模板名（通知管理员）
	SponsorEmailTypeApply = "sponsor_apply"
	// SponsorEmailTypeApproved 赞助审核通过邮件模板名（通知申请者）
	SponsorEmailTypeApproved = "sponsor_approved"
	// SponsorEmailTypeRejected 赞助审核拒绝邮件模板名（通知申请者）
	SponsorEmailTypeRejected = "sponsor_rejected"
)

const (
	// ContextCustomConfig 自定义配置上下文键
	ContextCustomConfig xCtx.ContextKey = "context_custom_config"
	// ContextEmailTemplate 内嵌邮件模板注入节点上下文键
	ContextEmailTemplate xCtx.ContextKey = "context_email_template"
	// ContextScreenshotManager 友链截图任务管理器上下文键
	ContextScreenshotManager xCtx.ContextKey = "context_screenshot_manager"
)
