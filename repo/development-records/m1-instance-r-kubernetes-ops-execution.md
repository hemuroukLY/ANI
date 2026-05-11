# M1-INSTANCE-R: Kubernetes Visual Ops Execution

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added provider-side visual ops execution:

- `pkg/adapters/runtime/kubernetes_instance_ops.go`
- `pkg/adapters/runtime/kubernetes_instance_ops_test.go`
- `pkg/bootstrap/deps.go`
- `pkg/bootstrap/server.go`
- `deploy/manifests/m1-instance-r/`
- `scripts/validate_kubernetes_ops_execution.py`
- `Makefile` target `validate-kubernetes-ops-execution`

Updated:

- `deploy/helm/ani-platform/values.yaml`
- `CLAUDE.md`
- `AGENTS.md`
- `ANI-06-开发计划.md`
- `ANI-13-开源组件松耦合适配器架构.md`
- `development-records/README.md`

## Behavior

Default behavior remains local/offline.

When configured with `WORKLOAD_OPS_PROVIDER=kubernetes_rest`,
`KubernetesInstanceOps` uses `KubernetesRESTClient` for:

- logs through the pod log subresource
- events through the namespace events API
- metrics through `metrics.k8s.io`
- terminal and exec through the pod exec subresource

Supported environment variables:

- `WORKLOAD_OPS_PROVIDER`
- `WORKLOAD_OPS_ENABLED`
- `KUBERNETES_API_HOST`
- `KUBERNETES_BEARER_TOKEN`
- `KUBERNETES_PROVIDER_FIELD_MANAGER`

Ops execution remains disabled by default.

## Validation

Passed:

```bash
make validate-kubernetes-ops-execution
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

Reason: this adds a provider ops adapter and bootstrap wiring while preserving
default local/offline behavior.

## Next Work Boundary

Recommended next options:

- Execute M1 real-provider integration regression profile.
- Start `M3-MODEL-A` after accepting the M1/M2/GPU/Runtime/Instance baseline.
