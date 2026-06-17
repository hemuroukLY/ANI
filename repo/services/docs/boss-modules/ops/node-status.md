# 节点状态

## 页面定位

`节点状态` 是 BOSS **资源池与基础设施** 域下的 **全平台 K8s/裸机节点巡检** 专页：查看节点角色、就绪态、资源余量、GPU 驱动、心跳与污点；提供维护标记、drain 入口（Phase 2 → maint-skills）。

Console **无** 直接对等页（平台专属）。关联：[`gpu-pool-management.md`](gpu-pool-management.md)、[`platform-health.md`](../health/platform-health.md)。

## 文档管理规则

- 本文是 `节点状态` 的主维护源
- PRD/SPEC 为辅助
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- **当前 v1.yaml 无** `GET /api/v1/nodes` 或等价平台节点 list — **TODO-YAML**
- 可引用：`GET /api/v1/observability/query`（节点指标 PromQL）、`GET /api/v1/gpu-inventory`（节点-GPU 绑定 **只读参考**）
- K8s `/readyz`、控制面元数据 — Phase 1 文档口径；**不得** 写成 Core 已冻结 path
- 平台 RBAC — **TODO-YAML**
- drain / 维护写操作 — Phase 2 → [`maint-skills.md`](../health/maint-skills.md)

## 页面职责

- 列表展示全平台节点：名、角色、状态、CPU/内存/磁盘/GPU、驱动、K8s 版本、最近心跳
- 详情抽屉：标签、污点、Pod 摘要、最近事件、指标曲线
- 支持筛选：角色、状态、AZ、GPU 型号
- **不承担** 在本页直接调用未声明的 K8s API 写操作

## 页面结构

```text
节点状态
├── 筛选（角色 / 状态 / AZ / GPU）
├── 节点表格
├── 详情抽屉
│   ├── 标签与污点
│   ├── Pod 列表摘要
│   ├── 最近事件
│   └── 指标曲线 → observability/query
└── 操作
    ├── 查看（只读）
    ├── 标记维护（Phase 2）
    └── drain（Phase 2 → maint-skills）
```

## 数据来源与分层约束

| 层 | 来源 | 用法 |
|---|---|---|
| Core | 平台 nodes list **TODO-YAML** | BOSS 正式数据源 |
| Core | `GET /api/v1/observability/query` | CPU/内存/磁盘/GPU 指标 |
| Core | `GET /api/v1/gpu-inventory` | 节点 GPU 绑定参考 |
| 运维 | K8s Node API / 主机 agent | Phase 1 口径；待映射 YAML |

### 关键边界

- 不得自造 `PlatformNode` schema 为已冻结
- 节点 `Ready`/`NotReady` 须与未来 YAML enum 对齐；当前用产品目标三态：`ready` / `not_ready` / `maintenance`
- PromQL **≠** nodes list

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| 节点表格 | nodes aggregate **待 YAML** | — |
| GPU 列 | gpu-inventory 按 `node_id` | gpu-pool-management |
| 指标曲线 | PromQL | observability/query |
| 事件 | K8s events **待 YAML** | platform-logs |

## BOSS 与 Console 分工

| 维度 | BOSS | Console |
|---|---|---|
| 范围 | 全平台节点 | 不暴露 |
| 操作 | drain/维护（Phase 2） | 无 |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL |
| GET | `/api/v1/gpu-inventory` | `listGPUInventory` | 节点 GPU 参考 |

| 能力 | 状态 |
|---|---|
| `GET /api/v1/platform/nodes` | **ADDED-TO-YAML** P1 (`listPlatformNodes`) |
| 节点 drain / cordon API | Phase 2 |
| K8s events 平台 list | **TODO-YAML** |

## 字段级定义

| 字段 | 说明 |
|---|---|
| `node_name` | 节点名 |
| `role` | control-plane / worker / gpu 等（产品 enum · 待 YAML） |
| `status` | `ready` / `not_ready` / `maintenance` |
| `cpu_allocatable` | 可分配 CPU |
| `memory_allocatable_gb` | 可分配内存 GB |
| `disk_pressure` | 布尔或等级 |
| `gpu_count` | 绑定 GPU 数 |
| `driver_version` | GPU 驱动版本 |
| `k8s_version` | K8s 版本 |
| `last_heartbeat_at` | 最近心跳 |
| `taints[]` | 污点列表 |
| `last_refreshed_at` | UI |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 正常 | 表格 + 抽屉 |
| nodes API 未就绪 | 「平台 nodes API 待 Core 合入」 |
| `not_ready` | 黄/红高亮 |
| PromQL 失败 | 指标列独立失败 |
| 403 | 无权限 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `cpu_allocatable` | aggregate 定义前标注产品目标 | millicores 或 cores |
| `memory_allocatable_gb` | 节点可分配内存 | GB |
| `gpu_count` | 节点 GPU 设备数 | integer |
| `last_heartbeat_at` | 最近心跳 | ISO 8601 date-time |
| `uptime_seconds` | 节点在线时长 | seconds；UI 可换算天/小时 |

## 状态与能力口径

| 状态 | 含义 |
|---|---|
| `ready` | 可调度 |
| `not_ready` | 不可用或 NotReady |
| `maintenance` | 维护/ cordon |

写操作（drain、cordon）— Phase 2 Skills。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台基础设施读 RBAC（**TODO-YAML**） | `403` |
| observability 读（指标） | `403` |
| 未认证 | `401` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 管理员 |
|---|---|---|---|
| 查看节点列表 | ✅ | ✅ | ✅ |
| 查看详情/指标 | ✅ | ✅ | ✅ |
| 标记维护 / drain | ❌ | Phase 2 | Phase 2 |
| 删除节点 | ❌ | ❌ | ❌ |

## 删除前置校验

**N/A**（本页 Phase 1 无 DELETE 节点）

## 接口冻结规则

### `GET /api/v1/observability/query`

- `scope:observability:read`；`errors`：`400`、`401`、`403`

### 平台 nodes list（待补）

<!-- TODO-YAML: GET /api/v1/nodes 或 /platform/nodes -->

- 须含 `node_name`、`status`、`role`、资源余量、心跳
- 合入前不得写入「已冻结」

## 使用规则

- 禁止伪造全 green 节点表
- K8s 直连仅 Phase 2 文档口径
- 与 gpu-inventory 的 `node_id` 关联须一致

## 待补能力边界

- 平台 nodes list API — **ADDED-TO-YAML**
- drain / cordon — Phase 2 → maint-skills
- Pod 列表平台 API — **TODO-YAML**

## 响应示例

### 平台 nodes 目标（**待 YAML**）

```json
{
  "items": [
    {
      "node_name": "gpu-worker-12",
      "role": "gpu-worker",
      "status": "ready",
      "cpu_allocatable": 64,
      "memory_allocatable_gb": 512,
      "gpu_count": 8,
      "driver_version": "550.54.15",
      "k8s_version": "1.29.4",
      "last_heartbeat_at": "2026-06-17T11:00:00Z"
    }
  ]
}
```

## 错误示例

### 无平台基础设施读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-ns-403-001"
}
```

### PromQL 为空

```json
{
  "code": "BAD_REQUEST",
  "message": "query is required and must not be empty",
  "request_id": "req-boss-ns-400-001"
}
```

## 相关模块

- [`platform-resource-pool.md`](platform-resource-pool.md)、[`gpu-pool-management.md`](gpu-pool-management.md)
- [`maint-skills.md`](../health/maint-skills.md)、[`platform-logs.md`](../health/platform-logs.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] 明确无 nodes list；observability 引用边界
- [x] 400 + 403 错误示例
- [ ] nodes YAML 合入后回写
- [x] PRD/SPEC/HTML 与本文同步
