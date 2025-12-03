/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明:版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息,请查看项目根目录下的LICENSE文件或访问:
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package request

// SponsorChannelAddReq 添加赞助渠道请求
type SponsorChannelAddReq struct {
	Name        string  `json:"name" binding:"required,min=1,max=50" example:"支付宝"`                      // 渠道名称
	Icon        *string `json:"icon" binding:"omitempty,max=500" example:"https://example.com/icon.png"` // 渠道图标 URL
	Description *string `json:"description" binding:"omitempty,max=1000" example:"支付宝赞助渠道"`              // 渠道描述
	SortOrder   int     `json:"sort_order" binding:"omitempty,min=0" example:"0"`                        // 排序值
}

// SponsorChannelUpdateReq 更新赞助渠道请求
type SponsorChannelUpdateReq struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=50" example:"支付宝"`                     // 渠道名称
	Icon        *string `json:"icon" binding:"omitempty,max=500" example:"https://example.com/icon.png"` // 渠道图标 URL
	Description *string `json:"description" binding:"omitempty,max=1000" example:""`                     // 渠道描述
	SortOrder   *int    `json:"sort_order" binding:"omitempty,min=0" example:"0"`                        // 排序值
}

// SponsorChannelStatusReq 状态切换请求
type SponsorChannelStatusReq struct {
	Status bool `json:"status" example:"true"` // 状态:true=启用,false=禁用
}

// SponsorChannelListReq 列表查询请求(不分页)
type SponsorChannelListReq struct {
	Status      *bool   `form:"status" binding:"omitempty" example:"true"`                                          // 状态过滤:true=启用,false=禁用,不传=全部
	Name        *string `form:"name" binding:"omitempty,max=50" example:"支付"`                                       // 名称模糊搜索
	OnlyEnabled *bool   `form:"only_enabled" binding:"omitempty" example:"true"`                                    // 仅查询启用的渠道
	OrderBy     *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// SponsorChannelPageReq 分页查询请求
type SponsorChannelPageReq struct {
	Page     int     `form:"page" binding:"omitempty,min=1" example:"1"`                                         // 页码,默认1
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`                           // 每页数量,默认10,最大100
	Status   *bool   `form:"status" binding:"omitempty" example:"true"`                                          // 状态过滤:true=启用,false=禁用,不传=全部
	Name     *string `form:"name" binding:"omitempty,max=50" example:"支付"`                                       // 名称模糊搜索
	OrderBy  *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}
