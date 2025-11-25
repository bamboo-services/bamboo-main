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

// LinkGroup 表示一个友链分组实体，用于对友情链接进行分类管理。
//
// 该类型包含分组的唯一标识符、名称、描述、排序等信息。
// 同时记录该分组的创建时间和更新时间，便于数据管理和审计。
type LinkGroup struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement:false;comment:分组唯一标识符"`                         // ID 使用 Snowflake ID 作为主键
	Name        string    `json:"name" gorm:"type:varchar(100);not null;comment:分组名称"`                              // 分组名称
	Description *string   `json:"description" gorm:"type:text;comment:分组描述"`                                        // 分组描述
	SortOrder   int       `json:"sort_order" gorm:"type:int;default:0;comment:分组排序"`                                // 分组排序
	Status      bool      `json:"status" gorm:"type:boolean;default:1;comment:分组状态（0: 禁用, 1: 启用）"`                  // 分组状态
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间

	// 关联关系
	LinksFKey []*LinkFriend `json:"links_f_key" gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链外键"` // 友链外键，关联 LinkFriend 类型
}

// BeforeCreate 在创建 LinkGroup 实例前自动调用，确保为 ID 字段生成唯一标识符。
func (g *LinkGroup) BeforeCreate(tx *gorm.DB) error {
	if g.ID == 0 {
		// 从 Context 中获取 Snowflake 节点
		if val := tx.Statement.Context.Value(xConsts.ContextSnowflakeNode.String()); val != nil {
			if node, ok := val.(*snowflake.Node); ok {
				g.ID = node.Generate().Int64()
			} else {
				return errors.New("上下文中的无效雪花节点")
			}
		} else {
			return errors.New("snowflake 节点在上下文中未出现")
		}
	}
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 在更新 LinkGroup 实例前自动调用，用于更新 UpdatedAt 字段为当前时间戳。
func (g *LinkGroup) BeforeUpdate(_ *gorm.DB) error {
	g.UpdatedAt = time.Now()
	return nil
}
