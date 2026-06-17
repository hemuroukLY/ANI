# K8s 集群 — 工作负载摘要

## 页面定位

K8s 集群详情下的**工作负载摘要**只读页（非完整 Helm/应用市场）。

## 文档管理规则

- Handler 契约：`tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` §TASK-CORE-012
- 集群 CRUD：`k8s-cluster-management.md`

## Core 层要求

<!-- ADDED-TO-YAML: Phase 2 2026-06-17 -->

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/k8s-clusters/{cluster_id}/workloads` | `listK8sClusterWorkloads` |

RBAC：与集群读权限一致（见 YAML `x-ani-rbac-scope`）。

## 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 集群存在 | 当前租户可访问 | `404` |
| 用户认证 | 已登录 | `401` |

## 操作可用性矩阵

| 操作 | 只读 | 集群管理员 |
|---|---|---|
| 查看工作负载摘要 | ✅ | ✅ |
| 修改工作负载 | ❌ | ❌（无独立写 API） |

## 接口冻结规则

### `GET /api/v1/k8s-clusters/{cluster_id}/workloads`

- 成功：`200 + K8sClusterWorkloadListResponse`（以 OpenAPI schema 为准）
- 错误：`401`、`403`、`404`（cluster 不存在）
- 分页：遵循 YAML query（若声明）
- 租户边界：仅当前租户可见集群

## 待补边界

- Pod 事件、Helm、应用市场 — **YAML 未声明**
- 创建集群 `422` 验通 — `POST /k8s-clusters` PreconditionFailed（同 TASK-CORE-012）

## 相关模块

- `k8s-cluster-management.md`
- TASK：`TASK-CORE-012`

## 验收标准

- [ ] workloads 路径与 Phase 2 YAML 一致
- [ ] 不把 Helm/事件页写成冻结能力
