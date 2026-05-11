# M1-E2E-B: Real Provider Integration Regression Profile

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added a unified regression profile for the Kubernetes REST provider path:

- `pkg/adapters/runtime/m1_real_provider_profile_test.go`
- `deploy/manifests/m1-e2e-b/`
- `scripts/validate_m1_real_provider_profile.py`
- `Makefile` target `validate-m1-real-provider-profile`

Updated:

- `CLAUDE.md`
- `AGENTS.md`
- `ANI-06-开发计划.md`
- `ANI-13-开源组件松耦合适配器架构.md`
- `development-records/README.md`

## Behavior

The profile uses a fake HTTP transport, not a live Kubernetes cluster, while
keeping the real adapter chain:

- `KubernetesRESTClient`
- `KubernetesProviderAdapter`
- `KubernetesLifecycleExecutor`
- `KubernetesInstanceOps`
- `WorkloadInstanceService`

It verifies one unified instance flow:

- create with server-side dry-run and server-side apply
- observe provider status
- lifecycle start/stop via provider scale operation
- logs through pod log subresource
- exec through pod exec subresource

## Validation

Passed:

```bash
make validate-m1-real-provider-profile
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

Reason: this adds regression coverage and validation tooling without changing
public protobuf or database schemas.

## Next Work Boundary

Recommended next options:

- Start `M3-MODEL-A`.
- Add an optional live-cluster smoke profile for Kubernetes REST provider.
