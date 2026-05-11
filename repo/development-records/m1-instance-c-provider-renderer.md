# M1-INSTANCE-C: Provider Dry-Run Renderer

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.
- Continues the instance foundation before `模块 3：模型管理平台`.

## Scope

Added provider dry-run rendering:

- `pkg/ports/workload_runtime.go` adds `WorkloadManifest` and `WorkloadRenderer`.
- `pkg/adapters/runtime/dryrun_renderer.go`
- `pkg/adapters/runtime/dryrun_renderer_test.go`
- `deploy/manifests/m1-instance-c/`
- `scripts/validate_instance_renderer.py`

Updated:

- `pkg/bootstrap/deps.go` wires `WorkloadRenderer`.
- `Makefile` adds `validate-instance-renderer` and wires it into `validate-infra`.

## Behavior

The renderer calls `PlanningRuntime` first, then emits reviewable provider
manifests:

- VM -> KubeVirt `VirtualMachine`.
- Container / GPU container / inference / notebook / sandbox -> Kubernetes `Deployment`.
- Batch job -> Kubernetes `Job`.

The renderer does not create cluster resources. Output is JSON-formatted
Kubernetes manifest content, which is valid YAML and suitable for review or
future server-side dry-run validation.

## Validation

Passed:

```bash
make validate-instance-renderer
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds a new renderer capability and dry-run provider mapping behind the
  instance runtime abstraction.

## Next Recommended Slice

- `M1-INSTANCE-D`: provider server-side dry-run / admission guardrail.
