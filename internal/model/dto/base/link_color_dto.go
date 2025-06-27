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

package base

import "github.com/gogf/gf/v2/os/gtime"

// LinkColorDTO
//
// 表示链接颜色的数据传输对象，用于描述颜色的详细信息。
// 包含颜色的唯一标识符、名称、值（如HEX值）、描述、创建时间和更新时间等字段。
type LinkColorDTO struct {
	ColorUuid      string      `json:"color_uuid"       description:"颜色唯一标识符" v:"required"`            // 颜色唯一标识符
	ColorName      string      `json:"color_name"       description:"颜色名称" v:"required"`               // 颜色名称
	ColorValue     string      `json:"color_value"      description:"颜色值（如HEX值：#FFFFFF）" v:"required"` // 颜色值（如HEX值：#FFFFFF）
	ColorDesc      string      `json:"color_desc"       description:"颜色描述"`                            // 颜色描述
	ColorCreatedAt *gtime.Time `json:"color_created_at" description:"颜色创建时间" v:"required"`             // 颜色创建时间
	ColorUpdatedAt *gtime.Time `json:"color_updated_at" description:"颜色更新时间" v:"required"`             // 颜色更新时间
}
