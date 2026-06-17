# 平台 KB Queries

## 页面定位

`平台 KB Queries` 是 BOSS **平台计量与结算** 域下的 **跨租户知识库查询次数** 专页：按时间范围展示各租户 KB 查询消耗排行、趋势与占比，供平台运营、财务做 KB 成本分析与计费规划。

本页属于 **Core / Metering** 视角下的 **平台 RBAC** 页面，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`。

Console [`usage-report.md`](../../console-modules/tenant/usage-report.md) 的「KB Queries」仅为 **当前 JWT 租户** 的 UI 预设视角；**当前 `MeteringUsageRecord.resource_type` enum 不含 KB 类**，Console 与本页在 Core 合入前均可能呈现 **规划空态**。

**可观测 vs 计量边界**：运行态 KB 调用监控见 [`../health/kb-monitoring.md`](../health/kb-monitoring.md)（PromQL / 平台 aggregate 待 YAML）；本页聚焦 **计费计量** 视角，**不得** 把 kb-monitoring 的指标直接当作 metering 正式契约。

## 文档管理规则

- 本文是 `平台 KB Queries` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-platform-kb-queries.md`](../../tasks/modules/prd/boss/metering/prd-boss-platform-kb-queries.md) 与 [`spec-boss-platform-kb-queries.md`](../../tasks/modules/spec/boss/metering/spec-boss-platform-kb-queries.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- 计量归属 **Core / Metering**（KB query 类 resource_type **待 YAML 入 enum**）
- `GET /api/v1/metering/usage` 路径已声明，但 `MeteringUsageRecord.resource_type` **当前不含** KB 类
- **禁止** 把 UI 名「KB Queries」硬编码为已冻结 `resource_type`
- 当前权威源 **未** 给 `/metering/usage` 声明 `operationId`，正文不得自造
- BOSS 跨租户 list/aggregate — **TODO-YAML**；不得把单租户 GET 直接写成 BOSS 正式契约
- **Services** `GET /api/v1/svc/knowledge-bases`（`listKnowledgeBases`）为租户 KB list；**不得** 当 BOSS 平台 metering 契约
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页查询侧暂无 422

<!-- TODO-YAML: resource_type 如 kb_query_count -->

## 页面职责

- 提供 **全平台** 时间范围内的 KB 查询次数排行与趋势（Core enum 合入后）
- 提供按 day/hour 时间桶的聚合视图（平台 API 合入后）
- 支持从排行行钻取到 KB 运营趋势、KB 监控
- 在 Core 未合入前，提供 **诚实的规划空态**，不伪造计量数据
- 为 Phase 2 账单/导出预留边界，当前不写结算金额

## 页面结构

- 首屏至少包含：`时间筛选区`、`平台 KPI 区`、`KB 查询趋势图`、`租户排行表`、`边界说明`
- 趋势图与排行表 **共享同一查询上下文**（start/end、group_by）
- 本页锁定 KB 类 `resource_type`（**待 YAML**，目标如 `kb_query_count`）
- **enum 未合入时**：KPI/图表/表格展示规划空态，**不** 伪造全 0 数据

```text
平台 KB Queries
├── 时间筛选（start_time / end_time）
├── 可选筛选（租户状态 — 平台 API 待 YAML）
├── KPI：平台 KB 查询总量（enum 合入后）
├── 趋势图（group_by=day|hour）
├── 租户排行表（tenant_id / quantity / unit / 占比）
└── 行内钻取 → kb-ops-trend / kb-monitoring
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/metering/usage` | 路径已声明；**当前无** KB query resource_type 映射 |
| Core | 平台 metering aggregate **TODO-YAML** | BOSS 正式 list/排行数据源 |
| Core | `GET /api/v1/observability/query` | **可观测参考**；见 kb-monitoring；**非** metering 契约 |
| Services | `GET /api/v1/svc/knowledge-bases` | 租户 KB list；**非** BOSS metering |

### 关键边界

- 「KB Queries」是 **UI 展示名**；正式 `resource_type` **待 YAML**（目标如 `kb_query_count`）
- [`../health/kb-monitoring.md`](../health/kb-monitoring.md) 为 **可观测** 域：PromQL 指标、解析失败、延迟；**不是** 本页 metering 数据源
- **不得** 把 Services `listKnowledgeBases` 写成 metering 正式 API
- 单租户 `/metering/usage` **不能** 直接作为 BOSS 跨租户正式契约
- 账单金额、发票、对账不在本页（Phase 2）

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 时间筛选区 | UI | `start_time`、`end_time` | 刷新 KPI/图表/表格 |
| 平台 KPI | 平台 aggregate **待 YAML** | 全平台 KB 查询次数汇总 | — |
| KB 查询趋势图 | 平台 aggregate **待 YAML** | `group_by=day\|hour` | 排行表联动 |
| 租户排行表 | 平台 aggregate **待 YAML** | `tenant_id` + `total_quantity` | 租户钻取 |
| 边界说明 | 规划项 | KB enum 待 Core 合入；可观测见 kb-monitoring | kb-ops-trend |

## BOSS 与 Console 分工

| 维度 | BOSS 平台 KB Queries（计量） | Console KB Queries 视角 |
|---|---|---|
| 范围 | 全平台多租户排行 | 单租户 |
| 查询 API | 平台 aggregate **待 YAML** | usage-report UI 预设（当前可能空态） |
| 可观测 | 跳转 kb-monitoring | 知识库域 Services |
| RBAC | 平台 metering 读 | 租户 JWT |
| 钻取 | kb-ops-trend / kb-monitoring | 知识库管理 |
| 账单 | Phase 2 | 不在 usage-report |

## 当前冻结事实

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/metering/usage` | 路径已声明；**当前无** KB query resource_type |

| 字段 | 冻结值 |
|---|---|
| `MeteringUsageRecord.resource_type` | **不含** KB 类 — **TODO-YAML**（目标如 `kb_query_count`） |
| `group_by`（租户 API） | `resource_type` / `az` / `day` / `hour`（KB 合入后适用） |

| 能力 | 状态 |
|---|---|
| KB query resource_type | **ADDED-TO-YAML** `kb_query_count` |
| 跨租户 KB 查询排行 | **TODO-YAML** |
| `group_by=tenant_id` | **TODO-YAML** |
| 平台 CSV 导出 | **TODO-YAML** |

## 字段级定义

### 查询字段（平台 API 目标 · 待 YAML 合入后对齐）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` | query | ✅ | date-time |
| `end_time` | query | ✅ | date-time |
| `resource_type` | query | 固定 | **待 YAML**（目标如 `kb_query_count`） |
| `tenant_id` | query | 可选 | 平台 RBAC 下筛选；**不得** 越权 |
| `group_by` | query | 可选 | 平台扩展 `tenant_id` / `day` / `hour` |
| `tenant_status` | query | 可选 | 产品筛选；**待 YAML** |

### 返回字段

| 字段 | 来源 | 说明 |
|---|---|---|
| `items[].tenant_id` | 平台 aggregate **待 YAML** | 租户标识 |
| `items[].resource_type` | `MeteringUsageRecord` **待 YAML** | 目标如 `kb_query_count` |
| `items[].total_quantity` | `MeteringUsageRecord` **待 YAML** | 聚合查询次数 |
| `items[].unit` | `MeteringUsageRecord` **待 YAML** | 以后端为准；目标 `queries` 或 `count` |
| `items[].period` | `MeteringUsageRecord` | 时间桶，可空 |
| `total` | list response | 条数 |
| `dev_profile` | `CoreDevProfileInfo` | 联调标记；非主展示 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `kb_queries_display` | 直接展示 `total_quantity`（当 unit 匹配时） |
| `share_pct` | 租户占平台总量百分比 |
| `rank` | 排行序号 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态（enum 已合入） | KPI + 趋势图 + 排行表 | 共享同一查询上下文 |
| **规划空态（enum 未合入）** | 「KB 查询计量 resource_type 待 Core 合入，当前暂无数据」 | **不** 伪造全 0 图表或假排行 |
| 无数据态（enum 已合入但无用量） | 「当前时间范围内暂无 KB 查询用量」 | 与规划空态文案区分 |
| 平台 API 未就绪 | 说明「跨租户 aggregate 待 Core 合入」 | 不得伪装为生产已上线 |
| 查询失败态 | KPI/图表/表格分别失败提示 + 重试 | 保留筛选条件 |
| 无权限态 | 403 提示，不渲染假数据 | 平台 RBAC 拒绝 |
| 钻取无租户权 | 行内禁用或 403 提示 | 不得越权查看 tenant 明细 |
| 可观测跳转 | 提供 kb-monitoring 深链 | 标注「运行态监控，非计费计量」 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `resource_type` | **待 YAML**（目标如 `kb_query_count`） | 合入前不得写入已冻结 enum |
| `total_quantity` | 后端聚合原值 | number (double) 或 integer |
| `unit` | 以后端为准；目标 `queries` 或 `count` | string |
| KB Queries（UI） | 与 `total_quantity` 一致 | 整数 |
| `period` | 与 `group_by` 时间桶一致 | string，可空 |
| `start_time` / `end_time` | ISO 8601 | date-time |

## 状态与能力口径

本页为 **只读查询页**，无资源状态机。KB 文档/索引 lifecycle 不在本页变更。

| 能力 | 说明 |
|---|---|
| 查询 | 只读（enum 合入后） |
| 导出 | **待 YAML** |
| KB 运行态监控 | 跳转 [`kb-monitoring.md`](../health/kb-monitoring.md)，非本页 API |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 metering aggregate 读 RBAC | 已授权（scope **TODO-YAML** 待 Core 定义） | `403 FORBIDDEN` |
| 租户 `GET /metering/usage` 只读参考 | JWT 租户上下文；**非** BOSS 跨租户正式方案 | `403`（operation 已声明，**无** `x-ani-rbac-scope`） |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `start_time` / `end_time` | 合法且 start < end | `400 BAD_REQUEST` |
| KB enum 已合入 | Core 扩展完成 | 规划空态（非 HTTP 错误） |
| 钻取单租户 | 具备该租户平台查看权 | `403` |
| 导出（待 YAML） | export RBAC | **待 YAML** |

## 操作可用性矩阵

| 操作 | 平台只读 | 平台运营 | 财务只读 | 说明 |
|---|---|---|---|---|
| 查看跨租户 KB 查询排行 | ✅（可空态） | ✅ | ✅ | aggregate 待 YAML |
| 切换 group_by | ✅ **待 YAML** | ✅ **待 YAML** | ✅ **待 YAML** | enum 合入后 |
| 钻取租户 | ✅ | ✅ | ✅ | 须 RBAC |
| 跳转 kb-monitoring | ✅ | ✅ | ✅ | 可观测，非 metering |
| 导出 CSV | ❌ | ✅ **待 YAML** | ✅ **待 YAML** | — |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### `GET /api/v1/metering/usage`（租户上下文 · 只读参考 · **非 BOSS 正式契约**）

- `operationId`：权威源 **未声明**
- `summary`：`查询租户用量统计`
- `tags`：`["Metering"]`
- `query.required`：`start_time`、`end_time`
- `query.optional`：`resource_type`、`group_by`
- **当前**：`MeteringUsageRecord.resource_type` enum **不含** KB 类；本页 **不得** 固定 `resource_type` 为已冻结值
- `group_by enum`：`resource_type`、`az`、`day`、`hour`（**不含** `tenant_id`）
- `success`：`200 + MeteringUsageResponse`
- `errors`：`400`、`401`、`403`

### 平台 aggregate（待补）

<!-- TODO-YAML: GET /api/v1/metering/usage/platform 或等价 -->

- 须支持 KB query resource_type 与 `group_by=tenant_id`
- 合入前不得写入「已冻结」正文

### 可观测参考（非 metering · 见 kb-monitoring）

- `GET /api/v1/observability/query`（`queryObservability`）— PromQL 指标查询
- **禁止** 把 observability 指标直接当作本页 metering 正式数据源

## 使用规则

- **不得** 自造 `kb_query_count` 为当前已冻结 enum 值
- **不得** 把 Services `listKnowledgeBases` 写成 metering 正式 API
- **不得** 把 kb-monitoring 的 PromQL 指标伪装为 metering 排行数据
- enum 未合入时，页面须展示规划空态，**禁止** 伪造全 0 排行
- 平台 aggregate 未上线时，**禁止** 生产环境用逐租户 JWT 轮询作为正式方案
- 边界说明中须区分 **可观测（kb-monitoring）** 与 **计费计量（本页）**

## 待补能力边界

- `MeteringUsageRecord.kb_query_count` — **ADDED-TO-YAML**
- 平台跨租户 metering API — **ADDED-TO-YAML** P1 (`GET /api/v1/metering/usage/platform`)
- 平台 CSV 导出 — **TODO-YAML**
- 账单金额 / 发票 — Phase 2
- 按 knowledge_base_id 拆分 — Phase 3

## 响应示例

### enum 未合入时的页面空态说明（非 API 响应）

当前 Core 无 KB query resource_type，BOSS 页面应展示：

> KB 查询计量 resource_type 待 Core 合入（目标如 `kb_query_count`），当前暂无跨租户 KB 查询排行数据。运行态监控请见「知识库监控」。

### 单租户 usage 查询成功（KB enum 合入后 · 只读参考 · 非 BOSS 正式契约）

```json
{
  "items": [
    {
      "resource_type": "kb_query_count",
      "total_quantity": 42000,
      "unit": "queries",
      "period": "day"
    }
  ],
  "total": 1,
  "dev_profile": {
    "mode": "real",
    "provider": "kubernetes_rest",
    "real_provider": true,
    "reason": null
  }
}
```

> 注：`kb_query_count` 为 **TODO-YAML** 目标值，合入前不得作为已冻结 enum 引用。

## 错误示例

### 时间范围非法

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be earlier than end_time",
  "request_id": "req-boss-kb-400-001"
}
```

### 无平台 metering 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-kb-403-001"
}
```

> **注**：适用于 **平台 metering aggregate（TODO-YAML）** 场景。`message` 为示例文案，**非** YAML 已冻结 `x-ani-rbac-scope` 名。租户 `GET /metering/usage` 的 `403` 为 operation 已声明返回码，**无** scope 字段。

## 相关模块

- [`../overview/kb-ops-trend.md`](../overview/kb-ops-trend.md)、[`../health/kb-monitoring.md`](../health/kb-monitoring.md)
- [`tenant-usage-billing.md`](../tenant/tenant-usage-billing.md)
- Console：[usage-report.md](../../console-modules/tenant/usage-report.md)

## 回填验收标准

- [x] 满配章节齐全（对照 `boss-full-depth-checklist.md`）
- [x] 明确 KB 不在当前 enum
- [x] 与 kb-monitoring 可观测边界清晰
- [x] 规划空态与无数据态区分清晰
- [x] 含响应示例与错误示例
- [ ] KB resource_type YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
