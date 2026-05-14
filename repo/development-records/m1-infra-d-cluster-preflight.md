# M1-INFRA-D: Cluster Preflight Validation Profile

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 1：基础设施底座`.

## Scope

Added cluster-side preflight validation artifacts:

- `deploy/manifests/m1-infra-d/00-preflight-rbac.yaml`
- `deploy/manifests/m1-infra-d/10-preflight-config.yaml`
- `deploy/manifests/m1-infra-d/20-preflight-job.yaml`
- `deploy/manifests/m1-infra-d/README.md`
- `deploy/helm/ani-platform/profiles/cluster-validation.yaml`

The preflight job checks Kubernetes API access, required namespaces, KubeOVN CRDs, bootstrap secrets, optional StorageClass names, and NetworkPolicy API availability.

## Validation

Passed:

```bash
make validate-infra
make test
make build
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds new deployment validation manifests/profile without changing existing runtime APIs.

## Next Recommended Slice

- `M1-INFRA-E`: GPU scheduling baseline or expanded offline cluster acceptance checks.
