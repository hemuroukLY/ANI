package runtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestLocalWorkloadReconcileControllerReconcileNowUpdatesStore(t *testing.T) {
	store := newReconcileMemoryStore()
	record := reconcileTestRecord(ports.WorkloadStateProvisioning)
	if err := store.UpsertStatus(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	controller := NewLocalWorkloadReconcileController(
		store,
		store,
		NewLocalProviderStatusReader(WithStatusReaderClock(func() time.Time { return time.Unix(210, 0) })),
		NewLocalStatusReconciler(WithReconcileClock(func() time.Time { return time.Unix(220, 0) })),
		ports.ReconcileControllerConfig{},
		WithReconcileControllerClock(func() time.Time { return time.Unix(200, 0) }),
	)

	result, err := controller.ReconcileNow(context.Background(), ports.ReconcileTarget{
		TenantID:   record.TenantID,
		InstanceID: record.InstanceID,
		Kind:       record.Kind,
		Provider:   record.Provider,
		State:      record.Status.State,
	})
	if err != nil {
		t.Fatalf("ReconcileNow() error = %v", err)
	}
	if !result.StateChanged || result.PreviousState != ports.WorkloadStateProvisioning || result.CurrentState != ports.WorkloadStateRunning {
		t.Fatalf("unexpected reconcile result: %+v", result)
	}
	updated, err := store.Get(context.Background(), record.TenantID, record.InstanceID)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status.State != ports.WorkloadStateRunning {
		t.Fatalf("stored state = %s, want running", updated.Status.State)
	}
}

func TestLocalWorkloadReconcileControllerMarksProviderMissing(t *testing.T) {
	store := newReconcileMemoryStore()
	record := reconcileTestRecord(ports.WorkloadStateRunning)
	if err := store.UpsertStatus(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	controller := NewLocalWorkloadReconcileController(
		store,
		store,
		missingProviderStatusReader{},
		NewLocalStatusReconciler(),
		ports.ReconcileControllerConfig{},
		WithReconcileControllerClock(func() time.Time { return time.Unix(300, 0) }),
	)

	result, err := controller.ReconcileNow(context.Background(), ports.ReconcileTarget{
		TenantID:   record.TenantID,
		InstanceID: record.InstanceID,
		Kind:       record.Kind,
		Provider:   record.Provider,
		State:      record.Status.State,
	})
	if err != nil {
		t.Fatalf("ReconcileNow() error = %v", err)
	}
	if !result.ProviderMissing || result.CurrentState != ports.WorkloadStateFailed {
		t.Fatalf("unexpected missing-provider result: %+v", result)
	}
	updated, err := store.Get(context.Background(), record.TenantID, record.InstanceID)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status.State != ports.WorkloadStateFailed || updated.Status.Reason != "ProviderResourceLost" {
		t.Fatalf("stored status = %+v, want failed ProviderResourceLost", updated.Status)
	}
}

func TestLocalWorkloadReconcileControllerRunOnceUsesTargetLister(t *testing.T) {
	store := newReconcileMemoryStore()
	record := reconcileTestRecord(ports.WorkloadStateProvisioning)
	if err := store.UpsertStatus(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	controller := NewLocalWorkloadReconcileController(
		store,
		store,
		NewLocalProviderStatusReader(),
		NewLocalStatusReconciler(),
		ports.ReconcileControllerConfig{MaxConcurrentReconciles: 1, StaleThresholdSeconds: 60},
		WithReconcileControllerClock(func() time.Time { return time.Unix(400, 0) }),
	)

	active, err := controller.runOnce(context.Background())
	if err != nil {
		t.Fatalf("runOnce() error = %v", err)
	}
	if !active {
		t.Fatalf("runOnce() active = false, want true for transient target")
	}
	if store.listRequests != 1 {
		t.Fatalf("ListReconcileTargets calls = %d, want 1", store.listRequests)
	}
}

type missingProviderStatusReader struct{}

func (missingProviderStatusReader) Observe(context.Context, ports.WorkloadProviderStatusRequest) (ports.WorkloadProviderObservation, error) {
	return ports.WorkloadProviderObservation{}, ports.ErrNotFound
}

type reconcileMemoryStore struct {
	records      map[string]ports.WorkloadInstanceRecord
	listRequests int
}

func newReconcileMemoryStore() *reconcileMemoryStore {
	return &reconcileMemoryStore{records: map[string]ports.WorkloadInstanceRecord{}}
}

func (s *reconcileMemoryStore) UpsertStatus(_ context.Context, record ports.WorkloadInstanceRecord) error {
	s.records[record.TenantID+"/"+record.InstanceID] = record
	return nil
}

func (s *reconcileMemoryStore) Get(_ context.Context, tenantID string, instanceID string) (ports.WorkloadInstanceRecord, error) {
	record, ok := s.records[tenantID+"/"+instanceID]
	if !ok {
		return ports.WorkloadInstanceRecord{}, ports.ErrNotFound
	}
	return record, nil
}

func (s *reconcileMemoryStore) List(_ context.Context, tenantID string, kind ports.WorkloadKind) ([]ports.WorkloadInstanceRecord, error) {
	var records []ports.WorkloadInstanceRecord
	for _, record := range s.records {
		if record.TenantID != tenantID {
			continue
		}
		if kind != "" && record.Kind != kind {
			continue
		}
		records = append(records, record)
	}
	return records, nil
}

func (s *reconcileMemoryStore) ListReconcileTargets(_ context.Context, request ports.ReconcileTargetListRequest) ([]ports.ReconcileTarget, error) {
	s.listRequests++
	if request.Limit == 0 {
		return nil, errors.New("limit must be defaulted by controller")
	}
	var targets []ports.ReconcileTarget
	for _, record := range s.records {
		targets = append(targets, ports.ReconcileTarget{
			TenantID:       record.TenantID,
			InstanceID:     record.InstanceID,
			Kind:           record.Kind,
			State:          record.Status.State,
			Provider:       record.Provider,
			LastObservedAt: record.UpdatedAt,
		})
	}
	return targets, nil
}

func reconcileTestRecord(state ports.WorkloadState) ports.WorkloadInstanceRecord {
	updatedAt := time.Unix(100, 0).UTC()
	return ports.WorkloadInstanceRecord{
		TenantID:     "tenant-a",
		InstanceID:   "inst-a",
		Name:         "vm-a",
		Kind:         ports.WorkloadKindVM,
		Provider:     "local",
		AuditID:      "11111111-1111-4111-8111-111111111111",
		ResourceRefs: []string{"VirtualMachine/tenant-a/vm-a"},
		Status: ports.WorkloadStatus{
			Ref: ports.WorkloadRef{
				TenantID:   "tenant-a",
				InstanceID: "inst-a",
				Kind:       ports.WorkloadKindVM,
				ProviderID: "vm-a",
			},
			State:     state,
			Reason:    "before reconcile",
			UpdatedAt: updatedAt,
		},
		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}
}

var _ ports.WorkloadInstanceStore = (*reconcileMemoryStore)(nil)
var _ ports.ReconcileTargetLister = (*reconcileMemoryStore)(nil)
