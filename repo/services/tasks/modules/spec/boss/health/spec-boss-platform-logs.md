# SPEC: BOSS platform-logs

> Technical specification derived from: `tasks/modules/prd/boss/health/prd-boss-platform-logs.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **运维与可观测 / 日志** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/health/prd-boss-platform-logs.md`
- 主维护源: `docs/boss-modules/health/platform-logs.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文）

**当前无日志 list/query 冻结路径。**

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/observability/query` | **PromQL 指标代理**，不得当作日志 API 引用 |

<!-- TODO-YAML: GET /api/v1/observability/logs 或等价 -->

### 2.3 Non-Frozen Capabilities

- 平台跨租户日志 list — **TODO-YAML** P1
- BOSS 专属 RBAC scope — 待 Core 确认

### 2.4 Core 同步边界（2026-06-17）

- 平台跨租户日志 list — **TODO-YAML P1**

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。

## 3. Page Scope

### 3.1 Page Responsibilities

- 平台服务日志检索页，支持按服务、租户、request_id 等维度排查。

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

见主维护源 `docs/boss-modules/health/platform-logs.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台级 BOSS API — **TODO-YAML**（交付页：**N/A**）
- Phase：**P1**

## 7. 主维护源与流程

- `docs/boss-modules/health/platform-logs.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
