package response

import (
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
)

// LinkFriendPageResponse 友情链接分页响应
type LinkFriendPageResponse struct {
	Data       []dto.LinkFriendNormalDTO `json:"data"`       // 友情链接数据列表
	Pagination base.PaginationInfo       `json:"pagination"` // 分页信息
}

// LinkFriendListResponse 友情链接列表响应（公开API）
type LinkFriendListResponse struct {
	Data []dto.LinkFriendSimpleDTO `json:"data"` // 友情链接简单数据列表
}

// LinkFriendDetailResponse 友情链接详情响应
type LinkFriendDetailResponse struct {
	Data dto.LinkFriendDetailDTO `json:"data"` // 友情链接详细信息
}

// LinkFriendCreateResponse 友情链接创建响应
type LinkFriendCreateResponse struct {
	Data dto.LinkFriendDetailDTO `json:"data"` // 新创建的友情链接信息
}

// LinkFriendUpdateResponse 友情链接更新响应
type LinkFriendUpdateResponse struct {
	Data dto.LinkFriendDetailDTO `json:"data"` // 更新后的友情链接信息
}