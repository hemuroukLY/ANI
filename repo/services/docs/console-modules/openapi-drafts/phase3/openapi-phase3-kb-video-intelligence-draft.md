# Phase 3 — 知识库视频智能 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`knowledge/kb-video-intelligence.md`  
> **TASK**：`TASK-SVC-018` 子项（知识库智能）  
> **原则**：长视频 ingest 异步；转写、章节、关键帧/OCR 结果只读查询。

---

## 1. 与对象存储 / 会议智能的边界

| 维度 | 对象存储 upload | 视频智能（本草案） |
|---|---|---|
| 层 | Core `POST /objects/upload` | Services `.../videos/ingest` |
| 职责 | 字节上传 | 视频理解 + 写入 KB |
| 任务 | Core AsyncTask | `kb.video.ingest` |
| 产出 | object_key | 转写、章节、可选 KB 文档 |

与 `kb-meeting-intelligence.md` 区分：本模块面向 **视频**（章节/OCR/关键帧）；会议模块面向 **音频+纪要+行动项**。

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

### IngestKnowledgeBaseVideoRequest

```yaml
IngestKnowledgeBaseVideoRequest:
  type: object
  required: [idempotency_key, title]
  properties:
    idempotency_key: { type: string, format: uuid }
    title:           { type: string }
    source:
      type: string
      enum: [object_upload, url, existing_document]
      default: object_upload
    object_key:      { type: string, nullable: true, description: Core 对象 key }
    source_url:      { type: string, format: uri, nullable: true }
    document_id:     { type: string, format: uuid, nullable: true }
    language:        { type: string, nullable: true, default: zh-CN }
    features:
      type: array
      items:
        type: string
        enum: [transcript, chapters, keyframes, ocr]
      description: 默认全选
    create_kb_document: { type: boolean, default: true }
```

### KBVideoChapter

```yaml
KBVideoChapter:
  type: object
  required: [title, start_sec]
  properties:
    title:      { type: string }
    start_sec:  { type: number, minimum: 0 }
    end_sec:    { type: number, minimum: 0, nullable: true }
    summary:    { type: string, nullable: true }
```

### KBVideoRecord

```yaml
KBVideoRecord:
  type: object
  required: [id, knowledge_base_id, title, status, created_at]
  properties:
    id:                { type: string, format: uuid }
    knowledge_base_id: { type: string, format: uuid }
    title:             { type: string }
    duration_sec:      { type: number, minimum: 0, nullable: true }
    status:            { type: string, enum: [ingesting, processing, completed, failed] }
    transcript:        { type: string, nullable: true }
    chapters:
      type: array
      items: { $ref: '#/components/schemas/KBVideoChapter' }
    keyframe_urls:
      type: array
      items: { type: string, format: uri }
      nullable: true
    ocr_text:          { type: string, nullable: true }
    linked_document_id: { type: string, format: uuid, nullable: true }
    task_id:            { type: string, format: uuid, nullable: true }
    created_at:         { type: string, format: date-time }
    updated_at:         { type: string, format: date-time, nullable: true }
    error_message:      { type: string, nullable: true }
```

### KBVideoListResponse

```yaml
KBVideoListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/KBVideoRecord' } }
    next_cursor: { type: string, nullable: true }
```

---

## 4. Operations（草案）

### 4.1 `POST /api/v1/svc/knowledge-bases/{kb_id}/videos/ingest`

- operationId: `ingestKnowledgeBaseVideo`
- Body: `IngestKnowledgeBaseVideoRequest`
- 成功：`202 + AsyncTask`（`task_type` 建议 `kb.video.ingest`）
- 错误：
  - `401`、`403`、`404`
  - `400`：source=object_upload 但缺 object_key
  - `422`：视频格式不支持 — 建议 `VIDEO_FORMAT_UNSUPPORTED`

### 4.2 `GET /api/v1/svc/knowledge-bases/{kb_id}/videos`

- operationId: `listKnowledgeBaseVideos`
- Query: `status`、`limit`、`cursor`
- 成功：`200 + KBVideoListResponse`

### 4.3 `GET /api/v1/svc/knowledge-bases/{kb_id}/videos/{video_id}`

- operationId: `getKnowledgeBaseVideo`
- 成功：`200 + KBVideoRecord`
- 错误：`404`

### 4.4 `GET /api/v1/svc/knowledge-bases/{kb_id}/videos/{video_id}/transcript`（3b）

- operationId: `getKnowledgeBaseVideoTranscript`
- 成功：`200` + `{ text, segments[] }`（大文件分片下载）
- 3a 可合并进 `getKnowledgeBaseVideo.transcript`

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 编辑者 | 说明 |
|---|---|---|---|
| 视频列表/详情 | 可用 | 可用 | GET |
| ingest 视频 | 不可用 | 可用 | POST；先 object upload 或选已有 doc |
| 查看转写/章节 | 可用 | 可用 | GET video |
| 跳转任务中心 | 可用 | 可用 | 长视频 AsyncTask |
| 跳转对象存储 | 跳转 | 跳转 | object_key 来源 |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","title":"培训录像","source":"object_upload","object_key":"'"$KEY"'"}' \
  "$BASE/api/v1/svc/knowledge-bases/$KB_ID/videos/ingest"

curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/knowledge-bases/$KB_ID/videos/$VIDEO_ID"
```

---

## 7. 评审检查清单

- [ ] ingest 与 Core objects/upload 分工清晰
- [ ] 长视频必须 202，禁止同步阻塞
- [ ] chapters 时间轴与 transcript 一致
- [ ] 合入后更新 `kb-video-intelligence.md`

---

## 相关文件

- `docs/console-modules/knowledge/kb-video-intelligence.md`
- `docs/console-modules/compute/storage/object-storage-upload.md`
- `docs/console-modules/alerts/async-task-center.md`
