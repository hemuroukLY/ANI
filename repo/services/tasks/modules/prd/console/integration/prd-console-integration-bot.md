# PRD: Console integration-bot

详文：`docs/console-modules/integration/integration-bot.md` · Phase 2 YAML 已声明项文档收口（2026-06-17）。

## 目标

- 为 **Services / Integrations / 企微·钉钉 Bot** 提供可维护的产品边界说明
- 明确 `POST .../integrations/bots` 已冻结、DELETE/PATCH 待补

## 用户故事（规划）

- 作为租户管理员，我希望创建企微/钉钉 Bot 集成，以便接收平台事件
- 作为只读用户，我希望查看 Bot 列表但无法创建

## 范围

- Bot 创建表单、integrations 列表筛选 — 见主维护源
- 不修改 ANI-main OpenAPI（由 Services 团队按 TASK-SVC-016 推进）

## 非目标

- 不自造 Bot DELETE/PATCH 为冻结能力
- Bot 消息投递日志（待补）
