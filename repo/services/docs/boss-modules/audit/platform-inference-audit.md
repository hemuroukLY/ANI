# 推理调用审计

## 页面定位

`推理调用审计` 是 BOSS **安全审计与合规** 域下的 **跨租户推理调用计量与状态审计** 专页：汇总全平台推理请求的 Token 用量、模型、调用方与请求状态摘要，**不展示 prompt / completion 全文**，供平台管理员、SRE 与计费运营做滥用排查与容量分析。

本页属于 **Core / Metering + Audit** 视角下的 **平台 RBAC** 页面。`POST /api/v1/metering/token-usage`（`reportTokenUsage`）为 Services **上报写入** 路径，返回 `TokenUsageReport` — **不是** BOSS 审计 list；跨租户 inference audit list — **ADDED-TO-YAML**。

Console **无直接对等页**（平台专属能力）。租户侧用量只读参考：`GET /api/v1/metering/usage`。

## 文档管理规则

- 本文是 `推理调用审计` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-platform-inference-audit.md`](../../tasks/modules/prd/boss/audit/prd-boss-platform-inference-audit.md) 与 [`spec-boss-platform-inference-audit.md`](../../tasks/modules/spec/boss/audit/spec-boss-platform-inference-audit.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- Token 计量上报归属 **Core** `POST /api/v1/metering/token-usage`（`scope:metering:write`）
- 租户用量统计：`GET /api/v1/metering/usage` — **只读参考**；JWT 租户上下文
- 平台跨租户 inference audit list — **ADDED-TO-YAML**
- **不得** 把 `reportTokenUsage` 响应当 list 数据源；**不得** 把 `GET /metering/usage` 聚合当作逐请求审计 list
- 平台 RBAC 鉴权；**不得** 信任未授权 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422` 仅 YAML 已声明 operation 可写冻结（§2.10）
- **禁止** 展示 prompt 全文 — 产品红线
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **全平台** 推理调用审计 list 与筛选（list API 待 YAML）
- 展示 Token 用量摘要：`input_tokens`、`output_tokens`、`total_tokens`、`model`、`source`
- 展示调用元数据：`tenant_id`、`request_id`、上报 `state`（`accepted` / `duplicate`）
- 支持按租户、模型、source、时间、request_id 筛选
- 支持 **`request_id` 跳转 Trace**（见 [`platform-trace.md`](../health/platform-trace.md)）
- 与租户 metering 用量页对照，明确 BOSS 跨租户逐事件 vs 聚合差异

## 页面结构

- 首屏至少包含：`筛选区`、`调用审计表格`、`详情抽屉`、`用量摘要卡片`、`边界说明`
- 无数据态、无权限态、list API 未就绪态须可区分

```text
推理调用审计
├── 筛选（tenant_id / model / source / 时间 / request_id — list 待 YAML）
├── 摘要卡片（总 tokens / 调用次数 — 待 YAML 聚合）
├── 调用表格（时间 / 租户 / model / tokens / request_id / state）
├── 详情抽屉（TokenUsageReport 字段 · 无 prompt）
│   └── 跳转 platform-trace
└── 边界：reportTokenUsage=写入；usage=租户聚合参考
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | platform inference audit list **TODO-YAML** | BOSS 正式跨租户逐事件 list |
| Core | `POST /api/v1/metering/token-usage` | **写入/report**；`TokenUsageReport` 为 202 响应，非 list |
| Core | `GET /api/v1/metering/usage` | 租户 **只读参考**；按 resource_type 聚合统计 |
| Services | 推理网关上报 | 经 Core `reportTokenUsage` 入库 |

### 关键边界

- `reportTokenUsage`：`202 + TokenUsageReport`；`state` 仅 `accepted` | `duplicate` — 表上报受理态，**非** HTTP 业务 status_code
- `GET /metering/usage`：必填 `start_time`/`end_time`；返回聚合 `items[]` — **非** 逐请求审计
- `TokenUsageReport` **无** `duration_ms`、`status_code`、`prompt` — UI **不得** 伪造（若需由 **待 YAML** 扩展 schema）
- BOSS list 未就绪时 **禁止** 逐租户 JWT 轮询 usage 或伪造 POST 查询

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 筛选区 | UI + list **待 YAML** | tenant / model / 时间 | 刷新 |
| 摘要卡片 | list 聚合 **待 YAML** | tokens 合计 | — |
| 调用表格 | list **待 YAML** | 逐事件 | 详情 |
| 详情抽屉 | get **待 YAML** 或 list item | TokenUsageReport 字段 | Trace |
| 边界说明 | 冻结事实 | report ≠ list | metering usage（租户） |

## BOSS 与 Console 分工

| 维度 | BOSS 推理调用审计 | Console / 租户 |
|---|---|---|
| 范围 | 全平台多租户逐事件 | 当前租户 |
| List | platform inference audit **待 YAML** P1 | 无对等页 |
| 用量聚合 | 平台摘要（待 YAML） | `GET /metering/usage` |
| 上报 | 不经过 BOSS UI | Services → `reportTokenUsage` |
| Prompt 全文 | **禁止展示** | **禁止展示** |
| RBAC | 平台 audit 读（待 YAML） | `scope:metering:write`（上报） |

## 当前冻结事实

| 方法 | 路径 | operationId | RBAC | BOSS 用法 |
|---|---|---|---|---|
| POST | `/api/v1/metering/token-usage` | `reportTokenUsage` | `scope:metering:write` | **写入**；非 list |
| GET | `/api/v1/metering/usage` | **未声明** | 租户 JWT | **只读参考**（聚合） |

| Schema | 字段（冻结） |
|---|---|
| `TokenUsageReport` | `id`, `tenant_id`, `source`, `model`, `input_tokens`, `output_tokens`, `total_tokens`, `request_id`, `instance_id`, `state`, `created_at` |
| `TokenUsageReport.state` | `accepted` / `duplicate` |
| `ReportTokenUsageRequest` | `idempotency_key`, `source`, `model`, `input_tokens`, `output_tokens`（+ 可选 `request_id`, `instance_id`, `occurred_at`, `labels`） |

| 能力 | 状态 |
|---|---|
| 平台跨租户 inference audit list | **ADDED-TO-YAML** P1 |
| `duration_ms` / HTTP `status_code` | **TODO-YAML**（当前 schema **无**） |
| 平台 Token 用量 CSV 导出 | Phase 2 |

## 字段级定义

### 查询字段（platform list 目标 · 待 YAML）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `tenant_id` | query | 可选 | 平台 RBAC 筛选 |
| `model` | query | 可选 | 模型名 |
| `source` | query | 可选 | 上报来源 |
| `request_id` | query | 可选 | 精确 Trace 关联 |
| `state` | query | 可选 | accepted / duplicate |
| `start_time` / `end_time` | query | 可选 | 时间窗 |
| `limit` / `cursor` | query | 可选 | 分页 |

### 查询字段（GET /metering/usage · 租户 **只读参考**）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` | query | ✅ | ISO 8601 |
| `end_time` | query | ✅ | ISO 8601 |
| `resource_type` | query | 可选 | 资源类型 |
| `group_by` | query | 可选 | resource_type / az / day / hour |

### 返回字段（TokenUsageReport · Core 已冻结 · 上报响应）

| 字段 | 来源 | 说明 |
|---|---|---|
| `id` | `TokenUsageReport` | 记录 ID |
| `tenant_id` | `TokenUsageReport` | 租户 |
| `source` | `TokenUsageReport` | 如 model-service |
| `model` | `TokenUsageReport` | 模型标识 |
| `input_tokens` | `TokenUsageReport` | 输入 tokens |
| `output_tokens` | `TokenUsageReport` | 输出 tokens |
| `total_tokens` | `TokenUsageReport` | 合计 |
| `request_id` | `TokenUsageReport` | 可空 |
| `instance_id` | `TokenUsageReport` | 可空 |
| `state` | `TokenUsageReport` | accepted / duplicate |
| `created_at` | `TokenUsageReport` | 创建时间 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `caller_display` | API Key / 用户 — **待 YAML** |
| `duration_display` | **待 YAML** schema |
| `status_code_display` | **待 YAML**；当前 schema **无** |
| `tokens_formatted` | 千分位 + 单位 tokens |
| `duplicate_badge` | state=duplicate 时标注「重复上报」 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 表格 + 详情 | list 待 YAML |
| 无数据态 | 「暂无推理调用记录」 | **不** 伪造 |
| list 未就绪 | 「跨租户 inference audit 待 Core 合入」 | 不得用 usage 聚合冒充 |
| `state=duplicate` | 灰色「重复上报」标签 | 非错误 |
| prompt / completion | **永不展示** | 安全红线 |
| `request_id` 非空 | 可跳转 Trace | — |
| 403 | 无平台读权限 | — |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `input_tokens` / `output_tokens` / `total_tokens` | int64 ≥ 0 | tokens |
| `model` / `source` | 后端原值 | string |
| `state` | accepted / duplicate | enum |
| `created_at` / `occurred_at` | ISO 8601 | date-time |
| `duration_ms` | **待 YAML** | milliseconds |
| `status_code` | **待 YAML** | HTTP integer |
| usage `total_quantity` | 聚合数值 | 与 `unit` 字段一致 |

## 状态与能力口径

### TokenUsageReport.state（上报受理 · 已冻结）

| 状态 | 含义 | UI |
|---|---|---|
| `accepted` | 首次受理 | 默认 |
| `duplicate` | 幂等重复 | 灰色标签 |

> **不是** 推理 HTTP 响应码；业务成功/失败若需展示 — **待 YAML** 扩展字段。

### 状态 × 操作可用性矩阵

| 操作 | list 待 YAML | list 已合入 |
|---|---|---|
| 逐事件 list | ⏳ | ✅ |
| 查看详情 | ⏳ | ✅ |
| 跳转 Trace | ✅（request_id） | ✅ |
| 租户 usage 聚合参考 | ✅ 只读（非本页主表） | ✅ |
| 导出 CSV | ⏳ Phase 2 | ⏳ |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 inference audit 读 RBAC | **TODO-YAML** | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| list 时间窗 | start < end | `400 BAD_REQUEST` |
| usage 参考查询 | start_time + end_time 必填 | `400 BAD_REQUEST` |
| reportTokenUsage（Services 上报） | `scope:metering:write` + idempotency_key | `403` / `400` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 跨租户调用 list | ✅ **待 YAML** | ✅ | ✅ | 主数据源 |
| 查看详情 | ✅ **待 YAML** | ✅ | ✅ | — |
| 跳转 Trace | ✅ | ✅ | ✅ | request_id |
| 查看租户 usage 聚合 | ✅ 参考 | ✅ | ✅ | 非逐事件 |
| 手动 reportTokenUsage | ❌ | ❌ | ❌ | Services 专用 |
| 导出 | ❌ | Phase 2 | Phase 2 | **待 YAML** |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。Token 用量记录为 append-only 计量事件；保留/归档策略 **待 YAML** 合规定义。

## 接口冻结规则

### `POST /api/v1/metering/token-usage`（Core · **已声明 · 写入非 list**）

- `operationId`：`reportTokenUsage`
- `summary`：`上报 Token 用量`
- `tags`：`["Metering"]`
- `x-ani-rbac-scope`：`scope:metering:write`
- `requestBody`：`ReportTokenUsageRequest`
- `success`：`202 + TokenUsageReport`
- `errors`：`400`、`401`、`403`
- BOSS 审计页 **不得** 调用此 API 做 list 查询

### `GET /api/v1/metering/usage`（Core · **只读参考 · 聚合**）

- `operationId`：**未声明**
- `summary`：`查询租户用量统计`
- `tags`：`["Metering"]`
- `parameters`：`start_time`、`end_time`（必填）；`resource_type`、`group_by`
- `success`：`200` + items 聚合
- `errors`：`400`、`401`、`403`
- **非** BOSS 跨租户逐事件 list；**禁止** JWT 轮询拼平台 dashboard

### 平台 inference audit list（待补）

<!-- ADDED-TO-YAML: GET /api/v1/platform/audit/inference-calls (`listPlatformInferenceAuditEvents`) -->

- 须 platform RBAC + cursor
- 须含 `tenant_id` 与 `TokenUsageReport` 对齐字段
- 扩展 `duration_ms` / `status_code` 须新 YAML 字段

## 使用规则

- **不得** 把 `reportTokenUsage` 或 `GET /metering/usage` 写成 BOSS audit list 已实现
- **不得** 展示 prompt / completion 全文
- **不得** 伪造 `duration_ms`、`status_code`（schema 当前 **无**）
- list 未就绪时 **禁止** 逐租户 usage 轮询冒充平台 list
- `state=duplicate` 为幂等语义，**非** 调用失败
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 平台跨租户 inference audit list — **ADDED-TO-YAML**
- `duration_ms` / HTTP status / caller 字段 — P1
- 与 [`platform-audit-log.md`](platform-audit-log.md) 事件关联 — P1
- 平台 CSV 导出 — Phase 2
- 异常检测（突增 tokens）— Phase 2

## 响应示例

### reportTokenUsage 成功（Core · **写入 · 非 list**）

```json
{
  "id": "tu-001",
  "tenant_id": "t-001",
  "source": "inference-gateway",
  "model": "llama-3-70b",
  "input_tokens": 512,
  "output_tokens": 128,
  "total_tokens": 640,
  "request_id": "req-inf-001",
  "instance_id": "inst-abc",
  "state": "accepted",
  "created_at": "2026-06-16T10:05:00Z"
}
```

### 平台 inference audit list 目标（**待 YAML**）

```json
{
  "items": [
    {
      "id": "tu-001",
      "tenant_id": "t-001",
      "source": "inference-gateway",
      "model": "llama-3-70b",
      "input_tokens": 512,
      "output_tokens": 128,
      "total_tokens": 640,
      "request_id": "req-inf-001",
      "state": "accepted",
      "created_at": "2026-06-16T10:05:00Z"
    }
  ],
  "next_cursor": "cursor-inf-audit-2",
  "total": 120034
}
```

### GET /metering/usage 引用（租户聚合 · **非逐事件**）

```json
{
  "items": [
    {
      "resource_type": "tokens",
      "total_quantity": 1500000,
      "unit": "tokens",
      "period": "2026-06-01/2026-06-16"
    }
  ],
  "total": 1
}
```

## 错误示例

### 时间窗非法（list · **待 YAML**）

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be before end_time",
  "request_id": "req-boss-inf-audit-400-001"
}
```

### 无平台 inference audit 读权限（**TODO-YAML**）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-inf-audit-403-001"
}
```

## 相关模块

- [`platform-audit-log.md`](platform-audit-log.md)
- [`platform-apikey-audit.md`](platform-apikey-audit.md)
- [`platform-trace.md`](../health/platform-trace.md)
- [`inference-monitoring.md`](../health/inference-monitoring.md)（若存在）
- [`compliance-export.md`](compliance-export.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] `TokenUsageReport` 字段与 v1.yaml 一致；state enum accepted/duplicate
- [x] reportTokenUsage=写入、usage=租户聚合只读参考已区分
- [x] 禁止 prompt 全文；无 duration_ms/status_code 已标注
- [x] 含字段展示规则、口径、状态矩阵、400+403 错误示例
- [x] 删除前置校验 N/A
- [ ] platform inference audit YAML 合入后回写冻结表
- [x] PRD/SPEC/HTML 与本文同步
