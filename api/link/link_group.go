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

package apiLink

import (
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
)

// GroupIDRequest 友链分组路径参数
type GroupIDRequest struct {
	ID xSnowflake.SnowflakeID `uri:"id" binding:"required"`
}

// GroupAddRequest 添加友链分组请求
type GroupAddRequest struct {
	GroupName  string `json:"group_name" binding:"required,min=1,max=100" example:"技术博客"`
	GroupDesc  string `json:"group_desc" binding:"omitempty,max=500" example:"技术相关的友情链接"`
	GroupOrder int    `json:"group_order" binding:"omitempty,min=0" example:"0"`
}

// GroupUpdateRequest 更新友链分组请求
type GroupUpdateRequest struct {
	GroupName   string `json:"group_name" binding:"omitempty,min=1,max=100" example:"技术博客"`
	GroupDesc   string `json:"group_desc" binding:"omitempty,max=500" example:"技术相关的友情链接"`
	GroupOrder  *int   `json:"group_order" binding:"omitempty,min=0" example:"0"`
	GroupStatus *int   `json:"group_status" binding:"omitempty,oneof=0 1" example:"1"`
}

// GroupSortRequest 分组排序请求
type GroupSortRequest struct {
	GroupIDs  []xSnowflake.SnowflakeID `json:"group_ids" binding:"required,min=1" validate:"required"` // 分组ID数组，按新的排序传入
	SortOrder *int                     `json:"sort_order" binding:"omitempty,min=0" example:"10"`      // 可选的起始排序值，不填则从0开始递增
}

// GroupStatusRequest 分组状态切换请求
type GroupStatusRequest struct {
	Status bool `json:"status" binding:"omitempty" example:"true"` // 状态：true=启用，false=禁用
}

// GroupDeleteRequest 删除分组请求
type GroupDeleteRequest struct {
	Force bool `json:"force" binding:"omitempty" example:"false"` // 是否强制删除：false=检查关联，true=强制删除并清空关联
}

// GroupListRequest 分组列表查询请求
type GroupListRequest struct {
	Status      *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Name        *string `form:"name" binding:"omitempty,max=100" example:"技术"`                                      // 名称模糊搜索
	WithLinks   *bool   `form:"with_links" binding:"omitempty" example:"false"`                                     // 是否包含友链列表
	OnlyEnabled *bool   `form:"only_enabled" binding:"omitempty" example:"true"`                                    // 仅查询启用的分组
	OrderBy     *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// GroupPageRequest 分组分页查询请求
type GroupPageRequest struct {
	Page     int     `form:"page" binding:"omitempty,min=1" validate:"min=1" example:"1"`                        // 页码，默认1
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100" validate:"min=1,max=100" example:"10"`  // 每页数量，默认10，最大100
	Status   *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Name     *string `form:"name" binding:"omitempty,max=100" example:"技术"`                                      // 名称模糊搜索
	OrderBy  *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// GroupAddResponse 添加友链分组响应
type GroupAddResponse struct {
	entity.LinkGroup
}

// GroupUpdateResponse 更新友链分组响应
type GroupUpdateResponse struct {
	entity.LinkGroup
}

// GroupDetailResponse 友链分组详情响应
type GroupDetailResponse struct {
	entity.LinkGroup
}

// GroupListResponse 友链分组列表响应
type GroupListResponse struct {
	Groups []entity.LinkGroup `json:"groups"`
}

// GroupPageResponse 友链分组分页响应
type GroupPageResponse struct {
	base.PaginationResponse[entity.LinkGroup]
}

// GroupSortResponse 友链分组排序响应
type GroupSortResponse struct {
	Count int `json:"count"` // 更新的分组数量
}

// GroupStatusResponse 友链分组状态切换响应
type GroupStatusResponse struct {
	Status bool `json:"status"` // 更新后的状态
}

// GroupDeleteConflictResponse 友链分组删除冲突响应
type GroupDeleteConflictResponse struct {
	Message      string                  `json:"message"`       // 错误消息
	ConflictInfo GroupDeleteConflictInfo `json:"conflict_info"` // 冲突信息
	Links        []entity.LinkFriend     `json:"links"`         // 冲突的友链列表（前10个）
}

// GroupDeleteConflictInfo 删除冲突的详细信息
type GroupDeleteConflictInfo struct {
	GroupID      xSnowflake.SnowflakeID `json:"group_id"`
	GroupName    string                 `json:"group_name"`    // 分组名称
	TotalLinks   int                    `json:"total_links"`   // 总友链数量
	DisplayCount int                    `json:"display_count"` // 显示的友链数量（最多10个）
}

// BuiltinInvalidGroupUpdateRequest 更新内置「已失效」分组配置请求
//
// PATCH 语义：仅更新非 nil 字段；名称必填（去除空白后不能为空），描述传空串即清空。
type BuiltinInvalidGroupUpdateRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100" example:"已失效"`           // 分组名称
	Description *string `json:"description" binding:"omitempty,max=500" example:"已失效友链，站点不可访问"` // 分组描述（空串清空）
}

// BuiltinInvalidGroupResponse 内置「已失效」分组配置响应
type BuiltinInvalidGroupResponse struct {
	entity.LinkGroup
}
