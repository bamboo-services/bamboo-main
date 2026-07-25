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

// SponsorChannel 表示一个赞助渠道实体，用于管理不同的赞助来源渠道。
//
// 该类型包含渠道的唯一标识符、名称、图标、描述、排序等信息。
// 同时记录该渠道的创建时间和更新时间，便于数据管理和审计。
//
// 注意: 赞助渠道通过外键关联到赞助记录表，一个渠道可以有多条赞助记录。
type SponsorChannel struct {
	xModels.BaseEntity         // 嵌入基础实体（SnowflakeID 主键 + timestamptz 时间戳）
	Name               string  `json:"name" gorm:"type:varchar(50);not null;uniqueIndex;comment:渠道名称"`            // 渠道名称（如：微信赞赏、支付宝、爱发电）
	Icon               *string `json:"icon,omitempty" gorm:"type:varchar(500);comment:渠道图标地址"`                    // 渠道图标URL
	Description        *string `json:"description,omitempty" gorm:"type:text;comment:渠道描述"`                       // 渠道描述说明
	SortOrder          int     `json:"sort_order" gorm:"type:int;default:0;comment:渠道排序"`                         // 排序值，数字越大越靠前
	Status             bool    `json:"status" gorm:"type:boolean;default:true;comment:渠道状态（false: 禁用, true: 启用）"` // 渠道状态

	// 关联关系
	SponsorsFKey []*SponsorRecord `json:"sponsors_f_key,omitempty" gorm:"foreignKey:ChannelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:赞助记录外键"` // 赞助记录外键，关联 SponsorRecord 类型
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *SponsorChannel) GetGene() xSnowflake.Gene {
	return bConst.GeneSponsorChannel
}
