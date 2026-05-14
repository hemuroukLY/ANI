# M1-INFRA-D Cluster Preflight Profile

This directory contains cluster-side validation manifests for attach-k8s and
offline deployments.

The preflight job checks:

- Kubernetes API reachability.
- Required namespaces.
- KubeOVN CRDs for `Vpc` and `Subnet`.
- Required bootstrap secrets.
- Optional StorageClass name after replacing `REPLACE_WITH_STORAGE_CLASS`.
- NetworkPolicy API availability.

These manifests are intentionally separate from the core install baseline. They
are operational validation tools, not long-running platform components.

Validation:

```bash
cd repo
make validate-infra
```

Apply manually in a target cluster after replacing placeholders:

```bash
kubectl apply -f deploy/manifests/m1-infra-d
kubectl -n ani-system logs job/ani-cluster-preflight
```
