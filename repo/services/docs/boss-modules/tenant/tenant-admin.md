# 租户管理员

## 页面定位

`租户管理员` 是 BOSS **租户与客户管理** 域下的 **平台侧 tenant-admin 账号治理** 页：查看、重置、轮换各租户的初始/主管理员凭据，供平台运营在客户失联或安全事件时紧急处置。

本页属于 **Core / Tenants + Admin** 视角下的 **平台 RBAC** 页面。平台 reset-admin API 规划随 `tenants` 子资源；**ADDED-TO-YAML**（`/platform/tenant-admins`、`/tenants/{id}/admin*`）。

Console [`tenant-management.md`](../../console-modules/tenant/tenant-management.md) 的 `inviteTenantMember` 面向 **租户内协作成员**；本页面向 **平台开通时创建的 tenant-admin** 及 **平台发起的凭据重置**，**不得** 混用 API 语义。

## 文档管理规则

- 本文是 `租户管理员` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- [`prd-boss-tenant-admin.md`](../../tasks/modules/prd/boss/tenant/prd-boss-tenant-admin.md) 与 [`spec-boss-tenant-admin.md`](../../tasks/modules/spec/boss/tenant/spec-boss-tenant-admin.md) 为辅助材料，不替代本文
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- 平台级管理员重置归属 **Core**（建议 `GET/POST /api/v1/tenants/{tenant_id}/admin*`）
- **不得** 用 `POST /api/v1/svc/tenant/members`（`inviteTenantMember`）冒充 BOSS reset-admin
- 重置/轮换须 `idempotency_key` + 平台 RBAC；须产生 **审计事件**（audit **TODO-YAML**）
- 敏感操作须 UI **二次确认**（产品行为；非 OpenAPI 字段）
- 统一错误结构；`422` 仅 YAML 已声明 operation 可写冻结（§2.10）
- 禁止自造 `ResetTenantAdminRequest` 为 **已冻结** schema（可写目标语义）
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **跨租户** tenant-admin list 与筛选（**ADDED-TO-YAML** `listTenantAdmins`）
- 提供 **重置密码**、**更换管理员** 抽屉（风险说明、新邮箱、通知方式、审计备注）
- 支持从 [`tenant-list.md`](tenant-list.md) 带 `tenant_id` 深链
- 可选 **只读跳转** Console 成员 list（须租户上下文 — BOSS 默认不代调用）
- 写入 [`platform-audit-log.md`](../audit/platform-audit-log.md) 审计链（API 待 YAML）
- **不承担** 邀请普通成员（Console POST members）

## 页面结构

- 首屏至少包含：`筛选区`、`管理员表格`、`边界说明`
- 重置/更换为 drawer；高风险须二次确认 modal
- API 未就绪时：只读占位 + 操作 disabled

```text
租户管理员
├── 筛选（租户 / status / 最近重置时间）
├── 管理员表格
│   ├── tenant / admin_email / role / last_login / MFA
│   └── 操作：重置密码 / 更换管理员 / 查看成员（Console 深链）
└── 重置/更换 drawer
    ├── 风险说明
    ├── 新邮箱或用户名
    ├── 通知方式（邮件 / 工单号）
    ├── audit_note（必填）
    └── 提交 POST/PUT + idempotency_key
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 / 来源 | 本页用法 |
|---|---|---|
| Core | tenants/admin **TODO-YAML** | reset/rotate 正式 API |
| Core | tenants list **TODO-YAML** | 表格 tenant 上下文 |
| Services | `GET /api/v1/svc/tenant/members` | **Console**；BOSS 仅深链参考 |
| Core | auth 会话吊销 | Phase 2 · **待 YAML** |
| Core | platform-audit-log | 操作审计 **TODO-YAML** |

### 关键边界

- `inviteTenantMember`：`InviteMemberRequest` 必填 `idempotency_key`（uuid）、`email`、`role` — 语义为 **邀请**，**不是** reset-admin
- `TenantMember` schema：`id`, `user_id`, `email`, `role`, `joined_at` — **无** platform reset 字段
- reset-admin **不得** 返回明文密码于 API 响应（产品：邮件/一次性链接 · 待 YAML）
- MFA 状态字段 **待 YAML**；合入前 UI 显示「待接入」

## 页面区块与数据来源映射

| 区块 | 主要来源 | 说明 | 跳转 |
|---|---|---|---|
| 筛选区 | admin list API **待 YAML** | tenant/status/time | 刷新表格 |
| 管理员表格 | GET admin list **待 YAML** | 跨租户主管理员 | 打开 drawer |
| 重置 drawer | POST admin:reset **待 YAML** | idempotency_key | audit |
| 更换 drawer | PUT admin **待 YAML** | 新 email 校验 | audit |
| 查看成员 | Console | 租户 JWT 上下文 | tenant-management |
| 边界说明 | 规划项 | invite ≠ reset | — |

## BOSS 与 Console 分工

| 操作 | BOSS 租户管理员 | Console 租户管理 |
|---|---|---|
| 重置 tenant-admin 密码 | ✅（待 YAML） | ❌ |
| 更换 tenant-admin 绑定用户 | ✅（待 YAML） | ❌ |
| 邀请普通成员 | ❌ | ✅ POST members |
| 查看成员 list | 深链 optional | ✅ GET members |
| SSO 绑定 | ❌ | ✅ PUT sso |
| 强制下线会话 | Phase 2 | ❌ |
| RBAC | 平台 tenant admin write | 租户管理员 JWT |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| — | `/api/v1/tenants/*/admin*` | — | **未冻结** |

| Services 参考（**非 reset-admin**） | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/svc/tenant/members` | `listTenantMembers` |
| POST | `/api/v1/svc/tenant/members` | `inviteTenantMember` |

| Schema | 说明 |
|---|---|
| `InviteMemberRequest` | 邀请成员；**≠** reset-admin |
| `TenantMember` | 成员记录；Console 只读参考 |

| 能力 | 状态 |
|---|---|
| 平台 reset/rotate admin | **ADDED-TO-YAML** |
| MFA 状态查询 | **TODO-YAML** Phase 2 |
| 会话吊销 | Phase 2 |

### 建议 TODO-YAML（非冻结）

| 能力 | 建议路径 | 归属 |
|---|---|---|
| 查询租户主管理员 | `GET /api/v1/tenants/{tenant_id}/admin` | Core |
| 重置管理员凭据 | `POST /api/v1/tenants/{tenant_id}/admin:reset` | Core |
| 指定新管理员 | `PUT /api/v1/tenants/{tenant_id}/admin` | Core |

## 字段级定义

### 查询字段（admin list · **ADDED-TO-YAML**）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `tenant_id` | query | ❌ | 筛选 |
| `status` | query | ❌ | 租户 status |
| `reset_after` | query | ❌ | 最近重置时间下限 |
| `limit` / `cursor` | query | ❌ | 分页 |

### 表格字段（页面目标 · **ADDED-TO-YAML**）

| 字段 | 说明 |
|---|---|
| `tenant_id` / `display_name` | 所属租户 |
| `admin_user_id` | 管理员用户 ID |
| `admin_email` | 登录邮箱 |
| `role` | tenant-admin 或等价角色 |
| `last_login_at` | 最近登录 |
| `mfa_enabled` | MFA 是否启用 **待 YAML** |
| `last_reset_at` | 最近平台重置时间 |
| `reset_by` | 平台操作者 |

### 重置/更换 — 请求字段（**ADDED-TO-YAML**）

| 字段 | 必填 | 说明 |
|---|---|---|
| `idempotency_key` | ✅ | 幂等 |
| `notification_channel` | ❌ | email / ticket_id |
| `new_email` | 更换时 ✅ | 须合法且未被占用 |
| `audit_note` | ✅ | 审计备注 |
| `force_logout` | ❌ | Phase 2 |

### 展示字段（UI 计算）

| 字段 | 说明 |
|---|---|
| `days_since_reset` | 距 last_reset_at |
| `risk_level` | 长期未登录 + 无 MFA → 高（UI 规则） |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| API 未就绪 | 表格占位 + 操作 disabled | **不** 伪造管理员数据 |
| 重置成功 | Toast + 更新 last_reset_at | 不展示明文密码 |
| 更换成功 | 刷新 admin_email 列 | — |
| 租户 deleted | 行内操作 disabled | 422/403 |
| 无写权限 | 隐藏重置/更换 | 403 |
| MFA 未接入 | 列显示「待 YAML」 | — |
| 二次确认取消 | 关闭 modal，不提交 | UI |

## 字段口径与单位

| 字段 | 口径 | 格式 |
|---|---|---|
| `admin_email` | email | string |
| `last_login_at` / `last_reset_at` | ISO 8601 | date-time |
| `role` | 固定 tenant-admin 展示 | string |
| `mfa_enabled` | boolean **待 YAML** | 是/否 |

## 状态与能力口径

本页 **不** 管理租户 `status`；deleted 租户禁止 reset（前置 422 **待 YAML**）。

| 操作 | 说明 |
|---|---|
| 重置密码 | 触发凭据轮换 + 通知 |
| 更换管理员 | 绑定新 user/email |
| 邀请成员 | **不在本页** → Console |

强制下线 — Phase 2（auth 联动 **待 YAML**）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 tenant admin 读 RBAC | 已授权（**TODO-YAML** admin scope 待 Core 定义） | `403 FORBIDDEN` |
| 平台 tenant admin 写 RBAC | 已授权（**TODO-YAML**） | `403 FORBIDDEN` |
| `idempotency_key` | 非空 | `400 BAD_REQUEST` |
| `audit_note` | 非空 | UI 拦截；`400` **待 YAML** |
| 租户非 deleted | status 允许 | `422` **待 YAML** |
| 更换管理员 | new_email 合法且未占用 | `409` / `400` |

## 操作可用性矩阵

| 操作 | 平台只读 | 平台运营 | 平台管理员 | 说明 |
|---|---|---|---|---|
| 查看管理员 list | ✅（待 YAML） | ✅ | ✅ | GET |
| 重置密码 | ❌ | ✅ | ✅ | POST reset |
| 更换管理员 | ❌ | ✅ | ✅ | PUT admin |
| 强制下线 | ❌ | Phase 2 | Phase 2 | auth |
| 邀请普通成员 | ❌ | ❌ | ❌ | Console |
| 查看 Console 成员 | 深链 | 深链 | 深链 | 租户 JWT |

## 删除前置校验与当前契约边界

本页 **无 DELETE 管理员**；「移除管理员」通过更换绑定用户实现（PUT **待 YAML**）。

**不得** 使用 `DELETE /svc/tenant/members/{member_id}` 作为平台 reset-admin 的替代。

## 接口冻结规则

### `POST /api/v1/svc/tenant/members`（Services · **非 reset-admin**）

- `operationId`：`inviteTenantMember`
- `body.required`：`idempotency_key`, `email`, `role`（`InviteMemberRequest`）
- `success`：`202 + TenantMember`
- `errors`：`400`、`409`
- **不得** 写成 reset-admin alias

### Core admin 子资源（**ADDED-TO-YAML**）

<!-- TODO-YAML: GET/POST/PUT /api/v1/tenants/{tenant_id}/admin* -->

- POST reset 须 `idempotency_key` + `audit_note`
- 合入前不得写入「已冻结」表

## 使用规则

- **禁止** 把 `inviteTenantMember` 绑定到「重置管理员」按钮
- **禁止** 在 API 响应或 UI 展示明文新密码
- reset/rotate **必须** 写 audit（API 待 YAML）
- 高风险操作 **必须** UI 二次确认
- BOSS **不得** 默认携带租户 JWT 调用 `listTenantMembers` 冒充平台 list

## 待补能力边界

- Core admin 子资源 API — **ADDED-TO-YAML**
- 与 Core `/api/v1/auth/*` 会话吊销 — Phase 2
- MFA 强制策略 — Phase 2
- 管理员 list — **ADDED-TO-YAML** `listTenantAdmins`

## 响应示例

### inviteTenantMember（Services · **非本页 · 仅边界对照**）

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440003",
  "user_id": "990e8400-e29b-41d4-a716-446655440004",
  "email": "collab@acme.example",
  "role": "member",
  "joined_at": "2026-06-16T12:00:00Z"
}
```

### reset-admin 目标响应（**待 YAML · 非已冻结**）

```json
{
  "tenant_id": "ten-acme",
  "admin_email": "admin@acme.example",
  "last_reset_at": "2026-06-16T14:00:00Z",
  "reset_by": "platform-op-bob",
  "notification_sent": true,
  "audit_id": "aud-20260616-002"
}
```

## 错误示例

### 无 platform admin 写权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-admin-403-001"
}
```

> **注**：适用于 **admin reset API（TODO-YAML）** 场景。`message` 为示例文案，**非** YAML 已冻结 scope 名。

### 缺少 idempotency_key（目标契约）

```json
{
  "code": "BAD_REQUEST",
  "message": "idempotency_key is required",
  "request_id": "req-boss-admin-400-001"
}
```

### 租户已删除

```json
{
  "code": "PRECONDITION_FAILED",
  "message": "cannot reset admin for deleted tenant",
  "request_id": "req-boss-admin-422-001"
}
```

## 相关模块

- [tenant-list.md](tenant-list.md)
- [platform-audit-log.md](../audit/platform-audit-log.md)
- Console：[tenant-management.md](../../console-modules/tenant/tenant-management.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 与 Console members / invite API 边界清晰
- [x] 无伪造 admin API 为已冻结
- [x] 含响应示例与错误示例（400 + 403 + 422）
- [x] 独立字段定义（查询/表格/请求/展示）
- [ ] admin YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
