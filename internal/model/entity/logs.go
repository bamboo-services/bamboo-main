// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Logs is the golang structure for table logs.
type Logs struct {
	LogUuid      string      `json:"log_uuid"       orm:"log_uuid"       description:"日志唯一标识符"` // 日志唯一标识符
	LogType      int         `json:"log_type"       orm:"log_type"       description:"日志类型"`    // 日志类型
	LogContent   string      `json:"log_content"    orm:"log_content"    description:"日志内容"`    // 日志内容
	LogCreatedAt *gtime.Time `json:"log_created_at" orm:"log_created_at" description:"日志创建时间"`  // 日志创建时间
}
