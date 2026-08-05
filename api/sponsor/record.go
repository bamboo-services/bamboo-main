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

package apiSponsor

import (
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/models/base"
)

// RecordIDRequest 赞助记录路径参数
type RecordIDRequest struct {
	ID xSnowflake.SnowflakeID `uri:"id" binding:"required"`
}

// RecordAddRequest 添加赞助记录请求
type RecordAddRequest struct {
	Nickname    string                  `json:"nickname" binding:"required,min=1,max=100" example:"张三"`                 // 赞助者昵称
	RedirectURL *string                 `json:"redirect_url" binding:"omitempty,max=500" example:"https://example.com"` // 跳转链接
	Amount      int64                   `json:"amount" binding:"required,min=1" example:"1000"`                         // 赞助金额(分)
	ChannelID   *xSnowflake.SnowflakeID `json:"channel_id" binding:"omitempty" example:"123456789"`                     // 赞助渠道 ID
	Message     *string                 `json:"message" binding:"omitempty,max=500" example:"感谢开发者"`                    // 留言信息
	SponsorAt   *time.Time              `json:"sponsor_at" binding:"omitempty" example:"2025-01-01T12:00:00Z"`          // 赞助时间
	SortOrder   int                     `json:"sort_order" binding:"omitempty,min=0" example:"0"`                       // 排序值
	IsAnonymous bool                    `json:"is_anonymous" binding:"omitempty" example:"false"`                       // 是否匿名
	IsHidden    bool                    `json:"is_hidden" binding:"omitempty" example:"false"`                          // 是否隐藏
}

// RecordUpdateRequest 更新赞助记录请求
type RecordUpdateRequest struct {
	Nickname    *string                 `json:"nickname" binding:"omitempty,min=1,max=100" example:"张三"`                // 赞助者昵称
	RedirectURL *string                 `json:"redirect_url" binding:"omitempty,max=500" example:"https://example.com"` // 跳转链接
	Amount      *int64                  `json:"amount" binding:"omitempty,min=1" example:"1000"`                        // 赞助金额(分)
	ChannelID   *xSnowflake.SnowflakeID `json:"channel_id" binding:"omitempty" example:"123456789"`                     // 赞助渠道 ID
	Message     *string                 `json:"message" binding:"omitempty,max=500" example:""`                         // 留言信息
	SponsorAt   *time.Time              `json:"sponsor_at" binding:"omitempty" example:""`                              // 赞助时间
	SortOrder   *int                    `json:"sort_order" binding:"omitempty,min=0" example:"0"`                       // 排序值
	IsAnonymous *bool                   `json:"is_anonymous" binding:"omitempty" example:"false"`                       // 是否匿名
	IsHidden    *bool                   `json:"is_hidden" binding:"omitempty" example:"false"`                          // 是否隐藏
}

// RecordPageRequest 分页查询请求(后台)
type RecordPageRequest struct {
	Page        int                     `form:"page" binding:"omitempty,min=1" example:"1"`                                                               // 页码,默认1
	PageSize    int                     `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`                                                 // 每页数量,默认10,最大100
	ChannelID   *xSnowflake.SnowflakeID `form:"channel_id" binding:"omitempty" example:"123456789"`                                                       // 渠道 ID 过滤
	Nickname    *string                 `form:"nickname" binding:"omitempty,max=100" example:"张"`                                                         // 昵称模糊搜索
	IsAnonymous *bool                   `form:"is_anonymous" binding:"omitempty" example:"false"`                                                         // 是否匿名过滤
	IsHidden    *bool                   `form:"is_hidden" binding:"omitempty" example:"false"`                                                            // 是否隐藏过滤
	Status      *int                    `form:"status" binding:"omitempty,oneof=0 1 2" example:"1"`                                                       // 审核状态过滤（0:待审核 1:已通过 2:已拒绝）
	OrderBy     *string                 `form:"order_by" binding:"omitempty,oneof=nickname amount sponsor_at sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string                 `form:"order" binding:"omitempty,oneof=asc desc" example:"desc"`                                                  // 排序方向
}

// RecordPublicPageRequest 公开分页查询请求(前台)
type RecordPublicPageRequest struct {
	Page      int                     `form:"page" binding:"omitempty,min=1" example:"1"`                                           // 页码,默认1
	PageSize  int                     `form:"page_size" binding:"omitempty,min=1,max=50" example:"20"`                              // 每页数量,默认20,最大50
	ChannelID *xSnowflake.SnowflakeID `form:"channel_id" binding:"omitempty" example:"123456789"`                                   // 渠道 ID 过滤
	OrderBy   *string                 `form:"order_by" binding:"omitempty,oneof=amount sponsor_at sort_order" example:"sort_order"` // 排序字段
	Order     *string                 `form:"order" binding:"omitempty,oneof=asc desc" example:"desc"`                              // 排序方向
}

// SponsorApplyRequest 访客自助申请赞助展示请求
//
// 面向游客与登录用户的公开申请入口：仅需赞助基本信息，联系邮箱必填（用于确认归属与结果通知），
// 金额与渠道由申请者如实填报、管理员审核时核验；is_hidden/sort_order 等管理员专属字段不在此开放。
type SponsorApplyRequest struct {
	Nickname    string                  `json:"nickname" binding:"required,min=1,max=100" example:"张三"`                 // 赞助者昵称
	Amount      int64                   `json:"amount" binding:"required,min=1" example:"1000"`                         // 赞助金额(分)
	ChannelID   *xSnowflake.SnowflakeID `json:"channel_id" binding:"omitempty" example:"123456789"`                     // 赞助渠道 ID
	Message     *string                 `json:"message" binding:"omitempty,max=500" example:"感谢开发者"`                    // 留言信息
	SponsorAt   *time.Time              `json:"sponsor_at" binding:"omitempty" example:"2025-01-01T12:00:00Z"`          // 赞助发生时间
	Email       string                  `json:"email" binding:"required,email,max=100" example:"admin@example.com"`      // 联系邮箱（用于归属确认）
	RedirectURL *string                 `json:"redirect_url" binding:"omitempty,url,max=500" example:"https://example.com"` // 点击昵称的跳转地址
	IsAnonymous bool                    `json:"is_anonymous" binding:"omitempty" example:"false"`                        // 是否匿名展示
	ApplyRemark *string                 `json:"apply_remark" binding:"omitempty,max=500" example:"感谢支持"`                // 申请者备注
}

// SponsorStatusRequest 更新赞助记录审核状态请求
type SponsorStatusRequest struct {
	SponsorStatus       int    `json:"sponsor_status" binding:"required,oneof=0 1 2" example:"1"` // 审核状态（0:待审核 1:已通过 2:已拒绝）
	SponsorReviewRemark string `json:"sponsor_review_remark" binding:"omitempty,max=500" example:"审核通过"` // 审核备注（拒绝原因）
}

// SponsorUserUpdateRequest 用户更新自己赞助申请请求
//
// 仅允许更新展示类基础字段，金额/渠道需与实际支付核验，不允许用户自行修改。
type SponsorUserUpdateRequest struct {
	Nickname    *string                 `json:"nickname" binding:"omitempty,min=1,max=100" example:"张三"`                // 赞助者昵称
	RedirectURL *string                 `json:"redirect_url" binding:"omitempty,url,max=500" example:"https://example.com"` // 点击昵称的跳转地址
	Message     *string                 `json:"message" binding:"omitempty,max=500" example:"感谢开发者"`                    // 留言信息
	SponsorAt   *time.Time              `json:"sponsor_at" binding:"omitempty" example:"2025-01-01T12:00:00Z"`          // 赞助发生时间
	IsAnonymous *bool                   `json:"is_anonymous" binding:"omitempty" example:"false"`                        // 是否匿名展示
	ApplyRemark *string                 `json:"apply_remark" binding:"omitempty,max=500" example:"感谢支持"`                // 申请者备注
}

// SponsorUserQueryRequest 用户查询自己赞助申请请求
type SponsorUserQueryRequest struct {
	Page          int  `form:"page" binding:"omitempty,min=1" example:"1"`            // 页码,默认1
	PageSize      int  `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"` // 每页数量,默认10,最大100
	SponsorStatus *int `form:"sponsor_status" binding:"omitempty,oneof=0 1 2" example:"1"` // 审核状态过滤
}

// SponsorChannelSimpleResponse 赞助渠道简要响应
type SponsorChannelSimpleResponse struct {
	ID   xSnowflake.SnowflakeID `json:"id"`
	Name string                 `json:"name"`
	Icon *string                `json:"icon"`
}

// RecordEntityResponse 赞助记录实体响应
type RecordEntityResponse struct {
	ID           xSnowflake.SnowflakeID        `json:"id"`
	Nickname     string                        `json:"nickname"`
	RedirectURL  *string                       `json:"redirect_url"`
	Amount       int64                         `json:"amount"`
	ChannelID    *xSnowflake.SnowflakeID       `json:"channel_id"`
	Message      *string                       `json:"message"`
	SponsorAt    *time.Time                    `json:"sponsor_at"`
	SortOrder    int                           `json:"sort_order"`
	IsAnonymous  bool                          `json:"is_anonymous"`
	IsHidden     bool                          `json:"is_hidden"`
	Status       int                           `json:"status"`        // 审核状态（0:待审核 1:已通过 2:已拒绝）
	Email        *string                       `json:"email,omitempty"`
	UserID       *xSnowflake.SnowflakeID       `json:"user_id,omitempty"`
	ApplyRemark  *string                       `json:"apply_remark,omitempty"`
	ReviewRemark *string                       `json:"review_remark,omitempty"`
	CreatedAt    time.Time                     `json:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at"`
	Channel      *SponsorChannelSimpleResponse `json:"channel,omitempty"`
}

// RecordPublicItemResponse 赞助记录公开条目响应
type RecordPublicItemResponse struct {
	ID          xSnowflake.SnowflakeID        `json:"id"`
	Nickname    string                        `json:"nickname"`
	RedirectURL *string                       `json:"redirect_url"`
	Amount      int64                         `json:"amount"`
	Message     *string                       `json:"message"`
	SponsorAt   *time.Time                    `json:"sponsor_at"`
	Channel     *SponsorChannelSimpleResponse `json:"channel,omitempty"`
}

// RecordAddResponse 添加记录响应
type RecordAddResponse struct {
	RecordEntityResponse
}

// RecordUpdateResponse 更新记录响应
type RecordUpdateResponse struct {
	RecordEntityResponse
}

// RecordDetailResponse 详情响应
type RecordDetailResponse struct {
	RecordEntityResponse
}

// RecordPageResponse 分页响应（后台）
type RecordPageResponse struct {
	base.PaginationResponse[RecordEntityResponse]
}

// RecordAdminPageResponse 分页响应（后台，附带待审核计数供管理入口徽章展示）
type RecordAdminPageResponse struct {
	base.PaginationResponse[RecordEntityResponse]
	PendingCount int64 `json:"pending_count"` // 待审核赞助数量
}

// RecordPublicPageResponse 公开分页响应（前台）
type RecordPublicPageResponse struct {
	base.PaginationResponse[RecordPublicItemResponse]
}
