# PRD: Console agent-orchestration

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/ai-native/agent-orchestration.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-orchestration-draft.md`

## 目标

- 定义多 Agent **工作流编排**（DAG、触发器、运行监控）的产品边界
- 输出可评审的 Services OpenAPI 草案（workflow CRUD + run 启动/查询）
- 为 TASK-SVC-018 Orchestration 子项提供合入 YAML 的前置材料

## 用户故事

- 作为编排管理员，我希望创建 workflow 并手动/定时触发运行
- 作为运维，我希望 run 返回 `task_id` 并跳转异步任务中心
- 作为开发者，我希望 workflow `status=active` 才能启动 run，非 active 返回 422

## 范围

- workflow list/create/get/patch；run start/list/get（3b 单 run 详情）
- Console 列表/编辑器/运行监控 UI 边界
- RBAC：`agent-workflows:read|write|run`

## 非目标

- `definition` DAG schema 细节（Phase 3b ADR）
- Gateway 流式节点类型
- 合入 ANI-main

## 成功标准

- [ ] 草案通过评审；run 202 与 Core AsyncTask 对齐
- [ ] 与 agent-session、agent-audit 跳转关系清晰
- [ ] 合入 YAML 后详文切换为冻结口径
