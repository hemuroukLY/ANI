package runtime

import (
	"context"

	"github.com/kubercloud/ani/pkg/ports"
)

type NotConfigured struct{}

func (NotConfigured) Capabilities(context.Context) (ports.WorkloadRuntimeCapabilities, error) {
	return ports.WorkloadRuntimeCapabilities{}, ports.ErrNotConfigured
}

func (NotConfigured) Create(context.Context, ports.WorkloadSpec) (ports.WorkloadRef, error) {
	return ports.WorkloadRef{}, ports.ErrNotConfigured
}

func (NotConfigured) Get(context.Context, ports.WorkloadRef) (ports.WorkloadStatus, error) {
	return ports.WorkloadStatus{}, ports.ErrNotConfigured
}

func (NotConfigured) ApplyLifecycle(context.Context, ports.WorkloadRef, ports.WorkloadLifecycleAction) (ports.WorkloadStatus, error) {
	return ports.WorkloadStatus{}, ports.ErrNotConfigured
}

func (NotConfigured) Delete(context.Context, ports.WorkloadRef) error {
	return ports.ErrNotConfigured
}

func (NotConfigured) List(context.Context, string, ports.WorkloadKind) ([]ports.WorkloadStatus, error) {
	return nil, ports.ErrNotConfigured
}
