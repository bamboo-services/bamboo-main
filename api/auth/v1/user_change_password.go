package v1

import "github.com/gogf/gf/v2/frame/g"

// UserChangePasswordReq
// 用户修改密码请求
type UserChangePasswordReq struct {
	g.Meta `path:"/api/v1/user/change-password" tags:"User" method:"PUT" summary:"用户修改密码"`
}

// UserChangePasswordRes
// 用户修改密码返回
type UserChangePasswordRes struct {
	g.Meta `mime:"application/json"`
}
