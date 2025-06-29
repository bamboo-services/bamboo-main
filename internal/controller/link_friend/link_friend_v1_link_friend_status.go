package link_friend

import (
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"

	"bamboo-main/api/link_friend/v1"
)

// LinkFriendStatus
//
// 修改友情链接状态的方法。接受上下文和请求参数，调用逻辑层实现状态的更新。
// 方法会校验友链当前状态，禁止对已审核的友链再次修改状态。
// 成功时返回更新后的结果响应。
func (c *ControllerV1) LinkFriendStatus(ctx context.Context, req *v1.LinkFriendStatusReq) (res *v1.LinkFriendStatusRes, err error) {
	blog.ControllerInfo(ctx, "LinkFriendStatus", "修改 %s 的友链状态为 %s", req.LinkUUID, req.Status)

	// 调用逻辑层处理
	iFriend := service.Friend()
	friendEntity, errorCode := iFriend.GetOneByUUID(ctx, req.LinkUUID)
	if errorCode != nil {
		return nil, errorCode
	}

	// 检查友链是否已经审核
	if friendEntity.LinkStatus != 0 {
		return nil, berror.ErrorAddData(&berror.ErrBadRequest, "友链已审核，无法修改状态")
	}

	// 更新友链状态
	errorCode = iFriend.UpdateStatus(ctx, req.LinkUUID, req.Status)
	if errorCode != nil {
		return nil, errorCode
	}
	return &v1.LinkFriendStatusRes{
		ResponseDTO: bresult.Success(ctx, "友链状态修改成功"),
	}, nil
}
