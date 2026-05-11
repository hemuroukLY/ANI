# M1-INFRA-F GPU Scheduling Preflight/E2E Hardening

This directory turns the M1-INFRA-E GPU scheduling contracts into an executable
cluster-side preflight profile.

It still does not install GPU Operator, HAMi, Volcano, or DCGM. Those components
remain provider-owned deploy profile dependencies so ANI does not hard-code
rapidly changing component implementations.

This slice provides:

- Least-privilege RBAC for GPU scheduling validation.
- A Kubernetes Job that validates GPU node labels, Volcano queues, namespaces,
  CRDs, and ANI GPU contract ConfigMaps.
- Optional strict checks for GPU allocatable resources and DCGM exporter service.
- A smoke workload template stored as ConfigMap data so operators must choose a
  customer-approved diagnostic image before running GPU workload tests.

Offline validation:

```bash
cd repo
make validate-infra
```

Cluster-side validation:

```bash
kubectl apply -f deploy/manifests/m1-infra-e
kubectl apply -f deploy/manifests/m1-infra-f
kubectl -n ani-system logs job/ani-gpu-e2e-preflight
```

Strict runtime options:

- Set `ANI_GPU_REQUIRE_HAMI_ALLOCATABLE=true` in the Job to require GPU
  allocatable resources to be visible on matching nodes.
- Set `ANI_GPU_REQUIRE_DCGM_SERVICE=true` in the Job to require a DCGM exporter
  service in `ani-gpu-system`.
