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
	LinksFKey []*LinkFriend `json:"links_f_key,omitempty" gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链外键"` // 友链外键，关联 LinkFriend 类型
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *LinkGroup) GetGene() xSnowflake.Gene {
	return bConst.GeneLinkGroup
}
