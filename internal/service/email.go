// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/XiaoLFeng/bamboo-utils/berror"
)

type (
	IEmail interface {
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
		SendMail(ctx context.Context, mailTo string, subject string, templateName string, needBCC bool, data map[string]interface{}) *berror.ErrorCode
		// SendMailByPasswordReset
		//
		// 发送重置密码的邮件；
		// 如果发送成功，则返回 nil；
		// 如果发送失败，则返回错误码；
		// 如果获取邮件模板失败，则返回错误码。
		SendMailByPasswordReset(ctx context.Context, toMail string) *berror.ErrorCode
	}
)

var (
	localEmail IEmail
)

func Email() IEmail {
	if localEmail == nil {
		panic("implement not found for interface IEmail, forgot register?")
	}
	return localEmail
}

func RegisterEmail(i IEmail) {
	localEmail = i
}
