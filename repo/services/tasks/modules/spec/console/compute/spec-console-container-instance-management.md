# SPEC: Console 容器实例

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-container-instance-management.md`
> Generated: 2026-06-13 | Scope: `Console / 算力与云资源 / 实例管理 / 容器实例`

## 1. Summary

### 1.1 What This SPEC Covers

本 SPEC 定义 `Console / 算力与云资源 / 实例管理 / 容器实例` 的技术边界、页面结构、字段映射、接口约束和回填规则。该页面是 `Console` 侧的租户资源管理页面，直接承接 `Core` 中统一实例资源的 `kind=container` 能力，不改变容器实例资源归属，不把容器实例对象错误写入 `Services` 契约。

本 SPEC 的重点不是重新设计一套新的容器平台 API，而是把现有 `Core v1.yaml` 中已经冻结的统一实例契约，收口为可直接维护和对齐的前端需求材料，并把旧 HTML 中过度宽泛的容器平台 IA 降回到待补边界。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/compute/prd-console-container-instance-management.md`
- User Stories covered: `US-001 ~ US-006`
- Functional Requirements covered: `FR-1 ~ FR-11`

### 1.3 Design Decisions Summary

| Decision | Choice | Rationale |
|---|---|---|
| 页面归属 | Console 资源管理页 | 属于租户侧自助使用页面 |
| 资源归属 | 容器实例归 Core | 当前正式契约通过 `/instances` 统一承接 |
| 资源区分方式 | `kind=container` | 当前无独立 `/containers` 路径 |
| 接口前缀 | `/api/v1/*` | 现有 `Core v1.yaml` 已冻结 `instances` 路径组 |
| 创建返回码 | `201 Created` | 严格对齐现有 Core 契约 |
| 生命周期返回码 | `200 OK` | 严格对齐现有 Core 契约 |
| 页面范围 | 只定义容器实例自身与直接嵌入摘要 | 不把命名空间、Deployment、Service、Ingress 等平台能力写入主契约 |
| Console/Exec 处理 | 不纳入当前模块已冻结能力 | 当前 Core 里无容器实例专属 `console/exec/terminal` 路径 |
| HTML 策略 | 只保留摘要与入口 | 模块详文作为主维护源 |

## 2. Architecture

### 2.1 System Context

```text
浏览器 Console
    -> Core Instance API (/api/v1/instances*)
        -> 容器实例查询
        -> 容器实例创建
        -> 容器实例生命周期动作
        -> 操作历史与操作详情
```

约束说明：

1. 本页默认直接读取 `Core` 的容器实例资源能力，不要求经由 `Services` 再包装资源对象
2. 本页可展示 `InstanceRecord.container`、`resource_refs`、`workload_identity`、卷摘要等实例内嵌信息，但这些信息的底层资源契约仍归各自 `Core` 资源域
3. 本页不重写工作负载、服务发现、配置管理、存储、弹性伸缩、应用市场、微服务治理的完整资源定义，只保留待补说明和模块边界

### 2.2 Component Design

- 页面壳：承载筛选区、统计区、实例表格、详情抽屉、动作区
- 查询层：封装容器实例列表、实例详情、操作历史、操作详情
- 动作层：封装创建实例和生命周期动作
- 权限层：根据实例读/写权限决定是否显示动作入口
- 错误层：统一处理 `ErrorResponse` 并展示 `request_id`

### 2.3 Module Interactions

```text
用户打开容器实例页面
  -> 查询 GET /api/v1/instances?kind=container
  -> 选择某个实例
      -> 查询 GET /api/v1/instances/{instance_id}
      -> 查询 GET /api/v1/instances/{instance_id}/operations
  -> 点击创建容器实例
      -> 提交 POST /api/v1/instances
      -> 返回 201 + CreateInstanceResponse
  -> 点击生命周期动作
      -> 提交 POST /api/v1/instances/{instance_id}/lifecycle
      -> 返回 200 + InstanceLifecycleResponse
  -> 点击某条操作记录
      -> 查询 GET /api/v1/instance-operations/{operation_id}
```

### 2.4 File Structure

```text
prototypes/ani-services-prototype-console.html                        [MODIFY]
prototypes/ani-services-prototype.html                                [MODIFY]
docs/console-modules/README.md                             [MODIFY]
docs/console-modules/compute/container-instance-management.md [NEW]
tasks/modules/prd/console/compute/prd-console-container-instance-management.md         [NEW]
tasks/modules/spec/console/compute/spec-console-container-instance-management.md        [NEW]
```

## 3. Data Model

### 3.1 View Model Definitions

本阶段以页面视图模型和现有 Core schema 映射为主，不新增持久化表设计。

命名约定说明：

- 前端页面视图模型可使用 camelCase
- API 路径参数、请求体和响应 JSON schema 以 `Core v1.yaml` 中的 snake_case 为准
- 如二者有映射差异，模块详文必须明确说明映射关系

```ts
interface ContainerInstanceListItem {
  id: string;
  name: string;
  state: "pending" | "provisioning" | "starting" | "running" | "stopping" | "stopped" | "failed" | "deleting" | "deleted";
  provider: string;
  image?: string;
  cpu?: string;
  memory?: string;
  replicas?: number;
  readyReplicas?: number;
  rolloutStatus?: "pending" | "progressing" | "healthy" | "degraded" | "rolled_back";
  endpoint?: string;
  createdAt: string;
  updatedAt: string;
  operationHint?: string;
}

interface ContainerInstanceDetailView {
  id: string;
  name: string;
  state: string;
  stateReason?: string;
  stateMessage?: string;
  provider: string;
  endpoint?: string;
  nodeName?: string;
  image?: string;
  cpu?: string;
  memory?: string;
  container?: {
    replicas: number;
    readyReplicas: number;
    revision?: string;
    rolloutStatus?: "pending" | "progressing" | "healthy" | "degraded" | "rolled_back";
    history: Array<{
      revision: string;
      image?: string;
      createdAt: string;
    }>;
  };
  resourceRefs: string[];
  workloadIdentity?: {
    keyId?: string;
    keyPrefix?: string;
    scopes: string[];
    active: boolean;
    createdAt?: string;
    revokedAt?: string;
  };
  volumes: Array<{
    name: string;
    kind: "root_disk" | "data_disk" | "shared_pvc" | "object_fuse" | "ephemeral";
    sizeGiB?: number;
    sourceRef?: string;
    mountPath?: string;
    readOnly?: boolean;
  }>;
  createdAt: string;
  updatedAt: string;
}

interface CreateContainerInstanceRequestView {
  name: string;
  image?: string;
  cpu?: string;
  memory?: string;
  replicas?: number;
  autoStart?: boolean;
  idempotencyKey: string;
}

interface ContainerLifecycleActionView {
  action: "start" | "stop" | "restart" | "delete" | "rollback";
  idempotencyKey: string;
  revision?: string;
}

interface ContainerInstanceOperationItem {
  id: string;
  operation: string;
  status: "accepted" | "in_progress" | "succeeded" | "failed" | "cancelled";
  requestedBy: string;
  idempotencyKey?: string;
  failureReason?: string;
  failureMessage?: string;
  retryEligible: boolean;
  createdAt: string;
  updatedAt: string;
}
```

### 3.2 Core Schema Mapping

| 页面视图 | Core Schema | 说明 |
|---|---|---|
| 容器实例列表行 | `InstanceRecord` | 使用 `kind=container` 过滤 |
| 容器实例详情 | `InstanceRecord` | 展示 container 相关字段子集 |
| 创建结果 | `CreateInstanceResponse` | 包含 `instance` 与 `operation_id` |
| 生命周期动作结果 | `InstanceLifecycleResponse` | 返回更新后的实例状态和 `operation_id` |
| 操作历史 | `InstanceOperation` | 支持列表和详情 |

### 3.3 Relationships

- `ContainerInstanceListItem` 与 `ContainerInstanceDetailView` 都映射自 `InstanceRecord`
- `ContainerInstanceDetailView.container` 由 `InstanceRecord.container` 提供，用于展示副本、修订和 rollout 摘要
- `ContainerInstanceDetailView.workloadIdentity` 由 `InstanceRecord.workload_identity` 提供，只展示绑定摘要，不展示密钥明文
- `ContainerInstanceDetailView.volumes` 是实例详情中的嵌入摘要，不代表本页定义了独立卷资源契约

### 3.4 Migration Plan

- 本 SPEC 不要求新增数据库表
- 本模块优先对齐现有 `Core v1.yaml`，不主动新增新的容器实例主路径
- 如后续需要容器实例专属的 logs/metrics/exec/terminal 等能力，应先扩充 `Core v1.yaml`
- 如后续需要更完整的工作负载、服务发现、配置、存储、弹性伸缩等契约，应在对应模块单独展开，不在本页重复定义

## 4. API Design

### 4.1 Endpoints

| Method | Path | operationId | tags | summary | Auth | Success |
|---|---|---|---|---|---|---|
| GET | `/api/v1/instances?kind=container` | `listInstances` | `Instances` | 查询实例列表 | Bearer/ApiKey | `200 + InstanceListResponse` |
| GET | `/api/v1/instances/{instance_id}` | `getInstance` | `Instances` | 查询实例详情 | Bearer/ApiKey | `200 + InstanceRecord` |
| POST | `/api/v1/instances` | `createInstance` | `Instances` | 创建实例 | Bearer/ApiKey | `201 + CreateInstanceResponse` |
| POST | `/api/v1/instances/{instance_id}/lifecycle` | `applyInstanceLifecycle` | `Instances` | 执行实例生命周期操作 | Bearer/ApiKey | `200 + InstanceLifecycleResponse` |
| GET | `/api/v1/instances/{instance_id}/operations` | `listInstanceOperations` | `Instances` | 查询实例操作历史 | Bearer/ApiKey | `200 + allOf(CursorPage + items: InstanceOperation[])` |
| GET | `/api/v1/instance-operations/{operation_id}` | `getInstanceOperation` | `Instances` | 查询单个操作详情 | Bearer/ApiKey | `200 + InstanceOperation` |

### 4.2 Request/Response Schemas

#### 通用要求

- `security: [{BearerAuth: []}, {ApiKeyAuth: []}]`
- 所有请求支持 `X-Request-Id`
- 不要求前端传 `X-Tenant-Id`
- 租户边界必须由认证上下文或中间件获取
- 仅在页面层映射字段，不重定义 Core 原有 schema

#### 容器实例列表查询

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| kind | string | 是 | 固定传 `container` |
| limit | integer | 否 | 分页大小 |
| cursor | string | 否 | Cursor 分页游标 |

成功响应：

- `200 + InstanceListResponse`

错误分支：

- `401 UNAUTHORIZED`
- `403 FORBIDDEN`

#### 容器实例详情查询

路径参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| instance_id | string | 是 | 实例 ID |

成功响应：

- `200 + InstanceRecord`

错误分支：

- `401 UNAUTHORIZED`
- `403 FORBIDDEN`
- `404 NOT_FOUND`

#### 创建容器实例

| 项 | 值 |
|---|---|
| operationId | `createInstance` |
| summary | `创建实例` |
| tags | `["Instances"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| requestBody.required | `name`、`kind`、`idempotency_key` |
| 关键字段 | `kind=container`、`image`、`cpu`、`memory`、`replicas`、`auto_start` |
| success | `201 Created + CreateInstanceResponse` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN`、`409 CONFLICT` |

示例请求：

```json
{
  "idempotency_key": "containercreate-20260613-001",
  "name": "container-demo-01",
  "kind": "container",
  "image": "nginx:1.27",
  "cpu": "2",
  "memory": "4Gi",
  "replicas": 2,
  "auto_start": true
}
```

#### 生命周期动作

| 项 | 值 |
|---|---|
| operationId | `applyInstanceLifecycle` |
| summary | `执行实例生命周期操作` |
| tags | `["Instances"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| path params | `instance_id` |
| requestBody.required | `action`、`idempotency_key` |
| success | `200 OK + InstanceLifecycleResponse` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND`、`409 CONFLICT` |

页面收口规则：

- 本轮作为容器实例页面主能力承接：`start`、`stop`、`restart`、`delete`、`rollback`
- `rollback` 可使用 `revision`
- 当前不把 `snapshot`、`attach_volume`、`detach_volume`、`exec`、`terminal` 写成容器实例页面已冻结主动作

示例请求：

```json
{
  "action": "rollback",
  "idempotency_key": "containerlifecycle-20260613-001",
  "revision": "rev-20260610-03"
}
```

#### 操作历史与详情

| 接口 | 说明 |
|---|---|
| `GET /api/v1/instances/{instance_id}/operations` | 查询指定实例的操作历史 |
| `GET /api/v1/instance-operations/{operation_id}` | 查询单个操作详情与时间线 |

列表接口支持：

- `limit`
- `cursor`

### 4.3 Non-Frozen Capabilities

以下能力当前**不能**写成已对齐 Core 的正式接口：

- 容器实例日志
- 容器实例事件
- 容器实例指标
- 容器实例终端 / exec
- 容器实例 Dashboard
- 命名空间
- Deployment / StatefulSet / DaemonSet / Job / CronJob
- Service / Ingress
- ConfigMap / Secret
- PVC / PV / StorageClass / 卷快照
- HPA / VPA / KEDA
- Helm / DevOps / 微服务治理

### 4.4 Operation Availability Matrix

状态×操作矩阵与 `## 创建前置条件` 以主维护源为准：

- `docs/console-modules/compute/container-instance-management.md`

摘要：`InstanceRecord.state` × `start/stop/restart/resize/delete/rollback`；不允许时 `422`（YAML 已举例 `INSTANCE_STATE_INVALID`）。

### 4.5 Conflict Handling

- `POST /api/v1/instances` 当前已冻结返回码包含 `409 CONFLICT`
- `POST /api/v1/instances/{instance_id}/lifecycle` 当前已冻结返回码包含 `409 CONFLICT`
- 容器实例页面应把 `409` 解释为状态冲突、重复创建冲突或目标 revision/action 当前不可执行
- 本轮不新增额外自定义冲突错误码

## 5. Business Logic

### 5.1 Core Flows

#### 列表加载

1. 页面打开时请求 `GET /api/v1/instances?kind=container`
2. 将 `InstanceRecord` 映射为列表行
3. 若用户点击某一行，加载实例详情和操作历史

#### 创建容器实例

1. 用户填写容器实例基础配置
2. 前端生成 `idempotency_key`
3. 提交 `POST /api/v1/instances`
4. 成功后返回 `instance + operation_id`
5. 页面展示创建中的实例并允许查看对应操作历史

#### 生命周期动作

1. 用户在实例行或详情页选择动作
2. 对危险动作先展示确认提示
3. 前端生成新的 `idempotency_key`
4. 提交 `POST /api/v1/instances/{instance_id}/lifecycle`
5. 根据返回的 `instance` 与 `operation_id` 刷新详情与操作历史

#### 操作排障

1. 用户打开实例操作历史
2. 查看最近一次失败或进行中的操作
3. 点击单条操作查看时间线和失败原因

### 5.2 Validation Rules

- 容器实例创建必须提供 `name`、`kind=container` 和 `idempotency_key`
- `replicas` 若填写，必须满足最小值 `1`
- 页面不接受 VM 专属的 `boot_image`、`ssh_username`、`ssh_key_ref` 作为容器实例首要字段
- 前端不可要求用户主动填写 `tenant_id`
- 页面不得继续使用旧 `/api/v1/console/*`、旧 `/api/v1/storage/*` 或 K8s 原生 API 作为正式契约

### 5.3 Container State & Rollout Handling

- 实例主状态统一使用 `InstanceRecord.state`
- 发布健康度与副本状态使用 `InstanceRecord.container.rollout_status`、`replicas`、`ready_replicas`
- 修订历史使用 `InstanceRecord.container.history[]`
- 不在本轮引入独立容器发布系统状态机

### 5.4 Edge Cases

- 容器实例详情查询时目标实例不存在
- 创建时 `image` 缺失或字段不合法
- 生命周期动作在当前状态下不允许执行
- `rollback` 指定的 revision 不存在
- 用户期望使用 `exec/logs`，但当前 Core 无正式接口

## 6. Error Handling

### 6.1 Error Taxonomy

| Error Code | HTTP Status | Condition | User Message |
|---|---|---|---|
| `BAD_REQUEST` | 400 | 请求字段不合法 | 请求参数不正确，请检查后重试 |
| `UNAUTHORIZED` | 401 | 未登录或凭证失效 | 登录已失效，请重新登录 |
| `FORBIDDEN` | 403 | 无实例访问或操作权限 | 当前账号没有该操作权限 |
| `NOT_FOUND` | 404 | 实例或操作记录不存在 | 目标资源不存在或已删除 |
| `CONFLICT` | 409 | 状态冲突、重复创建、回滚目标不可用 | 当前状态不允许执行该操作 |

### 6.2 Shared Error Format

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-20260613-001"
}
```

### 6.3 Failure Modes

- 容器实例页面不额外声明不存在的 `console/exec/logs` 错误返回
- 查询失败时仅影响对应列表区、详情区或历史区，不阻断整个页面框架
- 对待补能力只展示说明和边界，不伪造失败调用路径

## 7. Security

### 7.1 Authentication & Authorization

- 全部容器实例接口使用 `security: [{BearerAuth: []}, {ApiKeyAuth: []}]`
- 后端必须从认证上下文获取租户边界
- 前端不可信任也不显式传 `tenant_id / X-Tenant-Id`

### 7.2 Input Validation

- `idempotency_key` 必须非空
- `name` 必须符合长度要求
- `replicas` 必须符合 schema 约束
- `revision` 仅在 `rollback` 时填写

### 7.3 Data Protection

- 不展示平台侧全局容器集群运营信息
- 只展示当前租户可见的容器实例与嵌入摘要
- `workload_identity` 只展示绑定摘要，不展示密钥明文

## 8. Performance

### 8.1 Expected Load

- 列表以租户级分页查询为主
- 详情和操作历史按需加载

### 8.2 Optimization Strategy

- 列表统一沿用 `limit + cursor`
- 详情和操作历史按需加载
- 不将待补能力扩写为额外聚合接口

### 8.3 Database Considerations

- 本文档不新增数据库设计
- 数据模型以现有 `Core` schema 为准

## 9. Testing Strategy

### 9.1 Unit Tests

- 校验 `kind=container` 查询和创建字段映射
- 校验待补能力不会误显示为已冻结接口

### 9.2 Integration Tests

- 容器实例列表、详情、创建、生命周期、操作历史接口映射正确
- 创建请求带有 `idempotency_key`
- 页面不会把 `console/exec/logs/dashboard` 写成既有冻结契约
- 页面不会把 K8s 工作负载 IA 误写成容器实例正式路径

### 9.3 Edge Case Tests

- 查询不存在实例时展示 `NOT_FOUND`
- 创建参数非法时展示 `BAD_REQUEST`
- 生命周期动作冲突时展示 `CONFLICT`

### 9.4 Acceptance Criteria Mapping

| US/FR | Test | Type | Description |
|---|---|---|---|
| `US-001` | 容器实例列表契约校验 | integration | 验证 `/instances?kind=container` 的路径、字段和返回码对齐 |
| `US-002` | 容器实例详情契约校验 | integration | 验证 `InstanceRecord.container`、`workload_identity` 和 `resource_refs` 映射 |
| `US-003` | 容器实例创建契约校验 | integration | 验证 `CreateInstanceRequest` 必填字段、返回结构与 `201` |
| `US-004` | 生命周期动作契约校验 | integration | 验证 `lifecycle` 路径、`action/idempotency_key` 与 `409` 冲突语义 |
| `US-005` | 操作历史契约校验 | integration | 验证操作列表与单条详情路径、分页和返回码 |
| `US-006` | Core 边界校验 | unit | 验证未冻结容器平台能力不会写入正式接口表 |

## 10. Implementation Plan

### 10.1 Phases

1. 明确容器实例已冻结与未冻结边界
2. 生成 PRD 和 SPEC
3. 收口主维护源
4. 回填 HTML 摘要
5. 执行最终 Core 合规复审

### 10.2 Issue Mapping

| Issue | SPEC Sections | Priority | Depends On |
|---|---|---|---|
| 容器实例边界收口 | 1, 2, 4.3 | high | — |
| Core 已冻结实例接口整理 | 3, 4.1, 4.2 | high | 容器实例边界收口 |
| 主维护源落盘 | 5, 6, 7 | high | Core 已冻结实例接口整理 |
| HTML 摘要回填 | 2.4, 10.1 | medium | 主维护源落盘 |

### 10.3 Incremental Delivery

- 本轮先收口主维护材料，不直接扩写新的 `Core v1.yaml`
- 后续若要补容器实例专属 logs/metrics/exec/terminal 等能力，再进入 `YAML` 扩充阶段

## 11. Open Questions & Risks

### 11.1 Unresolved Questions

- 容器实例日志、指标、事件、终端、exec 是否作为下一轮 Core 扩充能力
- 副本调整后续是继续复用 `lifecycle`，还是补充更清晰的字段或动作
- 容器平台 Dashboard 是否应拆为独立模块

### 11.2 Technical Risks

| Risk | Impact | Mitigation |
|---|---|---|
| 旧 HTML 继续沿用 K8s 平台 IA | 会造成容器实例模块严重超写 Core 事实 | 统一降回待补能力，只保留当前 `/instances` 冻结能力 |
| 误把 VM Console 接口套给容器实例 | 会造成错误契约和错误页面动作 | 明确 `POST /instances/{instance_id}/console` 仅用于 VM |
| 误把 Deployment/Service/Ingress/Config 资源写入本页 | 会导致边界失控 | 在正文和 SPEC 中单列非目标和待补能力 |

### 11.3 Frozen Baseline

- 当前 `Core v1.yaml` 是容器实例资源的唯一权威来源
- 页面不新增独立 `Services` 容器实例聚合接口
- 本轮只收口 `Console / 容器实例` 子模块，不同步细化完整容器平台管理台
