package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	authv1 "github.com/kubercloud/ani/pkg/generated/pb/auth/v1"
	"github.com/kubercloud/ani/pkg/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const apiKeyEnv = "dev"

type apiKeyStore struct {
	db *pgxpool.Pool
}

type apiKeyPrincipal struct {
	TenantID uuid.UUID
	UserID   uuid.UUID
	Scopes   []string
}

func newAPIKeyStore(db *pgxpool.Pool) *apiKeyStore {
	return &apiKeyStore{db: db}
}

func (s *apiKeyStore) create(ctx context.Context, req *authv1.CreateAPIKeyRequest) (*authv1.CreateAPIKeyResponse, error) {
	tenantID, err := uuid.Parse(req.GetTenantId())
	if err != nil || tenantID == uuid.Nil {
		return nil, fmt.Errorf("invalid tenant_id")
	}
	if req.GetName() == "" {
		return nil, fmt.Errorf("name required")
	}
	var userID uuid.UUID
	if req.GetUserId() != "" {
		userID, err = uuid.Parse(req.GetUserId())
		if err != nil || userID == uuid.Nil {
			return nil, fmt.Errorf("invalid user_id")
		}
	}
	rateLimit := req.GetRateLimitRpm()
	if rateLimit <= 0 {
		rateLimit = 60
	}

	rawKey, err := generateAPIKey(tenantID)
	if err != nil {
		return nil, err
	}
	keyHash := hashAPIKey(rawKey)
	keyPrefix := prefixAPIKey(rawKey)
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID, UserID: userID})

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer rollbackTx(ctx, tx)
	if err := types.SetDBTenant(ctx, tx); err != nil {
		return nil, err
	}

	var keyID uuid.UUID
	var expiresAt any
	if req.GetExpiresAt() != nil {
		if err := req.GetExpiresAt().CheckValid(); err != nil {
			return nil, fmt.Errorf("invalid expires_at")
		}
		expiresAt = req.GetExpiresAt().AsTime()
	}
	var userIDArg any
	if userID != uuid.Nil {
		userIDArg = userID
	}
	err = tx.QueryRow(ctx, `
		INSERT INTO api_keys (
			tenant_id, user_id, name, key_hash, key_prefix, scopes, rate_limit_rpm, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, tenantID, userIDArg, req.GetName(), keyHash, keyPrefix, req.GetScopes(), rateLimit, expiresAt).Scan(&keyID)
	if err != nil {
		return nil, fmt.Errorf("insert api key: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &authv1.CreateAPIKeyResponse{
		KeyId:     keyID.String(),
		KeyValue:  rawKey,
		KeyPrefix: keyPrefix,
	}, nil
}

func (s *apiKeyStore) list(ctx context.Context, req *authv1.ListAPIKeysRequest) (*authv1.ListAPIKeysResponse, error) {
	tenantID, err := uuid.Parse(req.GetTenantId())
	if err != nil || tenantID == uuid.Nil {
		return nil, fmt.Errorf("invalid tenant_id")
	}
	var userID uuid.UUID
	if req.GetUserId() != "" {
		userID, err = uuid.Parse(req.GetUserId())
		if err != nil || userID == uuid.Nil {
			return nil, fmt.Errorf("invalid user_id")
		}
	}
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID, UserID: userID})
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer rollbackTx(ctx, tx)
	if err := types.SetDBTenant(ctx, tx); err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, key_prefix, scopes, rate_limit_rpm, created_at, expires_at, last_used_at,
			revoked_at IS NULL AND (expires_at IS NULL OR expires_at > NOW()) AS is_active
		FROM api_keys
		WHERE tenant_id=$1
	`
	args := []any{tenantID}
	if userID != uuid.Nil {
		query += " AND user_id=$2"
		args = append(args, userID)
	}
	query += " ORDER BY created_at DESC"

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	defer rows.Close()

	resp := &authv1.ListAPIKeysResponse{}
	for rows.Next() {
		var info authv1.APIKeyInfo
		var id uuid.UUID
		var createdAt time.Time
		var expiresAt pgtype.Timestamptz
		var lastUsedAt pgtype.Timestamptz
		if err := rows.Scan(&id, &info.Name, &info.KeyPrefix, &info.Scopes, &info.RateLimitRpm, &createdAt, &expiresAt, &lastUsedAt, &info.IsActive); err != nil {
			return nil, err
		}
		info.Id = id.String()
		info.CreatedAt = timestamppb.New(createdAt)
		if expiresAt.Valid {
			info.ExpiresAt = timestamppb.New(expiresAt.Time)
		}
		if lastUsedAt.Valid {
			info.LastUsedAt = timestamppb.New(lastUsedAt.Time)
		}
		resp.Keys = append(resp.Keys, &info)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *apiKeyStore) revoke(ctx context.Context, req *authv1.RevokeAPIKeyRequest) error {
	tenantID, err := uuid.Parse(req.GetTenantId())
	if err != nil || tenantID == uuid.Nil {
		return fmt.Errorf("invalid tenant_id")
	}
	keyID, err := uuid.Parse(req.GetKeyId())
	if err != nil || keyID == uuid.Nil {
		return fmt.Errorf("invalid key_id")
	}
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID})
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer rollbackTx(ctx, tx)
	if err := types.SetDBTenant(ctx, tx); err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, `
		UPDATE api_keys
		SET revoked_at=COALESCE(revoked_at, NOW())
		WHERE tenant_id=$1 AND id=$2
	`, tenantID, keyID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return types.ErrNotFound
	}
	return tx.Commit(ctx)
}

func (s *apiKeyStore) validate(ctx context.Context, rawKey string) (*apiKeyPrincipal, error) {
	tenantID, err := parseAPIKeyTenant(rawKey)
	if err != nil {
		return nil, err
	}
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID})
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer rollbackTx(ctx, tx)
	if err := types.SetDBTenant(ctx, tx); err != nil {
		return nil, err
	}

	var principal apiKeyPrincipal
	var userID pgtype.UUID
	err = tx.QueryRow(ctx, `
		SELECT tenant_id, user_id, scopes
		FROM api_keys
		WHERE tenant_id=$1
		  AND key_hash=$2
		  AND revoked_at IS NULL
		  AND (expires_at IS NULL OR expires_at > NOW())
	`, tenantID, hashAPIKey(rawKey)).Scan(&principal.TenantID, &userID, &principal.Scopes)
	if err != nil {
		return nil, err
	}
	if userID.Valid {
		principal.UserID = uuid.UUID(userID.Bytes)
	}
	if _, err := tx.Exec(ctx, `
		UPDATE api_keys SET last_used_at=NOW()
		WHERE tenant_id=$1 AND key_hash=$2
	`, tenantID, hashAPIKey(rawKey)); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &principal, nil
}

func generateAPIKey(tenantID uuid.UUID) (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	secret := base64.RawURLEncoding.EncodeToString(randomBytes)
	return "ani_" + apiKeyEnv + "_" + tenantID.String() + "_" + secret, nil
}

func parseAPIKeyTenant(rawKey string) (uuid.UUID, error) {
	parts := strings.SplitN(rawKey, "_", 4)
	if len(parts) != 4 || parts[0] != "ani" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		return uuid.Nil, fmt.Errorf("invalid api key format")
	}
	tenantID, err := uuid.Parse(parts[2])
	if err != nil || tenantID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("invalid api key tenant")
	}
	return tenantID, nil
}

func hashAPIKey(rawKey string) string {
	sum := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(sum[:])
}

func prefixAPIKey(rawKey string) string {
	if len(rawKey) <= 24 {
		return rawKey
	}
	return rawKey[:24]
}

func rollbackTx(ctx context.Context, tx pgx.Tx) {
	_ = tx.Rollback(ctx)
}
