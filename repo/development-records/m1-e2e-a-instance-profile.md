# M1-E2E-A: M1 End-to-End Integration Profile

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added an offline M1 end-to-end integration profile for the core instance
objects:

- VM
- container
- GPU container

Implemented:

- `pkg/adapters/runtime/m1_e2e_profile_test.go`
- `deploy/manifests/m1-e2e/`
- `scripts/validate_m1_e2e_profile.py`
- `Makefile` target `validate-m1-e2e`

Updated:

- `CLAUDE.md`
- `AGENTS.md`
- `ANI-06-开发计划.md`
- `development-records/README.md`

## Behavior

The profile wires the M1 contract chain through local/offline adapters:

- `WorkloadRuntime`
- `WorkloadRenderer`
- `WorkloadAdmission`
- `WorkloadPlanAuditStore`
- `WorkloadProviderDryRun`
- `WorkloadProviderApply`
- `WorkloadProviderStatusReader`
- `WorkloadStatusReconciler`
- `WorkloadInstanceStore`
- `WorkloadInstanceService`
- `WorkloadInstanceOps`

The test covers:

- create for VM, container, and GPU container
- Start/Stop/Restart/Resize/Delete lifecycle operations
- Get/List query operations
- container logs and terminal operations
- GPU container metrics and exec operations
- VM terminal rejection

The profile does not require a real Kubernetes cluster. Production execution
must replace the local provider gate with `KubernetesProviderAdapter` and
`KubernetesProviderClient`.

## Validation

Passed:

```bash
make validate-m1-e2e
```

Full repository validation should include:

```bash
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `PATCH`.

Reason: this slice adds integration profile coverage, manifests, and validation
tooling without changing public protobuf or database schemas.

## Next Work Boundary

Recommended next options:

- Implement a real-cluster provider execution profile based on
  `KubernetesProviderClient`.
- Start `M3-MODEL-A` after accepting the M1/M2/GPU/Runtime/Instance baseline.
