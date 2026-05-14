# M1-GPU-A Heterogeneous GPU Discovery and Scheduling Contract

This directory defines the heterogeneous GPU contract required before model
deployment starts creating GPU workloads.

It covers:

- Multi-vendor GPU labels for NVIDIA, Huawei Ascend, and Hygon DCU.
- Same-vendor mixed-model classification.
- Kernel, driver, runtime, and device-plugin compatibility state.
- Resource-name mapping for vendor device plugins and HAMi/vGPU resources.
- Scheduling decision fields consumed by workload runtime adapters.

ANI services must depend on `ports.GPUInventory`, not vendor SDKs or direct
device-plugin APIs. Vendor-specific logic belongs in adapters and deployment
profiles.

Offline validation:

```bash
cd repo
make validate-gpu-contracts
```
