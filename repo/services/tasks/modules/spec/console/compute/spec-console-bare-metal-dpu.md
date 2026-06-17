# SPEC: Console bare-metal-dpu

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-bare-metal-dpu.md`  
> Revised: 2026-06-17

## 1. Summary

### 1.1 What This SPEC Covers

收口 Console **算力 / 实例 / 裸金属·DPU** 的正式边界、字段映射口径与待补依赖。只对齐当前 OpenAPI 已冻结能力，不重新设计 API。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/compute/prd-console-bare-metal-dpu.md`
- 主维护源: `docs/console-modules/compute/bare-metal-dpu.md`

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/instances?kind=bare_metal` | `listInstances` | `200 + InstanceListResponse` | `scope:instances:read` |
| GET | `/api/v1/instances?kind=dpu_node` | `listInstances` | `200 + InstanceListResponse` | `scope:instances:read` |
| POST | `/api/v1/instances` | `createInstance` | `201 + CreateInstanceResponse` | `scope:instances:create` |

<!-- 模块完整 API 仍待 TODO-YAML -->


## 3. Page Scope

### 3.1 Page Responsibilities

- 展示与管理 **`kind=bare_metal` / `kind=dpu_node`** 统一实例的 Console 页；复用 Core instances 契约，非独立裸金属 CRUD 域。…

### 3.2 Non-Goals

- 不改变未冻结的规划路径
- 不把 handler stub 标注为已实现
- 非 Console API 页不自造 REST 契约

## 4. 创建前置条件与操作矩阵

见主维护源 `docs/console-modules/compute/bare-metal-dpu.md`（§创建前置条件、§操作可用性矩阵）。

## 5. 待补边界

- DPU 设备清单 `dpu-inventory` — YAML 顶部规划，paths **未声明**
- 裸金属专有 provisioning 字段 — 待 Core 扩 schema

## 6. 主维护源

- `docs/console-modules/compute/bare-metal-dpu.md`
- 流程：`docs/console-modules/governance/module-delivery-workflow.md`

OpenAPI 已声明 ≠ handler 已实现。
