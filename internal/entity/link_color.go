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

// LinkColor 表示一个友链颜色实体，用于友情链接的颜色主题管理。
//
// 该类型支持两种颜色模式：
// - Type=0 (普通模式): 需要设置 PrimaryColor、SubColor、HoverColor 三个颜色值
// - Type=1 (炫彩模式): 颜色字段全部为空，由前端实现炫彩效果
type LinkColor struct {
	xModels.BaseEntity         // 嵌入基础实体（SnowflakeID 主键 + timestamptz 时间戳）
	Name               string  `json:"name" gorm:"type:varchar(50);not null;comment:颜色名称"`                             // 颜色名称
	Type               int     `json:"type" gorm:"type:int;default:0;not null;comment:颜色类型（0: 普通, 1: 炫彩）"`             // 颜色类型
	PrimaryColor       *string `json:"primary_color,omitempty" gorm:"type:varchar(9);comment:主颜色（如#FF0000或#FF0000FF）"` // 主颜色
	SubColor           *string `json:"sub_color,omitempty" gorm:"type:varchar(9);comment:副颜色"`                         // 副颜色
	HoverColor         *string `json:"hover_color,omitempty" gorm:"type:varchar(9);comment:悬停颜色"`                      // 悬停颜色
	SortOrder          int     `json:"sort_order" gorm:"type:int;default:0;comment:颜色排序"`                              // 颜色排序
	Status             bool    `json:"status" gorm:"type:boolean;default:true;comment:颜色状态（false: 禁用, true: 启用）"`      // 颜色状态

	// 关联关系
	LinksFKey []*LinkFriend `json:"links_f_key,omitempty" gorm:"foreignKey:ColorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链外键"` // 友链外键，关联 LinkFriend 类型
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *LinkColor) GetGene() xSnowflake.Gene {
	return bConst.GeneLinkColor
}
