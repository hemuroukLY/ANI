# 安全与身份概览

## 页面定位

`安全与身份概览` 是 `Console / 安全与身份` 的租户侧总入口，用于收口当前租户可见的认证、会话与程序化接入凭证能力，并明确已冻结能力与待补边界。

当前模块属于 `Core / Auth`，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`，正式路径前缀为 `/api/v1`。

本页属于 `Console` 页面，不是 `BOSS` 的平台审计、合规导出或风险运营页。

## 文档管理规则

- 本文是 `安全与身份概览` 的主维护文档，页面定位、分层边界、入口关系和验收标准统一以本文为准
- `prototypes/ani-services-prototype-console.html` 只保留摘要、边界和入口
- `tasks/modules/prd/console/tenant/prd-console-security-identity-overview.md` 与 `tasks/modules/spec/console/tenant/spec-console-security-identity-overview.md` 作为辅助材料保留，不替代本文
- 如本文与辅助材料出现差异，先对照 `ANI-main/repo/api/openapi/v1.yaml`，再统一回写

## Core 层要求

- 当前可确认的安全域路径只来自 `Core / Auth`
- 已冻结路径必须使用 `/api/v1/auth/*`
- 不允许继续使用旧 `/api/v1/console/*` 路径
- 不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- `OIDC` 当前只承接 **Core 登录链路**；租户级 SSO **读写**见 `tenant/tenant-management.md`（Services `/tenant/sso`），不等于本页已冻结
- `API Key` 当前只承接列表、创建、吊销，不等于轮换、恢复、审计、导出已经冻结
- 当前权威源未给 `POST /auth/refresh` 声明 `operationId`，正文不得自造

## 页面职责

- 汇总当前已冻结的 `OIDC begin / token / refresh / logout / API Key`
- 把会话链路能力和程序化接入凭证能力分开展示
- 给 `API Key 管理` 子页提供总入口和边界说明
- 给 **租户管理** 子页提供总入口（成员、角色只读、SSO、Webhook 见 `tenant/tenant-management.md`）
- 对 `用户管理 / 角色与权限编辑 / LDAP / 密钥管理 / 国密加解密 / 网络安全策略` 维持待补说明
- 明确 `Console` 与 `BOSS` 在安全域中的分工

## 页面结构

- 首屏至少包含 `会话链路概览`、`程序化接入凭证概览`、`租户与接入入口`、`待补能力边界` 四个区块
- `会话链路概览` 至少覆盖 `OIDC begin / token / refresh / logout`
- `程序化接入凭证概览` 至少覆盖 `API Key 管理` 入口和摘要
- 总入口只负责摘要、导航和边界说明，不重写 `API Key 管理` 子页全量 CRUD 细节

## 数据来源与分层约束

### 数据来源划分

- `Core` 数据
  - OIDC 登录链路
  - Token 换取与刷新
  - Logout 吊销
  - API Key 列表、创建、吊销
- `Services` 数据
  - 本页不直接承接 Services 契约表；**租户管理**由独立模块维护（`/api/v1/svc/tenant/*`）

### 关键边界

- 安全与身份概览不等于“安全配置中心”
- 总入口只做已冻结能力摘要，不替代各子模块的完整定义
- `API Key 管理` 由独立主维护文档维护详细字段、接口冻结规则和示例
- `LDAP` 配置、用户角色编辑、审计导出、风控分析只保留为待补边界
- 租户级 SSO **读写**不在本页展开，见 `tenant/tenant-management.md`；**角色只读列表**亦在该模块收口

## 冻结事实表

### Frozen Paths

- `POST /api/v1/auth/oidc/begin`
- `POST /api/v1/auth/token`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/api-keys`
- `POST /api/v1/auth/api-keys`
- `DELETE /api/v1/auth/api-keys/{key_id}`

### Frozen Schemas

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

### Non-Frozen Capabilities

- 用户管理
- 角色与权限**编辑**（只读列表见 `tenant-management.md`）
- LDAP 配置中心
- API Key 轮换、恢复、审计、风险分析、合规导出
- 密钥管理、国密加解密、网络安全策略

以下能力已拆至独立模块：`docs/console-modules/tenant/tenant-management.md`（成员、角色只读、SSO、Webhook）

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 会话链路概览 | Core | 对齐 `Auth` 登录与 token 链路 | 登录流程说明 |
| 程序化接入凭证概览 | Core | 对齐 `API Key` 正式能力摘要 | `API Key 管理` |
| 租户与接入入口 | Services（摘要） | 成员/角色/SSO/Webhook 详文见租户管理模块 | `租户管理` |
| 待补能力说明 | 规划项 | 仅保留边界，不形成正式接口表 | 后续安全模块 |

## 字段级定义

### 会话链路字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `tenant_name` | `BeginOIDCLoginRequest.tenant_name` | 租户 slug，用于限定登录上下文 |
| `redirect_uri` | `BeginOIDCLoginRequest.redirect_uri` / `CompleteOIDCLoginRequest.redirect_uri` | OIDC 回调地址 |
| `authorization_url` | `BeginOIDCLoginResponse.authorization_url` | 授权跳转地址 |
| `state` | `BeginOIDCLoginResponse.state` / `CompleteOIDCLoginRequest.state` | 登录状态串 |
| `code` | `CompleteOIDCLoginRequest.code` | OIDC 返回授权码 |
| `access_token` | `TokenPairResponse.access_token` / `RefreshAccessTokenResponse.access_token` | 访问令牌 |
| `refresh_token` | `TokenPairResponse.refresh_token` / `RefreshAccessTokenRequest.refresh_token` | 刷新令牌 |
| `expires_in` | `TokenPairResponse.expires_in` / `RefreshAccessTokenResponse.expires_in` | 过期秒数 |
| `jti` | `LogoutRequest.jti` | 当前 token 的 JWT ID |

### API Key 摘要字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `key_prefix` | `APIKeyInfo.key_prefix` | API Key 前缀，不回显明文 |
| `is_active` | `APIKeyInfo.is_active` | 当前 key 是否仍可用 |
| `key_value` | `CreateAPIKeyResponse.key_value` | 仅创建成功时返回一次 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 会话链路正常态 | 展示链路步骤和能力摘要 | 不展示完整 token 明文 |
| API Key 正常态 | 展示子页入口和摘要口径 | 详细字段下沉到 `API Key 管理` |
| 首次无 API Key | 突出“创建首个 API Key”入口 | 不在总入口直接重写创建表单 |
| 会话异常 | 引导重新登录或刷新 token | 不只给静态报错 |
| 待补能力 | 说明态展示，不生成伪入口 | 避免误导为已冻结能力 |

## 字段口径与单位

| 字段 | 口径建议 | 单位/格式 |
|---|---|---|
| `tenant_name` | 标识登录上下文租户，不等于租户配置详情 | 字符串 |
| `redirect_uri` | 必须与注册回调地址一致 | URI |
| `authorization_url` | 一次登录过程中的跳转地址 | URI |
| `expires_in` | token 有效期时长 | 秒，整数 |
| `key_prefix` | 只作为凭据识别前缀，不可直接使用 | 字符串 |

## 状态与能力口径

### 会话与凭证边界

- 会话链路：`OIDC begin -> token -> refresh -> logout`
- 程序化接入：`API Key` 列表、创建、吊销
- 以上两类能力必须分开展示，不能混写成模糊的“安全设置”

### 待补能力边界

- 用户管理、LDAP
- API Key 审计、恢复、轮换、风险分析、合规导出
- 密钥管理、国密加解密、网络安全策略

以下能力已拆至独立模块：`docs/console-modules/tenant/tenant-management.md`（成员、角色、SSO、Webhook）

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| OIDC 登录 | Provider 配置有效 | `400 BAD_REQUEST` |
| Token 刷新 | `refresh_token` 未过期 | `401 UNAUTHORIZED` |
| API Key 创建 | 租户管理员权限 | `403 FORBIDDEN` |

租户成员/角色/SSO/Webhook 管理见 `tenant/tenant-management.md`（Services `/tenant/*`）。

## 操作可用性矩阵（页面入口）

| 操作 | 已登录成员 | 租户管理员 | 说明 |
|---|---|---|---|
| 查看安全与身份概览 | 可用 | 可用 | 总入口概览 |
| 理解 OIDC / Token 链路 | 可用 | 可用 | 会话能力说明 |
| 进入 API Key 管理 | 视权限而定 | 可用 | 进入子页而非在本页直接操作 |
| 查看待补边界说明 | 可用 | 可用 | 说明态内容 |

## 接口冻结规则

### `POST /api/v1/auth/oidc/begin`

- `operationId`: `beginOIDCLogin`
- `requestBody`: `BeginOIDCLoginRequest`
- `requestBody.required`: `tenant_name`、`redirect_uri`
- `success`: `200 + BeginOIDCLoginResponse`
- `errors`: `400`

### `POST /api/v1/auth/token`

- `operationId`: `completeOIDCLogin`
- `requestBody`: `CompleteOIDCLoginRequest`
- `requestBody.required`: `state`、`code`、`redirect_uri`
- `success`: `200 + TokenPairResponse`
- `errors`: `400`、`401`

### `POST /api/v1/auth/refresh`

- `operationId`: 当前权威源未声明
- `requestBody`: `RefreshAccessTokenRequest`
- `requestBody.required`: `refresh_token`
- `success`: `200 + RefreshAccessTokenResponse`
- `errors`: `400`、`401`

### `POST /api/v1/auth/logout`

- `operationId`: `logout`
- `requestBody`: `LogoutRequest`
- `requestBody.required`: `jti`
- `success`: `200 + RevokeStatusResponse`
- `errors`: `400`、`401`、`403`

### `GET /api/v1/auth/api-keys`

- `operationId`: `listAPIKeys`
- `success`: `200 + ListAPIKeysResponse`
- `errors`: `401`、`403`

### `POST /api/v1/auth/api-keys`

- `operationId`: `createAPIKey`
- `success`: `201 + CreateAPIKeyResponse`
- `errors`: `400`、`401`、`403`

### `DELETE /api/v1/auth/api-keys/{key_id}`

- `operationId`: `revokeAPIKey`
- `success`: `200 + {status: revoked}`
- `errors`: `400`、`401`、`403`、`404`

## 响应示例

### OIDC 登录启动成功

```json
{
  "authorization_url": "https://dex.example.com/auth?client_id=ani-console&state=oidc-state-001",
  "state": "oidc-state-001"
}
```

### TokenPair 返回成功

```json
{
  "access_token": "eyJhbGciOi...",
  "refresh_token": "refresh-token-001",
  "expires_in": 3600,
  "issued_at": "2026-06-15T11:00:00Z"
}
```

### Logout 成功

```json
{
  "status": "revoked"
}
```

## 错误示例

### OIDC 登录参数错误

```json
{
  "code": "BAD_REQUEST",
  "message": "redirect_uri is invalid",
  "request_id": "req-auth-400-001"
}
```

### 会话无效

```json
{
  "code": "UNAUTHORIZED",
  "message": "token is invalid or expired",
  "request_id": "req-auth-401-001"
}
```

## 回填前置依赖

- 后续若要补 `用户 / 角色 / LDAP / SSO 配置`，必须先形成正式 schema 与路径
- 后续若要补 `API Key 审计 / 导出 / 风险分析 / 恢复 / 轮换`，必须先扩充权威源
- 后续若要补 `密钥管理 / 国密加解密 / 网络安全策略`，必须先明确资源归属

## 回填验收标准

- 正文明确区分已冻结能力与待补能力
- 正文不再出现旧 `/api/v1/console/*` 路径
- 正文不要求前端显式传 `tenant_id / X-Tenant-Id`
- 正文不把 `LDAP / SSO 配置 / 审计导出 / 风险分析` 写成已冻结接口
- `PRD`、`SPEC`、HTML 摘要与本文一致

## 产品经理补充定义

### 目标用户与权限视角

- 已登录租户管理员：查看安全域概览、进入 API Key 子页、理解当前认证边界
- 普通成员：查看与自身权限一致的认证说明和 API Key 入口
- 只读用户：查看认证链路说明和边界提示，不显示无权写操作

### 页面结构补充

- 首屏至少包含 `会话链路概览`、`程序化接入凭证概览`、`租户与接入入口`、`待补能力边界` 四个区块
- 总入口只承接摘要和导航，不在本页重写 API Key 全量 CRUD 细节
- 页面需要把会话相关能力和凭证相关能力分开展示，避免信息混杂

### 核心任务流

1. 用户进入总入口后，快速理解当前已冻结的 `OIDC / token / refresh / logout / API Key` 能力
2. 用户根据需求进入 `API Key 管理` 子页处理程序化接入凭证
3. 用户从会话过期、权限不足或待补能力提示中，明确知道当前能做什么、不能做什么

### 页面状态与反馈

- 首次进入且无 API Key 时，页面应突出“创建首个 API Key”入口
- 会话相关异常必须引导到重新登录或刷新流程，而不是只给静态报错
- 待补能力以说明态出现，不出现可点击但无契约支撑的入口
- 权限不足时，页面保留概览，但不暴露无权子功能按钮

### 跨模块协同

- 与 `API Key 管理`、`开放与集成总入口` 协同，共同形成认证与接入闭环
- 与首页协同，用于承接安全域异常或凭证待处理事项
- LDAP、SSO 配置、审计导出、密钥管理等未冻结能力只保留边界说明

### 产品验收补充

- 用户必须能在一个页面内理解“当前已冻结的安全能力”和“仍待补的安全域能力”
- 总入口与 API Key 子页的命名、入口、面包屑必须一致
- 页面必须区分会话态问题与权限态问题，不混用文案
- 本页不得把总入口写成新的身份资源中心或配置中心
