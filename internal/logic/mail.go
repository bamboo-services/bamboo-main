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

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	xEmail "github.com/bamboo-services/bamboo-base-go/plugins/email"
)

// MailLogic 邮件业务逻辑
//
// 基于 bamboo-base-go 的 xEmail 插件同步发送邮件，邮件客户端经节点化注入上下文，
// 模板由框架模板系统渲染（模板经 go:embed 内嵌，startup 节点注入 xEmail 客户端）。
//
// 调用方应通过 xAsync.Async 将本方法的执行置于异步 goroutine 中，以避免
// SMTP 拨号阻塞 HTTP 请求上下文；ctx 应为 xAsync 解耦后的独立 context（携带
// 组件容器引用），而非 gin.Context 本身。
type MailLogic struct{}

// NewMailLogic 创建 MailLogic 实例。
func NewMailLogic() *MailLogic {
	return &MailLogic{}
}

// SendWithTemplate 一键发送（填入模板名称、变量，直接发送）
//
// 这是最常用的发送方式，自动完成模板渲染与发送操作。
//
// 参数说明:
//   - ctx: 上下文（建议为 xAsync 解耦后的独立 context）
//   - templateName: 模板名称（如 "approved", "rejected"）
//   - to: 收件人邮箱列表
//   - subject: 邮件主题
//   - variables: 模板变量
//
// 返回值:
//   - 错误信息
func (m *MailLogic) SendWithTemplate(ctx context.Context, templateName string, to []string, subject string, variables map[string]string) *xError.Error {
	client, xerr := xCtxUtil.GetEmailClient(ctx)
	if xerr != nil {
		return xerr
	}

	msg := &xEmail.Message{
		To:           to,
		Subject:      subject,
		Template:     templateName,
		TemplateData: variables,
	}

	if err := client.SendTemplate(ctx, msg); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "发送邮件失败", false, err)
	}

	xLog.WithName(xLog.NamedLOGC, "MAIL").Info(ctx, fmt.Sprintf("邮件已发送: To=%v, Template=%s", to, templateName))

	return nil
}
