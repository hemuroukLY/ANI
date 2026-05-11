# M1-GPU-A: Heterogeneous GPU Discovery and Scheduling Contract

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座` / `1.2 GPU 算力纳管`.

## Scope

Added heterogeneous GPU abstraction and contracts:

- `pkg/ports/gpu_inventory.go`
- `pkg/adapters/gpu/not_configured.go`
- `deploy/manifests/m1-gpu-a/`
- `deploy/helm/ani-platform/component-contracts/gpu-inventory.yaml`
- `scripts/validate_gpu_contracts.py`

The contract covers NVIDIA, Huawei Ascend, and Hygon DCU environments, including
same-vendor mixed GPU models, kernel/driver/runtime compatibility, device-plugin
resource names, runtime class, Volcano queue, and scheduling decision fields.

Business services must depend on `ports.GPUInventory`; vendor-specific logic
belongs in adapters, deploy profiles, preflight jobs, or bounded runtime modules.

## Validation

Passed:

```bash
make validate-gpu-contracts
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds a new capability port and infrastructure contracts before model runtime
  implementation.

## Next Recommended Slice

- `M1-RUNTIME-A`: workload runtime and user-visible instance abstraction.
