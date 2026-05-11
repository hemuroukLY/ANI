# M1-INSTANCE-A: Instance Fabric

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.
- Required foundation before `模块 3：模型管理平台`.

## Scope

Added first-class instance object, lifecycle, network, and storage contracts:

- Extended `pkg/ports/workload_runtime.go`.
- Updated `pkg/adapters/runtime/not_configured.go`.
- Added `deploy/manifests/m1-instance-a/`.
- Added `deploy/helm/ani-platform/component-contracts/instance-fabric.yaml`.
- Added `deploy/helm/ani-platform/profiles/instance-foundation.yaml`.
- Added `scripts/validate_instance_contracts.py`.
- Added `make validate-instance-contracts` and wired it into `make validate-infra`.

The contract covers:

- VM / cloud host instance.
- Traditional container or Pod instance.
- GPU container or Pod instance.
- Inference instance.
- Notebook instance.
- AI Agent sandbox instance.
- Batch job instance.

## Network Plan

- `tenant_vpc`: tenant business traffic. VM and Pod instances may share this
  plane when they belong to the same business system and need direct
  connectivity.
- `foundation_mesh`: platform-controlled east-west service connectivity outside
  tenant VPCs.
- `storage`: storage backend access for object storage, PVC provisioners,
  model cache, and datasets.
- `management`: control-plane operations, health checks, logs, metrics, SSH/VNC
  proxy, and runtime management.
- `public_ingress`: explicit gateway or ingress exposure.

This avoids forcing platform, storage, and management dependencies into a
nested tenant VPC while still allowing VM/container business interoperability.

## Storage Plan

Supported attachments:

- `root_disk`
- `data_disk`
- `shared_pvc`
- `object_fuse`
- `ephemeral`

Runtime adapters must resolve storage class, bucket/PVC references, mount mode,
and retention policy before scheduling.

## Validation

Passed:

```bash
make validate-instance-contracts
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Extends capability contracts and deployment profiles before user-facing
  runtime implementation.

## Next Recommended Slice

- `M1-INSTANCE-B`: minimal runtime adapter/controller implementation for
  instance planning and lifecycle validation.
