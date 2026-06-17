# SPEC: Console GPU 算力管理

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-gpu-management.md`
> Revised: 2026-06-15

## 1. Summary

### 1.1 What This SPEC Covers

本 SPEC 用于收口 `Console / 算力与云资源 / GPU 算力管理` 的权威源事实、页面技术边界、跨模块关系和后续回填前置条件。

本 SPEC 的重点不是设计一套新的 GPU API，而是明确：

- 当前哪些内容已经被权威源证实
- 当前哪些内容仍然只能作为页面目标或待补边界
- 文档如何避免把未来设计误写成当前冻结契约

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/compute/prd-console-gpu-management.md`
- User Stories covered: `US-001 ~ US-005`
- Functional Requirements covered: `FR-1 ~ FR-8`

## 2. Frozen Facts

### 2.1 Authority Sources

- `ANI-main/repo/api/openapi/v1.yaml`
- `ANI-main/repo/api/core-v1-compatibility-baseline.yaml`

### 2.2 Verified Facts

| 项 | 当前结论 | 证据 |
|---|---|---|
| 页面归属 | `Console` 租户侧资源管理页面 | 当前模块体系定义 |
| 资源归属 | `GPU` 资源归 `Core` | `v1.yaml` 顶部资源域说明**规划**了 `gpu-inventory`（路径尚未声明） |
| 统一实例关联 | `InstanceRecord.gpu` 仅用于 `gpu_container` 实例摘要 | `components.schemas.InstanceRecord.gpu` |
| 独立 GPU 路径 | 当前未冻结 | `v1.yaml` 中不存在独立 `GPU` path 段 |
| 独立 GPU schema | 当前未冻结 | `v1.yaml` 中不存在 `GpuOverviewResponse` 等 schema |

### 2.3 Non-Frozen Capabilities

以下能力当前只能写成待补边界，不能写成已冻结契约：

- GPU 总览独立查询
- GPU 型号分布独立查询
- GPU 节点与设备独立查询
- GPU 占用分布独立查询
- GPU 分配 / 回收 / 释放独立写操作

## 3. Page Scope

### 3.1 Page Responsibilities

- 作为租户侧 GPU 资源观察入口，定义页面结构和用户理解路径
- 表达总览、分布、节点、设备、占用关系和异常定位这几个区块
- 表达与 `平台概览`、`VM`、`容器实例`、`GPU 容器实例` 的联动关系
- 明确哪些能力当前仍需 `Core` 补齐正式接口

### 3.2 Non-Goals

- 不定义新的 `Services` GPU 资源契约
- 不定义平台侧 GPU 运营大盘
- 不定义底层调度算法、设备插件或资源池实现
- 不编造临时路径、schema、operationId 或返回码

## 4. Data and Field Boundaries

### 4.1 Page-Level Target Fields

以下字段是页面目标表达，不代表当前已经存在同名冻结 schema：

| 页面区块 | 目标字段 | 当前状态 |
|---|---|---|
| 资源总览 | `gpu_total`、`gpu_allocated`、`gpu_idle`、`avg_utilization`、`abnormal_device_count`、`refreshed_at` | 待 Core 补齐独立资源接口 |
| 型号分布 | `vendor`、`model`、`gpu_total`、`gpu_allocated`、`gpu_idle` | 待 Core 补齐独立资源接口 |
| 节点列表 | `node_name`、`model`、`gpu_count`、`status`、`utilization` | 待 Core 补齐独立资源接口 |
| 设备列表 | `device_id`、`status`、`alloc_state`、`owner_type`、`owner_id` | 待 Core 补齐独立资源接口 |
| 占用分布 | `scope_type`、`scope_id`、`scope_name`、`gpu_allocated` | 待 Core 补齐独立资源接口 |

### 4.2 Existing Safe References

当前可安全引用、但不等于 GPU 页面独立契约的既有能力：

| 类型 | 路径/Schema | 用途 |
|---|---|---|
| Core schema | `InstanceRecord.gpu` | 仅用于 `gpu_container` 调度和利用率摘要 |
| Core schema | `InstanceRecord.state_reason` | 可承接如 `InsufficientGPU` 之类实例状态原因 |
| Core 路径 | `/api/v1/instances` | 可作为关联实例跳转的事实来源 |
| Core 路径 | `/api/v1/observability/query` | 可作为后续监控联动参考，不等于已冻结 GPU 页面接口 |
| Core 路径（规划） | `GET /api/v1/gpu-inventory` | **待 Core 冻结**；仅见于 `v1.yaml` 资源域规划 |

## 5. Interaction Model

### 5.1 Page Structure

```text
GPU 算力管理
├── 资源总览
├── GPU 型号分布
├── GPU 节点列表
├── GPU 设备列表
├── 租户内占用分布
└── 资源动作（待补）
```

### 5.2 Cross-Module Interactions

| 来源区块 | 下一跳 | 当前说明 |
|---|---|---|
| 资源总览 | `GPU 算力管理` 内部各明细区块 | 当前只定义导航关系 |
| 占用分布 | `云主机 VM` / `容器实例` / `GPU 容器实例` | 用于定位占用对象 |
| 平台概览告警摘要 | `GPU 算力管理` | 用于带筛选进入资源排障 |

### 5.3 Action Area Rules

- 页面允许保留 `分配 GPU`、`回收 / 释放 GPU` 的入口位
- 当前入口位只能表达为 `待补能力`
- 在独立 GPU 写接口冻结前，不得写出正式路径、正式 `operationId`、正式响应 schema
- 如未来动作冻结为 `POST` 或其他有副作用方式，必须补充 `idempotency_key`、成功状态码和统一错误结构

## 6. Error and Status Rules

### 6.1 Unified Error Format

统一错误口径保持为：

```json
{"code":"UPPER_SNAKE","message":"...","request_id":"..."}
```

### 6.2 Page State Rules

| 场景 | 页面表达 |
|---|---|
| 当前无 GPU 数据 | 展示空态 |
| 指标未接通 | 展示 `数据暂不可用` |
| 资源异常 | 高亮异常节点、设备或占用关系 |
| 能力未冻结 | 展示 `待 Core 补齐`，不得伪造请求结果 |

## 7. Risks

| 风险 | 影响 | 处理方式 |
|---|---|---|
| 误把页面目标写成已冻结 API | 直接破坏 Core 对齐可信度 | 明确区分 `页面目标字段` 与 `正式 schema` |
| 误把 `gpu_container` 摘要当成 GPU inventory 正式资源 | 模块边界混乱 | 仅把 `InstanceRecord.gpu` 作为关联引用 |
| 误把动作入口写成已实现 | 用户和后续维护者误判能力成熟度 | 统一改写为待补边界 |

## 8. Core Alignment Prerequisites

- 若要把 `资源总览 / 型号分布 / 节点 / 设备 / 占用分布` 写成正式接口，必须先在 `v1.yaml` 中新增对应 path 和 schema
- 若要把 `分配 / 回收 / 释放` 写成正式动作，必须先在 `v1.yaml` 中冻结请求体、错误码和成功响应
- 在权威源补齐前，HTML、主维护文档和任何辅助材料都不得出现伪造的 `Gpu*Response`、`/api/v1/gpu-*` 或 `gpu-allocations*`
