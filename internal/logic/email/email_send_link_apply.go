/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package email

import (
	"bamboo-main/internal/consts"
	"bamboo-main/internal/dao"
	"bamboo-main/internal/model/dto/base"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/gogf/gf/v2/os/gtime"
	"sync"
)

// SendEmailLinkApply
//
// 发送友情链接申请邮件，包含发送申请成功邮件给申请者及申请通知邮件给管理员。
// applyLink 包含友链信息，如名称、URL、邮箱、头像、描述等。
// 如果邮件发送失败，则记录错误日志并返回错误码。
func (s *sEmail) SendEmailLinkApply(ctx context.Context, applyLink base.LinkFriendDTO) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "SendEmailLinkApply", "发送友情链接申请邮件，收件人：%s", applyLink.LinkEmail)

	// 数据准备
	getData := s.getBaseData(ctx)
	getData["link_name"] = applyLink.LinkName
	getData["link_url"] = applyLink.LinkUrl
	getData["link_email"] = applyLink.LinkEmail
	getData["link_avatar"] = applyLink.LinkAvatar
	getData["link_description"] = applyLink.LinkDesc
	getData["link_created_at"] = gtime.Now().Format("Y-m-d H:i:s")
	getData["link_uuid"] = applyLink.LinkUuid

	// 为双方发送邮件
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 发送申请成功邮件给申请者
	go func() {
		defer wg.Done()
		if err := s.SendMail(ctx, applyLink.LinkEmail, "友链申请确认", "friend_link_application_sent", false, getData); err != nil {
			blog.ServiceError(ctx, "SendEmailLinkApply", "发送申请成功邮件失败，收件人：%s, 错误：%v", applyLink.LinkEmail, err)
			return
		}
		blog.ServiceInfo(ctx, "SendEmailLinkApply", "发送申请成功邮件成功，收件人：%s", applyLink.LinkEmail)
	}()

	// 发送申请通知邮件给管理员
	go func() {
		defer wg.Done()
		getAdminEmail := dao.System.MustGetSystemValue(ctx, consts.SystemUserEmailKey)
		if err := s.SendMail(ctx, getAdminEmail, "友链申请审核", "friend_link_admin_notification", false, getData); err != nil {
			blog.ServiceError(ctx, "SendEmailLinkApply", "发送申请通知邮件失败，收件人：%s, 错误：%v", getAdminEmail, err)
			return
		}
		blog.ServiceInfo(ctx, "SendEmailLinkApply", "发送申请通知邮件成功，收件人：%s", getAdminEmail)
	}()
	wg.Wait()

	return nil
}
