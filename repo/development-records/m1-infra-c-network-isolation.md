# M1-INFRA-C Development Record

Date: 2026-05-11

Product plan mapping:
- `ANI-06 / 模块 1：基础设施底座`
- Slice: `M1-INFRA-C`
- Scope: KubeOVN tenant VPC and NetworkPolicy isolation templates.

Version impact:
- Current line: `v0.x` development.
- Release impact: `MINOR`.
- Reason: adds infrastructure network isolation templates and validation coverage.
- No formal version tag was created.

## Implemented

- Added tenant namespace template with ANI tenant labels.
- Added KubeOVN per-tenant VPC/Subnet template.
- Added tenant default-deny NetworkPolicy template.
- Added gateway-only ingress policy template for tenant services.
- Added tenant DNS egress policy template.
- Added AI Agent sandbox egress restriction templates.
- Extended `make validate-infra` to validate `m1-infra-c` manifests.

## Files Added

- `deploy/manifests/m1-infra-c/00-tenant-namespace-template.yaml`
- `deploy/manifests/m1-infra-c/10-kubeovn-tenant-vpc-template.yaml`
- `deploy/manifests/m1-infra-c/20-tenant-networkpolicy-template.yaml`
- `deploy/manifests/m1-infra-c/30-agent-sandbox-networkpolicy-template.yaml`
- `deploy/manifests/m1-infra-c/README.md`

## Files Updated

- `Makefile`
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

- KubeOVN CRDs are validated only as offline YAML in this slice; cluster-side CRD schema validation requires a KubeOVN-enabled cluster.
- CIDR allocation is still template-based. A tenant network allocator/controller is deferred.
- BGP peering and shared-services VPC routing are deferred.

## Next Boundary

Recommended next M1 slice:
- `M1-INFRA-D`: cluster-side KubeOVN validation profile or GPU scheduling baseline depending on available test infrastructure.
