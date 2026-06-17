# SPEC: BOSS gpu-monitoring

> Technical specification derived from: `tasks/modules/prd/boss/health/prd-boss-gpu-monitoring.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **运维与可观测 / GPU 监控** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/health/prd-boss-gpu-monitoring.md`
- 主维护源: `docs/boss-modules/health/gpu-monitoring.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文）

| 方法 | 路径 | operationId | 成功响应 | RBAC | BOSS 用法 |
|---|---|---|---|---|---|
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | `200` | `scope:gpu-inventory:read` | 租户可见范围 inventory；BOSS 平台 aggregate **待 YAML** |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | `200` | `scope:gpu-inventory:read` | 占用摘要；同上 |
| GET | `/api/v1/observability/query` | `queryObservability` | `200` | `scope:observability:read` | PromQL；**不是** GPU 大盘 API |

**`GPUInventoryDevice.state`**：`available` / `allocated` / `unavailable` / `maintenance`。

### 2.3 Non-Frozen Capabilities

- 平台跨租户 GPU 排行/告警联动 — 部分 **TODO-YAML**
- BOSS 专属 RBAC scope — 待 Core 确认

### 2.4 Core 同步边界（2026-06-17）

- `gpu-inventory*` paths **路径已声明**
- 平台跨租户排行/告警联动 — 部分 **TODO-YAML**

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。

## 3. Page Scope

### 3.1 Page Responsibilities

- 平台侧页面

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

见主维护源 `docs/boss-modules/health/gpu-monitoring.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台级 BOSS API — **TODO-YAML**（交付页：**N/A**）
- Phase：**P1**

## 7. 主维护源与流程

- `docs/boss-modules/health/gpu-monitoring.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
