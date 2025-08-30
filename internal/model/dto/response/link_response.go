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
