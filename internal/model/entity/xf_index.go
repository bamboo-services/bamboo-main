// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// XfIndex is the golang structure for table xf_index.
type XfIndex struct {
	Id        int         `json:"id"        ` // 主键
	Key       string      `json:"key"       ` // 键
	Value     string      `json:"value"     ` // 值
	CreatedAt *gtime.Time `json:"createdAt" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" ` // 修改时间
}
