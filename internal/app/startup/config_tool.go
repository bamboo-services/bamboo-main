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
	"sync"

	"github.com/bamboo-services/bamboo-base-go/application/v2"
	"github.com/bamboo-services/bamboo-base-go/application/v2/builtin"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/ctxkey"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
)

// ConfigTool 加载框架未覆盖的业务配置（如管理员邮箱），以 ModeTool 暴露 BambooConfig。
type ConfigTool struct {
	mu  sync.Mutex
	cfg *base.BambooConfig
}

// NewConfigTool 构造业务配置 Tool 插件信息与纯构造工厂。
func NewConfigTool() (application.PluginInfo, func() application.Plugin) {
	info := application.PluginInfo{
		Key:      constants.PluginConfigKey,
		Name:     "bamboo-config-tool",
		Mode:     application.ModeTool,
		Requires: []ctxkey.Key{builtin.LoggerKey},
		Provides: []ctxkey.Key{constants.ContextCustomConfig},
	}
	factory := func() application.Plugin {
		return &ConfigTool{}
	}
	return info, factory
}

// Start 在受管启动阶段从环境变量加载业务配置。
func (t *ConfigTool) Start(ctx context.Context) *apperr.Error {
	logger, err := application.Get[*logging.Logger](ctx, builtin.LoggerKey)
	if err != nil {
		return err
	}
	logger.WithName(logging.NamedINIT).Info(ctx, "加载邮件业务环境变量配置")

	cfg := &base.BambooConfig{
		Email: base.EmailConfig{
			AdminEmail: getEnvStringByKey("EMAIL_ADMIN_EMAIL", ""),
		},
	}

	t.mu.Lock()
	t.cfg = cfg
	t.mu.Unlock()
	return nil
}

// Stop 清理配置工具持有状态。
func (t *ConfigTool) Stop(_ context.Context) *apperr.Error {
	return nil
}

// Expose 向共享容器交付业务配置实例。
func (t *ConfigTool) Expose() map[ctxkey.Key]any {
	t.mu.Lock()
	defer t.mu.Unlock()
	return map[ctxkey.Key]any{
		constants.ContextCustomConfig: t.cfg,
	}
}

func getEnvStringByKey(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
