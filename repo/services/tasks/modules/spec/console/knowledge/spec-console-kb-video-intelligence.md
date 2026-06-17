# SPEC: Console kb-video-intelligence

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/knowledge/prd-console-kb-video-intelligence.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md`  
> Revised: 2026-06-17

## 1. Summary

视频 ingest → 转写/章节/OCR（规划）；202 AsyncTask + videos list/get。Core objects/upload 只读引用。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- Services: `ANI-main/repo/api/openapi/services/v1.yaml`
- Core 只读: `ANI-main/repo/api/openapi/v1.yaml`（objects/upload）

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| POST | `.../videos/ingest` | `ingestKnowledgeBaseVideo` | 202 | write |
| GET | `.../videos` | `listKnowledgeBaseVideos` | 200 | read |
| GET | `.../videos/{video_id}` | `getKnowledgeBaseVideo` | 200 | read |
| GET | `.../videos/{video_id}/transcript` | `getKnowledgeBaseVideoTranscript` | 200 | read（3b） |

### 2.3 Planned Schemas

- `IngestKnowledgeBaseVideoRequest`（features: transcript|chapters|keyframes|ocr）
- `KBVideoRecord`、`KBVideoChapter`、`KBVideoListResponse`

### 2.4 AsyncTask

- `task_type`: `kb.video.ingest`

## 3. Page Scope

- ingest（object_key / url / document_id）、章节时间轴、转写展示
- **Non-Goals**：视频 CDN 播放契约

## 4. 创建前置条件

见主维护源。`VIDEO_FORMAT_UNSUPPORTED` 建议 422。

## 5. Handler 验收（合入 YAML 后）

```bash
curl -X POST .../knowledge-bases/$KB/videos/ingest -d '{"idempotency_key":"...","title":"培训","object_key":"..."}'
curl .../knowledge-bases/$KB/videos/$VID
```

## 6. 主维护源

- `docs/console-modules/knowledge/kb-video-intelligence.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
