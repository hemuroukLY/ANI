# 角色与权限编辑

## 页面定位

租户角色权限的查看与编辑，扩展 `tenant-management.md` 只读角色列表。

## 文档管理规则

- 本文是角色权限编辑子模块主维护源
- 一级权威源：`services/v1.yaml` Tenant 组

## Services 层要求

<!-- ADDED-TO-YAML: PUT /api/v1/svc/tenant/roles/{role_id} (Phase 2) -->

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/svc/tenant/roles` | `listTenantRoles` |
| PUT | `/api/v1/svc/tenant/roles/{role_id}` | `updateTenantRole` |

请求：`UpdateTenantRoleRequest`（必填 `idempotency_key`、`permissions[]`，`permissions` minItems ≥ 1）。

响应：`200 + TenantRole`。

## 页面职责

- 角色列表 + 权限矩阵编辑
- 保存后展示 `TenantRole` 更新结果
- 不提供角色创建/删除（**当前 YAML 未声明** POST/DELETE roles）

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 租户管理员 | 更新角色权限 | `403 FORBIDDEN` |
| `role_id` | 存在 | `404 NOT_FOUND` |
| PUT 请求体 | `idempotency_key` + 非空 `permissions[]` | `400 BAD_REQUEST` |

**当前 YAML 未在 update 声明 `422`**。

## 操作可用性矩阵

| 操作 | 普通成员 | 租户管理员 |
|---|---|---|
| 查看角色 | ✅ | ✅ |
| 编辑权限 | ❌ | ✅ |
| 创建/删除角色 | ❌ | ❌（待 YAML） |

## 接口冻结规则

### `GET /api/v1/svc/tenant/roles`

- 成功：`200` + items（`TenantRole[]`）
- 错误：`401`
- **当前 YAML 未声明** `403`

### `PUT /api/v1/svc/tenant/roles/{role_id}`

- 成功：`200 + TenantRole`
- 错误：`400`、`401`、`403`、`404`

## 待补边界

- 角色创建 POST / 删除 DELETE — **TODO-YAML**
- 系统内置角色不可编辑 — 产品规则，待 RBAC 文档
- permission 枚举权威清单 — 以 Core/Services scope 注册表为准

## 验收标准

- [ ] 与 tenant-management 待补边界对齐
- [ ] PUT 必填 `idempotency_key`
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
