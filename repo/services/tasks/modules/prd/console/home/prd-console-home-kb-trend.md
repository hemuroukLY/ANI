# PRD: Console home-kb-trend

详文：`docs/console-modules/home/home-kb-trend.md` · Phase 2 文档收口（2026-06-17）。

## 目标

- 为 **首页 / 知识库调用趋势** 区块提供可维护的产品边界说明
- 明确无独立趋势 API、部分指标待产品/Metering 口径确认

## 用户故事（规划）

- 作为租户成员，我希望在首页看到知识库调用趋势摘要，以便了解使用情况
- 作为产品经理，我希望待确认指标在 UI 显式标注，避免伪造曲线

## 范围

- 聚合 `knowledge-bases` 列表与 `metering/usage`（filter 待确认）— 见主维护源
- 不修改 ANI-main OpenAPI（专用趋势 API 未声明）

## 非目标

- 不自造 `GET /knowledge-trend` 或成功率专用 API
- 跨租户平台运营视图、正式报表导出（见 usage-report）
