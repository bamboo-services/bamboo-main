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

// System 表示一个系统配置实体，用于存储键值对配置信息。
//
// 该类型包含配置的唯一标识符、键名、值等信息。
// 同时记录该配置的创建时间和更新时间，便于数据管理和审计。
type System struct {
	xModels.BaseEntity         // 嵌入基础实体（SnowflakeID 主键 + timestamptz 时间戳）
	Key                string  `json:"key" gorm:"type:varchar(100);not null;uniqueIndex;comment:配置键名"` // 配置键名
	Value              *string `json:"value,omitempty" gorm:"type:text;comment:配置值"`                   // 配置值
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *System) GetGene() xSnowflake.Gene {
	return xSnowflake.GeneConfig
}
