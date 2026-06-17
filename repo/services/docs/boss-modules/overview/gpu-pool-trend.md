# GPU 资源池态势

## 页面定位

`GPU 资源池态势` 是 BOSS **平台运营总览** 域下的 **全平台 GPU 资源池** 趋势明细页，展示 GPU 型号分布、分配率、利用率趋势、异常设备与跨租户占用排行。

Console 对照：[`home-gpu-utilization.md`](../../console-modules/home/home-gpu-utilization.md)（单租户首页区块）。深度专页：[`gpu-monitoring.md`](../health/gpu-monitoring.md)、[`gpu-pool-management.md`](../ops/gpu-pool-management.md)。

## 文档管理规则

- 本文是 `GPU 资源池态势` 的主维护源
- PRD/SPEC 为辅助；冲突以本文 + `v1.yaml` 为准
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- `GET /api/v1/gpu-inventory`（`listGPUInventory`）、`GET /api/v1/gpu-inventory/occupancy`（`getGPUOccupancy`）— YAML **已声明**；BOSS 上下文为租户 **只读参考**
- `GET /api/v1/observability/query` — PromQL；**非** inventory list
- 平台 GPU 跨租户 aggregate — **TODO-YAML**
- RBAC：`scope:gpu-inventory:read`（租户路径）；平台 scope **TODO-YAML**
- 禁止把租户 inventory 直接写成 BOSS 正式契约

## 页面职责

- 展示全平台 GPU 总量、已分配、空闲、异常数与利用率趋势
- 型号分布与 AZ/节点池维度钻取（aggregate 待 YAML）
- 跨租户 GPU 占用 Top N（**TODO-YAML**）
- 跳转 GPU 监控、GPU 资源池管理、平台 GPU-Hours
- 标明统计周期（非逐秒实时）

## 页面结构

```text
GPU 资源池态势
├── 顶部 KPI（总量 / 已分配 / 利用率 / 异常）
├── 型号分布图
├── 利用率趋势（7d / 30d）
├── 跨租户占用排行（待 YAML）
└── 跳转
    ├── GPU 监控
    ├── GPU 资源池管理
    └── 平台 GPU-Hours
```

## 数据来源与分层约束

| 层 | 路径 | 用法 |
|---|---|---|
| Core | `GET /api/v1/gpu-inventory` | 租户 **只读参考**；`GPUInventoryDevice` |
| Core | `GET /api/v1/gpu-inventory/occupancy` | 租户 **只读参考**；`GPUOccupancySummary` |
| Core | `GET /api/v1/observability/query` | 利用率/温度 PromQL |
| Core | 平台 GPU aggregate **TODO-YAML** | BOSS 正式大盘 |

### 关键边界

- `GPUInventoryDevice.state` enum：`available` / `allocated` / `unavailable` / `maintenance`
- PromQL 具体 metric 名称 **待运维模板 YAML**；文档不得写死未声明名称
- 不得用 Services deprecated `/gpu-containers*` 作为 GPU 事实来源

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| KPI 卡片 | aggregate + occupancy 参考 | — |
| 型号分布 | inventory 聚合 | gpu-monitoring |
| 利用率趋势 | PromQL 或 aggregate | observability/query |
| 租户排行 | **TODO-YAML** | tenant-list |
| 异常设备 | `state=unavailable|maintenance` | gpu-monitoring |

## BOSS 与 Console 分工

| 维度 | BOSS | Console home-gpu-utilization |
|---|---|---|
| 范围 | 全平台 GPU 池 | 当前租户可见 GPU |
| 排行 | 跨租户 Top N | 无 |
| 正式 API | 平台 aggregate 待 YAML | 租户 gpu-inventory |

## 当前冻结事实

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | `scope:gpu-inventory:read` |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | `scope:gpu-inventory:read` |
| GET | `/api/v1/observability/query` | `queryObservability` | `scope:observability:read` |

| 能力 | 状态 |
|---|---|
| 平台 GPU aggregate | **ADDED-TO-YAML** `getPlatformGPUTrend` |

## 字段级定义

| 字段 | Schema / 说明 |
|---|---|
| `gpu_total` | 设备总数 |
| `gpu_allocated` | `state=allocated` 或 occupancy |
| `gpu_available` | `state=available` |
| `gpu_unavailable` | `unavailable` + `maintenance` |
| `avg_utilization_pct` | 窗口平均；PromQL 或 aggregate |
| `model_distribution` | `{ model: count }` |
| `tenant_rank[]` | `tenant_id`, `gpu_allocated` — **待 YAML** |
| `last_refreshed_at` | UI |
| `time_window` | UI |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 正常 | KPI + 趋势 + 分布 |
| 监控不可用 | 「GPU 监控暂不可用」+ 保留跳转 |
| aggregate 未就绪 | 标注 TODO-YAML；不伪造排行 |
| 403 | 无权限提示 |
| 异常 GPU > 0 | 高亮 |

## 字段口径与单位

| 字段 | 口径 |
|---|---|
| `avg_utilization_pct` | 0–100 %；须标注窗口 |
| `GPUInventoryDevice.memory_mb` | 整数 MB |
| occupancy 比例 | 见 `GPUOccupancySummary` YAML |

## 状态与能力口径

### GPUInventoryDevice.state

| 状态 | 展示 |
|---|---|
| `available` | 空闲 |
| `allocated` | 已分配 |
| `unavailable` | 异常红标 |
| `maintenance` | 维护灰标 |

本页只读。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台读 RBAC（**TODO-YAML**） | `403` |
| `scope:gpu-inventory:read`（引用租户 API 时） | `403` |
| PromQL 非空 | `400` |
| 未认证 | `401` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 管理员 |
|---|---|---|---|
| 查看态势 | ✅ | ✅ | ✅ |
| 切换窗口 | ✅ | ✅ | ✅ |
| 跳转专页 | ✅ | ✅ | ✅ |
| drain/维护 GPU | ❌ | Phase 2 | Phase 2 |

## 删除前置校验

**N/A**

## 接口冻结规则

### `GET /api/v1/gpu-inventory`

- Query：`vendor`、`state`、`limit`、`cursor`
- `success`：`200` + `GPUInventoryListResponse`
- `errors`：`401`、`403`
- 当前 YAML 未声明 BadRequest response
- BOSS：**只读参考**；正式须 platform path

### `GET /api/v1/gpu-inventory/occupancy`

- `success`：`200` + `GPUOccupancySummary`
- `errors`：`401`、`403`

### `GET /api/v1/observability/query`

- `query` required；`errors`：`400`、`401`、`403`

### 平台 GPU aggregate（待补）

<!-- ADDED-TO-YAML: GET /api/v1/platform/trends/gpu -->

## 使用规则

- 利用率优先窗口聚合，非逐秒实时
- 禁止跨租户 JWT 轮询
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 平台 GPU aggregate + 租户排行 — **ADDED-TO-YAML**
- curated PromQL 模板 — **TODO-YAML**
- GPU drain — Phase 2 → maint-skills

## 响应示例

### occupancy 引用成功（租户 · **非 BOSS 正式**）

```json
{
  "total_gpus": 8,
  "allocated_gpus": 6,
  "available_gpus": 2,
  "allocation_rate": 0.75
}
```

### 平台 aggregate 目标（**待 YAML**）

```json
{
  "gpu_total": 512,
  "gpu_allocated": 401,
  "avg_utilization_pct": 72.5,
  "time_window": "7d",
  "top_tenants": [{ "tenant_id": "t-001", "gpu_allocated": 64 }]
}
```

## 错误示例

### 无 gpu-inventory 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: scope:gpu-inventory:read required",
  "request_id": "req-boss-gpt-403-001"
}
```

### PromQL 为空

```json
{
  "code": "BAD_REQUEST",
  "message": "query is required and must not be empty",
  "request_id": "req-boss-gpt-400-001"
}
```

## 相关模块

- [`operations-overview.md`](operations-overview.md)、[`resource-capacity-trend.md`](resource-capacity-trend.md)
- [`gpu-monitoring.md`](../health/gpu-monitoring.md)、[`platform-gpu-hours.md`](../metering/platform-gpu-hours.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] gpu-inventory 已声明路径 vs 平台 aggregate 分层清晰
- [x] 400 + 403 错误示例
- [x] 平台 aggregate YAML 已合入
- [x] PRD/SPEC/HTML 与本文同步
