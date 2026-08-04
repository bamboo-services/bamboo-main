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
	"encoding/json"
	"testing"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xModels "github.com/bamboo-services/bamboo-base-go/major/models"
	apiLink "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/repository"
)

// mustItem 经 JSON 反序列化构造排序条目（NullableSnowflakeID 为私有字段，仅能经 JSON 构造三态）
func mustItem(t *testing.T, raw string) apiLink.FriendSortItem {
	t.Helper()
	var item apiLink.FriendSortItem
	if err := json.Unmarshal([]byte(raw), &item); err != nil {
		t.Fatalf("反序列化排序条目 %q 失败: %v", raw, err)
	}
	return item
}

// newLink 构造带指定 ID 与分组的友链实体
func newLink(id int64, groupID *xSnowflake.SnowflakeID) entity.LinkFriend {
	return entity.LinkFriend{
		BaseEntity: xModels.BaseEntity{ID: xSnowflake.SnowflakeID(id)},
		GroupID:    groupID,
	}
}

// groupPtr 返回雪花 ID 指针
func groupPtr(id int64) *xSnowflake.SnowflakeID {
	v := xSnowflake.SnowflakeID(id)
	return &v
}

// sameGroupID 比较两个可空分组 ID（nil 与 nil 视为相等）
func sameGroupID(a, b *xSnowflake.SnowflakeID) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// TestBuildSortAssignments 表驱动验证全局序号 / 跨组 / null / 省略 / 重复 / 空切片六种语义
func TestBuildSortAssignments(t *testing.T) {
	tests := []struct {
		name    string
		links   []entity.LinkFriend
		items   []apiLink.FriendSortItem
		want    []repository.SortAssignment
		wantErr bool
	}{
		{
			name: "全局序号按载荷顺序分配 0..N-1",
			links: []entity.LinkFriend{
				newLink(1001, groupPtr(2001)),
				newLink(1002, groupPtr(2001)),
				newLink(1003, groupPtr(2002)),
			},
			items: []apiLink.FriendSortItem{
				mustItem(t, `{"id":1003}`),
				mustItem(t, `{"id":1001}`),
				mustItem(t, `{"id":1002}`),
			},
			want: []repository.SortAssignment{
				{ID: 1003, GroupID: groupPtr(2002), Order: 0},
				{ID: 1001, GroupID: groupPtr(2001), Order: 1},
				{ID: 1002, GroupID: groupPtr(2001), Order: 2},
			},
		},
		{
			name:  "跨组移动取 item 指定分组",
			links: []entity.LinkFriend{newLink(1001, groupPtr(2001))},
			items: []apiLink.FriendSortItem{mustItem(t, `{"id":1001,"group_id":2002}`)},
			want:  []repository.SortAssignment{{ID: 1001, GroupID: groupPtr(2002), Order: 0}},
		},
		{
			name:  "显式 null 置未分组",
			links: []entity.LinkFriend{newLink(1001, groupPtr(2001))},
			items: []apiLink.FriendSortItem{mustItem(t, `{"id":1001,"group_id":null}`)},
			want:  []repository.SortAssignment{{ID: 1001, GroupID: nil, Order: 0}},
		},
		{
			name:  "省略 group_id 保持原分组",
			links: []entity.LinkFriend{newLink(1001, groupPtr(2001))},
			items: []apiLink.FriendSortItem{mustItem(t, `{"id":1001}`)},
			want:  []repository.SortAssignment{{ID: 1001, GroupID: groupPtr(2001), Order: 0}},
		},
		{
			name:    "重复 ID 返回 error",
			links:   []entity.LinkFriend{newLink(1001, nil), newLink(1002, nil)},
			items:   []apiLink.FriendSortItem{mustItem(t, `{"id":1001}`), mustItem(t, `{"id":1001}`)},
			wantErr: true,
		},
		{
			name:  "空切片正常返回空",
			links: []entity.LinkFriend{},
			items: []apiLink.FriendSortItem{},
			want:  []repository.SortAssignment{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildSortAssignments(tt.links, tt.items)
			if tt.wantErr {
				if err == nil {
					t.Fatal("期望返回 error，实际为 nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("非预期 error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("赋值长度不符：got %d，want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("[%d] ID 不符：got %d，want %d", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Order != tt.want[i].Order {
					t.Errorf("[%d] Order 不符：got %d，want %d", i, got[i].Order, tt.want[i].Order)
				}
				if !sameGroupID(got[i].GroupID, tt.want[i].GroupID) {
					t.Errorf("[%d] GroupID 不符：got %v，want %v", i, got[i].GroupID, tt.want[i].GroupID)
				}
			}
		})
	}
}
