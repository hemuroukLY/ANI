package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// CorePlatformUserClient wraps Core SDK calls to /api/v1/admin/platform-users/*.
type CorePlatformUserClient interface {
	Create(ctx context.Context, in PlatformUserCreateInput) (PlatformUserDTO, error)
	List(ctx context.Context, filter PlatformUserListFilter) (PlatformUserListDTO, error)
	Get(ctx context.Context, userID uuid.UUID) (PlatformUserDTO, error)
	ChangeRole(ctx context.Context, userID uuid.UUID, role string) error
	ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	SetStatus(ctx context.Context, userID uuid.UUID, status string) error
	SoftDelete(ctx context.Context, userID uuid.UUID) error
	ListPlatformRoles(ctx context.Context) ([]PlatformRoleDTO, error)
}

type PlatformUserCreateInput struct {
	Email       string
	Username    string
	DisplayName string
	Role        string
	Password    string
}

type PlatformUserListFilter struct {
	Limit  int
	Cursor string
	Role   string
	Status string
	Source string
	Search string
}

type PlatformUserDTO struct {
	ID          string
	Email       string
	Username    string
	DisplayName *string
	Role        string
	Status      string
	Source      string
	LastLoginAt *time.Time
	CreatedAt   time.Time
}

type PlatformUserListDTO struct {
	Items      []PlatformUserDTO
	NextCursor string
}

type PlatformRoleDTO struct {
	Name        string
	Label       string
	Description string
	Permissions map[string]string
}
