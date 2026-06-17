# PRD: Console 推理服务

## 1. Introduction/Overview

`Console / 推理服务` 是租户侧推理部署与运行管理页面，用于管理推理服务列表、部署请求、详情查看和删除请求。

本模块属于 `Services` 层，权威源为 `ANI-main/repo/api/openapi/services/v1.yaml` 中的 `InferenceServices` 路径组。正式访问前缀为 `/api/v1/svc`，不允许把推理服务写入 `Core /api/v1/*`。

当前已确认的正式能力：

- `GET /api/v1/svc/inference-services`
- `POST /api/v1/svc/inference-services`
- `GET /api/v1/svc/inference-services/{service_id}`
- `DELETE /api/v1/svc/inference-services/{service_id}`

当前仍属于待补边界的能力：

- OpenAI 兼容 API 独立资源
- 限流与访问策略独立资源
- 日志 / 事件 / 指标独立资源
- 调用测试独立资源
- 生命周期子页独立后端路径

## 2. Goals

- 让租户用户查看推理服务列表和详情
- 让租户用户提交部署请求并跟踪状态变化
- 让租户用户提交删除请求并得到明确反馈
- 明确推理服务属于 `Services`，不属于 `Core`
- 明确协议、策略、观测等能力当前仍为待补边界

## 3. User Stories

### US-001: 查看推理服务列表

作为服务管理员，我希望查看当前租户的推理服务列表，以便快速确认已部署服务及其运行状态。

**Acceptance Criteria**

- [ ] 页面通过 `GET /api/v1/svc/inference-services` 查询推理服务列表
- [ ] 列表至少展示 `name`、`model`、`replicas`、`gpu_type`、`status`、`created_at`

### US-002: 部署推理服务

作为服务管理员，我希望基于已有模型提交推理服务部署请求，以便获得可运行的推理服务实例。

**Acceptance Criteria**

- [ ] 页面通过 `POST /api/v1/svc/inference-services` 提交部署请求
- [ ] 请求至少包含 `idempotency_key`、`name`、`model`
- [ ] 成功响应按 `202 + InferenceService` 处理，不伪装成同步部署完成
- [ ] 页面不把“一键部署模型”写成独立资源域

### US-003: 查看推理服务详情

作为运维成员，我希望查看单个推理服务详情，以便确认服务当前状态和配置摘要。

**Acceptance Criteria**

- [ ] 页面通过 `GET /api/v1/svc/inference-services/{service_id}` 查询详情
- [ ] 详情至少展示 `name`、`model`、`replicas`、`gpu_type`、`gpu_count_per_pod`、`max_concurrency`、`status`、`endpoint_url`、`created_at`

### US-004: 删除推理服务

作为服务管理员，我希望提交删除推理服务请求，以便停止并清理不再需要的服务。

**Acceptance Criteria**

- [ ] 页面通过 `DELETE /api/v1/svc/inference-services/{service_id}` 提交删除请求
- [ ] 成功响应按 `202` 处理，表示停止任务已提交
- [ ] 删除动作属于危险动作，页面必须有明确确认和回流反馈

### US-005: 保持边界清晰

作为模块维护人，我希望文档明确哪些能力已经冻结、哪些仍待补，避免把规划项写成正式契约。

**Acceptance Criteria**

- [ ] 页面正文明确推理服务属于 `Services`
- [ ] 页面正文不把推理服务对象写成 `Core` 资源
- [ ] 页面正文不把 OpenAI 兼容 API、限流策略、观测资源写成当前已冻结能力
- [ ] 页面正文使用真实成功返回码，不再使用模糊表述

## 4. Functional Requirements

- FR-1: 系统必须支持推理服务列表读取和详情读取
- FR-2: 系统必须支持推理服务部署请求提交，并要求 `idempotency_key`
- FR-3: 系统必须支持推理服务删除请求提交，并明确其为异步停止任务提交
- FR-4: 系统必须展示推理服务状态并仅使用权威源中的状态枚举
- FR-5: 系统必须支持只读用户查看列表和详情，但不显示部署与删除入口
- FR-6: 系统必须明确待补边界，不把未冻结能力写成正式接口

## 5. Business Rules

- 推理服务部署请求成功并不等于服务已 `running`
- 删除成功返回 `202` 仅表示停止 / 删除任务已提交，不等于服务已经完全移除
- 推理服务状态当前只承接权威源中已确认的枚举：`pending`、`deploying`、`running`、`stopping`、`stopped`、`failed`
- 推理服务依赖已有模型信息，但模型选择来源仍由 `模型中心` 提供
- 只读用户可以看列表和详情，但不能看到部署与删除入口

## 6. Non-Goals

- 不在本轮扩写独立 OpenAI 兼容协议资源
- 不在本轮扩写独立限流策略资源
- 不在本轮扩写日志 / 事件 / 指标独立资源
- 不在本轮扩写调用测试独立资源
- 不在本轮把生命周期拆成独立后端路径

## 7. Design Considerations

- 列表页突出“当前状态 + 上游模型 + 部署规模”
- 详情页突出“配置摘要 + 状态反馈 + endpoint 信息”
- 部署与删除都必须保留任务提交反馈，不允许误导为同步完成
- 异常态需要保留 `request_id` 以支持排障

## 8. Technical Considerations

- 权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 路径归属：`Services / InferenceServices`
- 服务前缀：`/api/v1/svc`
- 页面不要求前端显式传 `tenant_id`
- 部署请求必须包含 `idempotency_key`
- 部署接口返回 `202 + InferenceService`
- 删除接口返回 `202`

## 9. Success Metrics

- 用户能完成“查看服务列表 -> 提交部署 -> 查看服务详情 -> 提交删除”的闭环
- 服务部署与删除均具备明确的过程态提示
- 模块文档与 `services/v1.yaml` 的路径、返回码、schema 名称保持一致
- 文档中不再出现把删除写成同步成功或把部署写成 `201` 的错误口径

## 10. Open Questions

- 删除请求提交后，未来是否需要单独的任务跟踪页
- 列表是否需要未来补充健康检查或最近请求摘要字段
- `endpoint_url` 为空时，页面是否需要统一定义“部署中 / 未暴露”的文案规范

## 11. 回填前置依赖

- 如未来补充 OpenAI 兼容协议能力，必须先冻结正式路径与 schema
- 如未来补充限流策略、观测资源、调用测试，必须先确认资源归属
- 如未来补充删除任务历史页，必须先确认是否复用 `AsyncTask` 或其它正式资源
