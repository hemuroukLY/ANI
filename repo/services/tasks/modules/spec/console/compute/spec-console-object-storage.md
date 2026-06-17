# SPEC: Console 对象存储

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-object-storage.md`
> Revised: 2026-06-15 | Scope: `Console / 算力与云资源 / 存储管理 / 对象存储`

## 1. Summary

### 1.1 What This SPEC Covers

本 SPEC 用于收口 `Console / 算力与云资源 / 存储管理 / 对象存储` 的技术边界、数据模型、接口冻结规则、错误处理方式和待补能力边界。

目标是把对象存储模块整理成一份可直接用于对齐 `Core v1.yaml` 的稳定材料，并明确区分：

- 当前 `Core` 已冻结能力
- 当前仍待补充到 `Core` 的规划能力

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/compute/prd-console-object-storage.md`
- User Stories covered: `US-001 ~ US-005`
- Functional Requirements covered: `FR-1 ~ FR-8`

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 |
|---|---|---|---|
| GET | `/api/v1/objects` | `listStorageObjects` | `200 + StorageObjectListResponse` |
| POST | `/api/v1/objects` | `createStorageObject` | `201 + StorageObject` |
| GET | `/api/v1/objects/{object_id}` | `getStorageObject` | `200 + StorageObject` |
| DELETE | `/api/v1/objects/{object_id}` | `deleteStorageObject` | `200 + StorageObject` |

### 2.3 Verified Schemas

- `StorageObject`
- `StorageObjectListResponse`
- `CreateStorageObjectRequest`
- `StorageResourceState`

## 3. Scope and Boundary

### 3.1 Page Scope

- 页面属于 `Console / 算力与云资源 / 存储管理`
- 页面直接消费 `Core /api/v1/*`
- 页面不定义新的 `Services /api/v1/svc/*` 对象存储资源契约
- `VM`、容器、`GPU 容器` 等页面中的对象资源摘要和跳转只消费本模块已冻结的对象元数据定义

### 3.2 Non-Frozen Capabilities

以下能力当前不能写成正式 Core 契约：

- 对象上传
- 对象下载
- 桶资源
- 桶策略

## 4. Data Model

### 4.1 Key Entity Definition

| 实体 | 关键字段 | 说明 |
|---|---|---|
| `StorageObject` | `id`、`tenant_id`、`bucket`、`key`、`size_bytes`、`content_type`、`state`、`reason`、`created_at`、`updated_at` | 当前只承接对象元数据 |

### 4.2 Naming Rules

- 页面展示字段允许使用 camelCase 作为 view model
- API 路径参数、query 参数、请求体和响应 JSON schema 统一以现有 `Core v1.yaml` 命名为准
- 不把旧 HTML 中的自定义对象操作字段直接写成已冻结接口参数

## 5. API Rules

### 5.1 Frozen Endpoints

| Method | Path | Response |
|---|---|---|
| GET | `/api/v1/objects` | `200 + StorageObjectListResponse` |
| POST | `/api/v1/objects` | `201 + StorageObject` |
| GET | `/api/v1/objects/{object_id}` | `200 + StorageObject` |
| DELETE | `/api/v1/objects/{object_id}` | `200 + StorageObject` |

### 5.2 Per-Endpoint Constraints

#### `GET /api/v1/objects`

- `operationId`: `listStorageObjects`
- `query params`: `limit?`、`cursor?`
- `success`: `200 + StorageObjectListResponse`
- `errors`: `401`、`403`

#### `POST /api/v1/objects`

- `operationId`: `createStorageObject`
- `requestBody.required`: `bucket`、`key`、`idempotency_key`
- `success`: `201 + StorageObject`
- `errors`: `400`、`401`、`403`

#### `GET /api/v1/objects/{object_id}`

- `operationId`: `getStorageObject`
- `success`: `200 + StorageObject`
- `errors`: `401`、`403`、`404`

#### `DELETE /api/v1/objects/{object_id}`

- `operationId`: `deleteStorageObject`
- `success`: `200 + StorageObject`
- `errors`: `401`、`403`、`404`

### 5.3 Non-Frozen Paths

以下路径当前不能写成已对齐 Core 的正式接口：

- `/api/v1/objects/{object_id}:upload`
- `/api/v1/objects/{object_id}:download`
- `/api/v1/object-buckets`
- 桶策略相关路径

## 6. Operation and Error Rules

### 6.1 Operation Availability

| 资源 | list/get | create | delete | action | 说明 |
|---|---|---|---|---|---|
| 对象元数据 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可以提示用户先确认是否存在潜在引用影响 |

### 6.2 Deletion Boundary

- 页面可以提示用户先确认是否存在潜在引用影响，但当前权威源未冻结“引用冲突校验”作为正式后端契约
- 当前 `DELETE /api/v1/objects/{object_id}` 显式冻结返回码为 `200`、`401`、`403`、`404`
- 因此本轮文档不能把 `409 CONFLICT` 写成已冻结对象存储删除契约

### 6.3 Shared Error Format

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-xxx"
}
```

重点错误：

| Error Code | HTTP Status | Condition |
|---|---|---|
| `BAD_REQUEST` | 400 | 参数格式错误、字段长度非法 |
| `UNAUTHORIZED` | 401 | 未登录或凭证失效 |
| `FORBIDDEN` | 403 | 无访问或操作权限 |
| `NOT_FOUND` | 404 | `object_id` 不存在 |

## 7. Examples

### 7.1 CreateStorageObjectRequest

```json
{
  "idempotency_key": "idem-object-001",
  "bucket": "tenant-assets",
  "key": "datasets/train/file-001.jsonl",
  "size_bytes": 2048,
  "content_type": "application/jsonl"
}
```

### 7.2 StorageObject

```json
{
  "id": "obj-12ab",
  "tenant_id": "t-demo",
  "bucket": "tenant-assets",
  "key": "datasets/train/file-001.jsonl",
  "size_bytes": 2048,
  "content_type": "application/jsonl",
  "state": "available",
  "reason": null,
  "created_at": "2026-06-13T10:00:00Z",
  "updated_at": "2026-06-13T10:00:00Z"
}
```

## 8. Risks and Prerequisites

### 8.1 Risks

| 风险 | 影响 | 处理方式 |
|---|---|---|
| 旧路径仍沿用 `/api/v1/storage/*` | 会造成前后端继续对错路径 | 统一改为现有 `/api/v1/objects*` 路径 |
| 旧文档把“上传文件”直接写成对象创建能力 | 会误导实现把未冻结动作写成既有契约 | 明确当前只承接对象元数据创建 |
| 把未冻结能力或未冻结返回码写进正式契约 | 后续 Core 对齐会失真 | 在正文和 SPEC 中单列待补能力 |

### 8.2 Core Alignment Prerequisites

- 后续若要把 `对象上传 / 下载` 写入正式接口，必须先扩充 `Core v1.yaml`
- 后续若要增加 `桶资源` 或 `桶策略`，必须先补充正式 schema 与路径
- 后续若要扩展更细粒度的对象引用关系，也必须先确认 `StorageObject` schema 是否需要增补
