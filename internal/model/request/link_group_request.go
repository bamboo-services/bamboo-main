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

package request

// LinkGroupAddReq 添加友链分组请求
type LinkGroupAddReq struct {
	GroupName  string `json:"group_name" binding:"required,min=1,max=100" example:"技术博客"`
	GroupDesc  string `json:"group_desc" binding:"omitempty,max=500" example:"技术相关的友情链接"`
	GroupOrder int    `json:"group_order" binding:"omitempty,min=0" example:"0"`
}

// LinkGroupUpdateReq 更新友链分组请求
type LinkGroupUpdateReq struct {
	GroupName   string `json:"group_name" binding:"omitempty,min=1,max=100" example:"技术博客"`
	GroupDesc   string `json:"group_desc" binding:"omitempty,max=500" example:"技术相关的友情链接"`
	GroupOrder  *int   `json:"group_order" binding:"omitempty,min=0" example:"0"`
	GroupStatus *int   `json:"group_status" binding:"omitempty,oneof=0 1" example:"1"`
}

// LinkGroupSortReq 分组排序请求
type LinkGroupSortReq struct {
	GroupIDs  []int64 `json:"group_ids" binding:"required,min=1" validate:"required"` // 分组ID数组，按新的排序传入
	SortOrder *int    `json:"sort_order" binding:"omitempty,min=0" example:"10"`      // 可选的起始排序值，不填则从0开始递增
}

// LinkGroupStatusReq 分组状态切换请求
type LinkGroupStatusReq struct {
	Status bool `json:"status" binding:"omitempty" example:"true"` // 状态：true=启用，false=禁用
}

// LinkGroupDeleteReq 删除分组请求
type LinkGroupDeleteReq struct {
	Force bool `json:"force" binding:"omitempty" example:"false"` // 是否强制删除：false=检查关联，true=强制删除并清空关联
}

// LinkGroupListReq 分组列表查询请求
type LinkGroupListReq struct {
	Status      *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Name        *string `form:"name" binding:"omitempty,max=100" example:"技术"`                                      // 名称模糊搜索
	WithLinks   *bool   `form:"with_links" binding:"omitempty" example:"false"`                                     // 是否包含友链列表
	OnlyEnabled *bool   `form:"only_enabled" binding:"omitempty" example:"true"`                                    // 仅查询启用的分组
	OrderBy     *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order       *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}

// LinkGroupPageReq 分组分页查询请求
type LinkGroupPageReq struct {
	Page     int     `form:"page" binding:"omitempty,min=1" validate:"min=1" example:"1"`                        // 页码，默认1
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100" validate:"min=1,max=100" example:"10"`  // 每页数量，默认10，最大100
	Status   *int    `form:"status" binding:"omitempty,oneof=0 1" example:"1"`                                   // 状态过滤：0=禁用，1=启用，不传=全部
	Name     *string `form:"name" binding:"omitempty,max=100" example:"技术"`                                      // 名称模糊搜索
	OrderBy  *string `form:"order_by" binding:"omitempty,oneof=name sort_order created_at" example:"sort_order"` // 排序字段
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc" example:"asc"`                             // 排序方向
}
