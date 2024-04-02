// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// XfToken is the golang structure for table xf_token.
type XfToken struct {
	Id        int64       `json:"id"        ` // 主键
	UserUuid  string      `json:"userUuid"  ` // 用户 UUID
	UserToken string      `json:"userToken" ` // 用户 TOKEN
	CreatedAt *gtime.Time `json:"createdAt" ` // 创建时间
	ExpiredAt *gtime.Time `json:"expiredAt" ` // 过期时间
}
