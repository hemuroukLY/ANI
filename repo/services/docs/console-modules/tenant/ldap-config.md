# LDAP 目录对接

## 页面定位

企业 LDAP/AD 与租户身份目录的对接配置，区别于 OIDC SSO 登录流。

## 文档管理规则

- 本文是 LDAP 配置规划主维护源
- 与 `security-identity-overview.md` SSO 分流口径一致
- 不得把 SSO OIDC 配置写成 LDAP 已对齐

## Services 层要求

### 已有（SSO / OIDC）

- `GET /api/v1/svc/tenant/sso` — `getSsoConfig`
- `PUT /api/v1/svc/tenant/sso` — `updateSsoConfig`（含 `idempotency_key`）
- 响应：`SsoConfig` / `UpdateSsoConfigRequest`
- 错误：GET `401`/`403`；PUT `400`/`403`

### 待补（LDAP 专有）

<!-- TODO-YAML: LDAP bind URL、base DN、sync 策略、group mapping 等 -->

- **无**独立 `/tenant/ldap` 或 `/auth/ldap` 冻结路径

## Core 层要求

- Core `POST /api/v1/auth/oidc/begin` 等 — 登录流，非 LDAP 目录配置
- LDAP bind 若未来走 Core Auth，须先扩 YAML

## 页面职责

- 规划说明：LDAP vs OIDC 分工
- 配置表单占位 + TODO-YAML 标注
- OIDC 配置跳转 tenant/sso（已冻结）

## 创建前置条件

| 场景 | 要求 | 未满足响应 |
|---|---|---|
| 查看 OIDC SSO | 租户管理员 | `403` |
| 更新 OIDC SSO | `idempotency_key` + 合法 body | `400`/`403` |
| LDAP 配置保存 | **待 YAML** | — |

## 操作可用性矩阵

| 操作 | 普通成员 | 租户管理员 |
|---|---|---|
| 查看 LDAP 页（规划） | ❌ | ✅（占位） |
| 编辑 OIDC SSO | ❌ | ✅ |
| 编辑 LDAP | ❌ | ❌（待 YAML） |
| 触发 LDAP 同步 | ❌ | ❌（待 YAML） |

## 接口冻结规则

### `GET /api/v1/svc/tenant/sso`

- 成功：`200 + SsoConfig`
- 错误：`401`、`403`

### `PUT /api/v1/svc/tenant/sso`

- 成功：`200 + SsoConfig`
- 错误：`400`、`403`
- **当前 YAML 未声明** `401`、`404`、`422`

### LDAP 专有 API

- **无冻结 operation**

## 待补边界

- `/tenant/ldap` GET/PUT — **TODO-YAML**
- LDAP group → TenantRole 映射 — **TODO-YAML**
- 与 `user-management.md` User 资源关系

## 验收标准

- [ ] 与 security-identity-overview SSO 分流口径一致
- [ ] LDAP 能力标注 TODO-YAML，不冒充已实现
- [ ] OIDC 路径与 services/v1.yaml 一致
