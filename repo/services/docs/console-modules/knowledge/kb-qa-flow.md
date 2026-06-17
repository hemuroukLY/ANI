# 知识库问答流程

## 页面定位

`知识库问答流程` 是 `Console / 知识库与 AI 应用` 下的交互子流程，用于在选定知识库后完成**提问 → 召回 → 回答展示 → 会话延续**的完整 UI 路径。

本模块属于 **Services / KnowledgeBases**，与 `knowledge-base.md` 主模块共享同一 OpenAPI 资源组。

## 文档管理规则

- 本文专注**问答交互流程**与 session 边界；CRUD 见 `knowledge-base.md`
- 一级权威源：`services/v1.yaml` 中 `query`、`query/stream`、`sessions`、`citations`
- 流式方法必须为 **GET**，不得写成 POST
- OpenAPI 已声明 ≠ handler 已实现

## Services 层要求

### 普通问答

- `POST /api/v1/svc/knowledge-bases/{kb_id}/query`
- `operationId`: `queryKnowledgeBase`
- 请求必填：`question`、`idempotency_key`
- 请求可选：`session_id`、`top_k`、`score_threshold`
- 响应：`KBQueryResponse`
- 错误：`400`、`401`、`404`

### 流式问答

- `GET /api/v1/svc/knowledge-bases/{kb_id}/query/stream`
- `operationId`: `streamQueryKnowledgeBase`
- Query 必填：`question`
- Query 可选：`session_id`
- 响应：`200` + `text/event-stream`
- 错误：`400`、`401`
- **方法必须为 GET**，不得写成 POST

### 对话历史

- `GET /api/v1/svc/knowledge-bases/{kb_id}/sessions`
- `operationId`: `listKnowledgeBaseSessions`
- 响应：`KnowledgeBaseSessionListResponse`
- 分页：`limit`、`cursor`
- 错误：`401`、`403`、`404`

<!-- ADDED-TO-YAML: sessions / citations (Services v1.yaml, Phase 2 2026-06-17) -->

### 来源引用

- `GET /api/v1/svc/knowledge-bases/{kb_id}/citations`
- `operationId`: `listKnowledgeBaseCitations`
- 响应：`KnowledgeBaseCitationListResponse`
- 错误：`401`、`403`、`404`

## 页面职责

- 知识库选择与会话列表
- 提问输入（支持流式 / 非流式切换）
- 展示回答与 citations（document_id、snippet、score）
- 会话 `session_id` 延续（nullable 首问创建新会话）
- 不替代知识库文档上传与权限管理

## 页面结构

```text
知识库问答
├── 知识库选择器
├── 会话侧栏（sessions 列表）
├── 对话主区
│   ├── 用户问题
│   ├── 助手回答
│   └── 引用来源（citations）
└── 输入框 + 发送（JSON / SSE）
```

## 核心任务流

1. 用户选择 `kb_id`，可选加载历史 `sessions`
2. 用户输入 `question`，携带 `idempotency_key`（非流式）
3. 非流式：`POST /query` → 渲染 `KBQueryResponse`
4. 流式：`GET /query/stream` → SSE 增量渲染
5. 后续轮次携带返回的 `session_id`

## 创建前置条件

| 依赖项 | 要求 | 未满足响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 知识库存在 | 有效 `kb_id` | `404 NOT_FOUND` |
| 读权限 | 租户成员可读 | `403 FORBIDDEN`（sessions/citations） |
| 请求体合法 | `question` + `idempotency_key` | `400 BAD_REQUEST` |
| 知识库可问答 | 建议 `KnowledgeBase` 就绪态 | **当前 YAML 未在 query 声明 `422`** |

## 操作可用性矩阵

| 操作 | 只读用户 | 编辑者 |
|---|---|---|
| 查看会话与历史 | ✅ | ✅ |
| 发起问答（非流式） | ✅ | ✅ |
| 发起问答（流式） | ✅ | ✅ |
| 修改知识库权限 | ❌ | ✅（见 permissions 路径） |

## 接口冻结规则

### `POST /api/v1/svc/knowledge-bases/{kb_id}/query`

- 成功：`200 + KBQueryResponse`
- 错误：`400`、`401`、`404`
- requestBody.required：`question`、`idempotency_key`
- **当前 YAML 未声明** `403`、`422`

### `GET /api/v1/svc/knowledge-bases/{kb_id}/query/stream`

- 成功：`200` + `text/event-stream`
- 错误：`400`、`401`
- **当前 YAML 未声明** `403`、`404`、`422`

### `GET /api/v1/svc/knowledge-bases/{kb_id}/sessions`

- 成功：`200 + KnowledgeBaseSessionListResponse`
- 错误：`401`、`403`、`404`

### `GET /api/v1/svc/knowledge-bases/{kb_id}/citations`

- 成功：`200 + KnowledgeBaseCitationListResponse`
- 错误：`401`、`403`、`404`

## 待补边界

- query 前置条件 `422`（知识库未就绪 / 无索引文档）— **当前 YAML 未声明**
- 流式 SSE 事件 schema 细化 — 以 handler 实现为准，文档不自造 event type 枚举
- 单条 session 详情 GET — **当前 YAML 未声明**
- 趋势指标与 `metering/usage` 口径对齐 — 非本模块冻结路径

## 与平台概览的关系

- 首页「知识库调用趋势」区块为用量摘要；本页为明细交互
- 趋势指标待与 `metering/usage` 或 Services 统计口径对齐

## 验收标准

- [ ] 流式方法为 GET，与 YAML 一致
- [ ] `session_id` 语义与 nullable 首问一致
- [ ] 不自造 `ChatMessage` 等未声明 schema
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
