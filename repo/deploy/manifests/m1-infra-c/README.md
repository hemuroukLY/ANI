# M1-INFRA-C Network Isolation Baseline

This directory contains template manifests for:

- KubeOVN per-tenant VPC/Subnet.
- Per-tenant namespace labeling.
- Default-deny tenant NetworkPolicy.
- Gateway-only ingress to tenant services.
- AI Agent sandbox egress restrictions.

These files are templates. Replace:

- `ani-tenant-template`
- `REPLACE_WITH_TENANT_UUID`
- tenant CIDR blocks

before applying to a cluster.

Validation:

```bash
cd repo
make validate-infra
```

Deferred:

- BGP peer configuration.
- KubeOVN gateway peering between tenant VPCs and shared platform services.
- Runtime controller that materializes these templates per tenant.
