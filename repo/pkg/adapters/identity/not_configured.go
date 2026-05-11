package identity

import (
	"context"

	"github.com/kubercloud/ani/pkg/ports"
)

type NotConfigured struct{}

var _ ports.IdentityProvider = NotConfigured{}

func (NotConfigured) ProviderName() string {
	return "not_configured"
}

func (NotConfigured) ValidateToken(context.Context, string) (ports.IdentityClaims, error) {
	return ports.IdentityClaims{}, ports.ErrNotConfigured
}

func (NotConfigured) SyncPrincipal(context.Context, ports.IdentityClaims) error {
	return ports.ErrNotConfigured
}
