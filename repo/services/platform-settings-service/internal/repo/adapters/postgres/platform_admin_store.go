package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kubercloud/ani/services/platform-settings-service/internal/repo/ports"
)

// PostgresPlatformAdminStore is the PostgreSQL adapter skeleton for PlatformAdminStore.
// It only operates on audit_logs; concrete SQL lands in later issues.
type PostgresPlatformAdminStore struct{}

var _ ports.PlatformAdminStore = (*PostgresPlatformAdminStore)(nil)

// NewPostgresPlatformAdminStore returns a placeholder audit store.
func NewPostgresPlatformAdminStore() ports.PlatformAdminStore {
	return &PostgresPlatformAdminStore{}
}

func (s *PostgresPlatformAdminStore) CreateAudit(ctx context.Context, in ports.AuditCreateInput) error {
	_ = ctx
	_ = in
	return ports.ErrNotImplemented
}

func (s *PostgresPlatformAdminStore) ListAuditLogs(ctx context.Context, userID uuid.UUID, filter ports.AuditLogFilter) (ports.AuditLogListResult, error) {
	_ = ctx
	_ = userID
	_ = filter
	return ports.AuditLogListResult{}, ports.ErrNotImplemented
}
