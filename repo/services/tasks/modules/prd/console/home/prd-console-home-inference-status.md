# PRD: Console 首页推理服务状态

> Revised: 2026-06-17  
> 详文：`docs/console-modules/home/home-inference-status.md`

## 1. Overview

平台概览第三区块；Services `listInferenceServices` 按 status 聚合计数。

## 2. Goals

- 状态枚举与 InferenceService.status 一致
- failed > 0 高亮
- 「最近活跃」待确认处标注

## 3. User Stories

### US-001: 状态分布

验收：running/deploying/failed/stopped 计数正确。

### US-002: 跳转列表

验收：跳转 inference-service 列表。

## 4. Functional Requirements

- FR-1: 不自造首页专属 schema

## 5. Non-Goals

- 本页 CRUD

## 6. References

- 详文、SPEC
