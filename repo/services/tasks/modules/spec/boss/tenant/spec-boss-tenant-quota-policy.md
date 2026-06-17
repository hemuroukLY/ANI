# SPEC: BOSS tenant-quota-policy

> Derived from: `tasks/modules/prd/boss/tenant/prd-boss-tenant-quota-policy.md`  
> Revised: 2026-06-17

## 1. Summary

配额策略的 Core/Services 分层与 TODO-YAML 建议。

主维护源：`docs/boss-modules/tenant/tenant-quota-policy.md`

## 2. Frozen Facts

### 2.1 Authority

- `v1.yaml` — 无 quota paths
- `GET /api/v1/metering/usage` — 单租户用量只读参考

### 2.2 Suggested TODO-YAML

| 方法 | 建议路径 |
|---|---|
| GET | `/api/v1/tenants/{tenant_id}/quota` |
| PATCH | `/api/v1/tenants/{tenant_id}/quota` |

### 2.3 Quota Field Layers

| 字段 | 建议归属 |
|---|---|
| max_gpu_count, max_cpu, max_memory_gb, max_storage_gb | Core |
| max_inference_services, max_knowledge_bases | Services（待确认） |
| rate_limit_rps | Core Gateway（待确认） |

### 2.4 Core 同步边界（2026-06-17）

- 平台 quota CRUD / 跨租户 list — **TODO-YAML P1**
- 平台 RBAC 待 YAML；403 示例为 `permission denied` 参考

> 详文满配（Core）见主维护源；**路径已声明 ≠ RBAC 已冻结**；OpenAPI 已声明 ≠ handler 已实现。
## 3. Write Contract（待 YAML）

- PATCH 需 `idempotency_key`
- 成功：200/204（待冻结）
- 调低低于已用量：422（待 YAML 声明）

## 4. 前置条件与矩阵

见主维护源。

## 5. 待补

- Core quota CRUD
- 审计事件联动 audit 域
