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

package logic

import (
	"encoding/json"
	"fmt"
	"time"

	logcHelper "github.com/bamboo-services/bamboo-main/internal/logic/helper"
	dtoRedis "github.com/bamboo-services/bamboo-main/internal/models/dto/redis"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MailLogic 邮件业务逻辑
type MailLogic struct {
	TemplateService *logcHelper.MailTemplateLogic
	MaxRetry        int // 最大重试次数
}

func NewMailLogic() *MailLogic {
	return &MailLogic{TemplateService: &logcHelper.MailTemplateLogic{}, MaxRetry: 3}
}

// GetTemplate 获取邮件模板
//
// 参数说明:
//   - templateName: 模板名称（不含扩展名）
//
// 返回值:
//   - 模板 HTML 内容
//   - 错误信息
func (m *MailLogic) GetTemplate(templateName string) (string, error) {
	return m.TemplateService.GetTemplate(templateName)
}

// RenderTemplate 渲染模板（变量替换）
//
// 参数说明:
//   - templateName: 模板名称
//   - variables: 变量键值对
//
// 返回值:
//   - 渲染后的 HTML 内容
//   - 错误信息
func (m *MailLogic) RenderTemplate(templateName string, variables map[string]string) (string, error) {
	return m.TemplateService.RenderTemplate(templateName, variables)
}

// ListTemplates 列出所有可用模板
//
// 返回值:
//   - 模板名称列表
//   - 错误信息
func (m *MailLogic) ListTemplates() ([]string, error) {
	return m.TemplateService.ListTemplates()
}

// SendMail 发送邮件（加入队列）
//
// 直接发送已渲染的邮件内容
//
// 参数说明:
//   - ctx: Gin 上下文
//   - to: 收件人邮箱列表
//   - subject: 邮件主题
//   - body: 已渲染的 HTML 邮件正文
//
// 返回值:
//   - 错误信息
func (m *MailLogic) SendMail(ctx *gin.Context, to []string, subject, body string) *xError.Error {
	task := &dtoRedis.MailTaskDTO{
		ID:           uuid.New().String(),
		TemplateName: "",
		To:           to,
		Cc:           nil,
		Subject:      subject,
		Body:         body,
		Variables:    nil,
		RetryCount:   0,
		MaxRetry:     m.getMaxRetry(),
		CreatedAt:    time.Now(),
		NextRetryAt:  time.Now(),
	}

	if err := m.enqueueMailTask(ctx, task); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "邮件入队失败", false, err)
	}

	return nil
}

// SendMailWithCC 发送邮件（带抄送）
//
// 参数说明:
//   - ctx: Gin 上下文
//   - to: 收件人邮箱列表
//   - cc: 抄送邮箱列表
//   - subject: 邮件主题
//   - body: 已渲染的 HTML 邮件正文
//
// 返回值:
//   - 错误信息
func (m *MailLogic) SendMailWithCC(ctx *gin.Context, to, cc []string, subject, body string) *xError.Error {
	task := &dtoRedis.MailTaskDTO{
		ID:           uuid.New().String(),
		TemplateName: "",
		To:           to,
		Cc:           cc,
		Subject:      subject,
		Body:         body,
		Variables:    nil,
		RetryCount:   0,
		MaxRetry:     m.getMaxRetry(),
		CreatedAt:    time.Now(),
		NextRetryAt:  time.Now(),
	}

	if err := m.enqueueMailTask(ctx, task); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "邮件入队失败", false, err)
	}

	return nil
}

// SendWithTemplate 一键发送（填入模板名称、变量，直接发送）
//
// 这是最常用的发送方式，自动完成模板渲染和入队操作
//
// 参数说明:
//   - ctx: Gin 上下文
//   - templateName: 模板名称（如 "approved", "rejected"）
//   - to: 收件人邮箱列表
//   - subject: 邮件主题
//   - variables: 模板变量
//
// 返回值:
//   - 错误信息
func (m *MailLogic) SendWithTemplate(ctx *gin.Context, templateName string, to []string, subject string, variables map[string]string) *xError.Error {
	// 渲染模板
	body, err := m.TemplateService.RenderTemplate(templateName, variables)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "渲染邮件模板失败", false, err)
	}

	// 构建任务
	task := &dtoRedis.MailTaskDTO{
		ID:           uuid.New().String(),
		TemplateName: templateName,
		To:           to,
		Cc:           nil,
		Subject:      subject,
		Body:         body,
		Variables:    variables,
		RetryCount:   0,
		MaxRetry:     m.getMaxRetry(),
		CreatedAt:    time.Now(),
		NextRetryAt:  time.Now(),
	}

	// 入队
	if err := m.enqueueMailTask(ctx, task); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "邮件入队失败", false, err)
	}

	// 记录日志
	xLog.WithName(xLog.NamedLOGC, "MAIL").Info(ctx, fmt.Sprintf("邮件任务已入队: ID=%s, To=%v, Template=%s", task.ID, to, templateName))

	return nil
}

// SendWithTemplateAndCC 一键发送（带抄送）
//
// 参数说明:
//   - ctx: Gin 上下文
//   - templateName: 模板名称
//   - to: 收件人邮箱列表
//   - cc: 抄送邮箱列表
//   - subject: 邮件主题
//   - variables: 模板变量
//
// 返回值:
//   - 错误信息
func (m *MailLogic) SendWithTemplateAndCC(ctx *gin.Context, templateName string, to, cc []string, subject string, variables map[string]string) *xError.Error {
	// 渲染模板
	body, err := m.TemplateService.RenderTemplate(templateName, variables)
	if err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "渲染邮件模板失败", false, err)
	}

	// 构建任务
	task := &dtoRedis.MailTaskDTO{
		ID:           uuid.New().String(),
		TemplateName: templateName,
		To:           to,
		Cc:           cc,
		Subject:      subject,
		Body:         body,
		Variables:    variables,
		RetryCount:   0,
		MaxRetry:     m.getMaxRetry(),
		CreatedAt:    time.Now(),
		NextRetryAt:  time.Now(),
	}

	// 入队
	if err := m.enqueueMailTask(ctx, task); err != nil {
		return xError.NewError(ctx, xError.ServerInternalError, "邮件入队失败", false, err)
	}

	// 记录日志
	xLog.WithName(xLog.NamedLOGC, "MAIL").Info(ctx, fmt.Sprintf("邮件任务已入队: ID=%s, To=%v, CC=%v, Template=%s", task.ID, to, cc, templateName))

	return nil
}

// enqueueMailTask 将邮件任务加入 Redis 队列
func (m *MailLogic) enqueueMailTask(ctx *gin.Context, task *dtoRedis.MailTaskDTO) error {
	rdb := ctxUtil.GetRedisClient(ctx)
	if rdb == nil {
		return fmt.Errorf("Redis 客户端不可用")
	}

	// 序列化任务
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("序列化邮件任务失败: %w", err)
	}

	// LPUSH 入队（左进右出，FIFO）
	err = rdb.LPush(ctx.Request.Context(), constants.MailQueueKey, taskJSON).Err()
	if err != nil {
		return fmt.Errorf("邮件任务入队失败: %w", err)
	}

	return nil
}

// getMaxRetry 获取最大重试次数
func (m *MailLogic) getMaxRetry() int {
	if m.MaxRetry > 0 {
		return m.MaxRetry
	}
	return 3 // 默认值
}
