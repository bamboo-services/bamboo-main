// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package repository

import (
	"strings"
	"testing"
)

// TestSortRowTemplate 校验各重排 VALUES 行模板的占位符数与类型标注，将
// 「占位符与 args append 一一对应」的隐式耦合固化为可回归不变量。
func TestSortRowTemplate(t *testing.T) {
	cases := []struct {
		name   string
		tpl    string
		fields int
	}{
		{"友链(id, sort_order, group_id)", linkSortRowTemplate, 3},
		{"分组(id, sort_order)", groupSortRowTemplate, 2},
		{"颜色(id, sort_order)", colorSortRowTemplate, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := strings.Count(tc.tpl, "?"); got != tc.fields {
				t.Fatalf("期望 %d 个占位符，实际 %d（模板: %s）", tc.fields, got, tc.tpl)
			}
			if got := strings.Count(tc.tpl, "?::"); got != tc.fields {
				t.Fatalf("期望每个占位符均显式 cast（?:: 出现 %d 次，模板: %s）", got, tc.tpl)
			}
		})
	}
}
