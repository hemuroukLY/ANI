# 知识库管理

## 页面定位

`知识库管理` 是 `Console / 知识库与 AI 应用` 下的租户侧知识管理页面，用于管理知识库主对象、文档和问答能力。

当前模块属于 `Services / KnowledgeBases`，一级权威源为 `ANI-main/repo/api/openapi/services/v1.yaml`，正式路径前缀为 `/api/v1/svc`。

## 文档管理规则

- 本文是 `知识库管理` 的主维护文档，知识库页面定义、边界口径和验收标准统一以本文为准
- 如 `PRD`、`SPEC`、HTML 摘要与本文出现差异，先对照 `ANI-main/repo/api/openapi/services/v1.yaml`，再统一回写
- 文档管理与问答说明只承接当前已冻结能力，不新增向量化、权限、历史等未冻结资源承诺
- 页面补充可以增强交互流与状态表达，但不得把 `Services` 能力改写成 `Core` 资源

## Services 层要求

- 知识库属于 `Services`
- 正式路径使用 `/api/v1/svc/knowledge-bases*`
- 页面不要求前端显式传 `tenant_id`
- 文档上传与问答都承接知识库主资源能力
- 当前模块不把向量化、来源引用、对话历史、权限写成已冻结独立资源
- 创建知识库与同步问答需要承接 `idempotency_key`
- 错误结构统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

## 页面职责

- 展示知识库列表与详情
- 创建和删除知识库
- 管理知识库文档列表、上传和删除
- 在知识库上下文内执行同步或流式问答

## 页面结构

- 列表页至少包含 `知识库列表`、`创建入口`、`详情入口`
- 详情页至少包含 `基础信息`、`文档区`、`问答区`
- 文档区与问答区必须在同一知识库上下文内独立反馈状态，不扩写新的资源子页

## 数据来源与分层约束

### 数据来源划分

- `Services` 数据
  - 知识库列表
  - 知识库创建
  - 知识库详情
  - 知识库删除
  - 文档列表 / 上传 / 删除
  - 同步问答 / 流式问答
- `Core` 数据
  - 当前页面不承接 `Core` 知识库资源契约
  - 如后续需展示用量、任务、存储摘要，只能通过跳转或引用方式承接

### 关键边界

- 知识库不写成 `Core` 资源
- 文档上传成功 `202` 表示解析任务进行中，不等于文档已完成索引
- 同步问答与流式问答都挂接在知识库主资源下，不扩写独立会话资源
- 向量化、来源引用、对话历史、权限 — 子模块详文见 `kb-source-citation.md`、`kb-chat-history.md`、`kb-permissions.md`（Phase 2 YAML 已声明；handler stub）

## 字段级定义

### 知识库字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `id` | `KnowledgeBase.id` | 知识库唯一标识 |
| `name` | `KnowledgeBase.name` | 知识库名称 |
| `status` | `KnowledgeBase.status` | `active / rebuilding / deleted` |
| `doc_count` | `KnowledgeBase.doc_count` | 文档数量 |
| `created_at` | `KnowledgeBase.created_at` | 创建时间 |

### 文档字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `id` | `KBDocument.id` | 文档唯一标识 |
| `knowledge_base_id` | `KBDocument.knowledge_base_id` | 所属知识库 |
| `file_name` | `KBDocument.file_name` | 文件名 |
| `content_type` | `KBDocument.content_type` | 内容类型 |
| `size_bytes` | `KBDocument.size_bytes` | 文件大小 |
| `status` | `KBDocument.status` | `uploaded / parsing / indexed / failed / deleted` |
| `created_at` | `KBDocument.created_at` | 上传时间 |

### 问答字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `answer` | `KBQueryResponse.answer` | 问答结果 |
| `sources` | `KBQueryResponse.sources` | 来源片段 |
| `session_id` | `KBQueryResponse.session_id` | 会话标识 |
| `input_tokens` | `KBQueryResponse.input_tokens` | 输入 token 数 |
| `output_tokens` | `KBQueryResponse.output_tokens` | 输出 token 数 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 空态 | 展示创建入口与用途说明 | 不伪造示例知识库 |
| 列表正常态 | 展示知识库名称、状态、文档数量 | 字段必须来自 `KnowledgeBase` |
| 文档上传态 | 展示上传成功、解析中、失败反馈 | 不把 `202` 写成“已索引完成” |
| 文档删除态 | 删除成功前保留当前文档上下文 | 不静默移除造成误判 |
| 问答同步态 | 展示完整回答和来源片段 | 结果来自 `KBQueryResponse` |
| 问答流式态 | 展示增量输出和结束反馈 | 使用 `GET /query/stream` |

## 字段口径与单位

| 字段 | 口径建议 | 单位/格式 |
|---|---|---|
| `name` | 知识库名称 | 字符串 |
| `doc_count` | 当前知识库文档数量 | 整数 |
| `size_bytes` | 文档大小 | bytes，可格式化为 KB / MB |
| `created_at` | 服务端创建时间 | date-time |
| `input_tokens` | 输入 token 数 | 整数 |
| `output_tokens` | 输出 token 数 | 整数 |

## 状态与能力口径

### 知识库状态

- `active`：知识库可正常使用
- `rebuilding`：知识库重建中
- `deleted`：知识库已删除

### 文档状态

- `uploaded`
- `parsing`
- `indexed`
- `failed`
- `deleted`

## 待补能力边界

- 向量化与索引独立资源
- 来源引用 — 详文 `kb-source-citation.md` <!-- ADDED-TO-YAML: GET /api/v1/svc/knowledge-bases/{kb_id}/citations (Services v1.yaml, Phase 2 2026-06-17) -->
- 对话历史 — 详文 `kb-chat-history.md` <!-- ADDED-TO-YAML: GET /api/v1/svc/knowledge-bases/{kb_id}/sessions (Services v1.yaml, Phase 2 2026-06-17) -->
- 知识库权限 — 详文 `kb-permissions.md` <!-- ADDED-TO-YAML: PUT /api/v1/svc/knowledge-bases/{kb_id}/permissions (Services v1.yaml, Phase 2 2026-06-17) -->
- 文档智能 / 会议智能 / 视频智能

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录与租户上下文 | 已认证 | `401 UNAUTHORIZED` / `403 FORBIDDEN` |
| 知识库写权限 | 创建/上传/问答需编辑或管理员 | `403 FORBIDDEN` |
| `idempotency_key` | 创建知识库与同步问答必填 | `400 BAD_REQUEST` |
| 文档上传 | 使用 `multipart/form-data` | `400 BAD_REQUEST` |

## 操作可用性矩阵

| 操作 | 只读用户 | 知识管理员/编辑成员 | 说明 |
|---|---|---|---|
| 查看知识库列表 | 可用 | 可用 | 读取 `GET /api/v1/svc/knowledge-bases` |
| 查看知识库详情 | 可用 | 可用 | 读取 `GET /api/v1/svc/knowledge-bases/{kb_id}` |
| 创建知识库 | 不可用 | 可用 | 返回 `201 + KnowledgeBase` |
| 删除知识库 | 不可用 | 可用 | 返回 `204` |
| 查看文档列表 | 可用 | 可用 | 读取 `GET /api/v1/svc/knowledge-bases/{kb_id}/documents` |
| 上传文档 | 不可用 | 可用 | 返回 `202 + KBDocument` |
| 删除文档 | 不可用 | 可用 | 读取 `DELETE /documents/{doc_id}` |
| 同步问答 | 不可用 | 可用 | `POST /query` |
| 流式问答 | 不可用 | 可用 | `GET /query/stream` |

## 删除前置校验与当前契约边界

- 当前正式契约仅确认：
  - `DELETE /api/v1/svc/knowledge-bases/{kb_id}` 成功返回 `204`
  - `DELETE /api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}` 成功返回 `204`
- 当前权威源未声明删除时的引用冲突、软删除或回收站规则，正文不得自造这些逻辑
- 若后续权威源补充删除校验、任务跟踪或恢复语义，需再更新本文

## 接口冻结规则

### `GET /api/v1/svc/knowledge-bases`

- `operationId`: `listKnowledgeBases`
- `success`: `200 + {items: KnowledgeBase[]}`
- 当前权威源未单独列出错误返回码，正文不得擅自补写

### `POST /api/v1/svc/knowledge-bases`

- `operationId`: `createKnowledgeBase`
- `success`: `201 + KnowledgeBase`
- `requestBody.required`: `idempotency_key`、`name`
- `optional fields`: `description`、`embedding_model`、`chunk_size`、`top_k`

### `GET /api/v1/svc/knowledge-bases/{kb_id}`

- `operationId`: `getKnowledgeBase`
- `success`: `200 + KnowledgeBase`
- `errors`: `404`

### `DELETE /api/v1/svc/knowledge-bases/{kb_id}`

- `operationId`: `deleteKnowledgeBase`
- `success`: `204`
- `errors`: `404`

### `GET /api/v1/svc/knowledge-bases/{kb_id}/documents`

- `operationId`: `listKnowledgeBaseDocuments`
- `success`: `200 + {items: KBDocument[]}`

### `POST /api/v1/svc/knowledge-bases/{kb_id}/documents`

- `operationId`: `uploadKnowledgeBaseDocument`
- `success`: `202 + KBDocument`
- `requestBody.contentType`: `multipart/form-data`
- `requestBody.required`: `file`
- `success semantics`: 文档已上传，解析任务进行中

### `DELETE /api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}`

- `operationId`: `deleteKnowledgeBaseDocument`
- `success`: `204`
- `errors`: `404`
- 删除文档路径必须包含 `doc_id`

### `POST /api/v1/svc/knowledge-bases/{kb_id}/query`

- `operationId`: `queryKnowledgeBase`
- `success`: `200 + KBQueryResponse`
- `requestBody.required`: `question`、`idempotency_key`
- `optional fields`: `session_id`、`top_k`、`score_threshold`
- `errors`: `400`、`401`、`404`

### `GET /api/v1/svc/knowledge-bases/{kb_id}/query/stream`

- `operationId`: `streamQueryKnowledgeBase`
- `success`: `200 + text/event-stream`
- `query.required`: `question`
- `query.optional`: `session_id`
- `errors`: `400`、`401`
- 流式问答方法为 `GET`，不能写成 `POST`

## 响应示例

### 创建知识库成功

```json
{
  "id": "8aa3f627-897b-4166-82d7-b6c3c5e1f4a1",
  "name": "finance-kb",
  "status": "active",
  "doc_count": 0,
  "created_at": "2026-06-14T12:00:00Z"
}
```

### 文档上传成功

```json
{
  "id": "e9c8c0c8-1f1c-4cdb-9011-5b8f76b80e52",
  "knowledge_base_id": "8aa3f627-897b-4166-82d7-b6c3c5e1f4a1",
  "file_name": "product-spec.pdf",
  "content_type": "application/pdf",
  "size_bytes": 245760,
  "status": "uploaded",
  "created_at": "2026-06-14T12:10:00Z"
}
```

### 同步问答成功

```json
{
  "answer": "该产品支持多租户推理服务部署。",
  "sources": [
    {
      "doc_id": "e9c8c0c8-1f1c-4cdb-9011-5b8f76b80e52",
      "file_name": "product-spec.pdf",
      "page": 3,
      "content": "平台支持多租户推理服务部署。",
      "score": 0.91
    }
  ],
  "session_id": "91d41a84-f9d7-4ef6-a7f5-6b5aa1e9c979",
  "input_tokens": 128,
  "output_tokens": 64
}
```

## 错误示例

### 问答参数错误

```json
{
  "code": "BAD_REQUEST",
  "message": "question is required",
  "request_id": "req-kb-400-001"
}
```

### 知识库不存在

```json
{
  "code": "NOT_FOUND",
  "message": "knowledge base not found",
  "request_id": "req-kb-404-001"
}
```

## 回填前置依赖

- 如未来补充向量化、来源引用、对话历史、权限，必须先在 `services/v1.yaml` 中冻结路径和 schema
- 如未来补充智能增强能力，必须先明确资源归属
- 如未来补充问答会话资源，必须先明确 `session_id` 的正式关系

## 待确认项

- 文档上传后是否需要未来单独的解析任务历史页
- SSE 结束事件和异常事件是否需要未来统一事件协议
- `doc_count` 是否需要未来拆分为可检索文档数与总文档数

## 回填验收标准

- 正文路径全部位于 `Services / KnowledgeBases`
- 路径方法与 `services/v1.yaml` 保持一致，尤其是删除文档路径必须包含 `doc_id`，流式问答方法必须为 `GET`
- 正文只引用权威源真实存在的 schema：`KnowledgeBase`、`KBDocument`、`KBQueryResponse`、`ErrorResponse`
- 正文不把待补能力写成正式接口或正式 schema
- HTML 摘要、PRD、SPEC 与本文件一致

## 产品经理补充定义

### 目标用户与权限视角

- 知识管理员：查看知识库列表、创建知识库、上传文档、执行问答、删除知识库或文档
- 内容编辑成员：查看知识库详情、上传文档、发起问答
- 只读用户：仅查看知识库和文档列表，不显示写操作入口

### 页面结构补充

- 首屏至少包含 `知识库列表`、`创建入口`、`详情入口`
- 知识库详情至少包含 `基础信息`、`文档区`、`问答区` 三个主区块
- 文档管理与问答都作为知识库主资源内的能力呈现，不额外拆出新的业务资源页

### 核心任务流

1. 用户创建知识库后进入详情，确认基础信息并继续上传文档
2. 用户上传文档后，在文档区查看提交结果、当前可见文档列表和失败反馈
3. 用户在问答区执行同步问答或流式问答，并在当前详情上下文中查看结果
4. 用户删除文档或知识库时，如后端返回冲突或限制，页面按统一错误结构反馈

### 页面状态与反馈

- 空态：当前无知识库时，页面展示创建入口与用途说明
- 文档态：只承接上传已提交、列表可见、删除处理中和失败提示，不把向量化流程写成当前必有状态
- 问答态：同步问答展示完整结果，流式问答展示增量输出和结束反馈
- 局部失败态：列表、文档区、问答区彼此独立，单区失败不阻断整页

### 跨模块协同

- 与 `向量存储`、首页、用量页仅通过摘要和跳转协同，不重写底层资源契约
- 来源引用、对话历史、权限、文档智能等能力仍是待补边界，不形成现有页签承诺
- 如果首页展示知识库活跃趋势，只能回跳到本页或问答区，不替代详情页

### 产品验收补充

- 用户必须能完成“创建知识库 -> 上传文档 -> 发起问答 -> 查看结果”的基本闭环
- 文档区和问答区必须有独立的空态、失败态与成功回流
- 页面必须明确区分已冻结的文档/问答能力与未冻结的智能增强能力
- 本页不得把知识库写成 `Core` 资源，也不得把向量化等 Phase 2 能力写成现有正式资源
