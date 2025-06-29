package link_friend

import (
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"

	"bamboo-main/api/link_friend/v1"
)

// LinkFriendAdd
//
// 添加友情链接的方法。
// 接收上下文和请求参数，调用服务层逻辑实现添加友链功能，成功时返回结果响应。
func (c *ControllerV1) LinkFriendAdd(ctx context.Context, req *v1.LinkFriendAddReq) (res *v1.LinkFriendAddRes, err error) {
	blog.ControllerInfo(ctx, "LinkFriendAdd", "添加 %s 友链", req.Name)

	// 构建数据层
	newFriend := base.LinkFriendDTO{
		LinkName:         req.Name,
		LinkUrl:          req.URL,
		LinkAvatar:       req.Avatar,
		LinkDesc:         req.Description,
		LinkEmail:        req.Email,
		LinkGroupUuid:    req.Group,
		LinkColorUuid:    req.Color,
		LinkOrder:        req.Order,
		LinkReviewRemark: req.ReviewRemark,
	}

	// 调用逻辑层处理
	iFriend := service.Friend()
	errorCode := iFriend.AddFriend(ctx, newFriend)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkFriendAddRes{
		ResponseDTO: bresult.Success(ctx, "添加成功"),
	}, nil
}
