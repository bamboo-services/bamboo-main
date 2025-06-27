package link_group

import (
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"

	"bamboo-main/api/link_group/v1"
)

// LinkGroupPage
//
// 获取友链分组列表的控制器方法；
// 该方法接收一个 LinkGroupPageReq 请求参数，包含搜索内容；
// 返回一个 LinkGroupPageRes 响应结果，包含所有可用的友链分组信息。
func (c *ControllerV1) LinkGroupPage(ctx context.Context, req *v1.LinkGroupPageReq) (res *v1.LinkGroupPageRes, err error) {
	blog.ControllerInfo(ctx, "LinkGroupPage", "获取友链分组列表")

	// 调用逻辑层处理
	iGroup := service.Group()
	linkGroupPage, errorCode := iGroup.GetList(ctx, req.Search, req.Page, req.Size)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkGroupPageRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "获取成功", *linkGroupPage),
	}, nil
}
