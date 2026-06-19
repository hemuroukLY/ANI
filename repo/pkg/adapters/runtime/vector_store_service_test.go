package runtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestLocalVectorStoreServiceDevProfile(t *testing.T) {
	service := NewLocalVectorStoreService()

	store, err := service.CreateVectorStore(context.Background(), ports.VectorStoreCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "vector-store-a",
		Name:           "kb-main",
		Dimension:      3,
		Metric:         "cosine",
	})
	if err != nil {
		t.Fatalf("CreateVectorStore() error = %v", err)
	}
	if store.StoreID == "" || store.State != ports.VectorStoreReady || store.Metric != "cosine" {
		t.Fatalf("store = %#v, want ready cosine store", store)
	}
	replay, err := service.CreateVectorStore(context.Background(), ports.VectorStoreCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "vector-store-a",
		Name:           "kb-main-retry",
		Dimension:      99,
		Metric:         "l2",
	})
	if err != nil {
		t.Fatalf("CreateVectorStore replay error = %v", err)
	}
	if replay.StoreID != store.StoreID || replay.Dimension != store.Dimension {
		t.Fatalf("replay store = %#v, want original %#v", replay, store)
	}
	if _, err := service.GetVectorStore(context.Background(), ports.VectorStoreResourceGetRequest{TenantID: "tenant-b", ResourceID: store.StoreID}); err == nil {
		t.Fatalf("GetVectorStore from another tenant succeeded, want isolation error")
	}
	results, err := service.SearchVectorStore(context.Background(), ports.VectorStoreResourceSearchRequest{
		TenantID:   "tenant-a",
		ResourceID: store.StoreID,
		Vector:     []float32{0.1, 0.2, 0.3},
		TopK:       5,
	})
	if err != nil {
		t.Fatalf("SearchVectorStore() error = %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("results = %d, want empty local dev profile result", len(results))
	}
	deleted, err := service.DeleteVectorStore(context.Background(), ports.VectorStoreResourceGetRequest{TenantID: "tenant-a", ResourceID: store.StoreID})
	if err != nil {
		t.Fatalf("DeleteVectorStore() error = %v", err)
	}
	if deleted.State != ports.VectorStoreDeleted {
		t.Fatalf("deleted state = %q, want deleted", deleted.State)
	}
}

func TestLocalVectorStoreServiceSearchValidatesDimension(t *testing.T) {
	service := NewLocalVectorStoreService()
	store, err := service.CreateVectorStore(context.Background(), ports.VectorStoreCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "vector-store-b",
		Name:           "kb-main",
		Dimension:      3,
	})
	if err != nil {
		t.Fatalf("CreateVectorStore() error = %v", err)
	}

	_, err = service.SearchVectorStore(context.Background(), ports.VectorStoreResourceSearchRequest{
		TenantID:   "tenant-a",
		ResourceID: store.StoreID,
		Vector:     []float32{0.1, 0.2},
	})
	if err == nil {
		t.Fatalf("SearchVectorStore() error = nil, want dimension mismatch")
	}
}

func TestLocalVectorStoreServiceSearchRequiresReadyStore(t *testing.T) {
	service := NewLocalVectorStoreService()
	now := time.Now().UTC()
	service.stores["vst-pending"] = ports.VectorStoreRecord{
		TenantID:  "tenant-a",
		StoreID:   "vst-pending",
		Name:      "pending-store",
		Dimension: 3,
		Metric:    "cosine",
		State:     ports.VectorStorePending,
		Reason:    "index is still building",
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := service.SearchVectorStore(context.Background(), ports.VectorStoreResourceSearchRequest{
		TenantID:   "tenant-a",
		ResourceID: "vst-pending",
		Vector:     []float32{0.1, 0.2, 0.3},
	})
	if !errors.Is(err, ports.ErrFailedPrecondition) {
		t.Fatalf("SearchVectorStore error = %v, want ErrFailedPrecondition", err)
	}
}

func TestLocalVectorStoreServiceCanCreatePendingDevProfileStore(t *testing.T) {
	service := NewLocalVectorStoreService()
	store, err := service.CreateVectorStore(context.Background(), ports.VectorStoreCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "vector-store-pending",
		Name:           "pending-index",
		Dimension:      3,
	})
	if err != nil {
		t.Fatalf("CreateVectorStore() error = %v", err)
	}
	if store.State != ports.VectorStorePending {
		t.Fatalf("store state = %s, want pending", store.State)
	}
	_, err = service.SearchVectorStore(context.Background(), ports.VectorStoreResourceSearchRequest{
		TenantID:   "tenant-a",
		ResourceID: store.StoreID,
		Vector:     []float32{0.1, 0.2, 0.3},
	})
	if !errors.Is(err, ports.ErrFailedPrecondition) {
		t.Fatalf("SearchVectorStore error = %v, want ErrFailedPrecondition", err)
	}
}
