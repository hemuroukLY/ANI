package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	commonv1 "github.com/kubercloud/ani/pkg/generated/pb/common/v1"
	platformsettingsv1 "github.com/kubercloud/ani/pkg/generated/pb/platform_settings/v1"
	"github.com/kubercloud/ani/services/platform-settings-service/internal/repo/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeCoreClient struct {
	listFilter ports.PlatformUserListFilter
	listRes    ports.PlatformUserListDTO
	listErr    error

	createIn  ports.PlatformUserCreateInput
	createID  string
	createErr error
}

func (f *fakeCoreClient) Create(_ context.Context, in ports.PlatformUserCreateInput) (string, error) {
	f.createIn = in
	if f.createErr != nil {
		return "", f.createErr
	}
	if f.createID != "" {
		return f.createID, nil
	}
	return "", ports.ErrNotImplemented
}
func (f *fakeCoreClient) List(_ context.Context, filter ports.PlatformUserListFilter) (ports.PlatformUserListDTO, error) {
	f.listFilter = filter
	return f.listRes, f.listErr
}
func (f *fakeCoreClient) Get(context.Context, uuid.UUID) (ports.PlatformUserDTO, error) {
	return ports.PlatformUserDTO{}, ports.ErrNotImplemented
}
func (f *fakeCoreClient) ChangeRole(context.Context, uuid.UUID, string) error {
	return ports.ErrNotImplemented
}
func (f *fakeCoreClient) ResetPassword(context.Context, uuid.UUID, string) error {
	return ports.ErrNotImplemented
}
func (f *fakeCoreClient) SetStatus(context.Context, uuid.UUID, string) error {
	return ports.ErrNotImplemented
}
func (f *fakeCoreClient) SoftDelete(context.Context, uuid.UUID) error { return ports.ErrNotImplemented }
func (f *fakeCoreClient) ListPlatformRoles(context.Context) ([]ports.PlatformRoleDTO, error) {
	return nil, ports.ErrNotImplemented
}

type recordingAuditStore struct {
	logs     []ports.AuditLog
	createFn func(ctx context.Context, log ports.AuditLog) (uuid.UUID, error)
}

func (f *recordingAuditStore) Create(ctx context.Context, log ports.AuditLog) (uuid.UUID, error) {
	if f.createFn != nil {
		id, err := f.createFn(ctx, log)
		if err == nil {
			f.logs = append(f.logs, log)
		}
		return id, err
	}
	f.logs = append(f.logs, log)
	return uuid.New(), nil
}
func (f *recordingAuditStore) ListUserAuditLogs(context.Context, uuid.UUID, ports.AuditLogFilter) (ports.AuditLogListResult, error) {
	return ports.AuditLogListResult{}, ports.ErrNotImplemented
}

type fakeAuditStore struct{}

func (fakeAuditStore) Create(context.Context, ports.AuditLog) (uuid.UUID, error) {
	return uuid.Nil, nil
}
func (fakeAuditStore) ListUserAuditLogs(context.Context, uuid.UUID, ports.AuditLogFilter) (ports.AuditLogListResult, error) {
	return ports.AuditLogListResult{}, ports.ErrNotImplemented
}

func TestPlatformAdminService_Create_Success(t *testing.T) {
	t.Parallel()
	core := &fakeCoreClient{
		createID: "11111111-1111-1111-1111-111111111111",
	}
	audit := &recordingAuditStore{}
	svc := NewPlatformAdminService(core, audit)

	res, err := svc.CreatePlatformAdmin(context.Background(), &platformsettingsv1.CreatePlatformAdminRequest{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops", Role: "platform-ops",
		Password: "Abcd1234!",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if res.GetId() != core.createID || res.GetMessage() != "platform admin created" {
		t.Fatalf("res=%+v", res)
	}
	if core.createIn.Password != "Abcd1234!" {
		t.Fatalf("createIn=%+v", core.createIn)
	}
	if len(audit.logs) != 1 || audit.logs[0].Action != "platform_admin.create" || audit.logs[0].Result != "success" {
		t.Fatalf("audit=%+v", audit.logs)
	}
	if audit.logs[0].Details["target_id"] != core.createID {
		t.Fatalf("details=%+v", audit.logs[0].Details)
	}
	if _, ok := audit.logs[0].Details["role"]; ok {
		t.Fatal("success audit must not include role")
	}
	if _, ok := audit.logs[0].Details["password"]; ok {
		t.Fatal("password must not appear in audit details")
	}
}

func TestPlatformAdminService_Create_UsernameConflict(t *testing.T) {
	t.Parallel()
	core := &fakeCoreClient{createErr: ports.ErrUsernameAlreadyExists}
	svc := NewPlatformAdminService(core, &recordingAuditStore{})
	_, err := svc.CreatePlatformAdmin(context.Background(), &platformsettingsv1.CreatePlatformAdminRequest{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops", Role: "platform-ops", Password: "Abcd1234!",
	})
	st := status.Convert(err)
	if st.Code() != codes.AlreadyExists || !strings.HasPrefix(st.Message(), "USERNAME_ALREADY_EXISTS") {
		t.Fatalf("err=%v", err)
	}
}

func TestPlatformAdminService_Create_Validation(t *testing.T) {
	t.Parallel()
	audit := &recordingAuditStore{}
	svc := NewPlatformAdminService(&fakeCoreClient{}, audit)
	_, err := svc.CreatePlatformAdmin(context.Background(), &platformsettingsv1.CreatePlatformAdminRequest{
		Email: "bad", Username: "ops:x", DisplayName: "", Role: "any-role", Password: "short",
	})
	st := status.Convert(err)
	if st.Code() != codes.InvalidArgument || !strings.HasPrefix(st.Message(), "VALIDATION_FAILED") {
		t.Fatalf("err=%v", err)
	}
	if len(audit.logs) != 1 || audit.logs[0].Result != "failed" || audit.logs[0].Action != "platform_admin.create" {
		t.Fatalf("validation must write failed audit: %+v", audit.logs)
	}
	if audit.logs[0].Details["role"] != "any-role" {
		t.Fatalf("validation audit must include role: %+v", audit.logs[0].Details)
	}
	if _, ok := audit.logs[0].Details["password"]; ok {
		t.Fatal("password must not appear in audit details")
	}
}

func TestPlatformAdminService_Create_RoleRequired(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{}, &recordingAuditStore{})
	_, err := svc.CreatePlatformAdmin(context.Background(), &platformsettingsv1.CreatePlatformAdminRequest{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops", Role: "  ", Password: "Abcd1234!",
	})
	st := status.Convert(err)
	if st.Code() != codes.InvalidArgument || !strings.Contains(st.Message(), "role required") {
		t.Fatalf("err=%v", err)
	}
}

func TestPlatformAdminService_Create_RoleNotFoundFromCore(t *testing.T) {
	t.Parallel()
	core := &fakeCoreClient{createErr: ports.ErrRoleNotFound}
	svc := NewPlatformAdminService(core, &recordingAuditStore{})
	_, err := svc.CreatePlatformAdmin(context.Background(), &platformsettingsv1.CreatePlatformAdminRequest{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops", Role: "future-role", Password: "Abcd1234!",
	})
	st := status.Convert(err)
	if st.Code() != codes.NotFound || !strings.HasPrefix(st.Message(), "ROLE_NOT_FOUND") {
		t.Fatalf("err=%v", err)
	}
}

func TestPlatformAdminService_Create_AuditWriteErrorDoesNotFail(t *testing.T) {
	t.Parallel()
	core := &fakeCoreClient{
		createID: "11111111-1111-1111-1111-111111111111",
	}
	audit := &recordingAuditStore{
		createFn: func(context.Context, ports.AuditLog) (uuid.UUID, error) {
			return uuid.Nil, ports.ErrNotImplemented
		},
	}
	svc := NewPlatformAdminService(core, audit)
	res, err := svc.CreatePlatformAdmin(context.Background(), &platformsettingsv1.CreatePlatformAdminRequest{
		Email: "a@ani.io", Username: "admin", DisplayName: "Admin", Role: "platform-admin", Password: "Abcd1234!",
	})
	if err != nil || res.GetId() == "" {
		t.Fatalf("create must succeed despite audit error: err=%v res=%+v", err, res)
	}
}

func TestListPlatformAdmins_Success(t *testing.T) {
	t.Parallel()
	dn := "Ops"
	last := time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)
	core := &fakeCoreClient{
		listRes: ports.PlatformUserListDTO{
			Items: []ports.PlatformUserDTO{{
				ID: "11111111-1111-1111-1111-111111111111", Username: "local:ops",
				DisplayName: &dn, Role: "platform-ops", Status: "active", Source: "local",
				LastLoginAt: &last,
			}},
			NextCursor: "n1",
		},
	}
	svc := NewPlatformAdminService(core, fakeAuditStore{})
	res, err := svc.ListPlatformAdmins(context.Background(), &platformsettingsv1.ListPlatformAdminsRequest{
		Role: "platform-ops", Status: "active", Source: "local", Search: "ops",
		Page: &commonv1.CursorPageRequest{Limit: 10, Cursor: "c0"},
	})
	if err != nil {
		t.Fatalf("ListPlatformAdmins: %v", err)
	}
	if core.listFilter.Limit != 10 || core.listFilter.Cursor != "c0" || core.listFilter.Role != "platform-ops" {
		t.Fatalf("filter=%+v", core.listFilter)
	}
	if len(res.Items) != 1 || res.NextCursor != "n1" || res.Items[0].DisplayName != "Ops" {
		t.Fatalf("res=%+v", res)
	}
}

func TestListPlatformAdmins_CoreUnavailable(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{listErr: ports.ErrCoreUnavailable}, fakeAuditStore{})
	_, err := svc.ListPlatformAdmins(context.Background(), &platformsettingsv1.ListPlatformAdminsRequest{})
	st := status.Convert(err)
	if st.Code() != codes.Unavailable || st.Message() != "CORE_UNAVAILABLE" {
		t.Fatalf("err=%v code=%v msg=%q", err, st.Code(), st.Message())
	}
}

func TestMapDomainError_Table(t *testing.T) {
	t.Parallel()
	cases := []struct {
		err  error
		code codes.Code
		msg  string
	}{
		{ports.ErrPlatformUserNotFound, codes.NotFound, "PLATFORM_USER_NOT_FOUND"},
		{ports.ErrRoleNotFound, codes.NotFound, "ROLE_NOT_FOUND"},
		{ports.ErrUsernameAlreadyExists, codes.AlreadyExists, "USERNAME_ALREADY_EXISTS"},
		{ports.ErrLastPlatformAdmin, codes.FailedPrecondition, "LAST_PLATFORM_ADMIN"},
		{ports.ErrPasswordSameAsOld, codes.FailedPrecondition, "PASSWORD_SAME_AS_OLD"},
		{ports.ErrRoleChangeInvalid, codes.FailedPrecondition, "ROLE_CHANGE_INVALID"},
		{ports.ErrValidationFailed, codes.InvalidArgument, "VALIDATION_FAILED"},
		{ports.ErrCoreUnavailable, codes.Unavailable, "CORE_UNAVAILABLE"},
		{ports.ErrNotImplemented, codes.Unimplemented, "NOT_IMPLEMENTED"},
	}
	for _, tc := range cases {
		t.Run(tc.msg, func(t *testing.T) {
			err := mapDomainError(tc.err)
			st := status.Convert(err)
			if st.Code() != tc.code || st.Message() != tc.msg {
				t.Fatalf("got %v %q", st.Code(), st.Message())
			}
		})
	}
}

func TestUnimplementedRPCs(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{}, fakeAuditStore{})
	checks := []error{
		func() error { _, err := svc.GetPlatformAdmin(context.Background(), nil); return err }(),
		func() error { _, err := svc.ListPlatformAdminRoles(context.Background(), nil); return err }(),
		func() error { _, err := svc.UpdatePlatformAdminRole(context.Background(), nil); return err }(),
		func() error { _, err := svc.ResetPlatformAdminPassword(context.Background(), nil); return err }(),
		func() error { _, err := svc.DisablePlatformAdmin(context.Background(), nil); return err }(),
		func() error { _, err := svc.EnablePlatformAdmin(context.Background(), nil); return err }(),
		func() error { _, err := svc.DeletePlatformAdmin(context.Background(), nil); return err }(),
		func() error { _, err := svc.ListPlatformAdminAuditLogs(context.Background(), nil); return err }(),
	}
	for i, err := range checks {
		if status.Code(err) != codes.Unimplemented {
			t.Fatalf("case %d: %v", i, err)
		}
	}
}

func TestValidatePasswordComplexity(t *testing.T) {
	t.Parallel()
	if validatePasswordComplexity("Abcd1234!") != "" {
		t.Fatal("expected ok")
	}
	if validatePasswordComplexity("abcdefgh") == "" {
		t.Fatal("expected fail: only lower")
	}
	if validatePasswordComplexity("Ab1!") == "" {
		t.Fatal("expected fail: too short")
	}
}
