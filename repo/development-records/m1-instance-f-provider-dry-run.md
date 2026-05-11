# M1-INSTANCE-F: Provider Dry-Run Executor

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.
- Continues the instance foundation before `模块 3：模型管理平台`.

## Scope

Added provider dry-run execution boundary:

- `pkg/ports/workload_runtime.go` adds `WorkloadProviderDryRun`.
- `pkg/adapters/runtime/provider_dryrun.go`
- `pkg/adapters/runtime/provider_dryrun_test.go`
- `deploy/manifests/m1-instance-f/`
- `scripts/validate_instance_provider_dry_run.py`

Updated:

- `pkg/bootstrap/deps.go` wires `WorkloadDryRun`.
- `Makefile` adds `validate-instance-provider-dry-run` and wires it into `validate-infra`.

## Behavior

`LocalProviderDryRun` validates provider/kind/apiVersion mappings after local
admission:

- `kubevirt` -> `kubevirt.io/v1` `VirtualMachine`.
- `kubernetes` -> `apps/v1` `Deployment`.
- `kubernetes` -> `batch/v1` `Job`.

It rejects admission-denied requests, mixed provider batches, unknown providers,
and invalid provider/kind/API-version combinations.

It does not create cluster resources.

## Validation

Passed:

```bash
make validate-instance-provider-dry-run
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds a new dry-run capability and local implementation behind the instance
  runtime abstraction.

## Next Recommended Slice

- `M1-INSTANCE-G`: real provider apply/create execution switch, permission
  checks, and audit closure.
