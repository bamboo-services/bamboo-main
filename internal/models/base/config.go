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

package base

// BambooConfig 业务自定义配置结构。
//
// 自 v1.0.4 起，数据库与缓存配置已迁移至 bamboo-base-go 的 xOption 声明式配置
// （xOption.WithDatabase / xOption.WithCache 从环境变量自动装配）；
// 自邮件系统迁移至框架 xEmail 插件后，SMTP 连接配置亦由插件经 EMAIL_* 环境变量装配，
// 本结构仅保留框架未覆盖的业务自定义配置（目前仅邮件业务配置）。
//
// 该结构通过 startup 节点注入到上下文，供 logic 等模块经
// pkg/util/ctx.GetConfig 读取。
type BambooConfig struct {
	Email EmailConfig `mapstructure:"email" yaml:"email"`
}

// EmailConfig 邮件业务配置
//
// SMTP 连接相关配置（服务器、端口、认证、TLS 策略等）由框架 xEmail 插件直接读取
// EMAIL_* 环境变量，此处仅保留框架未覆盖的业务级配置。
type EmailConfig struct {
	AdminEmail string `mapstructure:"admin_email" yaml:"admin_email"` // 管理员邮箱（接收友链申请通知）
}
