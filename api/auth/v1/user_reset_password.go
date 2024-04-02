package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"xiaoMain/utility/result"
)

// UserResetPasswordReq
// 用户重置密码请求
type UserResetPasswordReq struct {
	g.Meta `path:"/api/v1/user/reset-password" tags:"User" method:"PUT" summary:"用户重置密码"`
}

// UserResetPasswordRes
// 用户重置密码返回
type UserResetPasswordRes struct {
	result.BaseResponse
	g.Meta `mime:"application/json"`
}
