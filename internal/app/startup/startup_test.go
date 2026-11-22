// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package startup

import (
	"sync"
	"testing"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/schema"
)

// TestBuildGormConfigSingularTableNames 校验 GORM 命名策略固定开启单数表名，
// 确保全部持久化实体在默认前缀下解析为存量单数表名，防止升级漂移回复数表。
func TestBuildGormConfigSingularTableNames(t *testing.T) {
	t.Setenv("DATABASE_PREFIX", "")

	cfg := BuildGormConfig()
	ns, ok := cfg.NamingStrategy.(schema.NamingStrategy)
	if !ok {
		t.Fatalf("NamingStrategy 类型 = %T，期望 schema.NamingStrategy", cfg.NamingStrategy)
	}
	if !ns.SingularTable {
		t.Fatal("SingularTable 应为 true")
	}
	if ns.TablePrefix != "bm_" {
		t.Fatalf("默认 TablePrefix = %q，期望 %q", ns.TablePrefix, "bm_")
	}

	cases := []struct {
		model any
		want  string
	}{
		{&entity.SystemUser{}, "bm_system_user"},
		{&entity.LinkGroup{}, "bm_link_group"},
		{&entity.LinkColor{}, "bm_link_color"},
		{&entity.LinkFriend{}, "bm_link_friend"},
		{&entity.SystemLog{}, "bm_system_log"},
		{&entity.System{}, "bm_system"},
		{&entity.SponsorChannel{}, "bm_sponsor_channel"},
		{&entity.SponsorRecord{}, "bm_sponsor_record"},
	}

	cacheStore := &sync.Map{}
	for _, tc := range cases {
		parsed, err := schema.Parse(tc.model, cacheStore, cfg.NamingStrategy)
		if err != nil {
			t.Fatalf("解析实体 %T 失败: %v", tc.model, err)
		}
		if parsed.Table != tc.want {
			t.Errorf("实体 %T 表名 = %q，期望 %q", tc.model, parsed.Table, tc.want)
		}
	}
}

// TestBuildGormConfigCustomPrefix 校验 DATABASE_PREFIX 非空时可自定义前缀且保持单数表名。
func TestBuildGormConfigCustomPrefix(t *testing.T) {
	t.Setenv("DATABASE_PREFIX", "custom_")

	cfg := BuildGormConfig()
	parsed, err := schema.Parse(&entity.LinkFriend{}, &sync.Map{}, cfg.NamingStrategy)
	if err != nil {
		t.Fatalf("解析实体失败: %v", err)
	}
	if parsed.Table != "custom_link_friend" {
		t.Fatalf("表名 = %q，期望 %q", parsed.Table, "custom_link_friend")
	}
}

// TestConfigureGinMode 校验 Gin 运行模式跟随 APP_DEBUG 切换，且不覆盖显式 GIN_MODE。
func TestConfigureGinMode(t *testing.T) {
	prevMode := gin.Mode()
	t.Cleanup(func() {
		gin.SetMode(prevMode)
	})

	t.Run("APP_DEBUG=false 切换为 release", func(t *testing.T) {
		t.Setenv(gin.EnvGinMode, "")
		t.Setenv("APP_DEBUG", "false")
		gin.SetMode(gin.DebugMode)

		ConfigureGinMode()
		if got := gin.Mode(); got != gin.ReleaseMode {
			t.Fatalf("gin.Mode() = %q，期望 %q", got, gin.ReleaseMode)
		}
	})

	t.Run("APP_DEBUG=true 切换为 debug", func(t *testing.T) {
		t.Setenv(gin.EnvGinMode, "")
		t.Setenv("APP_DEBUG", "true")
		gin.SetMode(gin.ReleaseMode)

		ConfigureGinMode()
		if got := gin.Mode(); got != gin.DebugMode {
			t.Fatalf("gin.Mode() = %q，期望 %q", got, gin.DebugMode)
		}
	})

	t.Run("未设置 APP_DEBUG 默认 release", func(t *testing.T) {
		t.Setenv(gin.EnvGinMode, "")
		t.Setenv("APP_DEBUG", "")
		gin.SetMode(gin.DebugMode)

		ConfigureGinMode()
		if got := gin.Mode(); got != gin.ReleaseMode {
			t.Fatalf("gin.Mode() = %q，期望 %q", got, gin.ReleaseMode)
		}
	})

	t.Run("显式 GIN_MODE 不被覆盖", func(t *testing.T) {
		t.Setenv(gin.EnvGinMode, gin.TestMode)
		t.Setenv("APP_DEBUG", "false")
		gin.SetMode(gin.TestMode)

		ConfigureGinMode()
		if got := gin.Mode(); got != gin.TestMode {
			t.Fatalf("gin.Mode() = %q，期望 %q", got, gin.TestMode)
		}
	})
}
