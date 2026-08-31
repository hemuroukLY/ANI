package service

import (
	"context"
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
}

func (f *fakeCoreClient) Create(context.Context, ports.PlatformUserCreateInput) (ports.PlatformUserDTO, error) {
	return ports.PlatformUserDTO{}, ports.ErrNotImplemented
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

type fakeAuditStore struct{}

func (fakeAuditStore) Create(context.Context, ports.AuditLog) (uuid.UUID, error) {
	return uuid.Nil, ports.ErrNotImplemented
}
func (fakeAuditStore) ListUserAuditLogs(context.Context, uuid.UUID, ports.AuditLogFilter) (ports.AuditLogListResult, error) {
	return ports.AuditLogListResult{}, ports.ErrNotImplemented
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
		{ports.ErrEmailAlreadyExists, codes.AlreadyExists, "EMAIL_ALREADY_EXISTS"},
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
		func() error { _, err := svc.CreatePlatformAdmin(context.Background(), nil); return err }(),
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
