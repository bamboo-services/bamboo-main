package response

import (
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
)

// LinkColorPageResponse 友链颜色分页响应
type LinkColorPageResponse struct {
	Data       []dto.LinkColorNormalDTO `json:"data"`       // 友链颜色数据列表
	Pagination base.PaginationInfo      `json:"pagination"` // 分页信息
}

// LinkColorListResponse 友链颜色列表响应
type LinkColorListResponse struct {
	Data []dto.LinkColorSimpleDTO `json:"data"` // 友链颜色简单数据列表
}

// LinkColorDetailResponse 友链颜色详情响应
type LinkColorDetailResponse struct {
	Data dto.LinkColorDetailDTO `json:"data"` // 友链颜色详细信息
}

// LinkColorCreateResponse 友链颜色创建响应
type LinkColorCreateResponse struct {
	Data dto.LinkColorDetailDTO `json:"data"` // 新创建的友链颜色信息
}

// LinkColorUpdateResponse 友链颜色更新响应
type LinkColorUpdateResponse struct {
	Data dto.LinkColorDetailDTO `json:"data"` // 更新后的友链颜色信息
}