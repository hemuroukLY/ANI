# 日志

## 页面定位

`日志` 是 BOSS **平台可观测与健康** 域下的 **跨租户服务日志检索** 专页：按时间、服务、租户、`request_id`、级别与关键词排查控制面与数据面问题，供 SRE 与平台运营做故障定位。

本页属于 **平台 RBAC** 视角。当前 `ANI-main/repo/api/openapi/v1.yaml` **无** logs list/query 路径 — 全文数据源 **TODO-YAML** 或外部日志栈（产品待选型）。

Console **无对等页**；活动告警详情可跳转本页预填 `request_id`（见 [`alert-rules.md`](alert-rules.md)）。链路排查入口见 [`platform-trace.md`](platform-trace.md)。

## 文档管理规则

- 本文是 `日志` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-platform-logs.md`](../../tasks/modules/prd/boss/health/prd-boss-platform-logs.md) 与 [`spec-boss-platform-logs.md`](../../tasks/modules/spec/boss/health/spec-boss-platform-logs.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- 当前 v1.yaml **无** logs list/query path — 本页正式契约 **TODO-YAML**
- `GET /api/v1/observability/query`（`queryObservability`）为 **PromQL 指标代理**，**不是** 日志 API — 正文必须明确区分
- 平台 RBAC 鉴权；跨租户日志须后端鉴权，**不得** 信任 body/query 中未授权的 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页暂无 422
- 禁止自造 `LogEntry` schema / operationId / 路径为已冻结
- OpenAPI 已声明 ≠ handler 已实现

<!-- TODO-YAML: GET /api/v1/observability/logs 或等价 -->

## 页面职责

- 提供 **全平台** 多维度日志检索（时间、服务、租户、`request_id`、level、关键词）
- 提供日志流列表与行展开详情（API 合入后）
- 支持从日志行 **`request_id`** 跳转 [`platform-trace.md`](platform-trace.md)
- 支持从 [`alert-rules.md`](alert-rules.md) 告警详情预填筛选条件
- 在 logs API 未合入前，提供 **诚实的未就绪态**，不伪造日志条目

## 页面结构

- 首屏至少包含：`筛选区`、`日志流列表`、`边界说明`、`跳转 Trace 入口`
- 筛选变更须触发列表刷新；保留筛选条件于失败态
- 无数据态、无权限态、API 已声明（OpenAPI ≠ handler 已实现）态、查询失败态须可区分

```text
日志
├── 筛选：时间 / 服务 / 租户 / request_id / level / 关键词
├── 日志流列表（分页或 cursor — 待 YAML）
│   └── 行展开：完整 message + 结构化字段
├── 行操作：复制 request_id / 跳转 Trace
└── 边界说明（PromQL ≠ 日志；Console 无对等页）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | logs list/query **TODO-YAML** | BOSS 正式日志数据源 |
| Core | `GET /api/v1/observability/query` | **仅** PromQL 指标；**禁止** 当日志 API |
| 外部 | Loki/ES 等 | 架构待决；合入前不得写为已冻结 |

### 关键边界

- **`observability/query` 是 PromQL，不是日志** — 本页 **不得** 用 PromQL 结果渲染日志列表
- Console **无** 对等日志页 — 不得引用不存在的 Console 详文
- Services `getInferenceServiceLogs` 为 **单租户单服务** 日志，**不能** 直接作为 BOSS 平台日志契约
- 跨租户日志检索须 **平台 RBAC** 后端鉴权，前端传入 `tenant_id` 仅作筛选意图

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 筛选区 | UI + logs API **待 YAML** | 时间/服务/tenant/request_id/level | 刷新列表 |
| 日志流列表 | logs API **待 YAML** | 分页或 cursor | 行展开 |
| 行展开 | logs API **待 YAML** | message + 结构化字段 | — |
| 跳转 Trace | UI | 携带 `request_id` | platform-trace |
| 边界说明 | 规划项 | PromQL ≠ 日志 | alert-rules |

## BOSS 与 Console 分工

| 维度 | BOSS 日志 | Console |
|---|---|---|
| 范围 | 全平台多租户日志 | **无对等页** |
| API | logs list **待 YAML** | — |
| 单服务日志 | 平台 logs API 合入后 | Services `getInferenceServiceLogs`（单服务） |
| 关联 Trace | platform-trace | — |
| 关联告警 | alert-rules | — |
| RBAC | 平台 logs 读（待 YAML） | — |

## 当前冻结事实

| 项 | 结论 |
|---|---|
| Logs list/query API | **ADDED-TO-YAML** `listObservabilityLogs` |
| `GET /api/v1/observability/query` | 已声明 — **仅 PromQL 指标** |
| `ObservabilityLogEntry` schema | **ADDED-TO-YAML** |

| 能力 | 状态 |
|---|---|
| 跨租户日志检索 | **ADDED-TO-YAML** |
| 实时 tail | Phase 2 |
| CSV 导出 | Phase 2 |
| Loki/ES 集成 | 架构待决 |

## 字段级定义

### 查询字段（logs API · **ADDED-TO-YAML**）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` | query | ✅ | date-time；检索窗口起点 |
| `end_time` | query | ✅ | date-time；检索窗口终点 |
| `service` | query | 可选 | 服务名筛选 |
| `tenant_id` | query | 可选 | 平台 RBAC 下筛选；**不得** 越权 |
| `request_id` | query | 可选 | 精确或前缀匹配；跳转 Trace 入口 |
| `level` | query | 可选 | debug / info / warn / error（页面目标 enum；正式待 YAML） |
| `keyword` | query | 可选 | 全文检索关键词 |
| `limit` | query | 可选 | 分页条数 |
| `cursor` | query | 可选 | 游标分页 |

### 返回字段（logs API · **ADDED-TO-YAML**）

| 字段 | 来源 | 说明 |
|---|---|---|
| `items[].timestamp` | LogEntry **待 YAML** | 日志时间 |
| `items[].level` | LogEntry **待 YAML** | 级别 |
| `items[].service` | LogEntry **待 YAML** | 服务名 |
| `items[].instance_id` | LogEntry **待 YAML** | 实例/Pod，可空 |
| `items[].tenant_id` | LogEntry **待 YAML** | 租户；须 RBAC |
| `items[].request_id` | LogEntry **待 YAML** | 关联 Trace |
| `items[].message` | LogEntry **待 YAML** | 正文 |
| `items[].fields` | LogEntry **待 YAML** | 结构化附加字段，可空 |
| `next_cursor` | list response | 下一页游标，可空 |
| `total` | list response | 总条数（若 API 提供） |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `level_badge_color` | 级别对应 UI 色 |
| `message_preview` | 列表截断预览 |
| `duration_in_window` | 筛选窗口时长展示 |
| `has_trace_link` | `request_id` 非空时可跳转 Trace |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 筛选 + 日志列表 + 分页 | 保留筛选条件 |
| 无数据态 | 「当前条件下暂无日志」 | **不** 伪造空行 |
| API 已声明（OpenAPI ≠ handler 已实现） | 说明「logs API 待 Core 合入」 | 不得伪装为生产已上线 |
| 查询失败态 | 列表区失败提示 + 重试 | 保留筛选条件 |
| 无权限态 | 403 提示，不渲染假数据 | 平台 RBAC 拒绝 |
| 无效 tenant 筛选 | 403 或空结果 + 说明 | 不得越权 |
| error 级别 | 高亮 / 告警色 | 与 info 区分 |
| PromQL 误用 | **禁止** 在本页调用 `observability/query` 渲染日志 | 边界说明常驻 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `timestamp` | ISO 8601 | date-time |
| `level` | 页面目标 enum；正式待 YAML | debug/info/warn/error |
| `start_time` / `end_time` | ISO 8601 | date-time |
| `request_id` | 与 Trace/错误响应 `request_id` 对齐 | string |
| `message` | 原文展示；列表可截断 | string |

## 状态与能力口径

本页为 **只读检索页**，无资源状态机。日志条目本身无 lifecycle 变更操作。

| 能力 | 说明 |
|---|---|
| 检索 | 只读（待 YAML） |
| 跳转 Trace | UI 导航 |
| 实时 tail | Phase 2 |
| 导出 | Phase 2 |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 logs 读 RBAC | 已授权 | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `start_time` / `end_time` | 合法且 start < end | `400 BAD_REQUEST` |
| 时间窗口上限 | 不超过产品限制（待 YAML） | `400 BAD_REQUEST` |
| 无效 tenant 筛选 | 无平台查看权 | `403 FORBIDDEN` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 检索日志 | ✅ **待 YAML** | ✅ | ✅ | logs API **ADDED-TO-YAML**（OpenAPI ≠ handler） |
| 按 request_id 筛选 | ✅ **待 YAML** | ✅ | ✅ | — |
| 跳转 Trace | ✅ | ✅ | ✅ | UI 导航 |
| 复制 request_id | ✅ | ✅ | ✅ | 客户端 |
| 导出 CSV | ❌ | Phase 2 | Phase 2 | **待 YAML** |
| 实时 tail | ❌ | Phase 2 | Phase 2 | — |
| 删除日志 | ❌ | ❌ | ❌ | 不在产品范围 |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### logs list/query（待补 · **当前未冻结**）

<!-- TODO-YAML: GET /api/v1/observability/logs 或等价 -->

- 预期能力：跨租户日志检索、cursor 分页、`request_id` 筛选
- 须平台 RBAC 鉴权
- 须返回统一 `ErrorResponse`
- **合入前不得写入「已冻结」正文**
- **不得** 自造 `operationId` 为当前事实

### `GET /api/v1/observability/query`（Core · **非日志 API · 禁止混用**）

- `operationId`：`queryObservability`
- `summary`：`PromQL 代理查询`
- `tags`：`["Observability"]`
- `x-ani-rbac-scope`：`scope:observability:read`
- `query.required`：`query`（PromQL）
- `success`：`200 + ObservabilityQueryResponse`
- `errors`：`400`、`401`、`403`
- **本页不得调用此接口渲染日志列表**

## 使用规则

- **不得** 把 PromQL `observability/query` 写成或用作日志查询接口
- **不得** 自造 `LogEntry` Core schema 为已冻结（待 YAML 合入前）
- **不得** 把 Services 单服务 logs 直接写成 BOSS 平台正式契约
- 从 alert-rules 跳转时须预填 `request_id`，不丢失时间窗口
- API 已声明（OpenAPI ≠ handler 已实现）时 **禁止** 伪造日志条目或 mock 生产数据
- Console 无对等页 — 文档不得引用虚构 Console 详文

## 待补能力边界

- 日志 list API — **ADDED-TO-YAML** `listObservabilityLogs`
- Loki/ES 集成选型 — 架构待决
- 实时 tail — Phase 2
- CSV 导出 — Phase 2
- 结构化字段索引 — Phase 3

## 响应示例

### 日志检索成功（页面目标 · **ADDED-TO-YAML**）

```json
{
  "items": [
    {
      "timestamp": "2026-06-16T10:05:23Z",
      "level": "error",
      "service": "ani-core",
      "instance_id": "core-pod-7f3a",
      "tenant_id": "t-001",
      "request_id": "req-abc-123",
      "message": "failed to provision GPU: InsufficientGPU",
      "fields": { "operation_id": "op-xyz-789" }
    }
  ],
  "next_cursor": "cursor-page-2",
  "total": 156
}
```

## 错误示例

### 时间范围非法

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be earlier than end_time",
  "request_id": "req-boss-logs-400-001"
}
```

### 无平台 logs 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: platform logs read required",
  "request_id": "req-boss-logs-403-001"
}
```

## 相关模块

- [platform-trace.md](platform-trace.md)、[alert-rules.md](alert-rules.md)
- [`inference-monitoring.md`](inference-monitoring.md)、[`kb-monitoring.md`](kb-monitoring.md)
- Console：**无对等页**（Services 单服务日志见 inference-observability）

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 明确 PromQL ≠ 日志 API
- [x] Console 无对等页已声明
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [ ] logs API YAML 合入后回写「已冻结路径」表与 LogEntry schema
- [x] PRD/SPEC/HTML 与本文同步
