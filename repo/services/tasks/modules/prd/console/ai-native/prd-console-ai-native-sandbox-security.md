# PRD: Console ai-native-sandbox-security

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/ai-native/ai-native-sandbox-security.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-ai-native-sandbox-security-draft.md`

## 目标

- 定义 **Agent 运行时 Sandbox 安全配置** 与算力域 `kind=sandbox` 实例的边界
- 输出可评审的 Services OpenAPI 草案（profile CRUD + 安全摘要）
- 为 TASK-SVC-018 Sandbox Security 子项提供合入 YAML 的前置材料

## 用户故事

- 作为安全管理员，我希望配置 isolation_level 与 egress 白名单
- 作为 Agent 使用者，我希望 `allow_tool_ids` 与 tool-permission 交集生效有 UI 说明
- 作为开发者，我希望禁止 Services deprecated `/sandboxes*`，Core 实例只读引用

## 范围

- sandbox-profiles list/get/put；security summary（3b）
- Console profile 编辑与安全摘要 UI
- RBAC：`agent-sandbox:read|write`
- Core 只读：`instances?kind=sandbox`、`security-events`

## 非目标

- Core 实例 CRUD（归属 `sandbox-instance-management.md`）
- BOSS 全局 baseline 下发细节
- 合入 ANI-main

## 成功标准

- [ ] 草案通过评审；与算力 Sandbox 分域清晰
- [ ] egress 422 与 agent-audit 扩展事件口径待联评
- [ ] 合入 YAML 后详文切换为冻结口径
