# BOSS overview 域 Phase 0 GAP 摘要

> **日期**：2026-06-17 · **详文**：`docs/boss-modules/overview/*.md`（6 模块）

## 已声明 path（只读参考 · 非 BOSS 平台 aggregate）

| 路径 | operationId | 模块 | 用法 |
|---|---|---|---|
| `GET /api/v1/gpu-inventory` | `listGPUInventory` | gpu-pool-trend | 租户 JWT 范围 inventory |
| `GET /api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | gpu-pool-trend | 占用参考 |
| `GET /api/v1/metering/usage` | **未声明** | resource-capacity-trend | 单租户用量参考 |
| `GET /api/v1/svc/inference-services` | `listInferenceServices` | inference-ops-trend | 租户推理 list 参考 |
| `GET /api/v1/svc/knowledge-bases` | `listKnowledgeBases` | kb-ops-trend | 租户 KB list 参考 |
| `GET /api/v1/observability/alert-rules` | `listObservabilityAlertRules` | platform-alerts-pending | **规则** CRUD，非 firing 事件 |
| `GET /api/v1/tasks/{task_id}` | **未声明** | platform-alerts-pending | 单查；**无** list |

## Core 缺口（BOSS 平台 · TODO-YAML）

| 优先级 | 能力 | 涉及模块 |
|---|---|---|
| **P1** | 平台 resource/capacity aggregate | operations-overview、resource-capacity-trend |
| **P1** | 平台 GPU pool aggregate / Top N | gpu-pool-trend |
| **P1** | 平台 inference ops aggregate | inference-ops-trend |
| **P1** | 平台 KB ops aggregate | kb-ops-trend |
| **P1** | firing 告警 events list | platform-alerts-pending |
| **P1** | `GET /api/v1/tasks` 跨租户 list | platform-alerts-pending、job-history |
| **P1** | 平台 pending 事项 aggregate | platform-alerts-pending |

## Phase 0 已修正（2026-06-17）

| 项 | 修正 |
|---|---|
| `platform-alerts-pending` 错误示例 | 补 400 筛选参数示例 |

## 架构待决

- 平台 overview KPI 单一 aggregate path vs 多 resource 聚合
- firing 事件与 alert-rules 是否分表/分 API
- `group_by` 是否允许 `tenant_id`（BOSS 需要，Console 不需要）

门禁：详文满配（Core）；全 BOSS 审计见 `boss-phase0-gap-index.md`
