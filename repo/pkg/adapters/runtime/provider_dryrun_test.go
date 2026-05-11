package runtime

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestLocalProviderDryRunAcceptsRenderedDeployment(t *testing.T) {
	renderer := NewKubernetesDryRunRenderer(NewPlanningRuntime())
	manifests, err := renderer.Render(context.Background(), ports.WorkloadSpec{
		TenantID: "tenant-a",
		Name:     "app-01",
		Kind:     ports.WorkloadKindContainer,
		Image:    "harbor/app:1",
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	admission, err := NewLocalAdmissionGuard().Review(context.Background(), manifests)
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}

	executor := NewLocalProviderDryRun(WithDryRunClock(func() time.Time {
		return time.Unix(100, 0)
	}))
	result, err := executor.DryRun(context.Background(), manifests, admission)
	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}
	if !result.Accepted {
		t.Fatalf("Accepted = false, reason = %s", result.Reason)
	}
	if result.Provider != "kubernetes" {
		t.Fatalf("Provider = %s, want kubernetes", result.Provider)
	}
}

func TestLocalProviderDryRunStopsWhenAdmissionDenied(t *testing.T) {
	result, err := NewLocalProviderDryRun().DryRun(context.Background(), []ports.WorkloadManifest{
		{Name: "bad", Kind: "Deployment", Provider: "kubernetes", Content: "{}"},
	}, ports.WorkloadAdmissionResult{Allowed: false, Reason: "denied"})
	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}
	if result.Accepted {
		t.Fatalf("Accepted = true, want false")
	}
	if !strings.Contains(result.Reason, "admission denied") {
		t.Fatalf("Reason = %q, want admission denied", result.Reason)
	}
}

func TestLocalProviderDryRunRejectsWrongProviderKind(t *testing.T) {
	manifest := ports.WorkloadManifest{
		Name:     "bad",
		Kind:     "Deployment",
		Provider: "kubevirt",
		Content: `{
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "name": "bad"
  }
}`,
	}
	result, err := NewLocalProviderDryRun().DryRun(context.Background(), []ports.WorkloadManifest{manifest}, ports.WorkloadAdmissionResult{Allowed: true})
	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}
	if result.Accepted {
		t.Fatalf("Accepted = true, want false")
	}
	if !strings.Contains(result.Reason, "kubevirt") {
		t.Fatalf("Reason = %q, want kubevirt provider failure", result.Reason)
	}
}
