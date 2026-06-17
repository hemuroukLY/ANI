# Phase 3 — 知识库会议智能 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`knowledge/kb-meeting-intelligence.md`  
> **TASK**：`TASK-SVC-018` 子项（知识库智能）  
> **原则**：会议 ingest 异步 `202 + AsyncTask`；产出纪要、行动项并可选写入 KB 文档。

---

## 1. 与文档智能 / 第三方集成的边界

| 维度 | 文档智能 | 会议智能（本草案） |
|---|---|---|
| 输入 | 已上传 KB 文档 | 录音/转写文件、外部会议链接 |
| 路径 | `.../documents/{id}/analyze` | `.../meetings/ingest` |
| 产出 | 摘要/实体/分块 | 纪要、行动项、可选 `KBDocument` |
| 集成 | — | Teams/Zoom 为 **source 枚举**；OAuth 走 `integration-third-party.md`（Phase 3b） |

与 `kb-doc-intelligence.md` 共享 NLP 管线（Services 实现说明，非 Console 契约）。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [KnowledgeBases]
x-ani-rbac-scope:
  read:  scope:knowledge-bases:read
  write: scope:knowledge-bases:write
```

---

## 3. Schemas（草案）

### IngestKnowledgeBaseMeetingRequest

```yaml
IngestKnowledgeBaseMeetingRequest:
  type: object
  required: [idempotency_key, title, source]
  properties:
    idempotency_key: { type: string, format: uuid }
    title:           { type: string }
    source:
      type: string
      enum: [file, teams, zoom, manual_transcript]
    file_ref:
      type: object
      nullable: true
      description: source=file 时必填；对象存储或已上传 file id
      properties:
        object_key: { type: string, nullable: true }
        document_id: { type: string, format: uuid, nullable: true }
    external_meeting_id: { type: string, nullable: true, description: teams/zoom 外部 ID }
    occurred_at:         { type: string, format: date-time, nullable: true }
    language:            { type: string, nullable: true, default: zh-CN }
    create_kb_document:  { type: boolean, default: true, description: 是否将纪要写入 KB 为新文档 }
```

### KBMeetingRecord

```yaml
KBMeetingRecord:
  type: object
  required: [id, knowledge_base_id, title, status, created_at]
  properties:
    id:                { type: string, format: uuid }
    knowledge_base_id: { type: string, format: uuid }
    title:             { type: string }
    source:            { type: string, enum: [file, teams, zoom, manual_transcript] }
    status:            { type: string, enum: [ingesting, processing, completed, failed] }
    occurred_at:       { type: string, format: date-time, nullable: true }
    summary:           { type: string, nullable: true }
    action_items:
      type: array
      items:
        type: object
        required: [text]
        properties:
          text:       { type: string }
          assignee:   { type: string, nullable: true }
          due_at:     { type: string, format: date-time, nullable: true }
          completed:  { type: boolean, default: false }
    linked_document_id: { type: string, format: uuid, nullable: true }
    task_id:            { type: string, format: uuid, nullable: true }
    created_at:         { type: string, format: date-time }
    updated_at:         { type: string, format: date-time, nullable: true }
    error_message:      { type: string, nullable: true }
```

### KBMeetingListResponse

```yaml
KBMeetingListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/KBMeetingRecord' } }
    next_cursor: { type: string, nullable: true }
```

---

## 4. Operations（草案）

### 4.1 `POST /api/v1/svc/knowledge-bases/{kb_id}/meetings/ingest`

- operationId: `ingestKnowledgeBaseMeeting`
- Body: `IngestKnowledgeBaseMeetingRequest`
- 成功：`202 + AsyncTask`（`task_type` 建议 `kb.meeting.ingest`）+ 响应头或 body 可选 `meeting_id`（**评审固定**：建议 body 扩展 `{ task: AsyncTask, meeting: KBMeetingRecord }` 或仅 AsyncTask + `resource_id=meeting_id`）
- 错误：
  - `401`、`403`、`404`（kb）
  - `400`：source=file 但缺 file_ref
  - `422`：外部会议不可达 — 建议 `MEETING_SOURCE_UNAVAILABLE`

### 4.2 `GET /api/v1/svc/knowledge-bases/{kb_id}/meetings`

- operationId: `listKnowledgeBaseMeetings`
- Query: `status`、`source`、`limit`（default 50）、`cursor`
- 成功：`200 + KBMeetingListResponse`

### 4.3 `GET /api/v1/svc/knowledge-bases/{kb_id}/meetings/{meeting_id}`

- operationId: `getKnowledgeBaseMeeting`
- 成功：`200 + KBMeetingRecord`
- 错误：`404`

### 4.4 `DELETE /api/v1/svc/knowledge-bases/{kb_id}/meetings/{meeting_id}`（3b）

- operationId: `deleteKnowledgeBaseMeeting`
- Body: `{ idempotency_key }`
- 成功：`204`
- 不自动删除 linked_document（**评审确认**）

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 编辑者 | 说明 |
|---|---|---|---|
| 会议列表/详情 | 可用 | 可用 | GET |
| ingest 会议 | 不可用 | 可用 | POST ingest |
| 查看纪要/行动项 | 可用 | 可用 | GET meeting |
| 跳转关联文档 | 可用 | 可用 | linked_document_id |
| 删除会议 | 不可用 | 可用 | DELETE（3b） |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","title":"周会","source":"file","file_ref":{"document_id":"'"$DOC_ID"'"}' \
  "$BASE/api/v1/svc/knowledge-bases/$KB_ID/meetings/ingest"

curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/knowledge-bases/$KB_ID/meetings/$MEETING_ID"
```

---

## 7. 评审检查清单

- [ ] ingest 与 documents upload 不混 path
- [ ] AsyncTask resource 指向 meeting_id
- [ ] Teams/Zoom OAuth 不阻塞 3a（file/manual 先行）
- [ ] 合入后更新 `kb-meeting-intelligence.md`

---

## 相关文件

- `docs/console-modules/knowledge/kb-meeting-intelligence.md`
- `docs/console-modules/knowledge/kb-doc-intelligence.md`
- `docs/console-modules/integration/integration-third-party.md`
