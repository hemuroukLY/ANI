# DEMO-INSTANCE-CONSOLE-A: Staged Instance Demo Console

## Status

Completed.

## Purpose

This slice adds a staged presentation experience for VM, container, and GPU
container instances. It is intentionally a demo/presentation layer that helps
humans see and validate the M1 instance capability earlier.

It must not change the M1 runtime, provider, lifecycle, or ops port contracts.

## Implementation

- `services/ani-gateway/internal/router/demo_instances.go`
- `services/ani-gateway/internal/router/demo_instances_test.go`
- `frontends/console/src/App.tsx`
- `frontends/console/src/demo/InstanceDemoPage.tsx`
- `frontends/console/src/main.tsx`
- `frontends/console/src/styles.css`

## API

- `GET /api/v1/demo/instances`
- `POST /api/v1/demo/instances`
- `GET /api/v1/demo/instances/{instance_id}`
- `POST /api/v1/demo/instances/{instance_id}/lifecycle`
- `GET /api/v1/demo/instances/{instance_id}/ops/{action}`

## Guardrails

- The demo API calls `WorkloadInstanceService`; it does not call Kubernetes,
  KubeVirt, or GPU scheduling APIs directly.
- The demo service uses a local in-memory store and local apply for presentation.
- Live-cluster execution remains gated by the existing provider configuration:
  `WORKLOAD_PROVIDER=kubernetes_rest`,
  `WORKLOAD_PROVIDER_APPLY_ENABLED=true`, and explicit Kubernetes credentials.
- GPU demo scheduling uses a fake inventory decision for presentation only.
  Production GPU discovery/scheduling remains behind the `GPUInventory` port.

## Validation

```bash
make validate-demo-instances
```

## Release Impact

`MINOR`: adds a presentation API and console page without changing M1 core
contracts.
