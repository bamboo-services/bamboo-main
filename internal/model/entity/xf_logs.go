// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// XfLogs is the golang structure for table xf_logs.
type XfLogs struct {
	Id        int64       `json:"id"        orm:"id"         ` // 主键
	Type      int         `json:"type"      orm:"type"       ` // 日志类型
	Log       string      `json:"log"       orm:"log"        ` // 日志内容
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 日志时间
}
