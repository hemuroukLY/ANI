# 平台资源池总览

## 页面定位

`平台资源池总览` 是 BOSS **资源池与基础设施** 域下的 **算力域 landing 聚合页**，面向平台管理员与 SRE，汇总全平台 GPU、节点、存储与网络基础设施的容量摘要与健康态势，并提供跳转至各子模块。

本页是 **BOSS** 页面，不是 Console 租户资源池页。Console 对照：[`resource-pool-overview.md`](../../console-modules/compute/resource-pool-overview.md)（单租户）。与 overview 域 [`resource-capacity-trend.md`](../overview/resource-capacity-trend.md) 分工：本页偏 **基础设施运维入口**，彼页偏 **运营总览趋势**。

## 文档管理规则

- 本文是 `平台资源池总览` 的主维护源
- PRD/SPEC 为辅助；冲突以本文 + OpenAPI 为准
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- **当前无** 独立 `GET /api/v1/resource-pools/platform` 冻结路径
- 可引用租户上下文 **只读参考**：`gpu-inventory`、`instances`、`volumes`、`networks/vpcs` 等
- BOSS 平台 resource-pool aggregate — **TODO-YAML**
- 平台 RBAC — **TODO-YAML**；禁止信任未授权 `tenant_id`
- 错误结构：`code` / `message` / `request_id`
- 禁止自造 `ResourcePoolPlatformOverview` 为已冻结 schema

## 页面职责

- 提供算力域四象限摘要：GPU 池、节点、存储、网络
- 为 [`gpu-pool-management.md`](gpu-pool-management.md)、[`node-status.md`](node-status.md)、[`storage-infrastructure.md`](storage-infrastructure.md)、[`network-infrastructure.md`](network-infrastructure.md) 提供深链
- 展示 `last_refreshed_at` 与区块级局部失败态
- **不承担** 各子资源 CRUD（跳转专页）

## 页面结构

```text
平台资源池总览
├── 顶部总判断（资源池健康 rollup）
├── 四象限摘要卡片
│   ├── GPU 资源池 → gpu-pool-management
│   ├── 节点状态 → node-status
│   ├── 存储基础设施 → storage-infrastructure
│   └── 网络基础设施 → network-infrastructure
└── 快捷跳转
    ├── GPU 监控（health 域）
    └── 资源池与容量态势（overview 域）
```

## 数据来源与分层约束

| 层 | 路径 | 用法 |
|---|---|---|
| Core | `GET /api/v1/gpu-inventory`、`/occupancy` | 租户 **只读参考** |
| Core | `GET /api/v1/instances` | 租户 **只读参考** |
| Core | `GET /api/v1/volumes` | 租户 **只读参考** |
| Core | `GET /api/v1/networks/vpcs` 等 | 租户 **只读参考** |
| Core | `GET /api/v1/observability/query` | PromQL |
| Core | 平台 pool aggregate **TODO-YAML** | BOSS 正式数据源 |

### 关键边界

- Console `resource-pool-overview` 为单租户；本页为全平台
- 不得逐租户 JWT 轮询冒充平台摘要
- 子模块字段口径以各 ops 专页为准

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| GPU 摘要 | gpu-inventory 参考 / aggregate | gpu-pool-management |
| 节点摘要 | nodes list **TODO-YAML** | node-status |
| 存储摘要 | volumes 参考 / aggregate | storage-infrastructure |
| 网络摘要 | networks 参考 / aggregate | network-infrastructure |

## BOSS 与 Console 分工

| 维度 | BOSS 本页 | Console resource-pool-overview |
|---|---|---|
| 范围 | 全平台基础设施 | 当前租户资源摘要 |
| 网络/存储 | 平台基础设施视角 | 租户 VPC/卷 |
| API | platform aggregate 待 YAML | 租户页面聚合 |

## 当前冻结事实

| 方法 | 路径 | operationId | BOSS 用法 |
|---|---|---|---|
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | 租户 **只读参考** |
| GET | `/api/v1/gpu-inventory/occupancy` | `getGPUOccupancy` | 租户 **只读参考** |
| GET | `/api/v1/instances` | `listInstances` | 租户 **只读参考** |
| GET | `/api/v1/volumes` | `listStorageVolumes` | 租户 **只读参考** |
| GET | `/api/v1/networks/vpcs` | `listNetworkVPCs` | 租户 **只读参考** |

| 能力 | 状态 |
|---|---|
| 平台 resource-pool aggregate | **ADDED-TO-YAML** P1 |
| 平台 nodes list | **ADDED-TO-YAML** `listPlatformNodes` |

## 字段级定义

| 字段 | 说明 |
|---|---|
| `pool_overall_status` | `ok` / `degraded` / `error`（UI rollup） |
| `gpu_total` | 全平台 GPU 数 |
| `gpu_utilization_pct` | 窗口利用率 |
| `node_ready_count` | ready 节点数 — **待 YAML** |
| `node_not_ready_count` | not ready — **待 YAML** |
| `storage_used_pct` | 存储使用率摘要 |
| `network_vpc_count` | VPC 总数（平台视角待 YAML） |
| `last_refreshed_at` | UI |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 正常 | 四卡片完整 |
| 单卡片 API 未就绪 | 「待 YAML」；不伪造 |
| 单卡片失败 | 局部失败 + 重试 |
| 403 | 无权限 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `gpu.utilization_pct` | 指定 `time_window` 内 GPU 利用率 | %，保留 1 位小数 |
| `nodes.ready` / `nodes.not_ready` | 平台节点 rollup | integer |
| `storage.used_pct` | 平台 aggregate 目标；合入前不得冒充冻结 | % |
| `network.vpc_count` | 平台网络资源计数 | integer |
| `last_refreshed_at` | 数据刷新时间 | ISO 8601 date-time |

## 状态与能力口径

本页 **只读** landing；无独立状态机。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台 BOSS 读 RBAC（**TODO-YAML**） | `403` |
| 未认证 | `401` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 管理员 |
|---|---|---|---|
| 查看总览 | ✅ | ✅ | ✅ |
| 跳转子模块 | ✅ | ✅ | ✅ |
| 在本页创建资源 | ❌ | ❌ | ❌ |

## 删除前置校验

**N/A**

## 接口冻结规则

### 本页无独立冻结 API

引用子模块 operation 须标注 **租户只读参考** vs **平台 TODO-YAML**。

### `GET /api/v1/gpu-inventory`（租户 · **只读参考**）

- `operationId`：`listGPUInventory`；`errors`：`401`、`403`
- 当前 YAML 未声明 BadRequest response

### 平台 pool aggregate（待补）

<!-- ADDED-TO-YAML: GET /api/v1/platform/resource-pools -->

## 使用规则

- 四象限统一刷新上下文
- 禁止伪造跨租户数字
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 平台 resource-pool aggregate — **ADDED-TO-YAML**
- 与 overview `resource-capacity-trend` 口径对齐 — 产品待确认

## 响应示例

### 平台 aggregate 目标（**待 YAML**）

```json
{
  "pool_overall_status": "ok",
  "last_refreshed_at": "2026-06-17T11:00:00Z",
  "gpu": { "total": 512, "utilization_pct": 76.2 },
  "nodes": { "ready": 48, "not_ready": 2 },
  "storage": { "used_pct": 68.5 },
  "network": { "vpc_count": 120 }
}
```

## 错误示例

### 无平台读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-prp-403-001"
}
```

### 非法时间窗口（平台 aggregate · **TODO-YAML** 目标）

```json
{
  "code": "BAD_REQUEST",
  "message": "time_window must be one of: 5m, 1h, 24h, 7d",
  "request_id": "req-boss-prp-400-001"
}
```

## 相关模块

- [`gpu-pool-management.md`](gpu-pool-management.md)、[`node-status.md`](node-status.md)
- [`storage-infrastructure.md`](storage-infrastructure.md)、[`network-infrastructure.md`](network-infrastructure.md)
- overview：[`resource-capacity-trend.md`](../overview/resource-capacity-trend.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] 无独立 API；子路径分层清晰
- [x] 400 + 403 错误示例（403 为主）
- [ ] platform aggregate YAML 合入后回写
- [x] PRD/SPEC/HTML 与本文同步
