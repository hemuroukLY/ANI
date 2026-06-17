# 块存储 — 卷快照

## 页面定位

块存储卷的快照列表与创建，作为 `block-storage.md` 子能力。

## 文档管理规则

- 本文是卷快照子模块主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- VM 实例快照见 `vm-snapshot-restore.md`

## Core 层要求

<!-- ADDED-TO-YAML: Phase 2 2026-06-17 -->

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| GET | `/api/v1/volumes/{volume_id}/snapshots` | `listVolumeSnapshots` | `scope:volumes:read` |
| POST | `/api/v1/volumes/{volume_id}/snapshots` | `createVolumeSnapshot` | `scope:volumes:create` |

Schema：`VolumeSnapshot`、`VolumeSnapshotListResponse`、`CreateVolumeSnapshotRequest`。

POST 成功：`202 + VolumeSnapshot`（异步接受，非同步完成）。

POST 请求必须带 `idempotency_key`（以 `CreateVolumeSnapshotRequest` 为准）。

## 页面职责

- 卷详情下快照 Tab：列表、创建、状态
- 创建后展示 `202` 响应体或关联 `AsyncTask`（若 handler 返回 task_id）
- 删除快照 — **当前 YAML 未声明 DELETE**

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| `volume_id` | 存在且租户可见 | `404 NOT_FOUND` |
| 读权限 | `scope:volumes:read` | `403 FORBIDDEN` |
| 写权限 | `scope:volumes:create` | `403 FORBIDDEN` |
| 请求体 | `idempotency_key` 等必填 | `400 BAD_REQUEST` |
| 卷状态允许快照 | 产品语义 | **当前 YAML 未声明 `422`** |
| 名称/并发冲突 | — | `409 CONFLICT`（create 已声明） |

## 操作可用性矩阵

| 操作 | 只读用户 | 存储管理员 |
|---|---|---|
| 列表 | ✅ | ✅ |
| 创建 | ❌ | ✅ |
| 删除 | ❌ | ❌（待 YAML） |
| 从快照恢复卷 | ❌ | ❌（待 YAML） |

## 接口冻结规则

### `GET /api/v1/volumes/{volume_id}/snapshots`

- 成功：`200 + VolumeSnapshotListResponse`
- 错误：`401`、`403`、`404`

### `POST /api/v1/volumes/{volume_id}/snapshots`

- 成功：`202 + VolumeSnapshot`
- 错误：`400`、`401`、`403`、`404`、`409`

## 待补边界

- 快照 DELETE — **TODO-YAML**
- 从快照创建新卷 — **TODO-YAML**
- 快照创建完成态轮询 — 可复用 `GET /tasks/{task_id}` 若 handler 关联

## 验收标准

- [ ] 202 语义与 YAML 一致
- [ ] 未声明 DELETE 不写冻结事实
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
