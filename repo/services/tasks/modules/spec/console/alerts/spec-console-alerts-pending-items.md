# SPEC: Console 告警与待处理事项

> Technical specification derived from: `tasks/modules/prd/console/alerts/prd-console-alerts-pending-items.md`  
> Revised: 2026-06-17

## 1. Summary

Console 页面层聚合高优先级风险与待办摘要；**无独立冻结 list API**。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths（引用，非本页专属）

| Method | Path | operationId | 成功响应 |
|---|---|---|---|
| GET | `/api/v1/tasks/{task_id}` | — | `200 + AsyncTask` |
| GET | `/api/v1/observability/alert-rules` | `listObservabilityAlertRules` | `200 + ObservabilityAlertRuleListResponse` |
| GET | `/api/v1/observability/query` | `queryObservability` | `200 + ObservabilityQueryResponse` |
| GET | `/api/v1/svc/inference-services` | `listInferenceServices` | `200` + items |

### 2.3 TODO-YAML

- `listAlerts` / `listPendingItems` 或等价只读聚合 API

## 3. Page Scope

### 3.1 Page Responsibilities

- 四类摘要计数 + 分类列表
- 跳转任务中心、推理、知识库、安全概览

### 3.2 Non-Goals

- 不自造 Alert schema
- 不承担告警规则编辑

## 4. Architecture

Console 并发拉取各来源 → 轻聚合 → 摘要 + 列表；冲突时以详文为准。

## 5. Error Rules

- 各引用 operation 错误码见详文 §接口冻结规则
- `422` 口径见 `../governance/module-delivery-workflow.md` §2.10

## 6. References

- 详文：`docs/console-modules/alerts/alerts-pending-items.md`
