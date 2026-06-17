# 租户管理

## 页面定位

`租户管理` 是 `Console / 设置 / 安全与身份` 域下的租户侧管理页面，用于成员邀请、角色查看、SSO 配置与 Webhook 管理。

当前模块属于 `Services / Tenant`，一级权威源为 `ANI-main/repo/api/openapi/services/v1.yaml`，正式路径前缀为 `/api/v1/svc/tenant/*`。

## 文档管理规则

- 本文件是 `租户管理` 的主维护源
- `tasks/modules/prd/console/tenant/prd-console-tenant-management.md` 与 `tasks/modules/spec/console/tenant/spec-console-tenant-management.md` 作为辅助材料
- 与 `security-identity-overview.md` 分工：安全总入口只做跳转，本页承接 `/tenant/*` 正式契约

## Services 层要求

- 租户管理属于 `Services`，不走 Core `/api/v1/auth/*`（API Key 除外）
- 写操作必须含 `idempotency_key`（POST/PUT/PATCH）
- 错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- 不要求前端显式传 `tenant_id`

## 已冻结路径（services/v1.yaml）

| 能力 | 方法 | 路径 | operationId |
|---|---|---|---|
| 成员列表 | GET | `/api/v1/svc/tenant/members` | `listTenantMembers` |
| 邀请成员 | POST | `/api/v1/svc/tenant/members` | `inviteTenantMember` |
| 移除成员 | DELETE | `/api/v1/svc/tenant/members/{member_id}` | `removeTenantMember` |
| 角色列表 | GET | `/api/v1/svc/tenant/roles` | `listTenantRoles` |
| SSO 查询 | GET | `/api/v1/svc/tenant/sso` | `getTenantSSOConfig` |
| SSO 更新 | PUT | `/api/v1/svc/tenant/sso` | `updateTenantSSOConfig` |
| Webhook 列表 | GET | `/api/v1/svc/tenant/webhooks` | `listWebhooks` |
| Webhook 创建 | POST | `/api/v1/svc/tenant/webhooks` | `createWebhook` |
| Webhook 删除 | DELETE | `/api/v1/svc/tenant/webhooks/{webhook_id}` | `deleteWebhook` |

## 创建前置条件

| 场景 | 依赖项 | 未满足时的 HTTP 响应 |
|---|---|---|
| 邀请成员 | 操作者为租户管理员 | `403 FORBIDDEN` |
| 邀请成员 | 邮箱未被占用 | `409 CONFLICT` |
| 更新 SSO | OIDC Provider 配置合法 | `400 BAD_REQUEST` |
| 创建 Webhook | URL 可达且 secret 已配置 | `422 PRECONDITION_FAILED`（**`createWebhook` YAML 已声明**，Services v1.yaml Phase 2 2026-06-17；具体 `code` 待 description 补充） |

## 操作可用性矩阵

| 操作 | 普通成员 | 租户管理员 | 说明 |
|---|---|---|---|
| 查看成员列表 | ✅ | ✅ | GET members |
| 邀请成员 | ❌ | ✅ | POST members → `202` |
| 移除成员 | ❌ | ✅ | DELETE member → `204` |
| 查看角色 | ✅ | ✅ | GET roles |
| 查看/更新 SSO | ❌ | ✅ | GET/PUT sso |
| 管理 Webhook | ❌ | ✅ | GET/POST/DELETE webhooks |

## 待补边界

- 角色编辑（当前只读列表）<!-- ADDED-TO-YAML: PUT /api/v1/svc/tenant/roles/{role_id} (Services v1.yaml, Phase 2 2026-06-17) --> — 详文 `role-permission-edit.md`
- 成员详情、Webhook 投递日志 — 详文 `tenant-member-detail.md`、`tenant-webhook-ops.md`（Phase 2 YAML）
- 审计与合规导出

## 验收标准

- 路径与 `services/v1.yaml` 一致
- 不与 Core `/api/v1/auth/api-keys*` 重复定义 API Key 能力
- HTML 摘要、PRD、SPEC 与本文件一致
