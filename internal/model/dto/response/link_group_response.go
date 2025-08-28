package response

import (
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
)

// LinkGroupPageResponse 友链分组分页响应
type LinkGroupPageResponse struct {
	Data       []dto.LinkGroupNormalDTO `json:"data"`       // 友链分组数据列表
	Pagination base.PaginationInfo      `json:"pagination"` // 分页信息
}

// LinkGroupListResponse 友链分组列表响应
type LinkGroupListResponse struct {
	Data []dto.LinkGroupSimpleDTO `json:"data"` // 友链分组简单数据列表
}

// LinkGroupDetailResponse 友链分组详情响应
type LinkGroupDetailResponse struct {
	Data dto.LinkGroupDetailDTO `json:"data"` // 友链分组详细信息
}

// LinkGroupCreateResponse 友链分组创建响应
type LinkGroupCreateResponse struct {
	Data dto.LinkGroupDetailDTO `json:"data"` // 新创建的友链分组信息
}

// LinkGroupUpdateResponse 友链分组更新响应
type LinkGroupUpdateResponse struct {
	Data dto.LinkGroupDetailDTO `json:"data"` // 更新后的友链分组信息
}