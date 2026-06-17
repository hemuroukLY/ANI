# 平台运营系统对接

## 页面定位

`平台运营系统对接` 是 BOSS **平台集成与通知** 域下的 **跨租户运营系统连接器** 专页：配置与治理 ANI 平台与 **计费、工单、CMDB、ITSM** 等外部运营系统的对接，提供全平台集成实例 list、连通性状态、配置摘要与跨租户同步视图，供平台运营、交付与 SRE 统一治理。

**本页 ≠ Console 租户第三方集成。** Services `GET/POST /api/v1/svc/integrations`（`listIntegrations` / `createIntegration`）为 **租户上下文** 业务系统集成（见 Console [`integration-third-party.md`](../../console-modules/integration/integration-third-party.md)）；本页为 **平台 RBAC + 跨租户 aggregate**，对接对象服务于 **平台级** 计费对账、工单同步、资产 CMDB 等能力。

本页属于 **Core / Integration** 视角下的 **平台 RBAC** 页面。一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`（**paths 段无** `/platform/integrations*`）；Services 仅声明租户 integrations list/create，**无** `{integration_id}` GET/PATCH/DELETE。

Console [`integration-webhook-overview.md`](../../console-modules/integration/integration-webhook-overview.md) 为接入域导航；[`integration-bot.md`](../../console-modules/integration/integration-bot.md) 为 Bot 子类型 — Bot **不** 纳入本页（见 [`enterprise-notification.md`](enterprise-notification.md)）。

## 文档管理规则

- 本文是 `平台运营系统对接` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-ops-system-integration.md`](../../tasks/modules/prd/boss/integration/prd-boss-ops-system-integration.md) 与 [`spec-boss-ops-system-integration.md`](../../tasks/modules/spec/boss/integration/spec-boss-ops-system-integration.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml` + `services/v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)
- Console 对照：[`integration-third-party.md`](../../console-modules/integration/integration-third-party.md)

## Core 层要求

- 平台运营系统对接归属 **Core**；平台跨租户 list/CRUD — **TODO-YAML** **P2**
- Services `/api/v1/svc/integrations*` — **租户只读参考**；**无** integration_id 级 path — BOSS **不得** 自造 tenant `{integration_id}` 路径为平台正式 API
- **禁止** 逐租户 JWT 轮询 `listIntegrations` 冒充平台跨租户 aggregate
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id` 越权筛选（平台连接器默认 **无** 单租户归属或为 platform-scope）
- 写操作 POST/PUT/PATCH/DELETE（待 YAML）须 `idempotency_key`（UUID）
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）；tenant `createIntegration` **当前 YAML 未声明 422**
- 禁止自造 `PlatformOpsIntegration` schema / operationId / 路径为已冻结
- OpenAPI 已声明 ≠ handler 已实现
- **禁止 JWT polling** 拼装平台 dashboard 跨租户视图

## 页面职责

- 提供 **全平台** 运营系统连接器 list 与多维筛选（list API 待 YAML P2）
- 支持创建/更新/启停平台级对接（计费/工单/CMDB provider）（写 API 待 YAML P2）
- 展示 `provider`、`status`（active/inactive/error）、配置摘要（**不** 暴露 secret 明文）
- 支持 **连通性测试**（test API 待 YAML P2）
- 提供 **跨租户同步摘要**（如对账任务数、未同步工单数 — aggregate 待 YAML P2）
- 与 [`incident-handling.md`](../health/incident-handling.md) 联动：工单系统对接可升级/同步 Incident（API 待 YAML）
- 与 [`platform-alerts-pending.md`](../overview/platform-alerts-pending.md) 联动：对接失败可产生平台待办
- 与 [`ops-webhook.md`](ops-webhook.md) / [`enterprise-notification.md`](enterprise-notification.md) 对照：出站通知 vs 双向系统对接边界

## 页面结构

- 首屏至少包含：`筛选区`、`连接器表格`、`创建/编辑抽屉`、`连通性测试`、`同步摘要卡片`、`边界说明`
- 无数据态、无权限态、API 未就绪态、status=error 态须可区分

```text
平台运营系统对接
├── 筛选（provider / status / 关键字 — list 待 YAML）
├── 连接器表格（名称 / provider / 状态 / 最近同步）
├── 创建 / 编辑抽屉（provider / name / config 摘要）
├── 连通性测试（test API 待 YAML P2）
├── 跨租户同步摘要（aggregate 待 YAML P2）
│   └── 跳转：incident-handling / platform-alerts-pending
└── 边界说明（租户 integrations 见 Console）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 / 模块 | 本页用法 |
|---|---|---|
| Core | 平台 ops integration list/CRUD **TODO-YAML** P2 | BOSS 正式数据源 |
| Core | 平台 ops integration test **TODO-YAML** P2 | 连通性测试 |
| Core | 跨租户 sync aggregate **TODO-YAML** P2 | 摘要卡片 |
| Services | `GET /api/v1/svc/integrations` | **租户只读参考** |
| Services | `POST /api/v1/svc/integrations` | **租户只读参考** |
| 产品 | incident-handling / platform-alerts-pending | 工单/Incident 与待办 |

### 关键边界

- Core v1.yaml **无** platform integrations path — 不得写「已实现平台对接 list」
- Services Integration 项：`id, name, provider, status` — **可参考** enum，平台 schema 须独立（含 provider 扩展：billing/ticket/cmdb 等）
- Services **无** `GET/PATCH/DELETE /integrations/{integration_id}` — 租户详情/更新/删除 **YAML 未声明**
- **禁止** 生产环境逐租户 JWT 轮询 aggregate
- OAuth 回调 / 第三方授权 — Gateway/Services **待规划**；本页 **不** 自造 OAuth path 为已冻结
- `config` 对象 additionalProperties — 租户 POST 已声明；平台 config **须** 密钥托管策略（UI 不回显 secret）

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 筛选区 | UI + list API **待 YAML** | provider / status | 刷新表格 |
| 连接器表格 | list API **待 YAML** | 平台连接器摘要 | 详情 / 编辑 |
| 创建/编辑 | write API **待 YAML** | provider / name / config | — |
| 连通性测试 | test API **待 YAML** | 单次 probe 结果 | — |
| 同步摘要 | aggregate **待 YAML** | 跨租户计数 | incident-handling |
| 边界说明 | 规划项 | 租户 integrations ≠ 本页 | integration-third-party |

## BOSS 与 Console 分工

| 维度 | BOSS 平台运营系统对接 | Console 第三方集成 |
|---|---|---|
| 范围 | 全平台运营系统（计费/工单/CMDB） | 当前租户业务 provider |
| List | platform ops integration list **待 YAML** P2 | `listIntegrations`（Services 已声明） |
| Create | platform POST **待 YAML** P2 | `createIntegration`（Services 已声明） |
| Get/Update/Delete | **待 YAML** P2 | **YAML 未声明** `{integration_id}` |
| status enum | active / inactive / error（目标对齐） | active / inactive / error（Services 已声明） |
| 跨租户视图 | 单 API aggregate **待 YAML** | 单租户 JWT |
| 连通性测试 | platform test **待 YAML** P2 | 待补 |
| RBAC | 平台 ops integration scope **待 YAML** | 租户集成 scope |
| Bot 子类型 | **不** 承担 | 见 integration-bot |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| — | `/platform/integrations*` 或等价 | **不存在** | Core 无 platform ops integration path |
| GET | `/api/v1/svc/integrations` | `listIntegrations` | **Services 租户参考** |
| POST | `/api/v1/svc/integrations` | `createIntegration` | **Services 租户参考** |
| — | `/integrations/{integration_id}` | **未声明** | 租户 get/patch/delete 均无 |

| 能力 | 状态 |
|---|---|
| 平台 ops integration list | **TODO-YAML** **P2** |
| 平台 ops integration create/update/delete | **TODO-YAML** **P2** |
| 连通性 test | **TODO-YAML** **P2** |
| 跨租户 sync aggregate | **TODO-YAML** **P2** |
| `PlatformOpsIntegration` schema | **TODO-YAML** **P2** |

### 页面目标字段（待 YAML 合入后对齐）

| 字段 | 说明 |
|---|---|
| `id` | 连接器 UUID |
| `name` | 展示名称 |
| `provider` | billing / jira / servicenow / cmdb / custom 等（待 enum） |
| `status` | active / inactive / error |
| `config_summary` | 非敏感配置摘要 |
| `created_at` | 创建时间 |
| `last_sync_at` | 最近同步时间；可空 |
| `last_sync_status` | success / failure；可空 |
| `linked_tenant_count` | 关联租户数（aggregate；可空） |

### Services 租户 Integration 参考（**非 BOSS 正式**）

#### listIntegrations / createIntegration 项

| 字段 | 说明 |
|---|---|
| `id` | UUID |
| `name` | string |
| `provider` | string |
| `status` | active / inactive / error |

#### createIntegration 请求

| 字段 | 必填 | 说明 |
|---|---|---|
| `idempotency_key` | 是 | UUID |
| `name` | 是 | 名称 |
| `provider` | 是 | provider 字符串 |
| `config` | 否 | object, additionalProperties |

## 字段级定义

### 查询字段（平台 list API 目标 · 待 YAML）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `provider` | query | 可选 | provider 筛选 |
| `status` | query | 可选 | active/inactive/error |
| `q` | query | 可选 | 名称关键字 |
| `limit` | query | 可选 | 分页 |
| `cursor` | query | 可选 | 游标 |

### 写入字段（create/update · 待 YAML）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `idempotency_key` | body | POST/PATCH 必填 | UUID |
| `name` | body | 必填 | 名称 |
| `provider` | body | 必填 | 运营系统 provider |
| `config` | body | 条件 | 连接配置；secret 字段加密存储 |
| `active` | body | 可选 | 映射 status 或独立字段（待 YAML） |

### 返回字段（PlatformOpsIntegration · **待 YAML**）

| 字段 | 来源 | 说明 |
|---|---|---|
| `id` | schema | UUID |
| `name` | schema | 名称 |
| `provider` | schema | provider |
| `status` | schema | active / inactive / error |
| `config_summary` | schema | 非敏感摘要 |
| `created_at` | schema | ISO 8601 |
| `updated_at` | schema | 可空 |
| `last_sync_at` | schema | 可空 |
| `last_sync_status` | schema | 可空 |
| `last_error` | schema | error 态摘要；可空 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `provider_label` | 计费 / Jira / CMDB 等中文名 |
| `status_badge` | active 绿 / inactive 灰 / error 红 |
| `sync_age_display` | 「3 小时前同步」 |
| `config_keys_display` | 已配置项 checklist（无 secret） |
| `pending_items_link` | 链接 platform-alerts-pending（若对接失败产生待办） |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 筛选 + 表格 + 详情/编辑 | list 待 YAML |
| 无数据态 | 「暂无平台运营系统对接」 | **不** 伪造行 |
| API 未就绪 | 「平台运营系统对接 API 待 Core 合入（P2）」 | 不得用 tenant integrations 冒充 |
| 无权限态 | 403 提示 | 平台 RBAC |
| `status=error` | 红色 + `last_error` | 引导测试或编辑 |
| `status=inactive` | 灰色停用 | — |
| `config` / secret | 仅「已配置 API Key」等 | **不** 明文 |
| provider 非法 | 前端 + 后端 enum | 400 |
| 跨租户摘要不可用 | 卡片「aggregate 待 YAML」 | **不** JWT 轮询伪造数字 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `provider` | 平台运营 provider enum（待 YAML） | string |
| `status` | active / inactive / error | enum（对齐 Services Integration） |
| `config` | JSON object；secret 键托管 | object |
| `created_at` / `updated_at` | 后端原值 | ISO 8601 date-time |
| `last_sync_at` | 最近同步完成时间 | ISO 8601；可空 |
| `linked_tenant_count` | 非负整数 | count（aggregate 待 YAML） |
| `last_sync_status` | success / failure | enum（待 YAML） |

## 状态与能力口径

### PlatformOpsIntegration.status（页面目标 · 待 YAML；参考 Services）

| 值 | 含义 | UI |
|---|---|---|
| `active` | 对接启用且最近一次 probe/sync 正常 | 绿色 |
| `inactive` | 已停用 | 灰色 |
| `error` | 连通性或同步失败 | 红色 |

### last_sync_status（aggregate · 待 YAML）

| 值 | 含义 | UI |
|---|---|---|
| `success` | 最近同步成功 | 绿色 |
| `failure` | 最近同步失败 | 红色 |

### 状态 × 操作可用性矩阵

| 操作 \ 数据就绪 | list 待 YAML | list 已合入 |
|---|---|---|
| 列表筛选 | ⏳ 未就绪态 | ✅ |
| 创建连接器 | ⏳ | ✅ |
| 编辑/启停 | ⏳ | ✅ |
| 连通性测试 | ⏳ P2 | ⏳ P2 |
| 查看 sync 摘要 | ⏳ P2 | ⏳ P2 |
| 删除连接器 | ⏳ | ✅（待 YAML DELETE） |
| Incident 工单同步 | ⏳ P2 | ⏳ P2 |

⏳ = API 待 YAML 合入。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 ops integration 读 RBAC | 已授权（**TODO-YAML**） | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| 创建写 RBAC | 平台 integration 写 **待 YAML** | `403 FORBIDDEN` |
| `idempotency_key` | POST 携带 UUID | `400 BAD_REQUEST` |
| `name` / `provider` | 非空合法 | `400 BAD_REQUEST` |
| `config` 必填项（按 provider） | 产品规则 **待 YAML** | `400 BAD_REQUEST` |
| 重复 provider+name（若禁止） | **待 YAML** | `409 CONFLICT` |

## 操作可用性矩阵

| 操作 | 平台只读 | 交付/运营 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|---|
| 列表（待 YAML） | ✅ | ✅ | ✅ | ✅ | platform list |
| 创建连接器 | ❌ | ✅ **待 YAML** | ✅ **待 YAML** | ✅ **待 YAML** | POST |
| 编辑/启停 | ❌ | ✅ **待 YAML** | ✅ **待 YAML** | ✅ **待 YAML** | PATCH |
| 连通性测试 | ❌ | ✅ **待 YAML** P2 | ✅ | ✅ | test POST |
| 查看 sync 摘要 | ✅ **待 YAML** | ✅ | ✅ | ✅ | aggregate |
| 删除连接器 | ❌ | ❌ | ✅ **待 YAML** | ✅ **待 YAML** | DELETE |
| Incident 同步 | ❌ | P2 | ✅ P2 | ✅ P2 | incident-handling |
| 租户 integrations | 跳转 Console | 跳转 | 跳转 | 跳转 | integration-third-party |

## 删除前置校验与当前契约边界

DELETE 平台 ops integration — **TODO-YAML** P2。预期前置校验：

| 校验项 | 未通过响应 | 说明 |
|---|---|---|
| 连接器存在 | `404 NOT_FOUND` | — |
| 平台写 RBAC | `403 FORBIDDEN` | — |
| 存在进行中的 sync 任务（若 YAML 声明） | `422 PreconditionFailed` | operation 声明后可写冻结 |
| 连接器为 Incident 默认工单路由（若 YAML 声明） | `409` / `422` | 须先解除绑定 |
| 计费对账周期内删除（若产品禁止） | `422` | 待产品 + YAML |

Services **无** tenant integration DELETE — **不得** 引用为平台 DELETE 契约。

## 接口冻结规则

### Core · 平台运营系统对接（待补 · **P2**）

<!-- TODO-YAML: GET/POST /api/v1/platform/ops-integrations 或等价 -->

- 须平台 RBAC
- list query：`provider`、`status`、`q`、`limit`、`cursor`
- create 须 `idempotency_key`、`name`、`provider`
- **合入前不得写入「已冻结」**

### Core · 连通性 test（待补 · **P2**）

<!-- TODO-YAML: POST .../ops-integrations/{id}/test -->

- 须 `idempotency_key`
- 错误：`400`、`401`、`403`、`404`

### Core · 跨租户 sync aggregate（待补 · **P2**）

<!-- TODO-YAML: GET /api/v1/platform/ops-integrations/sync-summary 或等价 -->

- **禁止** 用 tenant list 长度求和冒充

### `GET /api/v1/svc/integrations`（Services · **租户参考 · 非 BOSS 正式**）

- `operationId`：`listIntegrations`
- 项 required：`id, name, provider, status`
- status enum：`active`, `inactive`, `error`
- 成功：`200`
- 错误：`401`、`403`
- **BOSS 不得调用做平台 dashboard**

### `POST /api/v1/svc/integrations`（Services · **租户参考**）

- `operationId`：`createIntegration`
- 必填：`idempotency_key`, `name`, `provider`
- 可选：`config` (object)
- 成功：`201`
- 错误：`400`、`401`、`409`
- **当前 YAML 未声明 422**

## 使用规则

- **不得** 写 platform ops integration list 为已实现
- **不得** 把 Services tenant integrations 路径写成 BOSS 正式 API
- **禁止** 逐租户 JWT 轮询 `listIntegrations` 冒充平台 aggregate
- 租户 integration 详情/更新/删除 — YAML **未声明**；仅 Console 文档引用
- OAuth 授权流程 — **待规划**；不得自造 callback path 为已冻结
- 对接失败待办 — 见 [`platform-alerts-pending.md`](../overview/platform-alerts-pending.md)
- Incident 工单同步 — 见 [`incident-handling.md`](../health/incident-handling.md)
- 出站 Webhook / IM — 见 [`ops-webhook.md`](ops-webhook.md)、[`enterprise-notification.md`](enterprise-notification.md)
- `config` secret — UI **永不** 明文展示或日志打印
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 平台 ops integration list/CRUD — **TODO-YAML** **P2**
- 连通性 test — **TODO-YAML** **P2**
- 跨租户 sync aggregate — **TODO-YAML** **P2**
- provider enum（billing/jira/cmdb/…）— P2
- OAuth 第三方授权 — Gateway/Services 规划 P2
- Incident ↔ 工单双向同步 — P2 + incident-handling
- 租户 integration GET/PATCH/DELETE — Services **YAML 未声明**

## 响应示例

### 平台 ops integration list 成功（页面目标 · **待 YAML**）

```json
{
  "items": [
    {
      "id": "poi-001",
      "name": "集团 Jira 工单",
      "provider": "jira",
      "status": "active",
      "config_summary": {"base_url": "https://jira.example.com", "project_key": "ANI"},
      "created_at": "2026-05-20T09:00:00Z",
      "last_sync_at": "2026-06-17T07:00:00Z",
      "last_sync_status": "success",
      "linked_tenant_count": 48
    }
  ],
  "next_cursor": "cursor-poi-page-2",
  "total": 6
}
```

### 连通性 test 成功（页面目标 · **待 YAML P2**）

```json
{
  "probe_id": "probe-501",
  "status": "success",
  "latency_ms": 320,
  "tested_at": "2026-06-17T10:15:00Z"
}
```

## 错误示例

### provider 缺失（create · **待 YAML**）

```json
{
  "code": "BAD_REQUEST",
  "message": "provider is required",
  "request_id": "req-boss-poi-400-001"
}
```

### 无平台 ops integration 写权限（**TODO-YAML**）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-poi-403-001"
}
```

> **注**：适用于 **platform ops integration write/list（TODO-YAML）** 场景。

## 相关模块

- [`ops-webhook.md`](ops-webhook.md)（HTTPS 出站 · P1）
- [`enterprise-notification.md`](enterprise-notification.md)（IM 通知 · P2）
- [`incident-handling.md`](../health/incident-handling.md)（工单/Incident 同步）
- [`platform-alerts-pending.md`](../overview/platform-alerts-pending.md)（对接失败待办）
- Console：[`integration-third-party.md`](../../console-modules/integration/integration-third-party.md)、[`integration-webhook-overview.md`](../../console-modules/integration/integration-webhook-overview.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 明确 Core **无** platform ops integration path；tenant integrations **非** BOSS 正式
- [x] 含字段展示规则、字段口径与单位、状态与能力口径
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [x] 删除前置校验已声明（待 YAML DELETE）
- [ ] platform ops integration YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
