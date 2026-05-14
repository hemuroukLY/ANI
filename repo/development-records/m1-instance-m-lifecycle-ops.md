# M1-INSTANCE-M: Lifecycle and Visual Ops API

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added lifecycle and visual operations API coverage:

- `pkg/ports/workload_runtime.go` extends `WorkloadInstanceService` with
  Start/Stop/Restart/Resize/Delete/Ops.
- `pkg/ports/workload_runtime.go` adds `WorkloadInstanceOps`.
- `pkg/adapters/runtime/instance_ops.go`
- `pkg/adapters/runtime/instance_service.go`
- `pkg/adapters/runtime/instance_service_test.go`
- `deploy/manifests/m1-instance-m/`
- `scripts/validate_instance_lifecycle_ops.py`

Updated:

- `pkg/bootstrap/deps.go` wires `InstanceOps`.
- `Makefile` adds `validate-instance-lifecycle-ops` and wires it into
  `validate-infra`.

## Behavior

Lifecycle API coverage:

- Start
- Stop
- Restart
- Resize
- Delete

Visual ops API coverage:

- logs
- events
- metrics
- terminal
- exec

Ops are disabled by default. Terminal and exec are container-only and require a
running instance. Production implementations must route provider-specific calls
through adapter-owned packages.

## Validation

Passed:

```bash
make validate-instance-lifecycle-ops
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.
