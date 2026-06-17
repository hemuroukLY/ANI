# API Key 管理

## 页面定位

`API Key 管理` 是 `Console / 安全与身份` 下的租户侧凭证管理页面，用于查看、创建和吊销当前权限范围内的 API Key。

当前模块属于 `Core / Auth`，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`，正式路径前缀为 `/api/v1`。

本页是 `Console` 页面，不是 `BOSS / API Key 审计` 页面。

## 文档管理规则

- 本文是 `API Key 管理` 的主维护文档，页面定位、字段口径、交互边界和验收标准统一以本文为准
- 当 `PRD`、`SPEC`、HTML 摘要与本文存在差异时，先对照 `ANI-main/repo/api/openapi/v1.yaml`，再统一回写
- `prototypes/ani-services-prototype-console.html` 只保留页面摘要、模块边界和详细材料入口
- 未冻结能力只能写入待补边界，不得写成正式路径、正式 schema 或正式返回码

## Core 层要求

- API Key 能力属于 `Core / Auth`
- 正式路径使用 `/api/v1/auth/api-keys*`
- 不允许继续使用旧 `/api/v1/console/*`
- 当前页面不要求前端显式传 `tenant_id`
- 当前页面不把 API Key 审计、风险分析、合规导出写成已冻结能力
- 当前页面不把“再次查看明文 key”“恢复 API Key”写成现有正式能力
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

## 页面职责

- 展示当前租户可见的 API Key 列表
- 提供创建 API Key 入口
- 提供吊销 API Key 入口
- 在创建成功后一次性展示完整明文 `key_value`
- 明确接入前准备与安全边界，但不替代审计或合规模块

## 页面结构

- 列表页至少包含 `顶部动作区`、`API Key 列表`、`筛选区`、`空态或错误态反馈`
- 创建流程至少包含 `名称`、`scopes`、`user_id(可选)`、`rate_limit_rpm(可选)`、`expires_at(可选)`
- 成功创建后必须进入 `一次性明文展示态`
- 吊销动作必须在当前上下文内完成确认，不依赖未冻结的新子页

## 数据来源与分层约束

### 数据来源划分

- `Core` 数据
  - API Key 列表
  - API Key 创建
  - API Key 吊销
- `Services` 数据
  - 当前页面不承接 `Services` 资源契约

### 关键边界

- 本页直接承接 `Core / Auth` 能力，不通过 `Services` 再包装
- 当前正式资源只有列表、创建、吊销三条链路
- `key_value` 仅出现在 `CreateAPIKeyResponse` 中，不能写成列表字段或详情常驻字段
- 审计、风控、合规导出、明文再次查看只保留为待补边界

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 结果去向 |
|---|---|---|---|
| API Key 列表 | Core | 对齐 `ListAPIKeysResponse.items[] -> APIKeyInfo` | 创建入口 / 吊销确认 |
| 创建 API Key | Core | 对齐 `CreateAPIKeyRequest` | 一次性明文展示 |
| 创建成功反馈 | Core | 对齐 `CreateAPIKeyResponse` | 返回列表复核 |
| 吊销 API Key | Core | 对齐 `DELETE /api/v1/auth/api-keys/{key_id}` | 刷新列表 |
| 边界说明 | 规划项 | 标明审计、风控、合规不在当前模块内 | 后续安全模块 |

## 字段级定义

### 列表与详情字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `id` | `APIKeyInfo.id` | API Key 唯一标识 |
| `name` | `APIKeyInfo.name` | API Key 名称 |
| `key_prefix` | `APIKeyInfo.key_prefix` | 前缀，用于识别，不展示完整明文 |
| `scopes` | `APIKeyInfo.scopes` | 权限范围数组 |
| `rate_limit_rpm` | `APIKeyInfo.rate_limit_rpm` | 每分钟频率限制 |
| `created_at` | `APIKeyInfo.created_at` | 创建时间 |
| `expires_at` | `APIKeyInfo.expires_at` | 过期时间，可空 |
| `last_used_at` | `APIKeyInfo.last_used_at` | 最近使用时间，可空 |
| `is_active` | `APIKeyInfo.is_active` | 是否仍可用 |

### 创建请求字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `name` | `CreateAPIKeyRequest.name` | 名称，最长 `128` 字符 |
| `user_id` | `CreateAPIKeyRequest.user_id` | 可选；为空时使用当前认证用户 |
| `scopes` | `CreateAPIKeyRequest.scopes` | 至少 1 项，需满足正式 pattern |
| `rate_limit_rpm` | `CreateAPIKeyRequest.rate_limit_rpm` | `1 ~ 10000`，默认 `60` |
| `expires_at` | `CreateAPIKeyRequest.expires_at` | 过期时间 |

### 创建成功字段

| 字段 | 来源 schema | 说明 |
|---|---|---|
| `key_id` | `CreateAPIKeyResponse.key_id` | 新创建 key 的 ID |
| `key_value` | `CreateAPIKeyResponse.key_value` | 完整明文，仅创建时返回一次 |
| `key_prefix` | `CreateAPIKeyResponse.key_prefix` | 后续列表可见前缀 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 列表正常态 | 展示名称、前缀、scopes、状态、时间字段 | 字段来源必须可回指到 `APIKeyInfo` |
| 空态 | 展示“先创建首个 API Key”引导 | 不伪造示例 key |
| 创建成功态 | 显示一次性明文 `key_value` 和保存提示 | 关闭后不再回显 |
| 权限不足 | 保留列表壳或只读视图，隐藏无权动作 | 不伪造可点击入口 |
| 吊销后反馈 | 刷新列表并明确当前 key 已不可用 | 不需要额外自造状态机 |

## 字段口径与单位

| 字段 | 口径建议 | 单位/格式 |
|---|---|---|
| `name` | 用户可读名称，用于区分不同接入用途 | 字符串 |
| `key_prefix` | 只作为识别前缀，不等于可使用凭据 | 字符串 |
| `rate_limit_rpm` | 以单个 key 的请求频率限制为准 | RPM，整数 |
| `created_at` | 服务端创建时间 | date-time |
| `expires_at` | 超过该时间后凭据应视为失效 | date-time，可空 |
| `last_used_at` | 最近一次成功使用时间 | date-time，可空 |

## 状态与能力口径

### 凭据状态口径

- `is_active = true`：当前 key 仍可用
- `is_active = false`：当前 key 已不可用或已被吊销

### 待补能力边界

- API Key 更新
- API Key 恢复或轮换
- API Key 明文再次查看
- API Key 审计明细
- API Key 风险分析
- API Key 合规导出

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录与租户上下文 | 已认证 | `401 UNAUTHORIZED` / `403 FORBIDDEN` |
| API Key 写权限 | 创建/吊销需凭证管理员 | `403 FORBIDDEN` |
| `name`、`scopes` | 创建时必填；`scopes` 符合 `^scope:[a-z0-9_-]+:(\*|[a-z0-9_-]+)$` | `400 BAD_REQUEST` |

## 操作可用性矩阵

| 操作 | 只读用户 | 凭证管理员 | 说明 |
|---|---|---|---|
| 查看 API Key 列表 | 可用 | 可用 | 读取 `GET /api/v1/auth/api-keys` |
| 按 `user_id` 过滤 | 可用 | 可用 | 可选查询参数 |
| 创建 API Key | 不可用 | 可用 | 返回 `201 + CreateAPIKeyResponse` |
| 查看一次性明文 | 不可用 | 可用 | 仅限创建成功后当前反馈态 |
| 吊销 API Key | 不可用 | 可用 | 返回 `200 + {status: revoked}` |

## 接口冻结规则

### `GET /api/v1/auth/api-keys`

- `operationId`: `listAPIKeys`
- `summary`: `列出当前租户 API Key`
- `tags`: `["Auth"]`
- `query`: `user_id?`
- `success`: `200 + ListAPIKeysResponse`
- `errors`: `401`、`403`

### `POST /api/v1/auth/api-keys`

- `operationId`: `createAPIKey`
- `summary`: `创建 API Key`
- `tags`: `["Auth"]`
- `requestBody`: `CreateAPIKeyRequest`
- `requestBody.required`: `name`、`scopes`
- `requestBody.optional`: `user_id`、`rate_limit_rpm`、`expires_at`
- `scope pattern`: `^scope:[a-z0-9_-]+:(\*|[a-z0-9_-]+)$`
- `success`: `201 + CreateAPIKeyResponse`
- `errors`: `400`、`401`、`403`

### `DELETE /api/v1/auth/api-keys/{key_id}`

- `operationId`: `revokeAPIKey`
- `summary`: `吊销 API Key`
- `tags`: `["Auth"]`
- `path params`: `key_id`
- `success`: `200 + {status: revoked}`
- `errors`: `400`、`401`、`403`、`404`

## 响应示例

### API Key 列表成功

```json
{
  "items": [
    {
      "id": "key_001",
      "name": "cli-deploy",
      "key_prefix": "ani_live_x7Q",
      "scopes": ["scope:inference:*", "scope:models:read"],
      "rate_limit_rpm": 120,
      "created_at": "2026-06-15T09:00:00Z",
      "expires_at": "2026-12-31T23:59:59Z",
      "last_used_at": "2026-06-15T10:20:00Z",
      "is_active": true
    }
  ],
  "total": 1
}
```

### 创建 API Key 成功

```json
{
  "key_id": "key_002",
  "key_value": "ani_live_4mWJ8TnSg2uQpX9a",
  "key_prefix": "ani_live_4mW"
}
```

### 吊销 API Key 成功

```json
{
  "status": "revoked"
}
```

## 错误示例

### 创建请求参数错误

```json
{
  "code": "BAD_REQUEST",
  "message": "scopes must contain at least one valid scope",
  "request_id": "req-apikey-400-001"
}
```

### 吊销不存在的 API Key

```json
{
  "code": "NOT_FOUND",
  "message": "api key not found",
  "request_id": "req-apikey-404-001"
}
```

## 验收标准

- 正文路径、schema、返回码与现有 `ANI-main/repo/api/openapi/v1.yaml` 一致
- 正文明确 `key_value` 仅在创建成功时返回一次
- 正文不把 API Key 审计、风险分析、合规导出写成当前模块正式接口
- HTML 摘要、PRD、SPEC 与本文保持一致
- 本文可独立作为 `Console / API Key 管理` 的主维护材料

## 产品经理补充定义

### 目标用户与权限视角

- 有写权限的租户管理员或凭证管理员：查看、创建、吊销 API Key
- 普通成员：仅查看自己被授权可见的 API Key；是否可创建由后端权限控制
- 只读用户：仅查看列表与状态，不显示创建和吊销入口

### 默认视图与页面状态

- 默认展示 API Key 列表和创建入口，无 key 时展示“先创建首个 API Key”的空态
- 创建表单推荐按 `名称 -> scopes -> user_id(可选) -> 频率限制 -> 过期时间` 的顺序组织
- 创建成功后必须进入一次性明文展示态，并要求用户完成保存确认后再关闭
- 权限不足或 `403` 时，页面保留列表壳或只读壳，但不暴露无权动作

### 核心任务流

1. 用户进入列表查看现有 API Key 的名称、前缀、状态、过期时间和最近使用时间
2. 用户创建 API Key，并在成功弹窗中复制一次性明文 `key_value`
3. 用户返回列表确认仅展示 `key_prefix`，后续不再回显完整明文
4. 用户在确认风险后吊销 API Key，并通过结果反馈确认当前 key 已不可用

### 跨模块协同

- 与 `安全与身份概览`、`开放与集成总入口` 协同，作为接入前准备环节的子页
- 与 API 文档、SDK、CLI 接入说明协同，但不在本页重写接入教程
- 审计、风险分析、合规导出等未冻结能力只保留边界说明

### 产品验收补充

- 用户必须在一个闭环内完成“创建 -> 保存明文 -> 返回列表复核”
- 页面必须用显著文案强调 `key_value` 只返回一次
- 吊销动作必须有二次确认，并明确影响范围是当前 key 不再可用
- 本页不得出现再次查看明文、恢复 key 或伪造审计能力
