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

// SystemUser 表示一个系统用户实体，用于管理员用户的管理。
//
// 该类型包含用户的唯一标识符、用户名、密码、邮箱、角色等信息。
// 同时记录该用户的创建时间和更新时间，便于数据管理和审计。
type SystemUser struct {
	UUID        uuid.UUID  `json:"uuid" gorm:"primaryKey;type:uuid;comment:系统用户唯一标识符"`                               // UUID 使用 UUID 作为主键
	Username    string     `json:"username" gorm:"type:varchar(50);not null;uniqueIndex;comment:用户名"`                // 用户名
	Password    string     `json:"password" gorm:"type:varchar(255);not null;comment:密码哈希"`                          // 密码哈希
	Email       string     `json:"email" gorm:"type:varchar(100);not null;uniqueIndex;comment:邮箱"`                   // 邮箱
	Nickname    *string    `json:"nickname" gorm:"type:varchar(100);comment:昵称"`                                     // 昵称
	Avatar      *string    `json:"avatar" gorm:"type:varchar(500);comment:头像URL"`                                    // 头像URL
	Role        string     `json:"role" gorm:"type:varchar(20);default:'admin';comment:角色（admin, moderator）"`        // 角色
	Status      int        `json:"status" gorm:"type:int;default:1;comment:状态（0: 禁用, 1: 启用）"`                        // 状态
	LastLoginAt *time.Time `json:"last_login_at" gorm:"type:timestamp;comment:最后登录时间"`                               // 最后登录时间
	CreatedAt   time.Time  `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
	UpdatedAt   time.Time  `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间
}

// BeforeCreate 在创建 SystemUser 实例前自动调用，确保为 UUID 字段生成唯一标识符。
func (s *SystemUser) BeforeCreate(_ *gorm.DB) error {
	if s.UUID == uuid.Nil {
		newUUID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		s.UUID = newUUID
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 在更新 SystemUser 实例前自动调用，用于更新 UpdatedAt 字段为当前时间戳。
func (s *SystemUser) BeforeUpdate(_ *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
