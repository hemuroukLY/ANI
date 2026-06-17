# Phase 3 — Agent Session OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`ai-native/agent-session.md`  
> **TASK**：`TASK-SVC-018` 首子项  
> **原则**：对齐现有 `services/v1.yaml` 风格；与 `KnowledgeBaseSession` **分资源**，禁止混用 path。

---

## 1. 与知识库 Session 的边界

| 维度 | KB Session | Agent Session |
|---|---|---|
| 路径 | `/knowledge-bases/{kb_id}/sessions` | `/agent/sessions` |
| 用途 | RAG 问答续聊 | Agent 工具调用 / 多轮推理 |
| 父资源 | KnowledgeBase | 租户级 Agent 域 |
| schema 前缀 | `KnowledgeBaseSession` | `AgentSession` |
| 消息存储 | query/stream 侧 | 独立 messages 子资源（本草案） |

Console：`kb-chat-history.md` ≠ `agent-session.md`。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [AgentSessions]
x-ani-rbac-scope:
  read:  scope:agent-sessions:read
  write: scope:agent-sessions:write
```

---

## 3. Schemas（草案）

### AgentSession

```yaml
AgentSession:
  type: object
  required: [id, status, created_at]
  properties:
    id:              { type: string, format: uuid }
    title:           { type: string, nullable: true }
    status:          { type: string, enum: [active, archived] }
    agent_profile_id: { type: string, nullable: true, description: 关联 Agent 配置/编排，Phase 3 可选 }
    last_message_at: { type: string, format: date-time, nullable: true }
    created_at:      { type: string, format: date-time }
    updated_at:      { type: string, format: date-time, nullable: true }
```

### AgentSessionListResponse

```yaml
AgentSessionListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/AgentSession' } }
    next_cursor: { type: string, nullable: true }
```

### CreateAgentSessionRequest

```yaml
CreateAgentSessionRequest:
  type: object
  required: [idempotency_key]
  properties:
    idempotency_key:  { type: string, format: uuid }
    title:            { type: string, nullable: true }
    agent_profile_id: { type: string, nullable: true }
```

### UpdateAgentSessionRequest

```yaml
UpdateAgentSessionRequest:
  type: object
  required: [idempotency_key]
  properties:
    idempotency_key: { type: string, format: uuid }
    title:           { type: string, nullable: true }
    status:          { type: string, enum: [active, archived] }
```

### AgentMessage（messages 子资源 · Phase 3b）

```yaml
AgentMessage:
  type: object
  required: [id, role, created_at]
  properties:
    id:         { type: string, format: uuid }
    role:       { type: string, enum: [user, assistant, tool, system] }
    content:    { type: string }
    tool_call_id: { type: string, nullable: true }
    created_at: { type: string, format: date-time }

AgentMessageListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/AgentMessage' } }
    next_cursor: { type: string, nullable: true }
```

> 流式推理入口（SSE/WebSocket）**不在本 CRUD 草案**；建议 Gateway 或独立 `POST .../completions` 扩展，见 `agent-orchestration.md`。

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/svc/agent/sessions`

- operationId: `listAgentSessions`
- Query: `status`（active|archived）、`limit`（default 50, max 200）、`cursor`
- 成功：`200 + AgentSessionListResponse`
- 错误：`401`、`403`

### 4.2 `POST /api/v1/svc/agent/sessions`

- operationId: `createAgentSession`
- Body: `CreateAgentSessionRequest`
- 成功：`201 + AgentSession`
- 错误：`400`、`401`、`403`、`409`（幂等冲突）
- 422：若绑定 `agent_profile_id` 不存在 — **建议语义** `AGENT_PROFILE_NOT_FOUND`（待 Services 冻结举例）

### 4.3 `GET /api/v1/svc/agent/sessions/{session_id}`

- operationId: `getAgentSession`
- 成功：`200 + AgentSession`
- 错误：`401`、`403`、`404`

### 4.4 `PATCH /api/v1/svc/agent/sessions/{session_id}`

- operationId: `updateAgentSession`
- Body: `UpdateAgentSessionRequest`
- 成功：`200 + AgentSession`
- 归档 archived 后会话只读 — Console 置灰发送入口
- 错误：`400`、`401`、`403`、`404`

### 4.5 `DELETE /api/v1/svc/agent/sessions/{session_id}`（Phase 3b · 可选）

- operationId: `deleteAgentSession`
- 成功：`204` 或 `200`（以评审为准）
- 硬删 vs 软删 — 产品待定；Console 默认走 archive

### 4.6 `GET /api/v1/svc/agent/sessions/{session_id}/messages`

- operationId: `listAgentSessionMessages`
- Query: `limit`、`cursor`
- 成功：`200 + AgentMessageListResponse`

---

## 5. Console 操作可用性矩阵（草案语义）

| 操作 | 只读用户 | Agent 编辑者 | 说明 |
|---|---|---|---|
| 列表/详情 | 可用 | 可用 | GET |
| 创建会话 | 不可用 | 可用 | POST |
| 归档 | 不可用 | 可用 | PATCH status=archived |
| 删除 | 不可用 | 不可用 | 可选 API；默认不暴露 |
| 查看消息 | 可用 | 可用 | GET messages |
| 发送消息 | 不可用 | 可用 | 流式入口待 Gateway |

---

## 6. Handler 验收 curl（合入 YAML 后执行）

```bash
# 创建
curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","title":"demo"}' \
  "$BASE/api/v1/svc/agent/sessions"

# 列表
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/agent/sessions?status=active&limit=20"

# 归档
curl -sS -X PATCH -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","status":"archived"}' \
  "$BASE/api/v1/svc/agent/sessions/$SESSION_ID"
```

---

## 7. 评审检查清单

- [ ] 与 `KnowledgeBaseSession` 无 path/schema 冲突
- [ ] 全部写操作含 `idempotency_key`
- [ ] RBAC scope 写入 YAML
- [ ] 2xx 响应挂完整 schema
- [ ] 合入后更新 `agent-session.md` §接口冻结规则（去掉「规划」标记）

---

## 相关文件

- `docs/console-modules/ai-native/agent-session.md`
- `docs/console-modules/openapi-drafts/phase3/openapi-phase3-domain-draft.md` §1
- `tasks/modules/spec/console/ai-native/spec-console-agent-session.md`
