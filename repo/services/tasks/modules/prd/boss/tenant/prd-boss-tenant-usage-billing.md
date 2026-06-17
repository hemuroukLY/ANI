# PRD: BOSS 租户计费与用量

## 1. 页面定位

跨租户计量分析；主维护源 [`tenant-usage-billing.md`](../../docs/boss-modules/tenant/tenant-usage-billing.md)。

## 2. 权威源结论

- `GET /api/v1/metering/usage` — YAML 已声明，**单租户 JWT 上下文**
- BOSS 跨租户 API — **TODO-YAML**
- `POST /metering/token-usage` — 上报，非查询

## 3. Goals

- 平台查看各租户 GPU-Hours、Token 等视角
- 钻取单租户时口径与 Console usage-report 一致
- 不把 UI 视角写成独立 API

## 4. 用户故事

### US-001: 跨租户排行

平台运营按时间范围查看租户用量 Top N。

### US-002: 指标视角切换

同页切换 GPU-Hours / Token 等，不新增 REST 路径。

### US-003: 钻取单租户

带 tenant 上下文打开明细，口径对齐 Console usage-report。

### US-004: 结算边界

账单/发票不在本页；不冒充 billing-export。

## 5. 功能需求

- FR-1：时间筛选 + 租户筛选 + 排行 + 趋势
- FR-2：平台 metering TODO-YAML
- FR-3：group_by 扩展 tenant_id 待 YAML

## 6. 非目标

- 租户 JWT 轮询冒充平台 API
- 发票与对账（Phase 2）

## 7. 成功标准

- 与 Console usage-report 字段口径一致
- 平台 API 待补标注完整

## 10. 文档同步状态

- 同步日期：**2026-06-17**
- 详文状态：**满配（Core）** — [`docs/boss-modules/tenant/tenant-usage-billing.md`](../../docs/boss-modules/tenant/tenant-usage-billing.md)
- SPEC：[`spec-boss-tenant-usage-billing.md`](../spec/spec-boss-tenant-usage-billing.md)
- HTML 摘要：[`ani-services-prototype-boss.html`](../../prototypes/ani-services-prototype-boss.html) · `boss.tenant.usage`
- 权威顺序：详文 > SPEC > PRD > HTML 摘要
