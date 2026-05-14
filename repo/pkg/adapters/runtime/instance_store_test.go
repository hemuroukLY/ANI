package runtime

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestMetadataInstanceStoreUpsertsStatus(t *testing.T) {
	tx := &fakeMetadataTx{}
	store := NewMetadataInstanceStore(fakeMetadataStore{tx: tx}, WithInstanceStoreClock(func() time.Time {
		return time.Unix(600, 0)
	}))

	err := store.UpsertStatus(context.Background(), ports.WorkloadInstanceRecord{
		TenantID:     "5dbb1d01-0000-4000-8000-000000000001",
		InstanceID:   "inst_1",
		Name:         "app-01",
		Kind:         ports.WorkloadKindContainer,
		Provider:     "kubernetes",
		AuditID:      "5dbb1d01-0000-4000-8000-000000000002",
		ResourceRefs: []string{"kubernetes/Deployment/app-01"},
		Status: ports.WorkloadStatus{
			Ref: ports.WorkloadRef{
				TenantID:   "5dbb1d01-0000-4000-8000-000000000001",
				InstanceID: "inst_1",
				Kind:       ports.WorkloadKindContainer,
				ProviderID: "planning/container/tenant-a/1",
			},
			State:    ports.WorkloadStateRunning,
			Endpoint: "/instances/inst_1",
		},
	})
	if err != nil {
		t.Fatalf("UpsertStatus() error = %v", err)
	}
	if !strings.Contains(tx.sql, "INSERT INTO workload_instances") {
		t.Fatalf("sql = %q, want workload_instances insert", tx.sql)
	}
	if got, want := tx.args[2], "app-01"; got != want {
		t.Fatalf("name arg = %v, want %s", got, want)
	}
	if got, want := tx.args[8], "running"; got != want {
		t.Fatalf("state arg = %v, want %s", got, want)
	}
}

func TestMetadataInstanceStoreRejectsMissingInstanceID(t *testing.T) {
	store := NewMetadataInstanceStore(fakeMetadataStore{tx: &fakeMetadataTx{}})
	err := store.UpsertStatus(context.Background(), ports.WorkloadInstanceRecord{
		TenantID: "5dbb1d01-0000-4000-8000-000000000001",
		Name:     "app-01",
		Kind:     ports.WorkloadKindContainer,
		Status: ports.WorkloadStatus{
			State: ports.WorkloadStatePending,
		},
	})
	if err == nil {
		t.Fatalf("UpsertStatus() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "instanceID") {
		t.Fatalf("error = %q, want instanceID", err)
	}
}
