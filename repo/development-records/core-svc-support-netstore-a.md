# CORE-SVC-SUPPORT-NETSTORE-A — Core Services 支撑网络/存储/K8s Handler

> 批次类型：Feature batch
> 完成日期：2026-06-19
> 范围：仅 ANI Core，Tier1 local profile

## 背景

Sprint 12 目标是闭合 `api/openapi/v1.yaml` 已声明但网关尚未实现的 Core handler 缺口。本批覆盖 B2：网络路由、卷快照、文件系统挂载目标、K8s workloads，以及两个前置条件 422 行为。

## 完成内容

- 扩展 `ports.NetworkService`，新增 `CreateRoute` / `ListRoutes`；`runtime.LocalNetworkService` 提供带 `idempotency_key` 的 local route 状态与租户隔离。
- 扩展 `ports.StorageService`，新增 `CreateVolumeSnapshot` / `ListVolumeSnapshots` / `ListFilesystemMountTargets`；`runtime.LocalStorageService` 生成 local snapshot 与 mount target 元数据。
- 扩展 `ports.K8sClusterService`，新增 `ListWorkloads`；`runtime.NewLocalK8sClusterService()` 返回带筛选能力的 dev/local workload 摘要。
- Gateway 注册并实现 6 个 operationId：
  - `GET /api/v1/networks/routes`
  - `POST /api/v1/networks/routes`
  - `GET /api/v1/volumes/{volume_id}/snapshots`
  - `POST /api/v1/volumes/{volume_id}/snapshots`
  - `GET /api/v1/filesystems/{filesystem_id}/mount-targets`
  - `GET /api/v1/k8s-clusters/{cluster_id}/workloads`
- `searchVectorStore` 在向量库非 `ready` 时返回 `422 PRECONDITION_FAILED`。
- `createK8sCluster` 增加前置条件分支：local profile 遇到显式 real-provider 前置需求时返回 `422 PRECONDITION_FAILED`，不新建路由。
- `api/openapi/v1.yaml` 对齐 B2 list response：`items,total,next_cursor`；B2 新响应 schema 带 `dev_profile`。
- `validate_network_alpha_contract.py` / `validate_storage_alpha_contract.py` 增加 B2 path/schema/route/port token 覆盖；`validate_vector_alpha_contract.py` 覆盖 search 422。

## 验证

TDD 红测先行：

```bash
go test ./pkg/adapters/runtime ./services/ani-gateway/internal/router
```

红测阶段失败原因是缺少 B2 port 类型、local adapter 方法和 router 转换函数；实现后 targeted tests 通过。

完整门禁与 curl smoke 结果见本批提交记录和最终执行输出。

## 关键文件

- `api/openapi/v1.yaml`
- `pkg/ports/network_resources.go`
- `pkg/ports/storage_resources.go`
- `pkg/ports/k8s_clusters.go`
- `pkg/ports/errors.go`
- `pkg/adapters/runtime/network_service.go`
- `pkg/adapters/runtime/storage_service.go`
- `pkg/adapters/runtime/local_k8s_cluster_service.go`
- `pkg/adapters/runtime/k8s_cluster_proxy_forwarding_service.go`
- `pkg/adapters/runtime/vector_store_service.go`
- `services/ani-gateway/internal/router/network_resources.go`
- `services/ani-gateway/internal/router/storage_resources.go`
- `services/ani-gateway/internal/router/k8s_cluster_resources.go`
- `services/ani-gateway/internal/router/vector_store_resources.go`
- `scripts/validate_network_alpha_contract.py`
- `scripts/validate_storage_alpha_contract.py`
- `scripts/validate_vector_alpha_contract.py`

## 后续真实环境门禁关联

本批只完成 Tier1 local profile。Sprint 13 若推进真实 provider，必须沿用本批已建立的 port/handler 边界：

- 网络路由：从 `ports.NetworkService` 接 Kube-OVN route/provider adapter，并复用/扩展 Kube-OVN live gate。
- 卷快照与 mount target：从 `ports.StorageService` 接 Rook-Ceph/CSI/NFS provider adapter，并复用 Sprint 11 storage live evidence 入口。
- K8s workloads：从 `ports.K8sClusterService.ListWorkloads` 接 vCluster/Kubernetes API adapter，并复用 vCluster live gate。
- 向量检索前置：真实向量后端必须保持非 ready 返回 `PRECONDITION_FAILED`。

未执行 live gate 前，本批能力不得标记为 real-provider、runtime ready 或 production ready。
