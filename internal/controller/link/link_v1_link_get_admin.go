package link

import (
	"context"
	"xiaoMain/internal/service"

	"xiaoMain/api/link/v1"
)

// LinkGetAdmin
//
// # 获取链接
//
// 获取链接，包括友情链接、社交链接等。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - req: 请求参数，包括获取链接的请求参数。
//
// # 返回:
//   - res: 响应参数，返回获取链接的响应结果。
//   - err: 如果处理过程中发生错误，返回错误信息。
func (c *ControllerV1) LinkGetAdmin(ctx context.Context, req *v1.LinkGetAdminReq) (res *v1.LinkGetAdminRes, err error) {
	// 检查用户是否登录
	err = service.Auth().IsUserLogin(ctx)
	if err != nil {
		return nil, err
	}
	// 获取链接
	link, total, err := service.Link().GetLinkAdmin(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.LinkGetAdminRes{
		Links: link,
		Total: total,
	}, nil
}
