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

// LinkFriendBasicInfo 友链基础信息 - 避免循环依赖
type LinkFriendBasicInfo struct {
	UUID   string `json:"uuid"`   // 友链主键
	Name   string `json:"name"`   // 友链名称
	URL    string `json:"url"`    // 友链地址
	Avatar string `json:"avatar"` // 友链头像
}

// LinkGroupNormalDTO 友链分组标准DTO - 用于分页查询
type LinkGroupNormalDTO struct {
	UUID        string    `json:"uuid"`        // 分组主键
	Name        string    `json:"name"`        // 分组名称
	Description string    `json:"description"` // 分组描述
	SortOrder   int       `json:"sort_order"`  // 排序
	Status      int       `json:"status"`      // 状态
	LinkCount   int       `json:"link_count"`  // 友链数量
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
}

// LinkGroupDetailDTO 友链分组详细DTO - 用于详情查询
type LinkGroupDetailDTO struct {
	UUID        string    `json:"uuid"`        // 分组主键
	Name        string    `json:"name"`        // 分组名称
	Description string    `json:"description"` // 分组描述
	SortOrder   int       `json:"sort_order"`  // 排序
	Status      int       `json:"status"`      // 状态
	LinkCount   int       `json:"link_count"`  // 友链数量
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间

	// 关联信息 - 避免循环依赖，使用独立的结构
	Links []LinkFriendBasicInfo `json:"links,omitempty"` // 分组下的友链列表（可选）
}

// LinkGroupListDTO 友链分组列表DTO - 用于下拉选择等场景
type LinkGroupListDTO struct {
	UUID      string `json:"uuid"`       // 分组主键
	Name      string `json:"name"`       // 分组名称
	SortOrder int    `json:"sort_order"` // 排序
	Status    int    `json:"status"`     // 状态
	LinkCount int    `json:"link_count"` // 友链数量
}

// LinkGroupDTO 兼容性DTO，保持向后兼容
type LinkGroupDTO = LinkGroupDetailDTO

// LinkGroupDeleteConflictDTO 删除冲突DTO - 用于显示冲突的友链信息
type LinkGroupDeleteConflictDTO struct {
	UUID string `json:"uuid"` // 友链UUID
	Name string `json:"name"` // 友链名称
	URL  string `json:"url"`  // 友链地址
}
