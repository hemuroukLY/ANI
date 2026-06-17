# GPU 监控

## 页面定位

`GPU 监控` 是 BOSS **运维与可观测** 域下的 **平台级 GPU 运行态** 专页：全平台 GPU 总量、分配率、平均利用率、异常设备数、型号/节点维度分布与租户占用 Top N，供 SRE 容量规划与故障排查。

本页属于 **Core / GPU Inventory + Observability** 视角下的 **平台 RBAC** 页面，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`。

Console [`gpu-management.md`](../../console-modules/compute/gpu-management.md) 仅 **当前 JWT 租户** 可见范围内的 GPU inventory；本页须 **平台 RBAC + 跨租户 aggregate**，**不得** 逐租户切换 JWT 轮询冒充平台大盘。用量结算深链见 [`../metering/platform-gpu-hours.md`](../metering/platform-gpu-hours.md)。

## 文档管理规则

- 本文是 `GPU 监控` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-gpu-monitoring.md`](../../tasks/modules/prd/boss/health/prd-boss-gpu-monitoring.md) 与 [`spec-boss-gpu-monitoring.md`](../../tasks/modules/spec/boss/health/spec-boss-gpu-monitoring.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- GPU inventory 归属 **Core**；路径 `/api/v1/gpu-inventory*`（`listGPUInventory`、`getGPUOccupancy`）
- PromQL 代理归属 **Core** `GET /api/v1/observability/query` — **不等于** GPU 大盘聚合 API
- 租户级 `gpu-inventory` 为 **JWT 可见范围** 只读参考；BOSS 跨租户排行/趋势 — **TODO-YAML**
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id` 作为越权筛选依据
- 不要求 BOSS 前端显式传 `X-Tenant-Id` 做平台大盘
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）；本页查询侧暂无 422
- 禁止自造 schema / operationId / 路径 / 返回码
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **全平台** GPU 资源总览（总量 / 已分配 / 可用 / 平均利用率 / 异常设备数）
- 提供按型号、节点维度的分布图表与设备 list（平台 aggregate 待 YAML）
- 提供租户 GPU 占用 Top N 排行（平台 aggregate 待 YAML）
- 支持跳转 GPU 资源池（[`../ops/gpu-pool-management.md`](../ops/gpu-pool-management.md)）、告警规则（[`alert-rules.md`](alert-rules.md)）、用量（[`../metering/platform-gpu-hours.md`](../metering/platform-gpu-hours.md)）
- **不承担** GPU drain/分配/回收等写操作（Phase 2）

## 页面结构

- 首屏至少包含：`时间筛选区`（利用率趋势）、`平台 KPI 区`、`型号/节点分布`、`租户 Top N 表`、`边界说明`
- KPI 与排行 **共享同一查询上下文**
- 无数据态、无权限态、平台 API 未就绪态、PromQL 失败态须可区分

```text
GPU 监控
├── 时间筛选（start_time / end_time — 利用率趋势）
├── 可选筛选（vendor / model / state / 租户 — 平台 API 待 YAML）
├── KPI：gpu_total / allocated / available / avg_utilization / abnormal_count
├── 型号与节点维度图表
├── 设备 list（GPUInventoryDevice — 平台 list 待 YAML）
├── 租户占用 Top N
└── 跳转 → gpu-pool-management / alert-rules / platform-gpu-hours
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/gpu-inventory` | 租户上下文 **只读参考**；`GPUInventoryDevice` schema |
| Core | `GET /api/v1/gpu-inventory/occupancy` | 租户上下文占用摘要；`GPUOccupancySummary` |
| Core | `GET /api/v1/observability/query` | PromQL 利用率/温度等时序 |
| Core | 平台 GPU aggregate **TODO-YAML** | BOSS 正式跨租户 KPI/排行数据源 |
| Core | `GET /api/v1/metering/usage` | 用量深链参考；`instance_gpu_seconds` — 见 metering 专页 |

### 关键边界

- `listGPUInventory` 返回 **当前 JWT 租户可见** 设备，**不能** 直接作为 BOSS 跨租户正式契约
- `getGPUOccupancy` 为 **单租户上下文** 摘要，**不能** 直接作为平台 KPI
- `observability/query` 为 **PromQL 代理**，**不得** 替代 GPU inventory list
- GPU 分配/回收/drain 写操作 **不在本页**
- 「GPU-Hours」用量口径见 metering 域；本页聚焦 **运行态/占用**，不重复写结算金额

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 时间筛选区 | UI + PromQL | 利用率趋势窗口 | 刷新图表 |
| 平台 KPI | 平台 aggregate **待 YAML** | `GPUOccupancySummary` 口径参考 | — |
| 型号/节点分布 | aggregate 或租户 list 参考 | `vendor`/`model`/`node_name` | 设备 list |
| 设备 list | 平台 list **待 YAML** | `GPUInventoryDevice` 全字段 | 实例钻取 |
| 利用率趋势 | PromQL + curated **待 YAML** | `utilization_percent` 或 metrics | — |
| 租户 Top N | 平台 aggregate **待 YAML** | 跨租户占用排行 | tenant / gpu-hours |
| 边界说明 | 规划项 | 租户 path 仅为参考 | platform-gpu-hours |

## BOSS 与 Console 分工

| 维度 | BOSS GPU 监控 | Console GPU 算力管理 |
|---|---|---|
| 范围 | 全平台多租户热点 | 单租户 inventory |
| List API | 平台 aggregate **待 YAML** | `listGPUInventory` |
| 占用摘要 | 平台 aggregate **待 YAML** | `getGPUOccupancy` |
| 指标 | PromQL + curated **待 YAML** | 同左（租户 scope） |
| 用量结算 | 跳转 platform-gpu-hours | usage-report |
| 分配/回收 | Phase 2 | 待 Core 冻结 |
| RBAC | 平台 gpu-inventory + observability 读 | 租户 JWT |

## 当前冻结事实

| 方法 | 路径 | operationId | RBAC | 说明 |
|---|---|---|---|---|
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | `scope:gpu-inventory:read` | 租户可见范围；**非 BOSS 平台 list** |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | `scope:gpu-inventory:read` | 租户占用摘要 |
| GET | `/api/v1/observability/query` | `queryObservability` | `scope:observability:read` | PromQL |

| 字段 | 冻结值 |
|---|---|
| `GPUInventoryDevice.state` | `available` / `allocated` / `unavailable` / `maintenance` |
| `GPUOccupancySummary` | `total` / `allocated` / `available` / `avg_utilization_pct` |

| 能力 | 状态 |
|---|---|
| 跨租户 GPU KPI / 排行 / 设备 list | **TODO-YAML** |
| 平台 curated PromQL（GPU 利用率） | **TODO-YAML** |

## 字段级定义

### 查询字段 — `listGPUInventory`（租户参考 · YAML 已冻结）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `vendor` | query | ❌ | 厂商筛选 |
| `state` | query | ❌ | `available` / `allocated` / `unavailable` / `maintenance` |
| `limit` | query | ❌ | 默认 20；1–100 |
| `cursor` | query | ❌ | 分页游标 |

### 查询字段（平台 aggregate 目标 · 待 YAML）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `start_time` | query | ✅ | 利用率趋势起点 |
| `end_time` | query | ✅ | 利用率趋势终点 |
| `tenant_id` | query | ❌ | 平台 RBAC 筛选；**非** 越权依据 |
| `model` | query | ❌ | 型号筛选 |

### 返回字段 — `GPUInventoryDevice`（YAML 已冻结）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | string | 设备 ID |
| `vendor` | string | 厂商 |
| `model` | string | 型号 |
| `state` | enum | 四态 |
| `node_name` | string, nullable | 节点 |
| `memory_gib` | integer, nullable | 显存 GiB |
| `utilization_percent` | float 0–100, nullable | 利用率 |
| `allocated_to` | string, nullable | 占用实例 ID |

### 返回字段 — `GPUOccupancySummary`（YAML 已冻结 · 租户上下文）

| 字段 | 类型 | 说明 |
|---|---|---|
| `total` | integer ≥ 0 | GPU 总数 |
| `allocated` | integer ≥ 0 | 已分配 |
| `available` | integer ≥ 0 | 可用 |
| `avg_utilization_pct` | float, nullable | 平均利用率 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `abnormal_count` | `state=unavailable` 或 `maintenance` 计数 |
| `allocation_rate_pct` | `allocated / total × 100` |
| `tenant_rank` | Top N 序号 |
| `share_pct` | 租户占平台 GPU 百分比 |
| `by_model[]` / `by_node[]` | 维度拆分 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | KPI + 分布 + Top N + 趋势 | 共享查询上下文 |
| 无数据态 | 「当前时间范围内暂无 GPU 监控数据」 | **不** 伪造全 0 |
| 平台 API 未就绪 | 说明 aggregate 待 Core 合入 | 不得伪装生产已上线 |
| PromQL 失败 | 趋势区独立失败 + 重试 | 保留 KPI/list 区 |
| 查询失败态 | 各区块分别失败提示 | 保留筛选 |
| 无权限态 | 403 提示 | 不渲染假数据 |
| `unavailable` | 红色高亮 | 计入 abnormal_count |
| `maintenance` | 灰色/维护标签 | 不计入 available |
| `utilization_percent` null | 显示「—」 | 不显示 0% |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `GPUInventoryDevice.state` | YAML 四态 | string |
| `utilization_percent` | 0–100 | 百分比，保留 1 位小数 |
| `memory_gib` | 显存 | GiB 整数 |
| `avg_utilization_pct` | 租户摘要或 aggregate | 百分比 |
| `allocation_rate_pct` | UI 计算 | 百分比 |
| GPU-Hours 用量 | `instance_gpu_seconds` | 见 metering 专页 |

## 状态与能力口径

### GPUInventoryDevice.state

| 状态 | 含义 | 本页展示 |
|---|---|---|
| `available` | 可分配 | 计入 available |
| `allocated` | 已分配 | 计入 allocated |
| `unavailable` | 不可用/故障 | **高亮**；计入 abnormal |
| `maintenance` | 维护中 | 维护标签；不计入 available |

本页 **只读**；不在此页分配/回收 GPU。

| 能力 | 说明 |
|---|---|
| 查看监控 | 只读 |
| 创建告警规则 | 跳转 alert-rules |
| GPU drain/分配 | Phase 2 → maint-skills / ops |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 gpu-inventory 读 RBAC | 已授权 | `403 FORBIDDEN` |
| 平台 observability 读（PromQL） | 已授权 | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| PromQL `query` | 非空 | `400 BAD_REQUEST` |
| `start_time` / `end_time` | 合法且 start < end | `400 BAD_REQUEST` |
| 钻取单租户 | 平台租户查看权 | `403 FORBIDDEN` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 查看监控大盘 | ✅ | ✅ | ✅ | aggregate 待 YAML |
| 切换 vendor/state/时间 | ✅ | ✅ | ✅ | 平台 API 待 YAML |
| PromQL 自定义查询 | ✅ | ✅ | ✅ | queryObservability |
| 跳转 GPU 资源池 | ✅ | ✅ | ✅ | ops |
| 跳转 GPU-Hours 用量 | ✅ | ✅ | ✅ | metering |
| 创建告警规则 | ❌ | ✅ → alert-rules | ✅ | 非本页 API |
| GPU drain/分配 | ❌ | Phase 2 | Phase 2 | maint-skills |
| 导出 CSV | ❌ | Phase 2 | Phase 2 | **待 YAML** |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。

## 接口冻结规则

### `GET /api/v1/gpu-inventory`（租户上下文 · 只读参考 · **非 BOSS 正式契约**）

- `operationId`：`listGPUInventory`
- `summary`：`查询 GPU 设备清单`
- `tags`：`["GPU"]`
- `x-ani-rbac-scope`：`scope:gpu-inventory:read`
- `query`：`vendor`、`state`（四态 enum）、`limit`（1–100，默认 20）、`cursor`
- `success`：`200 + GPUInventoryListResponse`（`items` + `next_cursor`）
- `errors`：`401`、`403`
- **不得** 写成 BOSS 跨租户 list 正式契约

### `GET /api/v1/gpu-inventory/occupancy`（租户上下文 · **非平台 KPI**）

- `operationId`：`getGPUOccupancy`
- `summary`：`查询 GPU 占用分布摘要`
- `success`：`200 + GPUOccupancySummary`
- `errors`：`401`、`403`
- **不得** 直接作为平台 KPI 正式数据源

### `GET /api/v1/observability/query`（Core · PromQL · **非 inventory list**）

- `operationId`：`queryObservability`
- `query.required`：`query`
- `success`：`200 + ObservabilityQueryResponse`
- `errors`：`400`、`401`、`403`
- **不得** 替代 GPU inventory 或平台 aggregate

### 平台 GPU aggregate（待补）

<!-- TODO-YAML: GET /api/v1/gpu-inventory/platform 或等价 -->

- 须支持跨租户 KPI、Top N、设备 list
- 须平台 RBAC 鉴权
- 合入前不得写入「已冻结」正文

## 使用规则

- 页面不得把 `listGPUInventory` 暴露为 BOSS 平台跨租户 list 入口
- 不得把 `getGPUOccupancy` 直接渲染为平台 KPI 而不标注租户 scope 边界
- 不得把 PromQL 结果当作设备 list 渲染
- 平台 aggregate 未上线时，**禁止** 生产环境用逐租户 JWT 轮询作为正式方案
- UI 展示 `GPUInventoryDevice.state` 必须使用 YAML 四态
- 用量结算跳转 metering，不在本页混写 `instance_gpu_seconds` 排行逻辑

## 待补能力边界

- 平台跨租户 GPU aggregate — **ADDED-TO-YAML** `getPlatformGPUMonitoring`
- 平台 curated PromQL（GPU 利用率/温度）— **TODO-YAML**
- 与 [`../metering/platform-gpu-hours.md`](../metering/platform-gpu-hours.md) 口径对齐 — Phase 1
- GPU drain / 维护窗口 — Phase 2 → maint-skills
- 平台 CSV 导出 — Phase 2

## 响应示例

### 租户 occupancy 只读参考（Core · 非 BOSS 正式契约）

```json
{
  "total": 128,
  "allocated": 96,
  "available": 28,
  "avg_utilization_pct": 72.5
}
```

### 租户 inventory 只读参考（Core · 非 BOSS 正式契约）

```json
{
  "items": [
    {
      "id": "gpu-node-01-slot-0",
      "vendor": "nvidia",
      "model": "A100-80GB",
      "node_name": "gpu-node-01",
      "memory_gib": 80,
      "state": "allocated",
      "utilization_percent": 85.2,
      "allocated_to": "inst-550e8400-e29b"
    },
    {
      "id": "gpu-node-02-slot-1",
      "vendor": "nvidia",
      "model": "H100-80GB",
      "node_name": "gpu-node-02",
      "memory_gib": 80,
      "state": "unavailable",
      "utilization_percent": null,
      "allocated_to": null
    }
  ],
  "next_cursor": null
}
```

### PromQL 利用率查询成功（Core · 已冻结）

```json
{
  "query": "avg(gpu_utilization_percent)",
  "result_type": "vector",
  "results": [
    {
      "metric": {},
      "value": 68.3,
      "timestamp": "2026-06-16T10:00:00Z"
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
  "request_id": "req-boss-gpu-400-001"
}
```

### 无 gpu-inventory 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: scope:gpu-inventory:read required",
  "request_id": "req-boss-gpu-403-001"
}
```

## 相关模块

- [`../ops/gpu-pool-management.md`](../ops/gpu-pool-management.md)、[alert-rules.md](alert-rules.md)
- [`../metering/platform-gpu-hours.md`](../metering/platform-gpu-hours.md)
- Console：[gpu-management.md](../../console-modules/compute/gpu-management.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] `GPUInventoryDevice.state` 四态与 `v1.yaml` 一致
- [x] 区分租户只读参考、PromQL 与 BOSS 平台 aggregate
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [ ] 平台 GPU aggregate YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
