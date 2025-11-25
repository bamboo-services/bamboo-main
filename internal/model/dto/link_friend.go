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
	ID     int64  `json:"id"`     // 友链主键
	Name   string `json:"name"`   // 友链名称
	URL    string `json:"url"`    // 友链地址
	Avatar string `json:"avatar"` // 友链头像
}

// LinkFriendNormalDTO 友情链接标准DTO - 用于分页查询
// 包含常用信息，适用于表格展示、卡片列表等场景
type LinkFriendNormalDTO struct {
	ID          int64     `json:"id"`           // 友链主键
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
	ID           int64     `json:"id"`            // 友链主键
	Name         string    `json:"name"`          // 友链名称
	URL          string    `json:"url"`           // 友链地址
	Avatar       string    `json:"avatar"`        // 友链头像
	RSS          string    `json:"rss"`           // RSS地址
	Description  string    `json:"description"`   // 友链描述
	Email        string    `json:"email"`         // 联系邮箱
	GroupID      int64     `json:"group_id"`      // 所属分组ID
	ColorID      int64     `json:"color_id"`      // 颜色ID
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
