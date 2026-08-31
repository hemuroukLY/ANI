package runtime

import (
	"context"

	"github.com/google/uuid"
	"github.com/kubercloud/ani/pkg/ports"
)

// PostgresPlatformUserAdminStore 是 ports.PlatformUserAdminStore 的 PostgreSQL 适配器骨架。
// 方法仅声明；users / user_roles / roles 的 SQL 由后续 Core issue 落地。
// 与 PostgresUserAdmin 同属 runtime 包，经 MetadataStore 走平台侧 RLS bypass。
type PostgresPlatformUserAdminStore struct {
	store ports.MetadataStore
}

var _ ports.PlatformUserAdminStore = (*PostgresPlatformUserAdminStore)(nil)

// NewPostgresPlatformUserAdminStore 构造基于 MetadataStore（platform tx / RLS bypass）的占位实现。
func NewPostgresPlatformUserAdminStore(store ports.MetadataStore) *PostgresPlatformUserAdminStore {
	return &PostgresPlatformUserAdminStore{store: store}
}

// Create 插入平台账号（tenant_id IS NULL）并绑定平台角色。
func (s *PostgresPlatformUserAdminStore) Create(ctx context.Context, in ports.PlatformUserCreate) (ports.PlatformUserAdmin, error) {
	_ = s.store
	_ = ctx
	_ = in
	return ports.PlatformUserAdmin{}, ports.ErrNotImplemented
}

// List 游标分页返回平台账号列表（tenant_id IS NULL，排除已软删除）。
func (s *PostgresPlatformUserAdminStore) List(ctx context.Context, filter ports.PlatformUserFilter) (ports.PlatformUserListResult, error) {
	_ = s.store
	_ = ctx
	_ = filter
	return ports.PlatformUserListResult{}, ports.ErrNotImplemented
}

// Get 按 ID 返回单个平台账号（不含 password_hash）。
func (s *PostgresPlatformUserAdminStore) Get(ctx context.Context, userID uuid.UUID) (ports.PlatformUserAdmin, error) {
	_ = s.store
	_ = ctx
	_ = userID
	return ports.PlatformUserAdmin{}, ports.ErrNotImplemented
}

// ChangeRole 在事务内删除旧 user_roles 并插入新角色。
func (s *PostgresPlatformUserAdminStore) ChangeRole(ctx context.Context, userID uuid.UUID, newRole string) error {
	_ = s.store
	_ = ctx
	_ = userID
	_ = newRole
	return ports.ErrNotImplemented
}

// ResetPassword 校验新密码与旧密码不同后，哈希并更新 password_hash。
func (s *PostgresPlatformUserAdminStore) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	_ = s.store
	_ = ctx
	_ = userID
	_ = newPassword
	return ports.ErrNotImplemented
}

// SetStatus 更新 users.status（active/disabled），含最后管理员保护。
func (s *PostgresPlatformUserAdminStore) SetStatus(ctx context.Context, userID uuid.UUID, status string) error {
	_ = s.store
	_ = ctx
	_ = userID
	_ = status
	return ports.ErrNotImplemented
}

// SoftDelete 置 is_deleted=TRUE、deleted_at=now()、status=disabled。
func (s *PostgresPlatformUserAdminStore) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	_ = s.store
	_ = ctx
	_ = userID
	return ports.ErrNotImplemented
}

// CountActivePlatformAdmins 统计活跃 platform-admin 数量（排除 excludeUserID）。
func (s *PostgresPlatformUserAdminStore) CountActivePlatformAdmins(ctx context.Context, excludeUserID uuid.UUID) (int, error) {
	_ = s.store
	_ = ctx
	_ = excludeUserID
	return 0, ports.ErrNotImplemented
}

// ListPlatformRoles 返回平台内置角色（tenant_id IS NULL）。
func (s *PostgresPlatformUserAdminStore) ListPlatformRoles(ctx context.Context) ([]ports.PlatformRole, error) {
	_ = s.store
	_ = ctx
	return nil, ports.ErrNotImplemented
}
