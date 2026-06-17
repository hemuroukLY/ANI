# PRD: Console model-recommend-config

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/inference/model-recommend-config.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-recommend-config-draft.md`

## 目标

- 定义 **只读部署推荐**（GPU/副本/量化）与 inference deploy 的边界
- 输出 GET recommendations 草案
- 为 TASK-SVC-018 模型增强子项提供合入 YAML 的前置材料

## 用户故事

- 作为推理管理员，我希望按 chat/batch/embedding profile 查看多档推荐
- 作为用户，我希望一键预填部署表单但仍需确认提交 deploy
- 作为开发者，我希望无 POST apply API，推荐字段可映射 InferenceService create

## 范围

- GET `/models/{id}/recommendations?profile=...`
- Console explain、quota_sufficient、预填跳转

## 非目标

- 推荐算法实现细节
- 合入 ANI-main

## 成功标准

- [ ] 只读 GET；字段映射 inference-service 已文档化
- [ ] 合入 YAML 后详文切换为冻结口径
