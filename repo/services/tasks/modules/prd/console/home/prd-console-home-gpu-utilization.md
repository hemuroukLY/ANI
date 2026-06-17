# PRD: Console 首页 GPU 利用率

> Revised: 2026-06-17  
> 详文：`docs/console-modules/home/home-gpu-utilization.md`

## 1. Overview

平台概览第二区块；Core gpu-inventory + occupancy + 可选 observability/query。

## 2. Goals

- GPU 总量/分配/利用率摘要
- 跳转 GPU 管理
- 非 deprecated Services GPU 路径

## 3. User Stories

### US-001: 利用率摘要

验收：occupancy + inventory 口径与 platform-overview 一致。

### US-002: 趋势图

验收：有 observability 权限时可查 PromQL（模板待冻结）。

## 4. Functional Requirements

- FR-1: scope:gpu-inventory:read
- FR-2: handler stub ≠ 契约未声明

## 5. References

- 详文、SPEC、gpu-inventory-ui.md
