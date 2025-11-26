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

package servHelper

import (
	logcHelper "bamboo-main/internal/logic/helper"
)

// IMailTemplateService 邮件模板服务接口
type IMailTemplateService interface {
	// GetTemplate 获取模板原始内容
	GetTemplate(templateName string) (string, error)

	// RenderTemplate 渲染模板（变量替换）
	RenderTemplate(templateName string, variables map[string]string) (string, error)

	// ListTemplates 列出所有模板
	ListTemplates() ([]string, error)
}

// NewMailTemplateService 创建邮件模板服务实例
func NewMailTemplateService() IMailTemplateService {
	return &logcHelper.MailTemplateLogic{}
}
