# PRD: BOSS 租户管理员

## 1. 页面定位

平台 tenant-admin 重置与更换；主维护源 [`tenant-admin.md`](../../docs/boss-modules/tenant/tenant-admin.md)。

## 2. 权威源结论

- 平台 reset-admin — **TODO-YAML** Core
- `inviteTenantMember` — Console 邀请，**不等于** reset-admin

## 3. Goals

- 紧急重置 tenant-admin 凭据
- 更换主管理员账号
- 全操作可审计

## 4. 用户故事

### US-001: 查看各租户管理员

平台运营查看 admin_email、last_login、last_reset。

### US-002: 重置密码

二次确认 + idempotency_key + 审计备注；租户 deleted 时不可操作。

### US-003: 与 Console 分工

普通成员邀请只在 Console；BOSS 页不提供 invite 入口。

## 5. 功能需求

- FR-1：管理员表格 + 重置/更换抽屉
- FR-2：与 tenant-list「重置管理员」联动
- FR-3：禁止混用 Services members API

## 6. 非目标

- MFA 策略配置（Phase 2）
- 会话强制下线（Phase 2）

## 7. 成功标准

- 与 Console tenant-management 边界清晰

## 10. 文档同步状态

- 同步日期：**2026-06-17**
- 详文状态：**满配（Core）** — [`docs/boss-modules/tenant/tenant-admin.md`](../../docs/boss-modules/tenant/tenant-admin.md)
- SPEC：[`spec-boss-tenant-admin.md`](../spec/spec-boss-tenant-admin.md)
- HTML 摘要：[`ani-services-prototype-boss.html`](../../prototypes/ani-services-prototype-boss.html) · `boss.tenant.admin`
- 权威顺序：详文 > SPEC > PRD > HTML 摘要
