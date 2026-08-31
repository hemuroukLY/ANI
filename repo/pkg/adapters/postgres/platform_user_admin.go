package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kubercloud/ani/pkg/ports"
)

// PostgresPlatformUserAdminStore is the PostgreSQL adapter skeleton for PlatformUserAdminStore.
// Concrete SQL lands in issue #4; all methods return 501 Not Implemented for now.
type PostgresPlatformUserAdminStore struct{}

var _ ports.PlatformUserAdminStore = (*PostgresPlatformUserAdminStore)(nil)

// NewPostgresPlatformUserAdminStore returns a placeholder store implementation.
func NewPostgresPlatformUserAdminStore() ports.PlatformUserAdminStore {
	return &PostgresPlatformUserAdminStore{}
}

func (s *PostgresPlatformUserAdminStore) Create(ctx context.Context, in ports.PlatformUserCreate) (ports.PlatformUserAdmin, error) {
	_ = ctx
	_ = in
	return ports.PlatformUserAdmin{}, ports.ErrNotImplemented
}

func (s *PostgresPlatformUserAdminStore) List(ctx context.Context, filter ports.PlatformUserFilter) (ports.PlatformUserListResult, error) {
	_ = ctx
	_ = filter
	return ports.PlatformUserListResult{}, ports.ErrNotImplemented
}

func (s *PostgresPlatformUserAdminStore) Get(ctx context.Context, userID uuid.UUID) (ports.PlatformUserAdmin, error) {
	_ = ctx
	_ = userID
	return ports.PlatformUserAdmin{}, ports.ErrNotImplemented
}

func (s *PostgresPlatformUserAdminStore) ChangeRole(ctx context.Context, userID uuid.UUID, newRole string) error {
	_ = ctx
	_ = userID
	_ = newRole
	return ports.ErrNotImplemented
}

func (s *PostgresPlatformUserAdminStore) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	_ = ctx
	_ = userID
	_ = newPassword
	return ports.ErrNotImplemented
}

func (s *PostgresPlatformUserAdminStore) SetStatus(ctx context.Context, userID uuid.UUID, status string) error {
	_ = ctx
	_ = userID
	_ = status
	return ports.ErrNotImplemented
}

func (s *PostgresPlatformUserAdminStore) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	_ = ctx
	_ = userID
	return ports.ErrNotImplemented
}

func (s *PostgresPlatformUserAdminStore) CountActivePlatformAdmins(ctx context.Context, excludeUserID uuid.UUID) (int, error) {
	_ = ctx
	_ = excludeUserID
	return 0, ports.ErrNotImplemented
}

func (s *PostgresPlatformUserAdminStore) ListPlatformRoles(ctx context.Context) ([]ports.PlatformRole, error) {
	_ = ctx
	return nil, ports.ErrNotImplemented
}
