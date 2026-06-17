# PRD: Console 块存储

## 1. Introduction/Overview

`Console / 算力与云资源 / 存储管理 / 块存储` 用于帮助租户用户查看、创建、删除和理解当前权限范围内的块存储卷资源，并为 VM、容器、GPU 容器等工作负载提供统一的卷资源入口。

本轮 PRD 严格按现有 `Core v1.yaml` 约束生成，当前已核实的页面边界如下：

- `StorageVolume` 属于 `Core`
- 当前 `Core v1.yaml` 已冻结的块存储路径仅包括：
  - `/api/v1/volumes`
  - `/api/v1/volumes/{volume_id}`
- 当前 `Core v1.yaml` 尚未冻结卷快照、卷挂载、卷卸载等独立动作路径
- 本页属于 `Console` 租户侧资源管理页，不是 `BOSS` 的平台存储池运营页
- 本页可以保留快照、挂载/卸载等规划语义，但必须标记为“待补 Core 契约能力”，不能假装已经冻结

## 2. Goals

- 让租户在块存储页面内理解当前租户下卷资源的边界、状态和基础能力
- 严格以现有 `Core v1.yaml` 为准，沉淀块存储页面可直接用于 Core 对齐的主维护材料
- 清楚区分哪些块存储能力已经在 `Core` 冻结，哪些仍为后续待补能力
- 为 VM、容器、GPU 容器等页面提供统一的卷资源回指入口
- 保持 `Console` 租户视角，不混入 `BOSS` 平台存储池运营语义

## 3. User Stories

### US-001: 查看块存储总览
**Description:** 作为租户用户，我希望在块存储页面中快速理解当前卷资源范围和边界，以便确认哪些能力已经可用。

**Acceptance Criteria:**
- [ ] 页面明确这是 `Console` 租户侧的块存储资源页
- [ ] 页面明确当前只承接卷列表、详情、创建、删除
- [ ] 页面明确标出卷快照、挂载、卸载仍为待补能力
- [ ] 页面不混入平台侧全局存储池语义

### US-002: 查看块存储卷列表和详情
**Description:** 作为租户用户，我希望查看卷列表和卷详情，以便了解可用卷容量、状态和基础属性。

**Acceptance Criteria:**
- [ ] 页面提供块存储卷列表入口
- [ ] 页面提供卷详情入口
- [ ] 查询接口对齐 `GET /api/v1/volumes`
- [ ] 详情接口对齐 `GET /api/v1/volumes/{volume_id}`
- [ ] 页面字段口径对齐 `StorageVolume`

### US-003: 创建块存储卷
**Description:** 作为租户用户，我希望创建块存储卷，以便为工作负载准备基础持久化卷资源。

**Acceptance Criteria:**
- [ ] 页面提供创建块存储卷入口
- [ ] 创建接口对齐 `POST /api/v1/volumes`
- [ ] 创建请求体包含 `name`、`size_gib`、`idempotency_key`
- [ ] 可选字段只承接 `storage_class`
- [ ] 成功返回口径对齐 `201 + StorageVolume`

### US-004: 删除块存储卷
**Description:** 作为租户用户，我希望删除不再使用的块存储卷，以便清理租户内资源。

**Acceptance Criteria:**
- [ ] 页面提供删除块存储卷入口
- [ ] 删除接口对齐 `DELETE /api/v1/volumes/{volume_id}`
- [ ] 页面在删除前提示引用风险
- [ ] 文档不把 `409 CONFLICT` 写成当前已冻结删除契约
- [ ] 删除返回口径对齐当前 `Core v1.yaml`

### US-005: 保持块存储模块的 Core 边界清晰
**Description:** 作为 Service 团队成员，我希望块存储正文只声明当前 Core 已冻结的卷资源和待补能力边界，以便后续与 Core 团队对齐时没有歧义。

**Acceptance Criteria:**
- [ ] 页面正文明确卷资源归 `Core`
- [ ] 页面正文不继续使用旧 `/api/v1/storage/*` 或 `/api/v1/console/*` 路径
- [ ] 页面正文不要求前端显式传 `tenant_id / X-Tenant-Id`
- [ ] 页面正文明确快照、挂载、卸载当前不是已冻结 Core 路径
- [ ] 页面正文不把块存储资源写入 `Services /api/v1/svc/*`

## 4. Functional Requirements

- FR-1: 系统必须提供块存储页面总览，并明确已冻结能力与待补能力边界
- FR-2: 系统必须提供块存储卷的查询、详情、创建和删除入口
- FR-3: 系统必须从认证上下文获取租户边界，不依赖前端传入 `tenant_id`
- FR-4: 系统必须对创建接口应用 `idempotency_key`
- FR-5: 系统必须统一使用标准错误结构
- FR-6: 系统必须明确删除动作当前仅承接 `Core v1.yaml` 已冻结返回码
- FR-7: 系统必须让最终正文可直接作为块存储模块对齐 Core 的主维护材料

## 5. Non-Goals (Out of Scope)

- 不在本轮把卷快照写成已冻结 Core 契约
- 不在本轮把卷挂载/卸载动作写成已冻结 Core 契约
- 不在本轮实现卷性能扩缩容、卷克隆等额外动作
- 不在本轮实现 `BOSS` 侧平台存储池运营总览
- 不在本轮改变现有 `Core v1.yaml` 已冻结的返回码和主路径

## 6. Design Considerations

- 页面口径必须是“当前租户 / 当前权限范围内的块存储卷”，不能写成全平台存储池视角
- 总览页可保留快照、挂载、卸载等产品规划语义，但必须标记为待补能力
- VM、容器、GPU 容器页面里的“卷”入口应回指到本模块已冻结的卷资源定义
- 删除前应提示引用关系和风险影响

## 7. Technical Considerations

- 核心资源路径以 `Core v1.yaml` 现有卷路径组为准
- 块存储卷创建使用 `POST /api/v1/volumes`，成功返回 `201`
- 块存储卷删除使用 `DELETE /api/v1/volumes/{volume_id}`，成功返回 `200`
- 删除动作当前显式冻结的错误响应以 `openapi/v1.yaml` 为准，不把 `409` 写成既有契约
- 页面不得要求前端显式传 `tenant_id / X-Tenant-Id`
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

## 8. Success Metrics

- 用户能在 30 秒内识别当前租户下有哪些块存储卷资源
- 用户能在同一模块内完成卷资源的基础创建和查询
- 文档中关于块存储的路径、返回码、字段与现有 `Core v1.yaml` 不再冲突
- 文档不再把快照、挂载、卸载误写成已经冻结的 Core 契约

## 9. Open Questions

- 卷快照后续是否作为独立 `Core` 路径组补充
- 卷挂载/卸载后续是进入独立块存储动作路径，还是继续由工作负载侧承接
- 删除时的依赖冲突语义是否后续需要显式加入 `Core v1.yaml`

## 10. 回填前置依赖

- 后续若要把快照写成正式接口，需先补充到 `Core v1.yaml`
- 后续若要补挂载/卸载，应优先走 `Core v1.yaml` 扩充
- 若页面需要展示更细粒度的卷引用关系，应先确认 `StorageVolume` 当前 schema 是否足够支撑前端展示
