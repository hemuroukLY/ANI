# PRD: Console 异步任务中心

> Revised: 2026-06-17  
> 详文：`docs/console-modules/alerts/async-task-center.md`

## 1. Overview

`Console / 全局 / 异步任务中心` 统一展示 `202 + AsyncTask` 触发的异步操作进度，查询走 Core `GET /api/v1/tasks/{task_id}`。

## 2. 权威源结论

- Core：`AsyncTask` schema；`GET /api/v1/tasks/{task_id}`
- **待补**：`GET /api/v1/tasks` list（TODO-YAML）

## 3. Goals

- 展示已知 task_id 的任务状态与进度
- 支持按 `task_type`、`status` UI 筛选
- 关联资源跳转

## 4. User Stories

### US-001: 查看任务进度

验收标准：展示 `status`、`progress_pct`、`error_message`；与 YAML enum 一致。

### US-002: 轮询刷新

验收标准：running/pending 任务可轮询；间隔 ≥ 2s。

### US-003: 跳转关联资源

验收标准：`resource_type` + `resource_id` 可 deep link（若有）。

## 5. Functional Requirements

- FR-1: 单任务查询对齐 Core YAML
- FR-2: list 未冻结前仅展示会话内 task_id 集合
- FR-3: 不自造 `TaskListResponse`
- FR-4: 取消任务 UI 待 API 冻结前隐藏

## 6. Non-Goals

- 全量租户任务审计
- 任务取消（待 YAML）

## 7. References

- 详文、SPEC、module-delivery-workflow §2.10
