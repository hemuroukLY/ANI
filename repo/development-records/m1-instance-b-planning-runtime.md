# M1-INSTANCE-B: Planning Runtime

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.
- Continues the instance foundation before `模块 3：模型管理平台`.

## Scope

Added the first executable instance runtime adapter:

- `pkg/adapters/runtime/planning.go`
- `pkg/adapters/runtime/planning_test.go`
- `deploy/manifests/m1-instance-b/`
- `scripts/validate_instance_planner.py`

Updated:

- `pkg/bootstrap/deps.go` now wires `WorkloadRuntime` to `PlanningRuntime`.
- `pkg/ports/errors.go` adds `ErrInvalid`.
- `Makefile` adds `validate-instance-planner` and wires it into `validate-infra`.

## Behavior

`PlanningRuntime` validates and normalizes instance specs before a real provider
adapter creates resources:

- Tenant, name, kind, image, VM boot image, and root disk validation.
- Default network plane normalization by instance kind.
- Required `tenant_vpc` business connectivity.
- Storage attachment validation, including `object_fuse` source references.
- GPU container / inference dependency on `GPUInventory`.
- Lifecycle transitions for start, stop, restart, resize, and delete.

It does not create Kubernetes Pod, Deployment, Job, KubeVirt VM, PVC, VPC, or
Subnet resources.

## Validation

Passed:

```bash
make validate-instance-planner
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds a new runtime adapter implementation and validation behavior behind the
  existing `WorkloadRuntime` port.

## Next Recommended Slice

- `M1-INSTANCE-C`: Kubernetes/KubeVirt provider adapter boundary and dry-run
  rendering.
