# SPEC: Console 云主机 VM

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-vm-management.md`
> Revised: 2026-06-15

## 1. Summary

### 1.1 What This SPEC Covers

本 SPEC 用于收口 `Console / 算力与云资源 / 实例管理 / 云主机 VM` 的正式路径、字段映射、错误口径、跨模块边界和回填前置依赖。

本 SPEC 只对齐当前 `Core v1.yaml` 已冻结能力，不重新设计新的 VM API。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/compute/prd-console-vm-management.md`
- User Stories covered: `US-001 ~ US-007`
- Functional Requirements covered: `FR-1 ~ FR-12`

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | RBAC Scope |
|---|---|---|---|---|
| GET | `/api/v1/instances?kind=vm` | `listInstances` | `200 + InstanceListResponse` | `scope:instances:read` |
| GET | `/api/v1/instances/{instance_id}` | `getInstance` | `200 + InstanceRecord` | `scope:instances:read` |
| POST | `/api/v1/instances` | `createInstance` | `201 + CreateInstanceResponse` | `scope:instances:create` |
| POST | `/api/v1/instances/{instance_id}/lifecycle` | `applyInstanceLifecycle` | `200 + InstanceLifecycleResponse` | `scope:instances:update` |
| POST | `/api/v1/instances/{instance_id}/console` | `createInstanceConsoleSession` | `200 + InstanceConsoleSession` | `scope:instances:console` |
| GET | `/api/v1/instances/{instance_id}/operations` | `listInstanceOperations` | `200 + InstanceOperation[]` | `scope:instances:read` |
| GET | `/api/v1/instance-operations/{operation_id}` | `getInstanceOperation` | `200 + InstanceOperation` | `scope:instances:read` |

### 2.3 Verified Schemas

- `InstanceRecord`
- `InstanceListResponse`
- `CreateInstanceRequest`
- `CreateInstanceResponse`
- `InstanceLifecycleRequest`
- `InstanceLifecycleResponse`
- `CreateInstanceConsoleSessionRequest`
- `InstanceConsoleSession`
- `InstanceOperation`

### 2.4 创建前置条件与状态矩阵

见主维护源 `docs/console-modules/compute/vm-management.md`（§创建前置条件、§操作可用性矩阵）。`422` 的 `code` 口径见 `../governance/module-delivery-workflow.md` §2.10。

## 3. Page Scope

### 3.1 Page Responsibilities

- 展示 `kind=vm` 的实例列表和详情
- 承接创建实例、生命周期动作和短期 Console/VNC 会话
- 承接实例操作历史与失败原因查看
- 保留网络、存储、镜像、快照等关联资源的摘要和跳转

### 3.2 Non-Goals

- 不重写网络、存储、镜像、安全组等独立资源契约
- 不改变现有 `Core` 已冻结路径和返回码
- 不把 VM 动作改写成 `202 + AsyncTask`

## 4. Data Model Mapping

### 4.1 Page View to Core Schema

| 页面视图 | Core Schema | 说明 |
|---|---|---|
| VM 列表行 | `InstanceRecord` | 使用 `kind=vm` 过滤 |
| VM 详情 | `InstanceRecord` | 展示 VM 相关字段子集 |
| 创建结果 | `CreateInstanceResponse` | 包含 `instance` 与 `operation_id` |
| 生命周期动作结果 | `InstanceLifecycleResponse` | 返回更新后的实例状态和 `operation_id` |
| Console 会话 | `InstanceConsoleSession` | 返回短期会话信息 |
| 操作历史 | `InstanceOperation` | 支持列表与详情 |

### 4.2 Key Field Mapping

| 页面字段 | 权威源字段 | 说明 |
|---|---|---|
| `name` | `InstanceRecord.name` | 实例名称 |
| `state` | `InstanceRecord.state` | 当前状态 |
| `state_reason` | `InstanceRecord.state_reason` | 机器可读原因码 |
| `state_message` | `InstanceRecord.state_message` | 人类可读状态描述 |
| `provider` | `InstanceRecord.provider` | Provider 名称 |
| `node_name` | `InstanceRecord.node_name` | 节点信息 |
| `termination_protection` | `InstanceRecord.termination_protection` | 危险操作保护开关 |
| `ssh.*` | `InstanceRecord.ssh.*` | 仅连接元数据，不返回私钥 |
| `volumes[]` | `InstanceRecord.volumes[]` | 卷摘要，不等于独立卷契约 |
| `snapshots[]` | `InstanceRecord.snapshots[]` | 快照摘要，不等于独立快照模块 |

## 5. API Constraints

### 5.1 General Rules

- `security: [{BearerAuth: []}, {ApiKeyAuth: []}]`
- 所有请求支持 `X-Request-Id`
- 不要求前端传 `X-Tenant-Id`
- 租户边界由认证上下文或中间件获取
- 页面层可用 camelCase，但 API 请求和响应以 snake_case 为准

### 5.2 List and Detail

#### `GET /api/v1/instances?kind=vm`

- `kind` 固定传 `vm`
- 支持 `limit`、`cursor`
- 成功响应：`200 + InstanceListResponse`
- 错误分支：`401`、`403`

#### `GET /api/v1/instances/{instance_id}`

- 路径参数：`instance_id`
- 成功响应：`200 + InstanceRecord`
- 错误分支：`401`、`403`、`404`

### 5.3 Create VM

| 项 | 值 |
|---|---|
| operationId | `createInstance` |
| requestBody.required | `name`、`kind`、`idempotency_key` |
| 关键字段 | `kind=vm`、`cpu`、`memory`、`boot_image`、`ssh_username`、`ssh_key_ref`、`termination_protection` |
| success | `201 + CreateInstanceResponse` |
| error responses | `400`、`401`、`403`、`409` |

示例请求：

```json
{
  "idempotency_key": "vmcreate-20260613-001",
  "name": "vm-demo-01",
  "kind": "vm",
  "cpu": "4",
  "memory": "8Gi",
  "auto_start": true,
  "boot_image": "images/ubuntu-22.04.qcow2",
  "ssh_username": "ubuntu",
  "ssh_key_ref": "sshkey-001",
  "termination_protection": true
}
```

### 5.4 Lifecycle Actions

| 项 | 值 |
|---|---|
| operationId | `applyInstanceLifecycle` |
| requestBody.required | `action`、`idempotency_key` |
| success | `200 + InstanceLifecycleResponse` |
| error responses | `400`、`401`、`403`、`404`、`409` |

动作枚举：

- `start`
- `stop`
- `restart`
- `resize`
- `rebuild`
- `delete`
- `snapshot`
- `attach_volume`
- `detach_volume`
- `rollback`

约束：

- `stop`、`delete`、`rebuild` 等危险动作受 `termination_protection` 影响
- `resize` 可带 `cpu`、`memory`
- `snapshot` 可带 `snapshot_name`
- `attach_volume / detach_volume` 可带 `volume_id`
- `rollback` 可带 `revision`

### 5.5 Console / VNC Session

| 项 | 值 |
|---|---|
| operationId | `createInstanceConsoleSession` |
| requestBody | 可选；`protocol=console/vnc/novnc/serial` |
| success | `200 + InstanceConsoleSession` |
| error responses | `400`、`401`、`403`、`404` |

页面规则：

- 默认仅对 `running` 状态实例开放申请入口
- 只消费短期会话信息，不暴露 provider 长期凭据

### 5.6 Operation History

| 接口 | 用途 |
|---|---|
| `GET /api/v1/instances/{instance_id}/operations` | 查询指定实例的操作历史 |
| `GET /api/v1/instance-operations/{operation_id}` | 查询单个操作详情与时间线 |

## 6. State and Error Rules

### 6.1 Instance States

- `pending`
- `provisioning`
- `starting`
- `running`
- `stopping`
- `stopped`
- `failed`
- `deleting`
- `deleted`

### 6.2 Action Availability

| Action | stopped | running | failed | deleting | deleted |
|---|---|---|---|---|---|
| start | 可用 | 不可用 | 视策略决定 | 不可用 | 不可用 |
| stop | 不可用 | 可用 | 不可用 | 不可用 | 不可用 |
| restart | 不可用 | 可用 | 视策略决定 | 不可用 | 不可用 |
| resize | 建议关闭态优先 | 视 Core 能力决定 | 不可用 | 不可用 | 不可用 |
| delete | 可用 | 若开启保护则冲突 | 可用 | 不可用 | 不可用 |
| console | 不可用 | 可用 | 不可用 | 不可用 | 不可用 |

### 6.3 Error Format

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-20260613-001"
}
```

重点错误：

| Error Code | HTTP Status | 场景 |
|---|---|---|
| `BAD_REQUEST` | 400 | 字段缺失、枚举非法、参数不合法 |
| `UNAUTHORIZED` | 401 | 登录失效 |
| `FORBIDDEN` | 403 | 无实例读/写/console 权限 |
| `NOT_FOUND` | 404 | 实例或操作记录不存在 |
| `CONFLICT` | 409 | 状态冲突、终端保护、重复动作 |

## 7. Cross-Module Boundary

| 关联域 | 当前在 VM 页的处理方式 | 边界要求 |
|---|---|---|
| 网络 | 只展示关联关系和跳转入口 | 不重写 VPC / 子网 / 安全组契约 |
| 存储 | 只展示卷摘要和挂载关系 | 不重写块存储 / 对象存储 / 文件存储契约 |
| 镜像与快照 | 只展示当前实例相关摘要 | 不扩展成独立镜像或恢复模块 |
| 监控 | 仅保留联动或后续补充空间 | 优先复用 `observability`，不发明新主路径 |

## 8. Risks and Prerequisites

### 8.1 Risks

| 风险 | 影响 | 处理方式 |
|---|---|---|
| 旧材料混入其他资源域完整契约 | 模块边界失真 | 正文只保留关联摘要和跳转 |
| 页面想展示的字段超出现有 `InstanceRecord` | 容易自造契约 | 先按现有 schema 收口 |
| 把 GPU 模块的 `202 + AsyncTask` 套到 VM | 与现有 `Core` 返回码冲突 | 严格维持 `201 / 200 / 200` |

### 8.2 Core Alignment Prerequisites

- 当前 VM 页面主路径已在 `Core v1.yaml` 中存在，后续重点是保持页面正文与现有 schema 对齐
- 若后续需要新增 VM 监控摘要接口，应先确认是否复用 `GET /api/v1/observability/query`
- 若后续要把更多关联资源能力下沉成稳定契约，应优先在对应 `Core` 资源域扩充，而不是继续堆叠在 VM 页面内
