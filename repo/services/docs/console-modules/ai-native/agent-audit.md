# AI 原生 — Agent 行为审计

## 页面定位

记录并检索 Agent **工具调用、LLM 请求摘要、策略命中/拒绝** 的租户侧审计视图；不等同于平台级 `security/audit-log.md` 或 Core 实例 `security-events`。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md`。

## 文档管理规则

- 本文是 **Agent 行为审计** 的主维护源
- **规划权威源**（合入前）：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 审计事件 **append-only**；Console 不提供 DELETE
- PRD/SPEC：`tasks/modules/prd/console/ai-native/prd-console-agent-audit.md`、`tasks/modules/spec/console/ai-native/spec-console-agent-audit.md`
- TASK：`TASK-SVC-018` 子项（Agent Audit）

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/agent/audit-events` | `listAgentAuditEvents` | 3a |
| GET | `/api/v1/svc/agent/sessions/{session_id}/audit-events` | `listAgentSessionAuditEvents` | 3a |
| GET | `/api/v1/svc/agent/audit-events/{event_id}` | `getAgentAuditEvent` | 3b |

Schema（草案）：`AgentAuditEvent`、`AgentAuditEventListResponse`。

RBAC（草案）：`scope:agent-audit:read`（只读；写入由 Agent 运行时异步产生）。

<!-- TODO-YAML: 合入后删除「规划」标记 -->

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 审计读权限 | `scope:agent-audit:read` | `403 FORBIDDEN` |
| session_id（scoped list） | 会话存在且租户可见 | `404 NOT_FOUND` |

本模块无 Console 写 API。**当前 YAML 未声明 `422`**。

## 页面职责

- **全局审计**页：时间线 + 筛选（event_type、tool_id、severity、时间窗）
- **会话详情 · 审计 Tab**：`listAgentSessionAuditEvents`
- 展示 `summary` 一行摘要；**不展示**完整 prompt/completion 明文
- 跳转：`agent-session.md`、`agent-tool-permission.md`（policy_denied 上下文）

## 页面结构

```text
Agent 行为审计
├── 筛选区（类型 / 工具 / 严重度 / 时间）
├── 事件时间线
│   ├── tool_call / tool_result
│   ├── llm_request / llm_response（摘要）
│   └── policy_hit / policy_denied / error
└── 会话内 Tab（复用 scoped API）
```

## 操作可用性矩阵

| 操作 | 只读用户 | Agent 管理员 | 说明 |
|---|---|---|---|
| 全局列表 | 可用 | 可用 | GET audit-events |
| 会话 scoped 列表 | 可用 | 可用 | GET sessions/.../audit-events |
| 单条详情 | 可用 | 可用 | GET event（3b） |
| 导出 | 不可用 | 待定 | 3b；BOSS 边界 |
| 删除/修改 | 不可用 | 不可用 | append-only |

## 接口冻结规则

> **当前为规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md` §4。

### `GET /api/v1/svc/agent/audit-events`（规划）

- Query：`session_id`、`event_type`、`tool_id`、`severity`、`occurred_after`、`occurred_before`、`limit`、`cursor`
- 成功：`200 + AgentAuditEventListResponse`
- 错误：`401`、`403`
- 空列表：`200` + `items: []`

### `GET /api/v1/svc/agent/sessions/{session_id}/audit-events`（规划）

- 成功：`200 + AgentAuditEventListResponse`
- 错误：`401`、`403`、`404`

### `GET /api/v1/svc/agent/audit-events/{event_id}`（规划 · 3b）

- 成功：`200 + AgentAuditEvent`（redacted metadata）
- 错误：`404`

## 待补边界

- 导出 / 长期归档 — BOSS 或 `billing-export` 类报表边界
- PII 字段级 redact 规则 — 安全评审后写入 YAML description
- 与 Core `security/audit-log` 联合检索 — **不合并 API**；Console 仅做导航分流
- 实例 `security-events` — 算力 Sandbox 域，见 `sandbox-templates.md`

## 相关模块

- `ai-native/agent-session.md` — 会话上下文
- `ai-native/agent-tool-permission.md` — policy_denied 来源
- `ai-native/prompt-injection-guard.md` — policy_hit 来源
- `security/audit-log.md` — 平台审计分流

## 验收标准

- [ ] list 响应不含未 redact 的 prompt 全文
- [ ] 与平台 audit-log 路径不冲突
- [ ] 会话 Tab 与全局 list 字段口径一致
- [ ] 合入 YAML 后切换为冻结口径
- [ ] PRD/SPEC 已同步
