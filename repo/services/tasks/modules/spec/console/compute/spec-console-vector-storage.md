# SPEC: Console 向量存储

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-vector-storage.md`
> Generated: 2026-06-14 | Scope: `Console / 算力与云资源 / 存储管理 / 向量存储`

## 1. Summary

- Source: `tasks/modules/prd/console/compute/prd-console-vector-storage.md`
- User Stories covered: `US-001` ~ `US-004`
- Functional Requirements covered: `FR-1` ~ `FR-4`

## 2. Architecture

- 页面属于 `Console / 算力与云资源 / 存储管理`
- 页面直接消费 `Core /api/v1/vector-stores*`
- 页面不定义新的 `Services /api/v1/svc/*` 向量资源契约
- 知识库等页面中的向量资源摘要和跳转只消费本模块已冻结定义

## 3. Data Model

- `VectorStore`
- `VectorStoreListResponse`
- `CreateVectorStoreRequest`
- `VectorStoreSearchRequest`
- `VectorStoreSearchResponse`

## 4. API Design

| Method | Path | Description | Response |
|---|---|---|---|
| GET | `/api/v1/vector-stores` | 查询向量存储列表 | `200 + VectorStoreListResponse` |
| POST | `/api/v1/vector-stores` | 创建向量存储 | `201 + VectorStore` |
| GET | `/api/v1/vector-stores/{vector_store_id}` | 查询向量存储 | `200 + VectorStore` |
| DELETE | `/api/v1/vector-stores/{vector_store_id}` | 删除向量存储 | `200 + VectorStore` |
| POST | `/api/v1/vector-stores/{vector_store_id}/search` | 搜索向量存储 | `200 + VectorStoreSearchResponse` |

## 5. Business Logic

- 向量存储页默认展示当前租户可见向量资源
- `VectorStore.state` 使用 YAML 已冻结枚举：`pending` / `ready` / `failed` / `deleting` / `deleted`（**不是** `available`）
- 非 `ready` 状态时搜索入口置灰；`searchVectorStore` **当前 YAML 未声明 `422`**
- 总览区负责说明资源边界与能力状态
- 创建区只承接现有 `CreateVectorStoreRequest`
- 搜索区只承接现有 `VectorStoreSearchRequest`
- 对未冻结能力，只记录为“待补 Core 契约”，不能写成正式接口表

## 6. Error Handling

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-xxx"
}
```

## 7. Security

- 全部向量存储接口依赖标准认证
- 后端必须从认证上下文获取租户边界
- 前端不可信任也不显式传 `tenant_id / X-Tenant-Id`

## 8. Testing Strategy

- 校验向量存储列表、详情、创建、删除、搜索路径与返回码对齐 `Core v1.yaml`
- 校验创建请求带有 `idempotency_key`
- 校验待补能力不会误显示为已冻结接口
