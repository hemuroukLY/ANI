# PRD: Console tenant-webhook-ops

详文：`docs/console-modules/tenant/tenant-webhook-ops.md` · Phase 2 YAML 已声明项文档收口（2026-06-17）。

## 目标

- 为 **Services / Tenant / Webhook 投递运维** 提供可维护的产品边界说明
- 明确 `GET .../deliveries` 已冻结、手动 retry 未声明

## 用户故事（规划）

- 作为租户管理员，我希望查看 Webhook 投递记录与失败原因，以便排查集成问题
- 作为开发者，我希望文档不把 retry 写成冻结操作

## 范围

- 投递记录 Tab、分页与状态筛选 — 见主维护源
- Webhook CRUD 仍见 tenant-management.md
- 不修改 ANI-main OpenAPI（由 Services 团队按 TASK-SVC-015 推进）

## 非目标

- 不自造 `POST .../deliveries/{id}/retry`
- 投递 payload 明文重看、批量导出（待产品/安全定稿）
