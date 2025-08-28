package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SystemLog 表示一个系统日志实体，用于记录系统操作和事件。
//
// 该类型包含日志的唯一标识符、级别、模块、操作、消息等信息。
// 同时记录操作者、IP地址等审计信息。
type SystemLog struct {
	UUID       uuid.UUID  `json:"uuid" gorm:"primaryKey;type:char(36);comment:日志唯一标识符"`                           // UUID 使用 UUID 作为主键
	Level      string     `json:"level" gorm:"type:varchar(20);not null;comment:日志级别（INFO, WARN, ERROR）"`         // 日志级别
	Module     string     `json:"module" gorm:"type:varchar(50);not null;comment:日志模块"`                            // 日志模块
	Action     string     `json:"action" gorm:"type:varchar(100);not null;comment:操作动作"`                           // 操作动作
	Message    string     `json:"message" gorm:"type:text;not null;comment:日志消息"`                                   // 日志消息
	UserUUID   *uuid.UUID `json:"user_uuid" gorm:"type:char(36);comment:操作用户UUID"`                                 // 操作用户UUID
	IPAddress  *string    `json:"ip_address" gorm:"type:varchar(45);comment:操作IP地址"`                               // 操作IP地址
	UserAgent  *string    `json:"user_agent" gorm:"type:text;comment:用户代理"`                                         // 用户代理
	ExtraData  *string    `json:"extra_data" gorm:"type:text;comment:额外数据（JSON格式）"`                               // 额外数据
	CreatedAt  time.Time  `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间

	// 关联关系
	UserFKey *SystemUser `json:"user_f_key" gorm:"foreignKey:UserUUID;references:UUID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:用户外键"` // 用户外键
}

// BeforeCreate 在创建 SystemLog 实例前自动调用，确保为 UUID 字段生成唯一标识符。
func (l *SystemLog) BeforeCreate(_ *gorm.DB) error {
	if l.UUID == uuid.Nil {
		newUUID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		l.UUID = newUUID
	}
	l.CreatedAt = time.Now()
	return nil
}