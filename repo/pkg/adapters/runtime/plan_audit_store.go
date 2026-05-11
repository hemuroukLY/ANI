package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kubercloud/ani/pkg/ports"
)

type MetadataPlanAuditStore struct {
	store ports.MetadataStore
	now   func() time.Time
}

type PlanAuditOption func(*MetadataPlanAuditStore)

func WithAuditClock(now func() time.Time) PlanAuditOption {
	return func(store *MetadataPlanAuditStore) {
		if now != nil {
			store.now = now
		}
	}
}

func NewMetadataPlanAuditStore(store ports.MetadataStore, options ...PlanAuditOption) *MetadataPlanAuditStore {
	auditStore := &MetadataPlanAuditStore{
		store: store,
		now:   time.Now,
	}
	for _, option := range options {
		option(auditStore)
	}
	return auditStore
}

func (s *MetadataPlanAuditStore) RecordPlan(ctx context.Context, record ports.WorkloadPlanAuditRecord) (string, error) {
	if s.store == nil {
		return "", ports.ErrNotConfigured
	}
	if strings.TrimSpace(record.TenantID) == "" {
		return "", fmt.Errorf("%w: tenantID is required", ports.ErrInvalid)
	}
	if strings.TrimSpace(record.InstanceName) == "" {
		return "", fmt.Errorf("%w: instanceName is required", ports.ErrInvalid)
	}
	if record.WorkloadKind == "" {
		return "", fmt.Errorf("%w: workloadKind is required", ports.ErrInvalid)
	}

	id := uuid.NewString()
	createdAt := record.CreatedAt
	if createdAt.IsZero() {
		createdAt = s.now().UTC()
	}
	manifests, err := json.Marshal(record.Manifests)
	if err != nil {
		return "", fmt.Errorf("marshal rendered manifests: %w", err)
	}
	warnings, err := json.Marshal(record.AdmissionResult.Warnings)
	if err != nil {
		return "", fmt.Errorf("marshal admission warnings: %w", err)
	}

	err = s.store.WithTenantTx(ctx, func(ctx context.Context, tx ports.MetadataTx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO instance_plan_audits (
				id, tenant_id, user_id, instance_id, instance_name, workload_kind,
				provider, manifest_count, rendered_manifests, admission_allowed,
				admission_reason, admission_warnings, created_at
			)
			VALUES (
				$1::uuid, $2::uuid, NULLIF($3, '')::uuid, NULLIF($4, ''),
				$5, $6, NULLIF($7, ''), $8, $9::jsonb, $10, NULLIF($11, ''),
				$12::jsonb, $13
			)
		`, id, record.TenantID, record.UserID, record.InstanceID, record.InstanceName,
			string(record.WorkloadKind), record.Provider, len(record.Manifests), string(manifests),
			record.AdmissionResult.Allowed, record.AdmissionResult.Reason, string(warnings), createdAt)
		if err != nil {
			return fmt.Errorf("insert instance plan audit: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

var _ ports.WorkloadPlanAuditStore = (*MetadataPlanAuditStore)(nil)
