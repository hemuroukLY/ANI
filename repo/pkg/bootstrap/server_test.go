package bootstrap

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestConfigEnvironmentOverridesWorkloadReconcileController(t *testing.T) {
	t.Setenv("WORKLOAD_RECONCILE_CONTROLLER_ENABLED", "true")
	t.Setenv("WORKLOAD_RECONCILE_NORMAL_INTERVAL_SECONDS", "45")
	t.Setenv("WORKLOAD_RECONCILE_ACTIVE_INTERVAL_SECONDS", "7")
	t.Setenv("WORKLOAD_RECONCILE_STALE_THRESHOLD_SECONDS", "180")
	t.Setenv("WORKLOAD_RECONCILE_MAX_BATCH", "12")

	cfg := (Config{}).withEnvironmentOverrides()

	if !cfg.WorkloadReconcileControllerEnabled {
		t.Fatalf("WorkloadReconcileControllerEnabled = false, want true")
	}
	if cfg.WorkloadReconcileNormalInterval != 45 {
		t.Fatalf("WorkloadReconcileNormalInterval = %d, want 45", cfg.WorkloadReconcileNormalInterval)
	}
	if cfg.WorkloadReconcileActiveInterval != 7 {
		t.Fatalf("WorkloadReconcileActiveInterval = %d, want 7", cfg.WorkloadReconcileActiveInterval)
	}
	if cfg.WorkloadReconcileStaleThreshold != 180 {
		t.Fatalf("WorkloadReconcileStaleThreshold = %d, want 180", cfg.WorkloadReconcileStaleThreshold)
	}
	if cfg.WorkloadReconcileMaxBatch != 12 {
		t.Fatalf("WorkloadReconcileMaxBatch = %d, want 12", cfg.WorkloadReconcileMaxBatch)
	}
}

func TestStartWorkloadReconcileControllerRequiresOptIn(t *testing.T) {
	controller := &fakeWorkloadReconcileController{
		started: make(chan struct{}),
		stopped: make(chan struct{}),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	deps := &Deps{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Ports:  Capabilities{WorkloadController: controller},
	}

	if started := startWorkloadReconcileController(ctx, deps); started {
		t.Fatalf("startWorkloadReconcileController() = true, want false when disabled")
	}

	deps.WorkloadReconcileControllerEnabled = true
	if started := startWorkloadReconcileController(ctx, deps); !started {
		t.Fatalf("startWorkloadReconcileController() = false, want true when enabled")
	}
	select {
	case <-controller.started:
	case <-time.After(time.Second):
		t.Fatalf("controller did not start before context cancelled")
	}
	cancel()
	select {
	case <-controller.stopped:
	case <-time.After(time.Second):
		t.Fatalf("controller did not stop after context cancellation")
	}
}

type fakeWorkloadReconcileController struct {
	started chan struct{}
	stopped chan struct{}
}

func (c *fakeWorkloadReconcileController) Start(ctx context.Context) error {
	close(c.started)
	<-ctx.Done()
	close(c.stopped)
	return nil
}

func (*fakeWorkloadReconcileController) ReconcileNow(context.Context, ports.ReconcileTarget) (ports.ReconcileResult, error) {
	return ports.ReconcileResult{}, nil
}

var _ ports.WorkloadReconcileController = (*fakeWorkloadReconcileController)(nil)
