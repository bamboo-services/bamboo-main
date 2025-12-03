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

// SponsorRecordBasicInfo 赞助记录基础信息 - 避免循环依赖
type SponsorRecordBasicInfo struct {
	ID       int64  `json:"id"`       // 赞助记录主键
	Nickname string `json:"nickname"` // 赞助者昵称
	Amount   int64  `json:"amount"`   // 赞助金额
}

// SponsorRecordSimpleDTO 赞助记录简单DTO - 用于前台公开展示
type SponsorRecordSimpleDTO struct {
	ID          int64                    `json:"id"`                // 赞助记录主键
	Nickname    string                   `json:"nickname"`          // 赞助者昵称(匿名时显示"匿名用户")
	RedirectURL *string                  `json:"redirect_url"`      // 跳转链接
	Amount      int64                    `json:"amount"`            // 赞助金额
	Message     *string                  `json:"message"`           // 赞助留言
	SponsorAt   *time.Time               `json:"sponsor_at"`        // 赞助时间
	Channel     *SponsorChannelSimpleDTO `json:"channel,omitempty"` // 赞助渠道
}

// SponsorRecordNormalDTO 赞助记录标准DTO - 用于后台分页列表
type SponsorRecordNormalDTO struct {
	ID          int64                    `json:"id"`                // 赞助记录主键
	Nickname    string                   `json:"nickname"`          // 赞助者昵称
	RedirectURL *string                  `json:"redirect_url"`      // 跳转链接
	Amount      int64                    `json:"amount"`            // 赞助金额
	ChannelID   *int64                   `json:"channel_id"`        // 渠道ID
	Message     *string                  `json:"message"`           // 赞助留言
	SponsorAt   *time.Time               `json:"sponsor_at"`        // 赞助时间
	SortOrder   int                      `json:"sort_order"`        // 排序
	IsAnonymous bool                     `json:"is_anonymous"`      // 是否匿名
	IsHidden    bool                     `json:"is_hidden"`         // 是否隐藏
	CreatedAt   time.Time                `json:"created_at"`        // 创建时间
	UpdatedAt   time.Time                `json:"updated_at"`        // 更新时间
	Channel     *SponsorChannelSimpleDTO `json:"channel,omitempty"` // 赞助渠道
}

// SponsorRecordDetailDTO 赞助记录详细DTO - 用于详情查询
type SponsorRecordDetailDTO struct {
	SponsorRecordNormalDTO
}
