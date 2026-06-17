# PRD: Console integration-py-sdk

详文：`docs/console-modules/integration/integration-py-sdk.md` · P2 文档收口（2026-06-17）。

## 目标

- 为导航 **开放与集成 / Python SDK** 提供可维护的产品边界说明
- 明确需先补 OpenAPI：**N/A**（Phase: 非 Console API）

## 用户故事（规划）

- 作为租户管理员，我希望在 Console 查看/配置与本模块相关的能力边界
- 作为开发者，我希望从文档得知哪些 API 已冻结、哪些仍为 TODO-YAML

## 范围

- Console 页面定位、矩阵、验收标准 — 见主维护源
- 不修改 ANI-main OpenAPI（由 Core/Services 团队按 TASK 推进）

## 非目标

- 不自造 schema 为冻结事实
- 不替代 BOSS 平台运营能力（若适用）
