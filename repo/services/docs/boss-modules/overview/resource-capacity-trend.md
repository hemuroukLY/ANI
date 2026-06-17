# 资源池与容量态势

## 页面定位

`资源池与容量态势` 是 BOSS **平台运营总览** 域下的 **全平台资源池容量趋势** 明细页，展示 GPU、CPU、内存、存储与实例规模的跨租户容量摘要、利用率趋势与紧张度预警。

本页是 **BOSS** 页面，不是 Console 租户资源池页。Console 对照：[`resource-pool-overview.md`](../../console-modules/compute/resource-pool-overview.md)（单租户算力域 landing）。

深度专页：[`platform-resource-pool.md`](../ops/platform-resource-pool.md)、[`gpu-monitoring.md`](../health/gpu-monitoring.md)、[`resource-capacity-trend` 上级摘要](operations-overview.md)。

## 文档管理规则

- 本文是 `资源池与容量态势` 的主维护源
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-resource-capacity-trend.md`](../../tasks/modules/prd/boss/overview/prd-boss-resource-capacity-trend.md) 与 [`spec-boss-resource-capacity-trend.md`](../../tasks/modules/spec/boss/overview/spec-boss-resource-capacity-trend.md) 为辅助材料
- 流程：ANI-14 Phase 1 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- Core `/api/v1/*`；Services `/api/v1/svc/*`
- **当前无** 独立 `GET /api/v1/resource-pools/platform` 或等价 BOSS 冻结路径
- 可引用租户上下文 **只读参考**：`GET /api/v1/gpu-inventory`、`/gpu-inventory/occupancy`、`GET /api/v1/instances`（分 kind 计数）、`GET /api/v1/metering/usage`
- BOSS 跨租户容量 aggregate — **TODO-YAML**；**不得** 逐租户 JWT 轮询冒充平台大盘
- 平台 RBAC（**TODO-YAML**）；不信任 body/query 未授权 `tenant_id`
- 错误结构：`code` / `message` / `request_id`
- 禁止自造 `ResourcePoolPlatformSummary` schema 为已冻结

## 页面职责

- 展示全平台 GPU/CPU/内存/存储容量与利用率 **趋势**
- 识别容量紧张 AZ/节点池与 Top 租户占用（aggregate 待 YAML）
- 提供跳转：GPU 监控、平台资源池总览、计量专页
- 标明统计窗口与非实时口径
- **不承担** 资源创建/扩缩容（跳转 ops/health 专页）

## 页面结构

```text
资源池与容量态势
├── 顶部 KPI（GPU/CPU/内存/存储 总量·已用·利用率）
├── 趋势图（7d / 30d 窗口）
├── 紧张度排行（AZ / 租户 — 待 YAML）
├── 筛选（时间窗口、资源类型、AZ）
└── 跳转入口
    ├── GPU 监控
    ├── 平台资源池总览
    └── 平台 GPU-Hours 计量
```

## 数据来源与分层约束

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/gpu-inventory` | 租户 **只读参考**；`GPUInventoryDevice` |
| Core | `GET /api/v1/gpu-inventory/occupancy` | 租户 **只读参考**；`GPUOccupancySummary` |
| Core | `GET /api/v1/instances` | 分 `kind` 计数参考 |
| Core | `GET /api/v1/metering/usage` | 租户 **只读参考**；`instance_*` resource_type |
| Core | `GET /api/v1/observability/query` | PromQL 利用率曲线 |
| Core | 平台 capacity aggregate **TODO-YAML** | BOSS 正式数据源 |

### 关键边界

- Console `resource-pool-overview` 为单租户；本页为 **全平台**
- `metering/usage` **无** `group_by=tenant_id` — 平台扩展待 YAML
- 不得把 `instances` 租户 list 直接写成跨租户排行

## 页面区块与数据来源映射

| 区块 | 主要来源 | 说明 | 跳转 |
|---|---|---|---|
| GPU KPI | aggregate + gpu-inventory 参考 | 总量/已分配/利用率 | gpu-monitoring |
| CPU/内存 KPI | metering + PromQL | `instance_cpu_seconds` 等 | metering 专页 |
| 存储 KPI | metering **待 enum** | storage resource_type 待 YAML | platform-storage-gbdays |
| 趋势图 | aggregate 或 PromQL | 须标注 `time_window` | — |
| 紧张度排行 | **TODO-YAML** | 按 AZ/租户 | tenant-list |

## BOSS 与 Console 分工

| 维度 | BOSS 本页 | Console resource-pool-overview |
|---|---|---|
| 范围 | 全平台资源池容量 | 当前租户可见资源 |
| 排行 | 跨租户 Top N（待 YAML） | 无 |
| API | 平台 aggregate 待 YAML | 页面层租户聚合 |

## 当前冻结事实

| 方法 | 路径 | operationId | RBAC | BOSS 用法 |
|---|---|---|---|---|
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | `scope:gpu-inventory:read` | 租户 **只读参考** |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | `scope:gpu-inventory:read` | 租户 **只读参考** |
| GET | `/api/v1/instances` | `listInstances` | `scope:instances:read` | 租户 **只读参考** |
| GET | `/api/v1/metering/usage` | — | 租户 JWT | **只读参考**；无 scope |
| GET | `/api/v1/observability/query` | `queryObservability` | `scope:observability:read` | PromQL |

| 能力 | 状态 |
|---|---|
| 平台 capacity aggregate | **ADDED-TO-YAML** `getPlatformCapacity` |
| 平台 RBAC scope | **TODO-YAML** |

## 字段级定义

| 字段 | 说明 | 来源 |
|---|---|---|
| `gpu_total` | 全平台 GPU 设备数 | aggregate 目标 |
| `gpu_allocated` | 已分配 GPU | occupancy / aggregate |
| `gpu_utilization_pct` | 窗口平均利用率 | PromQL 或 aggregate |
| `cpu_hours_window` | 窗口 CPU 用量 | metering `instance_cpu_seconds` |
| `memory_gbhours_window` | 窗口内存 GB·小时 | metering |
| `storage_used_gbdays` | 存储 GB·天 | metering enum **待 YAML** |
| `instance_count_by_kind` | 分 kind 实例数 | instances 聚合 |
| `capacity_pressure_score` | 紧张度（产品计算） | UI · 非 API |
| `last_refreshed_at` | 刷新时间 | UI |
| `time_window` | `24h`/`7d`/`30d` | UI query |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 正常态 | KPI + 趋势完整渲染 |
| aggregate 未就绪 | 「平台 capacity API 待 Core 合入」 |
| 单指标 PromQL 失败 | 该曲线失败提示，KPI 可保留 |
| 无权限 | 403，不展示假数据 |
| 存储 enum 未冻结 | 存储卡片「待 YAML」 |

## 字段口径与单位

| 字段 | 单位 |
|---|---|
| `gpu_utilization_pct` | 0–100 % |
| `cpu_hours_window` | 小时（由 seconds 换算） |
| `memory_gbhours_window` | GB·小时 |
| `storage_used_gbdays` | GB·天 |

## 状态与能力口径

本页 **只读**；无资源状态机。紧张度为 UI 衍生指标，须在 aggregate YAML 前标注「产品计算」。

## 创建前置条件

| 依赖项 | 未满足响应 |
|---|---|
| 平台 BOSS 读 RBAC（**TODO-YAML**） | `403` |
| 已认证 | `401` |
| PromQL `query` 非空 | `400` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 |
|---|---|---|---|
| 查看态势 | ✅ | ✅ | ✅ |
| 切换时间窗口 | ✅ | ✅ | ✅ |
| 跳转 GPU/计量 | ✅ | ✅ | ✅ |
| 在本页扩缩容 | ❌ | ❌ | ❌ |
| 导出 CSV | Phase 2 | Phase 2 | Phase 2 |

## 删除前置校验

**N/A**（无 DELETE）

## 接口冻结规则

### `GET /api/v1/gpu-inventory`（租户 · **只读参考**）

- `operationId`：`listGPUInventory`；`success`：`200`；`errors`：`401`、`403`
- 当前 YAML 未声明 BadRequest response
- **非** BOSS 正式跨租户契约

### `GET /api/v1/metering/usage`（租户 · **只读参考**）

- 无 operationId；`success`：`200`；`errors`：`400`、`401`、`403`
- `group_by` enum **不含** `tenant_id`

### 平台 capacity aggregate（待补）

<!-- ADDED-TO-YAML: GET /api/v1/platform/capacity -->

## 使用规则

- 趋势须标注 `time_window` 与 `last_refreshed_at`
- 禁止伪造跨租户排行
- 租户 API 仅文档引用，不得作为 BOSS 生产架构

## 待补能力边界

- 平台 capacity aggregate — **ADDED-TO-YAML**
- 存储 metering `resource_type` enum — **TODO-YAML**
- CSV 导出 — Phase 2

## 响应示例

### 平台 aggregate 目标（**待 YAML**）

```json
{
  "time_window": "7d",
  "last_refreshed_at": "2026-06-17T10:00:00Z",
  "gpu": { "total": 512, "allocated": 401, "utilization_pct": 78.3 },
  "cpu_hours": 89200,
  "memory_gbhours": 1240000,
  "storage_gbdays": 45000,
  "top_tenants_by_gpu": [
    { "tenant_id": "t-001", "gpu_allocated": 64 }
  ]
}
```

## 错误示例

### 无平台读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-rct-403-001"
}
```

### metering 时间范围无效

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be before end_time",
  "request_id": "req-boss-rct-400-001"
}
```

## 相关模块

- [`operations-overview.md`](operations-overview.md)、[`gpu-pool-trend.md`](gpu-pool-trend.md)
- [`gpu-monitoring.md`](../health/gpu-monitoring.md)、[`platform-resource-pool.md`](../ops/platform-resource-pool.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] 租户路径标注只读参考；平台 aggregate TODO-YAML
- [x] 含 400 + 403 错误示例
- [ ] 平台 capacity YAML 合入后回写冻结表
- [x] PRD/SPEC/HTML 与本文同步
