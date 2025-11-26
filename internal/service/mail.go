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

package service

import (
	"bamboo-main/internal/logic"
	servHelper "bamboo-main/internal/service/helper"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	"github.com/gin-gonic/gin"
)

// IMailService 邮件服务接口
type IMailService interface {
	// GetTemplate 获取邮件模板原始内容
	GetTemplate(templateName string) (string, error)

	// RenderTemplate 渲染模板（变量替换）
	RenderTemplate(templateName string, variables map[string]string) (string, error)

	// ListTemplates 列出所有可用模板
	ListTemplates() ([]string, error)

	// SendMail 发送邮件（入队）
	SendMail(ctx *gin.Context, to []string, subject, body string) *xError.Error

	// SendMailWithCC 发送邮件（带抄送）
	SendMailWithCC(ctx *gin.Context, to, cc []string, subject, body string) *xError.Error

	// SendWithTemplate 一键发送（模板 + 变量 → 直接发送）
	SendWithTemplate(ctx *gin.Context, templateName string, to []string, subject string, variables map[string]string) *xError.Error

	// SendWithTemplateAndCC 一键发送（带抄送）
	SendWithTemplateAndCC(ctx *gin.Context, templateName string, to, cc []string, subject string, variables map[string]string) *xError.Error
}

// NewMailService 创建邮件服务实例
func NewMailService() IMailService {
	return &logic.MailLogic{
		TemplateService: servHelper.NewMailTemplateService(),
		MaxRetry:        3,
	}
}

// NewMailServiceWithRetry 创建邮件服务实例（自定义重试次数）
func NewMailServiceWithRetry(maxRetry int) IMailService {
	return &logic.MailLogic{
		TemplateService: servHelper.NewMailTemplateService(),
		MaxRetry:        maxRetry,
	}
}
