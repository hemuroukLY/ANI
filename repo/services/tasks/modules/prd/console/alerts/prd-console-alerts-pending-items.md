# PRD: Console 告警与待处理事项

> Revised: 2026-06-17  
> 详文主维护源：`docs/console-modules/alerts/alerts-pending-items.md`

## 1. Overview

`Console / 首页与总览 / 告警与待处理事项` 汇总租户需优先处理的风险与待办摘要，并提供跳转至来源模块的入口。本页为 **Console 页面层聚合**，当前无独立告警 list API。

## 2. 权威源结论

- Core：`GET /api/v1/tasks/{task_id}`、`GET /api/v1/observability/alert-rules*`、`GET /api/v1/observability/query`
- Services：`GET /api/v1/svc/inference-services` 等
- **待补**：`listAlerts` / `listPendingItems`（TODO-YAML）

## 3. Goals

- 在首屏可见区域展示高优先级告警与失败任务摘要
- 支持 Core + Services 多来源聚合（页面层）
- 明确标注当前无独立告警 list API 的边界
- 单来源失败时局部降级

## 4. User Stories

### US-001: 识别待处理风险

作为租户用户，我希望一眼看到需优先处理的事项数量。

验收标准：

- 展示四类摘要计数（告警 / 处理中 / 失败 / 待确认）
- 无数据时显示空态或「待接入」，不伪造 0

### US-002: 跳转处置

作为租户用户，我希望从待办项直接进入来源模块处置。

验收标准：

- 每条待办可跳转到对应模块（任务中心、推理服务、知识库等）
- 跳转目标为当前可用模块

### US-003: 局部失败降级

作为租户用户，当某一来源接口失败时，其他区块仍可用。

验收标准：

- 单来源错误不导致整页白屏
- 展示失败来源与 `request_id`（若有）

## 5. Functional Requirements

- FR-1: 支持摘要计数展示
- FR-2: 支持分类列表与局部失败态
- FR-3: 不得自造 `Alert` / `PendingItem` schema
- FR-4: 标注 `TODO-YAML` 待冻结聚合 API
- FR-5: 不得把 `alert-rules` 列表长度当作 firing 告警数

## 6. Non-Goals

- 告警规则 CRUD（归属 observability）
- 批量关闭告警（待专用 API）
- 平台全量告警运营（BOSS 范围）

## 7. References

- 详文：`docs/console-modules/alerts/alerts-pending-items.md`
- SPEC：`tasks/modules/spec/console/alerts/spec-console-alerts-pending-items.md`
- 工作流：`docs/console-modules/governance/module-delivery-workflow.md` §2.10
