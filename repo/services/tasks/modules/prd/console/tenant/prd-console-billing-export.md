# PRD: Console billing-export

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/tenant/billing-export.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-billing-export-draft.md`

## 目标

- 定义 **用量报表异步导出** 与 usage-report、BOSS 结算的边界
- 输出 POST billing/exports 草案（202 AsyncTask）
- 为 TASK-SVC-018 / Phase 3 §4 提供合入 YAML 前置材料

## 用户故事

- 作为租户管理员，我希望导出 CSV/PDF 用量明细
- 作为用户，我希望大 export 异步并在任务中心跟踪
- 作为开发者，我希望 export 不含 BOSS 发票金额

## 范围

- POST create export、GET export job；list（3b）
- metering query 参数对齐

## 非目标

- 发票、对账、结算 API
- 合入 ANI-main

## 成功标准

- [ ] 202 + billing.export task_type
- [ ] 合入 YAML 后详文切换为冻结口径
