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
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LinkFriend 表示一个友情链接实体，用于描述友情链接及其属性。
//
// 该类型包含友情链接的唯一标识符、名称、URL、头像、RSS、描述等信息。
// 同时记录该友情链接的创建时间和更新时间，便于数据管理和审计。
//
// 注意: 友情链接通过外键关联到分组和颜色表。
type LinkFriend struct {
	UUID         uuid.UUID  `json:"uuid" gorm:"primaryKey;type:char(36);comment:友链唯一标识符"`                             // UUID 使用 UUID 作为主键
	Name         string     `json:"name" gorm:"type:varchar(100);not null;comment:友链名称"`                              // 友链名称
	URL          string     `json:"url" gorm:"type:varchar(500);not null;comment:友链URL地址"`                            // 友链URL地址
	Avatar       *string    `json:"avatar" gorm:"type:varchar(500);comment:友链头像URL"`                                  // 友链头像URL
	RSS          *string    `json:"rss" gorm:"type:varchar(500);comment:友链RSS地址"`                                     // 友链RSS地址
	Description  *string    `json:"description" gorm:"type:text;comment:友链描述"`                                        // 友链描述
	Email        *string    `json:"email" gorm:"type:varchar(100);comment:友链联系邮箱"`                                    // 友链联系邮箱
	GroupUUID    *uuid.UUID `json:"group_uuid" gorm:"type:char(36);comment:所属分组UUID"`                                 // 所属分组UUID
	ColorUUID    *uuid.UUID `json:"color_uuid" gorm:"type:char(36);comment:友链颜色UUID"`                                 // 友链颜色UUID
	SortOrder    int        `json:"sort_order" gorm:"type:int;default:0;comment:友链排序"`                                // 友链排序
	Status       int        `json:"status" gorm:"type:int;default:0;comment:友链状态（0: 待审核, 1: 已通过, 2: 已拒绝）"`            // 友链状态
	IsFailure    int        `json:"is_failure" gorm:"type:int;default:0;comment:友链失效标志（0: 正常, 1: 失效）"`                // 友链失效标志
	FailReason   *string    `json:"fail_reason" gorm:"type:text;comment:友链失效原因"`                                      // 友链失效原因
	ApplyRemark  *string    `json:"apply_remark" gorm:"type:text;comment:申请者备注"`                                      // 申请者备注
	ReviewRemark *string    `json:"review_remark" gorm:"type:text;comment:审核备注"`                                      // 审核备注
	CreatedAt    time.Time  `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间

	// 关联关系
	GroupFKey *LinkGroup `json:"group_f_key" gorm:"foreignKey:GroupUUID;references:UUID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链分组外键"` // 友链分组外键
	ColorFKey *LinkColor `json:"color_f_key" gorm:"foreignKey:ColorUUID;references:UUID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链颜色外键"` // 友链颜色外键
}

// BeforeCreate 在创建 LinkFriend 实例前自动调用，确保为 UUID 字段生成唯一标识符。
func (l *LinkFriend) BeforeCreate(_ *gorm.DB) error {
	if l.UUID == uuid.Nil {
		newUUID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		l.UUID = newUUID
	}
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 在更新 LinkFriend 实例前自动调用，用于更新 UpdatedAt 字段为当前时间戳。
func (l *LinkFriend) BeforeUpdate(_ *gorm.DB) error {
	l.UpdatedAt = time.Now()
	return nil
}
