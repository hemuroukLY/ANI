# M1-INFRA-A Development Record

Date: 2026-05-11

Product plan mapping:
- `ANI-06 / 模块 1：基础设施底座`
- Slice: `M1-INFRA-A`
- Scope: infrastructure-as-code baseline for ANI platform dependencies and network isolation.

Version impact:
- Current line: `v0.x` development.
- Release impact: `MINOR`.
- Reason: adds deployable infrastructure manifests, Helm chart contract, and validation target.
- No formal version tag was created.

## Implemented

- Added `deploy/manifests/m1-infra-a` with:
  - `ani-system` namespace.
  - Platform dependency config placeholders.
  - Bootstrap secret placeholders for PostgreSQL, NATS, Redis, and MinIO endpoints.
  - Default-deny NetworkPolicy.
  - Initial control-plane ingress/egress NetworkPolicy seeds.
  - ServiceAccounts for gateway, auth, task, and model services.
- Added `deploy/helm/ani-platform` umbrella chart metadata and values contract.
- Added offline Kubernetes YAML validation script:
  - `scripts/validate_k8s_yaml.py`
- Added `make validate-infra`.

## Files Added

- `deploy/manifests/m1-infra-a/00-namespace.yaml`
- `deploy/manifests/m1-infra-a/10-platform-config.yaml`
- `deploy/manifests/m1-infra-a/20-networkpolicy.yaml`
- `deploy/manifests/m1-infra-a/30-serviceaccounts.yaml`
- `deploy/manifests/m1-infra-a/README.md`
- `deploy/helm/ani-platform/Chart.yaml`
- `deploy/helm/ani-platform/values.yaml`
- `deploy/helm/ani-platform/README.md`
- `scripts/validate_k8s_yaml.py`

## Files Updated

- `Makefile`
- `development-records/README.md`

## Validation

Commands run:

```bash
cd repo
make validate-infra
make gen-proto
make test
make build
```

Result:
- `make validate-infra`: passed.
- `make gen-proto`: passed.
- `make test`: passed.
- `make build`: passed.

## Remaining Risks

- Helm rendering templates are not implemented yet; the chart currently defines the contract and values surface only.
- Component deployments for CloudNativePG, NATS JetStream, Redis, MinIO, Milvus, and Harbor are placeholders and must be filled in later M1 slices.
- KubeOVN VPC/Subnet CRs, GPU Operator, HAMi, Volcano, and observability dashboards are deferred.
- Offline image mirror lockfiles and installer packaging are deferred.

## Next Boundary

Recommended next M1 slice:
- `M1-INFRA-B`: component Helm dependencies and install profiles for PostgreSQL, NATS JetStream, Redis, MinIO, Milvus, and Harbor.
