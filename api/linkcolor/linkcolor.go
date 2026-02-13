package apiLinkColor

import (
	"github.com/bamboo-services/bamboo-main/internal/model/base"
	"github.com/bamboo-services/bamboo-main/internal/model/dto"
)

// AddRequest 添加友链颜色请求
type AddRequest struct {
	ColorName    string  `json:"color_name" binding:"required,min=1,max=50" example:"炫彩红"`     // 颜色名称
	ColorType    int     `json:"color_type" binding:"oneof=0 1" example:"0"`                   // 颜色类型（0: 普通, 1: 炫彩）
	PrimaryColor *string `json:"primary_color" binding:"omitempty,hexcolor" example:"#FF0000"` // 主颜色
	SubColor     *string `json:"sub_color" binding:"omitempty,hexcolor" example:"#FF6600"`     // 副颜色
	HoverColor   *string `json:"hover_color" binding:"omitempty,hexcolor" example:"#FF3300"`   // 悬停颜色
	ColorOrder   int     `json:"color_order" binding:"omitempty,min=0" example:"0"`            // 颜色排序
}

// UpdateRequest 更新友链颜色请求
type UpdateRequest struct {
	ColorName    *string `json:"color_name" binding:"omitempty,min=1,max=50" example:"炫彩红"` // 颜色名称
	ColorType    *int    `json:"color_type" binding:"omitempty,oneof=0 1" example:"0"`      // 颜色类型（0: 普通, 1: 炫彩）
	PrimaryColor *string `json:"primary_color" binding:"omitempty" example:"#FF0000"`       // 主颜色（可传空字符串清空）
	SubColor     *string `json:"sub_color" binding:"omitempty" example:"#FF6600"`           // 副颜色（可传空字符串清空）
	HoverColor   *string `json:"hover_color" binding:"omitempty" example:"#FF3300"`         // 悬停颜色（可传空字符串清空）
	ColorOrder   *int    `json:"color_order" binding:"omitempty,min=0" example:"0"`         // 颜色排序
}

// SortRequest 颜色排序请求
type SortRequest struct {
	ColorIDs  []int64 `json:"color_ids" binding:"required,min=1" validate:"required"` // 颜色ID数组，按新的排序传入
	SortOrder *int    `json:"sort_order" binding:"omitempty,min=0" example:"10"`      // 可选的起始排序值，不填则从0开始递增
}

// StatusRequest 颜色状态切换请求
type StatusRequest struct {
	Status bool `json:"status" binding:"omitempty" example:"true"` // 状态：true=启用，false=禁用
}

// DeleteRequest 删除颜色请求
type DeleteRequest struct {
	Force bool `json:"force" binding:"omitempty" example:"false"` // 是否强制删除：false=检查关联，true=强制删除并清空关联
}

// ListRequest 颜色列表查询请求
type ListRequest struct {
	Status      *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Type        *int    `form:"type" binding:"omitempty,oneof=0 1" example:"0"`                                     // 类型过滤：0=普通，1=炫彩，不传=全部
	Name        *string `form:"name" binding:"omitempty,max=50" example:"红"`                                        // 名称模糊搜索
	OnlyEnabled *bool   `form:"only_enabled" binding:"omitempty" example:"true"`                                    // 仅查询启用的颜色
	OrderBy     *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// PageRequest 颜色分页查询请求
type PageRequest struct {
	Page     int     `form:"page" binding:"omitempty,min=1" validate:"min=1" example:"1"`                        // 页码，默认1
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100" validate:"min=1,max=100" example:"10"`  // 每页数量，默认10，最大100
	Status   *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Type     *int    `form:"type" binding:"omitempty,oneof=0 1" example:"0"`                                     // 类型过滤：0=普通，1=炫彩，不传=全部
	Name     *string `form:"name" binding:"omitempty,max=50" example:"红"`                                        // 名称模糊搜索
	OrderBy  *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// AddResponse 添加友链颜色响应
type AddResponse struct {
	dto.LinkColorDetailDTO
}

// UpdateResponse 更新友链颜色响应
type UpdateResponse struct {
	dto.LinkColorDetailDTO
}

// DetailResponse 友链颜色详情响应
type DetailResponse struct {
	dto.LinkColorDetailDTO
}

// ListResponse 友链颜色列表响应
type ListResponse struct {
	Colors []dto.LinkColorListDTO `json:"colors"`
}

// PageResponse 友链颜色分页响应
type PageResponse struct {
	base.PaginationResponse[dto.LinkColorNormalDTO]
}

// SortResponse 友链颜色排序响应
type SortResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"` // 更新的颜色数量
}

// StatusResponse 友链颜色状态切换响应
type StatusResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"` // 更新后的状态
}

// DeleteResponse 友链颜色删除响应
type DeleteResponse struct {
	Message string `json:"message"`
}

// DeleteConflictResponse 友链颜色删除冲突响应
type DeleteConflictResponse struct {
	Message      string                           `json:"message"`       // 错误消息
	ConflictInfo DeleteConflictInfo               `json:"conflict_info"` // 冲突信息
	Links        []dto.LinkColorDeleteConflictDTO `json:"links"`         // 冲突的友链列表（前10个）
}

// DeleteConflictInfo 删除冲突的详细信息
type DeleteConflictInfo struct {
	ColorID      int64  `json:"color_id"`      // 颜色ID
	ColorName    string `json:"color_name"`    // 颜色名称
	TotalLinks   int    `json:"total_links"`   // 总友链数量
	DisplayCount int    `json:"display_count"` // 显示的友链数量（最多10个）
}
