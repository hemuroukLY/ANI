# 网络基础设施

## 页面定位

`网络基础设施` 是 BOSS **资源池与基础设施** 域下的 **全平台网络控制面** 专页：汇总 VPC、子网、负载均衡与安全组规模的跨租户态势，以及底层 CNI/网关健康（aggregate 待 YAML）。

Console 对照：[`network-management.md`](../../console-modules/compute/network-management.md)（租户网络 CRUD）。

## 文档管理规则

- 本文是 `网络基础设施` 的主维护源
- PRD/SPEC 为辅助
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- 租户 **只读参考**（YAML 已声明路径，BOSS 非正式契约）：
  - `GET /api/v1/networks/vpcs`（`listNetworkVPCs`）
  - `GET /api/v1/networks/subnets`（`listNetworkSubnets`）
  - `GET /api/v1/networks/security-groups`（`listNetworkSecurityGroups`）
  - `GET /api/v1/networks/load-balancers`（`listNetworkLoadBalancers`）
- 平台 network infrastructure aggregate — **TODO-YAML**
- **不得** 逐租户 JWT 轮询统计全平台 VPC 数
- 平台 RBAC — **TODO-YAML**

## 页面职责

- 展示全平台 VPC/子网/LB/SG 数量与增长趋势
- CNI/网关/Underlay 健康摘要（**待 YAML** 或 PromQL）
- 跨租户网络资源 Top N（**TODO-YAML**）
- 跳转租户网络管理文档作字段口径参考
- **不承担** 租户 VPC CRUD

## 页面结构

```text
网络基础设施
├── 顶部 KPI（VPC / 子网 / LB / SG 计数）
├── 控制面健康（CNI / 网关 — 待 YAML）
├── 增长趋势
├── 跨租户排行（待 YAML）
└── 跳转
    └── Console network-management（口径参考）
```

## 数据来源与分层约束

| 层 | 路径 | 用法 |
|---|---|---|
| Core | `GET /api/v1/networks/vpcs` | 租户 **只读参考** |
| Core | `GET /api/v1/networks/subnets` | 租户 **只读参考** |
| Core | `GET /api/v1/networks/security-groups` | 租户 **只读参考** |
| Core | `GET /api/v1/networks/load-balancers` | 租户 **只读参考** |
| Core | `GET /api/v1/observability/query` | 网络组件 PromQL |
| Core | 平台 network aggregate **TODO-YAML** | BOSS 正式数据源 |

### 关键边界

- Console 网络 API 为 **租户边界**；BOSS 须 platform aggregate
- 计数口径：VPC 总数 ≠ 单租户 list 长度简单相加（禁止轮询）

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| 资源计数 KPI | aggregate **待 YAML** | — |
| 控制面健康 | health / PromQL | platform-health |
| 趋势 | aggregate 时间序列 | — |
| 租户排行 | **TODO-YAML** | tenant-list |

## BOSS 与 Console 分工

| 维度 | BOSS | Console network-management |
|---|---|---|
| 视角 | 全平台网络基础设施 | 当前租户 VPC/子网/LB |
| 操作 | 只读态势 | CRUD |
| API | platform aggregate | `networks/*` |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/networks/vpcs` | `listNetworkVPCs` | 租户 **只读参考** |
| GET | `/api/v1/networks/subnets` | `listNetworkSubnets` | 租户 **只读参考** |
| GET | `/api/v1/networks/security-groups` | `listNetworkSecurityGroups` | 租户 **只读参考** |
| GET | `/api/v1/networks/load-balancers` | `listNetworkLoadBalancers` | 租户 **只读参考** |
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL |

| 能力 | 状态 |
|---|---|
| 平台 network infrastructure aggregate | **ADDED-TO-YAML** P1 |
| CNI 健康 curated 指标 | **TODO-YAML** |

## 字段级定义

| 字段 | 说明 |
|---|---|
| `vpc_count` | 全平台 VPC 数 |
| `subnet_count` | 子网数 |
| `lb_count` | 负载均衡数 |
| `sg_count` | 安全组数 |
| `cni_status` | `ok` / `degraded` / `error` |
| `tenant_rank[]` | **待 YAML** |
| `last_refreshed_at` | UI |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 正常 | KPI + 趋势 |
| aggregate 未就绪 | 「待 YAML」 |
| `cni_status` degraded/error | 高亮 |
| 403 | 无权限 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `vpc_count` / `subnet_count` | 平台 aggregate 目标计数 | integer |
| `lb_count` / `sg_count` | 平台 aggregate 目标计数 | integer |
| `cni_status` | 三态健康 rollup | `ok` / `warn` / `critical` |
| 趋势字段 | 必须绑定 `time_window` | ISO 8601 时间窗 |

## 状态与能力口径

`cni_status` 对齐三态健康 rollup。本页 **只读**。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台基础设施读（**TODO-YAML**） | `403` |
| 未认证 | `401` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 管理员 |
|---|---|---|---|
| 查看网络基础设施 | ✅ | ✅ | ✅ |
| 在本页创建 VPC | ❌ | ❌ | ❌ |
| 跳转 Console 口径文档 | ✅ | ✅ | ✅ |

## 删除前置校验

**N/A**

## 接口冻结规则

### `GET /api/v1/networks/vpcs`（租户 · **只读参考**）

- `operationId`：`listNetworkVPCs`；`success`：`200`；`errors`：`401`、`403`
- 详参见 Console [`network-management.md`](../../console-modules/compute/network-management.md)

### 平台 network aggregate（待补）

<!-- ADDED-TO-YAML: GET /api/v1/platform/networks -->

## 使用规则

- 禁止跨租户 list 轮询冒充平台 KPI
- 引用 Console 冻结 operation 时保持成功码/错误码一致
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 平台 network aggregate — **ADDED-TO-YAML**
- Underlay 拓扑视图 — Phase 2

## 响应示例

### 平台 aggregate 目标（**待 YAML**）

```json
{
  "vpc_count": 120,
  "subnet_count": 480,
  "lb_count": 65,
  "sg_count": 310,
  "cni_status": "ok",
  "last_refreshed_at": "2026-06-17T11:00:00Z"
}
```

## 错误示例

### 无平台读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-ni-403-001"
}
```

### 租户 networks 无权限（引用场景）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-ni-403-002"
}
```

### 非法分页参数（平台 aggregate · **TODO-YAML** 目标）

```json
{
  "code": "BAD_REQUEST",
  "message": "limit must be between 1 and 100",
  "request_id": "req-boss-ni-400-001"
}
```

## 相关模块

- [`platform-resource-pool.md`](platform-resource-pool.md)
- Console：[`network-management.md`](../../console-modules/compute/network-management.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] networks 租户路径 vs 平台 aggregate 分层
- [x] 403 + 400 错误示例
- [ ] network aggregate YAML 合入后回写
- [x] PRD/SPEC/HTML 与本文同步
