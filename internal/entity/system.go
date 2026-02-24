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

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	"gorm.io/gorm"

	xCtx "github.com/bamboo-services/bamboo-base-go/defined/context"
)

// System 表示一个系统配置实体，用于存储键值对配置信息。
//
// 该类型包含配置的唯一标识符、键名、值等信息。
// 同时记录该配置的创建时间和更新时间，便于数据管理和审计。
type System struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement:false;comment:系统配置唯一标识符"`                       // ID 使用 Snowflake ID 作为主键
	Key       string    `json:"key" gorm:"type:varchar(100);not null;uniqueIndex;comment:配置键名"`                   // 配置键名
	Value     *string   `json:"value" gorm:"type:text;comment:配置值"`                                               // 配置值
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间
}

// BeforeCreate 在创建 System 实例前自动调用，确保为 ID 字段生成唯一标识符。
func (s *System) BeforeCreate(tx *gorm.DB) error {
	if s.ID == 0 {
		// 从 Context 中获取 Snowflake 节点
		if val := tx.Statement.Context.Value(xCtx.SnowflakeNodeKey); val != nil {
			if node, ok := val.(*xSnowflake.Node); ok {
				s.ID = node.MustGenerate().Int64()
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

// BeforeUpdate 在更新 System 实例前自动调用，用于更新 UpdatedAt 字段为当前时间戳。
func (s *System) BeforeUpdate(_ *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
