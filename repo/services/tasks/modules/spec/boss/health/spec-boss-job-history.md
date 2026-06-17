# SPEC: BOSS job-history

> Technical specification derived from: `tasks/modules/prd/boss/health/prd-boss-job-history.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **运维与可观测 / 任务历史** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/health/prd-boss-job-history.md`
- 主维护源: `docs/boss-modules/health/job-history.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文）

| 方法 | 路径 | operationId | 成功响应 | RBAC | BOSS 用法 |
|---|---|---|---|---|---|
| GET | `/api/v1/tasks/{task_id}` | — | `200` | **路径已声明**；RBAC 待 YAML | 单任务详情；**无 list** |

**AsyncTask.status 冻结枚举**：`pending` / `running` / `completed` / `failed` / `cancelled` / `dead_letter`（禁止 UI 使用 `succeeded`）。

<!-- TODO-YAML: GET /api/v1/tasks list with cursor + tenant filter -->

### 2.3 Non-Frozen Capabilities

- 平台级 `GET /api/v1/tasks` list（cursor、tenant_id、status、task_type）— **TODO-YAML** P1
- 平台 task 重试/取消 — **TODO-YAML** Phase 2
- BOSS 专属 RBAC scope — 待 Core/Services 确认
- Console `async-task-center.md` 为租户视角；客户端聚合已知 task_id

### 2.4 Core 同步边界（2026-06-17）

- `GET /api/v1/tasks/{task_id}` — **路径已声明**（无 operationId / 无 x-ani-rbac-scope）
- 平台 `GET /api/v1/tasks` list — **TODO-YAML P1**
- 403 示例限定 **list TODO-YAML** 场景；单查 YAML 仅声明 404

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。

## 3. Page Scope

### 3.1 Page Responsibilities

- 平台级异步任务历史页，查看所有租户后台任务状态，支持重试与 dead letter 排查。

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

见主维护源 `docs/boss-modules/health/job-history.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台级 BOSS API — **TODO-YAML**（交付页：**N/A**）
- Phase：**P1**

## 7. 主维护源与流程

- `docs/boss-modules/health/job-history.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
