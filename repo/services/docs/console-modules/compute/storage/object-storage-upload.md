# 对象存储 — 上传与下载

## 页面定位

对象存储的桶管理、对象上传与下载，扩展 `object-storage.md`（原仅元数据 CRUD）。

## 文档管理规则

- 本文是对象上传/下载子模块主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 元数据 CRUD 见 `object-storage.md`
- `../governance/module-delivery-workflow.md` §2.9：创建元数据 ≠ 上传文件

## Core 层要求

<!-- ADDED-TO-YAML: Phase 2 2026-06-17 -->

### 桶

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| GET | `/api/v1/buckets` | `listStorageBuckets` | `scope:objects:read` |
| POST | `/api/v1/buckets` | `createStorageBucket` | `scope:objects:create` |

### 上传

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| POST | `/api/v1/objects/upload` | `uploadStorageObject` | `scope:objects:create` |

成功 `202 + AsyncTask`；请求 `ObjectUploadRequest`（含 `idempotency_key`）。

### 下载

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| GET | `/api/v1/objects/{object_id}/download` | `downloadStorageObject` | `scope:objects:read` |

成功 `200 + ObjectDownloadResponse`（预签名 URL）。

### 元数据（主模块已有）

`GET/POST/DELETE /api/v1/objects*` — 见 `object-storage.md`。

## 页面职责

- 桶列表/创建；对象上传表单；下载链接/预签名
- 上传进度关联任务中心 `AsyncTask`
- 桶策略 — 待补，不得自造路径

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 读/写权限 | 对应 objects scope | `403 FORBIDDEN` |
| 创建桶 POST | `name`、`idempotency_key` | `400`；冲突 `409` |
| 上传 POST | 桶须存在 | `422 PRECONDITION_FAILED`（upload 已声明） |
| 上传 POST | `idempotency_key` 等必填 | `400 BAD_REQUEST` |
| 下载 GET | `object_id` 有效 | `404 NOT_FOUND` |

## 操作可用性矩阵

| 操作 | 只读用户 | 存储管理员 |
|---|---|---|
| 桶列表/对象列表 | ✅ | ✅ |
| 创建桶 | ❌ | ✅ |
| 上传对象 | ❌ | ✅ |
| 获取下载链接 | ✅ | ✅ |
| 桶策略编辑 | ❌ | ❌（待 YAML） |

## 接口冻结规则

### `GET /api/v1/buckets`

- 成功：`200 + StorageBucketListResponse`
- 错误：`401`、`403`

### `POST /api/v1/buckets`

- 成功：`201 + StorageBucket`
- 错误：`400`、`401`、`403`、`409`

### `POST /api/v1/objects/upload`

- 成功：`202 + AsyncTask`
- 错误：`400`、`401`、`403`、`422`

### `GET /api/v1/objects/{object_id}/download`

- 成功：`200 + ObjectDownloadResponse`
- 错误：`401`、`403`、`404`

## 待补边界

- 桶策略独立 API — **TODO-YAML**
- 桶 DELETE / 对象 multipart 大文件 — **TODO-YAML**
- 上传 `422` 的 `code` 举例（如 `BUCKET_NOT_FOUND`）— 待 Core description 补充

## 验收标准

- [ ] 区分元数据 CRUD 与 upload/download
- [ ] 桶作为资源已 Phase 2 声明
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
