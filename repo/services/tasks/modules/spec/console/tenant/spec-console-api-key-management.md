# SPEC: Console API Key 管理

> Technical specification derived from: `tasks/modules/prd/console/tenant/prd-console-api-key-management.md`
> Source of truth: `ANI-main/repo/api/openapi/v1.yaml`
> Scope: `Console / 安全与身份 / API Key`

## 1. Summary

本 SPEC 定义 `Console / 安全与身份 / API Key` 的技术边界、字段映射、页面承接范围和接口冻结规则。API Key 管理属于 `Core / Auth`，正式前缀为 `/api/v1`，不扩写为 `Services` 资源，也不扩写为平台审计或合规模块。

## 2. Source of Truth

### 2.1 Primary Authority

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Frozen Path Group

- `/auth/api-keys`
- `/auth/api-keys/{key_id}`

### 2.3 Frozen Schemas

- `CreateAPIKeyRequest`
- `CreateAPIKeyResponse`
- `APIKeyInfo`
- `ListAPIKeysResponse`
- `ErrorResponse`

## 3. API Design

### 3.1 Endpoint Matrix

| Method | Path | operationId | summary | Success | Error Responses |
|---|---|---|---|---|---|
| GET | `/api/v1/auth/api-keys` | `listAPIKeys` | 列出当前租户 API Key | `200 + ListAPIKeysResponse` | `401`、`403` |
| POST | `/api/v1/auth/api-keys` | `createAPIKey` | 创建 API Key | `201 + CreateAPIKeyResponse` | `400`、`401`、`403` |
| DELETE | `/api/v1/auth/api-keys/{key_id}` | `revokeAPIKey` | 吊销 API Key | `200 + {status: revoked}` | `400`、`401`、`403`、`404` |

### 3.2 Request Rules

#### `GET /api/v1/auth/api-keys`

- 可选查询参数：
  - `user_id: string`
- 成功响应：
  - `200 + ListAPIKeysResponse`

#### `POST /api/v1/auth/api-keys`

- 请求体：`CreateAPIKeyRequest`
- 必填字段：
  - `name`
  - `scopes`
- 可选字段：
  - `user_id`
  - `rate_limit_rpm`
  - `expires_at`
- 字段约束：
  - `name`: `string`, `maxLength: 128`
  - `scopes`: `array`, `minItems: 1`
  - `scopes[*]`: 匹配 `^scope:[a-z0-9_-]+:(\*|[a-z0-9_-]+)$`
  - `rate_limit_rpm`: `integer`, `minimum: 1`, `maximum: 10000`, `default: 60`
  - `expires_at`: `date-time`
- 成功响应：
  - `201 + CreateAPIKeyResponse`

#### `DELETE /api/v1/auth/api-keys/{key_id}`

- 路径参数：
  - `key_id: string`
- 成功响应：
  - `200 + {status: revoked}`

## 4. Data Model

### 4.1 Schema Mapping

| 页面场景 | 权威源 schema | 说明 |
|---|---|---|
| 列表行 | `APIKeyInfo` | 不回显完整明文，仅展示前缀与状态 |
| 列表响应 | `ListAPIKeysResponse` | 返回 `items` 与 `total` |
| 创建请求 | `CreateAPIKeyRequest` | 严格复用正式字段和约束 |
| 创建结果 | `CreateAPIKeyResponse` | 只在创建成功时返回一次明文 |
| 吊销结果 | inline object | `status` 仅允许 `revoked` |

### 4.2 APIKeyInfo

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `string` | API Key 唯一标识 |
| `name` | `string` | API Key 名称 |
| `key_prefix` | `string` | API Key 前缀，用于列表识别 |
| `scopes` | `string[]` | 权限范围数组 |
| `rate_limit_rpm` | `integer` | 每分钟频率限制 |
| `created_at` | `date-time \| null` | 创建时间 |
| `expires_at` | `date-time \| null` | 过期时间 |
| `last_used_at` | `date-time \| null` | 最近使用时间 |
| `is_active` | `boolean` | 是否仍可用 |

### 4.3 CreateAPIKeyResponse

| 字段 | 类型 | 说明 |
|---|---|---|
| `key_id` | `string` | 新创建的 key ID |
| `key_value` | `string` | 完整明文，仅创建时返回一次 |
| `key_prefix` | `string` | 后续列表可见前缀 |

## 5. Page Design Constraints

### 5.1 List Page

- 承接列表查看、创建入口和吊销入口
- 列表字段必须可回指到 `APIKeyInfo`
- 不在列表页承诺“查看明文”“恢复 key”“审计明细”等未冻结能力

### 5.2 Create Flow

- 创建表单仅承接 `CreateAPIKeyRequest` 正式字段
- 创建成功后必须进入一次性明文展示态
- 用户关闭成功弹窗后，页面只能继续展示 `key_prefix`

### 5.3 Revoke Flow

- 吊销属于危险动作，必须先确认再提交
- 吊销成功后刷新列表，不把“请求已发出”误写成其他状态语义
- 吊销结果仅承接 `{status: revoked}`，不自造更多返回字段

## 6. Rules

- API Key 管理归属 `Core / Auth`
- 正式路径必须使用 `/api/v1/auth/api-keys*`
- 页面不要求前端显式传 `tenant_id`
- 页面不缓存 `key_value`，不在后续任何视图中回显完整明文
- 审计、风控、合规导出当前只作为模块边界，不能写成正式接口

## 7. Error Handling

### 7.1 Unified Error Shape

```json
{"code":"UPPER_SNAKE","message":"...","request_id":"..."}
```

### 7.2 Expected Errors

- `GET /api/v1/auth/api-keys`: `401 + ErrorResponse`、`403 + ErrorResponse`
- `POST /api/v1/auth/api-keys`: `400 + ErrorResponse`、`401 + ErrorResponse`、`403 + ErrorResponse`
- `DELETE /api/v1/auth/api-keys/{key_id}`: `400 + ErrorResponse`、`401 + ErrorResponse`、`403 + ErrorResponse`、`404 + ErrorResponse`

### 7.3 Documentation Rule

- 当前权威源未冻结的能力，不允许写成正式路径、正式 schema 或正式返回码
- 不允许把“明文再次查看”“恢复 API Key”“批量导出”写成当前事实

## 8. Acceptance

- 文档中所有正式路径均位于 `/api/v1/auth/api-keys*`
- 文档中所有 schema、字段与 `v1.yaml` 保持一致
- 文档中明确 `key_value` 仅创建时返回一次
- 文档中不再出现 `[Assumption]`、旧 `console` 路径和未冻结能力伪装

## 9. Backfill Dependencies

- 如未来新增 API Key 轮换、恢复或编辑能力，需先冻结正式路径与 schema
- 如未来新增 API Key 审计、风险分析或合规导出，需先明确模块归属
- 如未来新增明文再次查看能力，需先完成安全审查并写回权威源
