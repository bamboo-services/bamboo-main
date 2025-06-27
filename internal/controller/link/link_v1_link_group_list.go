package link

import (
	"bamboo-main/api/link/v1"
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
)

// LinkGroupList
//
// 获取友链分组列表的控制器方法；
// 该方法接收一个 LinkGroupListReq 请求参数，包含搜索内容；
// 返回一个 LinkGroupListRes 响应结果，包含所有可用的友链分组信息。
func (c *ControllerV1) LinkGroupList(ctx context.Context, req *v1.LinkGroupListReq) (res *v1.LinkGroupListRes, err error) {
	blog.ControllerInfo(ctx, "LinkGroupList", "获取友链分组列表")

	// 调用逻辑层处理
	iGroup := service.Group()
	linkGroupPage, errorCode := iGroup.GetList(ctx, req.Search, req.Page, req.Size)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkGroupListRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "获取成功", *linkGroupPage),
	}, nil
}
