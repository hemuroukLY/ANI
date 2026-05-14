package ports

import "context"

type CommandTag struct {
	RowsAffected int64
}

type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type Row interface {
	Scan(dest ...any) error
}

type MetadataTx interface {
	Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
}

// MetadataStore owns tenant-scoped transactional metadata access.
// PostgreSQL is the default adapter, but service code should depend on this
// port when it does not require pgx-specific behavior.
type MetadataStore interface {
	Ping(ctx context.Context) error
	WithTenantTx(ctx context.Context, fn func(context.Context, MetadataTx) error) error
	WithPlatformTx(ctx context.Context, fn func(context.Context, MetadataTx) error) error
}
