# 知识库 — 来源引用

## 页面定位

展示知识库问答召回的**来源引用**列表（document_id、snippet、score），作为 `kb-qa-flow.md` 的只读明细页。

本页属于 **Services / KnowledgeBases** 子能力。

## 文档管理规则

- 本文是来源引用子页的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- inline 引用展示口径与 `kb-qa-flow.md` 必须一致
- TASK：`TASK-SVC-014`

## Services 层要求

<!-- ADDED-TO-YAML: Phase 2 2026-06-17 -->

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/svc/knowledge-bases/{kb_id}/citations` | `listKnowledgeBaseCitations` |

Query：`limit`、`cursor`（若 YAML 声明）；可选 `session_id`（以 OpenAPI 为准）。

响应：`KnowledgeBaseCitationListResponse`（`KnowledgeBaseCitation` items）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 知识库读权限 | readers 或 editors | `403 FORBIDDEN` |
| 知识库存在 | 当前租户可见 | `404 NOT_FOUND` |

本页无写接口。**当前 YAML 未声明 `422`**。

## 页面职责

- 按知识库展示引用列表；与问答页 citations 区块同源
- 支持按会话/时间筛选（仅当 YAML query 支持）
- 跳转对应文档详情（`documents/{doc_id}`，见 `knowledge-base.md`）

## 操作可用性矩阵

| 操作 | 只读用户 | 编辑者 | 管理员 |
|---|---|---|---|
| 查看引用列表 | 可用 | 可用 | 可用 |
| 导出引用 | 不可用 | 不可用 | 不可用（待补） |
| 编辑权限 | 不可用 | 见 `kb-permissions.md` | 可用 |

## 接口冻结规则

### `GET /api/v1/svc/knowledge-bases/{kb_id}/citations`

- 成功：`200 + KnowledgeBaseCitationListResponse`
- 错误：`401`、`403`、`404`
- 空列表：`200` + `items: []`（非 404）
- 租户边界：仅当前租户知识库

## 待补边界

- 引用详情单条 GET — **YAML 未声明** `{citation_id}`
- 引用导出/审计 — 待 Services 规划
- 与向量检索 debug 视图 — 非 Console 冻结能力

## 相关模块

- `kb-qa-flow.md` — 问答时 inline 展示引用
- `knowledge-base.md` — 文档 CRUD
- `kb-chat-history.md` — 会话上下文

## 验收标准

- [ ] 路径与 `services/v1.yaml` 一致
- [ ] handler stub ≠ 已实现
- [ ] 与 kb-qa-flow citations 字段口径一致
- [ ] 未自造 Citation 扩展 schema
