# SPEC: Console GPU 容器实例

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-gpu-container-instance-management.md`
> Generated: 2026-06-13 | Scope: `Console / 算力与云资源 / 实例管理 / GPU 容器实例`

## 1. Summary

本 SPEC 定义 `Console / 算力与云资源 / 实例管理 / GPU 容器实例` 的技术边界、字段映射、接口约束和回填规则。当前 GPU 容器实例必须复用 `Core / Instances` 的统一实例模型，并通过 `kind=gpu_container` 区分，不扩写独立 `/gpu-containers*` 资源。

## 2. Architecture

### 2.1 System Context

```text
浏览器 Console
    -> Core Instances API (/api/v1/instances*)
        -> GPU 容器实例列表(kind=gpu_container)
        -> GPU 容器实例详情
        -> GPU 容器实例创建
        -> 生命周期动作
        -> 操作历史与操作详情
```

约束说明：

1. GPU 容器实例与普通容器实例共享统一实例资源
2. 当前 GPU 差异主要体现在 `InstanceRecord.gpu` 摘要字段
3. 监控、日志、事件、exec、Dashboard 都不是当前已冻结路径

### 2.2 创建前置条件与状态矩阵

见主维护源 `docs/console-modules/compute/gpu-container-instance-management.md`。不得引用 deprecated `services/v1.yaml` `/gpu-containers*`。

## 3. Data Model

### 3.1 View Model Definitions

```ts
interface GPUContainerInstanceListItem {
  id: string;
  name: string;
  status: string;
  image?: string;
  cpu?: string;
  memory?: string;
  replicas?: number;
  readyReplicas?: number;
  revision?: string;
  rolloutStatus?: string;
  gpuVendor?: string;
  gpuModel?: string;
  gpuCount?: number;
  gpuUtilizationPercent?: number;
  createdAt?: string;
}

interface CreateGPUContainerInstanceViewRequest {
  idempotencyKey: string;
  name: string;
  kind: "gpu_container";
  image?: string;
  cpu?: string;
  memory?: string;
  replicas?: number;
  gpu?: {
    vendor?: string;
    model?: string;
    count?: number;
  };
}
```

### 3.2 Core Schema Mapping

| 页面视图 | Core Schema | 说明 |
|---|---|---|
| GPU 容器实例列表 | `InstanceListResponse.items[] -> InstanceRecord` | 列表查询需加 `kind=gpu_container` |
| GPU 容器实例详情 | `InstanceRecord` | 重点读取 `container` 与 `gpu` |
| 创建结果 | `CreateInstanceResponse` | 返回 `instance + operation_id` |
| 生命周期结果 | `InstanceLifecycleResponse` | 返回 `instance + operation_id` |
| 操作历史 | `CursorPage + items: InstanceOperation[]` | 操作时间线 |

## 4. API Design

### 4.1 Endpoints

| Method | Path | operationId | tags | summary | Success |
|---|---|---|---|---|---|
| GET | `/api/v1/instances?kind=gpu_container` | `listInstances` | `Instances` | 查询 GPU 容器实例列表 | `200 + InstanceListResponse` |
| GET | `/api/v1/instances/{instance_id}` | `getInstance` | `Instances` | 查询 GPU 容器实例详情 | `200 + InstanceRecord` |
| POST | `/api/v1/instances` | `createInstance` | `Instances` | 创建 GPU 容器实例 | `201 + CreateInstanceResponse` |
| POST | `/api/v1/instances/{instance_id}/lifecycle` | `applyInstanceLifecycle` | `Instances` | 执行实例生命周期动作 | `200 + InstanceLifecycleResponse` |
| GET | `/api/v1/instances/{instance_id}/operations` | `listInstanceOperations` | `Instances` | 查询操作历史 | `200 + CursorPage(items: InstanceOperation[])` |
| GET | `/api/v1/instance-operations/{operation_id}` | `getInstanceOperation` | `Instances` | 查询操作详情 | `200 + InstanceOperation` |

### 4.2 Request/Response Rules

- 列表查询必须带 `kind=gpu_container`
- 创建请求体使用 `CreateInstanceRequest`，`kind` 必须是 `gpu_container`
- 创建写操作必须提供 `idempotency_key`
- 生命周期动作当前页面只承接 `start`、`stop`、`restart`、`delete`、`rollback`
- 生命周期动作必须提供 `idempotency_key`

### 4.3 Field Mapping Rules

| 页面字段 | Core 字段 | 说明 |
|---|---|---|
| 状态 | `InstanceRecord.status` | 统一实例状态 |
| 镜像 | `InstanceRecord.image` | 运行时镜像引用 |
| 副本数 | `InstanceRecord.container.replicas` | GPU 容器副本数 |
| 就绪副本数 | `InstanceRecord.container.ready_replicas` | GPU 容器就绪摘要 |
| Revision | `InstanceRecord.container.revision` | 当前修订版本 |
| 发布状态 | `InstanceRecord.container.rollout_status` | 发布摘要 |
| GPU 厂商 | `InstanceRecord.gpu.vendor` | GPU 资源摘要 |
| GPU 型号 | `InstanceRecord.gpu.model` | GPU 资源摘要 |
| GPU 数量 | `InstanceRecord.gpu.count` | GPU 资源摘要 |
| GPU 利用率 | `InstanceRecord.gpu.utilization_percent` | GPU 状态摘要 |

### 4.4 Non-Frozen Capabilities

以下能力当前**不能**写成已对齐 Core 的正式接口：

- `/api/v1/gpu-containers`
- GPU 容器日志
- GPU 容器事件
- GPU 容器监控
- GPU 容器终端 / exec
- GPU 资源总览 Dashboard
- 版本历史独立资源

## 5. Business Logic

### 5.1 Core Flows

#### 列表加载

1. 页面打开时请求 `GET /api/v1/instances?kind=gpu_container`
2. 从统一实例记录中映射 `container` 与 `gpu` 摘要
3. 点击实例进入详情页

#### 创建实例

1. 用户填写镜像、CPU、内存、副本数、GPU 规格
2. 页面提交 `POST /api/v1/instances`
3. 请求体中传入 `kind=gpu_container`
4. 成功返回 `201 + CreateInstanceResponse`

#### 生命周期动作

1. 用户从详情页或列表触发动作
2. 页面提交 `POST /api/v1/instances/{instance_id}/lifecycle`
3. 成功返回 `InstanceLifecycleResponse`
4. 页面跳转或刷新操作历史

## 6. Error Handling

### 6.1 Shared Error Format

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-20260613-001"
}
```

### 6.2 Failure Modes

- 列表失败时仅影响列表区，不阻断页面骨架
- 创建冲突时允许 `409 CONFLICT`
- 生命周期动作冲突时允许 `409 CONFLICT`
- 查询详情或操作详情时需正确处理 `404 NOT_FOUND`

## 7. Testing Strategy

- 校验所有路径都来自 `Core / Instances`
- 校验列表查询必须使用 `kind=gpu_container`
- 校验未继续引用 `/gpu-containers*`
- 校验监控、日志、exec、Dashboard 均降级为待补边界

## 8. Implementation Plan

1. 生成 PRD 和 SPEC
2. 收口主维护文档
3. 回填两份 HTML 摘要
4. 执行最终 Core 合规复审
