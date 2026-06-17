# SPEC: BOSS gpu-pool-trend

> Technical specification derived from: `tasks/modules/prd/boss/overview/prd-boss-gpu-pool-trend.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **平台运营总览 / GPU 资源池态势** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/overview/prd-boss-gpu-pool-trend.md`
- 主维护源: `docs/boss-modules/overview/gpu-pool-trend.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文）

| 方法 | 路径 | operationId | 成功响应 | RBAC | BOSS 用法 |
|---|---|---|---|---|---|
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | `200` | `scope:gpu-inventory:read` | 只读参考；平台跨租户聚合 **待补** |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | `200` | `scope:gpu-inventory:read` | 只读参考；平台跨租户聚合 **待补** |

<!-- 上表为 Core YAML 已声明路径；BOSS 平台级 list/aggregate 仍待 TODO-YAML -->

### 2.3 Non-Frozen Capabilities

- 平台级 `GPU 资源池态势` 跨租户 list/aggregate API — **TODO-YAML**
- BOSS 专属 RBAC scope — 待 Core/Services 确认
- 见主维护源

### 2.4 Core 同步边界（2026-06-17）

- `gpu-inventory` YAML **已声明**；BOSS 正式契约须 platform path
- 租户路径 **只读参考**；PromQL **≠** inventory list

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

见主维护源 `docs/boss-modules/overview/gpu-pool-trend.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台级 BOSS API — **TODO-YAML**（交付页：**N/A**）
- Phase：**P1**

## 7. 主维护源与流程

- `docs/boss-modules/overview/gpu-pool-trend.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
