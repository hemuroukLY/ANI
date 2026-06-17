# SPEC: Console kb-chat-history

> Source: `tasks/modules/prd/console/knowledge/prd-console-kb-chat-history.md`  
> Revised: 2026-06-17

## 1. Summary

知识库下**会话列表**（session）只读页；支持从会话进入问答续聊。属于 **Services / KnowledgeBases** 子能力，与 `kb-qa-flow` 共享 session 契约。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | 权限 |
|---|---|---|---|---|
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/sessions` | `listKnowledgeBaseSessions` | `200 + KnowledgeBaseSessionListResponse` | readers 或 editors |

Query：`limit`、`cursor`（以 OpenAPI 为准）。

### 2.3 Verified Schemas

- `KnowledgeBaseSession`（`id`、`title`、`created_at` 等，以 YAML 为准）
- `KnowledgeBaseSessionListResponse`（以 YAML 为准）

## 3. Page Scope

- 会话侧栏/列表；点击进入 `kb-qa-flow` 并携带 `session_id`
- 展示会话标题、创建时间、最近活跃（若 schema 提供）
- 新会话创建发生在 `kb-qa-flow`（query/stream），不在本列表页

## 4. Non-Goals

- 会话 DELETE — **YAML 未声明** `DELETE /sessions/{session_id}`
- PATCH 重命名会话 — **YAML 未声明**
- 跨知识库会话聚合 — 非本页职责
- **当前 YAML 未声明 `422`**

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 知识库读权限 | readers 或 editors | `403 FORBIDDEN` |
| 知识库存在 | 当前租户可见 | `404 NOT_FOUND` |

## 6. 操作可用性矩阵

| 操作 | 只读用户 | 编辑者 | 管理员 |
|---|---|---|---|
| 查看会话列表 | ✅ | ✅ | ✅ |
| 进入续聊 | ✅ | ✅ | ✅ |
| 删除/重命名会话 | ❌ | ❌ | ❌（待补） |

## 7. 主维护源

- `docs/console-modules/knowledge/kb-chat-history.md`
- 相关：`kb-qa-flow.md`、`kb-source-citation.md`、`knowledge-base.md`
- TASK：`TASK-SVC-014`

## 8. Handler 验收（Services 团队）

```bash
curl -H "Authorization: Bearer $TOKEN" "$BASE/api/v1/svc/knowledge-bases/{kb_id}/sessions"
```

- 分页：遵循 YAML `limit`/`cursor`；排序以服务端默认为准
- 不把未声明的 DELETE 写成可用操作
- OpenAPI 已声明 ≠ handler 已实现
