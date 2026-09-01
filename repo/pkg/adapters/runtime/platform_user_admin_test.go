package runtime

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kubercloud/ani/pkg/ports"
)

func TestPostgresPlatformUserAdminStore_Create_Success(t *testing.T) {
	t.Parallel()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	roleID := uuid.MustParse("00000000-0000-0000-0000-000000000006")
	now := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	tx := &quotaFakeTx{}
	tx.enqueueRows(
		quotaFakeRow{values: []any{false}}, // username exists?
		quotaFakeRow{values: []any{roleID}},
		quotaFakeRow{values: []any{userID, now}},
	)
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: tx})

	out, err := store.Create(context.Background(), ports.PlatformUserCreate{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops",
		Role: "platform-ops", PasswordHash: "$2a$12$hashed",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if out.ID != userID || out.Username != "local:ops" || out.Role != "platform-ops" || out.Source != "local" {
		t.Fatalf("out=%+v", out)
	}
	if !hasExec(tx, "INSERT INTO user_roles") {
		t.Fatalf("missing role bind: %#v", tx.execSQLs)
	}
}

func TestPostgresPlatformUserAdminStore_Create_UsernameConflict(t *testing.T) {
	t.Parallel()
	tx := &quotaFakeTx{}
	tx.enqueueRows(quotaFakeRow{values: []any{true}})
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: tx})
	_, err := store.Create(context.Background(), ports.PlatformUserCreate{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops",
		Role: "platform-ops", PasswordHash: "hash",
	})
	if !errors.Is(err, ports.ErrUsernameAlreadyExists) {
		t.Fatalf("err=%v", err)
	}
}

func TestPostgresPlatformUserAdminStore_Create_UniqueViolationMapped(t *testing.T) {
	t.Parallel()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	roleID := uuid.MustParse("00000000-0000-0000-0000-000000000006")
	tx := &quotaFakeTx{}
	tx.enqueueRows(
		quotaFakeRow{values: []any{false}},
		quotaFakeRow{values: []any{roleID}},
		quotaFakeRow{err: &pgconn.PgError{Code: "23505"}},
	)
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: tx})
	_, err := store.Create(context.Background(), ports.PlatformUserCreate{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops",
		Role: "platform-ops", PasswordHash: "hash",
	})
	if !errors.Is(err, ports.ErrUsernameAlreadyExists) {
		t.Fatalf("err=%v", err)
	}
	_ = userID
}

func TestPostgresPlatformUserAdminStore_Create_RoleNotFound(t *testing.T) {
	t.Parallel()
	tx := &quotaFakeTx{}
	tx.enqueueRows(
		quotaFakeRow{values: []any{false}},
		quotaFakeRow{err: pgx.ErrNoRows},
	)
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: tx})
	_, err := store.Create(context.Background(), ports.PlatformUserCreate{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops",
		Role: "tenant-admin", PasswordHash: "hash",
	})
	if !errors.Is(err, ports.ErrRoleNotFound) {
		t.Fatalf("err=%v", err)
	}
}

func TestPostgresPlatformUserAdminStore_Get_NotFound(t *testing.T) {
	t.Parallel()
	tx := &quotaFakeTx{}
	tx.enqueueRows(quotaFakeRow{err: pgx.ErrNoRows})
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: tx})
	_, err := store.Get(context.Background(), uuid.MustParse("11111111-1111-1111-1111-111111111111"))
	if !errors.Is(err, ports.ErrPlatformUserNotFound) {
		t.Fatalf("err=%v", err)
	}
}

func TestPostgresPlatformUserAdminStore_Get_Success(t *testing.T) {
	t.Parallel()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	dn := "Ops"
	created := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	var lastLogin *time.Time
	tx := &quotaFakeTx{}
	tx.enqueueRows(quotaFakeRow{values: []any{
		userID, "ops@ani.io", "local:ops", &dn, "platform-ops", "active", lastLogin, created,
	}})
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: tx})
	out, err := store.Get(context.Background(), userID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if out.Username != "local:ops" || out.Role != "platform-ops" {
		t.Fatalf("out=%+v", out)
	}
}

func TestPostgresPlatformUserAdminStore_List_Success(t *testing.T) {
	t.Parallel()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	dn := "Ops"
	last := time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)
	created := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	tx := &quotaFakeTx{}
	tx.enqueueQuery(&quotaFakeRows{rows: []quotaFakeRow{
		{values: []any{userID, "ops@ani.io", "local:ops", &dn, "platform-ops", "active", &last, created}},
	}})
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: tx})
	res, err := store.List(context.Background(), ports.PlatformUserFilter{
		Limit: 20, Role: "platform-ops", Status: "active", Source: "local", Search: "ops",
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(res.Items) != 1 || res.Items[0].Role != "platform-ops" {
		t.Fatalf("res=%+v", res)
	}
	if len(tx.querySQLs) != 1 {
		t.Fatalf("queries=%v", tx.querySQLs)
	}
	sql := tx.querySQLs[0]
	if strings.Contains(sql, "u.email ILIKE") {
		t.Fatalf("search must not match email: %s", sql)
	}
	if !strings.Contains(sql, "REGEXP_REPLACE(u.username") {
		t.Fatalf("search must strip username prefix before ILIKE: %s", sql)
	}
}

func TestPostgresPlatformUserAdminStore_List_NextCursor(t *testing.T) {
	t.Parallel()
	makeRow := func(id string, created time.Time) quotaFakeRow {
		dn := "Ops"
		var lastLogin *time.Time
		return quotaFakeRow{values: []any{
			uuid.MustParse(id), "ops@ani.io", "local:ops", &dn, "platform-ops", "active", lastLogin, created,
		}}
	}
	t1 := time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	tx := &quotaFakeTx{}
	tx.enqueueQuery(&quotaFakeRows{rows: []quotaFakeRow{
		makeRow("11111111-1111-1111-1111-111111111111", t1),
		makeRow("22222222-2222-2222-2222-222222222222", t2),
	}})
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: tx})
	res, err := store.List(context.Background(), ports.PlatformUserFilter{Limit: 1})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(res.Items) != 1 || res.NextCursor == "" {
		t.Fatalf("res=%+v", res)
	}
}

func TestPostgresPlatformUserAdminStore_List_InvalidSource(t *testing.T) {
	t.Parallel()
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: &quotaFakeTx{}})
	_, err := store.List(context.Background(), ports.PlatformUserFilter{Source: "invalid"})
	if !errors.Is(err, ports.ErrValidationFailed) {
		t.Fatalf("err=%v", err)
	}
}

func TestPostgresPlatformUserAdminStore_List_InvalidStatus(t *testing.T) {
	t.Parallel()
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: &quotaFakeTx{}})
	_, err := store.List(context.Background(), ports.PlatformUserFilter{Status: "pending"})
	if !errors.Is(err, ports.ErrValidationFailed) {
		t.Fatalf("err=%v", err)
	}
}

func TestPostgresPlatformUserAdminStore_SetStatus_LastAdmin(t *testing.T) {
	t.Parallel()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	tx := &quotaFakeTx{}
	tx.enqueueRows(
		quotaFakeRow{values: []any{"platform-admin"}}, // current role
		quotaFakeRow{values: []any{int64(0)}},         // other active admins
	)
	store := NewPostgresPlatformUserAdminStore(&quotaFakeStore{tx: tx})
	err := store.SetStatus(context.Background(), userID, "disabled")
	if !errors.Is(err, ports.ErrLastPlatformAdmin) {
		t.Fatalf("err=%v", err)
	}
}

func TestInferPlatformSource(t *testing.T) {
	t.Parallel()
	cases := []struct {
		username string
		want     string
	}{
		{"oidc:alice", "third_party"},
		{"local:ops", "local"},
		{"ops", "unknown"},
		{"", "unknown"},
	}
	for _, tc := range cases {
		if got := inferPlatformSource(tc.username); got != tc.want {
			t.Fatalf("inferPlatformSource(%q)=%q want %q", tc.username, got, tc.want)
		}
	}
}

func TestValidatePlatformUserCreateFields(t *testing.T) {
	t.Parallel()
	err := validatePlatformUserCreateFields("bad", "a:b", "", "", "")
	if !errors.Is(err, ports.ErrValidationFailed) {
		t.Fatalf("err=%v", err)
	}
}
