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
)

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

	debug := false
	if val, ok := env.Lookup(env.Debug); ok && !val.IsEmpty() {
		if b, err := val.Bool(); err == nil {
			debug = b
		}
	}

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
