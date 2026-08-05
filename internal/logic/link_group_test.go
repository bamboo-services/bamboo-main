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

package logic

import (
	"testing"

	"github.com/bamboo-services/bamboo-main/internal/entity"
)

// TestBuiltinGroupValues 验证内置分组值切片保持固定顺序（首页 → 友链页）且不共享指针。
func TestBuiltinGroupValues(t *testing.T) {
	values := builtinGroupValues()
	if len(values) != 2 {
		t.Fatalf("内置分组数量应为 2，实际 %d", len(values))
	}
	if values[0].Name != "首页" || values[1].Name != "友链页" {
		t.Fatalf("内置分组顺序错误：%q → %q", values[0].Name, values[1].Name)
	}
	if len(entity.NewBuiltinGroups()) != len(values) {
		t.Fatal("值切片应复制虚拟对象而非共享引用")
	}
}

// TestBuiltinGroupsVisible 验证内置分组注入的过滤可见性谓词。
func TestBuiltinGroupsVisible(t *testing.T) {
	enabled := 1
	disabled := 0
	onlyOn := true
	onlyOff := false
	nameHit := "首页"
	nameHitSub := "友链"
	nameMiss := "技术"

	cases := []struct {
		name        string
		status      *int
		groupName   *string
		onlyEnabled *bool
		want        bool
	}{
		{"无过滤恒注入", nil, nil, nil, true},
		{"启用过滤命中", &enabled, nil, nil, true},
		{"禁用过滤排除", &disabled, nil, nil, false},
		{"仅启用命中", nil, nil, &onlyOn, true},
		{"不启用过滤排除", nil, nil, &onlyOff, false},
		{"名称精确命中内置", nil, &nameHit, nil, true},
		{"名称模糊命中内置", nil, &nameHitSub, nil, true},
		{"名称未命中排除", nil, &nameMiss, nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := builtinGroupsVisible(c.status, c.groupName, c.onlyEnabled); got != c.want {
				t.Fatalf("builtinGroupsVisible(%v, %v, %v) = %v，期望 %v",
					c.status, c.groupName, c.onlyEnabled, got, c.want)
			}
		})
	}
}
