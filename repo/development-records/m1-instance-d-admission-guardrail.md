# M1-INSTANCE-D: Admission Guardrail

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.
- Continues the instance foundation before `模块 3：模型管理平台`.

## Scope

Added local admission guardrail:

- `pkg/ports/workload_runtime.go` adds `WorkloadAdmission`.
- `pkg/adapters/runtime/admission.go`
- `pkg/adapters/runtime/admission_test.go`
- `deploy/manifests/m1-instance-d/`
- `scripts/validate_instance_admission.py`

Updated:

- `pkg/bootstrap/deps.go` wires `WorkloadAdmission`.
- `Makefile` adds `validate-instance-admission` and wires it into `validate-infra`.

## Behavior

The local guardrail reviews renderer output before future server-side dry-run or
real provider create/apply:

- Allows only KubeVirt `VirtualMachine`, Kubernetes `Deployment`, and
  Kubernetes `Job`.
- Requires tenant and instance labels.
- Requires `ani.kubercloud.io/render-mode=dry-run`.
- Requires network plane annotation.
- Denies `hostNetwork=true`.
- Denies privileged containers.

It does not create cluster resources.

## Validation

Passed:

```bash
make validate-instance-admission
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds a new admission capability and local guardrail behind the instance
  runtime abstraction.

## Next Recommended Slice

- `M1-INSTANCE-E`: instance plan/admission persistence and audit trail.
