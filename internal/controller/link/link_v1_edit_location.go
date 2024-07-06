package link

import (
	"context"
	"xiaoMain/internal/service"

	"xiaoMain/api/link/v1"
)

// EditLocation
//
// # 编辑位置
//
// 编辑位置，用于编辑位置信息。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - req: 请求参数，包括编辑位置的请求参数。
//
// # 返回:
//   - res: 响应参数，返回编辑位置的响应结果。
//   - err: 如果处理过程中发生错误，返回错误信息。
func (c *ControllerV1) EditLocation(ctx context.Context, req *v1.EditLocationReq) (res *v1.EditLocationRes, err error) {
	// 检查用户是否登录了
	err = service.Auth().IsUserLogin(ctx)
	if err != nil {
		return nil, err
	}
	// 对数据进行修改操作
	err = service.Link().EditLocation(ctx, *req)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
