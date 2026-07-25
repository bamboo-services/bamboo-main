// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package entity

import (
	"time"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xModels "github.com/bamboo-services/bamboo-base-go/major/models"
)

// SystemUser 表示一个系统用户实体，用于管理员用户的管理。
//
// 该类型包含用户的唯一标识符、用户名、密码、邮箱、角色等信息。
// 同时记录该用户的创建时间和更新时间，便于数据管理和审计。
type SystemUser struct {
	xModels.BaseEntity            // 嵌入基础实体（SnowflakeID 主键 + timestamptz 时间戳）
	OAuthUserID        *string    `json:"-" gorm:"type:varchar(255);uniqueIndex;comment:OAuth 用户唯一标识"`               // OAuth 用户唯一标识
	Username           string     `json:"username" gorm:"type:varchar(50);not null;uniqueIndex;comment:用户名"`         // 用户名
	Password           string     `json:"-" gorm:"type:varchar(255);not null;comment:密码哈希"`                          // 密码哈希
	Email              string     `json:"email" gorm:"type:varchar(100);not null;uniqueIndex;comment:邮箱"`            // 邮箱
	Nickname           *string    `json:"nickname,omitempty" gorm:"type:varchar(100);comment:昵称"`                    // 昵称
	Avatar             *string    `json:"avatar,omitempty" gorm:"type:varchar(500);comment:头像URL"`                   // 头像URL
	Role               string     `json:"role" gorm:"type:varchar(20);default:'admin';comment:角色（admin, moderator）"` // 角色
	Status             int        `json:"status" gorm:"type:int;default:1;comment:状态（0: 禁用, 1: 启用）"`                 // 状态
	EmailVerify        bool       `json:"email_verify" gorm:"type:boolean;default:false;comment:邮箱是否已验证"`            // 邮箱验证状态
	LastLoginAt        *time.Time `json:"last_login_at,omitempty" gorm:"type:timestamptz;comment:最后登录时间"`            // 最后登录时间
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *SystemUser) GetGene() xSnowflake.Gene {
	return xSnowflake.GeneUser
}
