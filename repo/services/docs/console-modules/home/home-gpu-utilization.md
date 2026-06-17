# 首页 — GPU 利用率

## 页面定位

`GPU 利用率` 是 `平台概览` 第二主题区块的明细页，用于展示 GPU 总量、分配与利用率摘要，帮助识别算力闲置或过载。

本页以 **Core GPU 库存与占用** 为主数据源，可选结合 **Core observability** 指标曲线。

## 文档管理规则

- 本文是「GPU 利用率」区块的主维护文档
- Phase 2 已在 Core YAML 声明 `gpu-inventory` 路径；handler 仍为 stub
- 不得使用 Services deprecated `/gpu-containers*` 作为 GPU 事实来源
- 与 `platform-overview.md` 第二区块、`gpu-inventory-ui.md` 口径一致

## Core 层要求

### GPU 设备清单

- `GET /api/v1/gpu-inventory`
- `operationId`: `listGPUInventory`
- Query：`vendor`、`state`（`available` / `allocated` / `unavailable` / `maintenance`）、`limit`、`cursor`
- 响应：`GPUInventoryListResponse`（`GPUInventoryDevice` items）
- RBAC：`scope:gpu-inventory:read`

<!-- ADDED-TO-YAML: GET /api/v1/gpu-inventory (Core v1.yaml, Phase 2 2026-06-17) -->

### GPU 占用摘要

- `GET /api/v1/gpu-inventory/occupancy`
- `operationId`: `getGPUOccupancy`
- 响应：`GPUOccupancySummary`
- RBAC：`scope:gpu-inventory:read`

<!-- ADDED-TO-YAML: GET /api/v1/gpu-inventory/occupancy (Core v1.yaml, Phase 2 2026-06-17) -->

### 利用率曲线（可选）

- `GET /api/v1/observability/query` — PromQL 代理
- RBAC：`scope:observability:read`
- 具体 PromQL 与 label 口径 **待运维模板冻结**；文档不得写死未声明的 metric 名称

## Services 层要求

- 本页不通过 Services 获取 GPU 物理库存
- GPU 容器实例占用见 Core `instances?kind=gpu_container`，不在本页展开 CRUD

## 页面职责

- 展示 GPU 总量、已分配、空闲、异常设备数
- 展示窗口平均利用率（occupancy 或 PromQL）
- 标明统计周期（非逐秒实时）
- 跳转 `GPU 算力管理` 与 `GPU 容器实例`

## 字段级定义

与 `platform-overview.md` §GPU 利用率 一致：

| 字段 | 说明 |
|---|---|
| GPU 总量 | 租户可见 GPU 设备数 |
| 已分配 GPU | `state=allocated` 或 occupancy 摘要 |
| 平均利用率 | 百分比 + 时间窗口 |
| 空闲 GPU | `state=available` 计数 |
| 异常 GPU 数 | `unavailable` / `maintenance` 等 |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| GPU 清单读权限 | `scope:gpu-inventory:read` | `403 FORBIDDEN` |
| 读指标权限（趋势图） | `scope:observability:read` | `403 FORBIDDEN` |
| PromQL 合法 | `query` 非空 | `400 BAD_REQUEST` |

本页无 POST/PUT；不涉及 `idempotency_key`。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员 |
|---|---|---|
| 查看利用率摘要 | ✅ | ✅ |
| 查看趋势图 | ✅（有 observability 权限） | ✅ |
| 跳转 GPU 清单 | ✅ | ✅ |
| 在本页调度 GPU | ❌ | ❌ |

## 接口冻结规则

### `GET /api/v1/gpu-inventory`

- 成功：`200 + GPUInventoryListResponse`
- 错误：`401`、`403`
- 租户边界：仅返回当前租户可见设备

### `GET /api/v1/gpu-inventory/occupancy`

- 成功：`200 + GPUOccupancySummary`
- 错误：`401`、`403`

### `GET /api/v1/observability/query`

- 成功：`200 + ObservabilityQueryResponse`
- 错误：`400`、`401`、`403`

## 待补边界

- GPU 单设备 metrics 独立路径 — YAML 未声明（勿写 `gpu-inventory/{id}/metrics` 为冻结契约）
- PromQL 模板与 label 集合 — 待运维文档冻结
- GPU 分配/回收写操作 — 待 Core 资源域规划

## 与相关模块的关系

- `gpu-management.md`：GPU 管理明细
- `gpu-inventory-ui.md`：设备清单完整列表
- `gpu-container-instance-management.md`：实例级 GPU 使用
- `platform-overview.md`：首页摘要区块

## 验收标准

- [ ] 路径与 Phase 2 Core YAML 一致
- [ ] 首页聚合口径引用 occupancy + inventory，不自造 schema
- [ ] 标注 handler 未实现 ≠ 契约未声明
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
