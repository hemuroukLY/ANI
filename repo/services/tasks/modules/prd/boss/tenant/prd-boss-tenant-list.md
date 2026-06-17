# PRD: BOSS 租户列表

## 1. 页面定位

BOSS / 租户与客户管理 / 租户列表，平台侧租户 CRUD 与生命周期入口。

主维护源：[`docs/boss-modules/tenant/tenant-list.md`](../../docs/boss-modules/tenant/tenant-list.md)

## 2. 权威源结论

- Core 规划：`ANI-02` §2.4 `/api/v1/tenants/` — 当前 **v1.yaml 无 paths**
- Services 已冻结：`/api/v1/svc/tenant/members` 等 — **Console 专用**，非 BOSS 租户创建
- 本轮 PRD 目标：收口产品边界，**不伪造** tenants API

## 3. Goals

- 平台运营可跨租户开通、停用、查看配额摘要
- 与 Console 租户管理职责零混淆
- 为 Core TODO-YAML（tenants CRUD）提供稳定产品基线

## 4. 用户故事

### US-001: 查看租户列表

作为平台运营，我希望筛选 active/suspended 租户并看到 GPU 配额与使用摘要，以便识别紧张租户。

验收：表格含 status、配额、last_active_at；未冻结 API 标注 TODO-YAML。

### US-002: 新建租户

作为平台运营，我希望通过向导创建租户并指定 tenant-admin 与初始配额。

验收：向导步骤完整；POST 须 idempotency_key；不得写成已实现 API。

### US-003: 生命周期管理

作为平台管理员，我希望停用/恢复租户并跳转配额或管理员页。

验收：状态×操作矩阵；suspended 与 active 转换前置条件写清。

### US-004: Core 边界

作为维护者，文档不得自造 TenantListResponse 或 `/api/v1/boss/tenants`。

## 5. 功能需求

- FR-1：跨租户列表 + 筛选 + 新建向导
- FR-2：行内操作：详情、配额、停用、恢复、重置管理员
- FR-3：Core tenants TODO-YAML 建议路径表
- FR-4：与 tenant-quota-policy、tenant-admin 跳转关系

## 6. 非目标

- 租户内成员邀请（Console）
- SSO 详细配置（Console PUT sso）
- 伪造已冻结 OpenAPI

## 7. 成功标准

- 详文、PRD、SPEC、HTML 摘要一致
- Phase 1 复审通过（ANI-14）

## 8. 开放问题

- deleted 租户是否允许恢复
- 创建租户时 SSO 是否同步写入 Services

## 10. 文档同步状态

- 同步日期：**2026-06-17**
- 详文状态：**满配（Core）** — [`docs/boss-modules/tenant/tenant-list.md`](../../docs/boss-modules/tenant/tenant-list.md)
- SPEC：[`spec-boss-tenant-list.md`](../spec/spec-boss-tenant-list.md)
- HTML 摘要：[`ani-services-prototype-boss.html`](../../prototypes/ani-services-prototype-boss.html) · `boss.tenant.list`
- 权威顺序：详文 > SPEC > PRD > HTML 摘要
