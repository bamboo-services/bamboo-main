// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// XfToken is the golang structure of table xf_token for DAO operations like Where/Data.
type XfToken struct {
	g.Meta    `orm:"table:xf_token, do:true"`
	Id        interface{} // 主键
	UserUuid  interface{} // 用户 UUID
	UserToken interface{} // 用户 TOKEN
	CreatedAt *gtime.Time // 创建时间
	ExpiredAt *gtime.Time // 过期时间
}
