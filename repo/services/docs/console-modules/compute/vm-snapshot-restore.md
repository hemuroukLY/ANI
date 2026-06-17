# VM 快照与恢复

## 页面定位

云主机 VM 详情下的快照创建、列表与回滚能力，复用统一实例 lifecycle 与卷快照模型。

父级：`vm-management.md`。

## 文档管理规则

- 本文是 VM 快照子模块主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 卷快照详文：`storage/block-storage-snapshot.md`

## Core 层要求

### 实例快照（lifecycle）

- `POST /api/v1/instances/{instance_id}/lifecycle`
- `operationId`: `applyInstanceLifecycle`
- `action` 枚举含 `snapshot`、`rollback`
- 请求必填：`action`、`idempotency_key`
- 请求可选：`snapshot_name`（snapshot）、`revision`（rollback）
- 成功：`200 + InstanceLifecycleResponse`
- 错误：`400`、`401`、`403`、`404`、`409`、`422`（YAML 已举例 `INSTANCE_STATE_INVALID`）

### 卷快照（关联存储）

<!-- ADDED-TO-YAML: /volumes/{volume_id}/snapshots (Phase 2) -->

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/volumes/{volume_id}/snapshots` | `listVolumeSnapshots` |
| POST | `/api/v1/volumes/{volume_id}/snapshots` | `createVolumeSnapshot` |

POST 成功：`202 + VolumeSnapshot`。

## 页面职责

- VM 详情「快照」区：触发 instance snapshot、展示 snapshots 列表（来自 Instance 内嵌或关联卷）
- 回滚操作二次确认；展示 lifecycle 异步反馈
- 卷级快照跳转块存储子模块

## 创建前置条件

| 场景 | 要求 | 未满足响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读权限 | `scope:instances:read` | `403` |
| 生命周期写权限 | `scope:instances:update` | `403` |
| instance snapshot | 实例状态允许 snapshot | `422`（YAML 已举例 `INSTANCE_STATE_INVALID`） |
| rollback | 目标快照/revision 存在 | 具体 `code` 待 Core 冻结 |
| 卷快照 POST | `volume_id` 有效 + `idempotency_key` | `404`/`400`/`409` |

## 操作可用性矩阵

| 操作 | stopped | running | deleting | 只读用户 | 运维 |
|---|---|---|---|---|---|
| 查看快照列表 | ✅ | ✅ | ❌ | ✅ | ✅ |
| instance snapshot | 视状态 | 视状态 | ❌ | ❌ | ✅ |
| rollback | 视状态 | 视状态 | ❌ | ❌ | ✅ |
| 卷快照创建 | — | — | — | ❌ | ✅（volumes:create） |

具体状态×动作矩阵见 `vm-management.md` §操作可用性矩阵。

## 接口冻结规则

### `POST /api/v1/instances/{instance_id}/lifecycle`（snapshot / rollback）

- 成功：`200 + InstanceLifecycleResponse`
- 错误：`400`、`401`、`403`、`404`、`409`、`422`

### `GET /api/v1/volumes/{volume_id}/snapshots`

- 成功：`200 + VolumeSnapshotListResponse`
- 错误：`401`、`403`、`404`

### `POST /api/v1/volumes/{volume_id}/snapshots`

- 成功：`202 + VolumeSnapshot`
- 错误：`400`、`401`、`403`、`404`、`409`

## 待补边界

- rollback 专用 `422` code 举例 — 待 Core YAML description 补充
- 快照删除 DELETE — **当前 YAML 未声明**（卷快照）
- 快照跨卷恢复 — 待存储域规划

## 验收标准

- [ ] snapshot/rollback 走 lifecycle，不自造 `/vm/snapshots` 路径
- [ ] 卷快照与块存储子模块口径一致
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
