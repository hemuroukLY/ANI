# SPEC: BOSS platform-resource-pool

> Technical specification derived from: `tasks/modules/prd/boss/ops/prd-boss-platform-resource-pool.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **资源池与基础设施 / 平台资源池总览** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/ops/prd-boss-platform-resource-pool.md`
- 主维护源: `docs/boss-modules/ops/platform-resource-pool.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文）

| 方法 | 路径 | operationId | BOSS 用法 |
|---|---|---|---|
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | 租户 **只读参考** |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | 租户 **只读参考** |
| GET | `/api/v1/instances` | `listInstances` | 租户 **只读参考** |
| GET | `/api/v1/volumes` | `listVolumes` | 租户 **只读参考** |
| GET | `/api/v1/networks/vpcs` | `listVPCs` | 租户 **只读参考** |

平台 resource-pool aggregate — **TODO-YAML P1**

### 2.3 Non-Frozen Capabilities

- 平台 pool aggregate、nodes list — **TODO-YAML** P1
- 平台 RBAC scope — **TODO-YAML**

### 2.4 Core 同步边界（2026-06-17）

- **无** 独立 `GET /api/v1/resource-pools/platform`
- 租户已声明路径仅 **只读参考**；禁止 JWT 轮询冒充平台 KPI

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。

## 3. Page Scope

### 3.1 Page Responsibilities

- 算力域 landing：GPU/节点/存储/网络四象限摘要 + 跳转 ops 子模块

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

见主维护源 `docs/boss-modules/ops/platform-resource-pool.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台级 BOSS API — **TODO-YAML**（交付页：**N/A**）
- Phase：**P1**

## 7. 主维护源与流程

- `docs/boss-modules/ops/platform-resource-pool.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
