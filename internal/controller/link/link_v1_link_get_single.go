package link

import (
	"context"
	"xiaoMain/internal/service"

	"xiaoMain/api/link/v1"
)

// LinkGetSingle
//
// # 获取单个链接
//
// 用于获取单个链接，用于展示出来；在一般情况下，该接口用于查询链接的详细信息用于修改某个链接
// 等的操作。
//
// # 参数
//   - ctx: 请求的上下文，用于管理超时和取消信号。
//   - req: 用户的请求，包含获取链接的详细信息。
//
// # 返回
//   - res: 发送给用户的响应。如果获取链接成功，它将返回成功的消息。
//   - err: 在获取链接过程中发生的任何错误。
func (c *ControllerV1) LinkGetSingle(
	ctx context.Context,
	req *v1.LinkGetSingleReq,
) (res *v1.LinkGetSingleRes, err error) {
	// 检查登录的权限
	err = service.Auth().IsUserLogin(ctx)
	if err != nil {
		return nil, err
	}
	// 根据请求的 ID 获取单个链接
	linkInfo, err := service.Link().GetSingleLink(ctx, req.ID)
	if err != nil {
		return nil, err
	} else {
		return &v1.LinkGetSingleRes{
			LinkList: *linkInfo,
		}, nil
	}
}
