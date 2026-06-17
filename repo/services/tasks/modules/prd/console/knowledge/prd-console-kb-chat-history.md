# PRD: Console kb-chat-history

详文：`docs/console-modules/knowledge/kb-chat-history.md` · Phase 2 YAML 已声明项文档收口（2026-06-17）。

## 目标

- 为 **Services / KnowledgeBases / 对话历史** 提供可维护的产品边界说明
- 明确 sessions 列表已冻结、DELETE/PATCH 待补

## 用户故事（规划）

- 作为知识库用户，我希望查看历史会话列表，以便继续之前的问答
- 作为开发者，我希望从文档得知 session 创建在 kb-qa-flow、列表页只读

## 范围

- 会话侧栏/列表、跳转 kb-qa-flow 续聊 — 见主维护源
- 不修改 ANI-main OpenAPI（由 Services 团队按 TASK-SVC-014 推进）

## 非目标

- 不自造 session DELETE/PATCH 为冻结能力
- 跨知识库会话聚合
