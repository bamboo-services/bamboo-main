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

package base

// BambooConfig 业务自定义配置结构。
//
// 自 v1.0.4 起，数据库与缓存配置已迁移至 bamboo-base-go 的 xOption 声明式配置
// （xOption.WithDatabase / xOption.WithCache 从环境变量自动装配），
// 本结构仅保留框架未覆盖的业务自定义配置（目前仅 Email）。
//
// 该结构通过 startup 节点注入到上下文，供 worker_mail 等模块经
// pkg/util/ctx.GetConfig 读取。
type BambooConfig struct {
	Email EmailConfig `mapstructure:"email" yaml:"email"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost    string `mapstructure:"smtp_host" yaml:"smtp_host"`       // SMTP 服务器地址
	SMTPPort    int    `mapstructure:"smtp_port" yaml:"smtp_port"`       // SMTP 端口（25/587/465）
	Username    string `mapstructure:"username" yaml:"username"`         // SMTP 用户名
	Password    string `mapstructure:"password" yaml:"password"`         // SMTP 密码
	FromEmail   string `mapstructure:"from_email" yaml:"from_email"`     // 发件人邮箱
	FromName    string `mapstructure:"from_name" yaml:"from_name"`       // 发件人名称
	AdminEmail  string `mapstructure:"admin_email" yaml:"admin_email"`   // 管理员邮箱（接收申请通知）
	WorkerCount int    `mapstructure:"worker_count" yaml:"worker_count"` // 工作协程数（默认4）
	MaxRetry    int    `mapstructure:"max_retry" yaml:"max_retry"`       // 最大重试次数（默认3）
	Timeout     int    `mapstructure:"timeout" yaml:"timeout"`           // 发送超时秒数（默认10）
	UseTLS      bool   `mapstructure:"use_tls" yaml:"use_tls"`           // 是否使用 TLS 直连（465 端口）
	UseStartTLS bool   `mapstructure:"use_starttls" yaml:"use_starttls"` // 是否使用 STARTTLS（587 端口）
}
