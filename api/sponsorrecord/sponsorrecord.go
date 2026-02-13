package apiSponsorRecord

import (
	"time"

	"github.com/bamboo-services/bamboo-main/internal/model/base"
	"github.com/bamboo-services/bamboo-main/internal/model/dto"
)

// AddRequest 添加赞助记录请求
type AddRequest struct {
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

// UpdateRequest 更新赞助记录请求
type UpdateRequest struct {
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

// PageRequest 分页查询请求(后台)
type PageRequest struct {
	Page        int     `form:"page" binding:"omitempty,min=1" example:"1"`                                                               // 页码,默认1
	PageSize    int     `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`                                                 // 每页数量,默认10,最大100
	ChannelID   *int64  `form:"channel_id" binding:"omitempty" example:"123456789"`                                                       // 渠道 ID 过滤
	Nickname    *string `form:"nickname" binding:"omitempty,max=100" example:"张"`                                                         // 昵称模糊搜索
	IsAnonymous *bool   `form:"is_anonymous" binding:"omitempty" example:"false"`                                                         // 是否匿名过滤
	IsHidden    *bool   `form:"is_hidden" binding:"omitempty" example:"false"`                                                            // 是否隐藏过滤
	OrderBy     *string `form:"order_by" binding:"omitempty,oneof=nickname amount sponsor_at sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc" example:"desc"`                                                  // 排序方向
}

// PublicPageRequest 公开分页查询请求(前台)
type PublicPageRequest struct {
	Page      int     `form:"page" binding:"omitempty,min=1" example:"1"`                                           // 页码,默认1
	PageSize  int     `form:"page_size" binding:"omitempty,min=1,max=50" example:"20"`                              // 每页数量,默认20,最大50
	ChannelID *int64  `form:"channel_id" binding:"omitempty" example:"123456789"`                                   // 渠道 ID 过滤
	OrderBy   *string `form:"order_by" binding:"omitempty,oneof=amount sponsor_at sort_order" example:"sort_order"` // 排序字段
	Order     *string `form:"order" binding:"omitempty,oneof=asc desc" example:"desc"`                              // 排序方向
}

// AddResponse 添加记录响应
type AddResponse struct {
	dto.SponsorRecordDetailDTO
}

// UpdateResponse 更新记录响应
type UpdateResponse struct {
	dto.SponsorRecordDetailDTO
}

// DetailResponse 详情响应
type DetailResponse struct {
	dto.SponsorRecordDetailDTO
}

// PageResponse 分页响应（后台）
type PageResponse struct {
	base.PaginationResponse[dto.SponsorRecordNormalDTO]
}

// PublicPageResponse 公开分页响应（前台）
type PublicPageResponse struct {
	base.PaginationResponse[dto.SponsorRecordSimpleDTO]
}

// DeleteResponse 删除响应
type DeleteResponse struct {
	Message string `json:"message"`
}
