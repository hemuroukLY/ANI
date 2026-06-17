# SPEC: Console ai-native-sandbox-security

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/ai-native/prd-console-ai-native-sandbox-security.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md`  
> Revised: 2026-06-17

## 1. Summary

Agent 运行时 Sandbox **安全配置**（规划）；与 Core `kind=sandbox` 实例分域。禁止 Services `/sandboxes*`。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- Services: `ANI-main/repo/api/openapi/services/v1.yaml`
- Core 只读: `ANI-main/repo/api/openapi/v1.yaml`（instances、security-events）

### 2.2 Planned Paths（Services）

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/svc/agent/sandbox-profiles` | `listAgentSandboxProfiles` | 200 | `scope:agent-sandbox:read` |
| GET | `.../{profile_id}` | `getAgentSandboxProfile` | 200 | read |
| PUT | `.../{profile_id}` | `updateAgentSandboxProfile` | 200 | write |
| GET | `/api/v1/svc/agent/sandbox-security/summary` | `getAgentSandboxSecuritySummary` | 200 | read（3b） |

### 2.3 Planned Schemas

- `AgentSandboxProfile`（`isolation_level`: standard | strict | locked）
- `AgentSandboxEgressRule`
- `AgentSandboxProfileUpdateRequest`
- `AgentSandboxSecuritySummary`（3b）

### 2.4 Core 只读引用

| Method | Path | 用途 |
|---|---|---|
| GET | `/api/v1/instances?kind=sandbox` | 关联实例展示 |
| GET | `/api/v1/instances/{id}/security-events` | 跳转安全事件 |

## 3. Page Scope

- Profile 列表/编辑、egress 规则、安全摘要、跳转算力 Sandbox
- **Non-Goals**：Core 实例 CRUD；deprecated Services sandboxes

## 4. 创建前置条件

无效 egress → `422 INVALID_EGRESS_RULE`。default profile 不可 DELETE。见主维护源。

## 5. Handler 验收（合入 YAML 后）

```bash
curl .../agent/sandbox-profiles
curl -X PUT .../agent/sandbox-profiles/$ID -d '{"idempotency_key":"...","isolation_level":"strict"}'
curl ".../instances?kind=sandbox"   # Core 只读
```

## 6. 主维护源

- `docs/console-modules/ai-native/ai-native-sandbox-security.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
