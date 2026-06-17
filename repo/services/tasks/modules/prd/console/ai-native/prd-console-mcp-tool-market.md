# PRD: Console mcp-tool-market

> Phase 3 详化 · 2026-06-17  
> 详文：`docs/console-modules/ai-native/mcp-tool-market.md`  
> OpenAPI 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-mcp-tool-market-draft.md`

## 目标

- 定义租户侧 **MCP 工具市场** 与工具权限、第三方集成的边界
- 输出可评审的 Services OpenAPI 草案（catalog + install/uninstall）
- 为 TASK-SVC-018 MCP 子项提供合入 YAML 的前置材料

## 用户故事

- 作为 Agent 管理员，我希望浏览并安装 MCP 工具包，安装后可在权限页配置 allow/deny
- 作为安全管理员，我希望看到包的 risk_level 与 publisher，便于审批
- 作为开发者，我希望 install 幂等与 tool id 命名（`mcp:{package_id}/{tool_name}`）在 SPEC 中明确

## 范围

- 市场 list/detail、install（3a）、uninstall（3b）
- Console catalog / 详情 / 安装 UI 边界
- RBAC：`agent-mcp:read|write`

## 非目标

- BOSS 社区包审核与计费（Console 只读 catalog）
- 合入 ANI-main（本阶段仅 service-ani 文档）
- 替代 `/integrations` 通用集成

## 成功标准

- [ ] 草案通过 Services 架构评审
- [ ] 安装后 tool_ids 与 `listAgentTools` 一致
- [ ] 合入 YAML 后详文切换为冻结口径
