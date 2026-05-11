package runtime

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestLocalInstanceOrchestratorCreatesAndReconciles(t *testing.T) {
	store := &fakeInstanceStore{}
	orchestrator := newTestInstanceOrchestrator(true, store)
	result, err := orchestrator.Create(context.Background(), ports.WorkloadInstanceCreateRequest{
		Spec: ports.WorkloadSpec{
			TenantID: "tenant-a",
			Name:     "app-01",
			Kind:     ports.WorkloadKindContainer,
			Image:    "harbor/app:1",
		},
		UserID:          "user-a",
		PermissionProof: "rbac:create:workload",
		RequestedAt:     time.Unix(500, 0),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.AuditID == "" {
		t.Fatalf("AuditID is empty")
	}
	if !result.Admission.Allowed {
		t.Fatalf("Admission.Allowed = false, reason = %s", result.Admission.Reason)
	}
	if !result.DryRun.Accepted {
		t.Fatalf("DryRun.Accepted = false, reason = %s", result.DryRun.Reason)
	}
	if !result.Apply.Applied {
		t.Fatalf("Apply.Applied = false, reason = %s", result.Apply.Reason)
	}
	if !result.Orchestrated {
		t.Fatalf("Orchestrated = false, want true")
	}
	if result.FinalStatus.State != ports.WorkloadStateRunning {
		t.Fatalf("FinalStatus.State = %s, want running", result.FinalStatus.State)
	}
	if store.upserts != 2 {
		t.Fatalf("store upserts = %d, want 2", store.upserts)
	}
	if store.last.Status.State != ports.WorkloadStateRunning {
		t.Fatalf("stored state = %s, want running", store.last.Status.State)
	}
}

func TestLocalInstanceOrchestratorStopsBeforeObservationWhenApplyDisabled(t *testing.T) {
	store := &fakeInstanceStore{}
	orchestrator := newTestInstanceOrchestrator(false, store)
	result, err := orchestrator.Create(context.Background(), ports.WorkloadInstanceCreateRequest{
		Spec: ports.WorkloadSpec{
			TenantID: "tenant-a",
			Name:     "app-01",
			Kind:     ports.WorkloadKindContainer,
			Image:    "harbor/app:1",
		},
		UserID:          "user-a",
		PermissionProof: "rbac:create:workload",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.Apply.Applied {
		t.Fatalf("Apply.Applied = true, want false")
	}
	if result.Orchestrated {
		t.Fatalf("Orchestrated = true, want false")
	}
	if result.FinalStatus.State != ports.WorkloadStatePending {
		t.Fatalf("FinalStatus.State = %s, want pending", result.FinalStatus.State)
	}
	if store.upserts != 1 {
		t.Fatalf("store upserts = %d, want 1", store.upserts)
	}
}

func TestLocalInstanceOrchestratorRequiresPermissionProof(t *testing.T) {
	_, err := newTestInstanceOrchestrator(true, &fakeInstanceStore{}).Create(context.Background(), ports.WorkloadInstanceCreateRequest{
		Spec: ports.WorkloadSpec{
			TenantID: "tenant-a",
			Name:     "app-01",
			Kind:     ports.WorkloadKindContainer,
			Image:    "harbor/app:1",
		},
		UserID: "user-a",
	})
	if err == nil {
		t.Fatalf("Create() error = nil, want permission proof error")
	}
	if !strings.Contains(err.Error(), "permission proof") {
		t.Fatalf("error = %q, want permission proof", err)
	}
}

func newTestInstanceOrchestrator(applyEnabled bool, store ports.WorkloadInstanceStore) *LocalInstanceOrchestrator {
	planner := NewPlanningRuntime()
	return NewLocalInstanceOrchestrator(
		planner,
		NewKubernetesDryRunRenderer(planner),
		NewLocalAdmissionGuard(),
		fakePlanAuditStore{},
		NewLocalProviderDryRun(),
		NewLocalProviderApply(WithProviderApplyEnabled(applyEnabled)),
		NewLocalProviderStatusReader(),
		NewLocalStatusReconciler(),
		WithInstanceStore(store),
	)
}

type fakePlanAuditStore struct{}

func (fakePlanAuditStore) RecordPlan(_ context.Context, record ports.WorkloadPlanAuditRecord) (string, error) {
	if record.TenantID == "" || record.InstanceName == "" || record.WorkloadKind == "" {
		return "", ports.ErrInvalid
	}
	return "audit-a", nil
}

var _ ports.WorkloadPlanAuditStore = fakePlanAuditStore{}

type fakeInstanceStore struct {
	upserts int
	last    ports.WorkloadInstanceRecord
}

func (s *fakeInstanceStore) UpsertStatus(_ context.Context, record ports.WorkloadInstanceRecord) error {
	s.upserts++
	s.last = record
	return nil
}

func (s *fakeInstanceStore) Get(context.Context, string, string) (ports.WorkloadInstanceRecord, error) {
	return s.last, nil
}

func (s *fakeInstanceStore) List(context.Context, string, ports.WorkloadKind) ([]ports.WorkloadInstanceRecord, error) {
	return []ports.WorkloadInstanceRecord{s.last}, nil
}

var _ ports.WorkloadInstanceStore = (*fakeInstanceStore)(nil)
