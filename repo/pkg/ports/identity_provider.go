package ports

import "context"

type IdentityClaims struct {
	Subject  string
	TenantID string
	Email    string
	Name     string
	Groups   []string
}

type IdentityProvider interface {
	ProviderName() string
	ValidateToken(ctx context.Context, token string) (IdentityClaims, error)
	SyncPrincipal(ctx context.Context, claims IdentityClaims) error
}
