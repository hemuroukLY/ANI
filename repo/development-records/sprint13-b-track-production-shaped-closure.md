# SPRINT13-B-TRACK-PRODUCTION-SHAPED-CLOSURE - S01-S04 production-shaped closure

> 记录类型：Sprint 13 B-track production-shaped closure
> 日期：2026-06-20
> 范围：仅 ANI Core S01-S04 production-shaped 代码、部署契约与门禁闭环；不改 Services，不推远端
> 状态：**production-shaped acceptance standard ready; historical S01-S04 lab evidence remains pending until rerun with production-shaped live mode**。

## 目标

把 S01-S04 从“只防止误标 production ready”推进到可生产验收的代码和门禁形态：

- Gateway Kubernetes REST client 支持标准 in-cluster ServiceAccount token、CA bundle 与 `KUBERNETES_SERVICE_HOST/PORT`，生产部署不再依赖本机 kubeconfig 或 kubectl proxy。
- S01 network route 元数据从 in-memory-only 提升为 metadata store 持久化路径，provider-backed route apply/observe 后可落库。
- 新增 production-shaped Gateway RBAC/profile，固定最小 ServiceAccount、ClusterRole、ClusterRoleBinding 与 S01-S07 B 轨通过标准。
- S02/S03/S04 live gate 增加 `--production-shaped` 模式；启用后拒绝 localhost/127、本机 proxy、port-forward/dev gateway 证据，并写入 `production_shape.status=passed` + `proof_items`。
- `validate-sprint13-b-track-production-shape` 从边界 guard 升级为生产形态契约门禁：passed evidence 必须有正向 proof_items，且不得使用 lab/local/dev transport。

## 本次代码闭环

| 切片 | 修复内容 | 生产验收影响 |
|---|---|---|
| S01 网络路由 Kube-OVN | `NetworkResourceStore.UpsertRoute`、`MetadataNetworkStore.UpsertRoute`、`LocalNetworkService.CreateRoute` 持久化 pending/success/failed route；新增 `network_routes` migration | route provider 不再只有内存态，具备持久 route metadata reconciliation 基础 |
| S01/S03/S04 Kubernetes provider | `KubernetesRESTClient` 支持 in-cluster ServiceAccount token/CA；Gateway network/storage/gpu runtime 透传 service host/port/token/CA file | Gateway 可用正式 ServiceAccount/RBAC 访问 Kubernetes API |
| S02 K8s workloads vCluster | `validate_vcluster_live_gate.py --production-shaped` 拒绝 `--proxy-server` 与本地 Gateway，要求非本地 metadata target server | 后续 S02 passed evidence 必须证明 per-cluster metadata target + TLS/token |
| S03 storage Rook-Ceph | `validate_storage_live_gate.py --production-shaped` 拒绝本地 Gateway，并写入 S03 production proof items | 后续 S03 passed evidence 必须经 production gateway / in-cluster RBAC 路径 |
| S04 GPU inventory/DCGM | `validate_gpu_inventory_live_gate.py --production-shaped` 拒绝 `--kubernetes-nodes-url` 与本地 DCGM/Prometheus URL | 后续 S04 passed evidence 必须经 in-cluster Kubernetes API 与集群 Service/Prometheus |
| S05-S07 标准 | `sprint13-production-shaped-gateway-profile.yaml` 写入 S05/S06/S07 `slice_proof_items` | 后续 B 轨沿用同一 production-shaped passed 标准 |

## 新增/更新门禁

```bash
make validate-sprint13-b-track-production-shape
```

该门禁现在同时检查：

- S01-S04 historical evidence 仍必须显式保留 `production_shape.status=pending`，不能误标 production ready。
- 若未来某切片把 `production_shape.status` 改为 `passed`，必须：
  - `transport_profile` 不含 lab/local/dev gateway/kubectl proxy/port-forward；
  - `missing_items` 为空；
  - `proof_items` 包含该切片要求的生产证明项。
- `deploy/real-k8s-lab/sprint13-production-shaped-gateway-profile.yaml` 必须覆盖 S01-S07 的 proof 标准。
- `deploy/real-k8s-lab/sprint13-production-shaped-gateway-rbac.yaml` 必须包含 Gateway ServiceAccount、ClusterRole、ClusterRoleBinding 和最小 Kubernetes/Kube-OVN/CSI/GPU 资源权限，且不得授予 wildcard resources。

## 重要边界

本批次没有把旧 S01-S04 lab evidence 改写为 production-shaped passed，因为尚未在正式 Gateway + in-cluster ServiceAccount/RBAC + metadata target / cluster Service 路径重新执行真实写操作。

S01-S04 现在达到的是 **production-shaped acceptance standard ready**：代码、部署契约、迁移和门禁已经具备生产验收形态；要把单个切片标 `production_shape.status=passed`，必须重新跑对应 `--production-shaped` live gate 并产出新的非敏感 evidence JSON。

## 验证

本批次提交前执行完整 Sprint 13 基线与 production-shaped 门禁；输出以提交记录为准。
