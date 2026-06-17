# SPEC: BOSS alert-rules

> Technical specification derived from: `tasks/modules/prd/boss/health/prd-boss-alert-rules.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **运维与可观测 / 告警规则** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/health/prd-boss-alert-rules.md`
- 主维护源: `docs/boss-modules/health/alert-rules.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文）

| 方法 | 路径 | operationId | 成功响应 | RBAC | BOSS 用法 |
|---|---|---|---|---|---|
| GET | `/api/v1/observability/alert-rules` | `listObservabilityAlertRules` | `200` | `scope:observability:read` | 租户上下文 CRUD 已声明；BOSS 跨租户 list **待 YAML** |
| POST | `/api/v1/observability/alert-rules` | `createObservabilityAlertRule` | `201` | `scope:observability:write` | 同上 |
| GET | `/api/v1/observability/alert-rules/{rule_id}` | `getObservabilityAlertRule` | `200` | `scope:observability:read` | 同上 |
| PATCH | `/api/v1/observability/alert-rules/{rule_id}` | `updateObservabilityAlertRule` | `200` | `scope:observability:write` | 同上 |
| DELETE | `/api/v1/observability/alert-rules/{rule_id}` | `deleteObservabilityAlertRule` | `200` | `scope:observability:write` | 同上 |

**Schema**：`ObservabilityAlertRule` — severity `info|warning|critical`；state `active|disabled|deleted`；POST/PATCH 须 `idempotency_key`。

### 2.3 Non-Frozen Capabilities

- 活动告警 firing events list/ack/silence — **TODO-YAML** P1
- BOSS 平台跨租户 rules 筛选 query — 待 RBAC 设计

### 2.4 Core 同步边界（2026-06-17）

- `CreateObservabilityAlertRuleRequest` **不含 `tenant_id`** — 创建时 `tenant_id` 由 JWT claims 写入
- BOSS 为其他租户建规则须 **platform API / TODO-YAML**，禁止切换租户 JWT 轮询
- 跨租户 rules list 筛选 — **TODO-YAML**

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。

## 3. Page Scope

### 3.1 Page Responsibilities

- 平台告警规则管理与当前告警处置页，承接 P0/P1/P2 告警确认、静默、分派与关闭。

### 3.2 Non-Goals

- 不定义 Console 租户自助契约
- 不把 YAML 未声明路径写成已实现
- 交付/白牌页不自造 REST schema

## 4. Architecture Constraints（ANI-14）

| 约束 | 要求 |
|---|---|
| Core 前缀 | `/api/v1/*` |
| Services 前缀 | `/api/v1/svc/*` |
| 租户边界 | JWT claims；不信任 body/query `tenant_id` |
| 写操作 | `idempotency_key` + 明确成功码 |
| 错误结构 | `code` / `message` / `request_id` |
| 422 | 仅 YAML 已声明 operation 可写冻结 |

## 5. 创建前置条件与操作矩阵

见主维护源 `docs/boss-modules/health/alert-rules.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台级 BOSS API — **TODO-YAML**（交付页：**N/A**）
- Phase：**P1**

## 7. 主维护源与流程

- `docs/boss-modules/health/alert-rules.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
