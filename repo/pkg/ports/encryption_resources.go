package ports

import "context"

type EncryptionKeyCreateRequest struct {
	TenantID       string
	IdempotencyKey string
	Name           string
	Algorithm      string
}

type EncryptionKeyGetRequest struct {
	TenantID string
	KeyID    string
}

type EncryptionKeyListRequest struct{ TenantID string }

type EncryptionSealRequest struct {
	TenantID       string
	IdempotencyKey string
	KeyID          string
	ObjectURI      string
}

type EncryptionUnsealTokenRequest struct {
	TenantID        string
	KeyID           string
	SealedObjectURI string
}

type EncryptionKeyRotateRequest struct {
	TenantID       string
	KeyID          string
	IdempotencyKey string
}

type EncryptionKeyRevokeRequest struct {
	TenantID       string
	KeyID          string
	IdempotencyKey string
	Reason         string
}

type EncryptionKeyRecord struct {
	KeyID     string
	TenantID  string
	Name      string
	Algorithm string
	State     string
	CreatedAt int64
	UpdatedAt int64
}

type EncryptionKeyRotationRecord struct {
	TenantID    string
	PreviousKey EncryptionKeyRecord
	RotatedKey  EncryptionKeyRecord
	RotationID  string
	RotatedAt   int64
}

type EncryptionSealRecord struct {
	KeyID           string
	TenantID        string
	ObjectURI       string
	SealedObjectURI string
	UnsealToken     string
	ExpiresAt       int64
	CreatedAt       int64
}

type EncryptionUnsealTokenRecord struct {
	KeyID           string
	TenantID        string
	SealedObjectURI string
	UnsealToken     string
	ExpiresAt       int64
	CreatedAt       int64
}

type EncryptionService interface {
	CreateKey(ctx context.Context, req EncryptionKeyCreateRequest) (EncryptionKeyRecord, error)
	GetKey(ctx context.Context, req EncryptionKeyGetRequest) (EncryptionKeyRecord, error)
	ListKeys(ctx context.Context, req EncryptionKeyListRequest) ([]EncryptionKeyRecord, error)
	DeleteKey(ctx context.Context, req EncryptionKeyGetRequest) (EncryptionKeyRecord, error)
	RotateKey(ctx context.Context, req EncryptionKeyRotateRequest) (EncryptionKeyRotationRecord, error)
	RevokeKey(ctx context.Context, req EncryptionKeyRevokeRequest) (EncryptionKeyRecord, error)
	Seal(ctx context.Context, req EncryptionSealRequest) (EncryptionSealRecord, error)
	CreateUnsealToken(ctx context.Context, req EncryptionUnsealTokenRequest) (EncryptionUnsealTokenRecord, error)
}
