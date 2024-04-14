// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// XfLogs is the golang structure of table xf_logs.sql for DAO operations like Where/Data.
type XfLogs struct {
	g.Meta    `orm:"table:xf_logs.sql, do:true"`
	Id        interface{} // 主键
	Type      interface{} // 日志类型
	Log       interface{} // 日志内容
	CreatedAt *gtime.Time // 日志时间
}
