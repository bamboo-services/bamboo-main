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
