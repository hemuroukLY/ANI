package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kubercloud/ani/pkg/ports"
)

type tokenBlocklist interface {
	Revoke(ctx context.Context, jti string, ttl time.Duration) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
}

type persistentTokenBlocklist struct {
	db    *pgxpool.Pool
	cache ports.CacheStore
	now   func() time.Time
}

func newTokenBlocklist(db *pgxpool.Pool, cache ports.CacheStore) *persistentTokenBlocklist {
	return &persistentTokenBlocklist{db: db, cache: cache, now: time.Now}
}

func (b *persistentTokenBlocklist) Revoke(ctx context.Context, jti string, ttl time.Duration) error {
	if jti == "" {
		return fmt.Errorf("jti required")
	}
	if ttl <= 0 {
		ttl = defaultAccessTokenTTL
	}
	expiresAt := b.now().Add(ttl)
	if b.db != nil {
		_, err := b.db.Exec(ctx, `
			INSERT INTO jwt_blocklist (jti, expires_at)
			VALUES ($1, $2)
			ON CONFLICT (jti) DO UPDATE
			SET expires_at = GREATEST(jwt_blocklist.expires_at, EXCLUDED.expires_at),
			    revoked_at = NOW()
		`, jti, expiresAt)
		if err != nil {
			return fmt.Errorf("persist jwt blocklist: %w", err)
		}
	}
	if b.cache != nil {
		if err := b.cache.Set(ctx, jwtBlocklistKey(jti), []byte("revoked"), ttl); err != nil {
			return fmt.Errorf("cache jwt blocklist: %w", err)
		}
	}
	return nil
}

func (b *persistentTokenBlocklist) IsRevoked(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}
	if b.cache != nil {
		blocked, err := b.cache.Exists(ctx, jwtBlocklistKey(jti))
		if err != nil {
			return false, fmt.Errorf("check jwt blocklist cache: %w", err)
		}
		if blocked {
			return true, nil
		}
	}
	if b.db == nil {
		return false, nil
	}

	var expiresAt time.Time
	err := b.db.QueryRow(ctx, `
		SELECT expires_at
		FROM jwt_blocklist
		WHERE jti=$1 AND expires_at > NOW()
	`, jti).Scan(&expiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check jwt blocklist store: %w", err)
	}
	if b.cache != nil {
		ttl := time.Until(expiresAt)
		if ttl > 0 {
			_ = b.cache.Set(ctx, jwtBlocklistKey(jti), []byte("revoked"), ttl)
		}
	}
	return true, nil
}
