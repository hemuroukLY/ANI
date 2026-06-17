# 租户管理 — 成员详情

## 页面定位

单个租户成员的只读详情页，扩展成员列表。

本页属于 **Services / Tenant** 子能力。

## 文档管理规则

- 本文是成员详情子页的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 列表/邀请/移除见 `tenant-management.md`
- TASK：`TASK-SVC-015`

## Services 层要求

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/svc/tenant/members/{member_id}` | `getTenantMember` |

Schema：`TenantMember`（以 OpenAPI 为准）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 租户读权限 | 租户成员或管理员 | `403 FORBIDDEN` |
| 成员存在 | 当前租户内 | `404 NOT_FOUND` |

本页无写接口。**当前 YAML 未声明 `422`**。

## 页面职责

- 展示成员 profile、角色、加入时间等（以 `TenantMember` schema 为准）
- 提供跳转：角色权限（`role-permission-edit.md`）、移除成员（`tenant-management.md`）
- 编辑成员角色 — **YAML 未声明** PATCH member

## 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 查看成员详情 | 可用 | 可用 |
| 编辑角色 | 不可用 | 不可用（待补 PATCH） |
| 移除成员 | 跳转 tenant-management | 可用 |

## 接口冻结规则

### `GET /api/v1/svc/tenant/members/{member_id}`

- 成功：`200 + TenantMember`
- 错误：`401`、`403`、`404`
- 租户边界：不可跨租户读取

## 待补边界

- `PATCH /tenant/members/{member_id}` — **YAML 未声明**
- 成员活动审计 — 待补
- 与用户目录 `user-management.md` 合并方案 — TODO-YAML

## 相关模块

- `tenant-management.md`
- `role-permission-edit.md`
- `user-management.md`（独立 User CRUD 待规划）

## 验收标准

- [ ] GET member 路径与 services/v1.yaml 一致
- [ ] 未自造 PATCH member 为冻结能力
- [ ] handler stub ≠ 已实现
