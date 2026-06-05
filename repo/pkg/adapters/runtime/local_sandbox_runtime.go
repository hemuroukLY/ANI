package runtime

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

type LocalSandboxRuntime struct {
	now       func() time.Time
	sequence  atomic.Uint64
	mu        sync.RWMutex
	instances map[string]ports.SandboxInstanceStatus
}

type SandboxRuntimeOption func(*LocalSandboxRuntime)

func WithSandboxRuntimeClock(now func() time.Time) SandboxRuntimeOption {
	return func(runtime *LocalSandboxRuntime) {
		if now != nil {
			runtime.now = now
		}
	}
}

func NewLocalSandboxRuntime(options ...SandboxRuntimeOption) *LocalSandboxRuntime {
	runtime := &LocalSandboxRuntime{
		now:       time.Now,
		instances: make(map[string]ports.SandboxInstanceStatus),
	}
	for _, option := range options {
		option(runtime)
	}
	return runtime
}

func (r *LocalSandboxRuntime) Create(_ context.Context, request ports.SandboxCreateRequest) (ports.SandboxInstanceStatus, error) {
	if strings.TrimSpace(request.TenantID) == "" {
		return ports.SandboxInstanceStatus{}, fmt.Errorf("%w: tenantID is required", ports.ErrInvalid)
	}
	if strings.TrimSpace(request.Name) == "" {
		return ports.SandboxInstanceStatus{}, fmt.Errorf("%w: sandbox name is required", ports.ErrInvalid)
	}
	config := normalizeSandboxConfig(request.Config)
	if err := validateSandboxConfig(config); err != nil {
		return ports.SandboxInstanceStatus{}, err
	}
	state := ports.SandboxStatePending
	if request.AutoStart {
		state = ports.SandboxStateRunning
	}
	now := firstNonZeroTime(request.CreatedAt, r.now().UTC())
	instance := ports.SandboxInstanceStatus{
		TenantID:   request.TenantID,
		InstanceID: "sandbox_" + strconv.FormatUint(r.sequence.Add(1), 10),
		Name:       request.Name,
		Kind:       ports.WorkloadKindSandbox,
		Provider:   "local_sandbox_runtime",
		State:      state,
		Config:     config,
		DevProfile: ports.DevProfileInfo{
			Mode:         "local",
			Provider:     "local-sandbox-runtime",
			RealProvider: false,
			Reason:       "local profile records Kata sandbox intent; it is not a real Kata provider execution",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances[sandboxKey(instance.TenantID, instance.InstanceID)] = instance
	return instance, nil
}

func (r *LocalSandboxRuntime) Get(_ context.Context, request ports.SandboxGetRequest) (ports.SandboxInstanceStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	instance, ok := r.instances[sandboxKey(request.TenantID, request.InstanceID)]
	if !ok {
		return ports.SandboxInstanceStatus{}, ports.ErrNotFound
	}
	return instance, nil
}

func (r *LocalSandboxRuntime) List(_ context.Context, request ports.SandboxListRequest) ([]ports.SandboxInstanceStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]ports.SandboxInstanceStatus, 0, len(r.instances))
	for _, instance := range r.instances {
		if request.TenantID != "" && instance.TenantID != request.TenantID {
			continue
		}
		items = append(items, instance)
	}
	return items, nil
}

func normalizeSandboxConfig(config ports.SandboxConfig) ports.SandboxConfig {
	if strings.TrimSpace(config.RuntimeClass) == "" {
		config.RuntimeClass = "sandbox-kata"
	}
	if config.SessionTimeout <= 0 {
		config.SessionTimeout = 30 * time.Minute
	}
	if strings.TrimSpace(string(config.NetworkEgressPolicy)) == "" {
		config.NetworkEgressPolicy = ports.SandboxNetworkEgressDenyAll
	}
	return config
}

func validateSandboxConfig(config ports.SandboxConfig) error {
	switch config.NetworkEgressPolicy {
	case ports.SandboxNetworkEgressDenyAll, ports.SandboxNetworkEgressAllowlist, ports.SandboxNetworkEgressInternet:
	default:
		return fmt.Errorf("%w: unsupported sandbox network egress policy %q", ports.ErrInvalid, config.NetworkEgressPolicy)
	}
	return nil
}

func sandboxKey(tenantID string, instanceID string) string {
	return tenantID + "/" + instanceID
}

var _ ports.SandboxRuntime = (*LocalSandboxRuntime)(nil)
