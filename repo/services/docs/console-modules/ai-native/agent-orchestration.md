# AI 原生 — Agent 编排平台

## 页面定位

多 Agent **工作流编排**（DAG、触发器、Human-in-the-loop）的设计与运行监控页。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md`。

## 文档管理规则

- 本文是 **Agent 编排平台** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- PRD/SPEC：`tasks/modules/prd/console/ai-native/prd-console-agent-orchestration.md`、`tasks/modules/spec/console/ai-native/spec-console-agent-orchestration.md`
- TASK：`TASK-SVC-018`

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/agent/workflows` | `listAgentWorkflows` | 3a |
| POST | `/api/v1/svc/agent/workflows` | `createAgentWorkflow` | 3a |
| GET | `/api/v1/svc/agent/workflows/{workflow_id}` | `getAgentWorkflow` | 3a |
| PATCH | `/api/v1/svc/agent/workflows/{workflow_id}` | `updateAgentWorkflow` | 3a |
| POST | `/api/v1/svc/agent/workflows/{workflow_id}/runs` | `startAgentWorkflowRun` | 3a |
| GET | `/api/v1/svc/agent/workflows/{workflow_id}/runs` | `listAgentWorkflowRuns` | 3a |
| GET | `/api/v1/svc/agent/workflow-runs/{run_id}` | `getAgentWorkflowRun` | 3b |

Schema（草案）：`AgentWorkflow`、`AgentWorkflowRun`、`CreateAgentWorkflowRequest`、`StartAgentWorkflowRunRequest`。

RBAC（草案）：`scope:agent-workflows:read|write|run`。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读/写/运行权限 | 对应 scope | `403` |
| 写操作幂等键 | POST/PATCH 带 `idempotency_key` | `400` |
| workflow 存在 | 有效 id | `404` |
| 启动 run | workflow `status=active` | `422`（建议 `WORKFLOW_NOT_ACTIVE`） |

## 页面职责

- 工作流列表/编辑器（definition DSL Phase 3b 细化）
- 触发器配置：manual / schedule / webhook
- 运行监控：run 状态、关联 `session_id`、`task_id`（跳转任务中心）
- Human-in-the-loop 审批（3b）

## 操作可用性矩阵

| 操作 | 只读用户 | 编排管理员 | 说明 |
|---|---|---|---|
| 查看 workflow/run | 可用 | 可用 | GET |
| 创建/编辑 | 不可用 | 可用 | POST/PATCH |
| 手动运行 | 不可用 | 可用 | POST runs（需 run scope） |
| 审批/取消 run | 不可用 | 可用 | 3b |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md` §4。

### `POST .../workflows/{id}/runs`（规划）

- 成功：`202 + AgentWorkflowRun`（含 `task_id`）
- 错误：`422`（非 active）

## 待补边界

- `definition` DAG schema — Phase 3b ADR
- 可视化 vs YAML DSL — Console 实现，API 存 JSON object
- 与 Gateway 流式 — 编排节点类型扩展

## 相关模块

- `ai-native/agent-session.md`
- `alerts/async-task-center.md`
- `ai-native/agent-audit.md`

## 验收标准

- [ ] run 202 与 Core AsyncTask 对齐
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
