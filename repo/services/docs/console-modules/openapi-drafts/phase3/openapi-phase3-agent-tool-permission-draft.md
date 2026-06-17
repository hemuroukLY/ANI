# Phase 3 — Agent Tool Permission OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`ai-native/agent-tool-permission.md`  
> **TASK**：`TASK-SVC-018` 子项（Session 详化 ✅ 之后）  
> **依赖**：Agent Session 运行时鉴权消费本契约；工具目录只读可关联 `mcp-tool-market.md`（Phase 3b）

---

## 1. 与租户 RBAC / KB 权限的边界

| 维度 | 租户 RBAC | KB 权限 | Agent 工具权限（本草案） |
|---|---|---|---|
| 模块 | `role-permission-edit.md` | `kb-permissions.md` | `agent-tool-permission.md` |
| 路径 | `/tenant/roles/{id}` | `/knowledge-bases/{kb_id}/permissions` | `/agent/tool-permissions` |
| 控制对象 | API **scope**（如 `instances:read`） | KB readers/editors | Agent 可调用的 **工具/MCP** |
| 运行时 | Gateway / Auth 中间件 | KB query | Agent 执行器 tool dispatch |

**禁止**把 `scope:agent-tools:*` 与 `UpdateTenantRoleRequest.permissions[]` 混为同一 PUT 体。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [AgentTools]
x-ani-rbac-scope:
  read:  scope:agent-tools:read
  write: scope:agent-tools:write
```

---

## 3. Schemas（草案）

### AgentToolDefinition（目录只读 · 3a）

```yaml
AgentToolDefinition:
  type: object
  required: [id, name, provider]
  properties:
    id:          { type: string, description: 工具稳定 ID，如 mcp:filesystem/read }
    name:        { type: string }
    provider:    { type: string, enum: [builtin, mcp, custom], description: 来源 }
    risk_level:  { type: string, enum: [low, medium, high], nullable: true }
    description: { type: string, nullable: true }
    enabled:     { type: boolean, description: 租户是否已在目录中启用展示 }
```

### AgentToolDefinitionListResponse

```yaml
AgentToolDefinitionListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/AgentToolDefinition' } }
    next_cursor: { type: string, nullable: true }
```

### AgentToolGrant（授权条目）

```yaml
AgentToolGrant:
  type: object
  required: [tool_id, effect]
  properties:
    tool_id:    { type: string }
    effect:     { type: string, enum: [allow, deny] }
    scope:      { type: string, enum: [tenant, role, agent_profile], default: tenant }
    subject_id: { type: string, nullable: true, description: role_id 或 agent_profile_id；tenant 级可空 }
```

### AgentToolPermissionsSnapshot（GET 响应）

```yaml
AgentToolPermissionsSnapshot:
  type: object
  required: [grants, updated_at]
  properties:
    grants:     { type: array, items: { $ref: '#/components/schemas/AgentToolGrant' } }
    updated_at: { type: string, format: date-time }
    version:    { type: integer, description: 乐观锁版本，可选 }
```

### AgentToolPermissionsUpdateRequest

```yaml
AgentToolPermissionsUpdateRequest:
  type: object
  required: [idempotency_key, grants]
  properties:
    idempotency_key: { type: string, format: uuid }
    grants:          { type: array, items: { $ref: '#/components/schemas/AgentToolGrant' } }
    version:         { type: integer, nullable: true, description: 若启用乐观锁则必填 }
```

> **deny 优先于 allow**：运行时解析顺序须在 Services 实现说明中固定；Console 矩阵 UI 须展示冲突提示。

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/svc/agent/tools`

- operationId: `listAgentTools`
- 说明：租户可见**工具目录**（builtin + 已订阅 MCP）；不等于 MCP 市场写操作
- Query: `provider`、`risk_level`、`limit`、`cursor`
- 成功：`200 + AgentToolDefinitionListResponse`
- 错误：`401`、`403`

### 4.2 `GET /api/v1/svc/agent/tool-permissions`

- operationId: `getAgentToolPermissions`
- 说明：返回当前租户工具授权快照（全量 grants）
- 成功：`200 + AgentToolPermissionsSnapshot`
- 错误：`401`、`403`

### 4.3 `PUT /api/v1/svc/agent/tool-permissions`

- operationId: `updateAgentToolPermissions`
- Body: `AgentToolPermissionsUpdateRequest`
- 成功：`200 + AgentToolPermissionsSnapshot`
- 错误：`400`（空 grants 非法 tool_id）、`401`、`403`、`409`（幂等/版本冲突）
- 422：引用不存在的 `tool_id` 或 `subject_id` — **建议语义** `TOOL_NOT_FOUND` / `SUBJECT_NOT_FOUND`（待 PreconditionFailed 举例）

### 4.4 `GET /api/v1/svc/agent/tool-permissions/effective`（3b · 可选）

- operationId: `getEffectiveAgentToolPermissions`
- Query: `role_id`、`agent_profile_id`（可选，用于 Console 预览「某角色有效工具集」）
- 成功：`200 + AgentToolPermissionsSnapshot`（已解析 deny/allow）

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 租户 Agent 管理员 | 说明 |
|---|---|---|---|
| 查看工具目录 | 可用 | 可用 | GET tools |
| 查看授权矩阵 | 可用 | 可用 | GET tool-permissions |
| 更新授权 | 不可用 | 可用 | PUT |
| 有效权限预览 | 可用 | 可用 | GET effective（3b） |
| MCP 市场上架 | 不可用 | 跳转 | `mcp-tool-market.md` |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
# 工具目录
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/tools?limit=50"

# 当前授权
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/tool-permissions"

# 更新（allow 单工具）
curl -sS -X PUT -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","grants":[{"tool_id":"mcp:fs/read","effect":"allow","scope":"tenant"}]}' \
  "$BASE/api/v1/svc/agent/tool-permissions"
```

---

## 7. 评审检查清单

- [ ] 与 `UpdateTenantRoleRequest` / KB permissions 无 path 冲突
- [ ] PUT 必填 `idempotency_key`
- [ ] `AgentToolGrant.scope` enum 与产品粒度一致
- [ ] deny/allow 解析规则写入 Services 实现说明
- [ ] 合入后更新 `agent-tool-permission.md` §接口冻结规则

---

## 相关文件

- `docs/console-modules/ai-native/agent-tool-permission.md`
- `docs/console-modules/ai-native/agent-session.md`
- `docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-session-draft.md`
- `tasks/modules/spec/console/ai-native/spec-console-agent-tool-permission.md`
