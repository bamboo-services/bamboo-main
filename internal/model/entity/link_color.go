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

// LinkColor 表示一个友链颜色实体，用于友情链接的颜色主题管理。
//
// 该类型支持两种颜色模式：
// - Type=0 (普通模式): 需要设置 PrimaryColor、SubColor、HoverColor 三个颜色值
// - Type=1 (炫彩模式): 颜色字段全部为空，由前端实现炫彩效果
type LinkColor struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement:false;comment:颜色唯一标识符"`                         // ID 使用 Snowflake ID 作为主键
	Name         string    `json:"name" gorm:"type:varchar(50);not null;comment:颜色名称"`                               // 颜色名称
	Type         int       `json:"type" gorm:"type:int;default:0;not null;comment:颜色类型（0: 普通, 1: 炫彩）"`               // 颜色类型
	PrimaryColor *string   `json:"primary_color" gorm:"type:varchar(9);comment:主颜色（如#FF0000或#FF0000FF）"`             // 主颜色
	SubColor     *string   `json:"sub_color" gorm:"type:varchar(9);comment:副颜色"`                                     // 副颜色
	HoverColor   *string   `json:"hover_color" gorm:"type:varchar(9);comment:悬停颜色"`                                  // 悬停颜色
	SortOrder    int       `json:"sort_order" gorm:"type:int;default:0;comment:颜色排序"`                                // 颜色排序
	Status       bool      `json:"status" gorm:"type:boolean;default:true;comment:颜色状态（false: 禁用, true: 启用）"`        // 颜色状态
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间

	// 关联关系
	LinksFKey []*LinkFriend `json:"links_f_key" gorm:"foreignKey:ColorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链外键"` // 友链外键，关联 LinkFriend 类型
}

// BeforeCreate 在创建 LinkColor 实例前自动调用，确保为 ID 字段生成唯一标识符。
func (c *LinkColor) BeforeCreate(tx *gorm.DB) error {
	if c.ID == 0 {
		// 从 Context 中获取 Snowflake 节点
		if val := tx.Statement.Context.Value(xConsts.ContextSnowflakeNode.String()); val != nil {
			if node, ok := val.(*snowflake.Node); ok {
				c.ID = node.Generate().Int64()
			} else {
				return errors.New("上下文中的无效雪花节点")
			}
		} else {
			return errors.New("snowflake 节点在上下文中未出现")
		}
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
