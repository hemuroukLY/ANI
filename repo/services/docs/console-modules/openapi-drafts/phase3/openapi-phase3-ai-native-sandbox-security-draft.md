# Phase 3 — AI Native Sandbox Security OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`ai-native/ai-native-sandbox-security.md`  
> **TASK**：`TASK-SVC-018` 子项  
> **依赖**：与 Core `kind=sandbox` **只读联动**；策略与 Prompt Guard 互补

---

## 1. 与算力 Sandbox 实例的边界

| 维度 | 算力 Sandbox 实例 | AI 原生 Sandbox 安全（本草案） |
|---|---|---|
| 模块 | `sandbox-instance-management.md` | `ai-native-sandbox-security.md` |
| 路径 | Core `/api/v1/instances?kind=sandbox` | Services `/agent/sandbox-profiles` |
| 职责 | VM/容器级实例 CRUD、lifecycle | **Agent 代码执行**隔离级别、egress、资源上限 |
| 安全事件 | Core `.../security-events` | 聚合展示 + Agent 域策略命中 |
| 模板 | `sandbox-templates.md` | 不替代模板 list |

**禁止**在 Services 新建 `/sandboxes*`（deprecated）。Agent 运行时选用 `sandbox_profile_id` 映射到 Core 实例或轻量执行环境 — **映射规则待 Services 实现说明**。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [AgentSandbox]
x-ani-rbac-scope:
  read:  scope:agent-sandbox:read
  write: scope:agent-sandbox:write
```

---

## 3. Schemas（草案）

### AgentSandboxIsolationLevel

```yaml
enum: [standard, strict, locked]
# standard — 网络受限 + 基础 syscall 过滤
# strict   — 无 outbound 除白名单
# locked   — 只读 FS + 无 network
```

### AgentSandboxEgressRule

```yaml
AgentSandboxEgressRule:
  type: object
  required: [destination, action]
  properties:
    destination: { type: string, description: CIDR、域名或 service tag }
    action:      { type: string, enum: [allow, deny] }
    port:        { type: integer, nullable: true }
```

### AgentSandboxProfile

```yaml
AgentSandboxProfile:
  type: object
  required: [id, name, isolation_level, enabled, updated_at]
  properties:
    id:               { type: string, format: uuid }
    name:             { type: string }
    is_default:       { type: boolean }
    enabled:          { type: boolean }
    isolation_level:  { type: string, enum: [standard, strict, locked] }
    max_cpu_millis:   { type: integer, nullable: true }
    max_memory_mib:   { type: integer, nullable: true }
    max_duration_sec: { type: integer, nullable: true }
    egress_rules:     { type: array, items: { $ref: '#/components/schemas/AgentSandboxEgressRule' } }
    allow_tool_ids:   { type: array, items: { type: string }, nullable: true, description: 与 tool-permission 交集生效 }
    updated_at:       { type: string, format: date-time }
```

### AgentSandboxProfileUpdateRequest

```yaml
AgentSandboxProfileUpdateRequest:
  type: object
  required: [idempotency_key]
  properties:
    idempotency_key:  { type: string, format: uuid }
    enabled:          { type: boolean, nullable: true }
    isolation_level:  { type: string, nullable: true }
    max_cpu_millis:   { type: integer, nullable: true }
    max_memory_mib:   { type: integer, nullable: true }
    max_duration_sec: { type: integer, nullable: true }
    egress_rules:     { type: array, items: { $ref: '#/components/schemas/AgentSandboxEgressRule' }, nullable: true }
    allow_tool_ids:   { type: array, items: { type: string }, nullable: true }
```

### AgentSandboxSecuritySummary（只读聚合 · 3b）

```yaml
AgentSandboxSecuritySummary:
  type: object
  properties:
    active_runs:        { type: integer }
    blocked_egress_24h: { type: integer }
    policy_violations_24h: { type: integer }
    linked_instance_ids: { type: array, items: { type: string }, description: 可选 Core sandbox 实例摘要 }
```

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/svc/agent/sandbox-profiles`

- operationId: `listAgentSandboxProfiles`
- Query: `limit`、`cursor`
- 成功：`200 + { items: AgentSandboxProfile[], next_cursor }`

### 4.2 `GET /api/v1/svc/agent/sandbox-profiles/{profile_id}`

- operationId: `getAgentSandboxProfile`
- 成功：`200`；`404`

### 4.3 `PUT /api/v1/svc/agent/sandbox-profiles/{profile_id}`

- operationId: `updateAgentSandboxProfile`
- Body: `AgentSandboxProfileUpdateRequest`
- 成功：`200 + AgentSandboxProfile`
- default profile 不可 DELETE（无 DELETE op）
- 422：egress 规则冲突 — 建议 `INVALID_EGRESS_RULE`

### 4.4 `GET /api/v1/svc/agent/sandbox-security/summary`（3b）

- operationId: `getAgentSandboxSecuritySummary`
- 成功：`200 + AgentSandboxSecuritySummary`
- 说明：Console 仪表盘；细节事件跳转 Core `security-events` 或 `agent-audit`

### 4.5 Core 只读引用（非本模块新增）

- `GET /api/v1/instances?kind=sandbox` — 关联实例列表
- `GET /api/v1/instances/{id}/security-events` — 算力域事件时间线

---

## 5. 与 Prompt Guard 联动

| 场景 | Guard | Sandbox Profile |
|---|---|---|
| Prompt 注入 | scan 文本 | — |
| 恶意代码执行 | — | isolation + egress |
| 工具调用 | tool-permission | allow_tool_ids 交集 |

Audit：`policy_denied`（工具）vs `policy_hit`（guard）vs sandbox **egress block**（建议 event_type `sandbox_egress_denied` 在 agent-audit 扩展 enum — 评审时同步）

---

## 6. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 安全管理员 | 说明 |
|---|---|---|---|
| 查看 profile 列表/详情 | 可用 | 可用 | GET |
| 编辑 profile | 不可用 | 可用 | PUT |
| 安全摘要 | 可用 | 可用 | GET summary（3b） |
| 查看 Core 实例事件 | 跳转 | 跳转 | sandbox-templates |
| 创建算力 Sandbox | 跳转 | 跳转 | sandbox-instance-management |

---

## 7. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/sandbox-profiles"

curl -sS -X PUT -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{
    "idempotency_key":"'"$(uuidgen)"'",
    "isolation_level":"strict",
    "egress_rules":[{"destination":"10.0.0.0/8","action":"allow"}]
  }' \
  "$BASE/api/v1/svc/agent/sandbox-profiles/$PROFILE_ID"
```

---

## 8. 评审检查清单

- [ ] 无 Services `/sandboxes*` 路径
- [ ] 与 Core instances 边界在详文写清
- [ ] egress 与 tool-permission 交集语义固定
- [ ] 合入后更新 `ai-native-sandbox-security.md`

---

## 相关文件

- `docs/console-modules/ai-native/ai-native-sandbox-security.md`
- `docs/console-modules/compute/sandbox-instance-management.md`
- `docs/console-modules/openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md`
