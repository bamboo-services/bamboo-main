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
	"bamboo-main/pkg/cerror"
	"context"
	"fmt"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
)

// SendMail
//
// 发送邮件；
// 如果发送成功，则返回 nil；
// 如果发送失败，则返回错误码；
// 如果获取邮件模板失败，则返回错误码。
//
// mailTo: 收件人邮箱地址；
// subject: 邮件主题；
// templateName: 邮件模板名称；
// needBCC: 是否需要抄送给管理员；
// data: 邮件模板数据，键值对形式，模板中使用 <%key%> 的方式引用。
func (s *sEmail) SendMail(ctx context.Context, mailTo, subject, templateName string, needBCC bool, data map[string]interface{}) *berror.ErrorCode {
	blog.ServiceInfo(ctx, "SendMail", "发送邮件，收件人：%s，主题：%s，模板名称：%s", mailTo, subject, templateName)
	var (
		EmailFromName     = dao.System.MustGetSystemValue(ctx, consts.SystemEmailFromNameKey)
		EmailFromAddress  = dao.System.MustGetSystemValue(ctx, consts.SystemEmailFromAddressKey)
		AdminEmail        = dao.System.MustGetSystemValue(ctx, consts.SystemUserEmailKey)
		EmailSmtpHost     = dao.System.MustGetSystemValue(ctx, consts.SystemEmailHostKey)
		EmailSmtpPort     = dao.System.MustGetSystemValue(ctx, consts.SystemEmailPortKey)
		EmailSmtpUsername = dao.System.MustGetSystemValue(ctx, consts.SystemEmailUsernameKey)
		EmailSmtpPassword = dao.System.MustGetSystemValue(ctx, consts.SystemEmailPasswordKey)
	)

	// 获取邮件模板
	getTemplate, errorCode := getMailTemplate(ctx, templateName, data)
	if errorCode != nil {
		blog.ServiceError(ctx, "SendMail", "获取邮件模板失败，模板名称：%s", templateName)
		return errorCode
	}

	fromStringBuilder := strings.Builder{}
	{
		fromStringBuilder.WriteString(EmailFromName)
		fromStringBuilder.WriteString(" <")
		fromStringBuilder.WriteString(EmailFromAddress)
		fromStringBuilder.WriteString(">")
	}

	newEmail := email.Email{
		To:      []string{mailTo},
		Subject: fmt.Sprintf("%s-%s", dao.System.MustGetSystemValue(ctx, consts.SystemNameKey), subject),
		From:    fromStringBuilder.String(),
		HTML:    getTemplate,
	}

	if needBCC && AdminEmail != mailTo {
		newEmail.Bcc = []string{AdminEmail}
	}

	mailErr := newEmail.Send(
		fmt.Sprintf("%s:%s", EmailSmtpHost, EmailSmtpPort),
		smtp.PlainAuth("", EmailSmtpUsername, EmailSmtpPassword, EmailSmtpHost),
	)
	if mailErr != nil {
		blog.ServiceError(ctx, "SendMail", "发送邮件失败，收件人：%s，错误信息：%v", mailTo, mailErr)
		return berror.ErrorAddData(cerror.ErrMailSend, mailErr.Error())
	}
	return nil
}

// getMailTemplate
//
// 获取邮件模板内容；
// 如果获取成功，则返回模板内容的字节切片；
// 如果获取失败，则返回错误码。
func getMailTemplate(ctx context.Context, templateName string, data map[string]interface{}) ([]byte, *berror.ErrorCode) {
	getTemplateContent := gres.GetContent("template/mail/" + templateName + ".html")
	if len(getTemplateContent) <= 0 {
		blog.ServiceError(ctx, "getMailTemplate", "获取邮件模板失败，模板名称：%s", templateName)
		return nil, berror.ErrorAddData(&berror.ErrResourceNotFound, "获取邮件模板失败")
	}

	// 定义邮件模板变量为 <%value%> 形式
	content := gconv.String(getTemplateContent)
	for key, value := range data {
		placeholder := "<% " + key + " %>"
		placeholderTrim := "<%" + key + "%>"
		if value == nil {
			value = ""
		}
		// 尝试替换
		content = strings.ReplaceAll(content, placeholder, gconv.String(value))
		// 尝试去除空格后的二次替换
		content = strings.ReplaceAll(content, placeholderTrim, gconv.String(value))
	}

	// 返回得到的邮件模板内容
	return []byte(content), nil
}
