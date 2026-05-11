# M1-INSTANCE-K: Kubernetes/KubeVirt Provider Adapter

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added provider adapter boundary:

- `pkg/adapters/runtime/kubernetes_provider_adapter.go`
- `pkg/adapters/runtime/kubernetes_provider_adapter_test.go`
- `deploy/manifests/m1-instance-k/`
- `scripts/validate_instance_provider_adapter.py`

## Behavior

`KubernetesProviderAdapter` implements:

- `WorkloadProviderDryRun`
- `WorkloadProviderApply`
- `WorkloadProviderStatusReader`

It delegates real provider work to an internal `KubernetesProviderClient`.
Apply remains disabled by default. Dry-run requires admission and provider
manifest validation before calling server-side dry-run.

## Validation

Passed:

```bash
make validate-instance-provider-adapter
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.
