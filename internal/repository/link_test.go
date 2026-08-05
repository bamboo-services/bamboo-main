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
	"reflect"
	"strings"
	"testing"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// TestBuildFailureUpdates 校验失效状态更新字段组装：
// 标记失效归入内置「已失效」分组，恢复时以单条 CASE 原子清空已失效分组，
// 未传入内置分组 ID 时不操作分组。
func TestBuildFailureUpdates(t *testing.T) {
	invalidID := xSnowflake.SnowflakeID(constants.BuiltinGroupInvalidID)

	cases := []struct {
		name         string
		isFailure    int
		failReason   string
		invalidGroup *xSnowflake.SnowflakeID
		wantGroupKey bool
		wantGroupVal any
	}{
		{
			name:         "标记失效归入已失效分组",
			isFailure:    constants.LinkFailBroken,
			failReason:   "站点不可访问",
			invalidGroup: &invalidID,
			wantGroupKey: true,
			wantGroupVal: invalidID,
		},
		{
			name:         "恢复以 CASE 清空已失效分组",
			isFailure:    constants.LinkFailNormal,
			invalidGroup: &invalidID,
			wantGroupKey: true,
			wantGroupVal: clause.Expr{
				SQL:  "CASE WHEN group_id = ? THEN NULL ELSE group_id END",
				Vars: []interface{}{invalidID},
			},
		},
		{
			name:         "未传分组 ID 不操作分组",
			isFailure:    constants.LinkFailBroken,
			invalidGroup: nil,
			wantGroupKey: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			updates := buildFailureUpdates(tc.isFailure, tc.failReason, tc.invalidGroup)

			if got, ok := updates["is_failure"]; !ok || got != tc.isFailure {
				t.Fatalf("is_failure = %v（存在=%v），期望 %v", got, ok, tc.isFailure)
			}
			if got, ok := updates["fail_reason"]; !ok || got != tc.failReason {
				t.Fatalf("fail_reason = %v（存在=%v），期望 %q", got, ok, tc.failReason)
			}

			got, hasGroup := updates["group_id"]
			if tc.wantGroupKey {
				if !hasGroup {
					t.Fatal("应包含 group_id 键")
				}
				wantExpr, wantIsExpr := tc.wantGroupVal.(clause.Expr)
				if wantIsExpr {
					gotExpr, gotIsExpr := got.(clause.Expr)
					if !gotIsExpr {
						t.Fatalf("group_id 应为 clause.Expr，实际 %T（%v）", got, got)
					}
					if gotExpr.SQL != wantExpr.SQL {
						t.Fatalf("CASE SQL = %q，期望 %q", gotExpr.SQL, wantExpr.SQL)
					}
					if len(gotExpr.Vars) != 1 || gotExpr.Vars[0] != wantExpr.Vars[0] {
						t.Fatalf("CASE 参数 = %v，期望 %v", gotExpr.Vars, wantExpr.Vars)
					}
				} else if got != tc.wantGroupVal {
					t.Fatalf("group_id = %v，期望 %v", got, tc.wantGroupVal)
				}
			} else if hasGroup {
				t.Fatalf("不应包含 group_id 键，实际 %v", got)
			}
		})
	}
}

// TestBuildEditResolveUpdates 校验修改位置/颜色申请的审核落库字段组装：
// 通过（apply=true）时以 COALESCE 将 expected_* 原子应用到正式字段；
// 拒绝（apply=false）时保持原值；两种情形均清空预期值、状态回已通过并写入审核备注。
func TestBuildEditResolveUpdates(t *testing.T) {
	wantGroupExpr := gorm.Expr("COALESCE(expected_group_id, group_id)")
	wantColorExpr := gorm.Expr("COALESCE(expected_color_id, color_id)")

	cases := []struct {
		name        string
		apply       bool
		remark      string
		wantApply   bool // 是否包含 group_id/color_id 键
		wantGroupID any
		wantColorID any
	}{
		{
			name:        "通过：应用预期值",
			apply:       true,
			remark:      "已调整",
			wantApply:   true,
			wantGroupID: wantGroupExpr,
			wantColorID: wantColorExpr,
		},
		{
			name:      "拒绝：保持原值",
			apply:     false,
			remark:    "维持现状",
			wantApply: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildEditResolveUpdates(tc.apply, tc.remark)

			// 恒清空预期值、状态回已通过、审核备注透传
			if _, ok := got["expected_group_id"]; !ok || got["expected_group_id"] != nil {
				t.Fatalf("期望 expected_group_id 键存在且为 nil，实际: %#v", got["expected_group_id"])
			}
			if _, ok := got["expected_color_id"]; !ok || got["expected_color_id"] != nil {
				t.Fatalf("期望 expected_color_id 键存在且为 nil，实际: %#v", got["expected_color_id"])
			}
			if got["status"] != constants.LinkStatusApproved {
				t.Fatalf("期望状态回已通过(%d)，实际: %v", constants.LinkStatusApproved, got["status"])
			}
			if got["review_remark"] != tc.remark {
				t.Fatalf("期望审核备注 %q，实际: %v", tc.remark, got["review_remark"])
			}

			_, hasGroup := got["group_id"]
			_, hasColor := got["color_id"]
			if hasGroup != tc.wantApply || hasColor != tc.wantApply {
				t.Fatalf("期望 group_id/color_id 键存在性为 %v/%v，实际 %v/%v", tc.wantApply, tc.wantApply, hasGroup, hasColor)
			}
			if tc.wantApply {
				if !reflect.DeepEqual(got["group_id"], tc.wantGroupID) {
					t.Fatalf("期望 group_id 为 %#v，实际: %#v", tc.wantGroupID, got["group_id"])
				}
				if !reflect.DeepEqual(got["color_id"], tc.wantColorID) {
					t.Fatalf("期望 color_id 为 %#v，实际: %#v", tc.wantColorID, got["color_id"])
				}
			}
		})
	}
}
