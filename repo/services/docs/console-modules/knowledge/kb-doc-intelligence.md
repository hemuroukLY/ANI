# 知识库 — 文档智能

## 页面定位

对已上传文档执行 **结构化解析、摘要、实体抽取、分块优化** 的配置与结果查看页，扩展 `knowledge-base.md` 文档 CRUD。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md`。

## 文档管理规则

- 本文是 **知识库 — 文档智能** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 不得把草案路径写成**已实现**
- PRD/SPEC：`tasks/modules/prd/console/knowledge/prd-console-kb-doc-intelligence.md`、`tasks/modules/spec/console/knowledge/spec-console-kb-doc-intelligence.md`
- TASK：`TASK-SVC-018` 子项（知识库智能）

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| POST | `/api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}/analyze` | `analyzeKnowledgeBaseDocument` | 3a |
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/documents/{doc_id}/intelligence` | `getKnowledgeBaseDocumentIntelligence` | 3a |
| GET | `.../intelligence/history` | `listKnowledgeBaseDocumentIntelligenceRuns` | 3b |

Schema（草案）：`AnalyzeKnowledgeBaseDocumentRequest`、`KBDocumentIntelligence`。

RBAC（草案）：`scope:knowledge-bases:read`、`scope:knowledge-bases:write`（与 documents 子资源一致）。

异步：`202 + AsyncTask`，`task_type` 建议 `kb.document.analyze`；轮询 Core `GET /tasks/{task_id}`。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 读/写权限 | knowledge-bases read/write | `403 FORBIDDEN` |
| analyze 幂等键 | POST 携带 `idempotency_key` | `400` |
| kb / doc 存在 | 有效 `kb_id`、`doc_id` | `404 NOT_FOUND` |
| 文档已索引 | upload/index 完成 | `422`（建议 `DOCUMENT_NOT_INDEXED`） |
| 分析未重复跑 | 无进行中的 analyze | `422`（建议 `ANALYSIS_ALREADY_RUNNING`） |

## 页面职责

- 文档详情侧栏：分析状态、摘要、实体列表
- 触发分析：选择 features（summary / entities / chunk_optimization）
- 进行中跳转 **异步任务中心**；完成后可选跳 **向量写入** re-index
- 不把 upload `202` 写成「智能分析已完成」

## 操作可用性矩阵

| 操作 | 只读用户 | 编辑者 | YAML 合入后 |
|---|---|---|---|
| 查看分析结果 | 可用 | 可用 | GET intelligence |
| 触发分析 | 不可用 | 可用 | POST analyze |
| 查看历史 run | 不可用 | 可用 | GET history（3b） |
| 跳转任务中心 | 可用 | 可用 | Core tasks |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md` §4。

### `POST .../documents/{doc_id}/analyze`（规划）

- 成功：`202 + AsyncTask`
- 错误：`400`、`404`、`422`

### `GET .../intelligence`（规划）

- 成功：`200 + KBDocumentIntelligence`
- 错误：`401`、`403`、`404`

## 待补边界

- 无 intelligence 时返回 `200 pending` vs `404` — 评审二选一
- chunk 优化后是否自动触发 `vector-store-write` — Services 实现说明
- 与 upload 侧 `kb.parse` / `kb.index` 任务去重策略

## 相关模块

- `knowledge/knowledge-base.md`
- `compute/storage/vector-store-write.md`
- `alerts/async-task-center.md`

## 验收标准

- [ ] analyze 202 与 AsyncTask 对齐
- [ ] 与 documents CRUD 无 path 冲突
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
