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

package dtoRedis

import "time"

// MailTaskDTO 邮件任务数据传输对象（Redis 队列消息格式）
type MailTaskDTO struct {
	ID           string            `json:"id"`            // 任务唯一ID（用于日志追踪）
	TemplateName string            `json:"template_name"` // 模板名称（如: apply, approved）
	To           []string          `json:"to"`            // 收件人邮箱列表
	Cc           []string          `json:"cc,omitempty"`  // 抄送邮箱列表
	Subject      string            `json:"subject"`       // 邮件主题
	Body         string            `json:"body"`          // 已渲染的邮件正文（HTML）
	Variables    map[string]string `json:"variables"`     // 模板变量（用于重试时重新渲染）
	RetryCount   int               `json:"retry_count"`   // 当前重试次数
	MaxRetry     int               `json:"max_retry"`     // 最大重试次数（默认 3）
	CreatedAt    time.Time         `json:"created_at"`    // 任务创建时间
	NextRetryAt  time.Time         `json:"next_retry_at"` // 下次重试时间
}
