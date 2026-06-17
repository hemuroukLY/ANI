# 知识库 — 会议智能

## 页面定位

导入会议录音/转写，生成 **纪要、行动项、知识库条目** 的能力页；ingest 异步，结果可关联 `KBDocument`。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md`。

## 文档管理规则

- 本文是 **知识库 — 会议智能** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- PRD/SPEC：`tasks/modules/prd/console/knowledge/prd-console-kb-meeting-intelligence.md`、`tasks/modules/spec/console/knowledge/spec-console-kb-meeting-intelligence.md`
- TASK：`TASK-SVC-018` 子项（知识库智能）

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| POST | `/api/v1/svc/knowledge-bases/{kb_id}/meetings/ingest` | `ingestKnowledgeBaseMeeting` | 3a |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/meetings` | `listKnowledgeBaseMeetings` | 3a |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/meetings/{meeting_id}` | `getKnowledgeBaseMeeting` | 3a |
| DELETE | `.../meetings/{meeting_id}` | `deleteKnowledgeBaseMeeting` | 3b |

Schema（草案）：`IngestKnowledgeBaseMeetingRequest`、`KBMeetingRecord`、`KBMeetingListResponse`。

RBAC（草案）：`scope:knowledge-bases:read|write`。

异步：`202 + AsyncTask`，`task_type` 建议 `kb.meeting.ingest`。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读/写权限 | knowledge-bases scope | `403` |
| ingest 幂等键 | POST 必填 | `400` |
| kb 存在 | 有效 `kb_id` | `404` |
| source=file | 提供 `file_ref` | `400` |
| 外部会议可达 | teams/zoom 源 | `422`（建议 `MEETING_SOURCE_UNAVAILABLE`） |

## 页面职责

- 会议列表：标题、来源、状态、发生时间
- ingest 表单：文件 / 手动转写 / 外部会议 ID（Teams/Zoom 3b 需集成）
- 详情：纪要、行动项 checklist、`linked_document_id` 跳转
- 进行中任务跳转 **异步任务中心**

## 操作可用性矩阵

| 操作 | 只读用户 | 编辑者 | 说明 |
|---|---|---|---|
| 列表/详情 | 可用 | 可用 | GET |
| ingest | 不可用 | 可用 | POST |
| 删除 | 不可用 | 可用 | DELETE（3b） |
| 跳转 KB 文档 | 可用 | 可用 | linked_document_id |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md` §4。

### `POST .../meetings/ingest`（规划）

- 成功：`202 + AsyncTask`
- 错误：`400`、`404`、`422`

### `GET .../meetings/{meeting_id}`（规划）

- 成功：`200 + KBMeetingRecord`

## 待补边界

- Teams/Zoom OAuth — `integration-third-party.md`（3b）
- 与 `kb-doc-intelligence` 共享 NLP 管线
- DELETE 是否级联 linked_document — 评审确认

## 相关模块

- `knowledge/kb-doc-intelligence.md`
- `knowledge/knowledge-base.md`
- `integration/integration-third-party.md`

## 验收标准

- [ ] ingest 与 documents path 不混用
- [ ] AsyncTask resource 指向 meeting
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
