# 存储基础设施

## 页面定位

`存储基础设施` 是 BOSS **资源池与基础设施** 域下的 **全平台存储后端与容量** 专页：汇总对象存储/块存储集群容量、使用率、健康态与租户占用分布（aggregate 待 YAML）。

Console 对照：[`storage-management.md`](../../console-modules/compute/storage-management.md)（租户卷 CRUD）。计量：[`platform-storage-gbdays.md`](../metering/platform-storage-gbdays.md)。

## 文档管理规则

- 本文是 `存储基础设施` 的主维护源
- PRD/SPEC 为辅助
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- 租户 **只读参考**：`GET /api/v1/volumes`（`listStorageVolumes`）、`GET /api/v1/volumes/{volume_id}`
- `GET /api/v1/metering/usage` — 存储 `resource_type` enum **TODO-YAML**（metering 域）
- 平台 storage infrastructure aggregate — **TODO-YAML**
- **不得** 把租户 volumes list 直接写成 BOSS 跨租户正式契约
- 平台 RBAC — **TODO-YAML**

## 页面职责

- 展示全平台存储池总容量、已用、可用、使用率趋势
- 块存储 vs 对象存储分组（产品分区 · 待 YAML）
- 跨租户存储占用 Top N（**TODO-YAML**）
- 跳转 Console 口径 volumes 文档、platform-storage-gbdays、platform-logs
- **不承担** 租户卷 CRUD

## 页面结构

```text
存储基础设施
├── 顶部 KPI（总容量 / 已用 / 使用率）
├── 存储池分组（块 / 对象 — 待 YAML）
├── 使用率趋势
├── 跨租户占用排行（待 YAML）
└── 跳转
    ├── 平台 Storage-GBDays
    └── 平台日志
```

## 数据来源与分层约束

| 层 | 路径 | 用法 |
|---|---|---|
| Core | `GET /api/v1/volumes` | 租户 **只读参考** |
| Core | `GET /api/v1/metering/usage` | 存储计量 enum **待 YAML** |
| Core | `GET /api/v1/observability/query` | 容量/IO PromQL |
| Core | 平台 storage aggregate **TODO-YAML** | BOSS 正式数据源 |

### 关键边界

- Console `volumes` 为租户资源；本页为 **基础设施/backend** 视角
- metering 存储 resource_type 未冻结前不得写死 enum 值

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| 容量 KPI | aggregate **待 YAML** | — |
| 趋势 | metering + PromQL | platform-storage-gbdays |
| 租户排行 | **TODO-YAML** | tenant-list |
| 健康告警 | observability / health | platform-health |

## BOSS 与 Console 分工

| 维度 | BOSS | Console storage-management |
|---|---|---|
| 视角 | 全平台存储后端 | 租户卷生命周期 |
| API | platform aggregate | `volumes` CRUD |
| 计量 | 跨租户 GB·天 | 租户 usage |

## 当前冻结事实

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| GET | `/api/v1/volumes` | `listStorageVolumes` | 租户 scope（见 YAML） |
| GET | `/api/v1/metering/usage` | — | 租户 JWT；无 scope |
| GET | `/api/v1/observability/query` | `queryObservability` | `scope:observability:read` |

| 能力 | 状态 |
|---|---|
| 平台 storage infrastructure aggregate | **ADDED-TO-YAML** P1 |
| metering storage `resource_type` | **TODO-YAML** |

## 字段级定义

| 字段 | 说明 |
|---|---|
| `total_capacity_tb` | 总容量 TB |
| `used_capacity_tb` | 已用 TB |
| `used_pct` | 使用率 % |
| `pool_name` | 存储池名（待 YAML） |
| `pool_type` | block / object |
| `health_status` | `ok` / `degraded` / `error` |
| `tenant_rank[]` | **待 YAML** |
| `last_refreshed_at` | UI |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 正常 | KPI + 趋势 |
| aggregate 未就绪 | 标注 TODO-YAML |
| metering enum 未冻结 | 用量卡片「待 YAML」 |
| `degraded`/`error` | 高亮 |
| 403 | 无权限 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `capacity_bytes` | 存储总容量 | bytes；UI 可换算 GB/TB |
| `used_bytes` | 已用容量 | bytes；UI 可换算 GB/TB |
| `used_pct` | `used_bytes / capacity_bytes * 100` | %，保留 1 位 |
| 趋势字段 | 必须绑定 `time_window` | ISO 8601 时间窗 |

## 状态与能力口径

`health_status` 对齐 `HealthCheck.status` 三态（aggregate 目标）。本页 Phase 1 **只读**。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台基础设施读（**TODO-YAML**） | `403` |
| metering 时间参数 | `400` |
| 未认证 | `401` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 管理员 |
|---|---|---|---|
| 查看存储基础设施 | ✅ | ✅ | ✅ |
| 跳转计量/日志 | ✅ | ✅ | ✅ |
| 在本页创建卷 | ❌ | ❌ | ❌ |

## 删除前置校验

**N/A**

## 接口冻结规则

### `GET /api/v1/volumes`（租户 · **只读参考**）

- `operationId`：`listStorageVolumes`；`success`：`200`；`errors`：`401`、`403`
- **非** BOSS 平台基础设施契约

### `GET /api/v1/metering/usage`（租户 · **只读参考**）

- `start_time`/`end_time` required；存储 `resource_type` — **TODO-YAML**

### 平台 storage aggregate（待补）

<!-- ADDED-TO-YAML: GET /api/v1/platform/storage -->

## 使用规则

- 区分租户卷与平台存储后端
- 禁止伪造跨租户占用
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 平台 storage aggregate — **ADDED-TO-YAML**
- 存储 metering enum — metering 域 **TODO-YAML**
- 存储池扩容操作 — Phase 2

## 响应示例

### 平台 aggregate 目标（**待 YAML**）

```json
{
  "total_capacity_tb": 1200,
  "used_capacity_tb": 820,
  "used_pct": 68.3,
  "pools": [
    { "pool_name": "ceph-block-1", "pool_type": "block", "health_status": "ok" }
  ],
  "last_refreshed_at": "2026-06-17T11:00:00Z"
}
```

## 错误示例

### 无平台读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-si-403-001"
}
```

### metering 时间范围无效

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be before end_time",
  "request_id": "req-boss-si-400-001"
}
```

## 相关模块

- [`platform-resource-pool.md`](platform-resource-pool.md)、[`platform-storage-gbdays.md`](../metering/platform-storage-gbdays.md)
- Console：[`storage-management.md`](../../console-modules/compute/storage-management.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] volumes 租户参考 vs 平台 aggregate 分层
- [x] 400 + 403 错误示例
- [ ] storage aggregate YAML 合入后回写
- [x] PRD/SPEC/HTML 与本文同步
