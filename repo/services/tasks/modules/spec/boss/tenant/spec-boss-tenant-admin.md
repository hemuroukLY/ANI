# SPEC: BOSS tenant-admin

> Derived from: `tasks/modules/prd/boss/tenant/prd-boss-tenant-admin.md`  
> Revised: 2026-06-17

## 1. Summary

平台 tenant-admin 治理与 Console members API 边界。

主维护源：`docs/boss-modules/tenant/tenant-admin.md`

## 2. Frozen Facts

### 2.1 BOSS Platform API

**当前无冻结路径。**

### 2.2 Suggested TODO-YAML（Core）

| 方法 | 建议路径 |
|---|---|
| GET | `/api/v1/tenants/{tenant_id}/admin` |
| POST | `/api/v1/tenants/{tenant_id}/admin:reset` |
| PUT | `/api/v1/tenants/{tenant_id}/admin` |

### 2.3 Console Reference（非 reset-admin）

| 方法 | 路径 | operationId | 语义 |
|---|---|---|---|
| POST | `/api/v1/svc/tenant/members` | `inviteTenantMember` | 邀请成员 |
| GET | `/api/v1/svc/tenant/members` | `listTenantMembers` | 成员列表 |

**禁止**将 invite 映射为 BOSS reset-admin。

### 2.4 Core 同步边界（2026-06-17）

- 平台 tenant-admin 操作 — **TODO-YAML P1**
- 不得把 Console `inviteTenantMember` 写成 BOSS 正式契约
- 403 示例为 `permission denied` 参考

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。
## 3. Security

- 重置须 idempotency_key + 审计
- 敏感操作 UI 二次确认（产品层）

## 4. 前置条件与矩阵

见主维护源。

## 5. 待补

- Core admin 子资源
- 会话吊销 Phase 2
