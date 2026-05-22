# M1-INFRA-A Kubernetes Baseline

This directory contains the first code-generated infrastructure baseline for:

- `ANI-06 / 模块 1：基础设施底座`
- Slice: `M1-INFRA-A`

The manifests are intentionally minimal and are safe to validate without a live cluster:

```bash
cd repo
make validate-infra
```

Included baseline:

- `ani-system` namespace.
- Platform config placeholders for required infrastructure dependencies.
- Default-deny NetworkPolicy.
- Narrow ingress/egress NetworkPolicy seeds for ANI control-plane services.
- ServiceAccounts for the initial Go control-plane services.

Deferred to later M1 slices:

- Native Kubernetes bootstrap automation.
- KubeOVN VPC/Subnet CRs.
- CloudNativePG, MinIO, Milvus, NATS, Redis, Harbor Helm dependencies.
- GPU Operator, HAMi, Volcano, Prometheus/Grafana dashboards.
- Offline installer packaging and image mirror lockfiles.
