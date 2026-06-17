# PRD: Console 批任务实例

> Revised: 2026-06-17  
> 详文：`docs/console-modules/compute/batch-job-instances.md`

## 1. Overview

统一实例架构 `kind=batch_job` 管理页；复用 `/api/v1/instances*`。

## 2. Goals

- 明确 list kind 枚举缺口（TODO-YAML）
- lifecycle 与 VM 一致

## 3. User Stories

US-001 列表（待 kind 扩展）；US-002 创建 batch_job；US-003 生命周期。

## 4. Functional Requirements

- FR-1: 不用 deprecated Services
- FR-2: create 422 见 module-delivery-workflow §2.10

## 5. References

- 详文、SPEC、vm-management.md
