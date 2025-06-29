package link_friend

import (
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"

	"bamboo-main/api/link_friend/v1"
)

// LinkFriendDelete
//
// 删除友情链接的控制器方法；接收删除请求参数，执行友链删除操作并返回处理结果。
// 方法会根据提供的 LinkUUID 查找并删除对应的友情链接记录。
func (c *ControllerV1) LinkFriendDelete(ctx context.Context, req *v1.LinkFriendDeleteReq) (res *v1.LinkFriendDeleteRes, err error) {
	blog.ControllerInfo(ctx, "LinkFriendDelete", "删除友情链接 %s", req.LinkUUID)

	// 调用逻辑层处理
	iFriend := service.Friend()

	// 检查友链是否存在
	_, errorCode := iFriend.GetOneByUUID(ctx, req.LinkUUID)
	if errorCode != nil {
		return nil, errorCode
	}

	// 执行删除操作
	errorCode = iFriend.Delete(ctx, req.LinkUUID)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkFriendDeleteRes{
		ResponseDTO: bresult.Success(ctx, "删除成功"),
	}, nil
}
