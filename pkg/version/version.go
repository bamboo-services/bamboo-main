// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package version

import (
	"strings"

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/env"
)

// Version 为构建期通过 -ldflags 注入的应用版本号，本地开发默认 "dev"。
var Version = "dev"

// Resolve 返回当前应用版本号（优先读取环境变量 APP_VERSION，其次回退构建期注入值）。
func Resolve() string {
	if val, ok := env.Lookup(env.AppVersion); ok && !val.IsEmpty() {
		if v := strings.TrimSpace(val.String()); v != "" {
			return v
		}
	}
	if v := strings.TrimSpace(Version); v != "" {
		return v
	}
	return "dev"
}

// ResolveEnvironment 返回当前运行环境标识（优先读取 APP_ENV，其次按 APP_DEBUG 推导）。
func ResolveEnvironment() string {
	if val, ok := env.Lookup(env.Env); ok && !val.IsEmpty() {
		if e := strings.TrimSpace(val.String()); e != "" {
			return e
		}
	}
	if val, ok := env.Lookup(env.Debug); ok && !val.IsEmpty() {
		if debug, err := val.Bool(); err == nil && debug {
			return "development"
		}
	}
	return "production"
}
