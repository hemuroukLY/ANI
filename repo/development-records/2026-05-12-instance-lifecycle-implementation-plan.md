# Instance Lifecycle Implementation Plan

Date: 2026-05-12

## Purpose

This plan decomposes the AWS-derived instance lifecycle reference into
implementation slices for ANI.

The current implementation proves the minimal M1 control chain:

```text
plan -> render -> admission -> audit -> dry-run -> gated apply -> observe
     -> reconcile -> persist -> service API -> demo console
```

That chain is correct and must be preserved. The next work is to enrich product
depth so VM, container, GPU container, inference, notebook, sandbox, and batch
instances behave like a production cloud platform rather than a minimal demo.

## Planning Principles

1. Keep one ANI product concept: "Instance".
2. Keep provider details behind ANI ports and adapters.
3. Add cross-cutting product semantics before per-kind feature expansion.
4. Every operation must have precheck, operation ID, state transition,
   timeline, failure reason, and audit.
5. Console click behavior is part of the backend contract, not frontend-only
   polish.
6. Mock pages may guide information architecture, but implementation readiness
   is judged only by production APIs, persistence, tests, and validation.

## Target Capability Layers

### Layer 0: Existing Baseline

Already completed:

- Instance kinds: VM, container, GPU container, inference, notebook,
  agent sandbox, batch job are present in the runtime abstraction.
- Create chain with planning, manifest render, admission, audit, dry-run,
  gated apply, observation, reconcile, and persistence.
- Lifecycle service entry points: start, stop, restart, resize, delete.
- Visual ops entry points: logs, events, metrics, terminal, exec, VM console.
- Kubernetes REST provider adapter and local/offline profile tests.
- Demo workspace with VM/container/GPU container flows.

Baseline limitation:

- It is a technical minimum viable chain, not full product behavior.

### Layer 1: Product Semantics Core

Cross-cutting features required before deep per-kind work:

- Expanded instance states and provider state mapping.
- State reason, state message, health summary, and recommended action.
- Operation object and operation timeline.
- Action precheck and disabled-action reason.
- Destructive-action confirmation metadata.
- Per-action support matrix by kind/provider/state.
- Instance event stream and event persistence.
- Audit and before/after spec diff for every operation.
- Idempotency key for create/update/lifecycle operations.

### Layer 2: Resource Domain Depth

Common resource domains required by all instance kinds:

- Compute profiles and quota impact.
- Image/model/runtime source.
- Network interfaces, security groups, ports, ingress, and egress policies.
- Storage attachments, retention policy, snapshots, backups, and restore.
- Tags/labels, owner, project, cost center.
- Metrics/logs/events/alarms.
- Connection/session metadata.

### Layer 3: Kind-Specific Depth

Per-kind product behavior:

- VM: SSH/RDP/managed session, VNC/serial console, console output,
  screenshot, snapshots, image creation, disk/network attach, termination
  protection, hibernate/suspend if supported.
- Container: replicas, rollout, rollback, deployment history, restart
  individual replica, image update, previous logs, exec audit.
- GPU container: GPU inventory explanation, queue position, GPU slice profile,
  utilization, memory, temperature, scheduling failure reason.
- Inference: model version, endpoint, autoscaling, traffic policy, canary,
  rollback, token limits, invocation logging, usage/cost.
- Notebook: app/runtime state separated from persistent workspace, browser
  session, idle shutdown, workspace deletion warning.
- Agent sandbox: restricted egress, tool policy, session isolation, command
  transcript.
- Batch/training: queue, priority, attempts, retry, timeout, dependencies,
  checkpoint/output links.

### Layer 4: Production Console

Frontend surfaces should be promoted only after API contracts exist:

- Instance list with state reason and actionable health.
- Detail tabs: overview, configuration, metrics, logs, events, network,
  storage, security, operations, audit.
- Create wizard with validation and review plan.
- Action buttons with precheck, confirmation, operation progress, and retry.
- Operation drawer/timeline.
- Failure explanation with recommended fix.

## Implementation Slices

### M1-INSTANCE-T: Instance Operation Semantics

Product mapping:

- `ANI-06 / 模块 1：基础设施底座 / Workload Runtime / Instance Fabric`.

Goal:

- Add the cross-cutting product contract for action precheck, operation
  tracking, timeline, state reason, disabled-action reason, and destructive
  confirmation metadata.

Backend changes:

- Extend `pkg/ports/workload_runtime.go` with:
  - expanded `WorkloadState` values where stable.
  - `WorkloadHealth`.
  - `WorkloadStateReason`.
  - `WorkloadOperation`.
  - `WorkloadOperationStep`.
  - `WorkloadActionPrecheck`.
  - `WorkloadActionSupport`.
  - `WorkloadDestructiveImpact`.
- Add `WorkloadOperationStore` port.
- Add metadata-backed operation persistence table:
  - `workload_instance_operations`.
  - `workload_instance_operation_steps`.
- Add local operation store adapter using `MetadataStore`.
- Update `WorkloadInstanceService` lifecycle methods to return operation
  metadata instead of only final status where needed.
- Keep existing synchronous test behavior as a compatibility path until Gateway
  routes are promoted.

Frontend/API changes:

- Demo/prod API contract should expose:
  - allowed actions.
  - disabled action reasons.
  - operation ID.
  - operation timeline.
  - destructive impact summary.

Validation:

- Unit tests for state/action matrix.
- Operation store persistence tests with RLS tenant context.
- Existing `make validate-infra`, `make test`, and `make build` pass.

Release impact:

- Expected `MINOR` because public internal contracts and schema expand.

### M1-INSTANCE-U: VM Production Operations

Goal:

- Bring VM operations closer to EC2-level management while preserving provider
  adapter boundaries.

Backend changes:

- Extend VM spec/status with:
  - termination protection.
  - boot image/template reference.
  - root disk retention.
  - data disk refs.
  - network interface refs.
  - login/connection profile.
  - console session capability.
- Extend `WorkloadInstanceOps` actions:
  - `vm_console`.
  - `vm_serial_console`.
  - `vm_screenshot`.
  - `vm_console_output`.
  - `vm_managed_session`.
- Add contracts for:
  - create snapshot.
  - create image/template.
  - attach/detach disk.
  - attach/detach network interface.

Provider work:

- KubeVirt adapter maps VNC/serial console/session metadata.
- OpenStack/VMware/cloud adapters can be added later without changing business
  service contracts.

Frontend/API changes:

- VM detail tabs gain:
  - console.
  - disks.
  - snapshots.
  - network interfaces.
  - protection and retention summary.

Validation:

- VM action precheck tests.
- Provider adapter tests using fake provider transport.
- Deletion requires confirmation metadata when persistent disks may be retained
  or destroyed.

Release impact:

- Expected `MINOR`.

### M1-INSTANCE-V: Container and GPU Container Operations

Goal:

- Add service/deployment-style operations for containers and GPU containers.

Backend changes:

- Add container deployment status:
  - desired replicas.
  - ready replicas.
  - unavailable replicas.
  - rollout revision.
  - deployment history.
  - per-replica health.
- Add actions:
  - scale replicas.
  - rollout image.
  - rollback.
  - pause/resume rollout.
  - restart replica.
  - view previous logs.
  - run one-off command.
- Add GPU-specific status:
  - vendor/model.
  - resource name.
  - slice profile.
  - queue.
  - scheduling reason.
  - allocation and utilization summary.

Provider work:

- Kubernetes adapter maps Deployment status, ReplicaSet revision, Pod events,
  previous logs, metrics, and GPU node allocation.
- GPUInventory remains the only business-facing GPU scheduling dependency.

Frontend/API changes:

- Container detail gets rollout and replica panels.
- GPU container detail gets scheduling and GPU utilization panels.
- Disabled actions explain provider or state limitations.

Validation:

- Local/offline profile for rollout and rollback contracts.
- Kubernetes fake transport tests for status mapping.
- GPU scheduling reason tests.

Release impact:

- Expected `MINOR`.

### M3-MODEL-A: Model Metadata and Object Storage Boundary

Goal:

- Implement production model metadata and object-store boundaries before
  inference lifecycle work.

Backend changes:

- Keep model files behind `ObjectStore`.
- Implement signed upload/download URL flows through `model-service`.
- Add object path policy:
  - tenant ID.
  - model ID.
  - version.
  - file name.
- Add checksum/size validation gates before version becomes ready.
- Ensure model metadata uses PostgreSQL RLS.

Provider work:

- Default object store adapter can be MinIO/S3-compatible, but SDK must remain
  in adapter/bootstrap packages.

Frontend/API changes:

- Model create/upload/version pages move from mock to production API.

Validation:

- Model-service unit tests.
- ObjectStore fake adapter tests.
- `make validate-architecture` confirms no MinIO/S3 SDK leak in business code.

Release impact:

- Expected `MINOR`.

### M3-INFERENCE-A: Inference Instance Lifecycle

Goal:

- Promote inference from future kind to production first-class instance.

Dependencies:

- `M1-INSTANCE-T`.
- `M3-MODEL-A`.

Backend changes:

- Add inference spec:
  - model version.
  - runtime engine.
  - replicas.
  - GPU/NPU profile.
  - endpoint mode.
  - autoscaling policy.
  - traffic policy.
  - token/request limits.
  - invocation logging policy.
- Implement create/update/delete through `WorkloadRuntime` and
  `WorkloadInstanceService`.
- Do not let model-service directly create Kubernetes Deployment resources.
- Add inference status:
  - endpoint URL.
  - route registration.
  - model loading status.
  - warmup status.
  - replica health.
  - traffic split.

Provider work:

- Inference operator or Kubernetes adapter renders workloads behind runtime
  ports.
- Metrics route through observability adapters.

Frontend/API changes:

- Inference detail tabs:
  - endpoint.
  - model.
  - replicas.
  - traffic.
  - metrics.
  - invocation logs.
  - usage/cost.

Validation:

- End-to-end contract test for model version -> inference instance -> endpoint
  status.
- Update/rollback state machine tests.

Release impact:

- Expected `MINOR`.

### M3-NOTEBOOK-A: Notebook Workspace Lifecycle

Goal:

- Make notebook/workspace behavior explicit and production-safe.

Dependencies:

- `M1-INSTANCE-T`.
- Object storage/PVC contracts from M1/M3.

Backend changes:

- Separate app/runtime state from workspace storage state.
- Add idle shutdown policy.
- Add browser session operation.
- Add persistent workspace retention/deletion confirmation.
- Add Git/object mount references.

Frontend/API changes:

- Notebook page exposes:
  - start/stop app.
  - open browser.
  - workspace storage.
  - idle policy.
  - delete workspace warning.

Validation:

- Stopping app preserves storage.
- Deleting workspace requires destructive confirmation metadata.

Release impact:

- Expected `MINOR`.

### M3-BATCH-A: Batch and Training Job Lifecycle

Goal:

- Implement AWS Batch-like lifecycle visibility for queued jobs and training.

Dependencies:

- `M1-INSTANCE-T`.
- GPU scheduling contracts.

Backend changes:

- Add job state mapping:
  - submitted.
  - pending.
  - runnable.
  - starting.
  - running.
  - retrying.
  - succeeded.
  - failed.
  - cancelled.
- Add queue, priority, retry, timeout, dependencies, attempts, and checkpoint
  paths.
- Add actions:
  - cancel.
  - stop.
  - retry.
  - clone and rerun.

Provider work:

- Kubernetes Job/Volcano/Kueue-style adapter maps queue and attempt state.

Frontend/API changes:

- Job detail page exposes:
  - queue position.
  - attempts.
  - logs.
  - output artifacts.
  - checkpoints.

Validation:

- State mapping tests.
- Retry/cancel/timeout tests.

Release impact:

- Expected `MINOR`.

### M4-CONSOLE-INSTANCE-A: Production Instance Console

Goal:

- Promote the Demo information architecture into production Console surfaces
  backed by real APIs.

Dependencies:

- `M1-INSTANCE-T`.
- At least one kind-specific depth slice.

Frontend changes:

- Replace mock instance resource-domain data with production endpoints.
- Add create wizard with review plan.
- Add operation drawer/timeline.
- Add disabled-action tooltips/reasons.
- Add destructive confirmations with retention and cost impact.
- Add detail tabs:
  - overview.
  - configuration.
  - metrics.
  - logs.
  - events.
  - network.
  - storage.
  - security.
  - operations.
  - audit.

Backend/API changes:

- Gateway routes expose stable REST contracts over the service layer.
- API must never expose provider manifests or provider SDK-specific objects.

Validation:

- `npm run build`.
- Demo/prod route smoke tests.
- API contract tests for action feedback.

Release impact:

- Expected `MINOR`.

## Recommended Execution Order

1. `M1-INSTANCE-T`: operation semantics and product feedback foundation.
2. `M1-INSTANCE-U`: VM production operations, because VM is the clearest Demo
   gap and maps well to customer expectations.
3. `M1-INSTANCE-V`: container/GPU rollout, metrics, and scheduling depth.
4. `M3-MODEL-A`: model metadata and object storage boundary.
5. `M3-INFERENCE-A`: inference instance lifecycle.
6. `M3-NOTEBOOK-A`: notebook/workspace lifecycle.
7. `M3-BATCH-A`: batch/training job lifecycle.
8. `M4-CONSOLE-INSTANCE-A`: production console promotion as APIs mature.

Rationale:

- Start with shared semantics so every later feature gets consistent action
  feedback, audit, and timeline behavior.
- VM depth gives the fastest visible improvement for Demo/product discussion.
- Container/GPU depth closes the AI infrastructure credibility gap.
- Model/inference work depends on object storage and model metadata readiness.
- Console promotion should follow API maturity, not lead it with mock-only
  behavior.

## P0 / P1 / P2 Product Priority

### P0

- State reason and health summary.
- Action precheck and disabled-action reason.
- Operation ID and timeline.
- Destructive action confirmation metadata.
- Logs/events/metrics on every instance.
- VM console/session contract.
- Container exec audit.
- GPU scheduling reason.
- Inference endpoint lifecycle basics.

### P1

- VM snapshots/backups/image creation.
- Container rollout history and rollback.
- GPU utilization panels.
- Notebook app/workspace separation.
- Batch attempts/retries/queue position.
- Inference autoscaling and traffic split.

### P2

- VM hibernate/suspend/recovery.
- Cross-AZ placement policies.
- Advanced alarms and scheduled events.
- Canary/blue-green inference rollout.
- Cost allocation and application-level token accounting.

## Data Model Additions

Minimum tables expected for `M1-INSTANCE-T`:

```sql
workload_instance_operations (
  id uuid primary key,
  tenant_id uuid not null,
  instance_id uuid not null,
  operation text not null,
  status text not null,
  requested_by uuid not null,
  idempotency_key text,
  precheck_json jsonb not null default '{}',
  destructive_impact_json jsonb not null default '{}',
  before_spec_json jsonb not null default '{}',
  after_spec_json jsonb not null default '{}',
  provider_refs_json jsonb not null default '[]',
  reason text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

workload_instance_operation_steps (
  id uuid primary key,
  tenant_id uuid not null,
  operation_id uuid not null,
  step_name text not null,
  status text not null,
  message text,
  started_at timestamptz,
  completed_at timestamptz,
  created_at timestamptz not null default now()
);
```

Both tables must use tenant RLS. Operation queries must not depend on
in-memory planner state.

## API Shape

Recommended REST surface for production Gateway:

```text
GET    /api/v1/instances
POST   /api/v1/instances
GET    /api/v1/instances/{id}
GET    /api/v1/instances/{id}/actions
POST   /api/v1/instances/{id}/actions/{action}/precheck
POST   /api/v1/instances/{id}/actions/{action}
GET    /api/v1/instances/{id}/operations
GET    /api/v1/instance-operations/{operation_id}
GET    /api/v1/instances/{id}/logs
GET    /api/v1/instances/{id}/events
GET    /api/v1/instances/{id}/metrics
POST   /api/v1/instances/{id}/sessions
```

Responses should include:

- request ID.
- operation ID when an operation is accepted.
- state and state reason.
- disabled actions and reasons.
- destructive impact summary when relevant.
- retry eligibility for failed operations.

## Validation Gate For Each Slice

Required before closing each implementation slice:

```bash
make gen-proto
make validate-infra
make test
make build
git diff --check
```

Frontend slices additionally require:

```bash
cd frontends/console
npm run build
```

Each slice must update:

- `repo/development-records/README.md`.
- A dedicated `repo/development-records/<slice>.md`.
- `ANI-06-开发计划.md` if product-plan wording or mapping changes.

## Release Impact

This document is planning only.

Expected release impact: `no-release-impact`.

