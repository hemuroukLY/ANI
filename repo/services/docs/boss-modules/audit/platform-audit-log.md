# 平台审计日志

## 页面定位

`平台审计日志` 是 BOSS **安全审计与合规** 域下的 **跨租户平台操作审计** 专页：检索全平台重要操作记录（谁、何时、对何资源、何种动作、结果与 `request_id`），供平台管理员、SRE 与合规人员做安全排查与取证入口。

本页属于 **Core / Audit** 视角下的 **平台 RBAC** 页面。一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`（OpenAPI `info.description` 提及 `audit` Tag 组，但 **paths 段无** `/audit*` 路径）。

Console [`audit-log.md`](../../console-modules/security/audit-log.md) 为 **单租户只读** 审计检索（同样待 Audit list YAML）；本页须 **平台 RBAC + 跨租户 list**，**不得** 逐租户 JWT 轮询或把 `GET /api/v1/observability/query` 当作审计事件数据源。

## 文档管理规则

- 本文是 `平台审计日志` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-platform-audit-log.md`](../../tasks/modules/prd/boss/audit/prd-boss-platform-audit-log.md) 与 [`spec-boss-platform-audit-log.md`](../../tasks/modules/spec/boss/audit/spec-boss-platform-audit-log.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)
- Console 对照：[`audit-log.md`](../../console-modules/security/audit-log.md)

## Core 层要求

- Audit 事件归属 **Core**；平台跨租户 list — **ADDED-TO-YAML**
- `GET /api/v1/observability/query`（`queryObservability`）为 **PromQL 指标代理**，`scope:observability:read` — **不是** 审计 list 替代
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id` 越权筛选
- 租户上下文路径（若未来 Console 合入）遵循 ANI-14：`tenant_id` 从 JWT claims 提取
- 写操作（导出任务触发等）— **TODO-YAML**；合入后须 `idempotency_key`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）
- 禁止自造 audit schema / operationId / 路径为已冻结
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **全平台** 审计事件 list 与多维筛选（list API 待 YAML）
- 展示事件详情：操作者、资源、结果、`request_id`、耗时等只读字段
- 支持按租户、用户、操作类型、资源类型、结果、时间范围筛选（list **ADDED-TO-YAML**）
- 支持 **`request_id` 跳转 Trace**（见 [`platform-trace.md`](../health/platform-trace.md)）
- 支持合规导出任务入口边界（写 API 待 YAML；异步任务见 [`compliance-export.md`](compliance-export.md)）
- 与 Console 审计日志对照，明确 BOSS 跨租户运维差异

## 页面结构

- 首屏至少包含：`筛选区`、`审计事件表格`、`详情抽屉`、`导出区（待 YAML）`、`边界说明`
- 无数据态、无权限态、list API 未就绪态、查询失败态须可区分

```text
平台审计日志
├── 筛选（tenant_id / actor / action / resource_type / result / 时间 — list **ADDED-TO-YAML**）
├── 审计事件表格（时间 / 租户 / 用户 / 操作 / 资源 / 结果 / request_id）
├── 详情抽屉（全字段 + duration_ms）
│   └── 跳转 Trace（request_id → platform-trace）
└── 操作：导出 / 合规包触发（写 API 待 YAML）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | 平台 audit events list **TODO-YAML** | BOSS 正式跨租户 list |
| Core | `GET /api/v1/observability/query` | **禁止** 作为审计数据源；仅 Observability 域参考 |
| Services | 无独立 audit path | 业务侧操作经 Core 审计管道入库（待 YAML 设计） |

### 关键边界

- v1.yaml **无** `/audit*` path — 正文不得写「已实现 audit list」
- `queryObservability` 返回 PromQL 时序/向量 — **与审计事件语义无关**
- 单租户 Console 审计与 BOSS 平台审计 **schema 应对齐**（待 Core 定义 `AuditEvent`）
- **禁止** 生产环境逐租户 JWT 轮询冒充平台 list
- 导出/取证大包生成走异步任务 — 见 [`compliance-export.md`](compliance-export.md) + `GET /api/v1/tasks/{task_id}`

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 筛选区 | UI + list API **待 YAML** | tenant / actor / action / 时间 | 刷新表格 |
| 事件表格 | list API **待 YAML** | 跨租户审计摘要 | 详情抽屉 |
| 详情抽屉 | list item 或 get **待 YAML** | 全字段只读 | Trace |
| 导出区 | 写 API **待 YAML** | 触发合规/CSV 导出 | compliance-export |
| 边界说明 | 规划项 | observability ≠ audit | audit-log（Console） |

## BOSS 与 Console 分工

| 维度 | BOSS 平台审计日志 | Console 审计日志 |
|---|---|---|
| 范围 | 全平台多租户 | 当前租户 |
| List | platform audit list **待 YAML** P1 | tenant audit list **待 YAML** |
| 导出/取证 | 平台合规导出（见 compliance-export） | 只读；完整包归 BOSS |
| Trace 跳转 | `request_id` → platform-trace | 同入口（Trace API 待 YAML） |
| RBAC | 平台 audit 读/写（待 YAML） | 租户 JWT + audit read scope |
| Observability | **不得** 用 PromQL 代理替代 | 同约束 |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL 代理；**非** audit list |
| GET | `/api/v1/audit/events` | `listAuditEvents` | list · **ADDED-TO-YAML** |
| GET | `/api/v1/audit/events/{event_id}` | `getAuditEvent` | 单条 · **ADDED-TO-YAML** | Tag 仅在 description 提及 |

| 能力 | 状态 |
|---|---|
| 平台跨租户 audit events list | **ADDED-TO-YAML** P1 |
| 单条 audit event get | **ADDED-TO-YAML** P1 |
| 平台 audit 导出 / 合规包触发 | **TODO-YAML** P2 |
| `AuditEvent` schema | **ADDED-TO-YAML** |

### 页面目标字段（待 YAML 合入后对齐）

| 字段 | 说明 |
|---|---|
| `timestamp` | 事件发生时间 |
| `tenant_id` | 所属租户 |
| `actor` | 操作主体（用户 / service account） |
| `source_ip` | 来源 IP |
| `action` | 操作类型 |
| `resource` | 资源标识（type + id） |
| `result` | success / failure 等 |
| `request_id` | 关联请求 ID |
| `duration_ms` | 操作耗时 |

## 字段级定义

### 查询字段（list API 目标 · 待 YAML 合入后对齐）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `tenant_id` | query | 可选 | 平台 RBAC 下筛选；**不得** 越权 |
| `actor` | query | 可选 | 用户 ID 或主体标识 |
| `action` | query | 可选 | 操作类型 enum（待 YAML） |
| `resource_type` | query | 可选 | 资源类型 |
| `result` | query | 可选 | 结果筛选 |
| `start_time` | query | 可选 | 时间下限 |
| `end_time` | query | 可选 | 时间上限 |
| `request_id` | query | 可选 | 精确关联 Trace |
| `limit` | query | 可选 | 分页条数 |
| `cursor` | query | 可选 | 游标分页 |

### 返回字段（AuditEvent · **待 YAML**）

| 字段 | 来源 | 说明 |
|---|---|---|
| `id` | `AuditEvent` | 事件 UUID |
| `timestamp` | `AuditEvent` | ISO 8601 |
| `tenant_id` | `AuditEvent` | 租户 ID |
| `actor` | `AuditEvent` | 操作主体 |
| `source_ip` | `AuditEvent` | 来源 IP；可空 |
| `action` | `AuditEvent` | 操作类型 |
| `resource_type` | `AuditEvent` | 资源类型 |
| `resource_id` | `AuditEvent` | 资源 ID |
| `result` | `AuditEvent` | 操作结果 |
| `request_id` | `AuditEvent` | 请求 ID；可空 |
| `duration_ms` | `AuditEvent` | 耗时毫秒；可空 |
| `metadata` | `AuditEvent` | 扩展键值；可空 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `actor_display` | 用户名 / SA 展示名 |
| `resource_display` | `resource_type` + `resource_id` 组合 |
| `result_badge` | success 绿 / failure 红 |
| `trace_link` | `request_id` 非空时可跳转 Trace |
| `duration_display` | `duration_ms` → 「123 ms」 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 筛选 + 表格 + 可开详情 | list **ADDED-TO-YAML** |
| 无数据态 | 「当前条件下暂无审计事件」 | **不** 伪造行 |
| list API 未就绪 | 说明「跨租户 audit list 待 Core 合入」 | 不得用 observability 冒充 |
| 无权限态 | 403 提示 | 平台 RBAC |
| `result=failure` | 红色高亮 + 展示 metadata 错误摘要 | — |
| `request_id` 非空 | 可点击跳转 Trace | Trace API 待 YAML 时链到搜索页 |
| 时间范围非法 | start ≥ end 时前端拦截 | 合入后 400 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `timestamp` | 后端原值 | ISO 8601 date-time |
| `duration_ms` | 非负整数 | milliseconds |
| `source_ip` | IPv4/IPv6 字符串 | string |
| `action` / `resource_type` | YAML enum（待合入） | string |
| `result` | YAML enum（待合入） | string |
| `request_id` | 与日志/Trace 一致 | string |

## 状态与能力口径

### AuditEvent.result（页面目标 · 待 YAML enum）

| 值 | 含义 | UI |
|---|---|---|
| `success` | 操作成功 | 绿色 |
| `failure` | 操作失败 | 红色高亮 |
| `denied` | 权限拒绝 | 橙色（若 YAML 声明） |

> 最终 enum 以 Core `AuditEvent` schema 为准；合入前不得写为已冻结。

### 状态 × 操作可用性矩阵

| 操作 \ 数据就绪 | list **ADDED-TO-YAML** | list 已合入 |
|---|---|---|
| 列表筛选 | ⏳ 未就绪态 | ✅ |
| 查看详情 | ⏳ | ✅ |
| 跳转 Trace | ✅（带 request_id） | ✅ |
| 导出 CSV | ⏳ Phase 2 | ⏳ Phase 2 |
| 触发合规包 | ⏳ 见 compliance-export | ⏳ |

⏳ = API 待 YAML 合入。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 audit list 读 RBAC | 已授权（**TODO-YAML** scope 待 Core 定义） | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| list 筛选时间 | start < end（list **ADDED-TO-YAML**） | `400 BAD_REQUEST` |
| 导出/写操作（待 YAML） | 写 RBAC + `idempotency_key` | `403` / `422` 待声明 |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 列表（待 YAML） | ✅ | ✅ | ✅ | platform audit list |
| 查看详情 | ✅ **待 YAML** | ✅ | ✅ | get event |
| 筛选 / 搜索 | ✅ **待 YAML** | ✅ | ✅ | — |
| 跳转 Trace | ✅ | ✅ | ✅ | request_id |
| 导出 CSV | ❌ | Phase 2 | Phase 2 | **待 YAML** |
| 触发合规导出 | ❌ | Phase 2 | Phase 2 | compliance-export |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。审计事件为 append-only 留存；物理删除/保留策略由合规产品定义（**待 YAML**），非本页 DELETE API。

## 接口冻结规则

### `GET /api/v1/observability/query`（Core · **已声明 · 非本页数据源**）

- `operationId`：`queryObservability`
- `summary`：`PromQL 代理查询`
- `tags`：`["Observability"]`
- `x-ani-rbac-scope`：`scope:observability:read`
- `parameters`：`query`（必填）、`time`、`timeout`
- `success`：`200 + ObservabilityQueryResponse`
- `errors`：`400`、`401`、`403`
- **不得** 在本页文档中将其描述为 audit events list

### 平台 audit events list（待补）

<!-- TODO-YAML: GET /api/v1/audit/events 或等价 platform list -->

- 须支持跨租户 list 与 cursor 分页
- 须平台 RBAC 鉴权
- query 预期：`tenant_id`、`actor`、`action`、`resource_type`、`result`、`start_time`、`end_time`、`request_id`、`limit`、`cursor`
- **合入前不得写入「已冻结」正文**

### 平台 audit 导出（待补）

<!-- TODO-YAML: POST /api/v1/audit/exports 或等价 -->

- 须声明 `idempotency_key`
- 异步任务返回 `task_id` — 轮询 `GET /api/v1/tasks/{task_id}`
- **合入前不得写入已实现**

## 使用规则

- **不得** 写 audit list 为已实现（v1.yaml 无 `/audit*`）
- **不得** 用 `queryObservability` 或 PromQL 结果渲染审计表格
- list API 未上线时，**禁止** 逐租户 JWT 轮询冒充平台 list
- `request_id` 跳转 Trace 时携带时间窗口（±15min 产品默认）
- 导出/写操作须待 YAML + `idempotency_key`
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 平台跨租户 audit events list/get — **ADDED-TO-YAML**
- `AuditEvent` schema 与 Console 对齐 — P1
- 平台 CSV/Parquet 导出 — Phase 2
- 与 API Key 审计、推理审计的事件 schema 映射 — Phase 2
- 保留周期 / 归档策略 — Phase 2 产品 + Core

## 响应示例

### 审计事件 list 成功（页面目标 · **待 YAML 合入后对齐**）

```json
{
  "items": [
    {
      "id": "ae-001",
      "timestamp": "2026-06-16T10:15:00Z",
      "tenant_id": "t-001",
      "actor": "user:alice@acme.com",
      "source_ip": "203.0.113.10",
      "action": "inference.deploy",
      "resource_type": "inference_service",
      "resource_id": "550e8400-e29b-41d4-a716-446655440000",
      "result": "success",
      "request_id": "req-inf-deploy-001",
      "duration_ms": 842
    }
  ],
  "next_cursor": "cursor-audit-page-2",
  "total": 4096
}
```

### observability 查询成功（**非 audit · 仅边界对照**）

```json
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": []
  }
}
```

## 错误示例

### 时间范围非法（list · **待 YAML**）

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be before end_time",
  "request_id": "req-boss-audit-400-001"
}
```

### 无平台 audit 读权限（**TODO-YAML**）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-audit-403-001"
}
```

> **注**：适用于 **platform audit list（TODO-YAML）** 场景。`message` 为示例文案，**非** YAML 已冻结 scope 名。

## 相关模块

- Console：[`audit-log.md`](../../console-modules/security/audit-log.md)
- [`platform-trace.md`](../health/platform-trace.md)（`request_id` 跳转）
- [`platform-apikey-audit.md`](platform-apikey-audit.md)、[`platform-inference-audit.md`](platform-inference-audit.md)
- [`compliance-export.md`](compliance-export.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 明确 v1.yaml **无** `/audit*`；observability query **非** audit 替代
- [x] 含字段展示规则、字段口径与单位、状态与能力口径
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [x] 删除前置校验 N/A 已声明
- [ ] audit list YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
