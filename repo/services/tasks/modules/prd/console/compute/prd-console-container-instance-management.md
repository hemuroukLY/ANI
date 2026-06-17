# PRD: Console 容器实例

## 1. Introduction/Overview

`Console / 算力与云资源 / 实例管理 / 容器实例` 用于帮助租户用户查看、创建、操作和排查自己权限范围内的容器实例。该页面属于 `Console` 侧资源管理页面，本质上是对 `Core` 统一实例资源中 `kind=container` 能力的前端承接。

本轮 PRD 严格按现有 `Core v1.yaml` 约束生成，当前已核实的页面边界如下：

- `容器实例` 当前没有独立 `/containers` Core 路径，正式契约复用 `/api/v1/instances` 并通过 `kind=container` 区分
- 当前已冻结的容器实例能力包括：列表、详情、创建、生命周期动作、操作历史、操作详情
- 当前 `Core v1.yaml` 没有冻结容器实例专属的 `logs / events / metrics / terminal / exec / dashboard` 路径
- 当前 `Core v1.yaml` 没有冻结命名空间、Deployment、StatefulSet、DaemonSet、Job、CronJob、Service、Ingress、ConfigMap、Secret、PVC、HPA、Helm 等容器平台资源为本页正式契约
- 容器实例页可以展示 `InstanceRecord.container` 摘要、`workload_identity` 摘要和 `resource_refs`，但不能把这些关联信息扩写成新的独立 Core 资源定义
- 本页属于 `Console` 租户侧资源管理页，不是 `BOSS` 的平台容器集群运营页

## 2. Goals

- 让租户在一个页面内完成容器实例的查询、创建、基础操作和排障
- 严格以现有 `Core v1.yaml` 为准，沉淀容器实例页面可直接用于 Core 对齐的主维护材料
- 明确容器实例页与 K8s 工作负载、服务发现、配置管理、存储、弹性伸缩、DevOps、应用市场等模块的边界
- 让页面文档可以直接映射到 `Core /instances*` 系列接口和操作历史接口
- 保持 `Console` 租户视角，不混入 `BOSS` 平台运营口径

## 3. User Stories

### US-001: 查看容器实例列表
**Description:** 作为租户用户，我希望查看当前租户下的容器实例列表，以便快速找到需要操作的实例。

**Acceptance Criteria:**
- [ ] 页面展示容器实例列表，并支持基于 `kind=container` 的查询视角
- [ ] 列表至少展示名称、状态、镜像、CPU、内存、副本、就绪副本、Provider、创建时间和最近操作线索
- [ ] 列表数据范围限定为当前租户可见实例
- [ ] 查询接口对齐 `GET /api/v1/instances?kind=container`

### US-002: 查看容器实例详情
**Description:** 作为租户用户，我希望查看单个容器实例的详细信息，以便判断当前运行状态、发布状态和关联资源。

**Acceptance Criteria:**
- [ ] 详情页或详情抽屉至少展示基础信息、状态、镜像、CPU、内存、副本、修订版本、rollout 状态、发布历史、访问地址和更新时间
- [ ] 详情中的数据字段对齐 `InstanceRecord`
- [ ] 页面可展示 `resource_refs`、`workload_identity` 和卷摘要，但不把它们写成独立资源契约
- [ ] 页面不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- [ ] 查询接口对齐 `GET /api/v1/instances/{instance_id}`

### US-003: 创建容器实例
**Description:** 作为有权限的租户用户，我希望创建新的容器实例，以便交付通用计算容器环境。

**Acceptance Criteria:**
- [ ] 页面提供创建入口，并要求填写名称、镜像、CPU、内存和副本等必要字段
- [ ] 创建接口对齐 `POST /api/v1/instances`
- [ ] 请求体包含 `name`、`kind=container`、`idempotency_key`
- [ ] 容器实例创建字段只承接 `image`、`cpu`、`memory`、`replicas`、`auto_start`
- [ ] 创建结果按现有 Core 契约返回 `201 Created`
- [ ] 返回结果包含 `instance` 和 `operation_id`
- [ ] 页面明确“创建容器实例”不等于创建完整 K8s 工作负载编排资源

### US-004: 执行容器实例生命周期动作
**Description:** 作为有权限的租户用户，我希望对容器实例执行基础生命周期动作，以便完成日常运维和版本回退。

**Acceptance Criteria:**
- [ ] 页面提供与当前 Core 冻结契约一致的生命周期动作入口
- [ ] 生命周期动作对齐 `POST /api/v1/instances/{instance_id}/lifecycle`
- [ ] 请求体包含 `action` 和 `idempotency_key`
- [ ] 页面可承接 `start`、`stop`、`restart`、`delete`、`rollback` 这组容器实例主动作
- [ ] 页面不把 `exec`、`terminal`、`logs` 写成现有已冻结动作
- [ ] 响应按现有 Core 契约返回 `200 OK + InstanceLifecycleResponse`

### US-005: 查看操作历史与失败原因
**Description:** 作为租户用户，我希望查看容器实例的操作历史和失败原因，以便快速判断问题发生在哪一步。

**Acceptance Criteria:**
- [ ] 页面提供实例操作历史入口
- [ ] 列表接口对齐 `GET /api/v1/instances/{instance_id}/operations`
- [ ] 明细接口对齐 `GET /api/v1/instance-operations/{operation_id}`
- [ ] 历史记录至少展示操作类型、状态、触发人、时间线步骤、失败原因和是否可重试

### US-006: 保持容器实例页面边界清晰
**Description:** 作为 Service 团队成员，我希望容器实例页面正文只定义当前 `Core` 已冻结的统一实例能力，以便后续和其他容器平台模块拆分维护时没有歧义。

**Acceptance Criteria:**
- [ ] 页面正文明确 `容器实例` 资源归 `Core`
- [ ] 页面正文不继续使用旧 `/api/v1/console/*`、旧 `/api/v1/storage/*` 或 K8s 原生 API 作为本页正式契约
- [ ] 页面正文不把 Deployment、StatefulSet、Service、Ingress、ConfigMap、Secret、PVC、HPA、Helm 等资源写成当前模块的主契约
- [ ] 页面正文不把 `logs / exec / terminal / dashboard` 写成当前 Core 已冻结路径
- [ ] 页面正文把工作负载平台能力定义为“待补能力/其他模块”，而不是当前模块的正式接口

## 4. Functional Requirements

- FR-1: 系统必须提供基于 `kind=container` 的容器实例列表查询视图
- FR-2: 系统必须提供容器实例详情视图，并对齐 `InstanceRecord`
- FR-3: 系统必须支持创建容器实例，并对齐 `CreateInstanceRequest / CreateInstanceResponse`
- FR-4: 系统必须支持对容器实例执行生命周期动作，并对齐 `InstanceLifecycleRequest / InstanceLifecycleResponse`
- FR-5: 系统必须提供容器实例操作历史与单个操作详情查询入口
- FR-6: 系统必须从认证上下文获取租户边界，不依赖前端传入 `tenant_id`
- FR-7: 系统必须对创建和生命周期动作应用 `idempotency_key`
- FR-8: 系统必须统一使用标准错误结构
- FR-9: 系统必须明确本页属于 `Console` 租户侧，不属于 `BOSS`
- FR-10: 系统必须明确当前页面只承接统一实例中的 `container` 能力，不承接完整容器平台 IA
- FR-11: 系统必须让最终正文可直接作为容器实例模块对齐 Core 的主维护材料

## 5. Non-Goals (Out of Scope)

- 不在本轮重写命名空间、Deployment、StatefulSet、DaemonSet、Job、CronJob 的完整资源契约
- 不在本轮重写 Service、Ingress、ConfigMap、Secret、PVC、PV、StorageClass、卷快照的完整资源契约
- 不在本轮实现容器实例专属的 `logs`、`events`、`metrics`、`terminal`、`exec`、`dashboard`
- 不在本轮实现 HPA、VPA、KEDA、Helm、DevOps、微服务治理等平台能力
- 不在本轮实现 `BOSS` 侧的平台容器集群运营总览
- 不在本轮改变现有 `Core v1.yaml` 已冻结的返回码和主路径

## 6. Design Considerations

- 页面口径必须是“当前租户 / 当前权限范围内的容器实例”，不能写成全平台 K8s 管理台视角
- 详情区域允许展示 `container.replicas / ready_replicas / revision / rollout_status / history`，但不要把这些字段扩写成独立发布系统契约
- 对于服务发现、配置管理、存储、弹性伸缩、DevOps、应用市场等能力，本页只保留待补说明或后续模块边界
- 危险动作应强调状态冲突和回滚修订版本选择
- 当前页面不提供 VM 专属 Console/VNC 会话能力

## 7. Technical Considerations

- 核心资源路径以 `Core v1.yaml` 现有 `instances` 路径组为准
- 容器实例列表使用 `GET /api/v1/instances?kind=container`
- 容器实例创建使用 `POST /api/v1/instances`，成功返回 `201`
- 生命周期动作使用 `POST /api/v1/instances/{instance_id}/lifecycle`，成功返回 `200`
- 操作历史使用 `GET /api/v1/instances/{instance_id}/operations` 和 `GET /api/v1/instance-operations/{operation_id}`
- 页面不得要求前端显式传 `tenant_id / X-Tenant-Id`
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- 当前页面不接入 `POST /api/v1/instances/{instance_id}/console`，因为该接口明确用于 VM Console/VNC 会话

## 8. Success Metrics

- 用户能在 30 秒内定位目标容器实例并判断其当前运行和发布状态
- 用户能在同一页面内完成容器实例的基础创建、生命周期动作和操作排障
- 文档中关于容器实例的路径、返回码、字段与现有 `Core v1.yaml` 不再冲突
- 页面正文不再混写 K8s 工作负载、配置、存储、Ingress、HPA、Helm 等其他平台模块的完整契约

## 9. Open Questions

- 容器实例专属的 `logs / events / metrics / terminal / exec` 后续是否进入 `Core v1.yaml`
- 容器实例的副本调整是否后续继续复用 `lifecycle`，还是新增更明确的动作或字段
- 若未来要支持容器平台 Dashboard，是否应作为独立模块而不是容器实例详情扩写

## 10. 回填前置依赖

- 后续若需要把 `logs / events / metrics / terminal / exec` 写成正式接口，需先补充到 `Core v1.yaml`
- 若后续需要把 Deployment、Service、Ingress、ConfigMap、Secret、PVC 等平台资源写成正式模块，应单独拆模块，不从本页扩写
- 若页面要展示更细粒度的发布时间线、回滚策略或副本策略，应先确认 `InstanceRecord.container` 和 `InstanceOperation` 当前 schema 是否已满足前端展示需要
