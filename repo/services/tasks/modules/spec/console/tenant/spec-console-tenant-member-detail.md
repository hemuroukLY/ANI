# SPEC: Console tenant-member-detail

> Source: `tasks/modules/prd/console/tenant/prd-console-tenant-member-detail.md`  
> Revised: 2026-06-17

## 1. Summary

单个租户成员的只读详情页，扩展成员列表。属于 **Services / Tenant** 子能力。列表/邀请/移除见 `tenant-management.md`。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | 权限 |
|---|---|---|---|---|
| GET | `/api/v1/svc/tenant/members/{member_id}` | `getTenantMember` | `200 + TenantMember` | 租户成员或管理员 |

### 2.3 Verified Schemas

- `TenantMember`（profile、角色、加入时间等，以 OpenAPI 为准）

## 3. Page Scope

- 展示成员 profile、角色、加入时间等（以 `TenantMember` schema 为准）
- 提供跳转：角色权限（`role-permission-edit.md`）、移除成员（`tenant-management.md`）

## 4. Non-Goals

- `PATCH /tenant/members/{member_id}` — **YAML 未声明**
- 成员活动审计 — 待补
- 与用户目录 `user-management.md` 合并方案 — TODO-YAML
- 本页无写接口；**当前 YAML 未声明 `422`**

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 租户读权限 | 租户成员或管理员 | `403 FORBIDDEN` |
| 成员存在 | 当前租户内 | `404 NOT_FOUND` |

## 6. 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 查看成员详情 | ✅ | ✅ |
| 编辑角色 | ❌ | ❌（待补 PATCH） |
| 移除成员 | 跳转 tenant-management | ✅ |

## 7. 主维护源

- `docs/console-modules/tenant/tenant-member-detail.md`
- 相关：`tenant-management.md`、`role-permission-edit.md`、`user-management.md`
- TASK：`TASK-SVC-015`

## 8. Handler 验收（Services 团队）

```bash
curl -H "Authorization: Bearer $TOKEN" "$BASE/api/v1/svc/tenant/members/{member_id}"
```

- 租户边界：不可跨租户读取
- 未自造 PATCH member 为冻结能力
- OpenAPI 已声明 ≠ handler 已实现
