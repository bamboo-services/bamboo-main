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

package entity

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"

	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
)

// SponsorChannel 表示一个赞助渠道实体，用于管理不同的赞助来源渠道。
//
// 该类型包含渠道的唯一标识符、名称、图标、描述、排序等信息。
// 同时记录该渠道的创建时间和更新时间，便于数据管理和审计。
//
// 注意: 赞助渠道通过外键关联到赞助记录表，一个渠道可以有多条赞助记录。
type SponsorChannel struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement:false;comment:渠道唯一标识符"`                         // ID 使用 Snowflake ID 作为主键
	Name        string    `json:"name" gorm:"type:varchar(50);not null;uniqueIndex;comment:渠道名称"`                   // 渠道名称（如：微信赞赏、支付宝、爱发电）
	Icon        *string   `json:"icon" gorm:"type:varchar(500);comment:渠道图标地址"`                                     // 渠道图标URL
	Description *string   `json:"description" gorm:"type:text;comment:渠道描述"`                                        // 渠道描述说明
	SortOrder   int       `json:"sort_order" gorm:"type:int;default:0;comment:渠道排序"`                                // 排序值，数字越大越靠前
	Status      bool      `json:"status" gorm:"type:boolean;default:true;comment:渠道状态（false: 禁用, true: 启用）"`        // 渠道状态
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间

	// 关联关系
	SponsorsFKey []*SponsorRecord `json:"sponsors_f_key" gorm:"foreignKey:ChannelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:赞助记录外键"` // 赞助记录外键，关联 SponsorRecord 类型
}

// BeforeCreate 在创建 SponsorChannel 实例前自动调用，确保为 ID 字段生成唯一标识符。
func (s *SponsorChannel) BeforeCreate(tx *gorm.DB) error {
	if s.ID == 0 {
		// 从 Context 中获取 Snowflake 节点
		if val := tx.Statement.Context.Value(xConsts.ContextSnowflakeNode.String()); val != nil {
			if node, ok := val.(*snowflake.Node); ok {
				s.ID = node.Generate().Int64()
			} else {
				return errors.New("上下文中的无效雪花节点")
			}
		} else {
			return errors.New("snowflake 节点在上下文中未出现")
		}
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 在更新 SponsorChannel 实例前自动调用，用于更新 UpdatedAt 字段为当前时间戳。
func (s *SponsorChannel) BeforeUpdate(_ *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
