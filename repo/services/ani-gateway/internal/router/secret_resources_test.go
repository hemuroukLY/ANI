package router

import (
	"context"
	"testing"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestSecretAPIDevProfileIdempotencyAndBinding(t *testing.T) {
	api := newSecretAPI()
	a, err := api.service.CreateSecret(context.Background(), ports.SecretCreateRequest{
		TenantID:       "t1",
		IdempotencyKey: "secret-1",
		Name:           "db-password",
		Type:           "opaque",
		Data:           map[string]string{"password": "secret-value", "username": "ani"},
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err := api.service.CreateSecret(context.Background(), ports.SecretCreateRequest{
		TenantID:       "t1",
		IdempotencyKey: "secret-1",
		Name:           "db-password",
		Type:           "opaque",
		Data:           map[string]string{"password": "secret-value", "username": "ani"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if a.SecretID != b.SecretID {
		t.Fatalf("want idempotent secret id, got %s != %s", a.SecretID, b.SecretID)
	}
	if len(a.Keys) != 2 || a.Keys[0] != "password" || a.Keys[1] != "username" {
		t.Fatalf("want sorted secret keys without values, got %#v", a.Keys)
	}
	resp := secretFromRecord(a)
	requireLocalCoreDevProfile(t, resp.DevProfile, "local-secret-service")

	binding, err := api.service.BindSecret(context.Background(), ports.SecretBindRequest{
		TenantID:   "t1",
		SecretID:   a.SecretID,
		TargetType: "instance",
		TargetID:   "inst-1",
		EnvPrefix:  "DB_",
	})
	if err != nil {
		t.Fatal(err)
	}
	if binding.State != "bound" {
		t.Fatalf("want bound state, got %s", binding.State)
	}
	bindingResp := secretBindingFromRecord(binding)
	requireLocalCoreDevProfile(t, bindingResp.DevProfile, "local-secret-service")
}
