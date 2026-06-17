# Phase 3 — 知识库文档智能 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`knowledge/kb-doc-intelligence.md`  
> **TASK**：`TASK-SVC-018` 子项（知识库智能）  
> **原则**：挂接 `KnowledgeBases` 资源组；异步分析返回 `202 + AsyncTask`，轮询 Core `GET /tasks/{task_id}`。

---

## 1. 与知识库主模块 / 向量写入的边界

| 维度 | 文档 CRUD（已冻结） | 文档智能（本草案） |
|---|---|---|
| 路径 | `.../documents` POST/GET/DELETE | `.../documents/{doc_id}/analyze`、`.../intelligence` |
| 触发时机 | 上传即触发基础解析（`202 + KBDocument`） | **增量**结构化分析：摘要、实体、分块优化 |
| 任务类型 | `kb.parse` / `kb.index`（upload 侧） | 建议 `kb.document.analyze` |
| 下游 | 向量索引 | `vector-store-write.md` 可选联动 re-index |

Console：上传成功 ≠ 智能分析已完成；分析结果只读页独立于文档列表。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [KnowledgeBases]
x-ani-rbac-scope:
  read:  scope:knowledge-bases:read
  write: scope:knowledge-bases:write
```

与 `knowledge-base.md` documents 子资源共用 scope；编辑者方可 POST analyze。

---

## 3. Schemas（草案）

### AnalyzeKnowledgeBaseDocumentRequest

```yaml
AnalyzeKnowledgeBaseDocumentRequest:
  type: object
  required: [idempotency_key]
  properties:
    idempotency_key: { type: string, format: uuid }
    features:
      type: array
      items:
        type: string
        enum: [summary, entities, chunk_optimization]
      description: 默认全选；可部分启用
    force_reanalyze: { type: boolean, default: false, description: 已有结果时是否强制重跑 }
```

### KBDocumentIntelligence

```yaml
KBDocumentIntelligence:
  type: object
  required: [document_id, status, updated_at]
  properties:
    document_id:   { type: string, format: uuid }
    status:        { type: string, enum: [pending, running, completed, failed] }
    summary:       { type: string, nullable: true }
    entities:
      type: array
      items:
        type: object
        required: [label, value]
        properties:
          label: { type: string }
          value: { type: string }
          confidence: { type: number, minimum: 0, maximum: 1, nullable: true }
    chunk_count:   { type: integer, minimum: 0, nullable: true }
    task_id:       { type: string, format: uuid, nullable: true, description: 进行中任务 ID }
    analyzed_at:   { type: string, format: date-time, nullable: true }
    updated_at:    { type: string, format: date-time }
    error_message: { type: string, nullable: true }
```

---

## 4. Operations（草案）

### 4.1 `POST /api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}/analyze`

- operationId: `analyzeKnowledgeBaseDocument`
- Path params: `kb_id`, `doc_id`
- Body: `AnalyzeKnowledgeBaseDocumentRequest`
- 成功：`202 + AsyncTask`（`task_type` 建议 `kb.document.analyze`；`resource_type=knowledge_base_document`，`resource_id=doc_id`）
- 错误：
  - `401`、`403`
  - `404`：kb 或 doc 不存在
  - `422`：文档未就绪（未 index 完成）— 建议 `DOCUMENT_NOT_INDEXED`；或分析已在跑 — `ANALYSIS_ALREADY_RUNNING`

### 4.2 `GET /api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}/intelligence`

- operationId: `getKnowledgeBaseDocumentIntelligence`
- 成功：`200 + KBDocumentIntelligence`
- 错误：`401`、`403`、`404`
- 说明：无结果时 `status=pending` 且 summary/entities 为空（或 `404` — **评审时二选一固定**）

### 4.3 `GET /api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}/intelligence/history`（3b）

- operationId: `listKnowledgeBaseDocumentIntelligenceRuns`
- Query: `limit`、`cursor`
- 成功：`200` + items（历史 run 摘要）
- Phase 3a 可不实现，Console 隐藏入口

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 编辑者 | 说明 |
|---|---|---|---|
| 查看分析结果 | 可用 | 可用 | GET intelligence |
| 触发分析 | 不可用 | 可用 | POST analyze |
| 跳转任务中心 | 可用 | 可用 | AsyncTask |
| 跳转向量写入 | 跳转 | 跳转 | chunk 优化后 re-index |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","features":["summary","entities"]}' \
  "$BASE/api/v1/svc/knowledge-bases/$KB_ID/documents/$DOC_ID/analyze"
# 期望 202；body 含 AsyncTask.id

curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/knowledge-bases/$KB_ID/documents/$DOC_ID/intelligence"
# 期望 200 + KBDocumentIntelligence
```

---

## 7. 评审检查清单

- [ ] 与 `POST .../documents` upload 202 语义不冲突
- [ ] AsyncTask 与 `async-task-center.md` 字段一致
- [ ] analyze 必填 `idempotency_key`
- [ ] 合入后更新 `kb-doc-intelligence.md`

---

## 相关文件

- `docs/console-modules/knowledge/kb-doc-intelligence.md`
- `docs/console-modules/knowledge/knowledge-base.md`
- `docs/console-modules/compute/storage/vector-store-write.md`
