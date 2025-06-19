// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Logs is the golang structure of table xf_logs for DAO operations like Where/Data.
type Logs struct {
	g.Meta       `orm:"table:xf_logs, do:true"`
	LogUuid      interface{} // 日志唯一标识符
	LogType      interface{} // 日志类型
	LogContent   interface{} // 日志内容
	LogCreatedAt *gtime.Time // 日志创建时间
}
