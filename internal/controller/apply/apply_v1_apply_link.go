package apply

import (
	"bamboo-main/api/apply/v1"
	"bamboo-main/internal/model/dto/base"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
	"github.com/google/uuid"
)

// ApplyLink 处理友情链接申请请求，验证输入数据的有效性，并将申请信息存储到数据库中。
// 同时发送通知邮件给申请者和管理员。返回操作结果。
func (c *ControllerV1) ApplyLink(ctx context.Context, req *v1.ApplyLinkReq) (res *v1.ApplyLinkRes, err error) {
	blog.ControllerInfo(ctx, "ApplyLink", "%s 申请友情链接 %s", req.Email, req.Name)

	// 检查分组是否存在
	getGroupEntity, errorCode := c.iGroup.GetOneByUUID(ctx, req.Group)
	if errorCode != nil {
		return nil, errorCode
	}
	// 检查颜色是否存在
	getColorEntity, errorCode := c.iColor.GetOneByUUID(ctx, req.Color)
	if errorCode != nil {
		return nil, errorCode
	}
	// 检查 LinkURL 是否存在
	_, errorCode = c.iFriend.GetOneByURL(ctx, req.URL)
	if errorCode != nil {
		if errorCode.GetErrorCode() != berror.ErrNotFound {
			return nil, errorCode
		}
	}

	// 构建数据层
	newApplyLink := base.LinkFriendDTO{
		LinkUuid:        uuid.New().String(),
		LinkName:        req.Name,
		LinkUrl:         req.URL,
		LinkAvatar:      req.Avatar,
		LinkRss:         req.Rss,
		LinkDesc:        req.Description,
		LinkEmail:       req.Email,
		LinkGroupUuid:   getGroupEntity.GroupUuid,
		LinkColorUuid:   getColorEntity.ColorUuid,
		LinkApplyRemark: req.ApplyRemark,
	}

	// 调用逻辑层处理
	errorCode = c.iApply.NewApplyLink(ctx, newApplyLink)
	if errorCode != nil {
		return nil, errorCode
	}
	// 发送通知邮件
	errorCode = c.iEmail.SendEmailLinkApply(ctx, newApplyLink)
	if errorCode != nil {
		return nil, errorCode
	}

	// 内容返回
	return &v1.ApplyLinkRes{
		ResponseDTO: bresult.Success(ctx, "申请友情链接成功"),
	}, nil
}
