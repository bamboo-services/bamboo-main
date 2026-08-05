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

// TestNewBuiltinGroup 验证内置「已失效」分组构造：ID 为保留值、恒启用。
func TestNewBuiltinGroup(t *testing.T) {
	group := NewBuiltinGroup("失联", nil)
	if group.ID != xSnowflake.SnowflakeID(1) {
		t.Fatalf("内置分组 ID 应为 1，实际 %v", group.ID)
	}
	if group.Name != "失联" {
		t.Errorf("内置分组名称 = %q，期望 %q", group.Name, "失联")
	}
	if group.SortOrder != 0 {
		t.Errorf("内置分组排序 = %d，期望 0", group.SortOrder)
	}
	if !group.Status {
		t.Errorf("内置分组应恒为启用")
	}
}

// TestNewDefaultBuiltinGroup 验证内置「已失效」分组默认值兜底。
func TestNewDefaultBuiltinGroup(t *testing.T) {
	group := NewDefaultBuiltinGroup()
	if group.ID != xSnowflake.SnowflakeID(1) {
		t.Fatalf("默认内置分组 ID 应为 1，实际 %v", group.ID)
	}
	if group.Name != "已失效" {
		t.Errorf("默认内置分组名称 = %q，期望 %q", group.Name, "已失效")
	}
	if group.Description != nil {
		t.Errorf("默认内置分组描述应为 nil，实际 %v", group.Description)
	}
}

// TestBuiltinGroupByID 验证按保留 ID 查找内置分组，非内置 ID 返回 nil。
func TestBuiltinGroupByID(t *testing.T) {
	if got := BuiltinGroupByID(1); got == nil || got.Name != "已失效" {
		t.Fatalf("BuiltinGroupByID(1) 应返回已失效分组，实际 %v", got)
	}
	if got := BuiltinGroupByID(2); got != nil {
		t.Fatalf("BuiltinGroupByID(2) 应为 nil，实际 %v", got)
	}
	if got := BuiltinGroupByID(3); got != nil {
		t.Fatalf("BuiltinGroupByID(3) 应为 nil，实际 %v", got)
	}
	if got := BuiltinGroupByID(xSnowflake.SnowflakeID(0)); got != nil {
		t.Fatalf("BuiltinGroupByID(0) 应为 nil，实际 %v", got)
	}
}

// TestIsBuiltinGroupID 验证内置「已失效」分组保留 ID 判定。
func TestIsBuiltinGroupID(t *testing.T) {
	cases := []struct {
		id   xSnowflake.SnowflakeID
		want bool
	}{
		{1, true},
		{2, false},
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
