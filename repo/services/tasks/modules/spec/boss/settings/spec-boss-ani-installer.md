# SPEC: BOSS ani-installer

> Technical specification derived from: `tasks/modules/prd/boss/settings/prd-boss-ani-installer.md`  
> Revised: 2026-06-17 · 流程：ANI-14 Phase 1

## 1. Summary

### 1.1 What This SPEC Covers

收口 BOSS **交付与安装 / ani-installer** 的权威源事实、平台级技术边界、Console 对照关系与 TODO-YAML 依赖。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/boss/settings/prd-boss-ani-installer.md`
- 主维护源: `docs/boss-modules/settings/ani-installer.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`
- `ANI-main/ANI-07-部署工程设计.md`（交付类权威源）

### 2.2 Verified Paths（BOSS 上下文）

**当前无本模块可直接引用的 BOSS 平台级冻结路径。**

<!-- TODO-YAML: 平台级 BOSS API -->

### 2.3 Non-Frozen Capabilities

- 平台级 `ani-installer` 跨租户 list/aggregate API — **TODO-YAML**
- BOSS 专属 RBAC scope — 待 Core/Services 确认
- 见主维护源

### 2.4 Core 同步边界（2026-06-17）

- 交付安装工作流 **不在 OpenAPI** — **TODO-YAML: N/A**
- 实际操作面：`ani-installer` CLI/TUI
- Console **无** 对等页

> 详文满配（Core）见主维护源；交付类 REST **TODO-YAML: N/A** 除非另行标注；OpenAPI 已声明 ≠ handler 已实现。

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

见主维护源 `docs/boss-modules/settings/ani-installer.md`（§创建前置条件、§操作可用性矩阵）。

## 6. 待补边界

- 平台级 BOSS API — **TODO-YAML**（交付页：**N/A**）
- Phase：**P1**

## 7. 主维护源与流程

- `docs/boss-modules/settings/ani-installer.md`
- `docs/boss-modules/governance/module-delivery-workflow.md`
- `ANI-main/ANI-14-API对齐与开发工作流.md`

OpenAPI 已声明 ≠ handler 已实现；BOSS 平台聚合 ≠ 租户级路径直接复用。
