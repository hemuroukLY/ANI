# SPEC: BOSS tenant-usage-billing

> Derived from: `tasks/modules/prd/boss/tenant/prd-boss-tenant-usage-billing.md`  
> Revised: 2026-06-17

## 1. Summary

跨租户计量与 Console usage-report 口径对齐。

主维护源：`docs/boss-modules/tenant/tenant-usage-billing.md`  
Console 对照：`docs/console-modules/tenant/usage-report.md`

## 2. Frozen Facts

### 2.1 Core — Tenant Context（只读参考）

| 方法 | 路径 | operationId | 成功 | 说明 |
|---|---|---|---|---|
| GET | `/api/v1/metering/usage` | —（未声明） | 200 | 单租户 JWT |
| POST | `/api/v1/metering/token-usage` | `reportTokenUsage` | 202 | 上报 |

Query（usage）：`start_time`、`end_time` required；`resource_type`、`group_by` optional。  
`group_by` enum：`resource_type`、`az`、`day`、`hour`。

Response inline：`items[]`（resource_type, total_quantity, unit, period）、`total`、`dev_profile`。

Errors：`400`、`401`、`403`。

### 2.2 BOSS Platform API

**当前无。** 跨租户 list/aggregate — **TODO-YAML**。

建议：`GET /api/v1/metering/usage/platform` 或扩展 query `group_by=tenant_id`（须平台 RBAC）。

### 2.3 UI Perspectives（非独立 API）

GPU-Hours、CPU-Hours、Memory-GBHours、Storage-GBDays、Input Tokens、Output Tokens、KB Queries — 同 Console usage-report §视角说明。

### 2.4 Core 同步边界（2026-06-17）

- 平台跨租户 aggregate — **TODO-YAML P1**
- 租户 `GET /api/v1/metering/usage` 为 JWT 上下文只读参考，**非** BOSS 正式平台契约
- 403 平台 list 场景示例为 `permission denied` 参考

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。

## 3. Boundaries

- 不得用逐租户 JWT 轮询作为 BOSS 正式架构
- `/metering/token-usage` 不是 BOSS 查询接口
- billing-export / 发票 / 对账 — 非本页

## 4. 前置条件与矩阵

见主维护源。

## 5. 待补

- 平台 metering aggregate YAML
- 与 BOSS metering 域口径统一
