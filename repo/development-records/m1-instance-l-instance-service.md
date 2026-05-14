# M1-INSTANCE-L: Instance Service API

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added business-facing instance service API layer:

- `pkg/ports/workload_runtime.go` adds `WorkloadInstanceService`.
- `pkg/adapters/runtime/instance_service.go`
- `pkg/adapters/runtime/instance_service_test.go`
- `deploy/manifests/m1-instance-l/`
- `scripts/validate_instance_service.py`

Updated:

- `pkg/bootstrap/deps.go` wires `InstanceService`.
- `Makefile` adds `validate-instance-service` and wires it into
  `validate-infra`.

## Behavior

`LocalInstanceService` supports Create/Get/List for:

- VM instances.
- Container instances.
- GPU container instances.

Create uses `WorkloadInstanceOrchestrator`. Get/List use
`WorkloadInstanceStore`. The service API does not expose provider manifests,
provider SDK objects, or provider-specific status details.

## Validation

Passed:

```bash
make validate-instance-service
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.
