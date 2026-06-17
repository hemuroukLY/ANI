# PRD: Console integration-webhook-overview

详文：`docs/console-modules/integration/integration-webhook-overview.md` · Phase 2 文档收口（2026-06-17）。

## 目标

- 为 **开放与集成 / Webhook 说明** 导航页提供可维护的产品边界说明
- 澄清接入域 Webhook 与租户 Webhook 分工，避免重复契约

## 用户故事（规划）

- 作为租户管理员，我希望理解 Webhook 与 Bot 的能力边界，以便选择正确配置入口
- 作为开发者，我希望文档不重复 tenant webhooks 契约、不自造 integrations/webhooks 路径

## 范围

- 导航说明、能力对比表、跳转子模块 — 见主维护源
- 不定义新 OpenAPI operation

## 非目标

- 不自造 `IntegrationWebhook` schema 或 `/integrations/webhooks` 路径
- 接入域统一 Webhook 注册、签名轮换 UI
