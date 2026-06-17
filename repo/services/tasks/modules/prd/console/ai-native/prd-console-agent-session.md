# PRD: Console agent-session

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/ai-native/agent-session.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-session-draft.md`

## 目标

- 定义租户侧 **Agent 会话** 与 KB 问答会话的产品边界
- 输出可评审的 Services OpenAPI 草案（CRUD + messages 只读）
- 为 TASK-SVC-018 首子项提供合入 YAML 的前置材料

## 用户故事

- 作为 Agent 使用者，我希望创建并归档多轮 Agent 会话，便于续聊与审计
- 作为租户管理员，我希望会话与 KB 问答分离，避免权限与数据模型混淆
- 作为开发者，我希望从 SPEC/草案得知 planned paths，合入前不误认为已实现

## 范围

- 会话 list/create/get/patch；messages list（3b）
- Console 列表/详情/归档 UI 边界
- RBAC：`agent-sessions:read|write`

## 非目标

- 流式推理 Gateway（归属 orchestration / Gateway 扩展）
- Agent Profile / MCP 市场（后续子模块）
- 合入 ANI-main（本阶段仅 service-ani 文档）

## 成功标准

- [ ] 草案通过 Services 架构评审
- [ ] 合入 YAML 后 agent-session 详文切换为冻结口径
- [ ] 与 kb-chat-history 无 schema 冲突
