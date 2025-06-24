// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LinkColor is the golang structure for table link_color.
type LinkColor struct {
	ColorUuid      string      `json:"color_uuid"       orm:"color_uuid"       description:"颜色唯一标识符"`            // 颜色唯一标识符
	ColorName      string      `json:"color_name"       orm:"color_name"       description:"颜色名称"`               // 颜色名称
	ColorValue     string      `json:"color_value"      orm:"color_value"      description:"颜色值（如HEX值：#FFFFFF）"` // 颜色值（如HEX值：#FFFFFF）
	ColorDesc      string      `json:"color_desc"       orm:"color_desc"       description:"颜色描述"`               // 颜色描述
	ColorCreatedAt *gtime.Time `json:"color_created_at" orm:"color_created_at" description:"颜色创建时间"`             // 颜色创建时间
	ColorUpdatedAt *gtime.Time `json:"color_updated_at" orm:"color_updated_at" description:"颜色更新时间"`             // 颜色更新时间
}
