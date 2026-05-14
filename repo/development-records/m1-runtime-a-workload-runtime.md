# M1-RUNTIME-A: Workload Runtime / Instance Abstraction

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.
- Required foundation before `模块 3：模型管理平台`.

## Scope

Added workload runtime abstraction and contracts:

- `pkg/ports/workload_runtime.go`
- `pkg/adapters/runtime/not_configured.go`
- `deploy/manifests/m1-runtime-a/`
- `deploy/helm/ani-platform/component-contracts/workload-runtime.yaml`
- `deploy/helm/ani-platform/profiles/runtime-foundation.yaml`
- `scripts/validate_runtime_contracts.py`

The abstraction covers:

- Traditional VM / cloud host instances.
- Traditional container instances.
- GPU container instances.
- Inference instances.
- Notebook instances.
- AI Agent sandbox instances.
- Batch job instances.

Model management and future operators must use `ports.WorkloadRuntime` rather
than directly creating Pod, Deployment, Job, RuntimeClass, or KubeVirt resources
from business workflows.

## Validation

Passed:

```bash
make validate-runtime-contracts
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds a new capability port and runtime contracts before user-facing runtime
  implementation.

## Next Recommended Slice

- `M3-MODEL-A`: model metadata and object-storage boundary now that GPU and
  runtime foundations exist.
