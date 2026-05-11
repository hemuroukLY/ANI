# M1-INSTANCE-O: Kubernetes REST Client

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added the first adapter-owned `KubernetesProviderClient` implementation:

- `pkg/adapters/runtime/kubernetes_rest_client.go`
- `pkg/adapters/runtime/kubernetes_rest_client_test.go`
- `deploy/manifests/m1-instance-o/`
- `scripts/validate_kubernetes_rest_client.py`
- `Makefile` target `validate-kubernetes-rest-client`

Updated:

- `CLAUDE.md`
- `AGENTS.md`
- `ANI-06-开发计划.md`
- `ANI-13-开源组件松耦合适配器架构.md`
- `development-records/README.md`

## Behavior

`KubernetesRESTClient` implements `KubernetesProviderClient` with standard
library HTTP:

- `ServerSideDryRun` issues `POST ...?dryRun=All`.
- `Apply` issues server-side apply `PATCH` with `fieldManager` and
  `force=true`.
- `Observe` issues `GET` and returns standard `WorkloadProviderObservation`.

Supported resources:

- Kubernetes `Deployment`
- Kubernetes `Job`
- KubeVirt `VirtualMachine`

The implementation remains adapter-owned and does not import client-go,
controller-runtime, or KubeVirt client libraries.

## Validation

Passed:

```bash
make validate-kubernetes-rest-client
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

Reason: this adds a new runtime adapter implementation, while public protobuf
and database schemas remain unchanged.

## Next Work Boundary

Recommended next options:

- Wire `KubernetesRESTClient` into bootstrap/configuration.
- Extend real provider lifecycle and visual ops execution.
- Start `M3-MODEL-A` after accepting the M1/M2/GPU/Runtime/Instance baseline.
