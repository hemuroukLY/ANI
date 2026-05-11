package service

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTValidatorValidateRS256Token(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tenantID := uuid.New()
	userID := uuid.New()
	issuedAt := time.Unix(1_700_000_000, 0)
	token := signTestJWT(t, key, map[string]any{
		"iss":   "ani-test",
		"sub":   userID.String(),
		"tid":   tenantID.String(),
		"uid":   userID.String(),
		"roles": []string{"tenant-admin"},
		"exp":   issuedAt.Add(time.Hour).Unix(),
		"iat":   issuedAt.Unix(),
		"jti":   "jwt-1",
	})

	validator, err := NewJWTValidator(JWTConfig{
		PublicKeyPEM: publicKeyPEM(t, &key.PublicKey),
		Issuer:       "ani-test",
	}, nil)
	if err != nil {
		t.Fatalf("NewJWTValidator: %v", err)
	}
	validator.now = func() time.Time { return issuedAt.Add(time.Minute) }

	claims, err := validator.Validate(context.Background(), token)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if claims.TenantID != tenantID {
		t.Fatalf("tenant id = %s, want %s", claims.TenantID, tenantID)
	}
	if claims.UserID != userID {
		t.Fatalf("user id = %s, want %s", claims.UserID, userID)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "tenant-admin" {
		t.Fatalf("roles = %v", claims.Roles)
	}
}

func TestJWTValidatorRejectsExpiredToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	issuedAt := time.Unix(1_700_000_000, 0)
	token := signTestJWT(t, key, map[string]any{
		"tid": uuid.NewString(),
		"uid": uuid.NewString(),
		"exp": issuedAt.Add(time.Minute).Unix(),
	})
	validator, err := NewJWTValidator(JWTConfig{PublicKeyPEM: publicKeyPEM(t, &key.PublicKey)}, nil)
	if err != nil {
		t.Fatalf("NewJWTValidator: %v", err)
	}
	validator.now = func() time.Time { return issuedAt.Add(2 * time.Minute) }

	if _, err := validator.Validate(context.Background(), token); err == nil {
		t.Fatal("expected expired token error")
	}
}

func signTestJWT(t *testing.T, key *rsa.PrivateKey, claims map[string]any) string {
	t.Helper()
	header := encodeJSON(t, map[string]any{"alg": "RS256", "typ": "JWT"})
	payload := encodeJSON(t, claims)
	signingInput := header + "." + payload
	digest := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func encodeJSON(t *testing.T, v any) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(data)
}

func publicKeyPEM(t *testing.T, key *rsa.PublicKey) string {
	t.Helper()
	data, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: data}))
}

func privateKeyPEM(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()
	data := x509.MarshalPKCS1PrivateKey(key)
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: data}))
}
