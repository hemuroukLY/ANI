# Trace

## 页面定位

`Trace` 是 BOSS **平台可观测与健康** 域下的 **分布式调用链排查** 专页：按 `request_id` / `trace_id` 查看跨服务 span 瀑布图、耗时分布与错误点，供 SRE 做全链路故障定位。

本页属于 **平台 RBAC** 视角。当前 `ANI-main/repo/api/openapi/v1.yaml` **无** trace list/search path — 全文数据源 **TODO-YAML** 或外部 OTel/Jaeger（产品待选型）。

Console **无对等页**。主入口来自 [`platform-logs.md`](platform-logs.md) 日志行的 **`request_id`** 跳转；告警详情亦可携带 `request_id` 进入（见 [`alert-rules.md`](alert-rules.md)）。

## 文档管理规则

- 本文是 `Trace` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-platform-trace.md`](../../tasks/modules/prd/boss/health/prd-boss-platform-trace.md) 与 [`spec-boss-platform-trace.md`](../../tasks/modules/spec/boss/health/spec-boss-platform-trace.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- 当前 v1.yaml **无** trace query/search path — 本页正式契约 **TODO-YAML**
- `GET /api/v1/observability/query` 为 **PromQL 指标代理**，**不是** Trace API
- 平台 RBAC 鉴权；跨租户 trace 须后端鉴权，**不得** 信任未授权 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页暂无 422
- span 字段在 YAML 合入前 **仅作页面目标**，不得写成 Core 已冻结 schema
- 禁止自造 trace schema / 路径 / operationId 为已冻结

<!-- TODO-YAML: GET /api/v1/observability/traces/{trace_id} 或 search by request_id -->

## 页面职责

- 提供 **`request_id` / `trace_id`** 搜索入口（API 合入后）
- 提供 Trace 瀑布图与服务调用链可视化
- 提供 span 耗时排行与错误 span 高亮
- 支持从错误 span 跳转回 [`platform-logs.md`](platform-logs.md)（携带 `request_id` + 时间窗口）
- 在 trace API 未合入前，提供 **诚实的未就绪态**，不伪造 span 数据

## 页面结构

- 首屏至少包含：`搜索区`、`Trace 瀑布图`、`服务调用链`、`耗时排行`、`边界说明`
- 无数据态、无权限态、API 未就绪态、404 态须可区分

```text
Trace
├── 搜索：request_id / trace_id（至少一项）
├── Trace 元信息（trace_id / 总耗时 / span 数 / tenant_id）
├── 瀑布图（spans 层级）
├── 服务调用链（按 service 聚合）
├── 耗时排行（Top spans）
├── 错误 span 高亮 → 跳转日志
└── 边界说明（PromQL ≠ Trace；span 字段待 YAML）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | trace query/search **TODO-YAML** | BOSS 正式 Trace 数据源 |
| Core | `GET /api/v1/observability/query` | **仅** PromQL；**禁止** 当 Trace API |
| 外部 | OTel/Jaeger 等 | 架构待决 |

### 关键边界

- **`observability/query` 是 PromQL，不是 Trace**
- 入口 **`request_id`** 来自日志页或错误响应 `request_id` 字段 — 须与日志口径对齐
- span 字段（`service`、`duration_ms`、`status` 等）YAML 合入前 **仅页面目标**
- **不得** 把 Jaeger/OTel 内部 schema 写成 Core 已冻结

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 搜索区 | UI + trace API **ADDED-TO-YAML** | request_id / trace_id | 触发查询 |
| Trace 元信息 | trace API **ADDED-TO-YAML** | trace 级摘要 | — |
| 瀑布图 | trace API **ADDED-TO-YAML** | spans 层级 | 行选中 |
| 服务调用链 | UI 聚合 | 按 service 分组 | — |
| 耗时排行 | UI 计算 | Top N spans | — |
| 错误 span | trace API **ADDED-TO-YAML** | status=error | platform-logs |
| 边界说明 | 规划项 | span 待 YAML | — |

## BOSS 与 Console 分工

| 维度 | BOSS Trace | Console |
|---|---|---|
| 范围 | 全平台 request 链路 | **无对等页** |
| 数据源 | Trace API **待 YAML** | — |
| 入口 | platform-logs `request_id` | — |
| 关联日志 | platform-logs | — |
| RBAC | 平台 trace 读（待 YAML） | — |

## 当前冻结事实

| 项 | 结论 |
|---|---|
| Trace query/search API | **ADDED-TO-YAML** `listObservabilityTraces`/`getObservabilityTrace` |
| `GET /api/v1/observability/query` | 已声明 — **仅 PromQL** |
| Trace/Span Core schema | **未冻结** |

| 能力 | 状态 |
|---|---|
| 按 request_id 搜索 | **ADDED-TO-YAML** |
| 按 trace_id 直查 | **ADDED-TO-YAML** |
| 导出 Trace | Phase 2 |
| OTel/Jaeger 集成 | Phase 2 |

## 字段级定义

### 查询字段（trace API 目标 · 待 YAML 合入后对齐）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `request_id` | query | 二选一 | 入口 ID；来自日志或错误响应 |
| `trace_id` | query | 二选一 | Trace 标识 |
| `start_time` | query | 可选 | 缩小搜索窗口 |
| `end_time` | query | 可选 | 缩小搜索窗口 |
| `tenant_id` | query | 可选 | 平台 RBAC 下筛选；**不得** 越权 |

### 返回字段（trace API 目标 · 待 YAML 合入后对齐）

| 字段 | 来源 | 说明 |
|---|---|---|
| `trace_id` | Trace **待 YAML** | Trace 标识 |
| `request_id` | Trace **待 YAML** | 入口 request_id |
| `tenant_id` | Trace **待 YAML** | 租户；须 RBAC |
| `total_duration_ms` | Trace **待 YAML** | 总耗时毫秒 |
| `span_count` | Trace **待 YAML** | span 总数 |
| `spans[].span_id` | Span **待 YAML** | span 标识 |
| `spans[].parent_span_id` | Span **待 YAML** | 父 span；根为空 |
| `spans[].service` | Span **待 YAML** | 服务名 |
| `spans[].operation` | Span **待 YAML** | 操作名，可空 |
| `spans[].start_time` | Span **待 YAML** | span 开始时间 |
| `spans[].duration_ms` | Span **待 YAML** | 耗时毫秒 |
| `spans[].status` | Span **待 YAML** | ok / error（页面目标 enum；正式待 YAML） |
| `spans[].attributes` | Span **待 YAML** | 附加属性，可空 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `waterfall_offset_ms` | 相对 trace 起点的横向偏移 |
| `depth_level` | span 层级深度 |
| `is_critical_path` | 是否关键路径（UI 启发式） |
| `error_span_highlight` | error status 高亮 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 搜索 + 瀑布图 + 调用链 | — |
| 无搜索条件 | 提示「请输入 request_id 或 trace_id」 | 不发起空查询 |
| 未找到 | 「未找到对应 Trace」 | 合法 ID 无数据 → `404`（待 YAML） |
| API 未就绪 | 说明「trace API 待 Core 合入」 | 不得伪造 span |
| 无权限态 | 403 提示 | 平台 RBAC 拒绝 |
| error span | 红色高亮 + 跳转日志 | 携带 request_id |
| PromQL 误用 | **禁止** 在本页调用 observability/query | 边界说明常驻 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `duration_ms` / `total_duration_ms` | 后端原值 | 毫秒 |
| `request_id` | 与日志/错误响应一致 | string |
| `trace_id` | W3C trace context 或产品 ID | string |
| `spans[].status` | 页面目标 ok/error；正式 enum 待 YAML | string |
| `start_time` | ISO 8601 | date-time |

## 状态与能力口径

### Span status（页面目标 · 正式 enum 待 YAML）

| 状态 | 含义 | 展示 |
|---|---|---|
| `ok` | 正常完成 | 默认色 |
| `error` | span 内错误 | **高亮**；可跳转日志 |

本页 **只读**；无 span 写操作。

| 能力 | 说明 |
|---|---|
| 搜索/查看 | 只读（待 YAML） |
| 跳转日志 | UI 导航 |
| 导出 | Phase 2 |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 trace 读 RBAC | 已授权 | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `request_id` / `trace_id` | 至少一项非空 | `400 BAD_REQUEST` |
| Trace 存在 | 合法 ID | `404 NOT_FOUND`（待 YAML） |
| 跨租户 trace | 须平台查看权 | `403 FORBIDDEN` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 按 request_id 搜索 | ✅ **待 YAML** | ✅ | ✅ | — |
| 按 trace_id 搜索 | ✅ **待 YAML** | ✅ | ✅ | — |
| 查看瀑布图/调用链 | ✅ **待 YAML** | ✅ | ✅ | — |
| 跳转关联日志 | ✅ | ✅ | ✅ | UI 导航 |
| 导出 Trace | ❌ | Phase 2 | Phase 2 | — |
| 修改 span | ❌ | ❌ | ❌ | 不在产品范围 |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### trace query/search（待补 · **当前未冻结**）

<!-- TODO-YAML: GET /api/v1/observability/traces/{trace_id} 或 GET /api/v1/observability/traces?request_id= -->

- 预期能力：按 `trace_id` 直查；按 `request_id` 搜索
- 须平台 RBAC 鉴权
- 须返回统一 `ErrorResponse`
- span 字段 **合入 YAML 前仅作页面目标**
- **合入前不得写入「已冻结」正文**

### `GET /api/v1/observability/query`（Core · **非 Trace API · 禁止混用**）

- `operationId`：`queryObservability`
- `summary`：`PromQL 代理查询`
- `tags`：`["Observability"]`
- `x-ani-rbac-scope`：`scope:observability:read`
- `query.required`：`query`（PromQL）
- `success`：`200 + ObservabilityQueryResponse`
- `errors`：`400`、`401`、`403`
- **本页不得调用此接口渲染 Trace**

## 使用规则

- **不得** 自造 `TraceSearchResponse` 为当前已冻结事实
- **不得** 把 Jaeger/OTel 内部 schema 写成 Core 已冻结
- span 字段 YAML 合入前 **仅** 在「页面目标」章节描述
- 入口 `request_id` 须与 platform-logs 口径一致
- API 未就绪时 **禁止** 伪造 span 瀑布图
- Console 无对等页 — 不得引用虚构 Console 详文

## 待补能力边界

- Trace query API — **ADDED-TO-YAML** P1 (`GET /api/v1/observability/traces`)
- OpenTelemetry/Jaeger 集成 — Phase 2
- 关键路径分析 — Phase 3
- Trace 导出 — Phase 2

## 响应示例

### Trace 查询成功（页面目标 · **待 YAML 合入后对齐**）

```json
{
  "trace_id": "trace-7f3a9b2c",
  "request_id": "req-abc-123",
  "tenant_id": "t-001",
  "total_duration_ms": 342,
  "span_count": 5,
  "spans": [
    {
      "span_id": "span-root",
      "parent_span_id": null,
      "service": "ani-gateway",
      "operation": "POST /v1/chat",
      "start_time": "2026-06-16T10:05:23.000Z",
      "duration_ms": 342,
      "status": "error",
      "attributes": { "http.status_code": "500" }
    },
    {
      "span_id": "span-gpu",
      "parent_span_id": "span-root",
      "service": "ani-scheduler",
      "operation": "allocate_gpu",
      "start_time": "2026-06-16T10:05:23.050Z",
      "duration_ms": 280,
      "status": "error",
      "attributes": { "reason": "InsufficientGPU" }
    }
  ]
}
```

## 错误示例

### 搜索条件缺失

```json
{
  "code": "BAD_REQUEST",
  "message": "either request_id or trace_id is required",
  "request_id": "req-boss-trace-400-001"
}
```

### 无平台 trace 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: platform trace read required",
  "request_id": "req-boss-trace-403-001"
}
```

## 相关模块

- [platform-logs.md](platform-logs.md)、[alert-rules.md](alert-rules.md)
- [`inference-monitoring.md`](inference-monitoring.md)、[`job-history.md`](job-history.md)
- Console：**无对等页**

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 明确 PromQL ≠ Trace API
- [x] span 字段标注为页面目标（待 YAML）
- [x] 入口 request_id 来自 platform-logs 已声明
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [ ] trace API YAML 合入后回写「已冻结路径」表与 Span schema
- [x] PRD/SPEC/HTML 与本文同步
