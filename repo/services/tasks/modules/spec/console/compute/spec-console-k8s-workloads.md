# SPEC: Console k8s-workloads

> Technical specification — Core handler 契约见 `tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` §TASK-CORE-012  
> Source: `tasks/modules/prd/console/compute/prd-console-k8s-workloads.md`  
> Revised: 2026-06-17

## 1. Summary

K8s 集群详情下的**工作负载摘要**只读页（非完整 Helm/应用市场）。集群 CRUD 见 `k8s-cluster-management.md`。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/k8s-clusters/{cluster_id}/workloads` | `listK8sClusterWorkloads` | `200 + K8sClusterWorkloadListResponse` | 与集群读权限一致（见 YAML `x-ani-rbac-scope`） |

### 2.3 Verified Schemas

- `K8sClusterWorkloadListResponse`（以 OpenAPI schema 为准）

## 3. Page Scope

- 集群详情下展示工作负载摘要列表
- 分页遵循 YAML query（若声明）
- 租户边界：仅当前租户可见集群

## 4. Non-Goals

- Pod 事件、Helm、应用市场 — **YAML 未声明**
- 修改工作负载 — 无独立写 API
- 创建集群 `422` 验通 — `POST /k8s-clusters` PreconditionFailed（同 TASK-CORE-012，非本页职责）

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户认证 | 已登录 | `401` |
| 集群存在 | 当前租户可访问 | `404` |

## 6. 操作可用性矩阵

| 操作 | 只读 | 集群管理员 |
|---|---|---|
| 查看工作负载摘要 | ✅ | ✅ |
| 修改工作负载 | ❌ | ❌（无独立写 API） |

## 7. 主维护源

- `docs/console-modules/compute/k8s-workloads.md`
- 父模块：`k8s-cluster-management.md`
- TASK：`TASK-CORE-012`

## 8. Handler 验收（Core 团队）

```bash
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/k8s-clusters/{cluster_id}/workloads"
```

- 不把 Helm/事件页写成冻结能力
- OpenAPI 已声明 ≠ handler 已实现
