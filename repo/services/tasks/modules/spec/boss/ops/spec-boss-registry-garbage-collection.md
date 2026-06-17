# SPEC: BOSS registry-garbage-collection

> Technical specification derived from: `tasks/modules/prd/boss/ops/prd-boss-registry-garbage-collection.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **资源池与基础设施 / 镜像仓库 / 垃圾回收** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/ops/prd-boss-registry-garbage-collection.md`
- 主维护源: `docs/boss-modules/ops/registry-garbage-collection.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（BOSS 上下文 · 体量只读参考）

**当前无 registry GC 冻结 path。**

| 方法 | 路径 | operationId | RBAC | BOSS 用法 |
|---|---|---|---|---|
| GET | `/api/v1/registry/projects` | `listRegistryProjects` | `scope:registry:read` | 项目 list 体量参考 |
| GET | `/api/v1/registry/projects/{project}/repositories` | `listRegistryRepositories` | `scope:registry:read` | `artifact_count` 参考 |
| GET | `/api/v1/registry/projects/{project}/repositories/{repository}/artifacts` | `listRegistryArtifacts` | `scope:registry:read` | `size_bytes` 汇总参考 |
| GET | `/api/v1/tasks/{task_id}` | — | — | 异步 GC **路径已声明**；无 list |

GC jobs list / run / schedule — **TODO-YAML P2**。

### 2.3 Non-Frozen Capabilities

- 全部 GC REST — **TODO-YAML** P2
- `GET /api/v1/tasks` list — **TODO-YAML**

### 2.4 Core 同步边界（2026-06-17）

- artifact `size_bytes` 汇总 **≠** GC 可回收空间
- `AsyncTask.status` 用 `completed`，禁 `succeeded`
- 409/422 冲突示例 — 待 YAML 声明后可标冻结

> 详文满配（Core）见主维护源。

## 3. Page Scope

### 3.1 Page Responsibilities

- 全平台 Harbor GC 历史、调度与手动触发（API 待 YAML）

### 3.2 Non-Goals

- 不自造 `RegistryGCJob` 为已冻结 schema
- 不把 artifact list 写成 GC 任务 list

## 4. Architecture Constraints（ANI-14）

| 约束 | 要求 |
|---|---|
| Core 前缀 | `/api/v1/*` |
| 写操作 | `idempotency_key` + 明确成功码 |
| 错误结构 | `code` / `message` / `request_id` |

## 5. 创建前置条件与操作矩阵

见主维护源 `docs/boss-modules/ops/registry-garbage-collection.md`。

## 6. 待补边界

- registry GC API — **TODO-YAML** P2
- Phase：**P2**

## 7. 主维护源与流程

- `docs/boss-modules/ops/registry-garbage-collection.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现。
