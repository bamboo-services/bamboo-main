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

// SponsorRecord 表示一条赞助记录实体，用于记录每一笔赞助的详细信息。
//
// 该类型包含赞助者昵称、跳转地址、赞助金额、赞助渠道、留言等信息。
// 支持匿名展示和隐藏功能，便于灵活管理赞助记录的前台展示。
//
// 注意: 金额以"分"为单位存储，避免浮点精度问题。展示时由前端转换为元。
type SponsorRecord struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement:false;comment:记录唯一标识符"`                         // ID 使用 Snowflake ID 作为主键
	Nickname    string     `json:"nickname" gorm:"type:varchar(100);not null;comment:赞助者昵称"`                         // 赞助者展示名称
	RedirectURL *string    `json:"redirect_url" gorm:"type:varchar(500);comment:赞助者跳转地址"`                            // 点击昵称时的跳转地址
	Amount      int64      `json:"amount" gorm:"type:bigint;not null;comment:赞助金额（单位：分）"`                            // 金额，单位为分（如 666 表示 6.66 元）
	ChannelID   *int64     `json:"channel_id" gorm:"comment:赞助渠道ID;index"`                                           // 关联的赞助渠道ID
	Message     *string    `json:"message" gorm:"type:text;comment:赞助留言"`                                            // 赞助者的留言内容
	SponsorAt   *time.Time `json:"sponsor_at" gorm:"type:timestamp;comment:赞助发生时间"`                                  // 实际赞助发生的时间
	SortOrder   int        `json:"sort_order" gorm:"type:int;default:0;comment:显示排序"`                                // 排序值，用于前台赞助墙排序
	IsAnonymous bool       `json:"is_anonymous" gorm:"type:boolean;default:false;comment:是否匿名展示"`                    // 为 true 时前台显示"匿名用户"
	IsHidden    bool       `json:"is_hidden" gorm:"type:boolean;default:false;comment:是否在前台隐藏"`                      // 为 true 时前台不展示该记录
	CreatedAt   time.Time  `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
	UpdatedAt   time.Time  `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间

	// 关联关系
	ChannelFKey *SponsorChannel `json:"channel_f_key" gorm:"foreignKey:ChannelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:赞助渠道外键"` // 赞助渠道外键，关联 SponsorChannel 类型
}

// BeforeCreate 在创建 SponsorRecord 实例前自动调用，确保为 ID 字段生成唯一标识符。
func (s *SponsorRecord) BeforeCreate(tx *gorm.DB) error {
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

// BeforeUpdate 在更新 SponsorRecord 实例前自动调用，用于更新 UpdatedAt 字段为当前时间戳。
func (s *SponsorRecord) BeforeUpdate(_ *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
