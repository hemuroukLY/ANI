# PRD: Console 容器实例可观测性

> Revised: 2026-06-17  
> 详文：`docs/console-modules/compute/container-observability.md`

## 1. Overview

实例 logs/events/metrics/exec 子路径；TASK-CORE-005。

## 2. Goals

- Tab 化观测
- exec 仅 running + scope:instances:exec
- exec 必填 idempotency_key

## 3. User Stories

US-001 日志分页；US-002 指标；US-003 exec 终端。

## 4. References

- 详文、SPEC、CORE-HANDLER-IMPLEMENTATION-GUIDE §TASK-CORE-005
