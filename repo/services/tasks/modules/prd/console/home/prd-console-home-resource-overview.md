# PRD: Console 首页资源使用概览

> Revised: 2026-06-17  
> 详文：`docs/console-modules/home/home-resource-overview.md`

## 1. Overview

平台概览第一区块明细页；Core 多接口页面层聚合，无独立 overview API。

## 2. Goals

- 实例/GPU/存储/网络摘要
- 与 platform-overview 口径一致
- 局部降级

## 3. User Stories

### US-001: 查看摘要

验收：并发 Core list；标注刷新时间。

### US-002: 跳转明细

验收：VM/GPU/存储/网络 deep link。

## 4. Functional Requirements

- FR-1: 不自造 ResourceOverviewResponse
- FR-2: GPU 用 gpu-inventory 非 deprecated Services

## 5. Non-Goals

- 独立 overview API（TODO-YAML）

## 6. References

- 详文、SPEC
