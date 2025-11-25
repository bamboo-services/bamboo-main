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

package response

import (
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
)

// LinkGroupAddResponse 添加友链分组响应
type LinkGroupAddResponse struct {
	dto.LinkGroupDetailDTO
}

// LinkGroupUpdateResponse 更新友链分组响应
type LinkGroupUpdateResponse struct {
	dto.LinkGroupDetailDTO
}

// LinkGroupDetailResponse 友链分组详情响应
type LinkGroupDetailResponse struct {
	dto.LinkGroupDetailDTO
}

// LinkGroupListResponse 友链分组列表响应
type LinkGroupListResponse struct {
	Groups []dto.LinkGroupListDTO `json:"groups"`
}

// LinkGroupPageResponse 友链分组分页响应
type LinkGroupPageResponse struct {
	base.PaginationResponse[dto.LinkGroupNormalDTO]
}

// LinkGroupSortResponse 友链分组排序响应
type LinkGroupSortResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"` // 更新的分组数量
}

// LinkGroupStatusResponse 友链分组状态切换响应
type LinkGroupStatusResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"` // 更新后的状态
}

// LinkGroupDeleteResponse 友链分组删除响应
type LinkGroupDeleteResponse struct {
	Message string `json:"message"`
}

// LinkGroupDeleteConflictResponse 友链分组删除冲突响应
type LinkGroupDeleteConflictResponse struct {
	Message      string                           `json:"message"`       // 错误消息
	ConflictInfo LinkGroupDeleteConflictInfo      `json:"conflict_info"` // 冲突信息
	Links        []dto.LinkGroupDeleteConflictDTO `json:"links"`         // 冲突的友链列表（前10个）
}

// LinkGroupDeleteConflictInfo 删除冲突的详细信息
type LinkGroupDeleteConflictInfo struct {
	GroupUUID    string `json:"group_uuid"`    // 分组UUID
	GroupName    string `json:"group_name"`    // 分组名称
	TotalLinks   int    `json:"total_links"`   // 总友链数量
	DisplayCount int    `json:"display_count"` // 显示的友链数量（最多10个）
}

// 为了保持一致性，定义一些类型别名（引用现有的MessageResponse）
type (
	// LinkGroupDeleteSimpleResponse 简单删除响应（成功时）
	LinkGroupDeleteSimpleResponse struct {
		Message string `json:"message"`
	}

	// LinkGroupStatusSimpleResponse 简单状态切换响应（不需要返回详细状态时）
	LinkGroupStatusSimpleResponse struct {
		Message string `json:"message"`
	}
)
