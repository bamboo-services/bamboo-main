// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LinkColor is the golang structure of table xf_link_color for DAO operations like Where/Data.
type LinkColor struct {
	g.Meta         `orm:"table:xf_link_color, do:true"`
	ColorUuid      interface{} // 颜色唯一标识符
	ColorName      interface{} // 颜色名称
	ColorValue     interface{} // 颜色值（如HEX值：#FFFFFF）
	ColorDesc      interface{} // 颜色描述
	ColorCreatedAt *gtime.Time // 颜色创建时间
	ColorUpdatedAt *gtime.Time // 颜色更新时间
}
