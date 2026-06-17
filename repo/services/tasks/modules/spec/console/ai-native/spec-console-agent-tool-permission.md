# SPEC: Console agent-tool-permission

> Phase 3 规划详化  
> Source: `tasks/modules/prd/console/ai-native/prd-console-agent-tool-permission.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md`  
> Revised: 2026-06-17

## 1. Summary

Agent 工具目录（只读）与授权矩阵（GET/PUT）；与 tenant RBAC 分资源。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/svc/agent/tools` | `listAgentTools` | 200 | `scope:agent-tools:read` |
| GET | `/api/v1/svc/agent/tool-permissions` | `getAgentToolPermissions` | 200 | read |
| PUT | `/api/v1/svc/agent/tool-permissions` | `updateAgentToolPermissions` | 200 | write |
| GET | `/api/v1/svc/agent/tool-permissions/effective` | `getEffectiveAgentToolPermissions` | 200 | read |

### 2.3 Planned Schemas

- `AgentToolDefinition`、`AgentToolDefinitionListResponse`
- `AgentToolGrant`（effect: allow|deny；scope: tenant|role|agent_profile）
- `AgentToolPermissionsSnapshot`、`AgentToolPermissionsUpdateRequest`

### 2.4 解析规则（草案）

- **deny 优先于 allow**
- PUT 全量替换 grants 快照

## 3. Page Scope

- 工具目录 Tab + 授权矩阵 + 保存
- **Non-Goals**：tenant roles PUT；MCP 市场写 API

## 4. 创建前置条件与操作矩阵

见主维护源。

## 5. Handler 验收（合入 YAML 后）

```bash
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/agent/tools"
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/agent/tool-permissions"
curl -X PUT .../agent/tool-permissions -d '{"idempotency_key":"...","grants":[...]}'
```

## 6. 主维护源

- `docs/console-modules/ai-native/agent-tool-permission.md`

OpenAPI **未合入** ≠ handler 已实现。
