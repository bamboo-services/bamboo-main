// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package version

import "testing"

func TestResolve(t *testing.T) {
	origVersion := Version
	t.Cleanup(func() {
		Version = origVersion
	})

	t.Run("优先使用 APP_VERSION 环境变量", func(t *testing.T) {
		Version = "v1.1.0-rc.1"
		t.Setenv("APP_VERSION", "v1.1.0-rc.2")
		if got := Resolve(); got != "v1.1.0-rc.2" {
			t.Fatalf("Resolve() = %q，期望 %q", got, "v1.1.0-rc.2")
		}
	})

	t.Run("回退构建期注入版本", func(t *testing.T) {
		Version = "v1.1.0-rc.1"
		t.Setenv("APP_VERSION", "")
		if got := Resolve(); got != "v1.1.0-rc.1" {
			t.Fatalf("Resolve() = %q，期望 %q", got, "v1.1.0-rc.1")
		}
	})

	t.Run("空白注入回退 dev", func(t *testing.T) {
		Version = "   "
		t.Setenv("APP_VERSION", "")
		if got := Resolve(); got != "dev" {
			t.Fatalf("Resolve() = %q，期望 %q", got, "dev")
		}
	})
}

func TestResolveEnvironment(t *testing.T) {
	t.Run("优先读取 APP_ENV", func(t *testing.T) {
		t.Setenv("APP_ENV", "staging")
		t.Setenv("APP_DEBUG", "true")
		if got := ResolveEnvironment(); got != "staging" {
			t.Fatalf("ResolveEnvironment() = %q，期望 %q", got, "staging")
		}
	})

	t.Run("APP_DEBUG=true 推导 development", func(t *testing.T) {
		t.Setenv("APP_ENV", "")
		t.Setenv("APP_DEBUG", "true")
		if got := ResolveEnvironment(); got != "development" {
			t.Fatalf("ResolveEnvironment() = %q，期望 %q", got, "development")
		}
	})

	t.Run("APP_DEBUG=false 或未设置推导 production", func(t *testing.T) {
		t.Setenv("APP_ENV", "")
		t.Setenv("APP_DEBUG", "false")
		if got := ResolveEnvironment(); got != "production" {
			t.Fatalf("ResolveEnvironment() = %q，期望 %q", got, "production")
		}
	})
}
