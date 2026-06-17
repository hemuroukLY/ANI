# PRD: Console GPU 容器实例

## 1. Introduction/Overview

`Console / 算力与云资源 / 实例管理 / GPU 容器实例` 用于帮助租户用户查看、创建和管理当前权限范围内的 GPU 容器实例。本模块属于 `Console` 租户侧实例管理页面，当前必须严格复用 `Core / Instances` 的统一实例资源，而不是引入独立 `/gpu-containers` 资源域。

本轮 PRD 严格按现有 `Core v1.yaml` 约束生成，当前已核实的页面边界如下：

- 当前 GPU 容器实例不拥有独立 `/gpu-containers` Core 路径，而是复用 `/api/v1/instances*` 并通过 `kind=gpu_container` 区分
- 当前已冻结的 GPU 容器实例能力包括：列表、详情、创建、生命周期动作、操作历史、操作详情
- `InstanceRecord.container` 与 `InstanceRecord.gpu` 是当前 GPU 容器实例的正式字段来源
- 当前 `Core v1.yaml` 没有冻结 GPU 容器专属 `logs / events / metrics / exec / dashboard / revisions` 路径
- 当前页面属于 `Console` 租户侧实例页，不是 `BOSS` 的平台 GPU 运维页

## 2. Goals

- 让租户用户通过统一实例能力管理 GPU 容器实例
- 明确 GPU 容器实例与普通容器实例共享统一资源模型，但拥有额外 GPU 字段摘要
- 沉淀一套严格对齐 `Core / Instances` 的主维护材料
- 明确旧原型中 `/gpu-containers*`、Dashboard、监控等能力当前都不是已冻结契约

## 3. User Stories

### US-001: 查看 GPU 容器实例列表
**Description:** 作为租户用户，我希望查看当前租户下的 GPU 容器实例列表，以便快速定位需要操作的实例。

**Acceptance Criteria:**
- [ ] 页面通过 `GET /api/v1/instances?kind=gpu_container` 查询列表
- [ ] 列表展示实例名称、状态、镜像、CPU、内存、副本数、GPU 型号/数量、创建时间
- [ ] 页面不要求前端显式传 `tenant_id`

### US-002: 查看 GPU 容器实例详情
**Description:** 作为租户用户，我希望查看 GPU 容器实例详情，以便了解当前部署状态和 GPU 摘要。

**Acceptance Criteria:**
- [ ] 页面通过 `GET /api/v1/instances/{instance_id}` 查询详情
- [ ] 详情至少展示 `container.replicas`、`ready_replicas`、`revision`、`rollout_status`
- [ ] 详情至少展示 `gpu.vendor`、`gpu.model`、`gpu.count`、`gpu.utilization_percent`

### US-003: 创建 GPU 容器实例
**Description:** 作为租户用户，我希望创建新的 GPU 容器实例，以便部署推理或训练前置工作负载。

**Acceptance Criteria:**
- [ ] 创建接口对齐 `POST /api/v1/instances`
- [ ] 请求体使用 `CreateInstanceRequest`，并传入 `kind=gpu_container`
- [ ] 页面承接 `image`、`cpu`、`memory`、`replicas`、`gpu.*`、`idempotency_key`
- [ ] 成功响应为 `201 + CreateInstanceResponse`

### US-004: 执行 GPU 容器实例生命周期动作
**Description:** 作为租户用户，我希望执行启动、停止、重启、删除或回滚，以便完成基础运维。

**Acceptance Criteria:**
- [ ] 生命周期接口对齐 `POST /api/v1/instances/{instance_id}/lifecycle`
- [ ] 页面可承接 `start`、`stop`、`restart`、`delete`、`rollback`
- [ ] 写操作必须带 `idempotency_key`
- [ ] 页面不把 `exec`、`logs`、`metrics` 写成当前已冻结动作

### US-005: 查看 GPU 容器实例操作历史
**Description:** 作为租户用户，我希望查看实例操作历史和单次操作详情，以便排查创建或变更结果。

**Acceptance Criteria:**
- [ ] 历史接口对齐 `GET /api/v1/instances/{instance_id}/operations`
- [ ] 详情接口对齐 `GET /api/v1/instance-operations/{operation_id}`
- [ ] 页面展示操作类型、状态、发起人、失败原因和步骤时间线

## 4. Functional Requirements

- FR-1: 系统必须通过统一实例模型展示 GPU 容器实例
- FR-2: 系统必须支持 GPU 容器实例的列表、详情、创建
- FR-3: 系统必须支持基础生命周期动作
- FR-4: 系统必须支持操作历史与操作详情查看
- FR-5: 系统必须禁止把未冻结的 GPU 容器专属路径写成正式契约

## 5. Non-Goals (Out of Scope)

- 不在本轮实现独立 `/api/v1/gpu-containers*` 资源
- 不在本轮实现 GPU 容器专属日志、事件、监控、终端 / exec
- 不在本轮实现 GPU 资源总览 Dashboard
- 不在本轮实现版本历史独立资源

## 6. Design Considerations

- 页面文案必须强调 GPU 容器实例属于统一实例 `kind=gpu_container`
- 列表与详情都应突出 `gpu` 和 `container` 两组摘要字段
- Dashboard、监控、日志等旧原型内容只能降为待补边界

## 7. Technical Considerations

- 权威路径以 `Core v1.yaml` 的 `Instances` 路径组为准
- 列表使用 `GET /api/v1/instances?kind=gpu_container`
- 创建使用 `POST /api/v1/instances`
- 生命周期使用 `POST /api/v1/instances/{instance_id}/lifecycle`
- 操作历史使用 `GET /api/v1/instances/{instance_id}/operations`
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

## 8. Success Metrics

- 用户能基于统一实例能力完成 GPU 容器实例的查看、创建和基础运维
- 文档不再出现独立 `/gpu-containers*` 正式路径
- HTML 摘要、PRD、SPEC 与主文档保持一致

## 9. Open Questions

- 后续是否要把 GPU 容器监控与事件拆为独立模块
- 后续是否要引入 GPU 容器专属 revision / deployment 资源
