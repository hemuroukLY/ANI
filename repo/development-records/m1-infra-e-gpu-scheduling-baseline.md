# M1-INFRA-E: GPU Scheduling Baseline

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座` / `1.2 GPU 算力纳管`.

## Scope

Added GPU scheduling baseline artifacts:

- `deploy/manifests/m1-infra-e/00-gpu-namespaces.yaml`
- `deploy/manifests/m1-infra-e/10-gpu-node-label-contract.yaml`
- `deploy/manifests/m1-infra-e/20-volcano-queue-template.yaml`
- `deploy/manifests/m1-infra-e/30-hami-device-plugin-contract.yaml`
- `deploy/manifests/m1-infra-e/40-dcgm-observability-contract.yaml`
- `deploy/manifests/m1-infra-e/50-gpu-preflight-config.yaml`
- `deploy/manifests/m1-infra-e/README.md`
- `deploy/helm/ani-platform/component-contracts/gpu-scheduling.yaml`
- `deploy/helm/ani-platform/profiles/gpu-scheduling.yaml`

The slice keeps GPU Operator, HAMi, Volcano, and DCGM integration as deployment contracts instead of hard-coding provider behavior into business services.

## Validation

Passed:

```bash
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds new infrastructure deployment contracts and validation manifests without changing runtime service APIs.

## Next Recommended Slice

- `M1-INFRA-F`: GPU scheduling preflight/e2e hardening against an attached Kubernetes environment.
- Or begin `模块 3` model management after M1/M2 baseline acceptance.
