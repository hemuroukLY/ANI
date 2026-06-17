# 批任务实例

## 页面定位

`批任务` 是统一实例架构下 `kind=batch_job` 的 Console 管理页，用于查看与管理批处理作业实例。

本模块属于 **Core / Instances**，与 VM、容器共享 `/api/v1/instances*` 契约。

## 文档管理规则

- 本文是批任务实例子模块主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 不得使用 deprecated Services 路径

## Core 层要求

- 资源归属 Core `/api/v1/instances*`
- `InstanceRecord.kind` / `instance_type` 枚举含 `batch_job`（schema 已声明）
- **列表 filter 缺口**：`listInstances` 的 `kind` query 枚举当前为 `[vm, container, gpu_container, sandbox]`，**不含 `batch_job`**

<!-- TODO-YAML: listInstances kind 枚举扩展 batch_job；或专用 listBatchJobs -->

### 冻结路径（与统一实例一致）

| 能力 | 路径 | operationId |
|---|---|---|
| 列表 | `GET /api/v1/instances?kind=batch_job`（**待 kind 枚举扩展**） | `listInstances` |
| 详情 | `GET /api/v1/instances/{instance_id}` | `getInstance` |
| 创建 | `POST /api/v1/instances`（body.kind=batch_job） | `createInstance` |
| 生命周期 | `POST /api/v1/instances/{instance_id}/lifecycle` | `applyInstanceLifecycle` |
| 操作历史 | `GET /api/v1/instances/{instance_id}/operations` | `listInstanceOperations` |

## 页面职责

- 展示批任务列表、详情、生命周期（复用 instances lifecycle）
- 与 VM/容器页相同的 operations 子路径模式
- 不提供批任务专属 CRUD 路径（除非 YAML 扩展）

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 读权限 | `scope:instances:read` | `403 FORBIDDEN` |
| 创建权限 | `scope:instances:create` | `403 FORBIDDEN` |
| POST 必填 | `name`、`kind`、`idempotency_key` | `400 BAD_REQUEST` |
| `kind=batch_job` | schema 枚举已含 | create 接受 body |
| 镜像/配额等 | 满足 Core 前置 | `422 PRECONDITION_FAILED`（create 已声明；YAML 已举例 `IMAGE_NOT_FOUND` 等） |

## 操作可用性矩阵

生命周期矩阵与 `vm-management.md` 一致，按 `InstanceRecord.state` 与 `action` 枚举判定。批任务无 Console/VNC 入口（除非产品另行定义）。

| 操作 | 只读用户 | 运维角色 |
|---|---|---|
| 列表/详情 | ✅ | ✅ |
| 创建 | ❌ | ✅ |
| 生命周期动作 | ❌ | ✅ |
| 操作历史 | ✅ | ✅ |

## 接口冻结规则

### `GET /api/v1/instances?kind=batch_job`

- 成功：`200 + InstanceListResponse`
- 错误：`401`、`403`
- **阻塞**：query `kind` 枚举尚未含 `batch_job` — 实现前不得声称 list 可用

### `POST /api/v1/instances`

- 成功：`201 + CreateInstanceResponse`
- 错误：`400`、`401`、`403`、`409`、`422`

### `POST /api/v1/instances/{instance_id}/lifecycle`

- 成功：`200 + InstanceLifecycleResponse`
- 错误：`400`、`401`、`403`、`404`、`409`、`422`（YAML 已举例 `INSTANCE_STATE_INVALID`）

### `GET /api/v1/instances/{instance_id}`

- 成功：`200 + InstanceRecord`
- 错误：`401`、`403`、`404`

## 待补边界

- `listInstances` kind 枚举扩展 `batch_job` — **TODO-YAML**
- 批任务专属字段（job spec、exit code）— 待 `InstanceRecord` schema 扩展
- 批任务日志/事件 — 复用 `container-observability.md` 实例子路径

## 相关模块

- `container-instance-management.md` — 同类统一实例 UI 模式
- `vm-management.md` — lifecycle 矩阵参考

## 验收标准

- [ ] 明确标注 list filter 与 schema enum 不一致
- [ ] 不使用 deprecated Services 路径
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
