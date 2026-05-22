package ports

import "context"

type SecretCreateRequest struct {
	TenantID       string
	IdempotencyKey string
	Name           string
	Type           string
	Data           map[string]string
}

type SecretGetRequest struct {
	TenantID string
	SecretID string
}

type SecretListRequest struct {
	TenantID string
}

type SecretBindRequest struct {
	TenantID   string
	SecretID   string
	TargetType string
	TargetID   string
	MountPath  string
	EnvPrefix  string
}

type SecretRecord struct {
	SecretID  string
	TenantID  string
	Name      string
	Type      string
	Keys      []string
	State     string
	CreatedAt int64
	UpdatedAt int64
}

type SecretBindingRecord struct {
	BindingID  string
	SecretID   string
	TenantID   string
	TargetType string
	TargetID   string
	MountPath  string
	EnvPrefix  string
	State      string
	CreatedAt  int64
}

type SecretService interface {
	CreateSecret(ctx context.Context, req SecretCreateRequest) (SecretRecord, error)
	GetSecret(ctx context.Context, req SecretGetRequest) (SecretRecord, error)
	ListSecrets(ctx context.Context, req SecretListRequest) ([]SecretRecord, error)
	DeleteSecret(ctx context.Context, req SecretGetRequest) (SecretRecord, error)
	BindSecret(ctx context.Context, req SecretBindRequest) (SecretBindingRecord, error)
}
