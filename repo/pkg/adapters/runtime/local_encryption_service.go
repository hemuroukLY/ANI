package runtime

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kubercloud/ani/pkg/ports"
)

type localEncryptionService struct {
	mu           sync.Mutex
	byID         map[string]ports.EncryptionKeyRecord
	idem         map[string]string
	rotationIdem map[string]ports.EncryptionKeyRotationRecord
	revokeIdem   map[string]ports.EncryptionKeyRecord
	sealIdem     map[string]ports.EncryptionSealRecord
}

func NewLocalEncryptionService() ports.EncryptionService {
	return &localEncryptionService{
		byID:         map[string]ports.EncryptionKeyRecord{},
		idem:         map[string]string{},
		rotationIdem: map[string]ports.EncryptionKeyRotationRecord{},
		revokeIdem:   map[string]ports.EncryptionKeyRecord{},
		sealIdem:     map[string]ports.EncryptionSealRecord{},
	}
}

func (s *localEncryptionService) CreateKey(_ context.Context, req ports.EncryptionKeyCreateRequest) (ports.EncryptionKeyRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.TenantID == "" || req.Name == "" || req.IdempotencyKey == "" {
		return ports.EncryptionKeyRecord{}, fmt.Errorf("%w: tenant_id/name/idempotency_key required", ports.ErrInvalid)
	}
	key := req.TenantID + ":" + req.IdempotencyKey
	if id, ok := s.idem[key]; ok {
		return s.byID[id], nil
	}
	now := time.Now().Unix()
	algo := req.Algorithm
	if algo == "" {
		algo = "SM4"
	}
	rec := ports.EncryptionKeyRecord{KeyID: "ekey-" + uuid.NewString(), TenantID: req.TenantID, Name: req.Name, Algorithm: algo, State: "active", CreatedAt: now, UpdatedAt: now}
	s.byID[rec.KeyID] = rec
	s.idem[key] = rec.KeyID
	return rec, nil
}
func (s *localEncryptionService) GetKey(_ context.Context, req ports.EncryptionKeyGetRequest) (ports.EncryptionKeyRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.byID[req.KeyID]
	if !ok || rec.TenantID != req.TenantID {
		return ports.EncryptionKeyRecord{}, ports.ErrNotFound
	}
	return rec, nil
}
func (s *localEncryptionService) ListKeys(_ context.Context, req ports.EncryptionKeyListRequest) ([]ports.EncryptionKeyRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []ports.EncryptionKeyRecord{}
	for _, r := range s.byID {
		if r.TenantID == req.TenantID {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt < out[j].CreatedAt })
	return out, nil
}
func (s *localEncryptionService) DeleteKey(_ context.Context, req ports.EncryptionKeyGetRequest) (ports.EncryptionKeyRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.byID[req.KeyID]
	if !ok || rec.TenantID != req.TenantID {
		return ports.EncryptionKeyRecord{}, ports.ErrNotFound
	}
	rec.State = "deleted"
	rec.UpdatedAt = time.Now().Unix()
	s.byID[req.KeyID] = rec
	return rec, nil
}

func (s *localEncryptionService) RotateKey(_ context.Context, req ports.EncryptionKeyRotateRequest) (ports.EncryptionKeyRotationRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.TenantID == "" || req.KeyID == "" || req.IdempotencyKey == "" {
		return ports.EncryptionKeyRotationRecord{}, fmt.Errorf("%w: tenant_id/key_id/idempotency_key required", ports.ErrInvalid)
	}
	idemKey := req.TenantID + ":" + req.IdempotencyKey
	if rec, ok := s.rotationIdem[idemKey]; ok {
		return rec, nil
	}
	current, err := s.requireActiveKey(req.TenantID, req.KeyID)
	if err != nil {
		return ports.EncryptionKeyRotationRecord{}, err
	}
	now := time.Now().Unix()
	previous := current
	previous.State = "rotated"
	previous.UpdatedAt = now
	rotated := ports.EncryptionKeyRecord{
		KeyID:     "ekey-" + uuid.NewString(),
		TenantID:  current.TenantID,
		Name:      current.Name + "-rotated",
		Algorithm: current.Algorithm,
		State:     "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.byID[previous.KeyID] = previous
	s.byID[rotated.KeyID] = rotated
	rec := ports.EncryptionKeyRotationRecord{
		TenantID:    current.TenantID,
		PreviousKey: previous,
		RotatedKey:  rotated,
		RotationID:  "erot-" + uuid.NewString(),
		RotatedAt:   now,
	}
	s.rotationIdem[idemKey] = rec
	return rec, nil
}

func (s *localEncryptionService) RevokeKey(_ context.Context, req ports.EncryptionKeyRevokeRequest) (ports.EncryptionKeyRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.TenantID == "" || req.KeyID == "" || req.IdempotencyKey == "" {
		return ports.EncryptionKeyRecord{}, fmt.Errorf("%w: tenant_id/key_id/idempotency_key required", ports.ErrInvalid)
	}
	idemKey := req.TenantID + ":" + req.IdempotencyKey
	if rec, ok := s.revokeIdem[idemKey]; ok {
		return rec, nil
	}
	rec, ok := s.byID[req.KeyID]
	if !ok || rec.TenantID != req.TenantID {
		return ports.EncryptionKeyRecord{}, ports.ErrNotFound
	}
	if rec.State == "deleted" {
		return ports.EncryptionKeyRecord{}, fmt.Errorf("%w: deleted encryption key cannot be revoked", ports.ErrConflict)
	}
	rec.State = "revoked"
	rec.UpdatedAt = time.Now().Unix()
	s.byID[req.KeyID] = rec
	s.revokeIdem[idemKey] = rec
	return rec, nil
}

func (s *localEncryptionService) Seal(_ context.Context, req ports.EncryptionSealRequest) (ports.EncryptionSealRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.TenantID == "" || req.IdempotencyKey == "" || req.KeyID == "" || req.ObjectURI == "" {
		return ports.EncryptionSealRecord{}, fmt.Errorf("%w: tenant_id/idempotency_key/key_id/object_uri required", ports.ErrInvalid)
	}
	idemKey := req.TenantID + ":" + req.IdempotencyKey
	if rec, ok := s.sealIdem[idemKey]; ok {
		return rec, nil
	}
	key, err := s.requireActiveKey(req.TenantID, req.KeyID)
	if err != nil {
		return ports.EncryptionSealRecord{}, err
	}
	now := time.Now().Unix()
	digest := sha256.Sum256([]byte(req.TenantID + ":" + req.KeyID + ":" + req.ObjectURI + ":" + req.IdempotencyKey))
	sealedURI := fmt.Sprintf("sealed://local/%s/%s", key.KeyID, hex.EncodeToString(digest[:12]))
	rec := ports.EncryptionSealRecord{
		KeyID:           key.KeyID,
		TenantID:        key.TenantID,
		ObjectURI:       req.ObjectURI,
		SealedObjectURI: sealedURI,
		UnsealToken:     "utok-" + uuid.NewString(),
		ExpiresAt:       now + 3600,
		CreatedAt:       now,
	}
	s.sealIdem[idemKey] = rec
	return rec, nil
}

func (s *localEncryptionService) CreateUnsealToken(_ context.Context, req ports.EncryptionUnsealTokenRequest) (ports.EncryptionUnsealTokenRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.TenantID == "" || req.KeyID == "" || req.SealedObjectURI == "" {
		return ports.EncryptionUnsealTokenRecord{}, fmt.Errorf("%w: tenant_id/key_id/sealed_object_uri required", ports.ErrInvalid)
	}
	key, err := s.requireActiveKey(req.TenantID, req.KeyID)
	if err != nil {
		return ports.EncryptionUnsealTokenRecord{}, err
	}
	now := time.Now().Unix()
	return ports.EncryptionUnsealTokenRecord{
		KeyID:           key.KeyID,
		TenantID:        key.TenantID,
		SealedObjectURI: req.SealedObjectURI,
		UnsealToken:     "utok-" + uuid.NewString(),
		ExpiresAt:       now + 3600,
		CreatedAt:       now,
	}, nil
}

func (s *localEncryptionService) requireActiveKey(tenantID string, keyID string) (ports.EncryptionKeyRecord, error) {
	key, ok := s.byID[keyID]
	if !ok || key.TenantID != tenantID {
		return ports.EncryptionKeyRecord{}, ports.ErrNotFound
	}
	if key.State != "active" {
		return ports.EncryptionKeyRecord{}, fmt.Errorf("%w: encryption key is not active", ports.ErrConflict)
	}
	return key, nil
}
