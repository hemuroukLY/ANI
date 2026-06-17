# BOSS health 域 Phase 0 GAP 摘要

> **日期**：2026-06-17 · **详文**：`docs/boss-modules/health/*.md`（10 模块）

## 已声明 path（只读参考）

| 路径 | operationId | 模块 | 备注 |
|---|---|---|---|
| `GET /api/v1/healthz` | `liveness` | platform-health | 健康探针 |
| `GET /api/v1/observability/query` | `queryObservability` | gpu/inference/kb 监控 | PromQL 指标代理，非 logs/trace/firing list |
| `GET /api/v1/gpu-inventory` | `listGPUInventory` | gpu-monitoring | 租户 scope |
| `GET /api/v1/observability/alert-rules` | `listObservabilityAlertRules` | alert-rules | 规则 CRUD |
| `GET /api/v1/tasks/{task_id}` | **未声明** | job-history | 单查 only |
| `GET /api/v1/svc/inference-services` | `listInferenceServices` | inference-monitoring | 租户参考 |
| `GET /api/v1/svc/knowledge-bases` | `listKnowledgeBases` | kb-monitoring | 租户参考 |

## Core 缺口（BOSS 平台 · TODO-YAML）

| 优先级 | 能力 | 模块 |
|---|---|---|
| **P1** | `GET /api/v1/tasks` 跨租户 list | job-history、platform-alerts |
| **P1** | 平台 logs search/list | platform-logs |
| **P1** | 平台 trace search | platform-trace |
| **P1** | firing 告警 events（非 rules） | alert-rules |
| **P1** | 平台 GPU/inference/KB 跨租户 aggregate | gpu/inference/kb monitoring |
| **P2** | Skill 编排 API | maint-skills |
| **P2** | Incident CRUD / timeline | incident-handling |

## 约束（文档已贯彻）

- `AsyncTask.status` 用 `completed`，**禁** `succeeded`
- alert-rules list **≠** firing 事件数
- observability/query **≠** 资源 inventory list

## 架构待决

- Logs/Trace 统一 observability 前缀 vs 外部 Jaeger/Loki
- 平台 tasks list RBAC scope
- Incident 与 alert/job 关联模型

## 文档状态

health 域 **10 模块** 已满配（Core）+ 四材料已同步。
