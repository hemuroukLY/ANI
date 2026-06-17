# PRD: Console 模型中心

## 1. Introduction/Overview

`Console / 模型中心` 是租户侧模型管理页面，用于管理模型元数据、模型详情、模型版本和第三方模型导入任务。

本模块属于 `Services` 层，权威源为 `ANI-main/repo/api/openapi/services/v1.yaml` 中的 `Models` 路径组。正式访问前缀为 `/api/v1/svc`，不允许把模型资源写入 `Core /api/v1/*`。

当前已确认的正式能力：

- `GET /api/v1/svc/models`
- `POST /api/v1/svc/models`
- `POST /api/v1/svc/models/import`
- `GET /api/v1/svc/models/{model_id}`
- `DELETE /api/v1/svc/models/{model_id}`
- `POST /api/v1/svc/models/{model_id}/versions`

当前仍属于待补边界的能力：

- 模型加密与密钥绑定
- 模型部署推荐配置
- 模型调用与使用情况
- 平台审核流或治理流

## 2. Goals

- 为租户用户提供统一的模型元数据管理入口
- 支持通过 HuggingFace / ModelScope 发起异步模型导入
- 支持为已有模型创建新版本
- 明确模型中心与推理服务之间的上下游关系
- 严格保持 `Services` 与 `Core` 的资源边界

## 3. User Stories

### US-001: 查看模型列表

作为模型管理员，我希望查看当前租户下的模型列表，以便快速了解已有模型、来源和当前状态。

**Acceptance Criteria**

- [ ] 页面通过 `GET /api/v1/svc/models` 获取模型列表
- [ ] 列表至少展示 `name`、`display_name`、`source`、`status`、`created_at`
- [ ] 列表支持按 `status` 过滤，并支持 `limit / cursor` 分页读取

### US-002: 创建模型元数据

作为模型管理员，我希望先创建模型元数据，再决定是否继续导入或补充模型版本。

**Acceptance Criteria**

- [ ] 页面通过 `POST /api/v1/svc/models` 创建模型元数据
- [ ] 创建表单至少承接 `idempotency_key`、`name`、`display_name`
- [ ] 可选字段仅限权威源已确认的 `description`、`capabilities`
- [ ] 成功后回流到模型详情或列表中的新记录

### US-003: 发起模型导入

作为模型管理员，我希望从 HuggingFace 或 ModelScope 发起导入，并看到真实的异步任务反馈。

**Acceptance Criteria**

- [ ] 页面通过 `POST /api/v1/svc/models/import` 提交导入
- [ ] 请求至少包含 `source`、`repo_id`、`idempotency_key`
- [ ] 导入结果按 `202 + AsyncTask` 呈现，不伪装为同步成功
- [ ] 页面保留导入任务反馈与失败提示

### US-004: 创建模型版本

作为模型管理员，我希望为既有模型创建新版本，以便后续部署引用明确的版本实体。

**Acceptance Criteria**

- [ ] 页面通过 `POST /api/v1/svc/models/{model_id}/versions` 创建模型版本
- [ ] 请求至少承接 `idempotency_key`、`version`、`format`、`storage_path`、`checksum_sha256`、`size_bytes`
- [ ] 成功返回 `201 + ModelVersion`
- [ ] 版本创建完成后，详情页可看到新版本结果

### US-005: 查看与删除模型

作为模型管理员，我希望查看模型详情并在必要时删除模型，以便清理无用或错误的模型记录。

**Acceptance Criteria**

- [ ] 页面通过 `GET /api/v1/svc/models/{model_id}` 查看模型详情
- [ ] 页面通过 `DELETE /api/v1/svc/models/{model_id}` 删除模型
- [ ] 删除动作只在后端返回成功后刷新，不伪造“已删除”
- [ ] 资源不存在时按统一错误结构提示

### US-006: 边界清晰

作为模块维护人，我希望文档明确哪些能力已经冻结、哪些仍待补，避免把规划项写成正式契约。

**Acceptance Criteria**

- [ ] 页面正文明确模型中心属于 `Services`
- [ ] 页面正文不把模型对象写成 `Core` 资源
- [ ] 页面正文不把模型加密、推荐配置、调用统计写成当前已冻结能力
- [ ] 页面正文不使用模糊成功码写法

## 4. Functional Requirements

- FR-1: 系统必须支持模型列表读取、模型详情读取和模型删除
- FR-2: 系统必须支持模型元数据创建，并要求 `idempotency_key`
- FR-3: 系统必须支持第三方模型异步导入，并返回 `AsyncTask`
- FR-4: 系统必须支持模型版本创建，并返回 `ModelVersion`
- FR-5: 系统必须支持只读用户查看列表和详情，但不显示写操作入口
- FR-6: 系统必须按统一错误结构展示失败信息：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- FR-7: 系统必须明确列出待补边界，不把未冻结能力写成正式接口

## 5. Business Rules

- 模型创建、模型删除、模型版本创建都围绕 `Model / ModelVersion` 主资源，不新增独立治理资源域
- 模型导入是异步任务，不等于模型已经 `ready`
- 模型状态当前只承接权威源中已确认的枚举：`pending`、`downloading`、`ready`、`error`、`deleted`
- 只读用户可以看列表与详情，但不能看到创建、导入、删除、版本创建入口
- 模型中心是推理服务的上游选择来源，但不在本页直接创建推理服务

## 6. Non-Goals (Out of Scope)

- 不在本轮定义模型加密与密钥绑定的正式接口
- 不在本轮定义模型部署推荐配置的正式接口
- 不在本轮定义模型调用统计的正式接口
- 不在本轮定义平台审核流、审批流或模型市场治理流程
- 不在本轮扩写模型版本历史独立资源页

## 7. Design Considerations

- 列表页突出“模型状态 + 来源 + 最新操作反馈”
- 详情页突出“基础信息 + 版本信息 + 最近导入反馈”
- 导入反馈必须保留任务态，不允许只给一次“提交成功”
- 删除属于危险动作，必须保留当前上下文并展示明确反馈

## 8. Technical Considerations

- 权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 路径归属：`Services / Models`
- 服务前缀：`/api/v1/svc`
- 页面不要求前端显式传 `tenant_id`
- 创建与版本写操作要求 `idempotency_key`
- 导入接口返回 `202 + AsyncTask`，并带 `Location` 响应头
- 错误响应复用 `ErrorResponse`

## 9. Success Metrics

- 用户能完成“创建模型元数据 -> 导入或创建版本 -> 查看详情反馈”的闭环
- 模型导入反馈具备可读的异步过程态与失败态
- 模块文档与 `services/v1.yaml` 的路径、返回码、schema 名称保持一致
- 文档中不再出现模糊成功码或 `[Assumption]` 口径

## 10. Open Questions

- 列表是否需要显式展示最近版本摘要，还是仅通过详情页承接
- 导入任务结果是否需要未来补充独立任务历史页
- 模型删除是否需要未来补充“被推理服务引用时的前置校验”正式契约

## 11. 回填前置依赖

- 如未来补充模型加密、推荐配置、调用统计，必须先在权威源中冻结正式路径与 schema
- 如未来补充模型审核流，必须先明确其资源归属是否仍属于 `Models`
- 如未来补充任务历史页，应先确认复用 `Core / tasks` 还是 `Services` 自身资源
