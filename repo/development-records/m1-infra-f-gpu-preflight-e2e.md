# M1-INFRA-F: GPU Preflight/E2E Hardening

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座` / `1.2 GPU 算力纳管`.

## Scope

Added executable GPU scheduling validation artifacts:

- `deploy/manifests/m1-infra-f/00-gpu-e2e-preflight-rbac.yaml`
- `deploy/manifests/m1-infra-f/10-gpu-e2e-preflight-config.yaml`
- `deploy/manifests/m1-infra-f/20-gpu-e2e-preflight-job.yaml`
- `deploy/manifests/m1-infra-f/30-gpu-smoke-workload-template.yaml`
- `deploy/manifests/m1-infra-f/README.md`
- `deploy/helm/ani-platform/profiles/gpu-scheduling-e2e.yaml`
- `scripts/validate_gpu_preflight_contract.py`

The cluster-side Job validates GPU node labels, Volcano queues, GPU namespaces,
required CRDs, and ANI GPU contract ConfigMaps. Optional strict checks can
require GPU allocatable resources and a DCGM exporter service once those
components are installed in a target environment.

The smoke workload is kept as a template so customer-approved diagnostic images
and site-specific resource keys are selected explicitly before execution.

## Validation

Passed:

```bash
make validate-gpu-preflight
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds new infrastructure validation manifests, Helm profile, and local contract
  tooling without changing service APIs.

## Next Recommended Slice

- `M3-MODEL-A`: model metadata and object-storage boundary design/implementation.
- Or `M1-INFRA-G`: optional offline install acceptance hardening.
