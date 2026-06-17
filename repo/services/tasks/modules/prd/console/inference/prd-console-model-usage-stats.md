# PRD: Console model-usage-stats

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/inference/model-usage-stats.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md`

## 目标

- 定义 **模型维度用量统计** 与租户用量报表的分工
- 输出 Core metering 扩展 + 可选 Services summary 草案
- 为 TASK-SVC-018 模型增强子项提供合入 YAML 的前置材料

## 用户故事

- 作为模型管理员，我希望在模型中心按 model_id 查看 Token 与请求趋势
- 作为运维，我希望按推理服务拆分同一模型的用量
- 作为开发者，我希望优先扩展 `GET /metering/usage` filter，避免重复计量库

## 范围

- Core：`model_id`、`inference_service_id` query 扩展；group_by 扩展（3b）
- Services 可选：`GET /models/{id}/usage-stats`
- Console Dashboard、跳转 usage-report

## 非目标

- 账单导出（billing-export）
- Console 调用 token-usage 上报 API

## 成功标准

- [ ] 与 usage-report 口径联评通过
- [ ] 合入 YAML 后详文切换为冻结口径
