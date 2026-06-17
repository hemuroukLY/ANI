# Phase 3 — Agent Orchestration OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`ai-native/agent-orchestration.md`  
> **TASK**：`TASK-SVC-018` 子项  
> **依赖**：Session、Tool Permission、Guard；运行态关联 Core `AsyncTask`

---

## 1. 与会话 / 任务中心的边界

| 维度 | Agent Session | Agent Orchestration | 异步任务中心 |
|---|---|---|---|
| 模块 | `agent-session.md` | **本文** | `async-task-center.md` |
| 资源 | 多轮对话 | **工作流定义 + 运行实例** | 任意 AsyncTask |
| 路径 | `/agent/sessions` | `/agent/workflows` | Core `/tasks/{id}` |
| 触发 | 用户聊天 | cron/webhook/手动 run | 202 异步操作 |

一次 **workflow run** 可创建/关联 `AgentSession`；长跑步骤返回 `operation_id` / `task_id` 供任务中心查询。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [AgentWorkflows]
x-ani-rbac-scope:
  read:  scope:agent-workflows:read
  write: scope:agent-workflows:write
  run:   scope:agent-workflows:run
```

---

## 3. Schemas（草案）

### AgentWorkflow

```yaml
AgentWorkflow:
  type: object
  required: [id, name, status, created_at]
  properties:
    id:          { type: string, format: uuid }
    name:        { type: string }
    description: { type: string, nullable: true }
    status:      { type: string, enum: [draft, active, archived] }
    trigger:     { $ref: '#/components/schemas/AgentWorkflowTrigger' }
    definition:  { type: object, additionalProperties: true, description: DAG/步骤 DSL；schema 待 Phase 3b 细化 }
    created_at:  { type: string, format: date-time }
    updated_at:  { type: string, format: date-time, nullable: true }
```

### AgentWorkflowTrigger

```yaml
AgentWorkflowTrigger:
  type: object
  required: [type]
  properties:
    type:  { type: string, enum: [manual, schedule, webhook] }
    cron:  { type: string, nullable: true }
    webhook_secret_ref: { type: string, nullable: true }
```

### AgentWorkflowRun

```yaml
AgentWorkflowRun:
  type: object
  required: [id, workflow_id, status, started_at]
  properties:
    id:          { type: string, format: uuid }
    workflow_id: { type: string, format: uuid }
    status:      { type: string, enum: [pending, running, completed, failed, cancelled] }
    session_id:  { type: string, format: uuid, nullable: true }
    task_id:     { type: string, format: uuid, nullable: true, description: Core AsyncTask }
    started_at:  { type: string, format: date-time }
    completed_at: { type: string, format: date-time, nullable: true }
    error_message: { type: string, nullable: true }
```

### CreateAgentWorkflowRequest

```yaml
CreateAgentWorkflowRequest:
  type: object
  required: [idempotency_key, name]
  properties:
    idempotency_key: { type: string, format: uuid }
    name:            { type: string }
    description:     { type: string, nullable: true }
    trigger:         { $ref: '#/components/schemas/AgentWorkflowTrigger' }
    definition:      { type: object, nullable: true }
```

### StartAgentWorkflowRunRequest

```yaml
StartAgentWorkflowRunRequest:
  type: object
  required: [idempotency_key]
  properties:
    idempotency_key: { type: string, format: uuid }
    input:           { type: object, nullable: true, description: 运行入参 }
    session_id:      { type: string, format: uuid, nullable: true }
```

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/svc/agent/workflows`

- operationId: `listAgentWorkflows`
- Query: `status`、`limit`、`cursor`
- 成功：`200 + { items: AgentWorkflow[], next_cursor }`

### 4.2 `POST /api/v1/svc/agent/workflows`

- operationId: `createAgentWorkflow`
- Body: `CreateAgentWorkflowRequest`
- 成功：`201 + AgentWorkflow`

### 4.3 `GET /api/v1/svc/agent/workflows/{workflow_id}`

- operationId: `getAgentWorkflow`
- 成功：`200 + AgentWorkflow`；`404`

### 4.4 `PATCH /api/v1/svc/agent/workflows/{workflow_id}`

- operationId: `updateAgentWorkflow`
- Body: `{ idempotency_key, name?, status?, trigger?, definition? }`
- 成功：`200 + AgentWorkflow`
- archived 后不可 start run

### 4.5 `POST /api/v1/svc/agent/workflows/{workflow_id}/runs`

- operationId: `startAgentWorkflowRun`
- Body: `StartAgentWorkflowRunRequest`
- 成功：`202 + AgentWorkflowRun`（含 `task_id`）
- 422：workflow 非 active — 建议 `WORKFLOW_NOT_ACTIVE`

### 4.6 `GET /api/v1/svc/agent/workflows/{workflow_id}/runs`

- operationId: `listAgentWorkflowRuns`
- Query: `status`、`limit`、`cursor`
- 成功：`200 + { items: AgentWorkflowRun[], next_cursor }`

### 4.7 `GET /api/v1/svc/agent/workflow-runs/{run_id}`（3b）

- operationId: `getAgentWorkflowRun`
- 成功：`200 + AgentWorkflowRun`

### 4.8 Human-in-the-loop（3b · 可选）

- `POST .../workflow-runs/{run_id}/approve` — 审批继续
- `POST .../workflow-runs/{run_id}/cancel` — 取消

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 编排管理员 | 说明 |
|---|---|---|---|
| 列表/详情 workflow | 可用 | 可用 | GET |
| 创建/编辑 | 不可用 | 可用 | POST/PATCH |
| 手动运行 | 不可用 | 可用（需 run scope） | POST runs |
| 查看 run 历史 | 可用 | 可用 | GET runs |
| 审批/取消 | 不可用 | 可用 | 3b |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","name":"daily-report","trigger":{"type":"manual"}}' \
  "$BASE/api/v1/svc/agent/workflows"

curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'"}' \
  "$BASE/api/v1/svc/agent/workflows/$WF_ID/runs"
# 期望 202 + run.task_id
```

---

## 7. 评审检查清单

- [ ] run 202 与 Core AsyncTask 字段对齐
- [ ] 与 `/agent/sessions` 无资源 ID 冲突
- [ ] DSL `definition` schema Phase 3b 单独 ADR
- [ ] 合入后更新 `agent-orchestration.md`

---

## 相关文件

- `docs/console-modules/ai-native/agent-orchestration.md`
- `docs/console-modules/ai-native/agent-session.md`
- `docs/console-modules/alerts/async-task-center.md`
