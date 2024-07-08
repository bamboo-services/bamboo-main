package link

import (
	"context"
	"xiaoMain/internal/service"

	"xiaoMain/api/link/v1"
)

// LinkEdit
//
// # 编辑链接
//
// 编辑链接，用于编辑链接的信息。
//
// # 参数
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - req: 用户的请求，包含编辑链接的详细信息。
//
// # 返回
//   - res: 发送给用户的响应。如果编辑链接成功，它将返回成功的消息。
//   - err: 在编辑链接过程中发生的任何错误。
func (c *ControllerV1) LinkEdit(ctx context.Context, req *v1.LinkEditReq) (res *v1.LinkEditRes, err error) {
	// 授权信息检查
	err = service.Auth().IsUserLogin(ctx)
	if err != nil {
		return nil, err
	}
	// 对信息内容进行修改操作
	err = service.Link().EditLink(ctx, *req)
	if err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}
