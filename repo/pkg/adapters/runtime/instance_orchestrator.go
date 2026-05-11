package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

type LocalInstanceOrchestrator struct {
	runtime    ports.WorkloadRuntime
	renderer   ports.WorkloadRenderer
	admission  ports.WorkloadAdmission
	audit      ports.WorkloadPlanAuditStore
	dryRun     ports.WorkloadProviderDryRun
	apply      ports.WorkloadProviderApply
	reader     ports.WorkloadProviderStatusReader
	reconciler ports.WorkloadStatusReconciler
	store      ports.WorkloadInstanceStore
	now        func() time.Time
}

type InstanceOrchestratorOption func(*LocalInstanceOrchestrator)

func WithInstanceOrchestratorClock(now func() time.Time) InstanceOrchestratorOption {
	return func(orchestrator *LocalInstanceOrchestrator) {
		if now != nil {
			orchestrator.now = now
		}
	}
}

func WithInstanceStore(store ports.WorkloadInstanceStore) InstanceOrchestratorOption {
	return func(orchestrator *LocalInstanceOrchestrator) {
		orchestrator.store = store
	}
}

func NewLocalInstanceOrchestrator(
	runtime ports.WorkloadRuntime,
	renderer ports.WorkloadRenderer,
	admission ports.WorkloadAdmission,
	audit ports.WorkloadPlanAuditStore,
	dryRun ports.WorkloadProviderDryRun,
	apply ports.WorkloadProviderApply,
	reader ports.WorkloadProviderStatusReader,
	reconciler ports.WorkloadStatusReconciler,
	options ...InstanceOrchestratorOption,
) *LocalInstanceOrchestrator {
	orchestrator := &LocalInstanceOrchestrator{
		runtime:    runtime,
		renderer:   renderer,
		admission:  admission,
		audit:      audit,
		dryRun:     dryRun,
		apply:      apply,
		reader:     reader,
		reconciler: reconciler,
		now:        time.Now,
	}
	for _, option := range options {
		option(orchestrator)
	}
	return orchestrator
}

func (o *LocalInstanceOrchestrator) Create(ctx context.Context, request ports.WorkloadInstanceCreateRequest) (ports.WorkloadInstanceCreateResult, error) {
	if err := o.validate(); err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}
	if request.UserID == "" {
		return ports.WorkloadInstanceCreateResult{}, fmt.Errorf("%w: user id is required for instance orchestration", ports.ErrInvalid)
	}
	if request.PermissionProof == "" {
		return ports.WorkloadInstanceCreateResult{}, fmt.Errorf("%w: permission proof is required for instance orchestration", ports.ErrInvalid)
	}

	ref, err := o.runtime.Create(ctx, request.Spec)
	if err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}
	current, err := o.runtime.Get(ctx, ref)
	if err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}
	manifests, err := o.renderer.Render(ctx, request.Spec)
	if err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}
	admission, err := o.admission.Review(ctx, manifests)
	if err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}
	provider := ""
	if len(manifests) > 0 {
		provider = manifests[0].Provider
	}
	auditID, err := o.audit.RecordPlan(ctx, ports.WorkloadPlanAuditRecord{
		TenantID:        request.Spec.TenantID,
		UserID:          request.UserID,
		InstanceID:      ref.InstanceID,
		InstanceName:    request.Spec.Name,
		WorkloadKind:    request.Spec.Kind,
		Provider:        provider,
		Manifests:       manifests,
		AdmissionResult: admission,
		CreatedAt:       firstNonZeroTime(request.RequestedAt, o.now().UTC()),
	})
	if err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}
	dryRun, err := o.dryRun.DryRun(ctx, manifests, admission)
	if err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}
	apply, err := o.apply.Apply(ctx, ports.WorkloadProviderApplyRequest{
		TenantID:        request.Spec.TenantID,
		UserID:          request.UserID,
		InstanceID:      ref.InstanceID,
		AuditID:         auditID,
		PermissionProof: request.PermissionProof,
		Operation:       ports.WorkloadLifecycleCreate,
		Manifests:       manifests,
		AdmissionResult: admission,
		DryRunResult:    dryRun,
		RequestedAt:     firstNonZeroTime(request.RequestedAt, o.now().UTC()),
	})
	if err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}

	result := ports.WorkloadInstanceCreateResult{
		Ref:         ref,
		AuditID:     auditID,
		Manifests:   manifests,
		Admission:   admission,
		DryRun:      dryRun,
		Apply:       apply,
		FinalStatus: current,
	}
	if o.store != nil {
		if err := o.store.UpsertStatus(ctx, instanceRecordFromResult(request.Spec, ref, auditID, provider, nil, current, firstNonZeroTime(request.RequestedAt, o.now().UTC()))); err != nil {
			return ports.WorkloadInstanceCreateResult{}, err
		}
	}
	if !apply.Applied {
		return result, nil
	}

	observation, err := o.reader.Observe(ctx, ports.WorkloadProviderStatusRequest{
		TenantID:    request.Spec.TenantID,
		InstanceID:  ref.InstanceID,
		Kind:        request.Spec.Kind,
		ApplyResult: apply,
		RequestedAt: firstNonZeroTime(request.RequestedAt, o.now().UTC()),
	})
	if err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}
	reconcile, err := o.reconciler.Reconcile(ctx, ports.WorkloadReconcileRequest{
		AuditID:     auditID,
		Current:     current,
		ApplyResult: apply,
		Observation: observation,
	})
	if err != nil {
		return ports.WorkloadInstanceCreateResult{}, err
	}

	result.Observation = observation
	result.Reconcile = reconcile
	result.FinalStatus = reconcile.Status
	result.Orchestrated = true
	if o.store != nil {
		if err := o.store.UpsertStatus(ctx, instanceRecordFromResult(request.Spec, ref, auditID, provider, apply.ResourceRefs, reconcile.Status, firstNonZeroTime(request.RequestedAt, o.now().UTC()))); err != nil {
			return ports.WorkloadInstanceCreateResult{}, err
		}
	}
	return result, nil
}

func (o *LocalInstanceOrchestrator) validate() error {
	if o.runtime == nil {
		return fmt.Errorf("%w: workload runtime is required for instance orchestration", ports.ErrNotConfigured)
	}
	if o.renderer == nil {
		return fmt.Errorf("%w: workload renderer is required for instance orchestration", ports.ErrNotConfigured)
	}
	if o.admission == nil {
		return fmt.Errorf("%w: workload admission is required for instance orchestration", ports.ErrNotConfigured)
	}
	if o.audit == nil {
		return fmt.Errorf("%w: workload plan audit is required for instance orchestration", ports.ErrNotConfigured)
	}
	if o.dryRun == nil {
		return fmt.Errorf("%w: workload provider dry-run is required for instance orchestration", ports.ErrNotConfigured)
	}
	if o.apply == nil {
		return fmt.Errorf("%w: workload provider apply is required for instance orchestration", ports.ErrNotConfigured)
	}
	if o.reader == nil {
		return fmt.Errorf("%w: workload provider status reader is required for instance orchestration", ports.ErrNotConfigured)
	}
	if o.reconciler == nil {
		return fmt.Errorf("%w: workload status reconciler is required for instance orchestration", ports.ErrNotConfigured)
	}
	return nil
}

var _ ports.WorkloadInstanceOrchestrator = (*LocalInstanceOrchestrator)(nil)

func instanceRecordFromResult(spec ports.WorkloadSpec, ref ports.WorkloadRef, auditID string, provider string, resourceRefs []string, status ports.WorkloadStatus, createdAt time.Time) ports.WorkloadInstanceRecord {
	status.Ref = ref
	return ports.WorkloadInstanceRecord{
		TenantID:     spec.TenantID,
		InstanceID:   ref.InstanceID,
		Name:         spec.Name,
		Kind:         spec.Kind,
		Provider:     provider,
		AuditID:      auditID,
		ResourceRefs: append([]string(nil), resourceRefs...),
		Status:       status,
		CreatedAt:    createdAt,
		UpdatedAt:    firstNonZeroTime(status.UpdatedAt, createdAt),
	}
}
