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
	"sync"

	"github.com/bamboo-services/bamboo-base-go/application/v2"
	"github.com/bamboo-services/bamboo-base-go/application/v2/builtin"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/ctxkey"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/plugins/database/v2"
	"github.com/bamboo-services/bamboo-main/internal/app/startup/prepare"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
)

// MigrationTool 在 Database 与 Snowflake 就绪后执行 AutoMigrate 与幂等种子数据初始化。
type MigrationTool struct {
	models   []any
	mu       sync.Mutex
	migrated bool
}

// NewMigrationTool 构造数据库迁移与种子初始化 Tool 插件。
func NewMigrationTool(models ...any) (application.PluginInfo, func() application.Plugin) {
	modelsCopy := make([]any, len(models))
	copy(modelsCopy, models)

	info := application.PluginInfo{
		Key:      constants.PluginMigrationKey,
		Name:     "bamboo-migration-tool",
		Mode:     application.ModeTool,
		Requires: []ctxkey.Key{ctxkey.DatabaseKey, ctxkey.SnowflakeNodeKey, builtin.LoggerKey},
		Provides: []ctxkey.Key{constants.ContextDatabaseMigrated},
	}
	factory := func() application.Plugin {
		return &MigrationTool{models: modelsCopy}
	}
	return info, factory
}

// Start 执行 GORM 表结构自动迁移与默认管理员/站点配置种子注入。
func (t *MigrationTool) Start(ctx context.Context) *apperr.Error {
	db, appErr := database.Get(ctx)
	if appErr != nil {
		return appErr
	}
	logger, appErr := application.Get[*logging.Logger](ctx, builtin.LoggerKey)
	if appErr != nil {
		return appErr
	}

	log := logger.WithName(logging.NamedINIT)
	log.Info(ctx, "执行数据库表结构迁移与种子初始化")

	if len(t.models) > 0 {
		if err := db.WithContext(ctx).AutoMigrate(t.models...); err != nil {
			return apperr.New(database.DatabaseError, "数据库自动迁移失败", err)
		}
	}

	if err := prepare.DefaultUser(ctx, db, logger); err != nil {
		return apperr.From(database.DatabaseError, "默认管理员种子初始化失败", err)
	}
	if err := prepare.DefaultInfo(ctx, db, logger); err != nil {
		return apperr.From(database.DatabaseError, "默认站点配置种子初始化失败", err)
	}

	t.mu.Lock()
	t.migrated = true
	t.mu.Unlock()
	return nil
}

// Stop 释放迁移工具状态。
func (t *MigrationTool) Stop(_ context.Context) *apperr.Error {
	return nil
}

// Expose 向共享容器发布迁移完成标记。
func (t *MigrationTool) Expose() map[ctxkey.Key]any {
	t.mu.Lock()
	defer t.mu.Unlock()
	return map[ctxkey.Key]any{
		constants.ContextDatabaseMigrated: t.migrated,
	}
}
