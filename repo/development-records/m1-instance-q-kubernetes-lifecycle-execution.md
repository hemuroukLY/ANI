# M1-INSTANCE-Q: Kubernetes Lifecycle Execution

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added provider-side lifecycle execution for existing instances:

- `pkg/ports/workload_runtime.go`
- `pkg/adapters/runtime/instance_service.go`
- `pkg/adapters/runtime/kubernetes_lifecycle_executor.go`
- `pkg/adapters/runtime/kubernetes_lifecycle_executor_test.go`
- `pkg/bootstrap/deps.go`
- `pkg/bootstrap/server.go`
- `deploy/manifests/m1-instance-q/`
- `scripts/validate_kubernetes_lifecycle_execution.py`
- `Makefile` target `validate-kubernetes-lifecycle-execution`

Updated:

- `deploy/helm/ani-platform/values.yaml`
- `CLAUDE.md`
- `AGENTS.md`
- `ANI-06-开发计划.md`
- `ANI-13-开源组件松耦合适配器架构.md`
- `development-records/README.md`

## Behavior

Added `WorkloadInstanceLifecycleExecutor` as the provider execution boundary for
already-created instances.

Default behavior remains local/offline. `LocalInstanceService` uses the
executor only when one is configured.

`KubernetesLifecycleExecutor` uses `KubernetesRESTClient` for:

- Deployment/Job start and stop through the scale subresource
- Deployment/Job restart through pod template annotation patch
- resource delete through Kubernetes `DELETE`
- KubeVirt VirtualMachine start/stop execution boundary

Supported environment variables:

- `WORKLOAD_LIFECYCLE_PROVIDER`
- `WORKLOAD_LIFECYCLE_APPLY_ENABLED`
- `KUBERNETES_API_HOST`
- `KUBERNETES_BEARER_TOKEN`
- `KUBERNETES_PROVIDER_FIELD_MANAGER`

Lifecycle execution remains disabled by default.

## Validation

Passed:

```bash
make validate-kubernetes-lifecycle-execution
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

Reason: this adds a new provider lifecycle execution port and adapter while
preserving default local/offline behavior.

## Next Work Boundary

Recommended next options:

- Extend visual ops provider execution.
- Start `M3-MODEL-A` after accepting the M1/M2/GPU/Runtime/Instance baseline.
