package auth

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"bamboo-main/api/auth/v1"
)

// AuthLogin
//
// 用户登录，需要用户提供用户名和密码。
//
// 参数
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - req: 用户的请求，包含登录的详细信息。
//
// 返回
//   - res: 发送给用户的响应。如果登录成功，它将返回成功的消息。
//   - err: 在登录过程中发生的任何错误。
func (c *ControllerV1) AuthLogin(ctx context.Context, req *v1.AuthLoginReq) (res *v1.AuthLoginRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
