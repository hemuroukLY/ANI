# 租户与客户管理（BOSS）

> **PRD/SPEC/HTML 同步**：tenant / health / metering 三域 21 模块已于 2026-06-17 与满配详文对齐。

平台侧租户治理域，与 Console `tenant/` 域分工：**BOSS 管开通/停用/配额/平台用量；Console 管租户内成员/SSO/Webhook**。

| 页面 | 详文 | PRD | SPEC |
|---|---|---|---|
| 租户列表 | [tenant-list.md](tenant-list.md) | [prd-boss-tenant-list.md](../../tasks/modules/prd/boss/tenant/prd-boss-tenant-list.md) | [spec-boss-tenant-list.md](../../tasks/modules/spec/boss/tenant/spec-boss-tenant-list.md) |
| 配额策略 | [tenant-quota-policy.md](tenant-quota-policy.md) | [prd-boss-tenant-quota-policy.md](../../tasks/modules/prd/boss/tenant/prd-boss-tenant-quota-policy.md) | [spec-boss-tenant-quota-policy.md](../../tasks/modules/spec/boss/tenant/spec-boss-tenant-quota-policy.md) |
| 租户管理员 | [tenant-admin.md](tenant-admin.md) | [prd-boss-tenant-admin.md](../../tasks/modules/prd/boss/tenant/prd-boss-tenant-admin.md) | [spec-boss-tenant-admin.md](../../tasks/modules/spec/boss/tenant/spec-boss-tenant-admin.md) |
| 租户计费与用量 | [tenant-usage-billing.md](tenant-usage-billing.md) | [prd-boss-tenant-usage-billing.md](../../tasks/modules/prd/boss/tenant/prd-boss-tenant-usage-billing.md) | [spec-boss-tenant-usage-billing.md](../../tasks/modules/spec/boss/tenant/spec-boss-tenant-usage-billing.md) |

## Console 对照

| BOSS 能力 | Console 对照 | 差异 |
|---|---|---|
| 租户列表 / 开通 | — | 平台 CRUD，Console 无对应页 |
| 配额策略 | — | 平台跨租户配额治理 |
| 租户管理员 | [tenant-management.md](../../console-modules/tenant/tenant-management.md) | BOSS 重置 tenant-admin；Console 邀请普通成员 |
| 租户计费与用量 | [usage-report.md](../../console-modules/tenant/usage-report.md) | BOSS 跨租户；Console 单租户 |

## 分层结论（当前）

- 平台租户资源：`Core /api/v1/tenants*` — **TODO-YAML**（产品规划见 ANI-02 §2.4，当前 `v1.yaml` 未声明 paths）
- 租户内成员/SSO：`Services /api/v1/svc/tenant/*` — 已冻结，**仅 Console** 使用，不得冒充 BOSS 平台租户 API
- 用量只读参考：`Core GET /api/v1/metering/usage` — 租户上下文；BOSS 跨租户聚合 **TODO-YAML**

## Phase 1 / 满配状态

**✅ 本域 4 模块均已满配（Core）文档**（2026-06-17）。主 API（`/tenants*`、quota、admin、平台 metering aggregate）仍为 **TODO-YAML**；`GET /api/v1/metering/usage` 为单租户只读参考。

满配检查清单：[`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)。
