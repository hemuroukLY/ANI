# M1-INSTANCE-N: Kubernetes Provider Execution Profile

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added a real-cluster provider execution profile without introducing a concrete
Kubernetes SDK dependency yet:

- `pkg/adapters/runtime/kubernetes_provider_execution_profile_test.go`
- `deploy/manifests/m1-instance-n/`
- `scripts/validate_kubernetes_provider_execution.py`
- `Makefile` target `validate-kubernetes-provider-execution`

Updated:

- `CLAUDE.md`
- `AGENTS.md`
- `ANI-06-开发计划.md`
- `ANI-13-开源组件松耦合适配器架构.md`
- `development-records/README.md`

## Behavior

The profile wires `KubernetesProviderAdapter` into
`WorkloadInstanceOrchestrator` as:

- `WorkloadProviderDryRun`
- `WorkloadProviderApply`
- `WorkloadProviderStatusReader`

It verifies:

- `KubernetesProviderClient.ServerSideDryRun` is invoked before apply
- dry-run evidence preserves `dryRun=All`
- apply is explicitly enabled and receives audit id and permission proof
- apply returns resource refs
- observe returns tenant/instance/provider/resource-ref aligned status
- final status is reconciled and persisted through `WorkloadInstanceStore`

The slice deliberately avoids importing client-go, dynamic client,
controller-runtime, or KubeVirt client. Those libraries may only appear in a
future adapter-owned `KubernetesProviderClient` implementation.

## Validation

Passed:

```bash
make validate-kubernetes-provider-execution
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

Reason: this slice adds tests, manifests, and validation tooling for the
existing provider adapter boundary without changing public protobuf or database
schemas.

## Next Work Boundary

Recommended next options:

- Implement the adapter-owned real `KubernetesProviderClient`.
- Start `M3-MODEL-A` after accepting the M1/M2/GPU/Runtime/Instance baseline.
