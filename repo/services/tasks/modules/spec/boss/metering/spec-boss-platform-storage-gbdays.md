# SPEC: BOSS platform-storage-gbdays

> Technical specification derived from: `tasks/modules/prd/boss/metering/prd-boss-platform-storage-gbdays.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **平台计量与结算 / 平台 Storage-GBDays** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/metering/prd-boss-platform-storage-gbdays.md`
- 主维护源: `docs/boss-modules/metering/platform-storage-gbdays.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文）

| 方法 | 路径 | operationId | 成功响应 | RBAC | BOSS 用法 |
|---|---|---|---|---|---|
| GET | `/api/v1/metering/usage` | `—` | `200` | 租户 JWT 上下文 | 路径已声明，但 **无** storage 类 `resource_type` |

**`MeteringUsageRecord.resource_type` enum 当前不含 storage** — 本页 resource_type 与平台 aggregate 均 **TODO-YAML**。

### 2.3 Non-Frozen Capabilities

- Storage `resource_type` enum — **待 YAML**
- 平台级跨租户 Storage 用量 list/aggregate — **TODO-YAML** P1

### 2.4 Core 同步边界（2026-06-17）

- Storage `resource_type` enum — **待 YAML**
- 平台跨租户 aggregate — **TODO-YAML P1**

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。

## 3. Page Scope

### 3.1 Page Responsibilities

- 跨租户存储用量（resource_type **待 YAML**）排行与趋势。

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

见主维护源 `docs/boss-modules/metering/platform-storage-gbdays.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台级 BOSS API — **TODO-YAML**（交付页：**N/A**）
- Phase：**P1**

## 7. 主维护源与流程

- `docs/boss-modules/metering/platform-storage-gbdays.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
