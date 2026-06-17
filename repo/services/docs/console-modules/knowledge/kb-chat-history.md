# 知识库 — 对话历史

## 页面定位

展示知识库下**会话列表**（session），支持从会话进入问答续聊。

本页属于 **Services / KnowledgeBases** 子能力，与 `kb-qa-flow.md` 共享 session 契约。

## 文档管理规则

- 本文是对话历史子页的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 会话创建/续聊入口以 `kb-qa-flow.md` 为准
- TASK：`TASK-SVC-014`

## Services 层要求

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/sessions` | `listKnowledgeBaseSessions` |

Query：`limit`、`cursor`（以 OpenAPI 为准）。

响应：`KnowledgeBaseSessionListResponse`。

Schema：`KnowledgeBaseSession`（`id`、`title`、`created_at` 等，以 YAML 为准）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 知识库读权限 | readers 或 editors | `403 FORBIDDEN` |
| 知识库存在 | 当前租户可见 | `404 NOT_FOUND` |

新会话创建发生在 `kb-qa-flow`（query/stream），不在本列表页。**当前 YAML 未声明 sessions DELETE**；**未声明 `422`**。

## 页面职责

- 会话侧栏/列表；点击进入 `kb-qa-flow` 并携带 `session_id`
- 展示会话标题、创建时间、最近活跃（若 schema 提供）
- 不提供会话 DELETE（**YAML 未声明**）

## 操作可用性矩阵

| 操作 | 只读用户 | 编辑者 | 管理员 |
|---|---|---|---|
| 查看会话列表 | 可用 | 可用 | 可用 |
| 进入续聊 | 可用 | 可用 | 可用 |
| 删除/重命名会话 | 不可用 | 不可用 | 不可用（待补） |

## 接口冻结规则

### `GET /api/v1/svc/knowledge-bases/{kb_id}/sessions`

- 成功：`200 + KnowledgeBaseSessionListResponse`
- 错误：`401`、`403`、`404`
- 分页：遵循 YAML `limit`/`cursor`
- 排序：以服务端默认为准，Console 不假定字段

## 待补边界

- `DELETE /sessions/{session_id}` — **YAML 未声明**
- `PATCH` 重命名会话 — **YAML 未声明**
- 跨知识库会话聚合 — 非本页职责

## 相关模块

- `kb-qa-flow.md`
- `kb-source-citation.md`
- `knowledge-base.md`

## 验收标准

- [ ] 与 kb-qa-flow 中 sessions 口径一致
- [ ] 路径与 services/v1.yaml 一致
- [ ] 不把未声明的 DELETE 写成可用操作
