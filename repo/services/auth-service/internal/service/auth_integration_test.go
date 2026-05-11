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

func TestOIDCRefreshValidateFlow(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	issuedAt := time.Unix(1_700_000_000, 0)
	cache := newMemoryCache()
	tenantID := uuid.New()
	userID := uuid.New()
	refreshStore := newMutableRefreshStore()

	issuer, err := NewJWTIssuer(JWTConfig{
		PrivateKeyPEM: privateKeyPEM(t, key),
		Issuer:        "ani-test",
	})
	if err != nil {
		t.Fatalf("NewJWTIssuer: %v", err)
	}
	issuer.now = func() time.Time { return issuedAt }
	validator, err := NewJWTValidator(JWTConfig{
		PublicKeyPEM: publicKeyPEM(t, &key.PublicKey),
		Issuer:       "ani-test",
	}, newTokenBlocklist(nil, cache))
	if err != nil {
		t.Fatalf("NewJWTValidator: %v", err)
	}
	validator.now = func() time.Time { return issuedAt.Add(time.Minute) }

	svc := &AuthService{
		jwt:           validator,
		issuer:        issuer,
		refreshTokens: refreshStore,
		blocklist:     newTokenBlocklist(nil, cache),
		oidc: newOIDCLoginManager(cache, JWTConfig{
			OIDCAuthURL:  "https://dex.example.test/auth",
			OIDCClientID: "ani-console",
		}, fakeOIDCSessionStore{
			principal: refreshPrincipal{TenantID: tenantID, UserID: userID, Roles: []string{"tenant-admin"}},
			token:     "refresh-from-oidc",
		}, issuer),
	}
	svc.oidc.exchanger = fakeOIDCExchanger{idToken: "id-token"}
	svc.oidc.verifier = fakeOIDCVerifier{claims: oidcClaims{
		Subject: "sub-1",
		Email:   "user@example.test",
		Groups:  []string{"ani-admins"},
	}}

	begin, err := svc.BeginOIDCLogin(context.Background(), &authv1.BeginOIDCLoginRequest{
		TenantName:  "tenant-a",
		RedirectUri: "https://console.example.test/callback",
	})
	if err != nil {
		t.Fatalf("BeginOIDCLogin: %v", err)
	}
	pair, err := svc.CompleteOIDCLogin(context.Background(), &authv1.CompleteOIDCLoginRequest{
		State:       begin.GetState(),
		Code:        "code-1",
		RedirectUri: "https://console.example.test/callback",
	})
	if err != nil {
		t.Fatalf("CompleteOIDCLogin: %v", err)
	}
	if pair.GetRefreshToken() != "refresh-from-oidc" {
		t.Fatalf("refresh token = %q", pair.GetRefreshToken())
	}
	refreshStore.tokens[pair.GetRefreshToken()] = refreshPrincipal{TenantID: tenantID, UserID: userID, Roles: []string{"tenant-admin"}}

	refreshed, err := svc.RefreshToken(context.Background(), &authv1.RefreshTokenRequest{RefreshToken: pair.GetRefreshToken()})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	tc, err := svc.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{Token: refreshed.GetAccessToken()})
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if tc.GetTenantId() != tenantID.String() || tc.GetUserId() != userID.String() {
		t.Fatalf("tenant/user = %s/%s", tc.GetTenantId(), tc.GetUserId())
	}
	if len(tc.GetRoles()) != 1 || tc.GetRoles()[0] != "tenant-admin" {
		t.Fatalf("roles = %v", tc.GetRoles())
	}
}

type mutableRefreshStore struct {
	tokens map[string]refreshPrincipal
}

func newMutableRefreshStore() *mutableRefreshStore {
	return &mutableRefreshStore{tokens: map[string]refreshPrincipal{}}
}

func (s *mutableRefreshStore) Validate(_ context.Context, rawToken string) (refreshPrincipal, error) {
	principal, ok := s.tokens[rawToken]
	if !ok {
		return refreshPrincipal{}, errInvalidJWT
	}
	return principal, nil
}
