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

// LinkFriendSimpleDTO 友情链接简单DTO - 用于列表查询
// 只包含最基本的信息，适用于下拉选择、简单列表等场景
type LinkFriendSimpleDTO struct {
	UUID   string `json:"uuid"`   // 友链主键
	Name   string `json:"name"`   // 友链名称
	URL    string `json:"url"`    // 友链地址
	Avatar string `json:"avatar"` // 友链头像
}

// LinkFriendNormalDTO 友情链接标准DTO - 用于分页查询
// 包含常用信息，适用于表格展示、卡片列表等场景
type LinkFriendNormalDTO struct {
	UUID        string    `json:"uuid"`         // 友链主键
	Name        string    `json:"name"`         // 友链名称
	URL         string    `json:"url"`          // 友链地址
	Avatar      string    `json:"avatar"`       // 友链头像
	Description string    `json:"description"`  // 友链描述
	Status      int       `json:"status"`       // 友链状态
	StatusText  string    `json:"status_text"`  // 友链状态文本
	IsFailure   int       `json:"is_failure"`   // 是否失效
	FailureText string    `json:"failure_text"` // 失效状态文本
	SortOrder   int       `json:"sort_order"`   // 排序
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}

// LinkFriendDetailDTO 友情链接详细DTO - 用于详情查询
// 包含所有信息，适用于详情页面、编辑表单等场景
type LinkFriendDetailDTO struct {
	UUID         string    `json:"uuid"`          // 友链主键
	Name         string    `json:"name"`          // 友链名称
	URL          string    `json:"url"`           // 友链地址
	Avatar       string    `json:"avatar"`        // 友链头像
	RSS          string    `json:"rss"`           // RSS地址
	Description  string    `json:"description"`   // 友链描述
	Email        string    `json:"email"`         // 联系邮箱
	GroupUUID    string    `json:"group_uuid"`    // 所属分组UUID
	ColorUUID    string    `json:"color_uuid"`    // 颜色UUID
	SortOrder    int       `json:"sort_order"`    // 排序
	Status       int       `json:"status"`        // 友链状态
	StatusText   string    `json:"status_text"`   // 友链状态文本
	IsFailure    int       `json:"is_failure"`    // 是否失效
	FailureText  string    `json:"failure_text"`  // 失效状态文本
	FailReason   string    `json:"fail_reason"`   // 失效原因
	ApplyRemark  string    `json:"apply_remark"`  // 申请备注
	ReviewRemark string    `json:"review_remark"` // 审核备注
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`    // 更新时间

	// 关联信息
	GroupInfo *LinkGroupSimpleDTO `json:"group_info,omitempty"` // 分组信息（可选）
	ColorInfo *LinkColorSimpleDTO `json:"color_info,omitempty"` // 颜色信息（可选）
}

// LinkFriendDTO 兼容性DTO，保持向后兼容
type LinkFriendDTO = LinkFriendDetailDTO

// LinkGroupSimpleDTO 友链分组简单DTO - 用于在LinkFriend中引用
type LinkGroupSimpleDTO struct {
	UUID string `json:"uuid"` // 分组主键
	Name string `json:"name"` // 分组名称
}

// LinkColorSimpleDTO 友链颜色简单DTO - 用于列表查询
type LinkColorSimpleDTO struct {
	UUID  string `json:"uuid"`  // 颜色主键
	Name  string `json:"name"`  // 颜色名称
	Value string `json:"value"` // 颜色值
}

// LinkColorNormalDTO 友链颜色标准DTO - 用于分页查询
type LinkColorNormalDTO struct {
	UUID      string    `json:"uuid"`       // 颜色主键
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
	UUID      string    `json:"uuid"`       // 颜色主键
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

// LinkColorDTO 兼容性DTO，保持向后兼容
type LinkColorDTO = LinkColorDetailDTO

// PaginationDTO 分页数据传输对象
type PaginationDTO[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// SystemUserSimpleDTO 系统用户简单DTO - 用于列表查询
type SystemUserSimpleDTO struct {
	UUID     string `json:"uuid"`     // 用户主键
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 昵称
}

// SystemUserNormalDTO 系统用户标准DTO - 用于分页查询
type SystemUserNormalDTO struct {
	UUID        string     `json:"uuid"`          // 用户主键
	Username    string     `json:"username"`      // 用户名
	Email       string     `json:"email"`         // 邮箱
	Nickname    string     `json:"nickname"`      // 昵称
	Role        string     `json:"role"`          // 角色
	Status      int        `json:"status"`        // 状态
	LastLoginAt *time.Time `json:"last_login_at"` // 最后登录时间
	CreatedAt   time.Time  `json:"created_at"`    // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`    // 更新时间
}

// SystemUserDetailDTO 系统用户详细DTO - 用于详情查询
type SystemUserDetailDTO struct {
	UUID        string     `json:"uuid"`          // 用户主键
	Username    string     `json:"username"`      // 用户名
	Email       string     `json:"email"`         // 邮箱
	Nickname    string     `json:"nickname"`      // 昵称
	Avatar      string     `json:"avatar"`        // 头像
	Role        string     `json:"role"`          // 角色
	Status      int        `json:"status"`        // 状态
	LastLoginAt *time.Time `json:"last_login_at"` // 最后登录时间
	CreatedAt   time.Time  `json:"created_at"`    // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`    // 更新时间
}

// SystemUserDTO 兼容性DTO，保持向后兼容
type SystemUserDTO = SystemUserDetailDTO
