# Phase 3 — Agent Audit OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`ai-native/agent-audit.md`  
> **TASK**：`TASK-SVC-018` 子项（Tool Permission 详化 ✅ 之后）  
> **依赖**：Agent Session、Tool Permission 运行时产生事件

---

## 1. 与平台审计 / 实例安全事件的边界

| 维度 | 平台审计 | 实例 security-events | Agent 行为审计（本草案） |
|---|---|---|---|
| 模块 | `security/audit-log.md` | `sandbox-templates.md` | `agent-audit.md` |
| 路径（规划） | Core observability audit | Core `/instances/{id}/security-events` | `/svc/agent/audit-events` |
| 事件来源 | 登录、Admin API | Sandbox 运行时安全 | tool_call、LLM 摘要、策略命中 |
| 读者 | 安全合规 | 算力运维 | Agent 运维 / 租户管理员 |
| PII | 合规脱敏 | 安全告警 | **必须**脱敏 prompt/response 正文 |

Console 不得把三类事件合并为单一 list API。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [AgentAudit]
x-ani-rbac-scope:
  read: scope:agent-audit:read
  # 审计只读；写入由运行时异步产生，无 Console POST
```

---

## 3. Schemas（草案）

### AgentAuditEvent

```yaml
AgentAuditEvent:
  type: object
  required: [id, event_type, occurred_at]
  properties:
    id:           { type: string, format: uuid }
    session_id:   { type: string, format: uuid, nullable: true }
    event_type:   { type: string, enum: [tool_call, tool_result, llm_request, llm_response, policy_hit, policy_denied, error] }
    tool_id:      { type: string, nullable: true }
    summary:      { type: string, description: 人类可读一行摘要；不含完整 prompt }
    actor_user_id: { type: string, nullable: true }
    severity:     { type: string, enum: [info, warning, critical], default: info }
    occurred_at:  { type: string, format: date-time }
    metadata:     { type: object, additionalProperties: true, description: 扩展字段；敏感键须后端 redact }
```

### AgentAuditEventListResponse

```yaml
AgentAuditEventListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/AgentAuditEvent' } }
    next_cursor: { type: string, nullable: true }
```

> **禁止**在 list 响应返回完整 `prompt`/`completion` 明文；详情如需扩展走 3b `GET .../audit-events/{id}` 且仍须 redact。

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/svc/agent/audit-events`

- operationId: `listAgentAuditEvents`
- Query:
  - `session_id`（uuid，可选）
  - `event_type`（enum，可选）
  - `tool_id`（可选）
  - `severity`（可选）
  - `occurred_after` / `occurred_before`（date-time，可选）
  - `limit`（default 50，max 200）、`cursor`
- 成功：`200 + AgentAuditEventListResponse`
- 错误：`401`、`403`

### 4.2 `GET /api/v1/svc/agent/sessions/{session_id}/audit-events`

- operationId: `listAgentSessionAuditEvents`
- 说明：会话详情 Tab 专用；等价于全局 list + `session_id` 过滤，但路径语义更清晰
- Query: `event_type`、`limit`、`cursor`
- 成功：`200 + AgentAuditEventListResponse`
- 错误：`404`（session 不存在或不可见）

### 4.3 `GET /api/v1/svc/agent/audit-events/{event_id}`（3b · 可选）

- operationId: `getAgentAuditEvent`
- 成功：`200 + AgentAuditEvent`（含 redacted metadata）
- 错误：`404`

### 4.4 导出（3b · 可选 · BOSS 边界）

- `POST /api/v1/svc/agent/audit-events/export` → `202 + AsyncTask`
- Console 租户页默认**不暴露**；若开放须单独 PRD

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 租户 Agent 管理员 | 说明 |
|---|---|---|---|
| 全局审计列表 | 可用 | 可用 | GET audit-events |
| 会话内审计 | 可用 | 可用 | GET sessions/{id}/audit-events |
| 单条详情 | 可用 | 可用 | GET event（3b） |
| 导出 | 不可用 | 待定 | 3b / BOSS |
| 删除/篡改 | 不可用 | 不可用 | 审计 append-only |

---

## 6. 与 Tool Permission / Session 联动

| 运行时行为 | 审计 event_type |
|---|---|
| 工具调用发起 | `tool_call` |
| 工具返回 | `tool_result` |
| 权限拒绝 | `policy_denied`（关联 `agent-tool-permission.md`） |
| Prompt 防护命中 | `policy_hit`（关联 `prompt-injection-guard.md`） |
| LLM 往返 | `llm_request` / `llm_response`（仅存 summary） |

会话详情页：消息时间线（`agent-session.md`）+ 审计 Tab（本模块）。

---

## 7. Handler 验收 curl（合入 YAML 后）

```bash
# 全局列表
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/audit-events?limit=50&event_type=tool_call"

# 会话 scoped
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/sessions/$SESSION_ID/audit-events?limit=50"

# 时间窗
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/audit-events?occurred_after=2026-06-01T00:00:00Z"
```

### 负向

| 用例 | 期望 |
|---|---|
| 无 agent-audit:read | 403 |
| 伪造 session_id | 404 |
| 响应含未 redact 的 prompt 全文 | **验收失败** |

---

## 8. 评审检查清单

- [ ] 与 Core `security/audit-log` 路径无冲突
- [ ] list 响应无敏感明文
- [ ] `event_type` enum 覆盖 tool + policy + llm
- [ ] session scoped list 与全局 list schema 一致
- [ ] 合入后更新 `agent-audit.md` §接口冻结规则

---

## 相关文件

- `docs/console-modules/ai-native/agent-audit.md`
- `docs/console-modules/ai-native/agent-session.md`
- `docs/console-modules/security/audit-log.md`
- `tasks/modules/spec/console/ai-native/spec-console-agent-audit.md`
