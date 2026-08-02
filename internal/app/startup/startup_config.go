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

	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/internal/service/screenshot"
	"github.com/bamboo-services/bamboo-main/internal/models/base"

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
)

// emailConfigInit 构造业务邮件配置并注入到上下文。
//
// 自邮件系统迁移至框架 xEmail 插件后，SMTP 连接配置由插件经 EMAIL_* 环境变量装配，
// 本节点仅负责框架未覆盖的邮件业务配置（管理员邮箱），供 logic 经
// pkg/util/ctx.GetConfig 读取。
func (r *reg) emailConfigInit(ctx context.Context) (any, error) {
	log := xLog.WithName(xLog.NamedINIT)
	log.Info(ctx, "加载邮件业务环境变量配置")

	cfg := &base.BambooConfig{
		Email: base.EmailConfig{
			AdminEmail: getEnvStringByKey("EMAIL_ADMIN_EMAIL", ""),
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

// screenshotManagerInit 构造友链截图任务管理器并注入到上下文。
//
// 截图服务配置经 SCREENSHOT_* 环境变量装配（env-first 约定），
// worker 常驻协程由 main.go 的 Runner 附加协程启动。
func (r *reg) screenshotManagerInit(ctx context.Context) (any, error) {
	log := xLog.WithName(xLog.NamedINIT)
	log.Info(ctx, "加载友链截图服务配置")

	cfg := screenshot.LoadConfig()
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return screenshot.NewManager(cfg, repository.NewLinkRepo(db, m), nil, nil), nil
}
