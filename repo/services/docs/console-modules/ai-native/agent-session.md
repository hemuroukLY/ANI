# AI 原生 — Agent 会话

## 页面定位

租户侧 **Agent 多轮会话** 列表、详情与生命周期管理页（创建、归档、删除），区别于知识库 `kb-chat-history.md` 的 RAG 问答会话。

本页属于 **Services / AgentSessions** 规划域（Phase 3）；OpenAPI **详化草案**见 `../openapi-drafts/phase3/openapi-phase3-agent-session-draft.md`。

## 文档管理规则

- 本文是 **Agent 会话** 的主维护源
- **规划权威源**（合入前）：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-session-draft.md`
- 合入后一级权威源切换为：`ANI-main/repo/api/openapi/services/v1.yaml`
- 不得把草案路径或 handler stub 写成**已实现**
- PRD/SPEC 同步：`tasks/modules/prd/console/ai-native/prd-console-agent-session.md`、`tasks/modules/spec/console/ai-native/spec-console-agent-session.md`
- TASK：`TASK-SVC-018`（Agent Session 为首子项）

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/agent/sessions` | `listAgentSessions` | 3a |
| POST | `/api/v1/svc/agent/sessions` | `createAgentSession` | 3a |
| GET | `/api/v1/svc/agent/sessions/{session_id}` | `getAgentSession` | 3a |
| PATCH | `/api/v1/svc/agent/sessions/{session_id}` | `updateAgentSession` | 3a |
| GET | `/api/v1/svc/agent/sessions/{session_id}/messages` | `listAgentSessionMessages` | 3b |
| DELETE | `/api/v1/svc/agent/sessions/{session_id}` | `deleteAgentSession` | 3b 可选 |

Schema（草案）：`AgentSession`、`AgentSessionListResponse`、`CreateAgentSessionRequest`、`UpdateAgentSessionRequest`、`AgentMessage`。

RBAC（草案）：`scope:agent-sessions:read`、`scope:agent-sessions:write`。

<!-- TODO-YAML: 上表路径合入 services/v1.yaml 后，本节改为「已冻结」并删除「规划」标记 -->

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 读/写权限 | `scope:agent-sessions:read` / `:write` | `403 FORBIDDEN` |
| 幂等键 | POST/PATCH 携带 `idempotency_key` | `400` |
| agent_profile_id | 若提供则 profile 存在 | `422`（建议语义 `AGENT_PROFILE_NOT_FOUND`；**待 YAML 声明 PreconditionFailed 后定稿**） |

## 页面职责

- 会话侧栏/列表：标题、状态、最近活跃时间
- 创建会话抽屉：标题、可选 Agent Profile
- 详情页：消息时间线（只读 list）；发送入口依赖 Gateway 流式能力（待编排模块）
- 归档操作：PATCH `status=archived`；归档后只读
- 跳转：`agent-audit.md`（工具调用轨迹）、`agent-tool-permission.md`

## 页面结构

```text
Agent 会话
├── 会话列表
│   ├── 筛选（active / archived）
│   └── 创建入口
├── 会话详情
│   ├── 消息时间线（listAgentSessionMessages）
│   ├── 工具调用摘要（跳转 agent-audit）
│   └── 归档 / 重命名
└── 空态 / YAML 未合入提示
```

## 操作可用性矩阵

| 操作 | 只读用户 | Agent 编辑者 | YAML 合入后 |
|---|---|---|---|
| 列表/详情 | 可用 | 可用 | GET |
| 创建 | 不可用 | 可用 | POST |
| 归档/重命名 | 不可用 | 可用 | PATCH |
| 查看消息 | 可用 | 可用 | GET messages |
| 硬删除 | 不可用 | 不可用 | DELETE 可选，默认 UI 不暴露 |
| 流式发送 | 不可用 | 待 Gateway | 非本 CRUD 草案 |

## 接口冻结规则

> **当前为规划草案**；合入 OpenAPI 后逐条改为已冻结事实。完整 operation 定义见 `../openapi-drafts/phase3/openapi-phase3-agent-session-draft.md` §4。

### `GET /api/v1/svc/agent/sessions`（规划）

- 成功：`200 + AgentSessionListResponse`
- 错误：`401`、`403`

### `POST /api/v1/svc/agent/sessions`（规划）

- 成功：`201 + AgentSession`
- 错误：`400`、`401`、`403`、`409`

### `GET /api/v1/svc/agent/sessions/{session_id}`（规划）

- 成功：`200 + AgentSession`
- 错误：`404`

### `PATCH /api/v1/svc/agent/sessions/{session_id}`（规划）

- 成功：`200 + AgentSession`
- 错误：`400`、`404`

## 待补边界

- 流式消息发送 / SSE — Gateway 或 `agent-orchestration.md`
- Agent Profile CRUD — 与 `agent-orchestration.md` 共享资源域
- 硬删除 vs 归档 — 产品默认 archive；DELETE 为 3b 可选项
- 与 KB `KnowledgeBaseSession` — **禁止**共用 path 或 schema

## 相关模块

- `ai-native/agent-audit.md` — 工具调用审计
- `ai-native/agent-tool-permission.md` — 工具授权
- `knowledge/kb-chat-history.md` — KB 问答会话（分资源）
- `../openapi-drafts/phase3/openapi-phase3-agent-session-draft.md` — OpenAPI 详化草案

## 验收标准

- [ ] 草案与主维护源 operation 一致
- [ ] 合入 YAML 前 UI 标注「规划态」
- [ ] 合入 YAML 后更新本节为冻结规则并移除规划标记
- [ ] 不与 KB sessions 混写契约
- [ ] PRD/SPEC/HTML 已同步
