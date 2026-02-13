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

	xSnowflake "github.com/bamboo-services/bamboo-base-go/snowflake"
	"gorm.io/gorm"

	xCtx "github.com/bamboo-services/bamboo-base-go/context"
)

// LinkFriend 表示一个友情链接实体，用于描述友情链接及其属性。
//
// 该类型包含友情链接的唯一标识符、名称、URL、头像、RSS、描述等信息。
// 同时记录该友情链接的创建时间和更新时间，便于数据管理和审计。
//
// 注意: 友情链接通过外键关联到分组和颜色表。
type LinkFriend struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement:false;comment:友链唯一标识符"`                         // ID 使用 Snowflake ID 作为主键
	Name         string    `json:"name" gorm:"type:varchar(100);not null;comment:友链名称"`                              // 友链名称
	URL          string    `json:"url" gorm:"type:varchar(500);not null;comment:友链URL地址"`                            // 友链URL地址
	Avatar       *string   `json:"avatar" gorm:"type:varchar(500);comment:友链头像URL"`                                  // 友链头像URL
	RSS          *string   `json:"rss" gorm:"type:varchar(500);comment:友链RSS地址"`                                     // 友链RSS地址
	Description  *string   `json:"description" gorm:"type:text;comment:友链描述"`                                        // 友链描述
	Email        *string   `json:"email" gorm:"type:varchar(100);comment:友链联系邮箱"`                                    // 友链联系邮箱
	GroupID      *int64    `json:"group_id" gorm:"comment:所属分组ID"`                                                   // 所属分组ID
	ColorID      *int64    `json:"color_id" gorm:"comment:友链颜色ID"`                                                   // 友链颜色ID
	SortOrder    int       `json:"sort_order" gorm:"type:int;default:0;comment:友链排序"`                                // 友链排序
	Status       int       `json:"status" gorm:"type:int;default:0;comment:友链状态（0: 待审核, 1: 已通过, 2: 已拒绝）"`            // 友链状态
	IsFailure    int       `json:"is_failure" gorm:"type:int;default:0;comment:友链失效标志（0: 正常, 1: 失效）"`                // 友链失效标志
	FailReason   *string   `json:"fail_reason" gorm:"type:text;comment:友链失效原因"`                                      // 友链失效原因
	ApplyRemark  *string   `json:"apply_remark" gorm:"type:text;comment:申请者备注"`                                      // 申请者备注
	ReviewRemark *string   `json:"review_remark" gorm:"type:text;comment:审核备注"`                                      // 审核备注
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间

	// 关联关系
	GroupFKey *LinkGroup `json:"group_f_key" gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链分组外键"` // 友链分组外键
	ColorFKey *LinkColor `json:"color_f_key" gorm:"foreignKey:ColorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:友链颜色外键"` // 友链颜色外键
}

// BeforeCreate 在创建 LinkFriend 实例前自动调用，确保为 ID 字段生成唯一标识符。
func (l *LinkFriend) BeforeCreate(tx *gorm.DB) error {
	if l.ID == 0 {
		// 从 Context 中获取 Snowflake 节点
		if val := tx.Statement.Context.Value(xCtx.SnowflakeNodeKey); val != nil {
			if node, ok := val.(*xSnowflake.Node); ok {
				l.ID = node.MustGenerate().Int64()
			} else {
				return errors.New("上下文中的无效雪花节点")
			}
		} else {
			return errors.New("snowflake 节点在上下文中未出现")
		}
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
