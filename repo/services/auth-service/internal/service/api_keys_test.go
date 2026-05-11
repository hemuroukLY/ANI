package service

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAPIKeyEmbedsTenantID(t *testing.T) {
	tenantID := uuid.New()
	key, err := generateAPIKey(tenantID)
	if err != nil {
		t.Fatalf("generateAPIKey: %v", err)
	}
	if !strings.HasPrefix(key, "ani_dev_"+tenantID.String()+"_") {
		t.Fatalf("key prefix = %q, want tenant embedded", key)
	}
	gotTenantID, err := parseAPIKeyTenant(key)
	if err != nil {
		t.Fatalf("parseAPIKeyTenant: %v", err)
	}
	if gotTenantID != tenantID {
		t.Fatalf("tenant id = %s, want %s", gotTenantID, tenantID)
	}
}

func TestHasScope(t *testing.T) {
	if !hasScope([]string{"scope:models:create"}, "models", "create") {
		t.Fatal("expected exact scope to allow")
	}
	if !hasScope([]string{"models:*"}, "models", "delete") {
		t.Fatal("expected resource wildcard scope to allow")
	}
	if hasScope([]string{"scope:tasks:get"}, "models", "get") {
		t.Fatal("unexpected scope allow")
	}
}
