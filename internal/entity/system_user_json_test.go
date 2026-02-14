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
