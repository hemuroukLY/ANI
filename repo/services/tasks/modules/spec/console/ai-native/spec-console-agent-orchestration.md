# SPEC: Console agent-orchestration

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/ai-native/prd-console-agent-orchestration.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md`  
> Revised: 2026-06-17

## 1. Summary

多 Agent 工作流 CRUD + run 启动/查询（规划）；run 202 对齐 Core AsyncTask。合入前标注**规划**。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/svc/agent/workflows` | `listAgentWorkflows` | 200 | `scope:agent-workflows:read` |
| POST | `/api/v1/svc/agent/workflows` | `createAgentWorkflow` | 201 | write |
| GET | `.../{workflow_id}` | `getAgentWorkflow` | 200 | read |
| PATCH | `.../{workflow_id}` | `updateAgentWorkflow` | 200 | write |
| POST | `.../runs` | `startAgentWorkflowRun` | 202 | run |
| GET | `.../runs` | `listAgentWorkflowRuns` | 200 | read |
| GET | `/api/v1/svc/agent/workflow-runs/{run_id}` | `getAgentWorkflowRun` | 200 | read（3b） |

### 2.3 Planned Schemas

- `AgentWorkflow`（`status`: draft | active | disabled）
- `AgentWorkflowRun`（含 `task_id`、`session_id` 可选）
- `CreateAgentWorkflowRequest`、`UpdateAgentWorkflowRequest`
- `StartAgentWorkflowRunRequest`

## 3. Page Scope

- 工作流列表/编辑、触发器、运行监控、Human-in-the-loop（3b）
- **Non-Goals**：definition DAG schema（3b ADR）；Gateway 流式

## 4. 创建前置条件

非 `active` workflow 启动 run → `422 WORKFLOW_NOT_ACTIVE`。见主维护源。

## 5. Handler 验收（合入 YAML 后）

```bash
curl -X POST .../agent/workflows -d '{"idempotency_key":"...","name":"demo"}'
curl -X POST .../agent/workflows/$ID/runs -d '{"idempotency_key":"..."}'
curl .../tasks/$TASK_ID   # Core AsyncTask
```

## 6. 主维护源

- `docs/console-modules/ai-native/agent-orchestration.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
