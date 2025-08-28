package response

import (
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
)

// LinkAddResponse 添加友情链接响应
type LinkAddResponse struct {
	dto.LinkFriendDTO
}

// LinkUpdateResponse 更新友情链接响应
type LinkUpdateResponse struct {
	dto.LinkFriendDTO
}

// LinkDetailResponse 友情链接详情响应
type LinkDetailResponse struct {
	dto.LinkFriendDTO
}

// LinkListResponse 友情链接列表响应
type LinkListResponse struct {
	base.PaginationResponse[dto.LinkFriendDTO]
}

// LinkPublicResponse 公开友情链接响应
type LinkPublicResponse struct {
	Links []dto.LinkFriendDTO `json:"links"`
}