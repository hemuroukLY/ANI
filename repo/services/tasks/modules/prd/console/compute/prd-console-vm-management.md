# PRD: Console 云主机 VM

## 1. 页面定位

`Console / 算力与云资源 / 实例管理 / 云主机 VM` 用于帮助租户用户查看、创建、操作和排查当前权限范围内的 VM 实例。

本页属于 `Console` 租户侧资源管理页面，底层资源和动作契约都归 `Core`，不属于 `BOSS` 平台运营页面，也不把 VM 资源写入 `Services` 契约。

## 2. 权威源结论

当前一级权威源以 `ANI-main/repo/api/openapi/v1.yaml` 为准，可确认：

- `GET /api/v1/instances` 通过 `kind=vm` 承接 VM 列表
- `GET /api/v1/instances/{instance_id}` 承接 VM 详情
- `POST /api/v1/instances` 承接 VM 创建，成功返回 `201 + CreateInstanceResponse`
- `POST /api/v1/instances/{instance_id}/lifecycle` 承接生命周期动作，成功返回 `200 + InstanceLifecycleResponse`
- `POST /api/v1/instances/{instance_id}/console` 承接短期 Console/VNC 会话，成功返回 `200 + InstanceConsoleSession`
- `GET /api/v1/instances/{instance_id}/operations` 与 `GET /api/v1/instance-operations/{operation_id}` 承接操作历史

## 3. Goals

- 让租户在一个页面内完成 VM 的查找、创建、生命周期操作和排障
- 严格以现有 `Core v1.yaml` 为准，沉淀可直接对齐 `Core` 的 VM 页面定义
- 明确 VM 页面与网络、存储、镜像、快照等关联模块的边界
- 保持 `Console` 租户视角，不混入平台实例池运营口径

## 4. 用户故事

### US-001: 查看 VM 列表

作为租户用户，我希望查看当前租户下的 VM 列表，以便快速找到需要操作的实例。

验收标准：

- 页面展示基于 `kind=vm` 的实例列表
- 列表至少展示名称、状态、规格、Provider、节点、创建时间和最近操作线索
- 列表数据范围限定为当前租户可见实例
- 查询接口对齐 `GET /api/v1/instances?kind=vm`

### US-002: 查看 VM 详情

作为租户用户，我希望查看单个 VM 的详细信息，以便判断当前运行状态和关联资源。

验收标准：

- 详情页或详情抽屉至少展示基础信息、状态、终端保护、SSH、卷、快照、节点和更新时间
- 详情中的字段对齐 `InstanceRecord`
- 页面不要求前端显式传 `tenant_id` 或 `X-Tenant-Id`
- 查询接口对齐 `GET /api/v1/instances/{instance_id}`

### US-003: 创建 VM

作为有权限的租户用户，我希望创建新的 VM 实例，以便交付计算环境。

验收标准：

- 页面提供创建入口，并要求填写名称、CPU、内存、镜像、SSH 等必要字段
- 创建接口对齐 `POST /api/v1/instances`
- 请求体包含 `idempotency_key`
- 创建结果按现有 `Core` 契约返回 `201 + CreateInstanceResponse`
- 返回结果包含 `instance` 和 `operation_id`

### US-004: 执行 VM 生命周期动作

作为有权限的租户用户，我希望对 VM 执行启动、停止、重启、变配等动作，以便完成日常运维。

验收标准：

- 页面提供启动、停止、重启、变配、重建、删除等动作入口
- 生命周期动作对齐 `POST /api/v1/instances/{instance_id}/lifecycle`
- 请求体包含 `action` 和 `idempotency_key`
- 若 VM 开启 `termination_protection`，危险动作必须体现冲突限制
- 响应按现有 `Core` 契约返回 `200 + InstanceLifecycleResponse`

### US-005: 创建 Console / VNC 会话

作为租户用户，我希望为运行中的 VM 创建临时 Console/VNC 会话，以便远程接入实例。

验收标准：

- 页面提供 `console / vnc / novnc / serial` 连接入口
- 接口对齐 `POST /api/v1/instances/{instance_id}/console`
- 会话响应至少包含 `session_id`、`protocol`、`connect_url`、`expires_at`
- 页面只消费短期会话信息，不暴露 provider 长期凭据

### US-006: 查看操作历史与排障信息

作为租户用户，我希望查看 VM 的操作历史和失败原因，以便快速判断问题发生在哪一步。

验收标准：

- 页面提供实例操作历史入口
- 列表接口对齐 `GET /api/v1/instances/{instance_id}/operations`
- 明细接口对齐 `GET /api/v1/instance-operations/{operation_id}`
- 历史记录至少展示操作类型、状态、触发人、时间线步骤、失败原因和是否可重试

### US-007: 保持 VM 页面边界清晰

作为模块维护者，我希望 VM 页面正文只定义 VM 自身和直接关联能力，以便后续与网络、存储等模块拆分维护时没有歧义。

验收标准：

- 页面正文明确 `VM` 资源归 `Core`
- 页面正文不继续使用旧 `/api/v1/console/*` 路径
- 页面正文不把 VPC、子网、块存储、对象存储重新写成 VM 资源契约
- 页面正文把网络、存储、镜像、安全组等定义为关联资源或跳转目标，而不是当前模块主契约

## 5. 功能需求

- FR-1：系统必须提供基于 `kind=vm` 的 VM 列表查询视图
- FR-2：系统必须提供 VM 详情视图，并对齐 `InstanceRecord`
- FR-3：系统必须支持创建 VM，并对齐 `CreateInstanceRequest / CreateInstanceResponse`
- FR-4：系统必须支持对 VM 执行生命周期动作，并对齐 `InstanceLifecycleRequest / InstanceLifecycleResponse`
- FR-5：系统必须支持为 VM 创建临时 Console/VNC 会话
- FR-6：系统必须提供实例操作历史与单个操作详情查询入口
- FR-7：系统必须从认证上下文获取租户边界，不依赖前端传入 `tenant_id`
- FR-8：系统必须对创建和生命周期动作应用 `idempotency_key`
- FR-9：系统必须统一使用标准错误结构
- FR-10：系统必须明确本页属于 `Console` 租户侧，不属于 `BOSS`
- FR-11：系统必须把网络、存储、镜像、快照等保持为关联资源，而不是在本页重新定义其底层契约
- FR-12：系统必须让最终正文可直接作为 VM 模块对齐 `Core` 的主维护材料

## 6. 非目标

- 不在本轮重写网络、VPC、子网、安全组的完整资源契约
- 不在本轮重写块存储、对象存储、文件存储的完整资源契约
- 不在本轮重写镜像、快照、SSH 密钥等独立模块文档
- 不在本轮实现 `BOSS` 侧的平台实例运营总览
- 不在本轮改变现有 `Core v1.yaml` 已冻结的返回码约定

## 7. 产品设计考虑

- 页面口径必须始终是当前租户权限范围内的 VM
- 详情区域允许展示卷、快照、SSH、节点和操作历史摘要，但不展开为完整子域契约
- 对于网络、存储、镜像、安全组等，本页只保留当前关联关系和跳转入口
- 危险动作必须强调终端保护、状态冲突和操作确认

## 8. 成功标准

- 用户能在 30 秒内定位目标 VM 并判断其当前状态
- 用户能在同一页面内完成常见生命周期动作和操作排障
- 文档中的路径、返回码和字段与现有 `Core v1.yaml` 保持一致
- 页面正文不再混写网络、存储、镜像等其他模块的完整契约

## 9. 开放问题

- 是否需要为 VM 页面单独提供监控聚合接口，还是继续通过 `observability` 能力组合展示
- VM 详情首版是否展示完整快照恢复记录，还是只展示快照摘要

## 10. 回填前置依赖

- 后续若需要把 VM 页面中尚未在 `Core v1.yaml` 冻结的监控摘要能力写入契约，应先确认是否复用现有 `/api/v1/observability/query`
- 若后续需要为 VM 页面新增更细粒度的只读接口，应优先补充到 `Core v1.yaml`，不写入 `Services v1.yaml`
- 若页面要把操作时间线、前置校验结果做更稳定的结构化展示，应先确认 `InstanceOperation` 当前 schema 是否满足前端展示需要
