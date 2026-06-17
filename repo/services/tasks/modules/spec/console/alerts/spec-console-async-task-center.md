# SPEC: Console 异步任务中心

> Source: `tasks/modules/prd/console/alerts/prd-console-async-task-center.md`  
> Revised: 2026-06-17

## 1. Summary

Core `AsyncTask` 统一任务模型；当前仅单任务 GET 冻结。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | 成功响应 | 错误 |
|---|---|---|---|
| GET | `/api/v1/tasks/{task_id}` | `200 + AsyncTask` | `404` |

### 2.3 Verified Schemas

- `AsyncTask`（status: pending/running/completed/failed/cancelled/dead_letter）

### 2.4 TODO-YAML

- `GET /api/v1/tasks` list + cursor
- 任务取消 API

## 3. Page Scope

- 任务列表（客户端聚合）、详情、轮询、资源跳转

## 4. Non-Goals

- 不自造 list schema
- 不宣称 handler 已实现

## 5. Data Mapping

| 页面字段 | Schema 字段 |
|---|---|
| 进度 | `progress_pct` |
| 类型 | `task_type` |
| 失败原因 | `error_message` |

## 6. References

- 详文：`docs/console-modules/alerts/async-task-center.md`
