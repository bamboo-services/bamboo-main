package apiLink

import (
	"encoding/json"
	"testing"
)

// TestNullableSnowflakeID_TriState 验证省略 / null / 值的三态解析语义
func TestNullableSnowflakeID_TriState(t *testing.T) {
	type req struct {
		GroupID NullableSnowflakeID `json:"link_group_id"`
		ColorID NullableSnowflakeID `json:"link_color_id"`
	}

	// 省略：两个字段都不出现
	var r1 req
	if err := json.Unmarshal([]byte(`{"link_name":"x"}`), &r1); err != nil {
		t.Fatal(err)
	}
	if r1.GroupID.Provided() || r1.ColorID.Provided() {
		t.Fatal("omitted fields should not be Provided")
	}

	// null：显式清空
	var r2 req
	if err := json.Unmarshal([]byte(`{"link_group_id":null,"link_color_id":null}`), &r2); err != nil {
		t.Fatal(err)
	}
	if !r2.GroupID.Provided() || r2.GroupID.Value() != nil {
		t.Fatal("null group should be Provided with nil Value")
	}
	if !r2.ColorID.Provided() || r2.ColorID.Value() != nil {
		t.Fatal("null color should be Provided with nil Value")
	}

	// 值：字符串与数字格式均兼容
	var r3 req
	if err := json.Unmarshal([]byte(`{"link_group_id":"1001","link_color_id":1002}`), &r3); err != nil {
		t.Fatal(err)
	}
	if !r3.GroupID.Provided() || r3.GroupID.Value() == nil || *r3.GroupID.Value() != 1001 {
		t.Fatalf("string value parse failed: %+v", r3.GroupID)
	}
	if !r3.ColorID.Provided() || r3.ColorID.Value() == nil || *r3.ColorID.Value() != 1002 {
		t.Fatalf("number value parse failed: %+v", r3.ColorID)
	}
}
