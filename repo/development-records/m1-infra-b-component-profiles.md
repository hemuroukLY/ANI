# M1-INFRA-B Development Record

Date: 2026-05-11

Product plan mapping:
- `ANI-06 / 模块 1：基础设施底座`
- Slice: `M1-INFRA-B`
- Scope: component install profiles and infrastructure dependency contracts.

Version impact:
- Current line: `v0.x` development.
- Release impact: `MINOR`.
- Reason: adds infrastructure install profiles and component contracts for required platform dependencies.
- No formal version tag was created.

## Implemented

- Extended `deploy/helm/ani-platform/values.yaml` with profile and component fields.
- Added install profile overlays:
  - `dev`
  - `attach-k8s`
  - `offline`
- Added component contracts for:
  - PostgreSQL / CloudNativePG
  - NATS JetStream
  - Redis
  - MinIO
  - Milvus
  - Harbor
- Added generic offline YAML validation:
  - `scripts/validate_yaml.py`
- Updated `make validate-infra` to validate both Kubernetes manifests and Helm/profile YAML.

## Files Added

- `deploy/helm/ani-platform/profiles/dev.yaml`
- `deploy/helm/ani-platform/profiles/attach-k8s.yaml`
- `deploy/helm/ani-platform/profiles/offline.yaml`
- `deploy/helm/ani-platform/component-contracts/postgresql.yaml`
- `deploy/helm/ani-platform/component-contracts/nats.yaml`
- `deploy/helm/ani-platform/component-contracts/redis.yaml`
- `deploy/helm/ani-platform/component-contracts/minio.yaml`
- `deploy/helm/ani-platform/component-contracts/milvus.yaml`
- `deploy/helm/ani-platform/component-contracts/harbor.yaml`
- `scripts/validate_yaml.py`

## Files Updated

- `Makefile`
- `deploy/helm/ani-platform/Chart.yaml`
- `deploy/helm/ani-platform/values.yaml`
- `deploy/helm/ani-platform/README.md`
- `development-records/README.md`

## Validation

Commands run:

```bash
cd repo
make gen-proto
make validate-infra
make test
make build
```

Result:
- `make gen-proto`: passed.
- `make validate-infra`: passed.
- `make test`: passed.
- `make build`: passed.

## Remaining Risks

- Public chart dependencies are not downloaded or rendered yet; versions are compatibility targets until a lockfile is added.
- Offline image and chart lockfiles are still deferred.
- KubeOVN, GPU Operator, HAMi, Volcano, and observability manifests are not in this slice.
- Cluster-level validation still requires a real or ephemeral K8s cluster.

## Next Boundary

Continue with `M2.2-AUTH-B`: wire ANI Gateway Auth/RBAC middleware to auth-service over gRPC.
