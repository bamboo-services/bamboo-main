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

package rediscache

import "time"

type MailTask struct {
	ID           string            `json:"id"`
	TemplateName string            `json:"template_name"`
	To           []string          `json:"to"`
	Cc           []string          `json:"cc,omitempty"`
	Subject      string            `json:"subject"`
	Body         string            `json:"body"`
	Variables    map[string]string `json:"variables"`
	RetryCount   int               `json:"retry_count"`
	MaxRetry     int               `json:"max_retry"`
	CreatedAt    time.Time         `json:"created_at"`
	NextRetryAt  time.Time         `json:"next_retry_at"`
}
