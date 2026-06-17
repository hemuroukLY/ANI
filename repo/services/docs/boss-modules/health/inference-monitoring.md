# 推理监控

## 页面定位

`推理监控` 是 BOSS **平台可观测与健康** 域下的 **跨租户推理服务运行态** 专页：全平台服务规模（按 `InferenceService.status` 分布）、QPS、延迟、错误率与租户排行，供 SRE 与平台运营排查部署失败与性能热点。

本页属于 **平台 RBAC** 视角，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`（Core 可观测）与 `ANI-main/repo/api/openapi/services/v1.yaml`（Services 推理资源）。

Console [`inference-service.md`](../../console-modules/inference/inference-service.md) 与 [`inference-observability.md`](../../console-modules/inference/inference-observability.md) 为 **单租户** Services 视角；本页 **不得** 逐租户切换 JWT 轮询冒充平台大盘。运营态势入口见 [`../overview/inference-ops-trend.md`](../overview/inference-ops-trend.md)。

## 文档管理规则

- 本文是 `推理监控` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-inference-monitoring.md`](../../tasks/modules/prd/boss/health/prd-boss-inference-monitoring.md) 与 [`spec-boss-inference-monitoring.md`](../../tasks/modules/spec/boss/health/spec-boss-inference-monitoring.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml` / `services/v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- **Services** 推理资源：`GET /api/v1/svc/inference-services`（`listInferenceServices`）— **租户 JWT 上下文**，BOSS 仅作只读参考
- **Core** 指标代理：`GET /api/v1/observability/query`（`queryObservability`）— **PromQL**，**不是** 推理服务 list API
- BOSS 跨租户规模/排行/趋势 — **TODO-YAML**；不得把单租户 Services list 直接写成 BOSS 正式契约
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id` 作为越权筛选依据
- 不要求 BOSS 前端显式传 `X-Tenant-Id` 做平台大盘（租户钻取由平台 API 鉴权）
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页查询侧暂无 422
- 禁止自造 schema 名称、路径、operationId、返回码
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **全平台** 推理服务规模概览（running / failed / deploying 等状态分布）
- 提供 QPS、P99 延迟、错误率趋势（PromQL 或平台 curated 指标，aggregate 待 YAML）
- 提供按租户的服务数量与失败数排行（平台 aggregate 待 YAML）
- 支持从排行行钻取到单租户推理服务（Console 深链或平台 drill-down 待 YAML）
- 为 [`alert-rules.md`](alert-rules.md) 与 [`../metering/platform-input-tokens.md`](../metering/platform-input-tokens.md) 提供运行态深链上下文

## 页面结构

- 首屏至少包含：`时间筛选区`、`平台 KPI 区`、`状态分布`、`QPS/P99/错误率趋势图`、`租户排行表`、`边界说明`
- 趋势图与排行表 **共享同一查询上下文**（start/end、时间桶）
- 无数据态、无权限态、平台 API 未就绪态、查询失败态须可区分

```text
推理监控
├── 时间筛选（start_time / end_time）
├── 可选筛选（status / model / 租户 — 平台 API 待 YAML）
├── KPI：running 数 / failed 数 / deploying 数 / 平台 QPS
├── 状态分布（InferenceService.status 六态）
├── 趋势图（QPS / P99 / error_rate — PromQL 或 aggregate）
├── 租户排行表（tenant_id / service_count / failed_count）
└── 行内钻取 → Console 推理服务 / inference-ops-trend
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Services | `GET /api/v1/svc/inference-services` | 租户上下文 **只读参考**；`InferenceService` schema |
| Core | `GET /api/v1/observability/query` | PromQL 指标代理；**非** 推理 list |
| Core | 平台 inference aggregate **TODO-YAML** | BOSS 正式跨租户规模/排行数据源 |

### 关键边界

- `listInferenceServices` 返回 **当前 JWT 租户** 可见服务，**不能** 直接作为 BOSS 跨租户正式契约
- `observability/query` 为 **PromQL 代理**，**不得** 替代推理服务 list 或平台 aggregate
- 推理变配、部署、删除等写操作 **不在本页**（跳转 Console 或 Services 路径）
- `dev_profile` 为开发/联调辅助字段，**不是** 用户业务展示主字段

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 时间筛选区 | UI + 平台 API **待 YAML** | `start_time`、`end_time` | 刷新 KPI/图表/表格 |
| 平台 KPI | 平台 aggregate **待 YAML** | 全平台服务数、失败数 | — |
| 状态分布 | aggregate 或 PromQL | `InferenceService.status` 六态计数 | — |
| QPS/P99/错误率趋势 | PromQL + 平台 curated **待 YAML** | `queryObservability` | 排行表联动 |
| 租户排行表 | 平台 aggregate **待 YAML** | `tenant_id` + 服务/失败计数 | 租户钻取 |
| 边界说明 | 规划项 | Services list 仅为租户参考 | inference-ops-trend |

## BOSS 与 Console 分工

| 维度 | BOSS 推理监控 | Console 推理域 |
|---|---|---|
| 范围 | 全平台多租户规模与排行 | 单租户 CRUD/运维 |
| List API | 平台 aggregate **待 YAML** | `listInferenceServices` |
| 指标 | PromQL + 平台 curated **待 YAML** | [`inference-observability.md`](../../console-modules/inference/inference-observability.md) |
| 日志 | 跳转 [`platform-logs.md`](platform-logs.md) | `getInferenceServiceLogs` |
| 变配/部署 | 只读监控；写操作不在本页 | PATCH/POST/DELETE 推理服务 |
| RBAC | 平台 observability 读 | 租户 JWT |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/svc/inference-services` | `listInferenceServices` | 单租户 list；**非** BOSS 平台 list |
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL；`scope:observability:read` |

| 字段 | 冻结值 |
|---|---|
| `InferenceService.status` | `pending` / `deploying` / `running` / `stopping` / `stopped` / `failed` |
| `ObservabilityQueryResponse.result_type` | `vector` / `matrix` / `scalar` / `string` |

| 能力 | 状态 |
|---|---|
| 跨租户 inference 规模/失败 aggregate | **TODO-YAML** |
| 跨租户 QPS/P99 排行 | **TODO-YAML** |
| 平台 curated PromQL 模板 | **TODO-YAML** |

## 字段级定义

### 查询字段（平台 aggregate 目标 · 待 YAML 合入后对齐）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` | query | ✅ | date-time；趋势窗口起点 |
| `end_time` | query | ✅ | date-time；趋势窗口终点 |
| `status` | query | 可选 | 筛选 `InferenceService.status` 枚举 |
| `tenant_id` | query | 可选 | 平台 RBAC 下筛选；**不得** 越权 |
| `model` | query | 可选 | 按模型名筛选；**待 YAML** |
| `group_by` | query | 可选 | `tenant_id` / `status` / `day` / `hour` — **待 YAML** |

### PromQL 查询字段（Core · 已冻结）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `query` | query | ✅ | PromQL 字符串；`minLength: 1` |
| `time` | query | 可选 | 单点查询时间；date-time |
| `timeout` | query | 可选 | 默认 `30s` |

### 返回字段（Services · 租户 list 只读参考）

| 字段 | 来源 | 说明 |
|---|---|---|
| `items[].id` | `InferenceService` | 服务 UUID |
| `items[].name` | `InferenceService` | 服务名 |
| `items[].model` | `InferenceService` | 关联模型 |
| `items[].status` | `InferenceService` | 六态枚举 |
| `items[].replicas` | `InferenceService` | 副本数 |
| `items[].gpu_type` | `InferenceService` | GPU 型号，可空 |
| `items[].gpu_count_per_pod` | `InferenceService` | 每 Pod GPU 数 |
| `items[].max_concurrency` | `InferenceService` | 最大并发 |
| `items[].endpoint_url` | `InferenceService` | 端点 URL，可空 |
| `items[].created_at` | `InferenceService` | 创建时间 |

### 返回字段（PromQL · Core 已冻结）

| 字段 | 来源 | 说明 |
|---|---|---|
| `query` | `ObservabilityQueryResponse` | 回显 PromQL |
| `result_type` | `ObservabilityQueryResponse` | vector/matrix/scalar/string |
| `results[].metric` | `ObservabilityQueryResponse` | label 键值对 |
| `results[].value` | `ObservabilityQueryResponse` | 数值 |
| `results[].timestamp` | `ObservabilityQueryResponse` | 时间戳，可空 |
| `dev_profile` | `CoreDevProfileInfo` | 联调标记；非主展示 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `service_count_by_status` | 各 status 计数；aggregate **待 YAML** |
| `qps_display` | 每秒请求数；来自 PromQL 或 aggregate |
| `p99_ms_display` | P99 延迟毫秒；展示层换算 |
| `error_rate_pct` | 错误率百分比 |
| `tenant_rank` | 租户排行序号 |
| `share_pct` | 租户占平台总量百分比 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | KPI + 状态分布 + 趋势图 + 排行表 | 共享同一查询上下文 |
| 无数据态 | 「当前时间范围内暂无推理监控数据」 | **不** 伪造全 0 图表 |
| 平台 API 未就绪 | 说明「跨租户 aggregate 待 Core/Services 合入」 | 不得伪装为生产已上线 |
| PromQL 失败 | 指标区独立失败提示 + 重试 | 保留 list/aggregate 区正常展示 |
| 查询失败态 | KPI/图表/表格分别失败提示 + 重试 | 保留筛选条件 |
| 无权限态 | 403 提示，不渲染假数据 | 平台 RBAC 拒绝 |
| 钻取无租户权 | 行内禁用或 403 提示 | 不得越权查看 tenant 明细 |
| failed 状态 | 高亮 / 告警色 | 与 `running` 区分 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `InferenceService.status` | YAML enum 六态 | string |
| `qps_display` | PromQL rate 或 aggregate | req/s，保留 2 位小数 |
| `p99_ms_display` | histogram_quantile 或 aggregate | 毫秒 |
| `error_rate_pct` | 错误数 / 总请求数 × 100 | 百分比，保留 2 位小数 |
| `start_time` / `end_time` | ISO 8601 | date-time |
| `progress_pct` | 不适用本页（AsyncTask 字段） | — |

## 状态与能力口径

### InferenceService.status 状态机（只读展示）

| 状态 | 含义 | 本页展示 |
|---|---|---|
| `pending` | 已创建，待调度 | 计入待部署 |
| `deploying` | 部署中 | 计入进行中；可能伴随 AsyncTask |
| `running` | 运行中 | 计入健康服务 |
| `stopping` | 停止中 | 计入过渡态 |
| `stopped` | 已停止 | 计入非运行 |
| `failed` | 部署/运行失败 | **高亮**；可钻取 error 上下文 |

本页 **只读**；不在此页 PATCH/POST 变更服务状态。

| 能力 | 说明 |
|---|---|
| 查询 | 只读 |
| 创建告警规则 | 跳转 [`alert-rules.md`](alert-rules.md) |
| 推理变配/重启 | 不在本页 |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 observability 读 RBAC | 已授权 | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `start_time` / `end_time` | 合法且 start < end | `400 BAD_REQUEST` |
| PromQL `query` | 非空 | `400 BAD_REQUEST` |
| 钻取单租户服务 | 具备该租户平台查看权 | `403 FORBIDDEN` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 查看规模与趋势 | ✅ | ✅ | ✅ | aggregate 待 YAML |
| 切换时间/status 筛选 | ✅ | ✅ | ✅ | 平台 API 待 YAML |
| PromQL 自定义查询 | ✅ | ✅ | ✅ | `queryObservability` |
| 钻取租户/服务 | ✅ | ✅ | ✅ | 须 RBAC |
| 创建告警规则 | ❌ | ✅ → alert-rules | ✅ | 非本页 API |
| 推理变配/重启 | ❌ | ❌ 不在本页 | ❌ | Console/Services |
| 导出 CSV | ❌ | Phase 2 | Phase 2 | **待 YAML** |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### `GET /api/v1/svc/inference-services`（租户上下文 · 只读参考 · **非 BOSS 正式契约**）

- `operationId`：`listInferenceServices`
- `summary`：`获取推理服务列表`
- `tags`：`["InferenceServices"]`
- 路径前缀：`/api/v1/svc/*`
- `query`：**无** 分页/筛选参数（当前 YAML）
- `success`：`200` + `{ items: InferenceService[] }`
- `errors`：当前 YAML GET 仅声明 `200`
- 认证/授权错误为 Services 网关通用推断，**未** 在本 operation responses 中冻结
- **不得** 写成 BOSS 跨租户 list 正式契约

### `GET /api/v1/observability/query`（Core · PromQL 代理 · **非推理 list**）

- `operationId`：`queryObservability`
- `summary`：`PromQL 代理查询`
- `tags`：`["Observability"]`
- `x-ani-rbac-scope`：`scope:observability:read`
- `query.required`：`query`（PromQL 字符串，`minLength: 1`）
- `query.optional`：`time`（date-time）、`timeout`（默认 `30s`）
- `success`：`200 + ObservabilityQueryResponse`
- `errors`：`400`、`401`、`403`
- **不得** 把本接口描述为推理服务 list 或平台 aggregate 替代

### 平台 inference aggregate（待补）

<!-- TODO-YAML: GET /api/v1/observability/inference/platform 或等价 -->

- 须支持跨租户 `tenant_id` 维度与 `InferenceService.status` 分布
- 须平台 RBAC 鉴权
- 合入前不得写入「已冻结」正文

## 使用规则

- 页面不得把 `listInferenceServices` 暴露为 BOSS 平台跨租户 list 入口
- 不得把 `observability/query` 的 PromQL 结果直接当作服务 list 渲染
- 图表与排行表必须同参刷新，避免口径漂移
- 平台 aggregate 未上线时，**禁止** 生产环境用逐租户 JWT 轮询作为正式方案
- UI 展示 `InferenceService.status` 必须使用 YAML 六态，不得自造状态值
- `dev_profile` 仅用于联调提示，不进入 SLA 报表

## 待补能力边界

- 平台跨租户 inference aggregate — **ADDED-TO-YAML** `getPlatformInferenceMonitoring`
- 平台 curated PromQL 模板（QPS/P99/error_rate）— **TODO-YAML**
- 与 [`../overview/inference-ops-trend.md`](../overview/inference-ops-trend.md) 口径对齐 — Phase 1
- 平台 CSV 导出 — Phase 2
- 按 `model_id` / AZ 维度拆分 — Phase 3

## 响应示例

### PromQL 查询成功（Core · 已冻结）

```json
{
  "query": "sum(rate(inference_requests_total[5m]))",
  "result_type": "vector",
  "results": [
    {
      "metric": { "tenant_id": "t-001" },
      "value": 42.5,
      "timestamp": "2026-06-16T10:00:00Z"
    }
  ],
  "dev_profile": {
    "mode": "real",
    "provider": "prometheus",
    "real_provider": true,
    "reason": null
  }
}
```

### 租户 list 只读参考（Services · 非 BOSS 正式契约）

```json
{
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "llama-70b-svc",
      "model": "llama-3-70b",
      "replicas": 2,
      "gpu_type": "A100-80G",
      "gpu_count_per_pod": 2,
      "max_concurrency": 8,
      "status": "running",
      "endpoint_url": "https://infer.internal/v1/chat",
      "created_at": "2026-06-15T08:00:00Z"
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
  "request_id": "req-boss-infer-400-001"
}
```

### 无平台 observability 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: scope:observability:read required",
  "request_id": "req-boss-infer-403-001"
}
```

## 相关模块

- [`../overview/inference-ops-trend.md`](../overview/inference-ops-trend.md)、[alert-rules.md](alert-rules.md)
- [`../metering/platform-input-tokens.md`](../metering/platform-input-tokens.md)
- Console：[inference-service.md](../../console-modules/inference/inference-service.md)、[inference-observability.md](../../console-modules/inference/inference-observability.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] `InferenceService.status` 六态与 `services/v1.yaml` 一致
- [x] 区分 Services 租户只读参考、PromQL 与 BOSS 平台 aggregate
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [ ] 平台 aggregate YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
