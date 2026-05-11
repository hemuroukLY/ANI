package ports

import "context"

type VectorCollectionRef struct {
	TenantID string
	KBID     string
}

type VectorRecord struct {
	ID       string
	Vector   []float32
	Metadata map[string]string
}

type VectorSearchQuery struct {
	Collection VectorCollectionRef
	Vector     []float32
	TopK       int
	Filter     map[string]string
}

type VectorSearchResult struct {
	ID       string
	Score    float32
	Metadata map[string]string
}

type VectorCollectionHealth struct {
	Ready  bool
	Reason string
}

type VectorStore interface {
	EnsureCollection(ctx context.Context, ref VectorCollectionRef, dimension int) error
	Upsert(ctx context.Context, ref VectorCollectionRef, records []VectorRecord) error
	Search(ctx context.Context, query VectorSearchQuery) ([]VectorSearchResult, error)
	Delete(ctx context.Context, ref VectorCollectionRef, ids []string) error
	Health(ctx context.Context, ref VectorCollectionRef) (VectorCollectionHealth, error)
}
