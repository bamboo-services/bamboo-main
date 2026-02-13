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

package apiSponsor

import (
	"time"

	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/models/dto"
)

// RecordAddRequest 添加赞助记录请求
type RecordAddRequest struct {
	Nickname    string     `json:"nickname" binding:"required,min=1,max=100" example:"张三"`                 // 赞助者昵称
	RedirectURL *string    `json:"redirect_url" binding:"omitempty,max=500" example:"https://example.com"` // 跳转链接
	Amount      int64      `json:"amount" binding:"required,min=1" example:"1000"`                         // 赞助金额(分)
	ChannelID   *int64     `json:"channel_id" binding:"omitempty" example:"123456789"`                     // 赞助渠道 ID
	Message     *string    `json:"message" binding:"omitempty,max=500" example:"感谢开发者"`                    // 留言信息
	SponsorAt   *time.Time `json:"sponsor_at" binding:"omitempty" example:"2025-01-01T12:00:00Z"`          // 赞助时间
	SortOrder   int        `json:"sort_order" binding:"omitempty,min=0" example:"0"`                       // 排序值
	IsAnonymous bool       `json:"is_anonymous" binding:"omitempty" example:"false"`                       // 是否匿名
	IsHidden    bool       `json:"is_hidden" binding:"omitempty" example:"false"`                          // 是否隐藏
}

// RecordUpdateRequest 更新赞助记录请求
type RecordUpdateRequest struct {
	Nickname    *string    `json:"nickname" binding:"omitempty,min=1,max=100" example:"张三"`                // 赞助者昵称
	RedirectURL *string    `json:"redirect_url" binding:"omitempty,max=500" example:"https://example.com"` // 跳转链接
	Amount      *int64     `json:"amount" binding:"omitempty,min=1" example:"1000"`                        // 赞助金额(分)
	ChannelID   *int64     `json:"channel_id" binding:"omitempty" example:"123456789"`                     // 赞助渠道 ID
	Message     *string    `json:"message" binding:"omitempty,max=500" example:""`                         // 留言信息
	SponsorAt   *time.Time `json:"sponsor_at" binding:"omitempty" example:""`                              // 赞助时间
	SortOrder   *int       `json:"sort_order" binding:"omitempty,min=0" example:"0"`                       // 排序值
	IsAnonymous *bool      `json:"is_anonymous" binding:"omitempty" example:"false"`                       // 是否匿名
	IsHidden    *bool      `json:"is_hidden" binding:"omitempty" example:"false"`                          // 是否隐藏
}

// RecordPageRequest 分页查询请求(后台)
type RecordPageRequest struct {
	Page        int     `form:"page" binding:"omitempty,min=1" example:"1"`                                                               // 页码,默认1
	PageSize    int     `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`                                                 // 每页数量,默认10,最大100
	ChannelID   *int64  `form:"channel_id" binding:"omitempty" example:"123456789"`                                                       // 渠道 ID 过滤
	Nickname    *string `form:"nickname" binding:"omitempty,max=100" example:"张"`                                                         // 昵称模糊搜索
	IsAnonymous *bool   `form:"is_anonymous" binding:"omitempty" example:"false"`                                                         // 是否匿名过滤
	IsHidden    *bool   `form:"is_hidden" binding:"omitempty" example:"false"`                                                            // 是否隐藏过滤
	OrderBy     *string `form:"order_by" binding:"omitempty,oneof=nickname amount sponsor_at sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc" example:"desc"`                                                  // 排序方向
}

// RecordPublicPageRequest 公开分页查询请求(前台)
type RecordPublicPageRequest struct {
	Page      int     `form:"page" binding:"omitempty,min=1" example:"1"`                                           // 页码,默认1
	PageSize  int     `form:"page_size" binding:"omitempty,min=1,max=50" example:"20"`                              // 每页数量,默认20,最大50
	ChannelID *int64  `form:"channel_id" binding:"omitempty" example:"123456789"`                                   // 渠道 ID 过滤
	OrderBy   *string `form:"order_by" binding:"omitempty,oneof=amount sponsor_at sort_order" example:"sort_order"` // 排序字段
	Order     *string `form:"order" binding:"omitempty,oneof=asc desc" example:"desc"`                              // 排序方向
}

// RecordAddResponse 添加记录响应
type RecordAddResponse struct {
	dto.SponsorRecordDetailDTO
}

// RecordUpdateResponse 更新记录响应
type RecordUpdateResponse struct {
	dto.SponsorRecordDetailDTO
}

// RecordDetailResponse 详情响应
type RecordDetailResponse struct {
	dto.SponsorRecordDetailDTO
}

// RecordPageResponse 分页响应（后台）
type RecordPageResponse struct {
	base.PaginationResponse[dto.SponsorRecordNormalDTO]
}

// RecordPublicPageResponse 公开分页响应（前台）
type RecordPublicPageResponse struct {
	base.PaginationResponse[dto.SponsorRecordSimpleDTO]
}

// RecordDeleteResponse 删除响应
type RecordDeleteResponse struct {
	Message string `json:"message"`
}
