package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PlatformUserAdminStore manages platform accounts and role bindings
// (users / user_roles / roles，tenant_id IS NULL) under platform RLS bypass.
// Implementations live in pkg/adapters/runtime (PostgresPlatformUserAdminStore).
//
// REST surface (Core OpenAPI, called by platform-settings-service CorePlatformUserClient):
//
//	POST   /admin/platform-users
//	GET    /admin/platform-users
//	GET    /admin/platform-users/{userId}
//	PUT    /admin/platform-users/{userId}/role
//	POST   /admin/platform-users/{userId}/reset-password
//	POST   /admin/platform-users/{userId}/disable
//	POST   /admin/platform-users/{userId}/enable
//	DELETE /admin/platform-users/{userId}
//	GET    /admin/platform-users/roles
type PlatformUserAdminStore interface {
	// Create inserts a platform account (tenant_id IS NULL) and binds a platform role.
	// passwordHash is pre-computed by the caller; Store does not hash.
	// Unknown role → ErrRoleNotFound; email/username conflict → ErrEmailAlreadyExists / ErrUsernameAlreadyExists.
	Create(ctx context.Context, in PlatformUserCreate) (PlatformUserAdmin, error)

	// List returns cursor-paginated platform accounts (tenant_id IS NULL, is_deleted=FALSE).
	List(ctx context.Context, filter PlatformUserFilter) (PlatformUserListResult, error)

	// Get returns one platform account by ID (no password_hash).
	// Missing / soft-deleted → ErrPlatformUserNotFound.
	Get(ctx context.Context, userID uuid.UUID) (PlatformUserAdmin, error)

	// ChangeRole deletes old user_roles and inserts the new role inside a transaction.
	// Unknown role → ErrRoleNotFound; illegal transition → ErrRoleChangeInvalid.
	ChangeRole(ctx context.Context, userID uuid.UUID, newRole string) error

	// ResetPassword validates the new password differs from the old, hashes, and updates password_hash.
	// Same as old → ErrPasswordSameAsOld; missing user → ErrPlatformUserNotFound.
	ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error

	// SetStatus updates users.status (active/disabled) with last-admin protection.
	// Disabling the last active platform-admin → ErrLastPlatformAdmin.
	SetStatus(ctx context.Context, userID uuid.UUID, status string) error

	// SoftDelete sets is_deleted=TRUE, deleted_at=now(), status=disabled.
	// Last active platform-admin → ErrLastPlatformAdmin.
	SoftDelete(ctx context.Context, userID uuid.UUID) error

	// CountActivePlatformAdmins counts active platform-admin rows, excluding excludeUserID.
	CountActivePlatformAdmins(ctx context.Context, excludeUserID uuid.UUID) (int, error)

	// ListPlatformRoles returns platform built-in roles (tenant_id IS NULL) from the roles table.
	ListPlatformRoles(ctx context.Context) ([]PlatformRole, error)
}

// PlatformUserAdmin is the Core platform-admin view (never includes password_hash).
// Named ...Admin to avoid clashing with PlatformLoginStore's PlatformUser (login view).
type PlatformUserAdmin struct {
	ID          uuid.UUID
	Email       string
	Username    string
	DisplayName *string
	Role        string // platform-admin | platform-ops | platform-readonly
	Status      string // active | disabled
	Source      string // local | third_party (inferred from username prefix)
	LastLoginAt *time.Time
	CreatedAt   time.Time
}

// PlatformUserCreate is the input for Create (password already bcrypt-hashed).
type PlatformUserCreate struct {
	Email        string
	Username     string // without prefix; Store prepends local:
	DisplayName  string
	Role         string
	PasswordHash string
}

// PlatformUserFilter is the cursor-page filter for List.
type PlatformUserFilter struct {
	Limit  int
	Cursor string
	Role   string
	Status string
	Source string // local | oidc (username prefix filter)
	Search string // email / username ILIKE
}

// PlatformUserListResult is a cursor page of platform admins.
type PlatformUserListResult struct {
	Items      []PlatformUserAdmin
	NextCursor string // "" = no more
}

// PlatformRole is a platform built-in role; Permissions is roles.permissions JSONB as-is.
type PlatformRole struct {
	Name        string
	Label       string
	Description string
	Permissions []map[string]any // resource/actions/scope entries (tenants / resource_pool / users / metering)
}
