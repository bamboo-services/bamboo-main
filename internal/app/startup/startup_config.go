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

package startup

import (
	"context"
	"os"
	"strconv"

	"github.com/bamboo-services/bamboo-main/internal/models/base"

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	xEnv "github.com/bamboo-services/bamboo-base-go/defined/env"
)

// emailConfigInit 构造业务邮件配置并注入到上下文。
//
// 自 v1.0.4 起，数据库与缓存配置由 xOption 从环境变量自动装配，
// 本节点仅负责框架未覆盖的 Email 业务配置，供 worker_mail 经
// pkg/util/ctx.GetConfig 读取。
func (r *reg) emailConfigInit(ctx context.Context) (any, error) {
	log := xLog.WithName(xLog.NamedINIT)
	log.Info(ctx, "加载邮件环境变量配置")

	cfg := &base.BambooConfig{
		Email: base.EmailConfig{
			SMTPHost:    xEnv.GetEnvString(xEnv.EmailHost, ""),
			SMTPPort:    xEnv.GetEnvInt(xEnv.EmailPort, 465),
			Username:    xEnv.GetEnvString(xEnv.EmailUser, ""),
			Password:    xEnv.GetEnvString(xEnv.EmailPass, ""),
			FromEmail:   xEnv.GetEnvString(xEnv.EmailFrom, ""),
			FromName:    getEnvStringByKey("EMAIL_FROM_NAME", "竹叶"),
			AdminEmail:  getEnvStringByKey("EMAIL_ADMIN_EMAIL", ""),
			WorkerCount: getEnvIntByKey("EMAIL_WORKER_COUNT", 4),
			MaxRetry:    getEnvIntByKey("EMAIL_MAX_RETRY", 3),
			Timeout:     getEnvIntByKey("EMAIL_TIMEOUT", 10),
			UseTLS:      getEnvBoolByKey("EMAIL_USE_TLS", true),
			UseStartTLS: getEnvBoolByKey("EMAIL_USE_STARTTLS", false),
		},
	}

	return cfg, nil
}

func getEnvStringByKey(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

func getEnvIntByKey(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvBoolByKey(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	boolValue, ok := xUtil.Parse().Bool(value)
	if !ok {
		return defaultValue
	}
	return boolValue
}
