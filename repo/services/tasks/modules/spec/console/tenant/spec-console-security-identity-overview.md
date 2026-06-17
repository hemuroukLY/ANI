# SPEC: Console 安全与身份概览

> Technical specification derived from: `tasks/modules/prd/console/tenant/prd-console-security-identity-overview.md`
> Source of truth: `ANI-main/repo/api/openapi/v1.yaml`
> Main doc: `docs/console-modules/tenant/security-identity-overview.md`
> Scope: `Console / 安全与身份 / 安全与身份概览`

## 1. Summary

本 SPEC 定义 `Console / 安全与身份 / 安全与身份概览` 的技术边界、页面承接范围、Auth 冻结路径和总入口约束。该页面属于租户侧安全域总入口，直接消费 `Core / Auth` 能力，不扩写新的 `Services` 资源，也不扩写配置中心或审计中心。

## 2. Source of Truth

### 2.1 Primary Authority

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Frozen Path Group

- `/auth/oidc/begin`
- `/auth/token`
- `/auth/refresh`
- `/auth/logout`
- `/auth/api-keys`
- `/auth/api-keys/{key_id}`

### 2.3 Frozen Schemas

- `BeginOIDCLoginRequest`
- `BeginOIDCLoginResponse`
- `CompleteOIDCLoginRequest`
- `TokenPairResponse`
- `RefreshAccessTokenRequest`
- `RefreshAccessTokenResponse`
- `LogoutRequest`
- `RevokeStatusResponse`
- `CreateAPIKeyRequest`
- `CreateAPIKeyResponse`
- `APIKeyInfo`
- `ListAPIKeysResponse`
- `ErrorResponse`

## 3. API Design

### 3.1 Endpoint Matrix

| Method | Path | operationId | summary | Success | Error Responses |
|---|---|---|---|---|---|
| POST | `/api/v1/auth/oidc/begin` | `beginOIDCLogin` | 发起 OIDC 登录 | `200 + BeginOIDCLoginResponse` | `400` |
| POST | `/api/v1/auth/token` | `completeOIDCLogin` | OIDC callback 换取 TokenPair | `200 + TokenPairResponse` | `400`、`401` |
| POST | `/api/v1/auth/refresh` | 当前权威源未声明 | 刷新 AccessToken | `200 + RefreshAccessTokenResponse` | `400`、`401` |
| POST | `/api/v1/auth/logout` | `logout` | 吊销当前 JWT JTI | `200 + RevokeStatusResponse` | `400`、`401`、`403` |
| GET | `/api/v1/auth/api-keys` | `listAPIKeys` | 列出当前租户 API Key | `200 + ListAPIKeysResponse` | `401`、`403` |
| POST | `/api/v1/auth/api-keys` | `createAPIKey` | 创建 API Key | `201 + CreateAPIKeyResponse` | `400`、`401`、`403` |
| DELETE | `/api/v1/auth/api-keys/{key_id}` | `revokeAPIKey` | 吊销 API Key | `200 + {status: revoked}` | `400`、`401`、`403`、`404` |

### 3.2 Request Rules

#### `POST /api/v1/auth/oidc/begin`

- 请求体：`BeginOIDCLoginRequest`
- 必填字段：
  - `tenant_name`
  - `redirect_uri`
- 成功响应：
  - `200 + BeginOIDCLoginResponse`

#### `POST /api/v1/auth/token`

- 请求体：`CompleteOIDCLoginRequest`
- 必填字段：
  - `state`
  - `code`
  - `redirect_uri`
- 成功响应：
  - `200 + TokenPairResponse`

#### `POST /api/v1/auth/refresh`

- 请求体：`RefreshAccessTokenRequest`
- 必填字段：
  - `refresh_token`
- 成功响应：
  - `200 + RefreshAccessTokenResponse`

#### `POST /api/v1/auth/logout`

- 请求体：`LogoutRequest`
- 必填字段：
  - `jti`
- 成功响应：
  - `200 + RevokeStatusResponse`

#### `GET /api/v1/auth/api-keys`

- 可选查询参数：
  - `user_id`
- 成功响应：
  - `200 + ListAPIKeysResponse`

#### `POST /api/v1/auth/api-keys`

- 请求体：`CreateAPIKeyRequest`
- 成功响应：
  - `201 + CreateAPIKeyResponse`

#### `DELETE /api/v1/auth/api-keys/{key_id}`

- 路径参数：
  - `key_id`
- 成功响应：
  - `200 + {status: revoked}`

## 4. Data Model

### 4.1 Session Flow Schemas

| schema | 关键字段 | 用途 |
|---|---|---|
| `BeginOIDCLoginRequest` | `tenant_name`、`redirect_uri` | 发起登录 |
| `BeginOIDCLoginResponse` | `authorization_url`、`state` | 返回授权跳转地址 |
| `CompleteOIDCLoginRequest` | `state`、`code`、`redirect_uri` | 换取 token |
| `TokenPairResponse` | `access_token`、`refresh_token`、`expires_in`、`issued_at?` | 登录完成结果 |
| `RefreshAccessTokenRequest` | `refresh_token` | 刷新 access token |
| `RefreshAccessTokenResponse` | `access_token`、`expires_in` | 刷新结果 |
| `LogoutRequest` | `jti` | 吊销当前 token |
| `RevokeStatusResponse` | `status=revoked` | 吊销结果 |

### 4.2 Programmatic Access Schemas

| schema | 关键字段 | 用途 |
|---|---|---|
| `CreateAPIKeyRequest` | `name`、`scopes`、`user_id?`、`rate_limit_rpm?`、`expires_at?` | 创建 API Key |
| `CreateAPIKeyResponse` | `key_id`、`key_value`、`key_prefix` | 创建成功返回一次性明文 |
| `APIKeyInfo` | `id`、`name`、`key_prefix`、`scopes`、`rate_limit_rpm`、`is_active` | API Key 列表项 |
| `ListAPIKeysResponse` | `items`、`total` | API Key 列表结果 |

### 4.3 Non-Frozen Data Areas

- 用户管理
- 角色与权限**编辑**（只读列表见 `tenant-management.md`）
- LDAP 配置中心
- API Key 审计、恢复、轮换、风险分析、合规导出
- 密钥管理、国密加解密、网络安全策略

### 4.4 Split to Separate Module

以下能力已拆至 `docs/console-modules/tenant/tenant-management.md`（Services `/api/v1/svc/tenant/*`）：

- 租户成员邀请/移除
- 角色只读列表
- 租户级 SSO 读写
- Webhook 管理

## 5. Page Design Constraints

### 5.1 Entry Page

- 总入口默认展示“会话链路摘要 + API Key 摘要 + 租户与接入入口 + 待补边界”
- 总入口负责导航与解释，不负责配置中心和审计中心
- 总入口不复制 `API Key 管理` 子页完整 CRUD 细节

### 5.2 Session vs Credential Segmentation

- `OIDC / token / refresh / logout` 必须归入会话链路说明
- `API Key` 必须归入程序化接入凭证说明
- 页面文案必须让用户知道当前问题属于哪一类链路

### 5.3 Child Page Relationship

- `API Key 管理` 使用独立主维护文档维护详细规则
- 总入口只保留其子页入口、能力摘要和安全边界
- 未冻结能力只能作为说明态，不形成可执行入口

## 6. Rules

- 安全与身份概览归属 `Core / Auth`
- 正式路径必须使用 `/api/v1/auth/*`
- 页面不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- 当前权威源未给 `POST /auth/refresh` 声明 `operationId`，正文不得自造
- 当前 `OIDC` 只代表 **Core 登录链路**；租户级 SSO **读写**见 `tenant-management.md`，不等于 LDAP 配置中心已经冻结
- 当前 `API Key` 只代表列表、创建、吊销，不代表审计、恢复或导出已经冻结

## 7. Error Handling

### 7.1 Unified Error Shape

```json
{"code":"UPPER_SNAKE","message":"...","request_id":"..."}
```

### 7.2 Expected Errors

- `POST /api/v1/auth/oidc/begin`: `400 + ErrorResponse`
- `POST /api/v1/auth/token`: `400 + ErrorResponse`、`401 + ErrorResponse`
- `POST /api/v1/auth/refresh`: `400 + ErrorResponse`、`401 + ErrorResponse`
- `POST /api/v1/auth/logout`: `400 + ErrorResponse`、`401 + ErrorResponse`、`403 + ErrorResponse`
- `GET /api/v1/auth/api-keys`: `401 + ErrorResponse`、`403 + ErrorResponse`
- `POST /api/v1/auth/api-keys`: `400 + ErrorResponse`、`401 + ErrorResponse`、`403 + ErrorResponse`
- `DELETE /api/v1/auth/api-keys/{key_id}`: `400 + ErrorResponse`、`401 + ErrorResponse`、`403 + ErrorResponse`、`404 + ErrorResponse`

## 8. Acceptance

- 文档中所有正式路径均位于 `/api/v1/auth/*`
- 文档中不再把 LDAP 配置中心、审计、合规写成已冻结接口
- 文档中不再把租户成员/SSO/Webhook 写成「全部待补」；上述 Services 能力见 `tenant-management.md`
- 文档中明确区分会话链路与程序化接入凭证链路
- 总入口与 `API Key 管理` 子页命名、入口、边界保持一致

## 9. Backfill Dependencies

- 如未来扩展 `LDAP / SSO` 配置中心，需先冻结正式路径与 schema
- 如未来扩展 `API Key` 审计、恢复、轮换、风控或导出，需先写回权威源
- 如未来扩展密钥管理与网络安全策略，需先明确资源归属与总入口边界
