# SPEC: Console Sandbox 实例

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-sandbox-instance-management.md`
> Generated: 2026-06-13 | Scope: `Console / 算力与云资源 / 实例管理 / Sandbox 实例`

## 1. Summary

本 SPEC 定义 `Console / 算力与云资源 / 实例管理 / Sandbox 实例` 的技术边界、字段映射、接口约束和回填规则。当前 Sandbox 实例必须复用 `Core / Instances` 的统一实例模型，并通过 `kind=sandbox` 区分，不扩写独立 `/sandboxes*` 资源。

## 2. Architecture

### 2.1 System Context

```text
浏览器 Console
    -> Core Instances API (/api/v1/instances*)
        -> Sandbox 实例创建(kind=sandbox)
        -> Sandbox 实例详情
        -> 生命周期动作
        -> 操作历史与操作详情
```

约束说明：

1. 当前 `CreateInstanceRequest.kind` 支持 `sandbox`
2. 当前 `GET /api/v1/instances` 已可通过 `kind=sandbox` 或 `instance_type=sandbox` 进入统一实例列表范围
3. Sandbox 差异主要体现在 `CreateInstanceRequest.sandbox_config` 与 `InstanceRecord.sandbox`

### 2.2 创建前置条件与状态矩阵

见主维护源 `docs/console-modules/compute/sandbox-instance-management.md`。不得引用 deprecated `services/v1.yaml` `/sandboxes*`。

## 3. Data Model

### 3.1 View Model Definitions

```ts
interface CreateSandboxInstanceViewRequest {
  idempotencyKey: string;
  name: string;
  kind: "sandbox";
  image?: string;
  cpu?: string;
  memory?: string;
  sandboxConfig?: {
    runtimeClass?: string;
    sessionTimeout?: string;
    networkEgressPolicy?: "deny_all" | "allowlist" | "internet";
  };
}

interface SandboxInstanceSummary {
  id: string;
  name: string;
  status: string;
  runtimeClass?: string;
  sessionTimeout?: string;
  networkEgressPolicy?: string;
  sessionState?: string;
  createdAt?: string;
}
```

### 3.2 Core Schema Mapping

| 页面视图 | Core Schema | 说明 |
|---|---|---|
| Sandbox 创建结果 | `CreateInstanceResponse` | 返回 `instance + operation_id` |
| Sandbox 详情 | `InstanceRecord` | 重点读取 `sandbox` 摘要字段 |
| Sandbox 生命周期结果 | `InstanceLifecycleResponse` | 返回 `instance + operation_id` |
| Sandbox 操作历史 | `CursorPage + items: InstanceOperation[]` | 通过实例 ID 读取 |

## 4. API Design

### 4.1 Endpoints

| Method | Path | operationId | tags | summary | Success |
|---|---|---|---|---|---|
| GET | `/api/v1/instances?kind=sandbox` | `listInstances` | `Instances` | 查询 Sandbox 实例列表（统一实例过滤） | `200 + InstanceListResponse` |
| POST | `/api/v1/instances` | `createInstance` | `Instances` | 创建 Sandbox 实例 | `201 + CreateInstanceResponse` |
| GET | `/api/v1/instances/{instance_id}` | `getInstance` | `Instances` | 查询 Sandbox 实例详情 | `200 + InstanceRecord` |
| POST | `/api/v1/instances/{instance_id}/lifecycle` | `applyInstanceLifecycle` | `Instances` | 执行生命周期动作 | `200 + InstanceLifecycleResponse` |
| GET | `/api/v1/instances/{instance_id}/operations` | `listInstanceOperations` | `Instances` | 查询操作历史 | `200 + CursorPage(items: InstanceOperation[])` |
| GET | `/api/v1/instance-operations/{operation_id}` | `getInstanceOperation` | `Instances` | 查询操作详情 | `200 + InstanceOperation` |

### 4.2 Request/Response Rules

- 创建请求体使用 `CreateInstanceRequest`，`kind` 必须是 `sandbox`
- `CreateInstanceRequest.sandbox_config` 引用 `SandboxConfig`
- 生命周期动作必须提供 `idempotency_key`
- 当前页面承接 `start`、`stop`、`restart`、`delete`
- 当前可以把 `GET /api/v1/instances?kind=sandbox` 写成正式查询（`listInstances.kind` 枚举已含 `sandbox`，2026-06-16 对齐）
- 当前**不能**把 `/api/v1/sandboxes*` 写成正式路径（services/v1.yaml 已 deprecated）

### 4.3 Field Mapping Rules

| 页面字段 | Core 字段 | 说明 |
|---|---|---|
| 运行时类 | `InstanceRecord.sandbox.runtime_class` | Sandbox 配置摘要 |
| 会话超时 | `InstanceRecord.sandbox.session_timeout` | Sandbox 配置摘要 |
| 出口策略 | `InstanceRecord.sandbox.network_egress_policy` | Sandbox 配置摘要 |
| 会话状态 | `InstanceRecord.sandbox.session_state` | `pending / running / expired / stopped` |
| 镜像 | `InstanceRecord.image` | 运行时镜像 |
| CPU | `InstanceRecord.cpu` | 实例规格 |
| 内存 | `InstanceRecord.memory` | 实例规格 |

### 4.4 Non-Frozen Capabilities

以下能力当前**不能**写成已对齐 Core 的正式接口：

- `/api/v1/sandboxes`
- Sandbox 模板列表
- Sandbox 安全事件
- Sandbox 监控指标
- Sandbox 安全总览
- Sandbox 暂停 / 恢复 / 延长存活专属动作
- Sandbox 独立终端能力

## 5. Business Logic

### 5.1 Core Flows

#### 创建实例

1. 用户填写 Sandbox 名称、镜像、资源规格和 `sandbox_config`
2. 页面提交 `POST /api/v1/instances`
3. 请求体中传入 `kind=sandbox`
4. 成功返回 `201 + CreateInstanceResponse`

#### 查看详情

1. 页面通过创建结果、最近实例记录或操作历史进入详情
2. 页面请求 `GET /api/v1/instances/{instance_id}`
3. 读取 `sandbox` 摘要字段并展示当前会话状态

#### 生命周期动作

1. 用户从详情页触发 `start`、`stop`、`restart`、`delete`
2. 页面提交 `POST /api/v1/instances/{instance_id}/lifecycle`
3. 成功后刷新详情和操作历史

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

- 创建冲突时允许 `409 CONFLICT`
- 生命周期动作冲突时允许 `409 CONFLICT`
- 查询详情或操作详情时需正确处理 `404 NOT_FOUND`

## 7. Testing Strategy

- 校验所有路径都来自 `Core / Instances`
- 校验未继续引用 `/sandboxes*`
- 校验已正确把 Sandbox 列表收口为统一实例查询，而非独立 Sandbox 资源组
- 校验模板、安全事件、监控和总览均降级为待补边界

## 8. Implementation Plan

1. 生成 PRD 和 SPEC
2. 收口主维护文档
3. 回填两份 HTML 摘要
4. 执行最终 Core 合规复审
