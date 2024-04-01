// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// XfIndex is the golang structure of table xf_index for DAO operations like Where/Data.
type XfIndex struct {
	g.Meta    `orm:"table:xf_index, do:true"`
	Id        interface{} // 主键
	Key       interface{} // 键
	Value     interface{} // 值
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 修改时间
}
