package link_friend

import (
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"

	"bamboo-main/api/link_friend/v1"
)

// LinkFriendPage
//
// 获取友情链接列表的方法。接收上下文和分页请求参数，调用服务层获取友情链接数据并返回结果。
// 方法会记录操作日志，并处理可能的错误情况。
func (c *ControllerV1) LinkFriendPage(ctx context.Context, req *v1.LinkFriendPageReq) (res *v1.LinkFriendPageRes, err error) {
	blog.ControllerInfo(ctx, "LinkFriendPage", "获取友情链接列表")

	// 调用逻辑层处理
	iFriend := service.Friend()
	getPageDTO, errorCode := iFriend.GetPage(ctx, req.Search, req.Page, req.Size, req.IsAll)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkFriendPageRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "获取成功", *getPageDTO),
	}, nil
}
