# 文件存储 — 挂载目标

## 页面定位

文件系统详情下的**挂载目标**只读列表，展示当前文件系统已声明或可观测的挂载目标摘要。

本页属于 **Core / Filesystems** 子能力，父级 CRUD 见 `file-storage.md`。

## 文档管理规则

- 本文是挂载目标子页的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- Handler 实现契约：`tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` §TASK-CORE-010
- `prototypes/ani-services-prototype-console.html` 只保留摘要与详文入口

## Core 层要求

<!-- ADDED-TO-YAML: Phase 2 2026-06-17 -->

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/filesystems/{filesystem_id}/mount-targets` | `listFilesystemMountTargets` |

Query：`limit`、`cursor`（若 YAML 声明）。

Schema：`FilesystemMountTarget`、`FilesystemMountTargetListResponse`（以 OpenAPI 为准）。

RBAC：与文件系统读权限一致（见 YAML `x-ani-rbac-scope`）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 文件系统读权限 | 对应 scope 已授予 | `403 FORBIDDEN` |
| 文件系统存在 | 当前租户可见 | `404 NOT_FOUND` |

本页无 POST/PUT 创建挂载目标接口。**当前 YAML 未声明 `422`**。

## 操作可用性矩阵

| 操作 | 只读用户 | 存储管理员 | 说明 |
|---|---|---|---|
| 查看挂载目标列表 | 可用 | 可用 | `GET .../mount-targets` |
| 创建/删除挂载目标 | 不可用 | 不可用 | YAML 未声明写路径 |
| 修改访问策略 | 不可用 | 不可用 | 待补能力 |

## 页面职责

- 在文件系统详情 Tab 展示挂载目标列表（目标 ID、类型、状态、关联实例摘要）
- 提供跳转：关联实例详情（若 `instance_id` 字段存在且可见）
- 不冒充「活跃挂载关系」实时拓扑（该能力待补）

## 接口冻结规则

### `GET /api/v1/filesystems/{filesystem_id}/mount-targets`

- 成功：`200 + FilesystemMountTargetListResponse`
- 错误：`401`、`403`、`404`（filesystem 不存在）
- 分页：遵循 YAML 中 `limit`/`cursor` 声明
- 租户边界：仅返回当前租户可见挂载目标

## 待补边界

- 挂载/卸载动作 — **YAML 未声明**（禁止写 `POST .../mount` 为冻结契约）
- 活跃挂载关系实时查询 — 待 Core 规划
- 访问策略（NFS export / POSIX ACL）细化 — 待补
- 单挂载目标 metrics — YAML 未声明

## 相关模块

- `file-storage.md` — 文件系统 CRUD
- `storage-management.md` — 存储总览入口
- TASK：`TASK-CORE-010`

## 验收标准

- [ ] 路径与 Phase 2 Core YAML 一致
- [ ] 正文不把 handler stub 写成已实现
- [ ] 未自造 mount/unmount 写接口
- [ ] 与 `file-storage.md` 分层不冲突
