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
	"encoding/json"
	"testing"
)

func TestSystemUserJSON_HidesSensitiveFields(t *testing.T) {
	oauthID := "oauth-uid"
	user := SystemUser{
		ID:          1,
		OAuthUserID: &oauthID,
		Username:    "admin",
		Password:    "hashed-password",
		Email:       "admin@example.com",
		Role:        "admin",
		Status:      1,
		EmailVerify: true,
	}

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("marshal user failed: %v", err)
	}

	var payload map[string]any
	if err = json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal user json failed: %v", err)
	}

	if _, exists := payload["password"]; exists {
		t.Fatal("expected password to be hidden in json output")
	}
	if _, exists := payload["oauth_user_id"]; exists {
		t.Fatal("expected oauth_user_id to be hidden in json output")
	}
	if payload["username"] != "admin" {
		t.Fatalf("expected username=admin, got=%v", payload["username"])
	}
	if payload["email"] != "admin@example.com" {
		t.Fatalf("expected email=admin@example.com, got=%v", payload["email"])
	}
}
