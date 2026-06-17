# PRD: Console Sandbox 实例

## 1. Introduction/Overview

`Console / 算力与云资源 / 实例管理 / Sandbox 实例` 用于帮助租户用户创建、查看和运维当前权限范围内的 Sandbox 实例。本模块属于 `Console` 租户侧实例页面，当前必须严格复用 `Core / Instances` 的统一实例资源，而不是引入独立 `/sandboxes` 资源域。

本轮 PRD 严格按现有 `Core v1.yaml` 约束生成，当前已核实的页面边界如下：

- 当前 Sandbox 实例不拥有独立 `/api/v1/sandboxes*` Core 路径，而是复用 `/api/v1/instances*` 并通过 `kind=sandbox` 区分
- 当前已冻结的 Sandbox 能力包括：统一实例列表过滤、创建、详情、生命周期动作、操作历史、操作详情
- 当前 Sandbox 列表通过 `GET /api/v1/instances` 的 `kind=sandbox` 进入统一实例查询范围，但不等于已经存在独立 `/api/v1/sandboxes*` 资源组
- `CreateInstanceRequest.sandbox_config` 与 `InstanceRecord.sandbox` 是当前 Sandbox 的正式字段来源
- 当前 `Core v1.yaml` 没有冻结 Sandbox 专属 `templates / security-events / metrics / terminal` 路径

## 2. Goals

- 让租户用户通过统一实例能力创建和管理 Sandbox 实例
- 明确当前 Sandbox 页面的正式能力范围是创建、详情、生命周期和操作追踪
- 沉淀一套严格对齐 `Core / Instances` 的主维护材料
- 明确旧原型中的 `/sandboxes*`、安全总览、模板市场等能力当前都不是已冻结契约

## 3. User Stories

### US-001: 创建 Sandbox 实例
**Description:** 作为租户用户，我希望创建 Sandbox 实例，以便获得受限和隔离的安全执行环境。

**Acceptance Criteria:**
- [ ] 创建接口对齐 `POST /api/v1/instances`
- [ ] 请求体使用 `CreateInstanceRequest`，并传入 `kind=sandbox`
- [ ] 页面承接 `sandbox_config.runtime_class`、`session_timeout`、`network_egress_policy`
- [ ] 成功响应为 `201 + CreateInstanceResponse`

### US-002: 查看 Sandbox 实例详情
**Description:** 作为租户用户，我希望查看 Sandbox 实例详情，以便了解当前会话和隔离配置状态。

**Acceptance Criteria:**
- [ ] 详情接口对齐 `GET /api/v1/instances/{instance_id}`
- [ ] 页面展示 `sandbox.runtime_class`、`session_timeout`、`network_egress_policy`、`session_state`
- [ ] 页面不扩写独立安全事件和模板资源

### US-003: 执行 Sandbox 生命周期动作
**Description:** 作为租户用户，我希望对 Sandbox 实例执行基础生命周期动作，以便控制会话运行状态。

**Acceptance Criteria:**
- [ ] 生命周期接口对齐 `POST /api/v1/instances/{instance_id}/lifecycle`
- [ ] 页面当前只承接 `start`、`stop`、`restart`、`delete`
- [ ] 写操作必须带 `idempotency_key`
- [ ] 页面不把 `extend / suspend / resume` 写成现有已冻结动作

### US-004: 查看 Sandbox 操作历史
**Description:** 作为租户用户，我希望查看 Sandbox 实例操作历史和单次操作详情，以便排查创建和状态变更结果。

**Acceptance Criteria:**
- [ ] 历史接口对齐 `GET /api/v1/instances/{instance_id}/operations`
- [ ] 详情接口对齐 `GET /api/v1/instance-operations/{operation_id}`
- [ ] 页面展示操作类型、状态、发起人、失败原因和步骤时间线

### US-005: 保持页面边界清晰
**Description:** 作为 Service 团队成员，我希望 Sandbox 页面只定义当前 Core 已冻结的实例能力，以便后续安全能力拆分时没有歧义。

**Acceptance Criteria:**
- [ ] 页面正文明确当前没有冻结独立 `/api/v1/sandboxes*` 路径
- [ ] 页面正文不把 `/api/v1/sandboxes*` 写成正式契约
- [ ] 页面正文不把模板列表、安全事件、监控指标、安全总览写成已冻结路径

## 4. Functional Requirements

- FR-1: 系统必须支持创建 Sandbox 实例
- FR-2: 系统必须支持读取 Sandbox 实例详情
- FR-3: 系统必须支持基础生命周期动作
- FR-4: 系统必须支持操作历史与操作详情
- FR-5: 系统必须明确 Sandbox 列表来自统一实例查询，而不是独立 Sandbox 资源组

## 5. Non-Goals (Out of Scope)

- 不在本轮实现独立 `/api/v1/sandboxes*` 资源
- 不在本轮实现独立 `/api/v1/sandboxes*` 列表接口
- 不在本轮实现模板列表、安全事件、监控指标、安全总览
- 不在本轮实现独立终端或会话延长动作

## 6. Design Considerations

- 页面必须强调 Sandbox 复用统一实例 `kind=sandbox`
- 由于当前列表能力来自统一实例查询，页面应以“列表、创建、详情和操作追踪”组成统一主链路
- 旧原型中的安全态势和模板内容必须降级为待补边界

## 7. Technical Considerations

- 权威路径以 `Core v1.yaml` 的 `Instances` 路径组为准
- 创建使用 `POST /api/v1/instances`
- 详情使用 `GET /api/v1/instances/{instance_id}`
- 生命周期使用 `POST /api/v1/instances/{instance_id}/lifecycle`
- 操作历史使用 `GET /api/v1/instances/{instance_id}/operations`
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

## 8. Success Metrics

- 用户能完成 Sandbox 实例创建、详情查看和基础状态控制
- 文档不再出现独立 `/sandboxes*` 正式路径
- HTML 摘要、PRD、SPEC 与主文档保持一致

## 9. Open Questions

- 后续是否需要把 Sandbox 从统一实例列表中拆成独立资源组
- 后续是否要把安全事件、模板与监控拆为独立模块
