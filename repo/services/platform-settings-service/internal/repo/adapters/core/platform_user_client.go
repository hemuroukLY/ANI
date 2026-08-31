package core

import (
	"context"

	"github.com/google/uuid"
	"github.com/kubercloud/ani/services/platform-settings-service/internal/repo/ports"
)

// CorePlatformUserClient is the Core SDK adapter skeleton for platform user admin APIs.
type CorePlatformUserClient struct{}

var _ ports.CorePlatformUserClient = (*CorePlatformUserClient)(nil)

// NewCorePlatformUserClient returns a placeholder Core SDK client.
func NewCorePlatformUserClient() ports.CorePlatformUserClient {
	return &CorePlatformUserClient{}
}

func (c *CorePlatformUserClient) Create(ctx context.Context, in ports.PlatformUserCreateInput) (ports.PlatformUserDTO, error) {
	_ = ctx
	_ = in
	return ports.PlatformUserDTO{}, ports.ErrNotImplemented
}

func (c *CorePlatformUserClient) List(ctx context.Context, filter ports.PlatformUserListFilter) (ports.PlatformUserListDTO, error) {
	_ = ctx
	_ = filter
	return ports.PlatformUserListDTO{}, ports.ErrNotImplemented
}

func (c *CorePlatformUserClient) Get(ctx context.Context, userID uuid.UUID) (ports.PlatformUserDTO, error) {
	_ = ctx
	_ = userID
	return ports.PlatformUserDTO{}, ports.ErrNotImplemented
}

func (c *CorePlatformUserClient) ChangeRole(ctx context.Context, userID uuid.UUID, role string) error {
	_ = ctx
	_ = userID
	_ = role
	return ports.ErrNotImplemented
}

func (c *CorePlatformUserClient) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	_ = ctx
	_ = userID
	_ = newPassword
	return ports.ErrNotImplemented
}

func (c *CorePlatformUserClient) SetStatus(ctx context.Context, userID uuid.UUID, status string) error {
	_ = ctx
	_ = userID
	_ = status
	return ports.ErrNotImplemented
}

func (c *CorePlatformUserClient) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	_ = ctx
	_ = userID
	return ports.ErrNotImplemented
}

func (c *CorePlatformUserClient) ListPlatformRoles(ctx context.Context) ([]ports.PlatformRoleDTO, error) {
	_ = ctx
	return nil, ports.ErrNotImplemented
}
