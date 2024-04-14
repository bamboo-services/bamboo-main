// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// XfLogs is the golang structure for table xf_logs.sql.
type XfLogs struct {
	Id        int64       `json:"id"        ` // 主键
	Type      int         `json:"type"      ` // 日志类型
	Log       string      `json:"log"       ` // 日志内容
	CreatedAt *gtime.Time `json:"createdAt" ` // 日志时间
}
