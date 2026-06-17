# PRD: BOSS 配额策略

## 1. 页面定位

跨租户配额查看与调整；主维护源 [`tenant-quota-policy.md`](../../docs/boss-modules/tenant/tenant-quota-policy.md)。

## 2. 权威源结论

- Core tenants/quota — **TODO-YAML**
- Core 与 Services 配额维度可能拆分 — 须在 SPEC 标注来源层
- `/metering/usage` 仅单租户只读参考

## 3. Goals

- 平台统一治理 GPU/CPU/内存/存储及业务配额
- 调整配额须审计备注与影响说明
- 调低配额前置条件（不低于已用量）

## 4. 用户故事

### US-001: 配额总览

作为平台运营，我希望看到全平台配额使用率与超额租户列表。

### US-002: 调整配额

作为平台运营，我在抽屉中修改配额并填写审计备注；调低时若低于已用量应被拒绝（422 待 YAML）。

### US-003: 与租户列表联动

从租户列表「编辑配额」须携带 tenant_id 上下文进入本页。

## 5. 功能需求

- FR-1：配额表格 + 概览 + 调整抽屉
- FR-2：Core/Services 配额字段分层
- FR-3：操作矩阵与前置条件

## 6. 非目标

- Console 自助改配额
- 批量配额模板（Phase 2）

## 7. 成功标准

- 无伪造 quota schema
- 调整抽屉必填字段定义完整

## 10. 文档同步状态

- 同步日期：**2026-06-17**
- 详文状态：**满配（Core）** — [`docs/boss-modules/tenant/tenant-quota-policy.md`](../../docs/boss-modules/tenant/tenant-quota-policy.md)
- SPEC：[`spec-boss-tenant-quota-policy.md`](../spec/spec-boss-tenant-quota-policy.md)
- HTML 摘要：[`ani-services-prototype-boss.html`](../../prototypes/ani-services-prototype-boss.html) · `boss.tenant.quota`
- 权威顺序：详文 > SPEC > PRD > HTML 摘要
