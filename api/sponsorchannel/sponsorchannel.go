package apiSponsorChannel

import (
	"github.com/bamboo-services/bamboo-main/internal/model/base"
	"github.com/bamboo-services/bamboo-main/internal/model/dto"
)

// AddRequest 添加赞助渠道请求
type AddRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=50" example:"支付宝"`                      // 渠道名称
	Icon        *string `json:"icon" binding:"omitempty,max=500" example:"https://example.com/icon.png"` // 渠道图标 URL
	Description *string `json:"description" binding:"omitempty,max=1000" example:"支付宝赞助渠道"`              // 渠道描述
	SortOrder   int     `json:"sort_order" binding:"omitempty,min=0" example:"0"`                        // 排序值
}

// UpdateRequest 更新赞助渠道请求
type UpdateRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=50" example:"支付宝"`                     // 渠道名称
	Icon        *string `json:"icon" binding:"omitempty,max=500" example:"https://example.com/icon.png"` // 渠道图标 URL
	Description *string `json:"description" binding:"omitempty,max=1000" example:""`                     // 渠道描述
	SortOrder   *int    `json:"sort_order" binding:"omitempty,min=0" example:"0"`                        // 排序值
}

// StatusRequest 状态切换请求
type StatusRequest struct {
	Status bool `json:"status" example:"true"` // 状态:true=启用,false=禁用
}

// ListRequest 列表查询请求(不分页)
type ListRequest struct {
	Status      *bool   `form:"status" binding:"omitempty" example:"true"`                                          // 状态过滤:true=启用,false=禁用,不传=全部
	Name        *string `form:"name" binding:"omitempty,max=50" example:"支付"`                                       // 名称模糊搜索
	OnlyEnabled *bool   `form:"only_enabled" binding:"omitempty" example:"true"`                                    // 仅查询启用的渠道
	OrderBy     *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// PageRequest 分页查询请求
type PageRequest struct {
	Page     int     `form:"page" binding:"omitempty,min=1" example:"1"`                                         // 页码,默认1
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`                           // 每页数量,默认10,最大100
	Status   *bool   `form:"status" binding:"omitempty" example:"true"`                                          // 状态过滤:true=启用,false=禁用,不传=全部
	Name     *string `form:"name" binding:"omitempty,max=50" example:"支付"`                                       // 名称模糊搜索
	OrderBy  *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// AddResponse 添加渠道响应
type AddResponse struct {
	dto.SponsorChannelDetailDTO
}

// UpdateResponse 更新渠道响应
type UpdateResponse struct {
	dto.SponsorChannelDetailDTO
}

// DetailResponse 详情响应
type DetailResponse struct {
	dto.SponsorChannelDetailDTO
}

// ListResponse 列表响应（不分页）
type ListResponse struct {
	Channels []dto.SponsorChannelListDTO `json:"channels"`
}

// PageResponse 分页响应
type PageResponse struct {
	base.PaginationResponse[dto.SponsorChannelNormalDTO]
}

// StatusResponse 状态切换响应
type StatusResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"` // 更新后的状态
}

// DeleteResponse 删除响应
type DeleteResponse struct {
	Message string `json:"message"`
}

// DeleteConflictResponse 删除冲突响应
type DeleteConflictResponse struct {
	Message      string                                `json:"message"`       // 错误消息
	ConflictInfo DeleteConflictInfo                    `json:"conflict_info"` // 冲突信息
	Sponsors     []dto.SponsorChannelDeleteConflictDTO `json:"sponsors"`      // 冲突的赞助记录列表（前10个）
}

// DeleteConflictInfo 删除冲突的详细信息
type DeleteConflictInfo struct {
	ChannelID     int64  `json:"channel_id"`     // 渠道ID
	ChannelName   string `json:"channel_name"`   // 渠道名称
	TotalSponsors int    `json:"total_sponsors"` // 总赞助记录数量
	DisplayCount  int    `json:"display_count"`  // 显示的赞助记录数量（最多10个）
}
