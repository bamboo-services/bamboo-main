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

package apiLink

import (
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/models/dto"
)

// ColorAddRequest 添加友链颜色请求
type ColorAddRequest struct {
	ColorName    string  `json:"color_name" binding:"required,min=1,max=50" example:"炫彩红"`     // 颜色名称
	ColorType    int     `json:"color_type" binding:"oneof=0 1" example:"0"`                   // 颜色类型（0: 普通, 1: 炫彩）
	PrimaryColor *string `json:"primary_color" binding:"omitempty,hexcolor" example:"#FF0000"` // 主颜色
	SubColor     *string `json:"sub_color" binding:"omitempty,hexcolor" example:"#FF6600"`     // 副颜色
	HoverColor   *string `json:"hover_color" binding:"omitempty,hexcolor" example:"#FF3300"`   // 悬停颜色
	ColorOrder   int     `json:"color_order" binding:"omitempty,min=0" example:"0"`            // 颜色排序
}

// ColorUpdateRequest 更新友链颜色请求
type ColorUpdateRequest struct {
	ColorName    *string `json:"color_name" binding:"omitempty,min=1,max=50" example:"炫彩红"` // 颜色名称
	ColorType    *int    `json:"color_type" binding:"omitempty,oneof=0 1" example:"0"`      // 颜色类型（0: 普通, 1: 炫彩）
	PrimaryColor *string `json:"primary_color" binding:"omitempty" example:"#FF0000"`       // 主颜色（可传空字符串清空）
	SubColor     *string `json:"sub_color" binding:"omitempty" example:"#FF6600"`           // 副颜色（可传空字符串清空）
	HoverColor   *string `json:"hover_color" binding:"omitempty" example:"#FF3300"`         // 悬停颜色（可传空字符串清空）
	ColorOrder   *int    `json:"color_order" binding:"omitempty,min=0" example:"0"`         // 颜色排序
}

// ColorSortRequest 颜色排序请求
type ColorSortRequest struct {
	ColorIDs  []int64 `json:"color_ids" binding:"required,min=1" validate:"required"` // 颜色ID数组，按新的排序传入
	SortOrder *int    `json:"sort_order" binding:"omitempty,min=0" example:"10"`      // 可选的起始排序值，不填则从0开始递增
}

// ColorStatusRequest 颜色状态切换请求
type ColorStatusRequest struct {
	Status bool `json:"status" binding:"omitempty" example:"true"` // 状态：true=启用，false=禁用
}

// ColorDeleteRequest 删除颜色请求
type ColorDeleteRequest struct {
	Force bool `json:"force" binding:"omitempty" example:"false"` // 是否强制删除：false=检查关联，true=强制删除并清空关联
}

// ColorListRequest 颜色列表查询请求
type ColorListRequest struct {
	Status      *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Type        *int    `form:"type" binding:"omitempty,oneof=0 1" example:"0"`                                     // 类型过滤：0=普通，1=炫彩，不传=全部
	Name        *string `form:"name" binding:"omitempty,max=50" example:"红"`                                        // 名称模糊搜索
	OnlyEnabled *bool   `form:"only_enabled" binding:"omitempty" example:"true"`                                    // 仅查询启用的颜色
	OrderBy     *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// ColorPageRequest 颜色分页查询请求
type ColorPageRequest struct {
	Page     int     `form:"page" binding:"omitempty,min=1" validate:"min=1" example:"1"`                        // 页码，默认1
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100" validate:"min=1,max=100" example:"10"`  // 每页数量，默认10，最大100
	Status   *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Type     *int    `form:"type" binding:"omitempty,oneof=0 1" example:"0"`                                     // 类型过滤：0=普通，1=炫彩，不传=全部
	Name     *string `form:"name" binding:"omitempty,max=50" example:"红"`                                        // 名称模糊搜索
	OrderBy  *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// ColorAddResponse 添加友链颜色响应
type ColorAddResponse struct {
	dto.LinkColorDetailDTO
}

// ColorUpdateResponse 更新友链颜色响应
type ColorUpdateResponse struct {
	dto.LinkColorDetailDTO
}

// ColorDetailResponse 友链颜色详情响应
type ColorDetailResponse struct {
	dto.LinkColorDetailDTO
}

// ColorListResponse 友链颜色列表响应
type ColorListResponse struct {
	Colors []dto.LinkColorListDTO `json:"colors"`
}

// ColorPageResponse 友链颜色分页响应
type ColorPageResponse struct {
	base.PaginationResponse[dto.LinkColorNormalDTO]
}

// ColorSortResponse 友链颜色排序响应
type ColorSortResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"` // 更新的颜色数量
}

// ColorStatusResponse 友链颜色状态切换响应
type ColorStatusResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"` // 更新后的状态
}

// ColorDeleteResponse 友链颜色删除响应
type ColorDeleteResponse struct {
	Message string `json:"message"`
}

// ColorDeleteConflictResponse 友链颜色删除冲突响应
type ColorDeleteConflictResponse struct {
	Message      string                           `json:"message"`       // 错误消息
	ConflictInfo ColorDeleteConflictInfo          `json:"conflict_info"` // 冲突信息
	Links        []dto.LinkColorDeleteConflictDTO `json:"links"`         // 冲突的友链列表（前10个）
}

// ColorDeleteConflictInfo 删除冲突的详细信息
type ColorDeleteConflictInfo struct {
	ColorID      int64  `json:"color_id"`      // 颜色ID
	ColorName    string `json:"color_name"`    // 颜色名称
	TotalLinks   int    `json:"total_links"`   // 总友链数量
	DisplayCount int    `json:"display_count"` // 显示的友链数量（最多10个）
}
