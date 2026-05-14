# M1-INSTANCE-G: Provider Apply Gate

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.
- Continues the instance foundation before `模块 3：模型管理平台`.

## Scope

Added provider apply/create execution gate:

- `pkg/ports/workload_runtime.go` adds `WorkloadProviderApply`.
- `pkg/adapters/runtime/provider_apply.go`
- `pkg/adapters/runtime/provider_apply_test.go`
- `deploy/manifests/m1-instance-g/`
- `scripts/validate_instance_provider_apply.py`

Updated:

- `pkg/bootstrap/deps.go` wires `WorkloadApply`.
- `Makefile` adds `validate-instance-provider-apply` and wires it into
  `validate-infra`.

## Behavior

`LocalProviderApply` is disabled by default and fails closed.

When explicitly enabled, it requires:

- tenant id, user id, instance id, and audit id.
- permission proof from the authorization layer.
- allowed admission result.
- accepted provider dry-run result.
- matching dry-run provider and manifest count.
- a create operation.
- provider/kind/apiVersion combinations already accepted by the dry-run gate.

The local gate returns deterministic resource refs for tests. It is not a
Kubernetes/KubeVirt client.

## Validation

Passed:

```bash
make validate-instance-provider-apply
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds a new controlled execution capability behind the instance runtime
  abstraction.

## Next Recommended Slice

- `M1-INSTANCE-H`: instance status write-back and lifecycle reconcile contract.
