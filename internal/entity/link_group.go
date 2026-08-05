// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package entity

import (
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xModels "github.com/bamboo-services/bamboo-base-go/major/models"

	bConst "github.com/bamboo-services/bamboo-main/pkg/constants"
)

// LinkGroup 表示一个友链分组实体，用于对友情链接进行分类管理。
//
// 该类型包含分组的唯一标识符、名称、描述、排序等信息。
// 同时记录该分组的创建时间和更新时间，便于数据管理和审计。
type LinkGroup struct {
	xModels.BaseEntity         // 嵌入基础实体（SnowflakeID 主键 + timestamptz 时间戳）
	Name               string  `json:"name" gorm:"type:varchar(100);not null;comment:分组名称"`                       // 分组名称
	Description        *string `json:"description,omitempty" gorm:"type:text;comment:分组描述"`                       // 分组描述
	SortOrder          int     `json:"sort_order" gorm:"type:int;default:0;comment:分组排序"`                         // 分组排序
	Status             bool    `json:"status" gorm:"type:boolean;default:true;comment:分组状态（false: 禁用, true: 启用）"` // 分组状态

	// 关联关系
	// 注意：constraint:"-" 跳过数据库外键约束——首页/友链页为内置分组（不落库），
	// 友链 group_id 会引用保留 ID，数据库层面需放行该引用值；分组关联的清理由业务层保证。
	LinksFKey []*LinkFriend `json:"links_f_key,omitempty" gorm:"foreignKey:GroupID;references:ID;constraint:-;comment:友链外键"` // 友链外键，关联 LinkFriend 类型
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *LinkGroup) GetGene() xSnowflake.Gene {
	return bConst.GeneLinkGroup
}

// NewBuiltinGroups 构造内置分组对象列表（固定顺序：首页 → 友链页）。
//
// 内置分组为系统预设位置（不落库）：分组列表接口与友链查询返回时以此虚拟对象表达选项，
// ID 固定为 constants.BuiltinGroupHomepageID / BuiltinGroupFriendsID（雪花 ID 空间之外的保留值）。
func NewBuiltinGroups() []*LinkGroup {
	return []*LinkGroup{
		{
			BaseEntity: xModels.BaseEntity{ID: bConst.BuiltinGroupHomepageID},
			Name:       "首页",
			SortOrder:  0,
			Status:     true,
		},
		{
			BaseEntity: xModels.BaseEntity{ID: bConst.BuiltinGroupFriendsID},
			Name:       "友链页",
			SortOrder:  1,
			Status:     true,
		},
	}
}

// BuiltinGroupByID 按保留 ID 返回内置分组对象；非内置 ID 返回 nil。
func BuiltinGroupByID(id xSnowflake.SnowflakeID) *LinkGroup {
	for _, group := range NewBuiltinGroups() {
		if group.ID == id {
			return group
		}
	}
	return nil
}

// IsBuiltinGroupID 判断是否为内置分组的保留 ID。
func IsBuiltinGroupID(id xSnowflake.SnowflakeID) bool {
	return id == bConst.BuiltinGroupHomepageID || id == bConst.BuiltinGroupFriendsID
}
