# SPEC: Console 存储管理

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-storage-management.md`
> Revised: 2026-06-15 | Scope: `Console / 算力与云资源 / 存储管理`

## 1. Summary

### 1.1 What This SPEC Covers

本 SPEC 用于收口 `Console / 算力与云资源 / 存储管理` 的技术边界、数据模型、接口冻结规则、错误处理方式和待补能力边界。

目标是把存储管理模块整理成一份可直接用于对齐 `Core v1.yaml` 的稳定材料，并明确区分：

- 当前 `Core` 已冻结能力
- 当前仍待补充到 `Core` 的规划能力

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/compute/prd-console-storage-management.md`
- User Stories covered: `US-001 ~ US-006`
- Functional Requirements covered: `FR-1 ~ FR-10`

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| 资源 | 已冻结路径 |
|---|---|
| 块存储 | `/api/v1/volumes`、`/api/v1/volumes/{volume_id}` |
| 文件存储 | `/api/v1/filesystems`、`/api/v1/filesystems/{filesystem_id}` |
| 对象元数据 | `/api/v1/objects`、`/api/v1/objects/{object_id}` |
| 向量存储 | `/api/v1/vector-stores`、`/api/v1/vector-stores/{vector_store_id}`、`/api/v1/vector-stores/{vector_store_id}/search` |

### 2.3 Verified Schemas

- `StorageVolume`
- `StorageVolumeListResponse`
- `CreateStorageVolumeRequest`
- `StorageFilesystem`
- `StorageFilesystemListResponse`
- `CreateStorageFilesystemRequest`
- `StorageObject`
- `StorageObjectListResponse`
- `CreateStorageObjectRequest`
- `VectorStore`
- `VectorStoreListResponse`
- `CreateVectorStoreRequest`
- `VectorStoreSearchRequest`
- `VectorStoreSearchHit`
- `VectorStoreSearchResponse`

## 3. Scope and Boundary

### 3.1 Page Scope

- 页面属于 `Console / 算力与云资源`
- 页面直接消费 `Core /api/v1/*`
- 页面不定义新的 `Services /api/v1/svc/*` 存储资源契约
- `VM`、容器、`GPU 容器` 页面中的 `卷 / 文件 / 对象 / 向量存储` 只消费本模块已冻结资源摘要和跳转入口

### 3.2 Non-Frozen Capabilities

以下能力当前不能写成正式 Core 契约：

- 卷快照
- 卷挂载 / 卸载
- 对象上传 / 下载
- 桶资源
- 桶策略
- 文件系统挂载目标

## 4. Data Model

### 4.1 Key Entity Definitions

| 实体 | 关键字段 | 说明 |
|---|---|---|
| `StorageVolume` | `id`、`name`、`size_gib`、`storage_class`、`state`、`reason`、`created_at`、`updated_at` | 块存储卷 |
| `StorageFilesystem` | `id`、`name`、`protocol`、`size_gib`、`endpoint`、`state`、`reason`、`created_at`、`updated_at` | 文件存储 |
| `StorageObject` | `id`、`bucket`、`key`、`size_bytes`、`content_type`、`state`、`reason`、`created_at`、`updated_at` | 当前只承接对象元数据 |
| `VectorStore` | `id`、`name`、`dimension`、`metric`、`state`、`reason`、`created_at`、`updated_at` | 向量存储 |

### 4.2 Naming Rules

- 页面展示字段允许使用 camelCase 作为 view model
- API 路径参数、query 参数、请求体和响应 JSON schema 统一以现有 `Core v1.yaml` 命名为准
- 不能把旧 HTML 中的自定义动作字段写成已冻结接口参数

## 5. API Rules

### 5.1 Frozen Endpoints

| Method | Path | Response |
|---|---|---|
| GET | `/api/v1/volumes` | `200 + StorageVolumeListResponse` |
| POST | `/api/v1/volumes` | `201 + StorageVolume` |
| GET | `/api/v1/volumes/{volume_id}` | `200 + StorageVolume` |
| DELETE | `/api/v1/volumes/{volume_id}` | `200 + StorageVolume` |
| GET | `/api/v1/filesystems` | `200 + StorageFilesystemListResponse` |
| POST | `/api/v1/filesystems` | `201 + StorageFilesystem` |
| GET | `/api/v1/filesystems/{filesystem_id}` | `200 + StorageFilesystem` |
| DELETE | `/api/v1/filesystems/{filesystem_id}` | `200 + StorageFilesystem` |
| GET | `/api/v1/objects` | `200 + StorageObjectListResponse` |
| POST | `/api/v1/objects` | `201 + StorageObject` |
| GET | `/api/v1/objects/{object_id}` | `200 + StorageObject` |
| DELETE | `/api/v1/objects/{object_id}` | `200 + StorageObject` |
| GET | `/api/v1/vector-stores` | `200 + VectorStoreListResponse` |
| POST | `/api/v1/vector-stores` | `201 + VectorStore` |
| GET | `/api/v1/vector-stores/{vector_store_id}` | `200 + VectorStore` |
| DELETE | `/api/v1/vector-stores/{vector_store_id}` | `200 + VectorStore` |
| POST | `/api/v1/vector-stores/{vector_store_id}/search` | `200 + VectorStoreSearchResponse` |

### 5.2 Create and Search Constraints

#### `POST /api/v1/volumes`

- `operationId`: `createStorageVolume`
- `requestBody.required`: `name`、`size_gib`、`idempotency_key`
- `success`: `201 + StorageVolume`
- `errors`: `400`、`401`、`403`

#### `POST /api/v1/filesystems`

- `operationId`: `createStorageFilesystem`
- `requestBody.required`: `name`、`size_gib`、`idempotency_key`
- `success`: `201 + StorageFilesystem`
- `errors`: `400`、`401`、`403`

#### `POST /api/v1/objects`

- `operationId`: `createStorageObject`
- `requestBody.required`: `bucket`、`key`、`idempotency_key`
- `success`: `201 + StorageObject`
- `errors`: `400`、`401`、`403`

#### `POST /api/v1/vector-stores`

- `operationId`: `createVectorStore`
- `requestBody.required`: `name`、`dimension`、`idempotency_key`
- `success`: `201 + VectorStore`
- `errors`: `400`、`401`、`403`

#### `POST /api/v1/vector-stores/{vector_store_id}/search`

- `operationId`: `searchVectorStore`
- `requestBody.required`: `vector`
- `success`: `200 + VectorStoreSearchResponse`
- `errors`: `400`、`401`、`403`、`404`

### 5.3 Non-Frozen Paths

以下路径当前不能写成已对齐 Core 的正式接口：

- `/api/v1/volumes/{volume_id}/snapshots`
- `/api/v1/volumes/{volume_id}:attach`
- `/api/v1/volumes/{volume_id}:detach`
- `/api/v1/filesystems/{filesystem_id}/mount-targets`
- `/api/v1/objects/{object_id}:upload`
- `/api/v1/objects/{object_id}:download`
- `/api/v1/object-buckets`

## 6. Operation and Error Rules

### 6.1 Operation Availability

| 资源 | list/get | create | delete | action | 说明 |
|---|---|---|---|---|---|
| 块存储卷 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可以提示用户先确认是否存在潜在引用影响 |
| 文件存储 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可以提示用户先确认是否存在活跃挂载影响 |
| 对象元数据 | 可用 | 可用 | 可用 | 无额外冻结动作 | 页面可以提示用户先确认是否存在潜在引用影响 |
| 向量存储 | 可用 | 可用 | 可用 | `search` 可用 | 搜索前目标向量库必须存在且可访问 |

### 6.2 Deletion Boundary

- 页面可以提示用户先确认是否存在潜在依赖影响，例如卷引用、文件系统活跃挂载、对象关键引用、向量库活动依赖，但当前权威源未冻结这些冲突校验作为统一正式后端契约
- 当前 `DELETE` 显式冻结返回码为 `200`、`401`、`403`、`404`
- 因此本轮文档不能把 `409 CONFLICT` 写成已冻结存储删除契约

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
| `BAD_REQUEST` | 400 | 参数格式错误、容量非法、向量维度错误 |
| `UNAUTHORIZED` | 401 | 未登录或凭证失效 |
| `FORBIDDEN` | 403 | 无访问或操作权限 |
| `NOT_FOUND` | 404 | 资源不存在或已删除 |

## 7. Examples

### 7.1 CreateStorageVolumeRequest

```json
{
  "idempotency_key": "idem-volume-001",
  "name": "volume-data-01",
  "size_gib": 200,
  "storage_class": "standard"
}
```

### 7.2 VectorStoreSearchRequest

```json
{
  "vector": [0.12, 0.93, 0.44],
  "top_k": 10,
  "filter": {
    "kb_id": "kb-001"
  }
}
```

## 8. Risks and Prerequisites

### 8.1 Risks

| 风险 | 影响 | 处理方式 |
|---|---|---|
| 旧路径仍沿用 `/api/v1/storage/*` | 会造成前后端继续对错路径 | 统一改为现有 `/api/v1/*` 存储资源路径 |
| 旧文档继续要求 `X-Tenant-Id` | 会与 Core 网关边界冲突 | 文档统一改为认证上下文获取 |
| 把未冻结能力或未冻结返回码写进正式契约 | 后续 Core 对齐会失真 | 在正文和 SPEC 中单列待补能力 |

### 8.2 Core Alignment Prerequisites

- 后续若要把 `快照` 写入正式接口，必须先扩充 `Core v1.yaml`
- 后续若要增加 `挂载 / 卸载`、`对象上传 / 下载`、`桶策略`、`文件系统挂载目标`，必须先补充正式 schema 与路径
- 后续若要扩展向量存储写入、索引刷新等动作，也必须先补充 `Core v1.yaml`
