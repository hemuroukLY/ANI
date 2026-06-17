# 向量存储 — 写入与索引

## 页面定位

向量库的文档写入与索引维护，扩展 `vector-storage.md`（已有 search）。

## 文档管理规则

- 本文是向量写入子模块主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 搜索见 `vector-storage.md` `POST .../search`

## Core 层要求

<!-- ADDED-TO-YAML: POST .../documents (Phase 2) -->

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| POST | `/api/v1/vector-stores/{vector_store_id}/documents` | `insertVectorStoreDocuments` | `scope:vector-stores:write` |
| POST | `/api/v1/vector-stores/{vector_store_id}/search` | `searchVectorStore` | `scope:vector-stores:search` |

写入请求：`VectorStoreDocumentInsertRequest`（含 `idempotency_key`）；响应 `202 + VectorStoreDocumentInsertResponse`。

## 页面职责

- 向量库详情：批量写入 UI、任务状态、search 入口
- 索引 rebuild / 批量导入 — **YAML 未声明** → TODO-YAML

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| `vector_store_id` | 存在 | `404 NOT_FOUND` |
| 写权限 | `scope:vector-stores:write` | `403 FORBIDDEN` |
| 请求体 | `idempotency_key` + documents | `400 BAD_REQUEST` |
| VectorStore 状态 | 须 `ready`（产品语义） | `422 PRECONDITION_FAILED`（write/search 已声明） |

## 操作可用性矩阵

| 操作 | 只读用户 | 向量库管理员 |
|---|---|---|
| search | ✅（search scope） | ✅ |
| 写入 documents | ❌ | ✅ |
| rebuild 索引 | ❌ | ❌（待 YAML） |

## 接口冻结规则

### `POST /api/v1/vector-stores/{vector_store_id}/documents`

- 成功：`202 + VectorStoreDocumentInsertResponse`
- 错误：`400`、`401`、`403`、`404`、`422`

### `POST /api/v1/vector-stores/{vector_store_id}/search`

- 成功：`200 + VectorStoreSearchResponse`
- 错误：`400`、`401`、`403`、`404`、`422`

## 待补边界

- 索引 rebuild / bulk import 独立路径 — **TODO-YAML**
- `422` 的 `code` 举例（如 `VECTOR_STORE_NOT_READY`）— 待 Core description 补充
- 写入进度与 `AsyncTask` 关联 — 以实现为准

## 验收标准

- [ ] write 与 search 分路径描述
- [ ] 不自造 bulk import schema
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
