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

package logcHelper

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"strings"

	mailTemplates "github.com/bamboo-services/bamboo-main/templates/mail"
)

// MailTemplateLogic 邮件模板逻辑
type MailTemplateLogic struct{}

// GetTemplate 获取邮件模板原始内容
//
// 参数说明:
//   - templateName: 模板名称（不含扩展名，如 "approved"）
//
// 返回值:
//   - 模板 HTML 内容
//   - 错误信息
func (m *MailTemplateLogic) GetTemplate(templateName string) (string, error) {
	fileName := templateName + ".html"
	content, err := mailTemplates.TemplateFS.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("模板 %s 不存在: %w", templateName, err)
	}
	return string(content), nil
}

// RenderTemplate 渲染模板（变量替换）
//
// 使用 Go 标准 html/template 语法进行变量替换
//
// 参数说明:
//   - templateName: 模板名称（不含扩展名）
//   - variables: 变量键值对（如 {"Username": "张三", "LinkName": "我的博客"}）
//
// 返回值:
//   - 渲染后的 HTML 内容
//   - 错误信息
func (m *MailTemplateLogic) RenderTemplate(templateName string, variables map[string]string) (string, error) {
	// 读取模板内容
	content, err := m.GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	// 解析模板
	tmpl, err := template.New(templateName).Parse(content)
	if err != nil {
		return "", fmt.Errorf("解析模板 %s 失败: %w", templateName, err)
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", fmt.Errorf("渲染模板 %s 失败: %w", templateName, err)
	}

	return buf.String(), nil
}

// ListTemplates 列出所有可用模板
//
// 返回值:
//   - 模板名称列表（不含扩展名）
//   - 错误信息
func (m *MailTemplateLogic) ListTemplates() ([]string, error) {
	var templates []string

	entries, err := fs.ReadDir(mailTemplates.TemplateFS, ".")
	if err != nil {
		return nil, fmt.Errorf("读取模板目录失败: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html") {
			// 移除 .html 扩展名
			name := strings.TrimSuffix(entry.Name(), ".html")
			templates = append(templates, name)
		}
	}

	return templates, nil
}
