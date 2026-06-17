# SPEC: Console filesystem-mount-targets

> Technical specification — Core handler 契约见 `tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` §TASK-CORE-010  
> Source: `tasks/modules/prd/console/compute/prd-console-filesystem-mount-targets.md`  
> Revised: 2026-06-17

## 1. Summary

文件系统详情下的**挂载目标**只读列表；展示当前文件系统已声明或可观测的挂载目标摘要。属于 **Core / Filesystems** 子能力，父级 CRUD 见 `file-storage.md`。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/filesystems/{filesystem_id}/mount-targets` | `listFilesystemMountTargets` | `200 + FilesystemMountTargetListResponse` | 与文件系统读权限一致（见 YAML `x-ani-rbac-scope`） |

Query：`limit`、`cursor`（若 YAML 声明）。

### 2.3 Verified Schemas

- `FilesystemMountTarget`（以 OpenAPI 为准）
- `FilesystemMountTargetListResponse`（以 OpenAPI 为准）

## 3. Page Scope

- 在文件系统详情 Tab 展示挂载目标列表（目标 ID、类型、状态、关联实例摘要）
- 提供跳转：关联实例详情（若 `instance_id` 字段存在且可见）
- 不冒充「活跃挂载关系」实时拓扑（该能力待补）

## 4. Non-Goals

- 挂载/卸载动作 — **YAML 未声明**（禁止写 `POST .../mount` 为冻结契约）
- 创建/删除挂载目标 — YAML 未声明写路径
- 活跃挂载关系实时查询 — 待 Core 规划
- 访问策略（NFS export / POSIX ACL）细化 — 待补
- 单挂载目标 metrics — YAML 未声明
- **当前 YAML 未声明 `422`**

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 文件系统读权限 | 对应 scope 已授予 | `403 FORBIDDEN` |
| 文件系统存在 | 当前租户可见 | `404 NOT_FOUND` |

## 6. 操作可用性矩阵

| 操作 | 只读用户 | 存储管理员 | 说明 |
|---|---|---|---|
| 查看挂载目标列表 | ✅ | ✅ | `GET .../mount-targets` |
| 创建/删除挂载目标 | ❌ | ❌ | YAML 未声明写路径 |
| 修改访问策略 | ❌ | ❌ | 待补能力 |

## 7. 主维护源

- `docs/console-modules/compute/storage/filesystem-mount-targets.md`
- 父模块：`file-storage.md`、`storage-management.md`
- TASK：`TASK-CORE-010`

## 8. Handler 验收（Core 团队）

```bash
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/filesystems/{filesystem_id}/mount-targets"
```

- 分页：遵循 YAML 中 `limit`/`cursor` 声明
- 未自造 mount/unmount 写接口
- OpenAPI 已声明 ≠ handler 已实现
