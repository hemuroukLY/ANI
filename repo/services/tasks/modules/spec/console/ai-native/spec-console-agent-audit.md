# SPEC: Console agent-audit

> Phase 3 规划详化  
> Source: `tasks/modules/prd/console/ai-native/prd-console-agent-audit.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-audit-draft.md`  
> Revised: 2026-06-17

## 1. Summary

Agent 工具调用与策略事件只读审计；append-only；PII redact。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/svc/agent/audit-events` | `listAgentAuditEvents` | 200 | `scope:agent-audit:read` |
| GET | `/api/v1/svc/agent/sessions/{session_id}/audit-events` | `listAgentSessionAuditEvents` | 200 | read |
| GET | `/api/v1/svc/agent/audit-events/{event_id}` | `getAgentAuditEvent` | 200 | read |

### 2.3 Planned Schemas

- `AgentAuditEvent`（event_type: tool_call|tool_result|llm_request|llm_response|policy_hit|policy_denied|error）
- `AgentAuditEventListResponse`

### 2.4 安全约束

- list/detail **不得**返回完整 prompt/completion 明文
- `summary` 为人类可读一行；`metadata` 须 redact

## 3. Page Scope

- 全局审计页 + Session 详情审计 Tab
- **Non-Goals**：平台 audit-log；导出（默认）

## 4. 创建前置条件与操作矩阵

见主维护源。

## 5. Handler 验收（合入 YAML 后）

```bash
curl -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/audit-events?event_type=tool_call&limit=50"

curl -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/sessions/$SESSION_ID/audit-events"
```

负例：响应含未 redact prompt → 验收失败。

## 6. 主维护源

- `docs/console-modules/ai-native/agent-audit.md`

OpenAPI **未合入** ≠ handler 已实现。
