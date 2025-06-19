// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// System is the golang structure for table system.
type System struct {
	SystemUuid  string      `json:"system_uuid"  orm:"system_uuid"  description:"系统UUID"` // 系统UUID
	SystemName  string      `json:"system_name"  orm:"system_name"  description:"系统名称"`   // 系统名称
	SystemValue *gjson.Json `json:"system_value" orm:"system_value" description:"系统值"`    // 系统值
	CreatedAt   *gtime.Time `json:"created_at"   orm:"created_at"   description:"创建时间"`   // 创建时间
	UpdatedAt   *gtime.Time `json:"updated_at"   orm:"updated_at"   description:"更新时间"`   // 更新时间
}
