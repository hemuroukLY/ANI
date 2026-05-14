package service

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	authv1 "github.com/kubercloud/ani/pkg/generated/pb/auth/v1"
)

func TestBeginOIDCLoginStoresStateAndBuildsAuthorizationURL(t *testing.T) {
	cache := newMemoryCache()
	manager := newOIDCLoginManager(cache, JWTConfig{
		OIDCAuthURL:  "https://dex.example.test/auth",
		OIDCClientID: "ani-console",
	}, nil, nil)

	resp, err := manager.Begin(context.Background(), &authv1.BeginOIDCLoginRequest{
		TenantName:  "tenant-a",
		RedirectUri: "https://console.example.test/callback",
	})
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if resp.GetState() == "" {
		t.Fatal("expected state")
	}
	if _, err := cache.Get(context.Background(), "oidc:state:"+resp.GetState()); err != nil {
		t.Fatalf("state was not stored: %v", err)
	}

	parsed, err := url.Parse(resp.GetAuthorizationUrl())
	if err != nil {
		t.Fatalf("parse authorization url: %v", err)
	}
	query := parsed.Query()
	if query.Get("client_id") != "ani-console" {
		t.Fatalf("client_id = %q", query.Get("client_id"))
	}
	if query.Get("redirect_uri") != "https://console.example.test/callback" {
		t.Fatalf("redirect_uri = %q", query.Get("redirect_uri"))
	}
	if query.Get("state") != resp.GetState() {
		t.Fatalf("state query = %q", query.Get("state"))
	}
}

func TestOIDCJWKSVerifierAcceptsMatchingKID(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	jwksData, err := json.Marshal(map[string]any{
		"keys": []map[string]any{rsaPublicJWK("kid-1", &key.PublicKey)},
	})
	if err != nil {
		t.Fatalf("marshal jwks: %v", err)
	}

	issuedAt := time.Unix(1_700_000_000, 0)
	token := signOIDCTestJWT(t, key, "kid-1", map[string]any{
		"iss":    "https://dex.example.test",
		"sub":    "sub-1",
		"aud":    "ani-console",
		"exp":    issuedAt.Add(time.Hour).Unix(),
		"email":  "user@example.test",
		"name":   "User",
		"groups": []string{"tenant-admin"},
	})
	verifier, err := newOIDCIDTokenVerifier(JWTConfig{
		OIDCIssuerURL: "https://dex.example.test",
		OIDCClientID:  "ani-console",
		OIDCJWKSURL:   "https://dex.example.test/keys",
	})
	if err != nil {
		t.Fatalf("newOIDCIDTokenVerifier: %v", err)
	}
	jwksVerifier := verifier.(*oidcJWKSVerifier)
	jwksVerifier.now = func() time.Time { return issuedAt.Add(time.Minute) }
	jwksVerifier.httpClient = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(string(jwksData))),
			Header:     make(http.Header),
		}, nil
	})}

	claims, err := verifier.Verify(context.Background(), token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.Subject != "sub-1" || claims.Email != "user@example.test" {
		t.Fatalf("claims = %#v", claims)
	}
	if len(claims.Groups) != 1 || claims.Groups[0] != "tenant-admin" {
		t.Fatalf("groups = %v", claims.Groups)
	}
}

func TestCompleteOIDCLoginRejectsMismatchedRedirectURI(t *testing.T) {
	cache := newMemoryCache()
	manager := newOIDCLoginManager(cache, JWTConfig{
		OIDCAuthURL:  "https://dex.example.test/auth",
		OIDCClientID: "ani-console",
	}, nil, nil)
	resp, err := manager.Begin(context.Background(), &authv1.BeginOIDCLoginRequest{
		TenantName:  "tenant-a",
		RedirectUri: "https://console.example.test/callback",
	})
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}

	if _, err := manager.Complete(context.Background(), &authv1.CompleteOIDCLoginRequest{
		State:       resp.GetState(),
		Code:        "code-1",
		RedirectUri: "https://evil.example.test/callback",
	}); err == nil {
		t.Fatal("expected redirect uri mismatch to fail")
	}
}

func TestCompleteOIDCLoginIssuesTokenPair(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	cache := newMemoryCache()
	issuer, err := NewJWTIssuer(JWTConfig{PrivateKeyPEM: privateKeyPEM(t, key)})
	if err != nil {
		t.Fatalf("NewJWTIssuer: %v", err)
	}
	tenantID := uuid.New()
	userID := uuid.New()
	manager := newOIDCLoginManager(cache, JWTConfig{
		OIDCAuthURL:  "https://dex.example.test/auth",
		OIDCClientID: "ani-console",
	}, fakeOIDCSessionStore{
		principal: refreshPrincipal{TenantID: tenantID, UserID: userID, Roles: []string{"user"}},
		token:     "refresh-1",
	}, issuer)
	manager.exchanger = fakeOIDCExchanger{idToken: "id-token"}
	manager.verifier = fakeOIDCVerifier{claims: oidcClaims{
		Subject: "sub-1",
		Email:   "user@example.test",
		Groups:  []string{"user"},
	}}

	begin, err := manager.Begin(context.Background(), &authv1.BeginOIDCLoginRequest{
		TenantName:  "tenant-a",
		RedirectUri: "https://console.example.test/callback",
	})
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	resp, err := manager.Complete(context.Background(), &authv1.CompleteOIDCLoginRequest{
		State:       begin.GetState(),
		Code:        "code-1",
		RedirectUri: "https://console.example.test/callback",
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if resp.GetAccessToken() == "" || resp.GetRefreshToken() != "refresh-1" {
		t.Fatalf("unexpected token pair: %#v", resp)
	}
}

type fakeOIDCExchanger struct {
	idToken string
}

func (f fakeOIDCExchanger) Exchange(context.Context, string, string) (oidcTokenResponse, error) {
	return oidcTokenResponse{IDToken: f.idToken}, nil
}

type fakeOIDCVerifier struct {
	claims oidcClaims
}

func (f fakeOIDCVerifier) Verify(context.Context, string) (oidcClaims, error) {
	return f.claims, nil
}

type fakeOIDCSessionStore struct {
	principal refreshPrincipal
	token     string
}

func (f fakeOIDCSessionStore) CreateSession(context.Context, string, oidcClaims) (refreshPrincipal, string, error) {
	return f.principal, f.token, nil
}

func rsaPublicJWK(kid string, key *rsa.PublicKey) map[string]any {
	return map[string]any{
		"kty": "RSA",
		"use": "sig",
		"alg": "RS256",
		"kid": kid,
		"n":   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
	}
}

func signOIDCTestJWT(t *testing.T, key *rsa.PrivateKey, kid string, claims map[string]any) string {
	t.Helper()
	header := encodeJSON(t, map[string]any{"alg": "RS256", "typ": "JWT", "kid": kid})
	payload := encodeJSON(t, claims)
	signingInput := header + "." + payload
	digest := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
