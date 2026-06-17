# 平台健康

## 页面定位

`平台健康` 是 BOSS **运维与可观测** 域下的 **控制面/数据面组件巡检** 入口：汇总网关、Auth、模型/知识库/计量/任务服务、数据库、对象存储、向量存储、消息队列与监控组件的运行态，供 SRE 快速判断平台是否处于 `ok` / `degraded` / `error`。

本页属于 **平台 RBAC** 视角，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`（Core 健康与可观测）。

Console 无对等页；租户侧异常摘要见 [`alerts-pending-items.md`](../../console-modules/alerts/alerts-pending-items.md)。本页 **不得** 把租户摘要页的数据源直接当作平台组件大盘。

## 文档管理规则

- 本文是 `平台健康` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-platform-health.md`](../../tasks/modules/prd/boss/health/prd-boss-platform-health.md) 与 [`spec-boss-platform-health.md`](../../tasks/modules/spec/boss/health/spec-boss-platform-health.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- **Liveness**：`GET /api/v1/healthz`（`liveness`）— **仅进程存活**，`status` 枚举仅 `ok`；**不检查依赖**，**不代表** 平台大盘健康
- **PromQL 代理**：`GET /api/v1/observability/query`（`queryObservability`）— 指标曲线；**不等于** platform-health 聚合 API
- **`HealthCheck` schema** 存在于 v1.yaml（`status`: `ok` / `degraded` / `error`；`checks` 嵌套依赖项），**当前无** `GET /platform/health` 或等价 list 路径 — 不得写成已实现
- BOSS 平台健康大盘 aggregate — **TODO-YAML**
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页查询侧暂无 422
- 禁止自造 schema 名称、路径、operationId、返回码
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **全平台** 组件分组健康总览（顶部 `ok` / `degraded` / `error` 汇总）
- 展示各组件的就绪副本、版本、P99 延迟、错误率（aggregate 或 PromQL，待 YAML）
- 支持从组件行跳转指标（`observability/query`）、日志（[`platform-logs.md`](platform-logs.md)）、Trace（[`platform-trace.md`](platform-trace.md)）
- 为 [`alert-rules.md`](alert-rules.md) 与 [`gpu-monitoring.md`](gpu-monitoring.md) 提供运行态上下文
- **不承担** 服务重启/扩缩容等写操作（Phase 2 → [`maint-skills.md`](maint-skills.md)）

## 页面结构

- 首屏至少包含：`顶部总览区`、`服务分组表格`、`刷新/自动刷新控件`、`边界说明`
- 总览与表格 **共享同一刷新上下文**（last_refreshed_at）
- 无数据态、无权限态、平台 API 未就绪态、单组件 probe 失败态须可区分

```text
平台健康
├── 顶部总览（ok / degraded / error — 对齐 HealthCheck.status）
├── 服务分组表格
│   ├── 网关 / Auth / 模型 / 知识库 / 计量 / 任务
│   └── 数据库 / 对象存储 / 向量存储 / MQ / 监控
├── 行内操作
│   ├── 查看指标 → observability/query
│   ├── 查看日志 → platform-logs
│   └── 查看 Trace → platform-trace
└── 边界说明（healthz ≠ 大盘；aggregate 待 YAML）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 / 来源 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/healthz` | 单组件 **liveness** probe 参考；**非** 大盘 |
| Core | `GET /api/v1/observability/query` | PromQL 指标；P99/error_rate 等 |
| Core | `HealthCheck` schema | 字段口径参考；**无** list path |
| Core | 平台 health aggregate **TODO-YAML** | BOSS 正式组件大盘数据源 |
| 运维侧 | K8s `/readyz` / 服务网格健康 | Phase 1 文档口径；**不得** 写成 Core 已冻结 path |

### 关键边界

- `healthz` 200 + `{ "status": "ok" }` **仅表示进程存活**，**不得** 写成「全平台 healthy」或「所有依赖正常」
- `HealthCheck.status`（`ok` / `degraded` / `error`）用于 **未来** 平台聚合响应；与 `healthz` 的 `status: ok` **不是同一契约**
- `observability/query` 为 **PromQL 代理**，**不得** 替代组件 list 或依赖 checks 结构
- `HealthCheck.checks` 嵌套项 `status` 为 `ok` / `fail`（**不是** `error`）— UI 映射时须区分
- 不得在前端伪造 `checks` 依赖探测结果

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 顶部总览 | 平台 aggregate **待 YAML** | 全平台 `HealthCheck.status`  rollup | — |
| 服务分组表格 | aggregate + PromQL + probe 参考 | 每行一组件；`service_name` 产品枚举 | 行内操作 |
| 就绪副本 | aggregate 或 K8s 元数据 **待 YAML** | `replicas_ready` / desired | — |
| P99 / 错误率 | PromQL 或 aggregate | 预置或 curated 查询 | observability/query |
| 版本 | 部署元数据 **待 YAML** | `version` 字符串 | — |
| 边界说明 | 规划项 | healthz / PromQL / aggregate 分工 | alert-rules |

## BOSS 与 Console 分工

| 维度 | BOSS 平台健康 | Console |
|---|---|---|
| 范围 | 全平台控制面/数据面组件 | 当前租户业务摘要 |
| 数据源 | 平台 aggregate + PromQL（待 YAML） | 各模块推断 + 局部任务 |
| healthz | 单组件 probe 参考 | 不暴露 |
| 操作 | 只读巡检；跳转日志/指标/Trace | 跳转业务模块 |
| RBAC | 平台 observability 读 | 租户 JWT |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/healthz` | `liveness` | 无鉴权；`status` 仅 `ok` |
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL；`scope:observability:read` |

| Schema | 冻结字段 | 说明 |
|---|---|---|
| `HealthCheck` | `status`: `ok` / `degraded` / `error` | **无** 对应 GET path |
| `HealthCheck.checks.*` | `status`: `ok` / `fail`；`latency_ms`；`error` | 嵌套依赖项 |
| healthz 200 body | `status`: `ok` only | **不等于** `HealthCheck` |

| 能力 | 状态 |
|---|---|
| 平台组件 health list / aggregate | **TODO-YAML** |
| 平台 curated PromQL（组件 P99/error_rate） | **TODO-YAML** |
| 与 K8s / 服务网格探针映射 | Phase 2 |

## 字段级定义

### 查询字段（平台 aggregate 目标 · 待 YAML 合入后对齐）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `component_group` | query | ❌ | 筛选分组：gateway / auth / services / data / mq / observability |
| `status` | query | ❌ | 筛选 `ok` / `degraded` / `error` |
| `refresh` | UI | ❌ | 手动/自动刷新间隔 |

### 返回字段（`HealthCheck` 口径 · 平台 aggregate 目标）

| 字段 | Schema 来源 | 说明 |
|---|---|---|
| `status` | `HealthCheck.status` | 组件或 rollup：`ok` / `degraded` / `error` |
| `version` | `HealthCheck.version` | 部署版本 |
| `checks` | `HealthCheck.checks` | 依赖项 map；每项 `status`/`latency_ms`/`error` |
| `service_name` | 页面目标 · 非 YAML | 产品组件名（见服务分组） |
| `replicas_ready` | 页面目标 · 待 YAML | ready/desired 摘要 |
| `p99_latency_ms` | PromQL 或 aggregate | 毫秒 |
| `error_rate` | PromQL 或 aggregate | 0–1 或百分比展示 |
| `last_restart_at` | 部署元数据 · 待 YAML | date-time |

### 服务分组（产品目标 · 非 API enum）

网关、Auth、模型服务、知识库服务、计量服务、任务服务、数据库、对象存储、向量存储、消息队列、监控组件。

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `overall_status` | 顶部 rollup：`ok` / `degraded` / `error` |
| `degraded_count` | `degraded` 组件数 |
| `error_count` | `error` 组件数 |
| `last_refreshed_at` | 本页最后刷新时间 |
| `check_failures[]` | 从 `checks` 提取 `status=fail` 的项 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 总览 + 分组表格完整渲染 | 共享 `last_refreshed_at` |
| 无数据态 | 「平台健康 aggregate 待 Core 合入」 | **不** 伪造全 green 表格 |
| 平台 API 未就绪 | 说明待 YAML；可展示 healthz 参考（带边界提示） | 不得伪装生产已上线 |
| 单组件 probe 失败 | 行内 `error` 标签 + checks 展开 | healthz 200 不等于依赖 ok |
| PromQL 失败 | 指标列独立失败提示 + 重试 | 保留组件存活列 |
| 查询失败态 | 总览/表格分别失败提示 + 重试 | 保留筛选 |
| 无权限态 | 403 提示，不渲染假数据 | 平台 RBAC 拒绝 |
| `degraded` | 黄色警告色 | 与 `error` 区分 |
| `error` | 红色高亮 | 优先展示在表格顶部 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `HealthCheck.status` | YAML enum 三态 | string |
| `checks.*.status` | `ok` / `fail` | string |
| `checks.*.latency_ms` | 依赖探测延迟 | 整数毫秒 |
| `healthz.status` | **仅** `ok` | 与 HealthCheck 不同契约 |
| `p99_latency_ms` | histogram_quantile 或 aggregate | 毫秒 |
| `error_rate` | 错误数 / 总请求数 | UI 可展示为 % |
| `replicas_ready` | ready / desired | `2/3` 格式 |

## 状态与能力口径

### HealthCheck.status（组件 / rollup）

| 状态 | 含义 | 本页展示 |
|---|---|---|
| `ok` | 组件与关键依赖正常 | 绿色 |
| `degraded` | 部分依赖异常或 SLO 逼近阈值 | 黄色 |
| `error` | 组件不可用或关键依赖 fail | 红色高亮 |

### checks.*.status（依赖项）

| 状态 | 含义 |
|---|---|
| `ok` | 依赖探测成功 |
| `fail` | 依赖探测失败；可读 `error` 字段 |

本页 **只读**；不在此页重启/扩缩容服务。

| 能力 | 说明 |
|---|---|
| 查看大盘 | 只读 |
| 跳转指标/日志/Trace | 只读深链 |
| 创建告警规则 | 跳转 [`alert-rules.md`](alert-rules.md) |
| 服务重启 | Phase 2 → maint-skills |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 observability 读 RBAC | 已授权 | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| PromQL `query` | 非空（指标钻取时） | `400 BAD_REQUEST` |
| 跳转日志/Trace | 对应模块读权限 | `403 FORBIDDEN` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 查看健康大盘 | ✅ | ✅ | ✅ | aggregate 待 YAML |
| 手动/自动刷新 | ✅ | ✅ | ✅ | — |
| 跳转指标（PromQL） | ✅ | ✅ | ✅ | `queryObservability` |
| 跳转日志/Trace | ✅ | ✅ | ✅ | 须 RBAC |
| 创建告警规则 | ❌ | ✅ → alert-rules | ✅ | 非本页 API |
| 重启/扩缩容服务 | ❌ | Phase 2 | Phase 2 | maint-skills |
| 导出 CSV | ❌ | Phase 2 | Phase 2 | **待 YAML** |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### `GET /api/v1/healthz`（Core · Liveness · **非大盘健康**）

- `operationId`：`liveness`
- `summary`：`Liveness probe（进程存活检查）`
- `tags`：`["Health"]`
- `security`：`[]`（无鉴权）
- `success`：`200` + `{ status: "ok", version?, build_at? }`
- `errors`：`503`（极少见）
- **不得** 把 200 写成全平台 healthy 或 `HealthCheck` 等价响应

### `GET /api/v1/observability/query`（Core · PromQL 代理 · **非 health list**）

- `operationId`：`queryObservability`
- `x-ani-rbac-scope`：`scope:observability:read`
- `query.required`：`query`（PromQL，`minLength: 1`）
- `query.optional`：`time`、`timeout`（默认 `30s`）
- `success`：`200 + ObservabilityQueryResponse`
- `errors`：`400`、`401`、`403`
- **不得** 替代平台 health aggregate 或 `HealthCheck.checks`

### `HealthCheck` schema（**无 path**）

- `required`：`status`、`checks`
- `status` enum：`ok` / `degraded` / `error`
- **ADDED-TO-YAML** `GET /api/v1/platform/health` (`getPlatformHealth`) — 不得写入「已冻结路径」表

### 平台 health aggregate（待补）

<!-- ADDED-TO-YAML: GET /api/v1/platform/health -->

- 须返回组件 list + rollup `HealthCheck.status`
- 须平台 RBAC 鉴权
- 合入前不得写入「已冻结」正文

## 使用规则

- 页面不得把 `healthz` 200 渲染为「全平台 ok」
- 不得把 PromQL 结果直接伪造为 `HealthCheck.checks`
- 总览与表格必须同参刷新，避免口径漂移
- 平台 aggregate 未上线时，**禁止** 生产环境用逐组件 healthz 轮询冒充正式大盘
- UI 展示 `HealthCheck.status` 必须使用 YAML 三态；checks 内使用 `ok`/`fail`
- `dev_profile` 仅用于联调提示，不进入 SLA 报表

## 待补能力边界

- 平台组件 health list / aggregate API — **ADDED-TO-YAML** `getPlatformHealth`
- 平台 curated PromQL（组件 P99/error_rate）— **TODO-YAML**
- 与 K8s `/readyz`、服务网格健康映射 — Phase 2
- 与 [`maint-skills.md`](maint-skills.md) 自动 remediation 联动 — Phase 2
- 平台 CSV 导出 — Phase 2

## 响应示例

### healthz 成功（Core · 已冻结 · **非大盘**）

```json
{
  "status": "ok",
  "version": "v0.8.0",
  "build_at": "2026-06-15T08:00:00Z"
}
```

### PromQL 查询成功（组件错误率 · Core · 已冻结）

```json
{
  "query": "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m]))",
  "result_type": "vector",
  "results": [
    {
      "metric": { "component": "gateway" },
      "value": 0.002,
      "timestamp": "2026-06-16T10:00:00Z"
    }
  ]
}
```

### 平台 aggregate 目标响应（**待 YAML · 非已冻结**）

```json
{
  "overall_status": "degraded",
  "items": [
    {
      "service_name": "gateway",
      "status": "ok",
      "version": "v0.8.0",
      "checks": {
        "postgres": { "status": "ok", "latency_ms": 12 },
        "redis": { "status": "ok", "latency_ms": 3 }
      },
      "replicas_ready": "3/3",
      "p99_latency_ms": 45,
      "error_rate": 0.001
    },
    {
      "service_name": "model-service",
      "status": "degraded",
      "version": "v0.8.0",
      "checks": {
        "vector-store": { "status": "fail", "latency_ms": null, "error": "connection timeout" }
      },
      "replicas_ready": "2/3"
    }
  ]
}
```

## 错误示例

### PromQL query 为空

```json
{
  "code": "BAD_REQUEST",
  "message": "query is required and must not be empty",
  "request_id": "req-boss-ph-400-001"
}
```

### 无平台 observability 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: scope:observability:read required",
  "request_id": "req-boss-ph-403-001"
}
```

## 相关模块

- [platform-logs.md](platform-logs.md)、[platform-trace.md](platform-trace.md)、[gpu-monitoring.md](gpu-monitoring.md)、[alert-rules.md](alert-rules.md)
- Console：[alerts-pending-items.md](../../console-modules/alerts/alerts-pending-items.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 区分 `healthz`（仅 ok）、`HealthCheck.status`（三态）与 aggregate 待 YAML
- [x] `HealthCheck.checks.*.status` 使用 `ok`/`fail` 与 YAML 一致
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [ ] 平台 health aggregate YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
