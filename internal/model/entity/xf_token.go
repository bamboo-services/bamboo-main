// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// XfToken is the golang structure for table xf_token.
type XfToken struct {
	Id           int64       `json:"id"           orm:"id"           ` // 主键
	UserUuid     string      `json:"userUuid"     orm:"user_uuid"    ` // 用户 UUID
	UserToken    string      `json:"userToken"    orm:"user_token"   ` // 用户 TOKEN
	UserIp       string      `json:"userIp"       orm:"user_ip"      ` // 用户 IP 地址
	UserAgent    string      `json:"userAgent"    orm:"user_agent"   ` // 用户 Agent
	Verification string      `json:"verification" orm:"verification" ` // 验证用户是否是唯一用户
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"   ` // 创建时间
	ExpiredAt    *gtime.Time `json:"expiredAt"    orm:"expired_at"   ` // 修改时间
}
