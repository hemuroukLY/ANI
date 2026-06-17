# API Key 审计

## 页面定位

`API Key 审计` 是 BOSS **安全审计与合规** 域下的 **跨租户 API Key 生命周期与调用审计** 专页：汇总全平台 Key 的创建、使用、吊销与异常调用模式，供平台管理员与安全团队排查凭证滥用与限流异常。

本页属于 **Core / Auth + Audit** 视角下的 **平台 RBAC** 页面。租户侧 Auth CRUD 已在 `v1.yaml` 冻结（`/api/v1/auth/api-keys*`）；**调用事件 list** **ADDED-TO-YAML**；Key 元数据只读参考 auth/api-keys。

Console [`api-key-audit.md`](../../console-modules/tenant/api-key-audit.md) 为 **单租户 Key 调用事件** 检索（草案路径 `GET .../audit-events` — **待 YAML，非冻结**）；本页须 **平台 RBAC + 跨租户聚合**，**不得** 逐租户 JWT 轮询 `listAPIKeys` 冒充平台审计 dashboard。

## 文档管理规则

- 本文是 `API Key 审计` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-platform-apikey-audit.md`](../../tasks/modules/prd/boss/audit/prd-boss-platform-apikey-audit.md) 与 [`spec-boss-platform-apikey-audit.md`](../../tasks/modules/spec/boss/audit/spec-boss-platform-apikey-audit.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)
- Console 对照：[`api-key-audit.md`](../../console-modules/tenant/api-key-audit.md)（草案路径见 `openapi-phase3-api-key-audit-draft.md`，**非冻结**）

## Core 层要求

- Auth API Key CRUD 归属 **Core** `/api/v1/auth/api-keys*`
- 平台跨租户 API Key 调用事件 list — **ADDED-TO-YAML** `listPlatformAPIKeyAuditEvents`
- 租户 Auth 路径为 **只读参考**：`listAPIKeys` / `createAPIKey` / `revokeAPIKey` — **非** BOSS 正式跨租户契约
- `APIKeyInfo` **不含** `call_count`、`rate_limit_hits` — UI **不得** 伪造这两列
- 平台 RBAC 鉴权；**不得** 信任未授权 `tenant_id` 越权
- 写操作（平台代吊销等）— **TODO-YAML** Phase 2；合入后须 `idempotency_key`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）
- 禁止自造 audit-events path / schema 为已冻结（Console 草案仍为 **规划**）
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **全平台** API Key 摘要 list 与跨租户筛选（platform list **ADDED-TO-YAML**）
- 展示 Key 元数据：`APIKeyInfo` 已冻结字段 + 审计扩展字段（待 YAML）
- 展示 **调用事件时间线**（platform audit-events list 待 YAML；Console 草案 `.../audit-events` 为单 Key 单租户）
- 支持异常检测视图：长期未用、突增调用、限流命中（限流命中字段 **待 YAML**）
- 支持跳转 **平台审计日志**、**Trace**（携带 `request_id`）
- 与 Console API Key 管理 / 审计对照，明确 BOSS 跨租户差异

## 页面结构

- 首屏至少包含：`筛选区`、`Key 摘要表格`、`调用事件时间线`、`详情抽屉`、`边界说明`
- 无数据态、无权限态、list API 未就绪态须可区分

```text
API Key 审计
├── 筛选（tenant_id / key_prefix / is_active / 时间 — list 待 YAML）
├── Key 摘要表格（prefix / 租户 / scopes / last_used_at / is_active）
├── 调用事件时间线（待 YAML · 非 APIKeyInfo 字段）
├── 详情抽屉（APIKeyInfo + 事件详情）
│   └── 跳转 platform-audit-log / platform-trace
└── 操作：平台代吊销（Phase 2 · 待 YAML）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/auth/api-keys` | 租户 **只读参考**；`listAPIKeys` → `APIKeyInfo` |
| Core | `POST /api/v1/auth/api-keys` | 租户创建；BOSS 不直接用于 list |
| Core | `DELETE /api/v1/auth/api-keys/{key_id}` | 租户吊销；BOSS 平台代吊销 **待 YAML** |
| Core | platform apikey audit list **TODO-YAML** | BOSS 正式跨租户 Key + 事件聚合 |
| Core | `GET .../audit-events`（Console 草案） | **草案/待 YAML**；单租户单 Key；**非** BOSS 正式源 |

### 关键边界

- `listAPIKeys` 返回当前 JWT 租户上下文 — **禁止** BOSS 逐租户轮询拼 dashboard
- `APIKeyInfo` 仅有 `last_used_at`，**无** `call_count` / `rate_limit_hits`
- Console 草案 `listApiKeyAuditEvents` path — **未合入 v1.yaml**，正文标 **草案/待 YAML**
- Key 明文 `key_value` 仅在 `createAPIKey` 201 响应出现一次 — 审计页 **永不展示**
- 调用事件 schema（`ApiKeyAuditEvent` 草案）与平台审计日志 schema 应对齐（待 Core）

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 筛选区 | UI + platform list **待 YAML** | tenant / prefix / active | 刷新 |
| Key 摘要表 | platform list **待 YAML** | `APIKeyInfo` 子集 + tenant_id | 详情 / 事件线 |
| 调用事件线 | audit-events **待 YAML** | 非 APIKeyInfo 字段 | Trace |
| 详情抽屉 | `APIKeyInfo` + event get | 只读 | api-key-management（Console） |
| 边界说明 | 规划项 | 草案 path 非冻结 | api-key-audit（Console） |

## BOSS 与 Console 分工

| 维度 | BOSS API Key 审计 | Console API Key 审计 |
|---|---|---|
| 范围 | 全平台多租户 | 当前租户 |
| Key list | platform audit list **待 YAML** P1 | `listAPIKeys`（租户 JWT） |
| 调用事件 | platform events list **待 YAML** | 草案 `GET .../audit-events` **待 YAML** |
| call_count / rate_limit_hits | **待 YAML**；当前 schema **无** | 同约束 |
| 代吊销 | Phase 2 平台写 API | 租户 `revokeAPIKey` |
| RBAC | 平台 audit 读（待 YAML） | `scope:auth:api-keys:audit:read`（草案） |

## 当前冻结事实

| 方法 | 路径 | operationId | RBAC | BOSS 用法 |
|---|---|---|---|---|
| GET | `/api/v1/auth/api-keys` | `listAPIKeys` | 租户 JWT | **只读参考** |
| POST | `/api/v1/auth/api-keys` | `createAPIKey` | 租户 JWT | **只读参考**（创建审计事件源） |
| DELETE | `/api/v1/auth/api-keys/{key_id}` | `revokeAPIKey` | 租户 JWT | **只读参考**（吊销事件源） |

| Schema | 字段（冻结） |
|---|---|
| `APIKeyInfo` | `id`, `name`, `key_prefix`, `scopes`, `rate_limit_rpm`, `created_at`, `expires_at`, `last_used_at`, `is_active` |

| 能力 | 状态 |
|---|---|
| 平台跨租户 API Key 审计 list | **ADDED-TO-YAML** P1 |
| 平台跨租户调用事件 list | **ADDED-TO-YAML** `listPlatformAPIKeyAuditEvents` |
| `call_count` / `rate_limit_hits` 字段 | **TODO-YAML**（当前 schema **无**） |
| Console `.../audit-events` | **草案/待 YAML**（Phase 3 草案） |
| 平台代吊销 | **TODO-YAML** Phase 2 |

## 字段级定义

### 查询字段（platform list 目标 · 待 YAML）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `tenant_id` | query | 可选 | 平台 RBAC 筛选 |
| `key_prefix` | query | 可选 | 前缀模糊 |
| `is_active` | query | 可选 | 布尔 |
| `scope` | query | 可选 | scope 包含筛选 |
| `last_used_before` | query | 可选 | 长期未用检测 |
| `start_time` / `end_time` | query | 可选 | 事件时间窗 |
| `limit` / `cursor` | query | 可选 | 分页 |

### 返回字段（APIKeyInfo · Core 已冻结 · 租户 path）

| 字段 | 来源 | 说明 |
|---|---|---|
| `id` | `APIKeyInfo` | Key UUID |
| `name` | `APIKeyInfo` | 显示名 |
| `key_prefix` | `APIKeyInfo` | 前缀（非明文） |
| `scopes` | `APIKeyInfo` | scope 数组 |
| `rate_limit_rpm` | `APIKeyInfo` | 每分钟限流 |
| `created_at` | `APIKeyInfo` | 创建时间 |
| `expires_at` | `APIKeyInfo` | 过期；可空 |
| `last_used_at` | `APIKeyInfo` | 最近调用；可空 |
| `is_active` | `APIKeyInfo` | 是否有效 |

### 返回字段（调用事件 · **待 YAML** · 草案 `ApiKeyAuditEvent`）

| 字段 | 来源 | 说明 |
|---|---|---|
| `event_id` | 草案 | 事件 ID |
| `key_id` | 草案 | 关联 Key |
| `timestamp` | 草案 | 调用时间 |
| `source_ip` | 草案 | 来源 IP |
| `path` | 草案 | 路径摘要 |
| `status_code` | 草案 | HTTP 状态 |
| `request_id` | 草案 | 关联 Trace |
| `call_count` | **待 YAML** | **非** 当前 `APIKeyInfo` 字段 |
| `rate_limit_hits` | **待 YAML** | **非** 当前 `APIKeyInfo` 字段 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `tenant_id_display` | platform list 返回；租户 path **无** |
| `status_label` | `is_active` → 有效/已吊销 |
| `scopes_summary` | scopes 折叠展示 |
| `idle_days` | now − `last_used_at` 天数 |
| `anomaly_flag` | 产品规则（待 YAML 指标） |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | Key 表 + 事件线 | platform list **ADDED-TO-YAML** |
| 无数据态 | 「暂无 API Key 审计数据」 | **不** 伪造 |
| list 未就绪 | 「跨租户 API Key 审计待 Core 合入」 | 不得 JWT 轮询 |
| `is_active=false` | 灰色行 + 「已吊销」 | — |
| `last_used_at` 空 | 「从未使用」 | — |
| `call_count` / `rate_limit_hits` | 列显示「待 YAML」或隐藏 | **不得** 填假数 |
| `key_value` | **永不展示** | 安全红线 |
| 403 | 无平台 audit 读权限 | — |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `rate_limit_rpm` | 1–10000（create 默认 60） | requests/minute |
| `last_used_at` / `created_at` / `expires_at` | ISO 8601 | date-time |
| `scopes` | `scope:{domain}:{action}` pattern | string[] |
| `key_prefix` | 脱敏前缀 | string |
| `call_count` | **待 YAML** | integer |
| `rate_limit_hits` | **待 YAML** | integer |
| `status_code` | HTTP 状态码 | integer |

## 状态与能力口径

### APIKeyInfo.is_active

| 值 | 含义 | 典型原因 |
|---|---|---|
| `true` | 有效 | 正常使用 |
| `false` | 已吊销 | `revokeAPIKey` 或过期策略 |

### 调用事件 result（草案 · 待 YAML）

| 维度 | 说明 |
|---|---|
| `status_code` 2xx | 成功 |
| `status_code` 4xx/5xx | 失败 / 限流 |
| 限流命中 | 依赖 **待 YAML** `rate_limit_hits` 或事件 metadata |

### 状态 × 操作可用性矩阵

| 操作 \ Key 状态 | is_active=true | is_active=false |
|---|---|---|
| 查看 Key 摘要 | ✅ **待 YAML** | ✅ |
| 查看调用事件 | ✅ **待 YAML** | ✅（历史） |
| 跳转 Trace | ✅（有 request_id） | ✅ |
| 平台代吊销 | ⏳ Phase 2 | ❌ 已吊销 |
| 导出 | ⏳ Phase 2 | ⏳ |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 apikey audit 读 RBAC | **TODO-YAML** | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| 事件时间窗 | start ≤ end | `400 BAD_REQUEST` |
| `key_id`（单 Key 事件） | 合法 ID | `404`（待 YAML） |
| 平台代吊销（Phase 2） | 写 RBAC + `idempotency_key` | `403` / `422` 待声明 |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 跨租户 Key list | ✅ **待 YAML** | ✅ | ✅ | platform list |
| 调用事件 list | ✅ **待 YAML** | ✅ | ✅ | 非 APIKeyInfo |
| 查看 Key 详情 | ✅ **待 YAML** | ✅ | ✅ | — |
| 跳转 Trace / 审计日志 | ✅ | ✅ | ✅ | request_id |
| 平台代吊销 | ❌ | Phase 2 | Phase 2 | **待 YAML** |
| 导出取证 | ❌ | Phase 2 | Phase 2 | compliance-export |

## 删除前置校验与当前契约边界

租户 `revokeAPIKey`（`DELETE /auth/api-keys/{key_id}`）为 **逻辑吊销**，非物理 DELETE 审计记录。

| 场景 | 校验 | 响应 |
|---|---|---|
| BOSS 本页 | **无 DELETE** 审计事件操作 | **N/A** |
| 租户 revoke（参考） | Key 存在且 `is_active=true` | `200`；已吊销可能 `404`/`422`（待 YAML 细化） |

BOSS 平台代吊销 — **TODO-YAML** Phase 2；合入前不得写为已实现。

## 接口冻结规则

### `GET /api/v1/auth/api-keys`（租户 · **只读参考**）

- `operationId`：`listAPIKeys`
- `summary`：`列出当前租户 API Key`
- `tags`：`["Auth"]`
- `parameters`：`user_id`（可选 query）
- `success`：`200 + ListAPIKeysResponse`
- `errors`：`401`、`403`
- **非** BOSS 跨租户正式契约；**禁止** 平台 dashboard JWT 轮询

### `POST /api/v1/auth/api-keys`（租户 · **只读参考**）

- `operationId`：`createAPIKey`
- `requestBody`：`CreateAPIKeyRequest`（`name`, `scopes` 必填；`rate_limit_rpm` 默认 60）
- `success`：`201 + CreateAPIKeyResponse`（含一次性 `key_value`）
- `errors`：`400`、`401`、`403`

### `DELETE /api/v1/auth/api-keys/{key_id}`（租户 · **只读参考**）

- `operationId`：`revokeAPIKey`
- `path.required`：`key_id`
- `success`：`200`
- `errors`：`400`、`401`、`403`、`404`

### Console 草案 `GET .../audit-events`（**草案/待 YAML · 非冻结**）

- 规划见 `docs/console-modules/openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md`
- `operationId`（草案）：`listApiKeyAuditEvents`
- RBAC（草案）：`scope:auth:api-keys:audit:read`
- **合入 v1.yaml 前不得写入「已冻结」正文**

### 平台跨租户 API Key 审计 list（待补）

<!-- ADDED-TO-YAML: GET /api/v1/platform/audit/api-keys (`listPlatformAPIKeyAuditEvents`) -->

- 须 platform RBAC + cursor 分页
- 须含 `tenant_id` 与 `APIKeyInfo` 字段
- 调用统计字段须新 schema — **不得** 假造在 `APIKeyInfo` 上

## 使用规则

- **不得** 把 `listAPIKeys` 写成 BOSS 跨租户 list 已实现
- **不得** 展示 `call_count` / `rate_limit_hits` 假数据（schema 当前 **无**）
- **不得** 展示 `key_value` 明文
- Console 草案 audit-events path 引用须标 **草案/待 YAML**
- 平台 list 未上线时 **禁止** JWT 逐租户轮询
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 平台跨租户 API Key 事件 list — **ADDED-TO-YAML**
- `call_count` / `rate_limit_hits` schema — P1
- Console `audit-events` 合入 v1.yaml — Phase 3
- 平台代吊销 — Phase 2 + `idempotency_key`
- 与 [`platform-audit-log.md`](platform-audit-log.md) 统一事件 schema — P1

## 响应示例

### listAPIKeys 引用（租户 · **非 BOSS 正式响应**）

```json
{
  "items": [
    {
      "id": "key-001",
      "name": "ci-deploy",
      "key_prefix": "ani_k_abc",
      "scopes": ["scope:inference:read"],
      "rate_limit_rpm": 120,
      "created_at": "2026-06-01T08:00:00Z",
      "expires_at": null,
      "last_used_at": "2026-06-16T09:30:00Z",
      "is_active": true
    }
  ],
  "total": 1
}
```

### 平台 API Key 审计 list 目标（**待 YAML**）

```json
{
  "items": [
    {
      "tenant_id": "t-001",
      "id": "key-001",
      "name": "ci-deploy",
      "key_prefix": "ani_k_abc",
      "scopes": ["scope:inference:read"],
      "rate_limit_rpm": 120,
      "last_used_at": "2026-06-16T09:30:00Z",
      "is_active": true
    }
  ],
  "next_cursor": null,
  "total": 42
}
```

## 错误示例

### 时间窗非法（audit-events · **待 YAML**）

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be before or equal to end_time",
  "request_id": "req-boss-apikey-400-001"
}
```

### 无平台 API Key 审计读权限（**TODO-YAML**）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-apikey-403-001"
}
```

## 相关模块

- Console：[`api-key-audit.md`](../../console-modules/tenant/api-key-audit.md)、[`api-key-management.md`](../../console-modules/tenant/api-key-management.md)（若存在）
- [`platform-audit-log.md`](platform-audit-log.md)
- [`platform-trace.md`](../health/platform-trace.md)
- [`compliance-export.md`](compliance-export.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] `APIKeyInfo` 字段与 v1.yaml 一致；**无** call_count/rate_limit_hits 已标注
- [x] 租户 auth paths **只读参考**；禁止 JWT 轮询
- [x] Console audit-events 标 **草案/待 YAML**
- [x] 含字段展示规则、口径、状态矩阵、400+403 错误示例
- [x] 删除前置校验（revoke 边界 + BOSS N/A）已声明
- [ ] platform apikey audit YAML 合入后回写冻结表
- [x] PRD/SPEC/HTML 与本文同步
