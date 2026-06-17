# SPEC: Console kb-doc-intelligence

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/knowledge/prd-console-kb-doc-intelligence.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md`  
> Revised: 2026-06-17

## 1. Summary

文档增量智能分析（规划）；analyze 202 + intelligence 只读。合入前标注**规划**。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| POST | `.../documents/{doc_id}/analyze` | `analyzeKnowledgeBaseDocument` | 202 | write |
| GET | `.../documents/{doc_id}/intelligence` | `getKnowledgeBaseDocumentIntelligence` | 200 | read |
| GET | `.../intelligence/history` | `listKnowledgeBaseDocumentIntelligenceRuns` | 200 | read（3b） |

路径前缀：`/api/v1/svc/knowledge-bases/{kb_id}`。

### 2.3 Planned Schemas

- `AnalyzeKnowledgeBaseDocumentRequest`（features: summary|entities|chunk_optimization）
- `KBDocumentIntelligence`

### 2.4 AsyncTask

- `task_type`: `kb.document.analyze`
- 查询：Core `GET /api/v1/tasks/{task_id}`

## 3. Page Scope

- 分析触发、结果展示、任务中心跳转
- **Non-Goals**：替代 upload 202 解析

## 4. 创建前置条件

见主维护源。`DOCUMENT_NOT_INDEXED`、`ANALYSIS_ALREADY_RUNNING` 建议 422。

## 5. Handler 验收（合入 YAML 后）

```bash
curl -X POST .../knowledge-bases/$KB/documents/$DOC/analyze -d '{"idempotency_key":"..."}'
curl .../knowledge-bases/$KB/documents/$DOC/intelligence
```

## 6. 主维护源

- `docs/console-modules/knowledge/kb-doc-intelligence.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
