package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PlatformAdminStore operates on audit_logs only; it must not touch users/roles/user_roles.
type PlatformAdminStore interface {
	CreateAudit(ctx context.Context, in AuditCreateInput) error
	ListAuditLogs(ctx context.Context, userID uuid.UUID, filter AuditLogFilter) (AuditLogListResult, error)
}

type AuditCreateInput struct {
	UserID    *uuid.UUID
	RequestID string
	Action    string // platform_admin.create | change_role | reset_password | disable | enable | delete
	Resource  string // platform_user
	Result    string // success | failed
	Details   map[string]any
	IPAddress string
	UserAgent string
}

type AuditLogFilter struct {
	Limit  int
	Cursor string
	Action string
	Result string // success | failed
}

type AuditLogListItem struct {
	ID        uuid.UUID
	Action    string
	Resource  string
	Result    string
	Details   map[string]any
	CreatedAt time.Time
}

type AuditLogListResult struct {
	Items      []AuditLogListItem
	NextCursor string
}
