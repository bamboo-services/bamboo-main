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

	// 检查颜色是否存在
	iColor := service.Color()
	getColorEntity, errorCode := iColor.GetOneByUUID(ctx, req.Color)
	if errorCode != nil {
		return nil, errorCode
	}
	// 检查分组是否存在
	iGroup := service.Group()
	getGroupEntity, errorCode := iGroup.GetOneByUUID(ctx, req.Group)
	if errorCode != nil {
		return nil, errorCode
	}

	// 构建数据层
	newFriend := base.LinkFriendDTO{
		LinkName:         req.Name,
		LinkUrl:          req.URL,
		LinkAvatar:       req.Avatar,
		LinkDesc:         req.Description,
		LinkEmail:        req.Email,
		LinkRss:          req.Rss,
		LinkGroupUuid:    getGroupEntity.GroupUuid,
		LinkColorUuid:    getColorEntity.ColorUuid,
		LinkOrder:        req.Order,
		LinkReviewRemark: req.ReviewRemark,
	}

	// 调用逻辑层处理
	iFriend := service.Friend()
	errorCode = iFriend.AddFriend(ctx, newFriend)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.LinkFriendAddRes{
		ResponseDTO: bresult.Success(ctx, "添加成功"),
	}, nil
}
