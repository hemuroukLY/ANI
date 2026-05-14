package runtime

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestLocalProviderApplyFailsClosedWhenDisabled(t *testing.T) {
	request := validProviderApplyRequest(t)
	result, err := NewLocalProviderApply().Apply(context.Background(), request)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if result.Applied {
		t.Fatalf("Applied = true, want false")
	}
	if !strings.Contains(result.Reason, "disabled") {
		t.Fatalf("Reason = %q, want disabled", result.Reason)
	}
}

func TestLocalProviderApplyAcceptsAuditedDryRunWhenEnabled(t *testing.T) {
	request := validProviderApplyRequest(t)
	result, err := NewLocalProviderApply(
		WithProviderApplyEnabled(true),
		WithApplyClock(func() time.Time { return time.Unix(200, 0) }),
	).Apply(context.Background(), request)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if !result.Applied {
		t.Fatalf("Applied = false, reason = %s", result.Reason)
	}
	if result.Provider != "kubernetes" {
		t.Fatalf("Provider = %s, want kubernetes", result.Provider)
	}
	if len(result.ResourceRefs) != 1 || !strings.Contains(result.ResourceRefs[0], "Deployment") {
		t.Fatalf("ResourceRefs = %#v, want Deployment ref", result.ResourceRefs)
	}
}

func TestLocalProviderApplyRejectsMissingAudit(t *testing.T) {
	request := validProviderApplyRequest(t)
	request.AuditID = ""
	_, err := NewLocalProviderApply(WithProviderApplyEnabled(true)).Apply(context.Background(), request)
	if err == nil {
		t.Fatalf("Apply() error = nil, want missing audit error")
	}
	if !strings.Contains(err.Error(), "audit id") {
		t.Fatalf("error = %q, want audit id", err)
	}
}

func TestLocalProviderApplyRejectsUnacceptedDryRun(t *testing.T) {
	request := validProviderApplyRequest(t)
	request.DryRunResult.Accepted = false
	_, err := NewLocalProviderApply(WithProviderApplyEnabled(true)).Apply(context.Background(), request)
	if err == nil {
		t.Fatalf("Apply() error = nil, want dry-run error")
	}
	if !strings.Contains(err.Error(), "dry-run") {
		t.Fatalf("error = %q, want dry-run", err)
	}
}

func validProviderApplyRequest(t *testing.T) ports.WorkloadProviderApplyRequest {
	t.Helper()
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
	dryRun, err := NewLocalProviderDryRun().DryRun(context.Background(), manifests, admission)
	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}
	return ports.WorkloadProviderApplyRequest{
		TenantID:        "tenant-a",
		UserID:          "user-a",
		InstanceID:      "instance-a",
		AuditID:         "audit-a",
		PermissionProof: "rbac:create:workload",
		Operation:       ports.WorkloadLifecycleCreate,
		Manifests:       manifests,
		AdmissionResult: admission,
		DryRunResult:    dryRun,
	}
}
