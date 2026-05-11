# M1-INSTANCE-P: Kubernetes Provider Bootstrap Wiring

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Wired provider selection into bootstrap/config:

- `pkg/bootstrap/server.go`
- `pkg/bootstrap/deps.go`
- `pkg/bootstrap/deps_test.go`
- `deploy/manifests/m1-instance-p/`
- `scripts/validate_kubernetes_bootstrap_wiring.py`
- `Makefile` target `validate-kubernetes-bootstrap-wiring`

Updated:

- `deploy/helm/ani-platform/values.yaml`
- `CLAUDE.md`
- `AGENTS.md`
- `ANI-06-开发计划.md`
- `ANI-13-开源组件松耦合适配器架构.md`
- `development-records/README.md`

## Behavior

Default behavior remains local/offline:

- `WorkloadProviderDryRun` -> `LocalProviderDryRun`
- `WorkloadProviderApply` -> `LocalProviderApply`
- `WorkloadProviderStatusReader` -> `LocalProviderStatusReader`

When configured with `WORKLOAD_PROVIDER=kubernetes_rest`,
`KubernetesRESTClient` is wrapped by `KubernetesProviderAdapter` and exposed as:

- `WorkloadProviderDryRun`
- `WorkloadProviderApply`
- `WorkloadProviderStatusReader`

Supported environment variables:

- `WORKLOAD_PROVIDER`
- `WORKLOAD_PROVIDER_APPLY_ENABLED`
- `KUBERNETES_API_HOST`
- `KUBERNETES_BEARER_TOKEN`
- `KUBERNETES_PROVIDER_FIELD_MANAGER`

Apply remains disabled by default.

## Validation

Passed:

```bash
make validate-kubernetes-bootstrap-wiring
```

Full repository validation should include:

```bash
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason: this adds provider configuration and bootstrap wiring while preserving
the default local provider behavior.

## Next Work Boundary

Recommended next options:

- Extend real lifecycle provider execution.
- Extend visual ops provider execution.
- Start `M3-MODEL-A` after accepting the M1/M2/GPU/Runtime/Instance baseline.
