# SPEC: Console kb-source-citation

> Source: `tasks/modules/prd/console/knowledge/prd-console-kb-source-citation.md`  
> Revised: 2026-06-17

## 1. Summary

知识库问答召回的**来源引用**只读明细页；展示 `document_id`、`snippet`、`score`，与 `kb-qa-flow` citations 区块同源。属于 **Services / KnowledgeBases** 子能力。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | 权限 |
|---|---|---|---|---|
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/citations` | `listKnowledgeBaseCitations` | `200 + KnowledgeBaseCitationListResponse` | readers 或 editors |

Query：`limit`、`cursor`（若 YAML 声明）；可选 `session_id`（以 OpenAPI 为准）。

### 2.3 Verified Schemas

- `KnowledgeBaseCitation`（`document_id`、`snippet`、`score` 等，以 YAML 为准）
- `KnowledgeBaseCitationListResponse`（`items`、`next_cursor` 等，以 YAML 为准）

## 3. Page Scope

- 按知识库展示引用列表；与问答页 citations 区块同源
- 支持按会话/时间筛选（仅当 YAML query 支持）
- 跳转对应文档详情（`documents/{doc_id}`，见 `knowledge-base.md`）

## 4. Non-Goals

- 引用详情单条 GET — **YAML 未声明** `{citation_id}`
- 引用导出/审计 — 待 Services 规划
- 与向量检索 debug 视图 — 非 Console 冻结能力
- 本页无写接口；**当前 YAML 未声明 `422`**

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 知识库读权限 | readers 或 editors | `403 FORBIDDEN` |
| 知识库存在 | 当前租户可见 | `404 NOT_FOUND` |

## 6. 操作可用性矩阵

| 操作 | 只读用户 | 编辑者 | 管理员 |
|---|---|---|---|
| 查看引用列表 | ✅ | ✅ | ✅ |
| 导出引用 | ❌ | ❌ | ❌（待补） |
| 编辑权限 | ❌ | 见 `kb-permissions.md` | ✅ |

## 7. 主维护源

- `docs/console-modules/knowledge/kb-source-citation.md`
- 相关：`kb-qa-flow.md`、`knowledge-base.md`、`kb-chat-history.md`
- TASK：`TASK-SVC-014`

## 8. Handler 验收（Services 团队）

```bash
curl -H "Authorization: Bearer $TOKEN" "$BASE/api/v1/svc/knowledge-bases/{kb_id}/citations"
```

- 空列表：`200` + `items: []`（非 404）
- 与 kb-qa-flow citations 字段口径一致
- OpenAPI 已声明 ≠ handler 已实现
