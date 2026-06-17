# 用户管理

## 页面定位

`用户管理` 指 **平台/租户用户目录** 的独立 CRUD 页，与 `tenant-management.md` 的 **成员邀请**（`/tenant/members`）不同。

## 文档管理规则

- 成员列表/邀请见 `tenant-management.md`
- 本文仅描述「用户对象」级管理边界
- 不得把 members 列表冒充全局 User 管理

## Services 层要求

**当前无**独立 `GET /api/v1/svc/users` 或 Core `/users*` 冻结路径。

<!-- TODO-YAML: 用户 CRUD 与 members 关系待产品定稿（合并 vs 独立 User 资源） -->

### 当前可引用能力

| 能力 | 路径 | 说明 |
|---|---|---|
| 租户成员 | `/api/v1/svc/tenant/members*` | 非完整 User CRUD |
| OIDC 登录用户 | Core `/api/v1/auth/*` | 认证，非用户目录 |

## Core 层要求

- 本页不定义新的 Core User 资源（待 YAML）
- 认证相关只读能力不替代用户目录

## 页面职责

- 规划态：说明与成员管理的分工
- 占位 UI + 明确 TODO-YAML
- 跳转 `tenant-management.md` 成员入口

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 租户管理员 | 管理成员权限 | `403 FORBIDDEN` |

独立 User CRUD 的 POST/PUT 前置条件 — **待 YAML 冻结后补充**（含 `idempotency_key`）。

## 操作可用性矩阵

| 操作 | 普通成员 | 租户管理员 |
|---|---|---|
| 查看成员（跳转 tenant） | ✅ | ✅ |
| 独立 User CRUD | ❌ | ❌（待 YAML） |
| LDAP 同步用户 | ❌ | ❌（见 ldap-config） |

## 接口冻结规则

**当前无本模块冻结 API。**

可引用（非 User 资源）：

### `GET /api/v1/svc/tenant/members`（若存在）

- 以 `services/v1.yaml` Tenant 组为准
- 用途：成员列表，非 User 对象 CRUD

## 待补边界

- 独立 User 资源 vs Members 合并方案 — **TODO-YAML**
- 用户禁用/删除与 SSO 身份关联 — 待产品定稿
- 与 `role-permission-edit.md` 的权限边界

## 验收标准

- [ ] TODO-YAML 明确
- [ ] 与安全域 HTML「用户管理（待补）」一致
- [ ] 不自造 User schema 为冻结事实
