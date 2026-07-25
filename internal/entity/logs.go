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
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xModels "github.com/bamboo-services/bamboo-base-go/major/models"
)

// SystemLog 表示一个系统日志实体，用于记录系统操作和事件。
//
// 该类型包含日志的唯一标识符、级别、模块、操作、消息等信息。
// 同时记录操作者、IP地址等审计信息。
//
// 注意：嵌入的 xModels.BaseEntity 会附带 UpdatedAt 字段；日志表通常不更新，
// 该字段保持零值即可，不影响业务。
type SystemLog struct {
	xModels.BaseEntity                         // 嵌入基础实体（SnowflakeID 主键 + timestamptz 时间戳）
	Level              string                  `json:"level" gorm:"type:varchar(20);not null;comment:日志级别（INFO, WARN, ERROR）"` // 日志级别
	Module             string                  `json:"module" gorm:"type:varchar(50);not null;comment:日志模块"`                   // 日志模块
	Action             string                  `json:"action" gorm:"type:varchar(100);not null;comment:操作动作"`                  // 操作动作
	Message            string                  `json:"message" gorm:"type:text;not null;comment:日志消息"`                         // 日志消息
	UserID             *xSnowflake.SnowflakeID `json:"user_id,omitempty" gorm:"comment:操作用户ID"`                                // 操作用户ID
	IPAddress          *string                 `json:"ip_address,omitempty" gorm:"type:varchar(45);comment:操作IP地址"`            // 操作IP地址
	UserAgent          *string                 `json:"user_agent,omitempty" gorm:"type:text;comment:用户代理"`                     // 用户代理
	ExtraData          *string                 `json:"extra_data,omitempty" gorm:"type:text;comment:额外数据（JSON格式）"`             // 额外数据

	// 关联关系
	UserFKey *SystemUser `json:"user_f_key,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:用户外键"` // 用户外键
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *SystemLog) GetGene() xSnowflake.Gene {
	return xSnowflake.GeneLog
}
