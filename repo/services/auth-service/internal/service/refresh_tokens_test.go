package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/google/uuid"
	authv1 "github.com/kubercloud/ani/pkg/generated/pb/auth/v1"
)

func TestRefreshTokenIssuesAccessToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	issuedAt := time.Unix(1_700_000_000, 0)
	tenantID := uuid.New()
	userID := uuid.New()
	refreshStore := memoryRefreshStore{
		token: "refresh-1",
		principal: refreshPrincipal{
			TenantID: tenantID,
			UserID:   userID,
			Roles:    []string{"tenant-admin"},
		},
	}
	issuer, err := NewJWTIssuer(JWTConfig{
		PrivateKeyPEM: privateKeyPEM(t, key),
		Issuer:        "ani-test",
	})
	if err != nil {
		t.Fatalf("NewJWTIssuer: %v", err)
	}
	issuer.now = func() time.Time { return issuedAt }
	svc := &AuthService{issuer: issuer, refreshTokens: refreshStore}

	resp, err := svc.RefreshToken(context.Background(), &authv1.RefreshTokenRequest{RefreshToken: "refresh-1"})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	if resp.GetAccessToken() == "" {
		t.Fatal("expected access token")
	}
	if resp.GetExpiresIn() != int32(defaultAccessTokenTTL.Seconds()) {
		t.Fatalf("expires_in = %d", resp.GetExpiresIn())
	}

	validator, err := NewJWTValidator(JWTConfig{
		PublicKeyPEM: publicKeyPEM(t, &key.PublicKey),
		Issuer:       "ani-test",
	}, nil)
	if err != nil {
		t.Fatalf("NewJWTValidator: %v", err)
	}
	validator.now = func() time.Time { return issuedAt.Add(time.Minute) }
	claims, err := validator.Validate(context.Background(), resp.GetAccessToken())
	if err != nil {
		t.Fatalf("Validate issued access token: %v", err)
	}
	if claims.TenantID != tenantID || claims.UserID != userID {
		t.Fatalf("claims tenant/user = %s/%s", claims.TenantID, claims.UserID)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "tenant-admin" {
		t.Fatalf("roles = %v", claims.Roles)
	}
}

type memoryRefreshStore struct {
	token     string
	principal refreshPrincipal
}

func (s memoryRefreshStore) Validate(_ context.Context, rawToken string) (refreshPrincipal, error) {
	if rawToken != s.token {
		return refreshPrincipal{}, errInvalidJWT
	}
	return s.principal, nil
}
