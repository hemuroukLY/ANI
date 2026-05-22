# M1-K8S-A · K8s 集群 Core API dev profile

- 完成日期：2026-05-21；2026-05-22 补齐 kubeconfig local profile
- 状态：✅

实现了 `/api/v1/k8s-clusters` 的 create/get/list/delete 路由、local service、租户隔离与幂等创建；2026-05-22 继续补齐 `GET /api/v1/k8s-clusters/{cluster_id}/kubeconfig` 的 local dev profile，返回模拟 vCluster kubeconfig、token、server、namespace、过期时间和 `dev_profile` 标记；并补充了 OpenAPI 合同、SDK 生成和网关单元测试。

## 真实边界

本批次仍是 Core dev/local profile 切片，不包含：

- 真实 vCluster provider 创建/升级/删除。
- 原生 K8s API 透明代理。
- 节点池管理。
- 后台 reconcile controller。
