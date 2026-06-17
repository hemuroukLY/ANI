# SPEC: BOSS tenant-list

> Derived from: `tasks/modules/prd/boss/tenant/prd-boss-tenant-list.md`  
> Revised: 2026-06-17

## 1. Summary

收口 BOSS 租户列表的 Core 分层、TODO-YAML 建议与 Console 边界。

主维护源：`docs/boss-modules/tenant/tenant-list.md`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/openapi/services/v1.yaml`（对照用）
- `ANI-02` §2.4（产品规划 `/api/v1/tenants/`）

### 2.2 Verified Paths（BOSS 正式）

**当前无。** 平台 tenants CRUD 未在 v1.yaml 声明。

### 2.3 Console-only Paths（不得冒充 BOSS tenants API）

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/svc/tenant/members` | `listTenantMembers` |
| POST | `/api/v1/svc/tenant/members` | `inviteTenantMember` |

### 2.4 Core 同步边界（2026-06-17）

- Core `/api/v1/tenants*` **未声明** — BOSS 正式 API 全部 **TODO-YAML P1**
- 平台 RBAC 鉴权待 YAML；**不得**信任未授权 `tenant_id`
- 403 示例仅作 `permission denied` 文案参考，**非**已冻结 scope 名

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。

### 2.5 Suggested TODO-YAML（Core）

| 方法 | 建议路径 | 用途 |
|---|---|---|
| GET | `/api/v1/tenants` | 列表 |
| POST | `/api/v1/tenants` | 创建 |
| GET | `/api/v1/tenants/{tenant_id}` | 详情 |
| PATCH | `/api/v1/tenants/{tenant_id}` | 状态/资料 |

## 3. Architecture Constraints

- Core 前缀 `/api/v1/*`
- POST/PATCH 需 `idempotency_key`
- 平台 RBAC；禁止信任前端 tenant_id 越权
- 422 按 module-delivery-workflow §2.10

## 4. Page Scope

- 列表、向导、生命周期、跳转 quota/admin
- Non-goals：成员邀请、SSO 详情

## 5. 前置条件与矩阵

见主维护源 §创建前置条件、§操作可用性矩阵、§状态与能力口径。

## 6. 待补边界

- Core tenants 全量 CRUD
- inference/kb 计数聚合

OpenAPI 已声明 ≠ handler 已实现。
