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

// LinkColor 表示一个友链颜色实体，用于友情链接的颜色主题管理。
//
// 该类型包含颜色的唯一标识符、名称、颜色值等信息。
// 同时记录该颜色的创建时间和更新时间，便于数据管理和审计。
type LinkColor struct {
	UUID      uuid.UUID `json:"uuid" gorm:"primaryKey;type:char(36);comment:颜色唯一标识符"`                             // UUID 使用 UUID 作为主键
	Name      string    `json:"name" gorm:"type:varchar(50);not null;comment:颜色名称"`                               // 颜色名称
	Value     string    `json:"value" gorm:"type:varchar(7);not null;comment:颜色值（如#FF0000）"`                      // 颜色值
	SortOrder int       `json:"sort_order" gorm:"type:int;default:0;comment:颜色排序"`                                // 颜色排序
	Status    int       `json:"status" gorm:"type:int;default:1;comment:颜色状态（0: 禁用, 1: 启用）"`                      // 颜色状态
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间

	// 关联关系
	LinksFKey []*LinkFriend `json:"links_f_key" gorm:"foreignKey:ColorUUID;references:UUID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链外键"` // 友链外键，关联 LinkFriend 类型
}

// BeforeCreate 在创建 LinkColor 实例前自动调用，确保为 UUID 字段生成唯一标识符。
func (c *LinkColor) BeforeCreate(_ *gorm.DB) error {
	if c.UUID == uuid.Nil {
		newUUID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		c.UUID = newUUID
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 在更新 LinkColor 实例前自动调用，用于更新 UpdatedAt 字段为当前时间戳。
func (c *LinkColor) BeforeUpdate(_ *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}
