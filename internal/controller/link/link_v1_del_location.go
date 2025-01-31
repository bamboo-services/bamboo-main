package link

import (
	"context"
	"xiaoMain/internal/service"

	"xiaoMain/api/link/v1"
)

func (c *ControllerV1) DelLocation(ctx context.Context, req *v1.DelLocationReq) (res *v1.DelLocationRes, err error) {
	// 授权头检查
	err = service.Auth().IsUserLogin(ctx)
	if err != nil {
		return nil, err
	}
	// 友链迁移
	err = service.Link().LocationMove(ctx, req.ID, req.MoveID)
	if err != nil {
		return nil, err
	}
	// 位置删除
	err = service.Link().DelLocation(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
