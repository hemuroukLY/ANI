# M1-INSTANCE-I: Provider Status Reader and Instance Orchestrator

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.
- Continues the instance foundation before `模块 3：模型管理平台`.

## Scope

Added provider status reader and business-facing instance orchestration API:

- `pkg/ports/workload_runtime.go` adds `WorkloadProviderStatusReader` and
  `WorkloadInstanceOrchestrator`.
- `pkg/adapters/runtime/provider_status_reader.go`
- `pkg/adapters/runtime/provider_status_reader_test.go`
- `pkg/adapters/runtime/instance_orchestrator.go`
- `pkg/adapters/runtime/instance_orchestrator_test.go`
- `deploy/manifests/m1-instance-i/`
- `scripts/validate_instance_orchestrator.py`

Updated:

- `pkg/bootstrap/deps.go` wires `WorkloadStatus` and `WorkloadInstances`.
- `Makefile` adds `validate-instance-orchestrator` and wires it into
  `validate-infra`.

## Behavior

`LocalProviderStatusReader` normalizes applied provider resource refs into a
`WorkloadProviderObservation`.

`LocalInstanceOrchestrator` sequences the create path through ANI ports:

- `WorkloadRuntime`
- `WorkloadRenderer`
- `WorkloadAdmission`
- `WorkloadPlanAuditStore`
- `WorkloadProviderDryRun`
- `WorkloadProviderApply`
- `WorkloadProviderStatusReader`
- `WorkloadStatusReconciler`

If provider apply is disabled, orchestration stops before status observation and
returns the planned workload status.

Business services should call `WorkloadInstanceOrchestrator` instead of manually
sequencing provider renderer, dry-run, apply, status reader, or reconcile calls.

## Validation

Passed:

```bash
make validate-instance-orchestrator
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds new instance orchestration and provider status reader capabilities behind
  the workload runtime abstraction.

## Next Recommended Slice

- `M1-INSTANCE-J`: instance persistence/query API contract or real provider
  adapter integration.
