# Phase 3 知识库智能 ×3 — 文档验收记录

> **生成日期**：2026-06-17  
> **用途**：知识库智能三子模块 **文档层验收**（详文 + PRD/SPEC + OpenAPI 草案）。  
> **约束**：本阶段**不修改** `ANI-main/**`；YAML 合入与 handler 在 ANI-main PR 后回填。  
> **索引**：`docs/console-modules/knowledge/README.md` · **TASK**：`TASK-SVC-018` 子项

---

## 总表

| 模块 | 详文 | PRD/SPEC | OpenAPI 草案 | 文档验收 | YAML 合入 | Handler |
|---|---|---|---|---|---|---|
| 文档智能 | `kb-doc-intelligence.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md` | ✅ | ☐ | ☐ |
| 会议智能 | `kb-meeting-intelligence.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-kb-meeting-intelligence-draft.md` | ✅ | ☐ | ☐ |
| 视频智能 | `kb-video-intelligence.md` | ✅ | `../openapi-drafts/phase3/openapi-phase3-kb-video-intelligence-draft.md` | ✅ | ☐ | ☐ |

**文档验收 ✅ 定义**

- 详文含 Services 分层、创建前置条件、操作矩阵、规划口径接口冻结规则
- 异步操作统一 `202 + AsyncTask`；`task_type` 分别为 `kb.document.analyze`、`kb.meeting.ingest`、`kb.video.ingest`
- RBAC 与 knowledge-bases documents 一致：`scope:knowledge-bases:read|write`
- 视频 ingest 引用 Core `objects/upload`；不与 upload `202` 解析语义冲突

---

> **整域联评**：`tasks/phase3/PHASE3-JOINT-REVIEW-AGENDA.md`（场次 B：KB + 模型）

---

## 联评议程（域内摘要 · 完整版见整域议程 §4.2）

1. **AsyncTask 统一性**：三模块 task_type、resource_id 映射、async-task-center 展示
2. **文档智能 vs upload**：index 完成前置条件、`DOCUMENT_NOT_INDEXED` 422
3. **会议 vs 视频**：输入源、产出 schema、与 doc-intelligence 管线分工
4. **3a / 3b 切分**：history、DELETE meeting、transcript 分片 — 首版可隐藏 3b UI

---

## 文档层 curl 预览（合入 YAML 后）

```bash
curl -X POST .../knowledge-bases/$KB/documents/$DOC/analyze -d '{"idempotency_key":"..."}'
curl -X POST .../knowledge-bases/$KB/meetings/ingest -d '{"idempotency_key":"...","title":"周会","source":"file",...}'
curl -X POST .../knowledge-bases/$KB/videos/ingest -d '{"idempotency_key":"...","title":"培训","object_key":"..."}'
curl "$BASE/api/v1/tasks/$TASK_ID"
```

---

## 签核

- [x] ×3 详文 + PRD/SPEC + 草案（2026-06-17）
- [x] `knowledge/README.md`、`../openapi-drafts/phase3/openapi-phase3-domain-draft.md` §2 更新
- [ ] Services 架构联评通过 — **整域**：`PHASE3-JOINT-REVIEW-AGENDA.md` §7
- [ ] 分批合入 `services/v1.yaml`
- [ ] Handler 运行时验通（ANI-main）

---

## 相关文件

- `docs/console-modules/openapi-drafts/phase3/openapi-phase3-domain-draft.md`
- `tasks/phase3/acceptance/PHASE3-AI-NATIVE-ACCEPTANCE-RECORD.md`（AI 原生 ×7 同级验收）
- `tasks/execution/SERVICES-TEAM-TASKS.md` — TASK-SVC-018
