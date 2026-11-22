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
	"log/slog"
	"os"
	"strings"

	"github.com/bamboo-services/bamboo-base-go/application/v2/builtin"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/env"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func isAppDebug() bool {
	if val, ok := env.Lookup(env.Debug); ok && !val.IsEmpty() {
		if b, err := val.Bool(); err == nil {
			return b
		}
	}
	return false
}

// ConfigureGinMode 根据 APP_DEBUG 同步 Gin 全局运行模式。
// 若已显式设置 GIN_MODE 环境变量，则尊重用户设置不予覆盖。
func ConfigureGinMode() {
	if raw, ok := os.LookupEnv(gin.EnvGinMode); ok && strings.TrimSpace(raw) != "" {
		return
	}
	if isAppDebug() {
		gin.SetMode(gin.DebugMode)
		return
	}
	gin.SetMode(gin.ReleaseMode)
}

// BuildGormConfig 从环境变量装配 GORM 配置，固定开启单数表名策略以兼容存量数据表。
func BuildGormConfig() *gorm.Config {
	tablePrefix := "bm_"
	if val, ok := env.Lookup(env.DatabasePrefix); ok && !val.IsEmpty() {
		tablePrefix = val.String()
	}
	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
	}
}

// BuildLoggerConfig 从环境变量装配内置受管日志组件配置。
func BuildLoggerConfig() builtin.LoggerConfig {
	dir := "logs"
	if val, ok := env.Lookup(env.LogPath); ok && !val.IsEmpty() {
		dir = val.String()
	}

	maxSizeMB := 10
	if val, ok := env.Lookup(env.LogMaxSize); ok && !val.IsEmpty() {
		if n, err := val.Int(); err == nil && n > 0 {
			maxSizeMB = n
		}
	}

	debug := isAppDebug()

	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	if val, ok := env.Lookup(env.LogLevel); ok && !val.IsEmpty() {
		switch strings.ToLower(strings.TrimSpace(val.String())) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn", "warning":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}

	return builtin.LoggerConfig{
		Rotator: logging.RotatorConfig{
			Dir:      dir,
			BaseName: "bamboo-main",
			Ext:      ".log",
			MaxSize:  int64(maxSizeMB) * 1024 * 1024,
		},
		Console: os.Stdout,
		Level:   level,
		Debug:   debug,
	}
}

// ParseSnowflakeIDs 从环境变量读取雪花数据中心 ID 与节点 ID（默认 0 / 1）。
func ParseSnowflakeIDs() (datacenterID int64, nodeID int64) {
	datacenterID = 0
	if val, ok := env.Lookup(env.SnowflakeDatacenterID); ok && !val.IsEmpty() {
		if n, err := val.Int64(); err == nil {
			datacenterID = n
		}
	}
	nodeID = 1
	if val, ok := env.Lookup(env.SnowflakeNodeID); ok && !val.IsEmpty() {
		if n, err := val.Int64(); err == nil {
			nodeID = n
		}
	}
	return datacenterID, nodeID
}
