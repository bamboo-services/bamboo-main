/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明:版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息,请查看项目根目录下的LICENSE文件或访问:
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package response

import (
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
)

// SponsorChannelAddResponse 添加渠道响应
type SponsorChannelAddResponse struct {
	dto.SponsorChannelDetailDTO
}

// SponsorChannelUpdateResponse 更新渠道响应
type SponsorChannelUpdateResponse struct {
	dto.SponsorChannelDetailDTO
}

// SponsorChannelDetailResponse 详情响应
type SponsorChannelDetailResponse struct {
	dto.SponsorChannelDetailDTO
}

// SponsorChannelListResponse 列表响应（不分页）
type SponsorChannelListResponse struct {
	Channels []dto.SponsorChannelListDTO `json:"channels"`
}

// SponsorChannelPageResponse 分页响应
type SponsorChannelPageResponse struct {
	base.PaginationResponse[dto.SponsorChannelNormalDTO]
}

// SponsorChannelStatusResponse 状态切换响应
type SponsorChannelStatusResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"` // 更新后的状态
}

// SponsorChannelDeleteResponse 删除响应
type SponsorChannelDeleteResponse struct {
	Message string `json:"message"`
}

// SponsorChannelDeleteConflictResponse 删除冲突响应
type SponsorChannelDeleteConflictResponse struct {
	Message      string                                `json:"message"`       // 错误消息
	ConflictInfo SponsorChannelDeleteConflictInfo      `json:"conflict_info"` // 冲突信息
	Sponsors     []dto.SponsorChannelDeleteConflictDTO `json:"sponsors"`      // 冲突的赞助记录列表（前10个）
}

// SponsorChannelDeleteConflictInfo 删除冲突的详细信息
type SponsorChannelDeleteConflictInfo struct {
	ChannelID     int64  `json:"channel_id"`     // 渠道ID
	ChannelName   string `json:"channel_name"`   // 渠道名称
	TotalSponsors int    `json:"total_sponsors"` // 总赞助记录数量
	DisplayCount  int    `json:"display_count"`  // 显示的赞助记录数量（最多10个）
}
