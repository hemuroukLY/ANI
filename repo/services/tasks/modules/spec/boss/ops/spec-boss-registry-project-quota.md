# SPEC: BOSS registry-project-quota

> Technical specification derived from: `tasks/modules/prd/boss/ops/prd-boss-registry-project-quota.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **资源池与基础设施 / 镜像仓库 / 项目配额** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/ops/prd-boss-registry-project-quota.md`
- 主维护源: `docs/boss-modules/ops/registry-project-quota.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文）

| 方法 | 路径 | operationId | RBAC | BOSS 用法 |
|---|---|---|---|---|
| GET | `/api/v1/registry/projects` | `listRegistryProjects` | `scope:registry:read` | 租户 **只读参考** |
| POST | `/api/v1/registry/projects` | `createRegistryProject` | `scope:registry:write` | 租户创建 |

`RegistryProject` schema **不含 quota 字段**。平台 registry quota list/PATCH — **TODO-YAML P2**。

### 2.3 Non-Frozen Capabilities

- 平台 registry quota list/PATCH — **TODO-YAML** P2
- BOSS 跨租户 projects list — **TODO-YAML** P2

### 2.4 Core 同步边界（2026-06-17）

- Harbor 存储 quota **≠** tenant GPU quota（`tenant-quota-policy`）
- 租户 `registry/projects` — **只读参考**；禁止 JWT 轮询
- 422 仅 YAML 已声明 operation 可标冻结

> 详文满配（Core）见主维护源；OpenAPI 已声明 ≠ handler 已实现。

## 3. Page Scope

### 3.1 Page Responsibilities

- 全平台 Harbor 项目配额查看与调整（写 API 待 YAML）

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

见主维护源 `docs/boss-modules/ops/registry-project-quota.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台 registry quota API — **TODO-YAML** P2
- Phase：**P2**

## 7. 主维护源与流程

- `docs/boss-modules/ops/registry-project-quota.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
