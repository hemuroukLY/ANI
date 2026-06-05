package runtime

import (
	"context"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestLocalSandboxRuntimeCreatesRunningSessionWithDevProfile(t *testing.T) {
	runtime := NewLocalSandboxRuntime(WithSandboxRuntimeClock(func() time.Time {
		return time.Unix(2100, 0).UTC()
	}))

	instance, err := runtime.Create(context.Background(), ports.SandboxCreateRequest{
		TenantID: "tenant-a",
		Name:     "agent-session",
		Config: ports.SandboxConfig{
			RuntimeClass:        "sandbox-kata",
			SessionTimeout:      45 * time.Minute,
			NetworkEgressPolicy: ports.SandboxNetworkEgressDenyAll,
		},
		AutoStart: true,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if instance.Kind != ports.WorkloadKindSandbox || instance.State != ports.SandboxStateRunning {
		t.Fatalf("kind/state = %s/%s, want sandbox/running", instance.Kind, instance.State)
	}
	if instance.Config.RuntimeClass != "sandbox-kata" || instance.Config.SessionTimeout != 45*time.Minute || instance.Config.NetworkEgressPolicy != ports.SandboxNetworkEgressDenyAll {
		t.Fatalf("config = %+v, want request config", instance.Config)
	}
	if instance.DevProfile.Mode != "local" || instance.DevProfile.Provider != "local-sandbox-runtime" || instance.DevProfile.RealProvider {
		t.Fatalf("dev profile = %+v, want local non-real marker", instance.DevProfile)
	}

	got, err := runtime.Get(context.Background(), ports.SandboxGetRequest{
		TenantID:   "tenant-a",
		InstanceID: instance.InstanceID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.InstanceID != instance.InstanceID || got.State != ports.SandboxStateRunning {
		t.Fatalf("got = %+v, want stored running instance", got)
	}
}

func TestLocalSandboxRuntimeDefaultsToKataAndPendingWhenNotAutoStarted(t *testing.T) {
	runtime := NewLocalSandboxRuntime()

	instance, err := runtime.Create(context.Background(), ports.SandboxCreateRequest{
		TenantID: "tenant-a",
		Name:     "agent-session",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if instance.Config.RuntimeClass != "sandbox-kata" {
		t.Fatalf("runtime class = %q, want sandbox-kata", instance.Config.RuntimeClass)
	}
	if instance.Config.NetworkEgressPolicy != ports.SandboxNetworkEgressDenyAll {
		t.Fatalf("egress = %s, want deny_all", instance.Config.NetworkEgressPolicy)
	}
	if instance.Config.SessionTimeout != 30*time.Minute {
		t.Fatalf("timeout = %s, want 30m", instance.Config.SessionTimeout)
	}
	if instance.State != ports.SandboxStatePending {
		t.Fatalf("state = %s, want pending", instance.State)
	}
}
