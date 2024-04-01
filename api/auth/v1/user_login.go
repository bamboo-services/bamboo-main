package v1

import (
	"develop/utility/resultUtil"
	"github.com/gogf/gf/v2/frame/g"
)

// UserLoginReq
// 用户登陆请求
type UserLoginReq struct {
	g.Meta `path:"/api/v1/user/login" tags:"User" method:"POST" summary:"用户登陆"`
}

// UserLoginRes
// 用户登陆返回
type UserLoginRes struct {
	resultUtil.BaseResponse
	g.Meta `mime:"application/json"`
}
