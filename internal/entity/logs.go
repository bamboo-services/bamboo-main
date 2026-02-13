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

// SystemLog 表示一个系统日志实体，用于记录系统操作和事件。
//
// 该类型包含日志的唯一标识符、级别、模块、操作、消息等信息。
// 同时记录操作者、IP地址等审计信息。
type SystemLog struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement:false;comment:日志唯一标识符"`                         // ID 使用 Snowflake ID 作为主键
	Level     string    `json:"level" gorm:"type:varchar(20);not null;comment:日志级别（INFO, WARN, ERROR）"`           // 日志级别
	Module    string    `json:"module" gorm:"type:varchar(50);not null;comment:日志模块"`                             // 日志模块
	Action    string    `json:"action" gorm:"type:varchar(100);not null;comment:操作动作"`                            // 操作动作
	Message   string    `json:"message" gorm:"type:text;not null;comment:日志消息"`                                   // 日志消息
	UserID    *int64    `json:"user_id" gorm:"comment:操作用户ID"`                                                    // 操作用户ID
	IPAddress *string   `json:"ip_address" gorm:"type:varchar(45);comment:操作IP地址"`                                // 操作IP地址
	UserAgent *string   `json:"user_agent" gorm:"type:text;comment:用户代理"`                                         // 用户代理
	ExtraData *string   `json:"extra_data" gorm:"type:text;comment:额外数据（JSON格式）"`                                 // 额外数据
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间

	// 关联关系
	UserFKey *SystemUser `json:"user_f_key" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:用户外键"` // 用户外键
}

// BeforeCreate 在创建 SystemLog 实例前自动调用，确保为 ID 字段生成唯一标识符。
func (l *SystemLog) BeforeCreate(tx *gorm.DB) error {
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
	return nil
}
