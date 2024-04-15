// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// XfIndex is the golang structure for table xf_index.
type XfIndex struct {
	Id        int         `json:"id"        orm:"id"         ` // 主键
	Key       string      `json:"key"       orm:"key"        ` // 键
	Value     string      `json:"value"     orm:"value"      ` // 值
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 修改时间
}
