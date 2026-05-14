package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

type LocalProviderStatusReader struct {
	now func() time.Time
}

type ProviderStatusReaderOption func(*LocalProviderStatusReader)

func WithStatusReaderClock(now func() time.Time) ProviderStatusReaderOption {
	return func(reader *LocalProviderStatusReader) {
		if now != nil {
			reader.now = now
		}
	}
}

func NewLocalProviderStatusReader(options ...ProviderStatusReaderOption) *LocalProviderStatusReader {
	reader := &LocalProviderStatusReader{now: time.Now}
	for _, option := range options {
		option(reader)
	}
	return reader
}

func (r *LocalProviderStatusReader) Observe(_ context.Context, request ports.WorkloadProviderStatusRequest) (ports.WorkloadProviderObservation, error) {
	if request.TenantID == "" {
		return ports.WorkloadProviderObservation{}, fmt.Errorf("%w: tenant id is required for provider status observation", ports.ErrInvalid)
	}
	if request.InstanceID == "" {
		return ports.WorkloadProviderObservation{}, fmt.Errorf("%w: instance id is required for provider status observation", ports.ErrInvalid)
	}
	if request.Kind == "" {
		return ports.WorkloadProviderObservation{}, fmt.Errorf("%w: workload kind is required for provider status observation", ports.ErrInvalid)
	}
	if !request.ApplyResult.Applied {
		return ports.WorkloadProviderObservation{}, fmt.Errorf("%w: provider apply must be applied before provider status observation", ports.ErrInvalid)
	}
	if request.ApplyResult.Provider == "" {
		return ports.WorkloadProviderObservation{}, fmt.Errorf("%w: apply provider is required for provider status observation", ports.ErrInvalid)
	}
	if len(request.ApplyResult.ResourceRefs) == 0 {
		return ports.WorkloadProviderObservation{}, fmt.Errorf("%w: apply resource refs are required for provider status observation", ports.ErrInvalid)
	}

	return ports.WorkloadProviderObservation{
		TenantID:     request.TenantID,
		InstanceID:   request.InstanceID,
		Kind:         request.Kind,
		Provider:     request.ApplyResult.Provider,
		ResourceRefs: append([]string(nil), request.ApplyResult.ResourceRefs...),
		Phase:        phaseAfterApply(request.ApplyResult.Operation),
		Reason:       "observed by local provider status reader",
		ObservedAt:   r.now().UTC(),
	}, nil
}

func phaseAfterApply(operation ports.WorkloadLifecycleAction) string {
	switch operation {
	case ports.WorkloadLifecycleCreate, ports.WorkloadLifecycleStart, ports.WorkloadLifecycleRestart:
		return "Running"
	case ports.WorkloadLifecycleStop:
		return "Stopped"
	case ports.WorkloadLifecycleDelete:
		return "Deleted"
	default:
		return "Provisioning"
	}
}

var _ ports.WorkloadProviderStatusReader = (*LocalProviderStatusReader)(nil)
