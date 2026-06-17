# PRD: Console 推理服务可观测性

> Revised: 2026-06-17  
> 详文：`docs/console-modules/inference/inference-observability.md`

## 1. Overview

推理服务详情下的日志、指标 Tab；日志走 Services，指标走 Core PromQL 代理。

## 2. Goals

- 日志 cursor 分页
- 指标经 `GET /api/v1/observability/query`
- 事件 Tab 待 YAML 前隐藏或空态

## 3. User Stories

### US-001: 查看日志

验收：`getInferenceServiceLogs` 对齐 YAML。

### US-002: 查看指标

验收：PromQL 查询；403 无 observability 权限时提示。

## 4. Functional Requirements

- FR-1: 不自造 `.../metrics` Services 路径
- FR-2: 事件标注 TODO-YAML

## 5. Non-Goals

- 集群级 observability 控制台

## 6. References

- 详文、SPEC
