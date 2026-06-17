# 运营总览

## 页面定位

`运营总览` 是 BOSS **平台运营总览** 域下的 **首页聚合页**，面向平台管理员、SRE、运维、运营与交付人员，汇总全平台运行态、租户态、资源池态、服务态、风险态与用量态势。

本页是 **BOSS** 页面，不是 Console 租户使用页。一级权威源为 `ANI-main/repo/api/openapi/v1.yaml` 与 `services/v1.yaml`；**当前无** 独立 BOSS 首页聚合 API。

Console 对照：[`platform-overview.md`](../../console-modules/home/platform-overview.md)。Console 为单租户五主题区块；本页为 **平台级七主题区块**，不得把 Console 首页契约直接冒充 BOSS 正式平台 API。

## 文档管理规则

- 本文是 `运营总览` 的主维护源；页面定位、聚合边界、字段口径、跳转规则和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-operations-overview.md`](../../tasks/modules/prd/boss/overview/prd-boss-operations-overview.md) 与 [`spec-boss-operations-overview.md`](../../tasks/modules/spec/boss/overview/spec-boss-operations-overview.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml` / `services/v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)
- 子主题明细页：[`resource-capacity-trend.md`](resource-capacity-trend.md)、[`gpu-pool-trend.md`](gpu-pool-trend.md)、[`inference-ops-trend.md`](inference-ops-trend.md)、[`kb-ops-trend.md`](kb-ops-trend.md)、[`platform-alerts-pending.md`](platform-alerts-pending.md)

## Core 层要求

- 本页是 **BOSS** 平台侧 **UI 聚合** 能力；底层资源仍遵守 Core/Services 分层
- Core 资源类接口必须使用 `/api/v1/*`；Services 业务类接口必须使用 `/api/v1/svc/*`
- 不允许继续使用旧的 `/api/v1/boss/*` 或 `/api/v1/console/*` 路径
- **当前不冻结** 独立的 `运营总览` 聚合 API；可行方式：并发读取各子模块已声明或待补接口，在页面层做轻聚合
- BOSS 平台级读操作须基于 **平台 RBAC** 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id` 作为跨租户筛选依据
- 租户上下文内已声明路径（如 `GET /api/v1/metering/usage`、`GET /api/v1/gpu-inventory`）在 BOSS 文档中仅作 **只读参考**，**非** BOSS 正式平台契约
- 写操作 POST/PUT/PATCH 默认要求 `idempotency_key`（UUID）
- 错误结构统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）
- 禁止自造 schema 名称、路径、operationId、返回码
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **全平台** 七主题摘要总览（平台状态、租户、资源、服务、告警、交付、用量）
- 为各子主题提供「查看详情」深链至 overview 域明细页或 health/tenant/metering 专页
- 展示 **聚合刷新时间**（`last_refreshed_at`），避免用户误把摘要当实时逐秒值
- 支持 **区块级局部失败/空态**，单区块 API 未就绪不得拖垮整页
- **不承担** 各子域 CRUD、规则编辑、批量处置（跳转对应专页）

## 页面结构

- 首屏至少包含：`顶部平台总判断`、`七主题摘要卡片`、`刷新控件`、`边界说明`
- 各主题卡片 **共享同一刷新上下文**（`last_refreshed_at`）

```text
运营总览
├── 顶部总判断（ok / degraded / error — 聚合自平台健康等）
├── 七主题摘要
│   ├── 平台状态 → platform-health
│   ├── 租户概览 → tenant-list
│   ├── 资源概览 → resource-capacity-trend
│   ├── 服务概览 → inference-ops-trend + kb-ops-trend
│   ├── 告警摘要 → platform-alerts-pending
│   ├── 交付状态 → settings 域（待收口）
│   └── 用量趋势 → tenant-usage-billing / metering 专页
└── 快捷操作
    ├── 新建租户 → tenant-list
    ├── 查看告警 → platform-alerts-pending / alert-rules
    └── 查看平台健康 → platform-health
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 / 来源 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/healthz` | 单组件 liveness **参考**；**非** 大盘健康 |
| Core | `GET /api/v1/observability/query` | PromQL 指标摘要 |
| Core | `GET /api/v1/gpu-inventory`、`/gpu-inventory/occupancy` | 租户上下文 **只读参考**；资源摘要 |
| Core | `GET /api/v1/metering/usage` | 租户上下文 **只读参考**；用量趋势 **非** 平台 aggregate |
| Core | `GET /api/v1/tasks/{task_id}` | 单任务 **路径已声明**；**无** list |
| Core | `GET /api/v1/observability/alert-rules` | 告警**规则** CRUD；**非** firing 事件 |
| Services | `GET /api/v1/svc/inference-services` | 租户上下文 **只读参考** |
| Services | `GET /api/v1/svc/knowledge-bases` | 租户上下文 **只读参考** |
| 平台 aggregate | **TODO-YAML** | BOSS 正式跨租户首页数据源 |

### 关键边界

- **ADDED-TO-YAML** `GET /api/v1/platform/overview`
- 不得把租户级 `metering/usage` 轮询冒充平台用量大盘
- 不得把 `alert-rules` 列表长度写成「当前 P0/P1 告警数」
- `GET /api/v1/tasks/{task_id}` **路径已声明**（无 operationId / 无 x-ani-rbac-scope）；**无** `GET /api/v1/tasks` list
- 子主题字段口径以各 **overview 明细页** 与 **health/tenant/metering 专页** 为准；本页只做摘要，不重复展开字段表

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 顶部总判断 | 平台 health aggregate **待 YAML** + 子区块 rollup | `HealthCheck.status` 三态 rollup | platform-health |
| 平台状态 | platform-health 摘要 | 组件 ok/degraded/error 计数 | platform-health |
| 租户概览 | tenants list **待 YAML** | 总数 / active / suspended | tenant-list |
| 资源概览 | gpu-inventory + instances 聚合 **待 YAML** | GPU/CPU/内存/存储摘要 | resource-capacity-trend |
| 服务概览 | inference + KB aggregate **待 YAML** | 推理服务数、KB 数、模型数 | inference-ops-trend、kb-ops-trend |
| 告警摘要 | firing events **待 YAML** + 任务/服务推断 | P0/P1 计数；无 API 时标注待接入 | platform-alerts-pending |
| 交付状态 | settings 域文档 **待收口** | 最近安装/验收/故障 | ani-installer 等 |
| 用量趋势 | metering platform aggregate **待 YAML** | GPU-Hours / Token / KB Queries | metering 专页、tenant-usage-billing |

## BOSS 与 Console 分工

| 维度 | BOSS 运营总览 | Console 平台概览 |
|---|---|---|
| 视角 | 全平台、多租户、资源池 | 当前租户 |
| 主题数 | 七主题（含租户/交付/平台告警） | 五主题 |
| 数据源 | 平台 aggregate + 子模块并发（待 YAML） | 租户上下文并发 |
| 租户 API | 只读参考，**非** 正式契约 | 正式租户契约 |
| 新建租户 | 快捷入口 → tenant-list | 不在首页 |
| RBAC | 平台 RBAC（**TODO-YAML**） | 租户 JWT + scope |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/healthz` | `liveness` | 无鉴权；**非** 大盘 |
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL；`scope:observability:read` |
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | 租户上下文；BOSS **只读参考** |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | 同上 |
| GET | `/api/v1/observability/alert-rules` | `listObservabilityAlertRules` | **规则**；非 firing |
| GET | `/api/v1/tasks/{task_id}` | — | **路径已声明**；无 list |
| GET | `/api/v1/metering/usage` | — | 租户 JWT；**无** x-ani-rbac-scope |
| GET | `/api/v1/svc/inference-services` | `listInferenceServices` | 租户 Services |
| GET | `/api/v1/svc/knowledge-bases` | `listKnowledgeBases` | 租户 Services |

| 能力 | 状态 |
|---|---|
| BOSS 首页独立 aggregate API | **TODO-YAML** |
| 平台 tenants list | **TODO-YAML** |
| 平台 firing 告警 list | **TODO-YAML** |
| 平台 metering 跨租户 aggregate | **TODO-YAML** |
| 平台 tasks list | **TODO-YAML** |

## 字段级定义

### 展示字段（UI 聚合 · 非 API 字段）

| 字段 | 说明 | 来源区块 |
|---|---|---|
| `overall_status` | 顶部 rollup：`ok` / `degraded` / `error` | 平台状态 |
| `platform_component_error_count` | `error` 组件数 | platform-health |
| `tenant_total` | 租户总数 | 租户概览 |
| `tenant_active` | `active` 租户数 | 租户概览 |
| `tenant_suspended` | `suspended` 租户数 | 租户概览 |
| `gpu_total` | 全平台 GPU 设备数 | 资源概览 |
| `gpu_allocated` | 已分配 GPU | 资源概览 |
| `inference_service_count` | 推理服务总数 | 服务概览 |
| `knowledge_base_count` | 知识库总数 | 服务概览 |
| `alert_p0_count` | P0 当前告警数 | 告警摘要 |
| `alert_p1_count` | P1 当前告警数 | 告警摘要 |
| `pending_task_count` | 待处理任务摘要 | 告警摘要 |
| `gpu_hours_trend` | 窗口 GPU-Hours | 用量趋势 |
| `token_trend` | Input+Output Tokens 窗口值 | 用量趋势 |
| `kb_queries_trend` | KB Queries 窗口值 | 用量趋势 |
| `last_refreshed_at` | 本页最后聚合刷新时间 | 全页 |

### 查询字段（UI · 非 API）

| 字段 | 说明 |
|---|---|
| `time_window` | 趋势默认窗口：24h / 7d / 30d |
| `auto_refresh` | 自动刷新间隔（可选） |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 七主题卡片完整渲染 | 共享 `last_refreshed_at` |
| 单区块无数据 | 该卡片「待接入」或「暂无数据」 | **不** 隐藏整块 |
| 单区块 API 失败 | 卡片内失败提示 + 重试 | 其他区块照常 |
| 平台 aggregate 未就绪 | 卡片标注 TODO-YAML | **禁止** 伪造全绿数字 |
| 租户 API 不可用于平台 | 不展示假跨租户值 | 须等平台 API |
| 无权限态 | 403 提示，不渲染假数据 | 平台 RBAC（**TODO-YAML**） |
| P0 告警 > 0 | 红色高亮 | 跳转 platform-alerts-pending |
| 顶部 `degraded` | 黄色横幅 | 与 `error` 区分 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `overall_status` | 对齐 `HealthCheck.status` | `ok` / `degraded` / `error` |
| `gpu_hours_trend` | `instance_gpu_seconds` 窗口聚合 | 见 metering 专页 |
| `token_trend` | input + output tokens 窗口和 | 整数 |
| `kb_queries_trend` | KB query 计量；enum **待 YAML** | 整数 |
| `last_refreshed_at` | 页面层聚合时间 | ISO 8601 |
| 利用率类 | 窗口聚合，非逐秒实时 | % + 时间窗口说明 |

## 状态与能力口径

### 顶部 overall_status

| 状态 | 含义 | 本页展示 |
|---|---|---|
| `ok` | 各主题无 critical 异常 | 绿色 |
| `degraded` | 部分主题异常或 SLO 逼近阈值 | 黄色 |
| `error` | 存在 P0 或关键组件不可用 | 红色 |

本页 **只读** 摘要；写操作跳转专页。

| 能力 | 说明 |
|---|---|
| 查看首页 | 只读 |
| 区块刷新 | 只读 |
| 新建租户 | 跳转 tenant-list（待 tenants API YAML） |
| 导出运营报表 | Phase 2 **待 YAML** |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 BOSS 读 RBAC | 已授权（**TODO-YAML** scope） | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| 子区块 PromQL | `query` 非空（若用） | `400 BAD_REQUEST` |
| 跳转写操作专页 | 对应模块写权限 | `403 FORBIDDEN` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 / 运营 | 说明 |
|---|---|---|---|---|
| 查看运营总览 | ✅ | ✅ | ✅ | 聚合待 YAML |
| 手动/自动刷新 | ✅ | ✅ | ✅ | 区块级 |
| 跳转子主题明细 | ✅ | ✅ | ✅ | overview 域 |
| 新建租户 | ❌ | ❌ | ✅ → tenant-list | tenants API 待 YAML |
| 查看告警/健康 | ✅ | ✅ | ✅ | 深链 |
| 导出运营报表 | ❌ | Phase 2 | Phase 2 | **待 YAML** |
| 在本页编辑规则/配额 | ❌ | ❌ | ❌ | 跳转专页 |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### 本页无独立冻结 API

- **不得** 自造 `GET /api/v1/boss/overview` 或 `OperationsOverviewResponse`
- 引用子模块 operation 时须标注 **租户只读参考** vs **平台 TODO-YAML**

### `GET /api/v1/healthz`（引用 · **非大盘**）

- `operationId`：`liveness`；`success`：`200`；**不得** 写成全平台 healthy

### `GET /api/v1/observability/query`（引用 · PromQL）

- `operationId`：`queryObservability`；`x-ani-rbac-scope`：`scope:observability:read`
- `errors`：`400`、`401`、`403`

### `GET /api/v1/tasks/{task_id}`（引用 · **路径已声明**）

- `getTask` · `scope:tasks:read` · security 已声明
- `GET /api/v1/tasks` list — **ADDED-TO-YAML** `listTasks`

### 平台首页 aggregate（待补）

<!-- ADDED-TO-YAML: GET /api/v1/platform/overview (`getPlatformOverview`) -->

- 须平台 RBAC；合入前不得写入「已冻结」正文

## 使用规则

- 七主题统一在本页展示，不拆成七个独立首页
- 每卡片须提供「查看详情」深链
- 单区块失败不得伪造其他区块数据
- 必须展示 `last_refreshed_at` 与趋势 `time_window`
- 不得把 Console `platform-overview` 的租户 API 直接当 BOSS 平台数据源
- `dev_profile` 仅联调提示，不进 SLA 报表

## 待补能力边界

- BOSS 首页只读 aggregate API — **ADDED-TO-YAML** P1
- 平台 tenants / tasks / firing alerts / metering aggregate — **ADDED-TO-YAML**
- 交付状态与 settings 域详文收口 — 依赖 settings 模块
- 运营报表 CSV 导出 — Phase 2

## 响应示例

### 页面层聚合目标（**待 YAML · 非已冻结**）

```json
{
  "overall_status": "degraded",
  "last_refreshed_at": "2026-06-17T10:00:00Z",
  "time_window": "24h",
  "platform_status": { "error_count": 1, "degraded_count": 2 },
  "tenants": { "total": 128, "active": 120, "suspended": 8 },
  "resources": { "gpu_total": 512, "gpu_allocated": 401 },
  "services": { "inference_count": 340, "kb_count": 89 },
  "alerts": { "p0": 2, "p1": 5, "pending_tasks": 3 },
  "usage": {
    "gpu_hours": 12450.5,
    "tokens": 8900000,
    "kb_queries": 120000
  }
}
```

### healthz 引用成功（**非本页正式响应**）

```json
{
  "status": "ok",
  "version": "v0.8.0"
}
```

## 错误示例

### 无平台 BOSS 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-ov-403-001"
}
```

> 注：403 示例为平台 RBAC 拒绝参考；BOSS 专属 scope 名称 **待 YAML 冻结**，不得自造 scope 写入正文作为已冻结事实。

### PromQL query 为空（子区块指标）

```json
{
  "code": "BAD_REQUEST",
  "message": "query is required and must not be empty",
  "request_id": "req-boss-ov-400-001"
}
```

## 相关模块

- Overview 明细：[`resource-capacity-trend.md`](resource-capacity-trend.md)、[`gpu-pool-trend.md`](gpu-pool-trend.md)、[`inference-ops-trend.md`](inference-ops-trend.md)、[`kb-ops-trend.md`](kb-ops-trend.md)、[`platform-alerts-pending.md`](platform-alerts-pending.md)
- Health：[`platform-health.md`](../health/platform-health.md)、[`alert-rules.md`](../health/alert-rules.md)
- Tenant：[`tenant-list.md`](../tenant/tenant-list.md)、[`tenant-usage-billing.md`](../tenant/tenant-usage-billing.md)
- Console：[`platform-overview.md`](../../console-modules/home/platform-overview.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 明确无独立 BOSS 首页 API；子路径引用标注只读参考 vs TODO-YAML
- [x] `GET /tasks/{task_id}` 为路径已声明；list 为 TODO-YAML
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（展示/查询）
- [ ] 平台 overview aggregate YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
