# 知识库模块索引

Console **知识库与 AI 应用**域子模块主维护源。主资源 `knowledge-base.md` 已冻结 CRUD；智能 ×3 为 Phase 3 **详化草案**（2026-06-17）。

## 模块清单

| 模块 | 详文 | OpenAPI | 说明 |
|---|---|---|---|
| 知识库管理 | [knowledge-base.md](./knowledge-base.md) | **已冻结** | CRUD + documents + query |
| 问答流程 | [kb-qa-flow.md](./kb-qa-flow.md) | **已冻结** | query / stream |
| 来源引用 | [kb-source-citation.md](./kb-source-citation.md) | Phase 2 已声明 | citations list |
| 对话历史 | [kb-chat-history.md](./kb-chat-history.md) | Phase 2 已声明 | sessions list |
| 知识库权限 | [kb-permissions.md](./kb-permissions.md) | Phase 2 已声明 | permissions PUT |
| 文档智能 | [kb-doc-intelligence.md](./kb-doc-intelligence.md) | **详化草案** | [openapi-phase3-kb-doc-intelligence-draft.md](../openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md) |
| 会议智能 | [kb-meeting-intelligence.md](./kb-meeting-intelligence.md) | **详化草案** | [openapi-phase3-kb-meeting-intelligence-draft.md](../openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md) |
| 视频智能 | [kb-video-intelligence.md](./kb-video-intelligence.md) | **详化草案** | [openapi-phase3-kb-video-intelligence-draft.md](../openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md) |

## 维护规则

- 路径前缀：`/api/v1/svc/knowledge-bases/{kb_id}/*`
- 智能 ingest/analyze 均 **202 + AsyncTask**；任务查询走 Core `GET /tasks/{task_id}`
- 评审前禁止将草案 path 写入「接口冻结规则」为已冻结事实
- TASK：Phase 2 扩展 → `TASK-SVC-014`；智能 ×3 → `TASK-SVC-018` 子项

## 评审建议顺序

1. 文档智能（依赖 documents 已 index）
2. 会议智能 / 视频智能（可并行；共享 AsyncTask 模式）

## 相关

- Backlog：`../governance/console-undefined-features-backlog.md` P2-08～P2-10
- 整域索引：`../openapi-drafts/phase3/openapi-phase3-domain-draft.md` §2
- 异步任务：`../alerts/async-task-center.md`
