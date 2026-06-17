# PRD: Console agent-tool-permission

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/ai-native/agent-tool-permission.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-agent-tool-permission-draft.md`

## 目标

- 定义 Agent **工具/MCP 授权**与租户 **API scope RBAC** 的产品边界
- 输出工具目录 + 授权矩阵 PUT 的可评审 OpenAPI 草案
- 支撑 Agent 运行时 tool dispatch 鉴权（消费方：Session / Orchestration）

## 用户故事

- 作为租户 Agent 管理员，我希望配置哪些 MCP/内置工具可被 Agent 调用
- 作为安全负责人，我希望对 high-risk 工具默认 deny，并支持角色级覆盖
- 作为开发者，我希望 deny 优先于 allow 的规则在文档与 API 中一致

## 范围

- GET tools 目录、GET/PUT tool-permissions 快照
- Console 授权矩阵 UI、冲突提示
- RBAC：`agent-tools:read|write`

## 非目标

- 租户 Role 的 `scope:*` 编辑（见 `role-permission-edit.md`）
- MCP 市场上架写操作（见 `mcp-tool-market.md`）
- 本阶段合入 ANI-main

## 成功标准

- [ ] 草案通过 Services 评审
- [ ] 与 Session 详化草案无 path 冲突
- [ ] 合入 YAML 后详文切换冻结口径
