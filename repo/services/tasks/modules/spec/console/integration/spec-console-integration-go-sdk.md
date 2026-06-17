# SPEC: Console integration-go-sdk

> Technical specification derived from: `tasks/modules/prd/console/integration/prd-console-integration-go-sdk.md`  
> Revised: 2026-06-17

## 1. Summary

### 1.1 What This SPEC Covers

收口 Console **开放与集成 / Go SDK** 的正式边界、字段映射口径与待补依赖。只对齐当前 OpenAPI 已冻结能力，不重新设计 API。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/integration/prd-console-integration-go-sdk.md`
- 主维护源: `docs/console-modules/integration/integration-go-sdk.md`

## 2. Frozen Facts

### 2.1 Authority Source

- 仓内 SDK 源码（只读引用）（文档交付页无 YAML 权威源）

### 2.2 Verified Paths

**N/A** — 非 Console API 交付页（TODO-YAML: N/A）。

### 2.3 仓内源码

- `ANI-main/repo/sdks/core/go/anisdk/client.go`
- `ANI-main/repo/sdks/core/go/examples/basic/main.go`
- `ANI-main/repo/sdks/services/go/anisdk/client.go`
- `ANI-main/repo/sdks/services/go/examples/basic/main.go`

## 3. Page Scope

### 3.1 Page Responsibilities

- **非 Console API 交付页**：指向仓内 Go SDK 源码与示例，帮助租户自助接入 Core/Services HTTP API。…

### 3.2 Non-Goals

- 不改变未冻结的规划路径
- 不把 handler stub 标注为已实现
- 非 Console API 页不自造 REST 契约

## 4. 创建前置条件与操作矩阵

见主维护源 `docs/console-modules/integration/integration-go-sdk.md`（§创建前置条件、§操作可用性矩阵）。

## 5. 待补边界

- SDK 版本与 OpenAPI 生成流水线 — 见 repo Makefile
- <!-- TODO-YAML: N/A -->

## 6. 主维护源

- `docs/console-modules/integration/integration-go-sdk.md`
- 流程：`docs/console-modules/governance/module-delivery-workflow.md`

OpenAPI 已声明 ≠ handler 已实现。
