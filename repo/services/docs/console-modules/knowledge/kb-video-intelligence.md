# 知识库 — 视频智能

## 页面定位

对视频素材执行 **转写、章节切分、关键帧/OCR** 并写入知识库的 Phase 3 页；长视频 ingest 异步。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md`。

## 文档管理规则

- 本文是 **知识库 — 视频智能** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- PRD/SPEC：`tasks/modules/prd/console/knowledge/prd-console-kb-video-intelligence.md`、`tasks/modules/spec/console/knowledge/spec-console-kb-video-intelligence.md`
- TASK：`TASK-SVC-018` 子项（知识库智能）

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| POST | `/api/v1/svc/knowledge-bases/{kb_id}/videos/ingest` | `ingestKnowledgeBaseVideo` | 3a |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/videos` | `listKnowledgeBaseVideos` | 3a |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/videos/{video_id}` | `getKnowledgeBaseVideo` | 3a |
| GET | `.../transcript` | `getKnowledgeBaseVideoTranscript` | 3b |

Schema（草案）：`IngestKnowledgeBaseVideoRequest`、`KBVideoRecord`、`KBVideoChapter`、`KBVideoListResponse`。

RBAC（草案）：`scope:knowledge-bases:read|write`。

Core 只读引用：`POST /api/v1/objects/upload`（上传字节）→ ingest 引用 `object_key`。

异步：`202 + AsyncTask`，`task_type` 建议 `kb.video.ingest`。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读/写权限 | knowledge-bases scope | `403` |
| ingest 幂等键 | POST 必填 | `400` |
| kb 存在 | 有效 `kb_id` | `404` |
| object_upload 源 | 提供 `object_key` | `400` |
| 视频格式 | 支持列表内 | `422`（建议 `VIDEO_FORMAT_UNSUPPORTED`） |

## 页面职责

- 视频列表与详情：转写、章节时间轴、关键帧预览
- ingest：先对象存储 upload 或选已有 document
- 长任务轮询 **异步任务中心**
- 与会议智能分入口（视频 vs 音频纪要）

## 操作可用性矩阵

| 操作 | 只读用户 | 编辑者 | 说明 |
|---|---|---|---|
| 列表/详情 | 可用 | 可用 | GET |
| ingest | 不可用 | 可用 | POST |
| 大 transcript 分片 | 可用 | 可用 | GET transcript（3b） |
| 跳转对象存储 | 跳转 | 跳转 | object_key |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md` §4。

### `POST .../videos/ingest`（规划）

- 成功：`202 + AsyncTask`
- 错误：`400`、`404`、`422`

### `GET .../videos/{video_id}`（规划）

- 成功：`200 + KBVideoRecord`

## 待补边界

- transcript 过大时的分页/分片 API（3b）
- keyframe 存储 CDN vs 对象存储 URL
- 与 `kb-meeting-intelligence` 音频-only 输入边界

## 相关模块

- `knowledge/kb-meeting-intelligence.md`
- `compute/storage/object-storage-upload.md`
- `alerts/async-task-center.md`

## 验收标准

- [ ] ingest 必须 202，禁止同步长视频阻塞
- [ ] 与 Core objects/upload 分工清晰
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
