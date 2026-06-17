# PRD: Console kb-doc-intelligence

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/knowledge/kb-doc-intelligence.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-kb-doc-intelligence-draft.md`

## 目标

- 定义 **文档智能分析**（摘要/实体/分块）与 upload/index 流水线的产品边界
- 输出可评审的 Services OpenAPI 草案（analyze + intelligence 只读）
- 为 TASK-SVC-018 知识库智能子项提供合入 YAML 的前置材料

## 用户故事

- 作为知识库编辑者，我希望对已索引文档触发增量分析并查看摘要与实体
- 作为运维，我希望分析异步返回 task_id 并在任务中心跟踪进度
- 作为开发者，我希望文档未 index 完成时 analyze 返回 422，而非误报 202

## 范围

- POST analyze、GET intelligence（3a）；history list（3b）
- Console 分析触发与结果展示 UI
- RBAC：knowledge-bases read/write

## 非目标

- 替换 upload 侧基础解析（`POST .../documents` 202）
- 合入 ANI-main（本阶段仅 service-ani 文档）

## 成功标准

- [ ] 草案通过 Services 架构评审
- [ ] AsyncTask 与 async-task-center 一致
- [ ] 合入 YAML 后详文切换为冻结口径
