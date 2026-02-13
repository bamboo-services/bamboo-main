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

package dto

import (
	"time"
)

// SponsorChannelSimpleDTO 赞助渠道简单DTO - 用于下拉选择、记录中引用
type SponsorChannelSimpleDTO struct {
	ID   int64   `json:"id"`   // 渠道主键
	Name string  `json:"name"` // 渠道名称
	Icon *string `json:"icon"` // 渠道图标
}

// SponsorChannelNormalDTO 赞助渠道标准DTO - 用于分页列表
type SponsorChannelNormalDTO struct {
	ID           int64     `json:"id"`            // 渠道主键
	Name         string    `json:"name"`          // 渠道名称
	Icon         *string   `json:"icon"`          // 渠道图标
	Description  *string   `json:"description"`   // 渠道描述
	SortOrder    int       `json:"sort_order"`    // 排序
	Status       bool      `json:"status"`        // 状态
	SponsorCount int       `json:"sponsor_count"` // 关联赞助记录数
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`    // 更新时间
}

// SponsorChannelDetailDTO 赞助渠道详细DTO - 用于详情查询
type SponsorChannelDetailDTO struct {
	SponsorChannelNormalDTO
}

// SponsorChannelListDTO 赞助渠道列表DTO - 用于下拉列表
type SponsorChannelListDTO struct {
	ID           int64   `json:"id"`            // 渠道主键
	Name         string  `json:"name"`          // 渠道名称
	Icon         *string `json:"icon"`          // 渠道图标
	SortOrder    int     `json:"sort_order"`    // 排序
	Status       bool    `json:"status"`        // 状态
	SponsorCount int     `json:"sponsor_count"` // 关联赞助记录数
}

// SponsorChannelDeleteConflictDTO 删除冲突DTO - 删除冲突时返回关联记录
type SponsorChannelDeleteConflictDTO struct {
	ID       int64  `json:"id"`       // 赞助记录ID
	Nickname string `json:"nickname"` // 昵称
	Amount   int64  `json:"amount"`   // 金额
}
