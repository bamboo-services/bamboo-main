package info

import (
	"context"
	"xiaoMain/internal/service"

	"xiaoMain/api/info/v1"
)

// InfoUser
//
// # 获取用户信息
//
// 获取当前登录用户的信息。
//
// # 参数
//   - ctx  		上下文
//   - req  		请求参数
//
// # 响应
//   - res  		响应参数
//   - err  		错误信息
func (c *ControllerV1) InfoUser(ctx context.Context, req *v1.InfoUserReq) (res *v1.InfoUserRes, err error) {
	// 获取 token 信息进行解析是否登录有效
	err = service.Auth().IsUserLogin(ctx)
	if err != nil {
		return nil, err
	}
	// 获取用户信息
	current, err := service.User().GetUserCurrent(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.InfoUserRes{
		UserCurrentDTO: *current,
	}, nil
}
