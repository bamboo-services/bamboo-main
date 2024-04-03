// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// XfToken is the golang structure for table xf_token.
type XfToken struct {
	Id           int64       `json:"id"           ` // 主键
	UserUuid     string      `json:"userUuid"     ` // 用户 UUID
	UserToken    string      `json:"userToken"    ` // 用户 TOKEN
	UserIp       string      `json:"userIp"       ` // 用户 IP 地址
	UserAgent    string      `json:"userAgent"    ` // 用户 Agent
	Verification string      `json:"verification" ` // 验证用户是否是唯一用户
	CreatedAt    *gtime.Time `json:"createdAt"    ` // 创建时间
	ExpiredAt    *gtime.Time `json:"expiredAt"    ` // 修改时间
}
