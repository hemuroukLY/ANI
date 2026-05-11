package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubercloud/ani/pkg/ports"
)

type LocalInstanceService struct {
	orchestrator ports.WorkloadInstanceOrchestrator
	store        ports.WorkloadInstanceStore
	lifecycle    ports.WorkloadInstanceLifecycleExecutor
	ops          ports.WorkloadInstanceOps
}

type InstanceServiceOption func(*LocalInstanceService)

func WithInstanceLifecycleExecutor(lifecycle ports.WorkloadInstanceLifecycleExecutor) InstanceServiceOption {
	return func(service *LocalInstanceService) {
		service.lifecycle = lifecycle
	}
}

func NewLocalInstanceService(orchestrator ports.WorkloadInstanceOrchestrator, store ports.WorkloadInstanceStore, ops ports.WorkloadInstanceOps) *LocalInstanceService {
	return &LocalInstanceService{
		orchestrator: orchestrator,
		store:        store,
		ops:          ops,
	}
}

func NewLocalInstanceServiceWithOptions(orchestrator ports.WorkloadInstanceOrchestrator, store ports.WorkloadInstanceStore, ops ports.WorkloadInstanceOps, options ...InstanceServiceOption) *LocalInstanceService {
	service := NewLocalInstanceService(orchestrator, store, ops)
	for _, option := range options {
		option(service)
	}
	return service
}

func (s *LocalInstanceService) Create(ctx context.Context, request ports.WorkloadInstanceCreateRequest) (ports.WorkloadInstanceCreateResult, error) {
	if s.orchestrator == nil {
		return ports.WorkloadInstanceCreateResult{}, ports.ErrNotConfigured
	}
	if strings.TrimSpace(request.Spec.TenantID) == "" {
		return ports.WorkloadInstanceCreateResult{}, fmt.Errorf("%w: tenantID is required", ports.ErrInvalid)
	}
	if request.Spec.Kind != ports.WorkloadKindVM &&
		request.Spec.Kind != ports.WorkloadKindContainer &&
		request.Spec.Kind != ports.WorkloadKindGPUContainer {
		return ports.WorkloadInstanceCreateResult{}, fmt.Errorf("%w: instance service supports vm, container, and gpu_container create", ports.ErrUnsupported)
	}
	if strings.TrimSpace(request.UserID) == "" {
		return ports.WorkloadInstanceCreateResult{}, fmt.Errorf("%w: user id is required", ports.ErrInvalid)
	}
	if strings.TrimSpace(request.PermissionProof) == "" {
		return ports.WorkloadInstanceCreateResult{}, fmt.Errorf("%w: permission proof is required", ports.ErrInvalid)
	}
	return s.orchestrator.Create(ctx, request)
}

func (s *LocalInstanceService) Get(ctx context.Context, request ports.WorkloadInstanceGetRequest) (ports.WorkloadInstanceRecord, error) {
	if s.store == nil {
		return ports.WorkloadInstanceRecord{}, ports.ErrNotConfigured
	}
	if strings.TrimSpace(request.TenantID) == "" {
		return ports.WorkloadInstanceRecord{}, fmt.Errorf("%w: tenantID is required", ports.ErrInvalid)
	}
	if strings.TrimSpace(request.InstanceID) == "" {
		return ports.WorkloadInstanceRecord{}, fmt.Errorf("%w: instanceID is required", ports.ErrInvalid)
	}
	return s.store.Get(ctx, request.TenantID, request.InstanceID)
}

func (s *LocalInstanceService) List(ctx context.Context, request ports.WorkloadInstanceListRequest) ([]ports.WorkloadInstanceRecord, error) {
	if s.store == nil {
		return nil, ports.ErrNotConfigured
	}
	if strings.TrimSpace(request.TenantID) == "" {
		return nil, fmt.Errorf("%w: tenantID is required", ports.ErrInvalid)
	}
	return s.store.List(ctx, request.TenantID, request.Kind)
}

func (s *LocalInstanceService) Start(ctx context.Context, request ports.WorkloadInstanceLifecycleRequest) (ports.WorkloadInstanceRecord, error) {
	request.Action = ports.WorkloadLifecycleStart
	return s.applyLifecycle(ctx, request)
}

func (s *LocalInstanceService) Stop(ctx context.Context, request ports.WorkloadInstanceLifecycleRequest) (ports.WorkloadInstanceRecord, error) {
	request.Action = ports.WorkloadLifecycleStop
	return s.applyLifecycle(ctx, request)
}

func (s *LocalInstanceService) Restart(ctx context.Context, request ports.WorkloadInstanceLifecycleRequest) (ports.WorkloadInstanceRecord, error) {
	request.Action = ports.WorkloadLifecycleRestart
	return s.applyLifecycle(ctx, request)
}

func (s *LocalInstanceService) Resize(ctx context.Context, request ports.WorkloadInstanceResizeRequest) (ports.WorkloadInstanceRecord, error) {
	lifecycle := ports.WorkloadInstanceLifecycleRequest{
		TenantID:        request.TenantID,
		InstanceID:      request.InstanceID,
		Action:          ports.WorkloadLifecycleResize,
		UserID:          request.UserID,
		PermissionProof: request.PermissionProof,
		RequestedAt:     request.RequestedAt,
	}
	return s.applyLifecycle(ctx, lifecycle)
}

func (s *LocalInstanceService) Delete(ctx context.Context, request ports.WorkloadInstanceLifecycleRequest) (ports.WorkloadInstanceRecord, error) {
	request.Action = ports.WorkloadLifecycleDelete
	return s.applyLifecycle(ctx, request)
}

func (s *LocalInstanceService) Ops(ctx context.Context, request ports.WorkloadInstanceOpsRequest) (ports.WorkloadInstanceOpsResult, error) {
	if s.store == nil || s.ops == nil {
		return ports.WorkloadInstanceOpsResult{}, ports.ErrNotConfigured
	}
	record, err := s.Get(ctx, ports.WorkloadInstanceGetRequest{
		TenantID:   request.TenantID,
		InstanceID: request.InstanceID,
	})
	if err != nil {
		return ports.WorkloadInstanceOpsResult{}, err
	}
	return s.ops.Run(ctx, request, record)
}

func (s *LocalInstanceService) applyLifecycle(ctx context.Context, request ports.WorkloadInstanceLifecycleRequest) (ports.WorkloadInstanceRecord, error) {
	if s.store == nil {
		return ports.WorkloadInstanceRecord{}, ports.ErrNotConfigured
	}
	if strings.TrimSpace(request.TenantID) == "" || strings.TrimSpace(request.InstanceID) == "" {
		return ports.WorkloadInstanceRecord{}, fmt.Errorf("%w: tenantID and instanceID are required", ports.ErrInvalid)
	}
	if strings.TrimSpace(request.UserID) == "" || strings.TrimSpace(request.PermissionProof) == "" {
		return ports.WorkloadInstanceRecord{}, fmt.Errorf("%w: user id and permission proof are required", ports.ErrInvalid)
	}
	record, err := s.store.Get(ctx, request.TenantID, request.InstanceID)
	if err != nil {
		return ports.WorkloadInstanceRecord{}, err
	}
	next, err := transition(record.Status.State, request.Action)
	if err != nil {
		return ports.WorkloadInstanceRecord{}, err
	}
	if s.lifecycle != nil {
		result, err := s.lifecycle.Apply(ctx, request, record)
		if err != nil {
			return ports.WorkloadInstanceRecord{}, err
		}
		if !result.Accepted {
			record.Status.Reason = result.Reason
			if !result.CheckedAt.IsZero() {
				record.Status.UpdatedAt = result.CheckedAt.UTC()
				record.UpdatedAt = result.CheckedAt.UTC()
			}
			if err := s.store.UpsertStatus(ctx, record); err != nil {
				return ports.WorkloadInstanceRecord{}, err
			}
			return record, nil
		}
	}
	record.Status.State = next
	record.Status.Reason = "lifecycle " + string(request.Action) + " requested"
	if !request.RequestedAt.IsZero() {
		record.Status.UpdatedAt = request.RequestedAt.UTC()
		record.UpdatedAt = request.RequestedAt.UTC()
	}
	if err := s.store.UpsertStatus(ctx, record); err != nil {
		return ports.WorkloadInstanceRecord{}, err
	}
	return record, nil
}

var _ ports.WorkloadInstanceService = (*LocalInstanceService)(nil)
