package core

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	anisdk "github.com/kubercloud/ani-sdks/core-go/anisdk"
	"github.com/kubercloud/ani/services/platform-settings-service/internal/repo/ports"
)

func TestCorePlatformUserClient_List_ParamsAndDecode(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/admin/platform-users" {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("limit") != "10" || q.Get("cursor") != "c1" || q.Get("role") != "platform-ops" ||
			q.Get("status") != "active" || q.Get("source") != "local" || q.Get("search") != "ops" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{
				{
					"id": "11111111-1111-1111-1111-111111111111", "email": "a@example.com",
					"username": "local:ops", "display_name": "Ops", "role": "platform-ops",
					"status": "active", "source": "local", "created_at": "2026-08-31T00:00:00Z",
				},
			},
			"next_cursor": "c2",
		})
	}))
	defer srv.Close()

	client := &CorePlatformUserClient{sdk: anisdk.NewClient(strings.TrimRight(srv.URL, "/")+"/api/v1", "")}
	out, err := client.List(context.Background(), ports.PlatformUserListFilter{
		Limit: 10, Cursor: "c1", Role: "platform-ops", Status: "active", Source: "local", Search: "ops",
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(out.Items) != 1 || out.NextCursor != "c2" || out.Items[0].Username != "local:ops" {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestCorePlatformUserClient_Create(t *testing.T) {
	t.Parallel()

	id := "22222222-2222-2222-2222-222222222222"
	var sawCreate, sawGet bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/admin/platform-users":
			sawCreate = true
			_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "message": "ok"})
		case r.Method == http.MethodGet:
			sawGet = true
			t.Fatalf("Create must not call GET: %s %s", r.Method, r.URL.Path)
		default:
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	client := &CorePlatformUserClient{sdk: anisdk.NewClient(strings.TrimRight(srv.URL, "/")+"/api/v1", "")}
	out, err := client.Create(context.Background(), ports.PlatformUserCreateInput{
		Email: "b@example.com", Username: "admin", DisplayName: "Admin", Role: "platform-admin", Password: "Abcd1234!",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if !sawCreate || sawGet || out != id {
		t.Fatalf("create=%v get=%v out=%q", sawCreate, sawGet, out)
	}
}

func TestMapSDKError_BusinessCodes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		code string
		want error
	}{
		{"PLATFORM_USER_NOT_FOUND", ports.ErrPlatformUserNotFound},
		{"ROLE_NOT_FOUND", ports.ErrRoleNotFound},
		{"USERNAME_ALREADY_EXISTS", ports.ErrUsernameAlreadyExists},
		{"LAST_PLATFORM_ADMIN", ports.ErrLastPlatformAdmin},
		{"PASSWORD_SAME_AS_OLD", ports.ErrPasswordSameAsOld},
		{"ROLE_CHANGE_INVALID", ports.ErrRoleChangeInvalid},
		{"VALIDATION_FAILED", ports.ErrValidationFailed},
		{"NOT_FOUND", ports.ErrCoreUnavailable},
	}
	for _, tc := range cases {
		t.Run(tc.code, func(t *testing.T) {
			err := mapSDKError(anisdk.NewAPIError(tc.code, "detail", "req", nil))
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestCorePlatformUserClient_Get_Success(t *testing.T) {
	t.Parallel()

	id := "33333333-3333-3333-3333-333333333333"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/admin/platform-users/"+id {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": id, "email": "ops@ani.io", "username": "local:ops", "display_name": "Ops",
			"role": "platform-ops", "status": "active", "source": "local",
			"created_at": "2026-08-31T00:00:00Z",
		})
	}))
	defer srv.Close()

	client := &CorePlatformUserClient{sdk: anisdk.NewClient(strings.TrimRight(srv.URL, "/")+"/api/v1", "")}
	out, err := client.Get(context.Background(), uuid.MustParse(id))
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if out.ID != id || out.Username != "local:ops" || out.Role != "platform-ops" {
		t.Fatalf("out=%+v", out)
	}
}

func TestCorePlatformUserClient_Get_404CoreUnavailable(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{"code": "NOT_FOUND", "message": "no route"})
	}))
	defer srv.Close()

	client := &CorePlatformUserClient{sdk: anisdk.NewClient(strings.TrimRight(srv.URL, "/")+"/api/v1", "")}
	_, err := client.Get(context.Background(), uuid.MustParse("33333333-3333-3333-3333-333333333333"))
	if err == nil || !errors.Is(err, ports.ErrCoreUnavailable) {
		t.Fatalf("expected ErrCoreUnavailable, got %v", err)
	}
}

func TestCorePlatformUserClient_ListPlatformRoles_TODO(t *testing.T) {
	t.Parallel()
	client := &CorePlatformUserClient{}
	_, err := client.ListPlatformRoles(context.Background())
	if !errors.Is(err, ports.ErrNotImplemented) {
		t.Fatalf("got %v", err)
	}
}
