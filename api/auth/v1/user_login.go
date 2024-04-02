package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"xiaoMain/utility/result"
)

// UserLoginReq
// 用户登陆请求
type UserLoginReq struct {
	g.Meta `path:"/api/v1/user/login" tags:"User" method:"POST" summary:"用户登陆" json:"g.Meta"`
	User   string `json:"user" v:"required|regex:^[0-9A-Za-z-_]+$	#只允许输入0-9、A-Z、a-Z 以及 - 和 _"`
	Pass   string `json:"pass" v:"required							#请输入密码"`
}

// UserLoginRes
// 用户登陆返回
type UserLoginRes struct {
	result.BaseResponse
	g.Meta `mime:"application/json"`
}
