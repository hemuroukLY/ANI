package runtime

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestLocalProviderStatusReaderBuildsObservationFromApplyResult(t *testing.T) {
	reader := NewLocalProviderStatusReader(WithStatusReaderClock(func() time.Time {
		return time.Unix(400, 0)
	}))
	observation, err := reader.Observe(context.Background(), ports.WorkloadProviderStatusRequest{
		TenantID:   "tenant-a",
		InstanceID: "instance-a",
		Kind:       ports.WorkloadKindContainer,
		ApplyResult: ports.WorkloadProviderApplyResult{
			Applied:      true,
			Provider:     "kubernetes",
			Operation:    ports.WorkloadLifecycleCreate,
			ResourceRefs: []string{"kubernetes/Deployment/app-01"},
		},
	})
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if observation.Phase != "Running" {
		t.Fatalf("Phase = %q, want Running", observation.Phase)
	}
	if observation.Provider != "kubernetes" {
		t.Fatalf("Provider = %q, want kubernetes", observation.Provider)
	}
	if len(observation.ResourceRefs) != 1 {
		t.Fatalf("ResourceRefs = %#v, want one ref", observation.ResourceRefs)
	}
}

func TestLocalProviderStatusReaderRejectsUnappliedResult(t *testing.T) {
	_, err := NewLocalProviderStatusReader().Observe(context.Background(), ports.WorkloadProviderStatusRequest{
		TenantID:   "tenant-a",
		InstanceID: "instance-a",
		Kind:       ports.WorkloadKindContainer,
		ApplyResult: ports.WorkloadProviderApplyResult{
			Applied: false,
		},
	})
	if err == nil {
		t.Fatalf("Observe() error = nil, want unapplied error")
	}
	if !strings.Contains(err.Error(), "applied") {
		t.Fatalf("error = %q, want applied", err)
	}
}
