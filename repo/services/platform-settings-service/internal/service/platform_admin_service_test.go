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

	getID  uuid.UUID
	getRes ports.PlatformUserDTO
	getErr error

	rolesRes []ports.PlatformRoleDTO
	rolesErr error

	permsRes ports.PlatformUserPermissionsDTO
	permsErr error

	changeRoleID    uuid.UUID
	changeRoleVal   uuid.UUID
	changeRoleErr   error
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
func (f *fakeCoreClient) Get(_ context.Context, id uuid.UUID) (ports.PlatformUserDTO, error) {
	f.getID = id
	if f.getErr != nil {
		return ports.PlatformUserDTO{}, f.getErr
	}
	if f.getRes.ID != "" {
		return f.getRes, nil
	}
	return ports.PlatformUserDTO{}, ports.ErrNotImplemented
}
func (f *fakeCoreClient) ChangeRole(_ context.Context, id uuid.UUID, roleID uuid.UUID) error {
	f.changeRoleID = id
	f.changeRoleVal = roleID
	if f.changeRoleErr != nil {
		return f.changeRoleErr
	}
	return nil
}
func (f *fakeCoreClient) ResetPassword(context.Context, uuid.UUID, string) error {
	return ports.ErrNotImplemented
}
func (f *fakeCoreClient) SetStatus(context.Context, uuid.UUID, string) error {
	return ports.ErrNotImplemented
}
func (f *fakeCoreClient) SoftDelete(context.Context, uuid.UUID) error { return ports.ErrNotImplemented }
func (f *fakeCoreClient) ListPlatformRoles(context.Context) ([]ports.PlatformRoleDTO, error) {
	if f.rolesErr != nil {
		return nil, f.rolesErr
	}
	return f.rolesRes, nil
}
func (f *fakeCoreClient) GetPlatformUserPermissions(context.Context, uuid.UUID) (ports.PlatformUserPermissionsDTO, error) {
	if f.permsErr != nil {
		return ports.PlatformUserPermissionsDTO{}, f.permsErr
	}
	if f.permsRes.UserID != "" || f.permsRes.Role != "" {
		return f.permsRes, nil
	}
	return ports.PlatformUserPermissionsDTO{}, ports.ErrNotImplemented
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
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops", RoleId: "00000000-0000-0000-0000-000000000006",
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
	if _, ok := audit.logs[0].Details["role_id"]; ok {
		t.Fatal("success audit must not include role_id")
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
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops", RoleId: "00000000-0000-0000-0000-000000000006", Password: "Abcd1234!",
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
		Email: "bad", Username: "ops:x", DisplayName: "", RoleId: "not-a-uuid", Password: "short",
	})
	st := status.Convert(err)
	if st.Code() != codes.InvalidArgument || !strings.HasPrefix(st.Message(), "VALIDATION_FAILED") {
		t.Fatalf("err=%v", err)
	}
	if len(audit.logs) != 1 || audit.logs[0].Result != "failed" || audit.logs[0].Action != "platform_admin.create" {
		t.Fatalf("validation must write failed audit: %+v", audit.logs)
	}
	if audit.logs[0].Details["role_id"] != "not-a-uuid" {
		t.Fatalf("validation audit must include role_id: %+v", audit.logs[0].Details)
	}
	if _, ok := audit.logs[0].Details["password"]; ok {
		t.Fatal("password must not appear in audit details")
	}
}

func TestPlatformAdminService_Create_RoleIDRequired(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{}, &recordingAuditStore{})
	_, err := svc.CreatePlatformAdmin(context.Background(), &platformsettingsv1.CreatePlatformAdminRequest{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops", RoleId: "  ", Password: "Abcd1234!",
	})
	st := status.Convert(err)
	if st.Code() != codes.InvalidArgument || !strings.Contains(st.Message(), "role_id required") {
		t.Fatalf("err=%v", err)
	}
}

func TestPlatformAdminService_Create_RoleNotFoundFromCore(t *testing.T) {
	t.Parallel()
	core := &fakeCoreClient{createErr: ports.ErrRoleNotFound}
	svc := NewPlatformAdminService(core, &recordingAuditStore{})
	_, err := svc.CreatePlatformAdmin(context.Background(), &platformsettingsv1.CreatePlatformAdminRequest{
		Email: "ops@ani.io", Username: "ops", DisplayName: "Ops", RoleId: "00000000-0000-0000-0000-000000000006", Password: "Abcd1234!",
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
		Email: "a@ani.io", Username: "admin", DisplayName: "Admin", RoleId: "00000000-0000-0000-0000-000000000006", Password: "Abcd1234!",
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
				DisplayName: &dn, RoleID: "00000000-0000-0000-0000-000000000006", Role: "platform-ops", Status: "active", Source: "local",
				LastLoginAt: &last,
			}},
			NextCursor: "n1",
		},
	}
	svc := NewPlatformAdminService(core, fakeAuditStore{})
	res, err := svc.ListPlatformAdmins(context.Background(), &platformsettingsv1.ListPlatformAdminsRequest{
		RoleId: "00000000-0000-0000-0000-000000000006", Status: "active", Source: "local", Search: "ops",
		Page: &commonv1.CursorPageRequest{Limit: 10, Cursor: "c0"},
	})
	if err != nil {
		t.Fatalf("ListPlatformAdmins: %v", err)
	}
	if core.listFilter.Limit != 10 || core.listFilter.Cursor != "c0" || core.listFilter.RoleID != "00000000-0000-0000-0000-000000000006" {
		t.Fatalf("filter=%+v", core.listFilter)
	}
	if len(res.Items) != 1 || res.NextCursor != "n1" || res.Items[0].DisplayName != "Ops" {
		t.Fatalf("res=%+v", res)
	}
	if res.Items[0].Username != "ops" || res.Items[0].Source != "local" {
		t.Fatalf("item username/source unexpected: %+v", res.Items[0])
	}
}

func TestListPlatformAdmins_EmptyAndLimitClamp(t *testing.T) {
	t.Parallel()
	core := &fakeCoreClient{listRes: ports.PlatformUserListDTO{Items: nil, NextCursor: ""}}
	svc := NewPlatformAdminService(core, fakeAuditStore{})
	res, err := svc.ListPlatformAdmins(context.Background(), &platformsettingsv1.ListPlatformAdminsRequest{
		Page: &commonv1.CursorPageRequest{Limit: 500},
	})
	if err != nil {
		t.Fatalf("ListPlatformAdmins: %v", err)
	}
	if core.listFilter.Limit != 100 {
		t.Fatalf("limit clamp want 100 got %d", core.listFilter.Limit)
	}
	if len(res.Items) != 0 || res.NextCursor != "" {
		t.Fatalf("empty list want items=[] next_cursor=\"\" got %+v", res)
	}
}

func TestListPlatformAdmins_SourceInferOIDC(t *testing.T) {
	t.Parallel()
	core := &fakeCoreClient{
		listRes: ports.PlatformUserListDTO{
			Items: []ports.PlatformUserDTO{{
				ID: "22222222-2222-2222-2222-222222222222", Username: "oidc:alice",
				RoleID: "00000000-0000-0000-0000-000000000008", Role: "platform-readonly", Status: "active", Source: "",
			}},
		},
	}
	svc := NewPlatformAdminService(core, fakeAuditStore{})
	res, err := svc.ListPlatformAdmins(context.Background(), &platformsettingsv1.ListPlatformAdminsRequest{
		Source: "third_party",
	})
	if err != nil {
		t.Fatalf("ListPlatformAdmins: %v", err)
	}
	if core.listFilter.Source != "third_party" {
		t.Fatalf("source filter=%q", core.listFilter.Source)
	}
	if res.Items[0].Username != "alice" || res.Items[0].Source != "third_party" {
		t.Fatalf("item=%+v", res.Items[0])
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

func TestListPlatformAdmins_InvalidSource(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{}, fakeAuditStore{})
	_, err := svc.ListPlatformAdmins(context.Background(), &platformsettingsv1.ListPlatformAdminsRequest{
		Source: "oidc",
	})
	st := status.Convert(err)
	if st.Code() != codes.InvalidArgument || !strings.Contains(st.Message(), "VALIDATION_FAILED") {
		t.Fatalf("err=%v", err)
	}
}

func TestListPlatformAdmins_InvalidStatus(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{}, fakeAuditStore{})
	_, err := svc.ListPlatformAdmins(context.Background(), &platformsettingsv1.ListPlatformAdminsRequest{
		Status: "pending",
	})
	st := status.Convert(err)
	if st.Code() != codes.InvalidArgument || !strings.Contains(st.Message(), "VALIDATION_FAILED") {
		t.Fatalf("err=%v", err)
	}
}

func TestGetPlatformAdmin_Success(t *testing.T) {
	t.Parallel()
	const id = "11111111-1111-1111-1111-111111111111"
	dn := "Ops"
	last := time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)
	created := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	core := &fakeCoreClient{
		getRes: ports.PlatformUserDTO{
			ID: id, Email: "ops@ani.io", Username: "local:ops", DisplayName: &dn,
			RoleID: "00000000-0000-0000-0000-000000000006", Role: "platform-ops", Status: "active", Source: "local",
			LastLoginAt: &last, CreatedAt: created,
		},
	}
	svc := NewPlatformAdminService(core, fakeAuditStore{})
	res, err := svc.GetPlatformAdmin(context.Background(), &platformsettingsv1.GetPlatformAdminRequest{UserId: id})
	if err != nil {
		t.Fatalf("GetPlatformAdmin: %v", err)
	}
	if core.getID.String() != id {
		t.Fatalf("getID=%s", core.getID)
	}
	if res.GetId() != id || res.GetEmail() != "ops@ani.io" || res.GetUsername() != "ops" {
		t.Fatalf("detail identity=%+v", res)
	}
	if res.GetDisplayName() != "Ops" || res.GetRoleId() != "00000000-0000-0000-0000-000000000006" || res.GetRole() != "platform-ops" || res.GetStatus() != "active" || res.GetSource() != "local" {
		t.Fatalf("detail fields=%+v", res)
	}
	if res.GetLastLoginAt() == nil || !res.GetLastLoginAt().AsTime().Equal(last) {
		t.Fatalf("last_login_at=%v", res.GetLastLoginAt())
	}
	if res.GetCreatedAt() == nil || !res.GetCreatedAt().AsTime().Equal(created) {
		t.Fatalf("created_at=%v", res.GetCreatedAt())
	}
}

func TestGetPlatformAdmin_NotFound(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{getErr: ports.ErrPlatformUserNotFound}, fakeAuditStore{})
	_, err := svc.GetPlatformAdmin(context.Background(), &platformsettingsv1.GetPlatformAdminRequest{
		UserId: "11111111-1111-1111-1111-111111111111",
	})
	st := status.Convert(err)
	if st.Code() != codes.NotFound || !strings.HasPrefix(st.Message(), "PLATFORM_USER_NOT_FOUND") {
		t.Fatalf("err=%v code=%v msg=%q", err, st.Code(), st.Message())
	}
}

func TestGetPlatformAdmin_InvalidUserID(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{}, fakeAuditStore{})
	_, err := svc.GetPlatformAdmin(context.Background(), &platformsettingsv1.GetPlatformAdminRequest{UserId: "not-a-uuid"})
	st := status.Convert(err)
	if st.Code() != codes.InvalidArgument || !strings.HasPrefix(st.Message(), "VALIDATION_FAILED") {
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

func TestPlatformAdminService_ListPlatformAdminRoles(t *testing.T) {
	t.Parallel()
	core := &fakeCoreClient{
		rolesRes: []ports.PlatformRoleDTO{{
			ID: "00000000-0000-0000-0000-000000000006", Name: "platform-ops",
			Permissions: []map[string]any{
				{"resource": "tenants", "actions": []any{"*"}, "scope": "platform"},
			},
		}},
	}
	svc := NewPlatformAdminService(core, fakeAuditStore{})
	res, err := svc.ListPlatformAdminRoles(context.Background(), &platformsettingsv1.ListPlatformAdminRolesRequest{})
	if err != nil {
		t.Fatalf("ListPlatformAdminRoles: %v", err)
	}
	if len(res.Items) != 1 || res.Items[0].GetName() != "platform-ops" {
		t.Fatalf("res=%+v", res)
	}
	perms := res.Items[0].GetPermissions()
	if len(perms) != 1 || perms[0].AsMap()["resource"] != "tenants" {
		t.Fatalf("permissions=%v", perms)
	}
}

func TestPlatformAdminService_GetPlatformAdminPermissions(t *testing.T) {
	t.Parallel()
	const id = "11111111-1111-1111-1111-111111111111"
	core := &fakeCoreClient{
		permsRes: ports.PlatformUserPermissionsDTO{
			UserID: id, RoleID: "00000000-0000-0000-0000-000000000006", Role: "platform-readonly",
			Permissions: []map[string]any{{"resource": "metering", "actions": []any{"read"}, "scope": "platform"}},
		},
	}
	svc := NewPlatformAdminService(core, fakeAuditStore{})
	res, err := svc.GetPlatformAdminPermissions(context.Background(), &platformsettingsv1.GetPlatformAdminPermissionsRequest{UserId: id})
	if err != nil {
		t.Fatalf("GetPlatformAdminPermissions: %v", err)
	}
	if res.GetUserId() != id || res.GetRole() != "platform-readonly" || len(res.GetPermissions()) != 1 {
		t.Fatalf("res=%+v", res)
	}
}

func TestPlatformAdminService_GetPlatformAdminPermissions_NotFound(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{permsErr: ports.ErrPlatformUserNotFound}, fakeAuditStore{})
	_, err := svc.GetPlatformAdminPermissions(context.Background(), &platformsettingsv1.GetPlatformAdminPermissionsRequest{
		UserId: "11111111-1111-1111-1111-111111111111",
	})
	st := status.Convert(err)
	if st.Code() != codes.NotFound || !strings.HasPrefix(st.Message(), "PLATFORM_USER_NOT_FOUND") {
		t.Fatalf("err=%v", err)
	}
}

func TestPlatformAdminService_UpdatePlatformAdminRole_Success(t *testing.T) {
	t.Parallel()
	const id = "11111111-1111-1111-1111-111111111111"
	const roleID = "00000000-0000-0000-0000-000000000006"
	audit := &recordingAuditStore{}
	core := &fakeCoreClient{
		getRes: ports.PlatformUserDTO{ID: id, Role: "platform-ops"},
		rolesRes: []ports.PlatformRoleDTO{
			{ID: roleID, Name: "platform-readonly"},
		},
	}
	svc := NewPlatformAdminService(core, audit)
	res, err := svc.UpdatePlatformAdminRole(context.Background(), &platformsettingsv1.UpdatePlatformAdminRoleRequest{
		UserId: id, RoleId: roleID, IdempotencyKey: "44444444-4444-4444-4444-444444444444",
	})
	if err != nil {
		t.Fatalf("UpdatePlatformAdminRole: %v", err)
	}
	if res.GetId() != id || core.changeRoleVal != uuid.MustParse(roleID) {
		t.Fatalf("res=%+v change=%v", res, core.changeRoleVal)
	}
	if len(audit.logs) != 1 || audit.logs[0].Action != "platform_admin.change_role" || audit.logs[0].Result != "success" {
		t.Fatalf("audit=%+v", audit.logs)
	}
	if audit.logs[0].Details["old_role"] != "platform-ops" || audit.logs[0].Details["new_role"] != "platform-readonly" {
		t.Fatalf("details=%v", audit.logs[0].Details)
	}
}

func TestPlatformAdminService_UpdatePlatformAdminRole_LastAdmin(t *testing.T) {
	t.Parallel()
	const id = "11111111-1111-1111-1111-111111111111"
	const roleID = "00000000-0000-0000-0000-000000000006"
	audit := &recordingAuditStore{}
	svc := NewPlatformAdminService(&fakeCoreClient{
		getRes:        ports.PlatformUserDTO{ID: id, Role: "platform-admin"},
		rolesRes:      []ports.PlatformRoleDTO{{ID: roleID, Name: "platform-ops"}},
		changeRoleErr: ports.ErrLastPlatformAdmin,
	}, audit)
	_, err := svc.UpdatePlatformAdminRole(context.Background(), &platformsettingsv1.UpdatePlatformAdminRoleRequest{
		UserId: id, RoleId: roleID,
	})
	st := status.Convert(err)
	if st.Code() != codes.FailedPrecondition || !strings.HasPrefix(st.Message(), "LAST_PLATFORM_ADMIN") {
		t.Fatalf("err=%v", err)
	}
	if len(audit.logs) != 1 || audit.logs[0].Result != "failed" {
		t.Fatalf("audit=%+v", audit.logs)
	}
	if audit.logs[0].Details["old_role"] != "platform-admin" || audit.logs[0].Details["new_role"] != "platform-ops" {
		t.Fatalf("details=%v", audit.logs[0].Details)
	}
}

func TestPlatformAdminService_UpdatePlatformAdminRole_RoleNotFound(t *testing.T) {
	t.Parallel()
	const id = "11111111-1111-1111-1111-111111111111"
	svc := NewPlatformAdminService(&fakeCoreClient{
		getRes:        ports.PlatformUserDTO{ID: id, Role: "platform-ops"},
		changeRoleErr: ports.ErrRoleNotFound,
	}, &recordingAuditStore{})
	_, err := svc.UpdatePlatformAdminRole(context.Background(), &platformsettingsv1.UpdatePlatformAdminRoleRequest{
		UserId: id, RoleId: "00000000-0000-0000-0000-000000000006",
	})
	st := status.Convert(err)
	if st.Code() != codes.NotFound || !strings.HasPrefix(st.Message(), "ROLE_NOT_FOUND") {
		t.Fatalf("err=%v", err)
	}
}

func TestPlatformAdminService_UpdatePlatformAdminRole_UserNotFound(t *testing.T) {
	t.Parallel()
	svc := NewPlatformAdminService(&fakeCoreClient{getErr: ports.ErrPlatformUserNotFound}, &recordingAuditStore{})
	_, err := svc.UpdatePlatformAdminRole(context.Background(), &platformsettingsv1.UpdatePlatformAdminRoleRequest{
		UserId: "11111111-1111-1111-1111-111111111111", RoleId: "00000000-0000-0000-0000-000000000006",
	})
	st := status.Convert(err)
	if st.Code() != codes.NotFound || !strings.HasPrefix(st.Message(), "PLATFORM_USER_NOT_FOUND") {
		t.Fatalf("err=%v", err)
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
