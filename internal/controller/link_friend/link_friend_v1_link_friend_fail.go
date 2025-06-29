package link_friend

import (
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"

	"bamboo-main/api/link_friend/v1"
)

// LinkFriendFail
//
// 更新友情链接失败状态的方法。接受上下文和请求参数，调用逻辑层实现失败状态的更新。
// 方法会校验友链当前状态，记录失败原因。
// 成功时返回更新后的结果响应。
func (c *ControllerV1) LinkFriendFail(ctx context.Context, req *v1.LinkFriendFailReq) (res *v1.LinkFriendFailRes, err error) {
	blog.ControllerInfo(ctx, "LinkFriendFail", "更新 %s 的友链失败状态为 %v，原因：%s", req.LinkUUID, req.IsFail, req.Reason)

	// 调用逻辑层处理
	iFriend := service.Friend()
	friendEntity, errorCode := iFriend.GetOneByUUID(ctx, req.LinkUUID)
	if errorCode != nil {
		return nil, errorCode
	}

	// 检查友链是否已经审核
	if friendEntity.LinkStatus == 0 {
		return nil, berror.ErrorAddData(&berror.ErrBadRequest, "友链未审核，无法更新失败状态")
	}

	// 更新友链失败状态
	errorCode = iFriend.UpdateFailStatus(ctx, req.LinkUUID, req.IsFail, req.Reason)
	if errorCode != nil {
		return nil, errorCode
	}
	return &v1.LinkFriendFailRes{
		ResponseDTO: bresult.Success(ctx, "友链失败状态更新成功"),
	}, nil
}
