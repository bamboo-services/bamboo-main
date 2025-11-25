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

package dto

import (
	"time"
)

// LinkColorSimpleDTO 友链颜色简单DTO - 用于列表查询
type LinkColorSimpleDTO struct {
	ID    int64  `json:"id"`    // 颜色主键
	Name  string `json:"name"`  // 颜色名称
	Value string `json:"value"` // 颜色值
}

// LinkColorNormalDTO 友链颜色标准DTO - 用于分页查询
type LinkColorNormalDTO struct {
	ID        int64     `json:"id"`         // 颜色主键
	Name      string    `json:"name"`       // 颜色名称
	Value     string    `json:"value"`      // 颜色值
	SortOrder int       `json:"sort_order"` // 排序
	Status    int       `json:"status"`     // 状态
	LinkCount int       `json:"link_count"` // 使用此颜色的友链数量
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// LinkColorDetailDTO 友链颜色详细DTO - 用于详情查询
type LinkColorDetailDTO struct {
	ID        int64     `json:"id"`         // 颜色主键
	Name      string    `json:"name"`       // 颜色名称
	Value     string    `json:"value"`      // 颜色值
	SortOrder int       `json:"sort_order"` // 排序
	Status    int       `json:"status"`     // 状态
	LinkCount int       `json:"link_count"` // 使用此颜色的友链数量
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间

	// 关联信息
	Links []LinkFriendSimpleDTO `json:"links,omitempty"` // 使用此颜色的友链列表（可选）
}
