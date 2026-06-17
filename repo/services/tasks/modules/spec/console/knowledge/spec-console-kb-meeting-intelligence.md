# SPEC: Console kb-meeting-intelligence

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/knowledge/prd-console-kb-meeting-intelligence.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md`  
> Revised: 2026-06-17

## 1. Summary

会议 ingest → 纪要/行动项（规划）；202 AsyncTask + meetings CRUD 只读。合入前标注**规划**。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| POST | `.../meetings/ingest` | `ingestKnowledgeBaseMeeting` | 202 | write |
| GET | `.../meetings` | `listKnowledgeBaseMeetings` | 200 | read |
| GET | `.../meetings/{meeting_id}` | `getKnowledgeBaseMeeting` | 200 | read |
| DELETE | `.../meetings/{meeting_id}` | `deleteKnowledgeBaseMeeting` | 204 | write（3b） |

### 2.3 Planned Schemas

- `IngestKnowledgeBaseMeetingRequest`（source: file|teams|zoom|manual_transcript）
- `KBMeetingRecord`、`KBMeetingListResponse`

### 2.4 AsyncTask

- `task_type`: `kb.meeting.ingest`

## 3. Page Scope

- ingest 表单、会议列表/详情、行动项、linked_document 跳转
- **Non-Goals**：Teams/Zoom OAuth（3b）

## 4. 创建前置条件

见主维护源 §创建前置条件。

## 5. Handler 验收（合入 YAML 后）

```bash
curl -X POST .../knowledge-bases/$KB/meetings/ingest -d '{"idempotency_key":"...","title":"周会","source":"file",...}'
curl .../knowledge-bases/$KB/meetings/$ID
```

## 6. 主维护源

- `docs/console-modules/knowledge/kb-meeting-intelligence.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
