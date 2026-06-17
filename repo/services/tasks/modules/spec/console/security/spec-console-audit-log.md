# SPEC: Console audit-log

> Technical specification derived from: `tasks/modules/prd/console/security/prd-console-audit-log.md`  
> Revised: 2026-06-17

## 1. Summary

### 1.1 What This SPEC Covers

收口 Console **安全与身份 / 审计日志** 的正式边界、字段映射口径与待补依赖。只对齐当前 OpenAPI 已冻结能力，不重新设计 API。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/security/prd-console-audit-log.md`
- 主维护源: `docs/console-modules/security/audit-log.md`

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

**当前无本模块冻结路径。**

<!-- TODO-YAML -->

## 3. Page Scope

### 3.1 Page Responsibilities

- 租户侧 **平台操作审计** 检索页（谁、何时、对何资源、何种动作）；Console 只读，合规导出见 BOSS/`compliance.md`。…

### 3.2 Non-Goals

- 不改变未冻结的规划路径
- 不把 handler stub 标注为已实现
- 非 Console API 页不自造 REST 契约

## 4. 创建前置条件与操作矩阵

见主维护源 `docs/console-modules/security/audit-log.md`（§创建前置条件、§操作可用性矩阵）。

## 5. 待补边界

- Audit list API — **TODO-YAML**（Tag 已存在）
- 与 Agent 审计、API Key 审计分工

## 6. 主维护源

- `docs/console-modules/security/audit-log.md`
- 流程：`docs/console-modules/governance/module-delivery-workflow.md`

OpenAPI 已声明 ≠ handler 已实现。
