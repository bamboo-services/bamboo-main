// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LinkGroup is the golang structure for table link_group.
type LinkGroup struct {
	GroupUuid      string      `json:"group_uuid"       orm:"group_uuid"       description:"分组唯一标识符"` // 分组唯一标识符
	GroupName      string      `json:"group_name"       orm:"group_name"       description:"分组名称"`    // 分组名称
	GroupDesc      string      `json:"group_desc"       orm:"group_desc"       description:"分组描述"`    // 分组描述
	GroupOrder     int         `json:"group_order"      orm:"group_order"      description:"分组排序"`    // 分组排序
	GroupCreatedAt *gtime.Time `json:"group_created_at" orm:"group_created_at" description:"分组创建时间"`  // 分组创建时间
	GroupUpdatedAt *gtime.Time `json:"group_updated_at" orm:"group_updated_at" description:"分组更新时间"`  // 分组更新时间
}
