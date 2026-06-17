# 运维与可观测（BOSS）

> **PRD/SPEC/HTML 同步**：tenant / health / metering 三域 21 模块已于 2026-06-17 与满配详文对齐。

平台 SRE/运维域：健康巡检、监控、日志、Trace、告警规则、任务历史、运维 Skills、故障处理。

| 页面 | 详文 | Console 对照 |
|---|---|---|
| 平台健康 | [platform-health.md](platform-health.md) | — |
| GPU 监控 | [gpu-monitoring.md](gpu-monitoring.md) | [gpu-management.md](../../console-modules/compute/gpu-management.md) |
| 推理监控 | [inference-monitoring.md](inference-monitoring.md) | 推理服务域 |
| 知识库监控 | [kb-monitoring.md](kb-monitoring.md) | 知识库域 |
| 日志 | [platform-logs.md](platform-logs.md) | — |
| Trace | [platform-trace.md](platform-trace.md) | — |
| 告警规则 | [alert-rules.md](alert-rules.md) | [alerts-pending-items.md](../../console-modules/alerts/alerts-pending-items.md) |
| 运维 Skills | [maint-skills.md](maint-skills.md) | — |
| 任务历史 | [job-history.md](job-history.md) | [async-task-center.md](../../console-modules/alerts/async-task-center.md) |
| 故障处理 | [incident-handling.md](incident-handling.md) | — |

## Core 已冻结（租户/平台只读参考）

| 能力 | 路径 |
|---|---|
| Liveness | `GET /api/v1/healthz` |
| PromQL 代理 | `GET /api/v1/observability/query` |
| GPU inventory | `GET /api/v1/gpu-inventory*` |
| 告警规则 CRUD | `/api/v1/observability/alert-rules*` |
| 单任务查询 | `GET /api/v1/tasks/{task_id}` |

## BOSS 待补（TODO-YAML）

- 平台服务健康聚合、跨租户日志/Trace list、活动告警事件 list、`GET /api/v1/tasks` list、运维 Skill 编排

## Phase 1 / 满配状态

**✅ 本域 10 模块均已满配（Core）文档**（2026-06-17）；C 层（maint-skills、incident-handling）主 API 仍为 **TODO-YAML P2**，实现依赖 YAML 合入。

| 层级 | 模块 |
|---|---|
| A 层 | platform-health、gpu-monitoring、alert-rules |
| B 层 | inference-monitoring、kb-monitoring、platform-logs、platform-trace、job-history |
| C 层 | maint-skills、incident-handling |

满配检查清单：[`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)。
