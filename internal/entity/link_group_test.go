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

package entity

import (
	"testing"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
)

// TestNewBuiltinGroups 验证内置分组构造：固定顺序 首页 → 友链页，ID 为保留值、恒启用。
func TestNewBuiltinGroups(t *testing.T) {
	builtins := NewBuiltinGroups()
	if len(builtins) != 2 {
		t.Fatalf("内置分组数量应为 2，实际 %d", len(builtins))
	}

	want := []struct {
		name      string
		id        xSnowflake.SnowflakeID
		sortOrder int
		status    bool
	}{
		{"首页", 1, 0, true},
		{"友链页", 2, 1, true},
	}
	for i, w := range want {
		got := builtins[i]
		if got.ID != w.id {
			t.Errorf("第 %d 个内置分组 ID = %v，期望 %v", i, got.ID, w.id)
		}
		if got.Name != w.name {
			t.Errorf("第 %d 个内置分组名称 = %q，期望 %q", i, got.Name, w.name)
		}
		if got.SortOrder != w.sortOrder {
			t.Errorf("第 %d 个内置分组排序 = %d，期望 %d", i, got.SortOrder, w.sortOrder)
		}
		if !got.Status {
			t.Errorf("第 %d 个内置分组应恒为启用", i)
		}
	}
}

// TestBuiltinGroupByID 验证按保留 ID 查找内置分组，非内置 ID 返回 nil。
func TestBuiltinGroupByID(t *testing.T) {
	if got := BuiltinGroupByID(1); got == nil || got.Name != "首页" {
		t.Fatalf("BuiltinGroupByID(1) 应返回首页分组，实际 %v", got)
	}
	if got := BuiltinGroupByID(2); got == nil || got.Name != "友链页" {
		t.Fatalf("BuiltinGroupByID(2) 应返回友链页分组，实际 %v", got)
	}
	if got := BuiltinGroupByID(3); got != nil {
		t.Fatalf("BuiltinGroupByID(3) 应为 nil，实际 %v", got)
	}
	if got := BuiltinGroupByID(xSnowflake.SnowflakeID(0)); got != nil {
		t.Fatalf("BuiltinGroupByID(0) 应为 nil，实际 %v", got)
	}
}

// TestIsBuiltinGroupID 验证内置分组保留 ID 判定。
func TestIsBuiltinGroupID(t *testing.T) {
	cases := []struct {
		id   xSnowflake.SnowflakeID
		want bool
	}{
		{1, true},
		{2, true},
		{0, false},
		{3, false},
		{xSnowflake.SnowflakeID(1735689600000000000), false}, // 雪花 ID 空间的真实值
	}
	for _, c := range cases {
		if got := IsBuiltinGroupID(c.id); got != c.want {
			t.Errorf("IsBuiltinGroupID(%v) = %v，期望 %v", c.id, got, c.want)
		}
	}
}
