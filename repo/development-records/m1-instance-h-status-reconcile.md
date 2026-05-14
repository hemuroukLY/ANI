# M1-INSTANCE-H: Status Reconcile

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.
- Continues the instance foundation before `模块 3：模型管理平台`.

## Scope

Added status write-back and lifecycle reconcile boundary:

- `pkg/ports/workload_runtime.go` adds `WorkloadStatusReconciler`.
- `pkg/adapters/runtime/status_reconciler.go`
- `pkg/adapters/runtime/status_reconciler_test.go`
- `deploy/manifests/m1-instance-h/`
- `scripts/validate_instance_status_reconcile.py`

Updated:

- `pkg/bootstrap/deps.go` wires `WorkloadReconcile`.
- `Makefile` adds `validate-instance-status-reconcile` and wires it into
  `validate-infra`.

## Behavior

`LocalStatusReconciler` converts normalized provider observations into ANI
`WorkloadState`.

It requires:

- audit id.
- applied provider apply result.
- current workload tenant and instance id.
- observation tenant, instance, kind, and provider correlation.
- observation resource refs that match provider apply resource refs.

The local mapper covers provisioning, running, stopping, stopped, deleting,
deleted, and failed states.

Provider-specific status readers belong in adapters. Business services must not
poll Kubernetes, KubeVirt, or customer cloud status APIs directly.

## Validation

Passed:

```bash
make validate-instance-status-reconcile
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds a new status reconcile capability behind the instance runtime
  abstraction.

## Next Recommended Slice

- `M1-INSTANCE-I`: provider status reader and instance orchestration API
  contract.
