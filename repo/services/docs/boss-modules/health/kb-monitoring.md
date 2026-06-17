# 知识库监控

## 页面定位

`知识库监控` 是 BOSS **平台可观测与健康** 域下的 **跨租户知识库运行态** 专页：全平台 KB 查询量、解析/索引失败数、延迟分布与租户活跃排行，供 SRE 与平台运营排查 KB 链路问题。

本页属于 **平台 RBAC** 视角，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`（Core 可观测）与 `ANI-main/repo/api/openapi/services/v1.yaml`（Services KB 资源）。

Console 知识库域为 **单租户** `GET /api/v1/svc/knowledge-bases`（`listKnowledgeBases`）；本页 **不得** 逐租户切换 JWT 轮询冒充平台大盘。

**可观测 vs 计量边界**：KB 查询 **计费计量** 见 [`../metering/platform-kb-queries.md`](../metering/platform-kb-queries.md)；本页聚焦 **运行态监控**（PromQL / 平台 aggregate），**不得** 把 metering 的 `resource_type` 直接当作监控指标契约。

## 文档管理规则

- 本文是 `知识库监控` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-kb-monitoring.md`](../../tasks/modules/prd/boss/health/prd-boss-kb-monitoring.md) 与 [`spec-boss-kb-monitoring.md`](../../tasks/modules/spec/boss/health/spec-boss-kb-monitoring.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml` / `services/v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- **Services** KB 资源：`GET /api/v1/svc/knowledge-bases`（`listKnowledgeBases`）— **租户 JWT 上下文**，BOSS 仅作只读参考
- **Core** 指标代理：`GET /api/v1/observability/query`（`queryObservability`）— **PromQL**，**不是** KB list API
- Core metering **当前无** KB query 类 `resource_type` — 计量见 metering 域 **TODO-YAML**
- BOSS 跨租户 KB 调用/失败 aggregate — **TODO-YAML**
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页查询侧暂无 422
- 禁止自造 KB 监控 list API / schema 为已冻结

## 页面职责

- 提供 **全平台** KB 查询调用趋势与失败数（PromQL 或平台 aggregate 待 YAML）
- 提供解析/索引链路失败计数与 P99 延迟（aggregate 或 PromQL）
- 提供按租户的 KB 活跃排行（平台 aggregate 待 YAML）
- 支持从排行行钻取到单租户 KB（Console 深链或平台 drill-down 待 YAML）
- 与 [`../overview/kb-ops-trend.md`](../overview/kb-ops-trend.md) 运营态势口径对齐；计量排行跳转 [`platform-kb-queries.md`](../metering/platform-kb-queries.md)

## 页面结构

- 首屏至少包含：`时间筛选区`、`平台 KPI 区`、`调用/失败趋势图`、`租户活跃排行表`、`边界说明`
- 趋势图与排行表 **共享同一查询上下文**（start/end、时间桶）
- 无数据态、无权限态、平台 API 未就绪态、查询失败态须可区分

```text
知识库监控
├── 时间筛选（start_time / end_time）
├── 可选筛选（kb_id / tenant_id — 平台 API 待 YAML）
├── KPI：平台查询总量 / 失败数 / 解析 P99
├── 调用趋势图（query_count — aggregate 或 PromQL）
├── 失败趋势图（parse/index failure — aggregate 或 PromQL）
├── 租户活跃排行表（tenant_id / query_count / failure_count）
└── 行内钻取 → kb-ops-trend / platform-kb-queries / Console KB
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Services | `GET /api/v1/svc/knowledge-bases` | 租户上下文 **只读参考**；`KnowledgeBase` schema |
| Core | `GET /api/v1/observability/query` | PromQL 指标代理；**非** KB list |
| Core | 平台 KB metrics aggregate **TODO-YAML** | BOSS 正式跨租户监控数据源 |
| Core / Metering | `GET /api/v1/metering/usage` | **不在本页**；见 metering 域 |

### 关键边界

- `listKnowledgeBases` 返回 **当前 JWT 租户** 可见 KB，**不能** 直接作为 BOSS 跨租户正式契约
- `observability/query` 为 **PromQL 代理**，**不得** 替代 KB list 或平台 aggregate
- **不得** 把 metering `resource_type`（待 YAML）与 PromQL 指标混为同一冻结路径
- KB 重索引、删除等写操作 **不在本页**（Phase 2 → [`maint-skills.md`](maint-skills.md)）

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 时间筛选区 | UI + 平台 API **待 YAML** | `start_time`、`end_time` | 刷新 KPI/图表/表格 |
| 平台 KPI | 平台 aggregate **待 YAML** | 全平台 query/failure 汇总 | — |
| 调用趋势图 | PromQL 或 aggregate | query_count 时间序列 | 排行表联动 |
| 失败趋势图 | PromQL 或 aggregate | parse/index failure | alert-rules |
| 租户排行表 | 平台 aggregate **待 YAML** | `tenant_id` + 活跃/失败计数 | 租户钻取 |
| 边界说明 | 规划项 | metering 与 observability 分离 | platform-kb-queries |

## BOSS 与 Console 分工

| 维度 | BOSS KB 监控 | Console 知识库 |
|---|---|---|
| 范围 | 全平台 KB 调用/失败趋势 | 单租户 KB CRUD/文档 |
| List API | 平台 aggregate **待 YAML** | `listKnowledgeBases` |
| 指标 | PromQL + 平台 curated **待 YAML** | 文档状态 / AsyncTask |
| 计量排行 | 跳转 platform-kb-queries | usage-report UI 视角 |
| 重索引等写操作 | Phase 2 → maint-skills | Services KB 路径 |
| RBAC | 平台 observability 读 | 租户 JWT |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/svc/knowledge-bases` | `listKnowledgeBases` | 单租户 list；**非** BOSS 平台 list |
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL；`scope:observability:read` |

| 字段 | 冻结值 |
|---|---|
| `KnowledgeBase.status` | `active` / `rebuilding` / `deleted` |
| `KBDocument.status` | `uploaded` / `parsing` / `indexed` / `failed` / `deleted` |

| 能力 | 状态 |
|---|---|
| 平台 KB 调用/失败 aggregate | **TODO-YAML** |
| 跨租户活跃排行 | **TODO-YAML** |
| Core metering KB `resource_type` | **TODO-YAML**（metering 域） |

## 字段级定义

### 查询字段（平台 aggregate 目标 · 待 YAML 合入后对齐）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` | query | ✅ | date-time |
| `end_time` | query | ✅ | date-time |
| `tenant_id` | query | 可选 | 平台 RBAC 下筛选；**不得** 越权 |
| `kb_id` | query | 可选 | 单 KB 钻取；**待 YAML** |
| `metric` | query | 可选 | `query_count` / `parse_failure` / `index_failure` — **待 YAML** |
| `group_by` | query | 可选 | `tenant_id` / `kb_id` / `day` / `hour` — **待 YAML** |

### PromQL 查询字段（Core · 已冻结）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `query` | query | ✅ | PromQL 字符串；`minLength: 1` |
| `time` | query | 可选 | date-time |
| `timeout` | query | 可选 | 默认 `30s` |

### 返回字段（Services · 租户 list 只读参考）

| 字段 | 来源 | 说明 |
|---|---|---|
| `items[].id` | `KnowledgeBase` | KB UUID |
| `items[].name` | `KnowledgeBase` | 知识库名 |
| `items[].status` | `KnowledgeBase` | active/rebuilding/deleted |
| `items[].doc_count` | `KnowledgeBase` | 文档数 |
| `items[].created_at` | `KnowledgeBase` | 创建时间 |

### 返回字段（PromQL · Core 已冻结）

| 字段 | 来源 | 说明 |
|---|---|---|
| `query` | `ObservabilityQueryResponse` | 回显 PromQL |
| `result_type` | `ObservabilityQueryResponse` | vector/matrix/scalar/string |
| `results[].metric` | `ObservabilityQueryResponse` | label 键值对 |
| `results[].value` | `ObservabilityQueryResponse` | 数值 |
| `results[].timestamp` | `ObservabilityQueryResponse` | 可空 |
| `dev_profile` | `CoreDevProfileInfo` | 联调标记 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `query_count_display` | 平台查询总量；aggregate **待 YAML** |
| `failure_count_display` | 解析/索引失败总数 |
| `parse_latency_p99_ms` | 解析 P99 延迟毫秒 |
| `tenant_rank` | 租户活跃排行序号 |
| `active_kb_count` | 活跃 KB 数量 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | KPI + 双趋势图 + 排行表 | 共享同一查询上下文 |
| 无数据态 | 「当前时间范围内暂无 KB 监控数据」 | **不** 伪造全 0 图表 |
| 平台 API 未就绪 | 说明「跨租户 aggregate 待 Core 合入」 | 不得伪装为生产已上线 |
| PromQL 失败 | 指标区独立失败提示 | list/aggregate 区可独立展示 |
| 无权限态 | 403 提示，不渲染假数据 | 平台 RBAC 拒绝 |
| 钻取无租户权 | 行内禁用或 403 提示 | 不得越权 |
| metering 未合入 | 排行旁提示「计费计量见 platform-kb-queries」 | 不混口径 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `query_count_display` | 后端聚合或 PromQL rate | 次 / 时间桶 |
| `failure_count_display` | parse + index 失败合计 | 整数 |
| `parse_latency_p99_ms` | histogram_quantile 或 aggregate | 毫秒 |
| `KnowledgeBase.status` | YAML enum | string |
| `KBDocument.status` | YAML enum（文档级钻取参考） | string |
| `start_time` / `end_time` | ISO 8601 | date-time |

## 状态与能力口径

### KnowledgeBase.status（只读展示）

| 状态 | 含义 | 本页展示 |
|---|---|---|
| `active` | 正常可用 | 计入活跃 KB |
| `rebuilding` | 重建索引中 | 计入过渡态；可能伴随 AsyncTask |
| `deleted` | 已删除 | 默认不纳入活跃统计 |

### KBDocument.status（文档级钻取参考）

| 状态 | 含义 |
|---|---|
| `uploaded` | 已上传 |
| `parsing` | 解析中 |
| `indexed` | 已索引 |
| `failed` | 解析/索引失败 — **高亮** |
| `deleted` | 已删除 |

本页 **只读**；重索引等写操作不在本页。

| 能力 | 说明 |
|---|---|
| 查询 | 只读 |
| 创建告警规则 | 跳转 alert-rules |
| 重索引 | Phase 2 → maint-skills |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 observability 读 RBAC | 已授权 | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `start_time` / `end_time` | 合法且 start < end | `400 BAD_REQUEST` |
| PromQL `query` | 非空 | `400 BAD_REQUEST` |
| 钻取单租户 KB | 具备该租户平台查看权 | `403 FORBIDDEN` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 查看趋势与排行 | ✅ | ✅ | ✅ | aggregate 待 YAML |
| 切换时间筛选 | ✅ | ✅ | ✅ | 平台 API 待 YAML |
| PromQL 自定义查询 | ✅ | ✅ | ✅ | `queryObservability` |
| 钻取租户/KB | ✅ | ✅ | ✅ | 须 RBAC |
| 跳转 metering 排行 | ✅ | ✅ | ✅ | platform-kb-queries |
| 创建告警规则 | ❌ | ✅ → alert-rules | ✅ | 非本页 API |
| 重索引等写操作 | ❌ | Phase 2 | Phase 2 | maint-skills |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### `GET /api/v1/svc/knowledge-bases`（租户上下文 · 只读参考 · **非 BOSS 正式契约**）

- `operationId`：`listKnowledgeBases`
- `summary`：`知识库列表`
- `tags`：`["KnowledgeBases"]`
- 路径前缀：`/api/v1/svc/*`
- `query`：**无** 分页/筛选参数（当前 YAML）
- `success`：`200` + `{ items: KnowledgeBase[] }`
- **不得** 写成 BOSS 跨租户 list 正式契约

### `GET /api/v1/observability/query`（Core · PromQL 代理 · **非 KB list**）

- `operationId`：`queryObservability`
- `summary`：`PromQL 代理查询`
- `tags`：`["Observability"]`
- `x-ani-rbac-scope`：`scope:observability:read`
- `query.required`：`query`
- `query.optional`：`time`、`timeout`
- `success`：`200 + ObservabilityQueryResponse`
- `errors`：`400`、`401`、`403`
- **不得** 把本接口描述为 KB list 或平台 aggregate 替代

### 平台 KB metrics aggregate（待补）

<!-- TODO-YAML: GET /api/v1/observability/knowledge-bases/platform 或等价 -->

- 须支持跨租户 `tenant_id` / `kb_id` 维度
- 须与 metering 域 `resource_type` **分离** 定义
- 合入前不得写入「已冻结」正文

## 使用规则

- 不得把 `listKnowledgeBases` 暴露为 BOSS 平台跨租户 list 入口
- 不得把 metering `resource_type` 与 PromQL 指标混为同一 API
- 不得自造 `PlatformKBMetricsResponse` 为已冻结 schema
- 平台 aggregate 未上线时，**禁止** 逐租户 JWT 轮询作为正式方案
- 图表与排行表必须同参刷新

## 待补能力边界

- 平台 KB metrics 聚合 API — **ADDED-TO-YAML** `getPlatformKBMonitoring`
- Core metering KB `resource_type` enum — metering 域 **TODO-YAML**
- 与 [`../overview/kb-ops-trend.md`](../overview/kb-ops-trend.md) 口径对齐
- 平台 CSV 导出 — Phase 2
- 按 embedding_model 维度拆分 — Phase 3

## 响应示例

### PromQL 查询成功（Core · 已冻结）

```json
{
  "query": "sum(rate(kb_query_total[5m]))",
  "result_type": "vector",
  "results": [
    {
      "metric": { "tenant_id": "t-002" },
      "value": 128.0,
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
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "name": "product-docs",
      "status": "active",
      "doc_count": 42,
      "created_at": "2026-06-10T08:00:00Z"
    }
  ]
}
```

## 错误示例

### 时间范围非法

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be earlier than end_time",
  "request_id": "req-boss-kb-400-001"
}
```

### 无平台 observability 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: scope:observability:read required",
  "request_id": "req-boss-kb-403-001"
}
```

## 相关模块

- [`../overview/kb-ops-trend.md`](../overview/kb-ops-trend.md)、[`../metering/platform-kb-queries.md`](../metering/platform-kb-queries.md)
- [alert-rules.md](alert-rules.md)、[maint-skills.md](maint-skills.md)
- Console：知识库域详文（Services `KnowledgeBases`）

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] `KnowledgeBase.status` / `KBDocument.status` 与 `services/v1.yaml` 一致
- [x] 与 metering platform-kb-queries 边界分离
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [ ] 平台 aggregate YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
