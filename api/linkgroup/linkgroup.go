package apiLinkGroup

import (
	"github.com/bamboo-services/bamboo-main/internal/model/base"
	"github.com/bamboo-services/bamboo-main/internal/model/dto"
)

// AddRequest 添加友链分组请求
type AddRequest struct {
	GroupName  string `json:"group_name" binding:"required,min=1,max=100" example:"技术博客"`
	GroupDesc  string `json:"group_desc" binding:"omitempty,max=500" example:"技术相关的友情链接"`
	GroupOrder int    `json:"group_order" binding:"omitempty,min=0" example:"0"`
}

// UpdateRequest 更新友链分组请求
type UpdateRequest struct {
	GroupName   string `json:"group_name" binding:"omitempty,min=1,max=100" example:"技术博客"`
	GroupDesc   string `json:"group_desc" binding:"omitempty,max=500" example:"技术相关的友情链接"`
	GroupOrder  *int   `json:"group_order" binding:"omitempty,min=0" example:"0"`
	GroupStatus *int   `json:"group_status" binding:"omitempty,oneof=0 1" example:"1"`
}

// SortRequest 分组排序请求
type SortRequest struct {
	GroupIDs  []int64 `json:"group_ids" binding:"required,min=1" validate:"required"` // 分组ID数组，按新的排序传入
	SortOrder *int    `json:"sort_order" binding:"omitempty,min=0" example:"10"`      // 可选的起始排序值，不填则从0开始递增
}

// StatusRequest 分组状态切换请求
type StatusRequest struct {
	Status bool `json:"status" binding:"omitempty" example:"true"` // 状态：true=启用，false=禁用
}

// DeleteRequest 删除分组请求
type DeleteRequest struct {
	Force bool `json:"force" binding:"omitempty" example:"false"` // 是否强制删除：false=检查关联，true=强制删除并清空关联
}

// ListRequest 分组列表查询请求
type ListRequest struct {
	Status      *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Name        *string `form:"name" binding:"omitempty,max=100" example:"技术"`                                      // 名称模糊搜索
	WithLinks   *bool   `form:"with_links" binding:"omitempty" example:"false"`                                     // 是否包含友链列表
	OnlyEnabled *bool   `form:"only_enabled" binding:"omitempty" example:"true"`                                    // 仅查询启用的分组
	OrderBy     *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// PageRequest 分组分页查询请求
type PageRequest struct {
	Page     int     `form:"page" binding:"omitempty,min=1" validate:"min=1" example:"1"`                        // 页码，默认1
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100" validate:"min=1,max=100" example:"10"`  // 每页数量，默认10，最大100
	Status   *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Name     *string `form:"name" binding:"omitempty,max=100" example:"技术"`                                      // 名称模糊搜索
	OrderBy  *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// AddResponse 添加友链分组响应
type AddResponse struct {
	dto.LinkGroupDetailDTO
}

// UpdateResponse 更新友链分组响应
type UpdateResponse struct {
	dto.LinkGroupDetailDTO
}

// DetailResponse 友链分组详情响应
type DetailResponse struct {
	dto.LinkGroupDetailDTO
}

// ListResponse 友链分组列表响应
type ListResponse struct {
	Groups []dto.LinkGroupListDTO `json:"groups"`
}

// PageResponse 友链分组分页响应
type PageResponse struct {
	base.PaginationResponse[dto.LinkGroupNormalDTO]
}

// SortResponse 友链分组排序响应
type SortResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"` // 更新的分组数量
}

// StatusResponse 友链分组状态切换响应
type StatusResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"` // 更新后的状态
}

// DeleteResponse 友链分组删除响应
type DeleteResponse struct {
	Message string `json:"message"`
}

// DeleteConflictResponse 友链分组删除冲突响应
type DeleteConflictResponse struct {
	Message      string                           `json:"message"`       // 错误消息
	ConflictInfo DeleteConflictInfo               `json:"conflict_info"` // 冲突信息
	Links        []dto.LinkGroupDeleteConflictDTO `json:"links"`         // 冲突的友链列表（前10个）
}

// DeleteConflictInfo 删除冲突的详细信息
type DeleteConflictInfo struct {
	GroupUUID    string `json:"group_uuid"`    // 分组UUID
	GroupName    string `json:"group_name"`    // 分组名称
	TotalLinks   int    `json:"total_links"`   // 总友链数量
	DisplayCount int    `json:"display_count"` // 显示的友链数量（最多10个）
}

// 为了保持一致性，定义一些类型别名（引用现有的MessageResponse）
type (
	// DeleteSimpleResponse 简单删除响应（成功时）
	DeleteSimpleResponse struct {
		Message string `json:"message"`
	}

	// StatusSimpleResponse 简单状态切换响应（不需要返回详细状态时）
	StatusSimpleResponse struct {
		Message string `json:"message"`
	}
)
