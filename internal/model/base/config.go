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

// BambooConfig 主配置结构
type BambooConfig struct {
	Xlf       BMConfig        `mapstructure:"xlf" yaml:"xlf"`
	Snowflake SnowflakeConfig `mapstructure:"snowflake" yaml:"snowflake"`
	Database  DatabaseConfig  `mapstructure:"database" yaml:"database"`
	NoSQL     NoSQLConfig     `mapstructure:"nosql" yaml:"nosql"`
	Email     EmailConfig     `mapstructure:"email" yaml:"email"`
}

// BMConfig 应用配置
type BMConfig struct {
	Debug  bool         `mapstructure:"debug" yaml:"debug"`
	Server ServerConfig `mapstructure:"server" yaml:"server"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port" yaml:"port"`
}

// SnowflakeConfig 表示用于配置 Snowflake ID 生成器的结构。
type SnowflakeConfig struct {
	NodeID *int64 `mapstructure:"node_id" yaml:"node_id"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Pass     string `mapstructure:"pass" yaml:"pass"`
	Name     string `mapstructure:"name" yaml:"name"`
	Prefix   string `mapstructure:"prefix" yaml:"prefix"`
	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode"`
	TimeZone string `mapstructure:"timezone" yaml:"timezone"`
}

// NoSQLConfig Redis配置
type NoSQLConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Pass     string `mapstructure:"pass" yaml:"pass"`
	Database int    `mapstructure:"database" yaml:"database"`
	Prefix   string `mapstructure:"prefix" yaml:"prefix"`
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
