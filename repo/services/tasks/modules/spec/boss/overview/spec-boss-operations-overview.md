# SPEC: BOSS operations-overview

> Technical specification derived from: `tasks/modules/prd/boss/overview/prd-boss-operations-overview.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **平台运营总览 / 运营总览** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/overview/prd-boss-operations-overview.md`
- 主维护源: `docs/boss-modules/overview/operations-overview.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文）

| 方法 | 路径 | operationId | BOSS 用法 |
|---|---|---|---|
| GET | `/api/v1/healthz` | `liveness` | 单组件参考；**非** 大盘 |
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL 摘要 |
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | 租户 **只读参考** |
| GET | `/api/v1/observability/alert-rules` | `listObservabilityAlertRules` | **规则**；非 firing |
| GET | `/api/v1/tasks/{task_id}` | — | **路径已声明**；无 list |
| GET | `/api/v1/metering/usage` | — | 租户 **只读参考** |

**本页无独立冻结 API** — 平台首页 aggregate **TODO-YAML**。

### 2.3 Non-Frozen Capabilities

- 平台级 `运营总览` 跨租户 aggregate API — **TODO-YAML** P1
- 平台 tenants / tasks list / firing alerts / metering aggregate — **TODO-YAML**
- BOSS 专属 RBAC scope — 待 Core/Services 确认

### 2.4 Core 同步边界（2026-06-17）

- **无** `GET /api/v1/boss/overview` — 页面层轻聚合
- 租户级已声明路径仅 **只读参考**，**非** BOSS 正式平台契约
- `GET /api/v1/tasks/{task_id}` — **路径已声明**；list — **TODO-YAML**
- `alert-rules` = 规则，**非** firing 事件计数
- 403 示例为 `permission denied` 参考；scope 名称 **待 YAML**

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。

## 3. Page Scope

### 3.1 Page Responsibilities

- BOSS 首页聚合页，面向平台管理员、SRE、运维、运营与交付人员，汇总全平台运行态、租户态、资源池态与风险态。

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

见主维护源 `docs/boss-modules/overview/operations-overview.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台级 BOSS API — **TODO-YAML**（交付页：**N/A**）
- Phase：**P1**

## 7. 主维护源与流程

- `docs/boss-modules/overview/operations-overview.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
