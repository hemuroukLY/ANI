package ports

import (
	"context"
	"time"
)

type DevProfileInfo struct {
	Mode         string
	Provider     string
	RealProvider bool
	Reason       string
}

type SandboxNetworkEgressPolicy string

const (
	SandboxNetworkEgressDenyAll   SandboxNetworkEgressPolicy = "deny_all"
	SandboxNetworkEgressAllowlist SandboxNetworkEgressPolicy = "allowlist"
	SandboxNetworkEgressInternet  SandboxNetworkEgressPolicy = "internet"
)

type SandboxState string

const (
	SandboxStatePending SandboxState = "pending"
	SandboxStateRunning SandboxState = "running"
	SandboxStateExpired SandboxState = "expired"
	SandboxStateStopped SandboxState = "stopped"
)

type SandboxConfig struct {
	RuntimeClass        string
	SessionTimeout      time.Duration
	NetworkEgressPolicy SandboxNetworkEgressPolicy
}

type SandboxCreateRequest struct {
	TenantID  string
	Name      string
	Image     string
	Config    SandboxConfig
	AutoStart bool
	CreatedAt time.Time
}

type SandboxGetRequest struct {
	TenantID   string
	InstanceID string
}

type SandboxListRequest struct {
	TenantID string
}

type SandboxInstanceStatus struct {
	TenantID   string
	InstanceID string
	Name       string
	Kind       WorkloadKind
	Provider   string
	State      SandboxState
	Config     SandboxConfig
	DevProfile DevProfileInfo
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// SandboxRuntime owns ANI sandbox session intent and state. It does not expose
// Kubernetes, Kata, RuntimeClass, Pod, or CRI provider SDK objects.
type SandboxRuntime interface {
	Create(ctx context.Context, request SandboxCreateRequest) (SandboxInstanceStatus, error)
	Get(ctx context.Context, request SandboxGetRequest) (SandboxInstanceStatus, error)
	List(ctx context.Context, request SandboxListRequest) ([]SandboxInstanceStatus, error)
}
