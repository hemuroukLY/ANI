# SPEC: Console agent-session

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/ai-native/prd-console-agent-session.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-session-draft.md`  
> Revised: 2026-06-17

## 1. Summary

Agent 多轮会话 CRUD（规划）；与 `KnowledgeBaseSession` 分资源。合入 `services/v1.yaml` 前，本文与主维护源均标注**规划**。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/svc/agent/sessions` | `listAgentSessions` | 200 | `scope:agent-sessions:read` |
| POST | `/api/v1/svc/agent/sessions` | `createAgentSession` | 201 | `scope:agent-sessions:write` |
| GET | `/api/v1/svc/agent/sessions/{session_id}` | `getAgentSession` | 200 | read |
| PATCH | `/api/v1/svc/agent/sessions/{session_id}` | `updateAgentSession` | 200 | write |
| GET | `.../messages` | `listAgentSessionMessages` | 200 | read |

### 2.3 Planned Schemas

- `AgentSession`（`status`: active | archived）
- `AgentSessionListResponse`
- `CreateAgentSessionRequest` / `UpdateAgentSessionRequest`
- `AgentMessage` / `AgentMessageListResponse`（3b）

## 3. Page Scope

- 会话列表、创建、详情、归档、消息只读时间线
- **Non-Goals**：KB sessions；流式 Gateway；Boss 全平台 Agent 运营

## 4. 创建前置条件

见主维护源 §创建前置条件。`422` 的 `agent_profile_id` 口径待 YAML PreconditionFailed 声明。

## 5. 操作可用性矩阵

见主维护源 §操作可用性矩阵。

## 6. Handler 验收（合入 YAML 后）

```bash
curl -X POST .../agent/sessions -d '{"idempotency_key":"...","title":"demo"}'
curl .../agent/sessions?status=active
curl -X PATCH .../agent/sessions/$ID -d '{"idempotency_key":"...","status":"archived"}'
```

## 7. 主维护源

- `docs/console-modules/ai-native/agent-session.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-session-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
