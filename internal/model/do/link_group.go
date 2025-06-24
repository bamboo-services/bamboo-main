// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LinkGroup is the golang structure of table xf_link_group for DAO operations like Where/Data.
type LinkGroup struct {
	g.Meta         `orm:"table:xf_link_group, do:true"`
	GroupUuid      interface{} // 分组唯一标识符
	GroupName      interface{} // 分组名称
	GroupDesc      interface{} // 分组描述
	GroupOrder     interface{} // 分组排序
	GroupCreatedAt *gtime.Time // 分组创建时间
	GroupUpdatedAt *gtime.Time // 分组更新时间
}
