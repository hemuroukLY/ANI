# 首页 — 资源使用概览

## 页面定位

`资源使用概览` 是 `平台概览` 第一主题区块的**明细视角页**，用于展示当前租户下算力与存储资源的摘要用量，并引导进入各资源管理模块。

本页为 **Console UI 聚合页**，当前**无**独立冻结的 `GET /console/overview` 或 `/resource-overview` API。

## 文档管理规则

- 本文是「资源使用概览」区块的主维护文档
- 与 `platform-overview.md` 第一区块字段口径保持一致
- 不得自造 `ResourceOverviewResponse` schema
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 聚合逻辑在 Console 页面层，不写入 OpenAPI 新资源

## Core 层要求

本页通过并发调用各 Core 列表 / 摘要接口，在页面层轻聚合：

| 摘要字段 | 建议来源 API | RBAC |
|---|---|---|
| 实例总数 / 运行中 | `GET /api/v1/instances`（按 `kind` 筛选或汇总） | `scope:instances:read` |
| GPU 已分配 | `GET /api/v1/gpu-inventory` + `GET /api/v1/gpu-inventory/occupancy` | `scope:gpu-inventory:read` |
| 存储已使用 | `GET /api/v1/volumes`、`/filesystems`、`/objects` 等列表聚合 | 各子模块 scope |
| 网络 / 集群摘要 | `GET /api/v1/networks/vpcs`、`/k8s-clusters` 计数 | `scope:networks:read` 等 |

<!-- ADDED-TO-YAML: gpu-inventory 路径 (Core v1.yaml, Phase 2) -->

- 不允许继续使用 deprecated Services GPU 路径
- 页面不要求前端显式传 `tenant_id`

## Services 层要求

- 本页**不**消费 Services 业务 API 作为资源事实来源
- 推理 / 知识库用量不属于本区块（见 `home-inference-status.md`、`home-kb-trend.md`）

## 页面职责

- 展示实例、GPU、存储、网络 / 集群摘要
- 标明刷新时间与统计口径
- 提供跳转：VM、GPU、存储、网络、K8s
- 单来源失败时局部降级

## 字段级定义

与 `platform-overview.md` §资源使用概览 一致：

| 字段 | 说明 |
|---|---|
| 实例总数 | 当前租户可见实例总量 |
| 运行中实例数 | `state=running`（或各 kind 等价态）计数 |
| GPU 已分配量 | 已分配 / 总量（occupancy 或 inventory state） |
| 存储已使用量 | 各卷 / 文件系统 / 对象元数据聚合 |
| 资源更新时间 | 区块刷新时间 |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 各 Core 读权限 | 对应 scope | `403 FORBIDDEN`（单来源局部降级） |

本页无 POST/PUT 写操作，不涉及 `idempotency_key`。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员 |
|---|---|---|
| 查看摘要 | ✅ | ✅ |
| 刷新区块 | ✅ | ✅ |
| 跳转明细模块 | ✅ | ✅ |
| 在本页创建资源 | ❌ | ❌ |

## 接口冻结规则

本页**无独立冻结 API**。聚合引用的 Core operation 如下：

### `GET /api/v1/instances`

- 成功：`200 + InstanceListResponse`
- 错误：`401`、`403`
- Query 可选：`kind`、`limit`、`cursor`

### `GET /api/v1/gpu-inventory`

- 成功：`200 + GPUInventoryListResponse`
- 错误：`401`、`403`

### `GET /api/v1/gpu-inventory/occupancy`

- 成功：`200 + GPUOccupancySummary`
- 错误：`401`、`403`

### `GET /api/v1/volumes` / `GET /api/v1/filesystems` / `GET /api/v1/objects`

- 成功：各 `*ListResponse`
- 错误：`401`、`403`

### `GET /api/v1/networks/vpcs`

- 成功：`200 + NetworkVPCListResponse`
- 错误：`401`、`403`

## 待补边界

- 独立 `GET /resource-overview` 聚合 API — **TODO-YAML**
- 存储「已使用量」跨卷/对象/文件系统统一口径 — 待产品冻结
- 实例 `kind=batch_job` 等 list filter 缺口 — 见 `batch-job-instances.md`

## 与平台概览的关系

- 首页嵌入摘要；本页为「查看详情」目标之一
- 两页字段定义必须同源，不得出现两套口径

## 验收标准

- [ ] 未声明独立 overview API
- [ ] GPU 数据引用 `gpu-inventory`，非 deprecated Services 路径
- [ ] 空态与「真实 0」区分
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
