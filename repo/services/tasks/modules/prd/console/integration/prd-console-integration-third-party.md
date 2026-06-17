# PRD: Console integration-third-party

详文：`docs/console-modules/integration/integration-third-party.md` · Phase 2 YAML 已声明项文档收口（2026-06-17）。

## 目标

- 为 **Services / Integrations / 第三方集成** 提供可维护的产品边界说明
- 明确 list/create 已冻结、按 id 的 GET/PATCH/DELETE 未声明

## 用户故事（规划）

- 作为租户管理员，我希望创建并查看第三方集成，以便对接业务系统
- 作为安全要求，我希望 secret 不在 UI 明文展示

## 范围

- 集成列表、创建、状态展示、详情抽屉 — 见主维护源
- 不修改 ANI-main OpenAPI（由 Services 团队按 TASK-SVC-016 推进）

## 非目标

- 不自造 `integration_id` 写路径
- 连通性测试 API、OAuth 回调（待 Gateway/Services 规划）
