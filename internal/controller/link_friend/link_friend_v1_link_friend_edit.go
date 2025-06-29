package link_friend

import (
	"bamboo-main/api/link_friend/v1"
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
)

// LinkFriendEdit
//
// 编辑友情链接的方法。接收上下文和请求参数，调用服务层实现更新现有友链信息的功能。
// 该方法首先检查分组和颜色是否存在，然后获取要编辑的友链信息，最后更新友链信息并返回结果。
func (c *ControllerV1) LinkFriendEdit(ctx context.Context, req *v1.LinkFriendEditReq) (res *v1.LinkFriendEditRes, err error) {
	blog.ControllerInfo(ctx, "LinkFriendEdit", "编辑 %s 友链", req.LinkUUID)

	// 检查分组是否存在
	iGroup := service.Group()
	getGroupEntity, errorCode := iGroup.GetOneByUUID(ctx, req.Group)
	if errorCode != nil {
		return nil, errorCode
	}
	// 检查颜色是否存在
	iColor := service.Color()
	getColorEntity, errorCode := iColor.GetOneByUUID(ctx, req.Color)
	if errorCode != nil {
		return nil, errorCode
	}
	// 检查 LinkUUID 是否存在
	iFriend := service.Friend()
	getLinkEntity, errorCode := iFriend.GetOneByUUID(ctx, req.LinkUUID)
	if errorCode != nil {
		return nil, errorCode
	}

	// 构建数据层
	editFriendEntity := base.LinkFriendDTO{
		LinkUuid:         getLinkEntity.LinkUuid,
		LinkName:         req.Name,
		LinkUrl:          req.URL,
		LinkAvatar:       req.Avatar,
		LinkRss:          req.Rss,
		LinkDesc:         req.Description,
		LinkEmail:        req.Email,
		LinkGroupUuid:    getGroupEntity.GroupUuid,
		LinkColorUuid:    getColorEntity.ColorUuid,
		LinkOrder:        req.Order,
		LinkReviewRemark: req.ReviewRemark,
	}

	// 调用逻辑层处理
	errorCode = iFriend.Update(ctx, editFriendEntity)
	if errorCode != nil {
		return nil, errorCode
	}
	return &v1.LinkFriendEditRes{
		ResponseDTO: bresult.Success(ctx, "编辑成功"),
	}, nil
}
