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

// LinkColorAddResponse 添加友链颜色响应
type LinkColorAddResponse struct {
	dto.LinkColorDetailDTO
}

// LinkColorUpdateResponse 更新友链颜色响应
type LinkColorUpdateResponse struct {
	dto.LinkColorDetailDTO
}

// LinkColorDetailResponse 友链颜色详情响应
type LinkColorDetailResponse struct {
	dto.LinkColorDetailDTO
}

// LinkColorListResponse 友链颜色列表响应
type LinkColorListResponse struct {
	Colors []dto.LinkColorListDTO `json:"colors"`
}

// LinkColorPageResponse 友链颜色分页响应
type LinkColorPageResponse struct {
	base.PaginationResponse[dto.LinkColorNormalDTO]
}

// LinkColorSortResponse 友链颜色排序响应
type LinkColorSortResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"` // 更新的颜色数量
}

// LinkColorStatusResponse 友链颜色状态切换响应
type LinkColorStatusResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"` // 更新后的状态
}

// LinkColorDeleteResponse 友链颜色删除响应
type LinkColorDeleteResponse struct {
	Message string `json:"message"`
}

// LinkColorDeleteConflictResponse 友链颜色删除冲突响应
type LinkColorDeleteConflictResponse struct {
	Message      string                           `json:"message"`       // 错误消息
	ConflictInfo LinkColorDeleteConflictInfo      `json:"conflict_info"` // 冲突信息
	Links        []dto.LinkColorDeleteConflictDTO `json:"links"`         // 冲突的友链列表（前10个）
}

// LinkColorDeleteConflictInfo 删除冲突的详细信息
type LinkColorDeleteConflictInfo struct {
	ColorID      int64  `json:"color_id"`      // 颜色ID
	ColorName    string `json:"color_name"`    // 颜色名称
	TotalLinks   int    `json:"total_links"`   // 总友链数量
	DisplayCount int    `json:"display_count"` // 显示的友链数量（最多10个）
}
