/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package logic

import (
	"context"
	"fmt"

	"github.com/bamboo-services/bamboo-base-go/application/v2"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/ctxkey"
	"github.com/bamboo-services/bamboo-base-go/plugins/email/v2"
)

// MailLogic 邮件业务逻辑
//
// 基于 bamboo-base-go email/v2 插件发送邮件，邮件客户端经容器注入上下文，
// 模板由框架模板系统渲染（模板经 go:embed 内嵌，由 email.WithTemplateFS 注入）。
//
// 调用方应通过 async.Submit 将本方法的执行置于受管异步任务中，以避免
// SMTP 拨号阻塞 HTTP 请求上下文。
type MailLogic struct{}

// NewMailLogic 创建 MailLogic 实例。
func NewMailLogic() *MailLogic {
	return &MailLogic{}
}

// SendWithTemplate 一键发送（填入模板名称、变量，直接发送）
func (m *MailLogic) SendWithTemplate(ctx context.Context, templateName string, to []string, subject string, variables map[string]string) *apperr.Error {
	client, xerr := application.Get[*email.Client](ctx, ctxkey.EmailClientKey)
	if xerr != nil {
		return xerr
	}

	msg := &email.Message{
		To:           to,
		Subject:      subject,
		Template:     templateName,
		TemplateData: variables,
	}

	if err := client.SendTemplate(ctx, msg); err != nil {
		return apperr.New(email.EmailError, "发送邮件失败", err)
	}

	resolveLogicNamedLogger(ctx, "MAIL").Info(ctx, fmt.Sprintf("邮件已发送: To=%v, Template=%s", to, templateName))
	return nil
}
