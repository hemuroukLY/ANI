# SPEC: Console 知识库

> Technical specification derived from: `tasks/modules/prd/console/knowledge/prd-console-knowledge-base.md`
> Source of truth: `ANI-main/repo/api/openapi/services/v1.yaml`
> Scope: `Console / 知识库与 AI 应用 / 知识库管理`

## 1. Summary

本 SPEC 定义 `Console / 知识库与 AI 应用 / 知识库管理` 的技术边界、字段映射、页面承接范围和接口冻结规则。知识库属于 `Services / KnowledgeBases`，正式前缀为 `/api/v1/svc`，不扩写为 `Core` 资源。

## 2. Source of Truth

### 2.1 Primary Authority

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Frozen Path Group

- `/knowledge-bases`
- `/knowledge-bases/{kb_id}`
- `/knowledge-bases/{kb_id}/documents`
- `/knowledge-bases/{kb_id}/documents/{doc_id}`
- `/knowledge-bases/{kb_id}/query`
- `/knowledge-bases/{kb_id}/query/stream`

### 2.3 Frozen Schemas

- `KnowledgeBase`
- `KBDocument`
- `KBQueryResponse`
- `ErrorResponse`

## 3. API Design

### 3.1 Endpoint Matrix

| Method | Path | operationId | summary | Success | Error Responses |
|---|---|---|---|---|---|
| GET | `/api/v1/svc/knowledge-bases` | `listKnowledgeBases` | 知识库列表 | `200 + {items: KnowledgeBase[]}` | 当前权威源未单独列出 |
| POST | `/api/v1/svc/knowledge-bases` | `createKnowledgeBase` | 创建知识库 | `201 + KnowledgeBase` | 当前权威源未单独列出 |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}` | `getKnowledgeBase` | 获取知识库详情 | `200 + KnowledgeBase` | `404` |
| DELETE | `/api/v1/svc/knowledge-bases/{kb_id}` | `deleteKnowledgeBase` | 删除知识库 | `204` | `404` |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/documents` | `listKnowledgeBaseDocuments` | 文档列表 | `200 + {items: KBDocument[]}` | 当前权威源未单独列出 |
| POST | `/api/v1/svc/knowledge-bases/{kb_id}/documents` | `uploadKnowledgeBaseDocument` | 上传文档 | `202 + KBDocument` | 当前权威源未单独列出 |
| DELETE | `/api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}` | `deleteKnowledgeBaseDocument` | 删除知识库文档 | `204` | `404` |
| POST | `/api/v1/svc/knowledge-bases/{kb_id}/query` | `queryKnowledgeBase` | 知识库问答 | `200 + KBQueryResponse` | `400`、`401`、`404` |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/query/stream` | `streamQueryKnowledgeBase` | 知识库问答 SSE 流式 | `200 + text/event-stream` | `400`、`401` |

### 3.2 Request Rules

#### `POST /api/v1/svc/knowledge-bases`

- 必填字段：
  - `idempotency_key`
  - `name`
- 可选字段：
  - `description`
  - `embedding_model`
  - `chunk_size`
  - `top_k`
- 字段约束：
  - `idempotency_key`: `uuid`
  - `embedding_model`: 默认 `bge-m3`
  - `chunk_size`: 默认 `512`
  - `top_k`: 默认 `5`

#### `POST /api/v1/svc/knowledge-bases/{kb_id}/documents`

- 路径参数：
  - `kb_id`: `uuid`
- 请求体：
  - `multipart/form-data`
  - 必填字段 `file`
- 成功响应：
  - `202 + KBDocument`
  - 语义为文档已上传，解析任务进行中

#### `POST /api/v1/svc/knowledge-bases/{kb_id}/query`

- 路径参数：
  - `kb_id`: `uuid`
- 必填字段：
  - `question`
  - `idempotency_key`
- 可选字段：
  - `session_id`
  - `top_k`
  - `score_threshold`
- 字段约束：
  - `question`: `1-2000` 字符
  - `idempotency_key`: 最大长度 `128`
  - `top_k`: `1-20`，默认 `5`
  - `score_threshold`: `0.0-1.0`，默认 `0.3`

#### `GET /api/v1/svc/knowledge-bases/{kb_id}/query/stream`

- 路径参数：
  - `kb_id`: `uuid`
- 查询参数：
  - `question`，必填，`1-2000` 字符
  - `session_id`，可选，`uuid`

## 4. Data Model

### 4.1 KnowledgeBase

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `uuid` | 知识库唯一标识 |
| `name` | `string` | 知识库名称 |
| `status` | `string` | `active | rebuilding | deleted` |
| `doc_count` | `integer` | 文档数量 |
| `created_at` | `date-time` | 创建时间 |

### 4.2 KBDocument

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `uuid` | 文档唯一标识 |
| `knowledge_base_id` | `uuid` | 所属知识库 |
| `file_name` | `string` | 文件名 |
| `content_type` | `string \| null` | 内容类型 |
| `size_bytes` | `integer` | 文件大小 |
| `status` | `string` | `uploaded | parsing | indexed | failed | deleted` |
| `created_at` | `date-time` | 上传时间 |

### 4.3 KBQueryResponse

| 字段 | 类型 | 说明 |
|---|---|---|
| `answer` | `string` | 问答结果 |
| `sources` | `array` | 来源片段数组 |
| `session_id` | `uuid` | 会话标识 |
| `input_tokens` | `integer` | 输入 token 数 |
| `output_tokens` | `integer` | 输出 token 数 |

## 5. Page Design Constraints

### 5.1 List Page

- 承接知识库列表、创建入口、详情入口
- 列表字段必须来自 `KnowledgeBase`
- 不在列表页承诺向量化、权限、历史资源的正式子页

### 5.2 Detail Page

- 承接知识库基础信息、文档区、问答区
- 文档区和问答区在同一知识库上下文内独立反馈
- 不把未冻结能力扩写成独立资源页

### 5.3 State Handling

- 空态：无知识库时展示创建入口
- 文档过程态：`uploaded`、`parsing`、`indexed`、`failed`
- 知识库状态：`active`、`rebuilding`、`deleted`
- 流式问答必须区分开始、增量输出和结束

## 6. Rules

- 知识库业务对象归 `Services`
- 文档管理与问答都属于知识库主资源的已冻结能力
- 向量化、来源引用、对话历史、权限当前只作为模块边界
- 删除文档正式路径必须是 `/documents/{doc_id}`，不能省略 `doc_id`
- 流式问答正式方法是 `GET`，不能写成 `POST`

## 7. Error Handling

### 7.1 Unified Error Shape

```json
{"code":"UPPER_SNAKE","message":"...","request_id":"..."}
```

### 7.2 Expected Errors

- 知识库详情：`404 + ErrorResponse`
- 知识库删除：`404 + ErrorResponse`
- 文档删除：`404 + ErrorResponse`
- 同步问答：`400`、`401`、`404`
- 流式问答：`400`、`401`

### 7.3 Documentation Rule

- 不允许把文档上传写成同步索引完成
- 不允许把流式问答方法写成 `POST`
- 不允许把删除文档路径写成 `/documents`

## 8. Acceptance

- 文档中所有正式路径均位于 `/api/v1/svc/knowledge-bases*`
- 文档中的路径方法与 `services/v1.yaml` 保持一致
- 文档中不再出现错误的文档删除路径或错误的流式问答方法
- 文档中不把待补能力写成正式接口或正式 schema

## 9. Backfill Dependencies

- 如未来新增向量化、来源引用、对话历史、权限，需先冻结路径与 schema
- 如未来新增智能增强能力，需先明确资源归属
- 如未来新增会话资源，需先明确 `session_id` 的正式关系
