# PRD: Console 网络管理

## 1. Introduction/Overview

`Console / 算力与云资源 / 网络管理` 用于帮助租户用户查看和管理当前权限范围内的网络资源，承接 `Core` 中已经冻结的 VPC、子网、安全组和负载均衡能力，并为 VM、容器、GPU 容器等工作负载提供网络关联入口。

本轮 PRD 严格按现有 `Core v1.yaml` 约束生成，当前已核实的页面边界如下：

- `VPC`、`子网`、`安全组`、`负载均衡` 资源对象本身属于 `Core`
- 这些资源的查询、创建、详情、删除统一对齐 `Core /api/v1/networks/*`
- 当前 `Core v1.yaml` **尚未冻结** 独立的路由表资源、子网 IP 分配列表、安全组规则子资源 CRUD、绑定/解绑安全组到实例或网卡等路径
- 本页可以继续保留“路由”“IP 分配”“细粒度规则管理”“绑定关系”这些规划语义，但在主维护材料中必须标记为“待补 Core 契约”，不能假装已经冻结
- 页面属于 `Console` 租户侧资源管理页，不是 `BOSS` 平台网络运营页

## 2. Goals

- 让租户在一个页面内完成网络域核心资源的查看、创建、详情和删除
- 严格以现有 `Core v1.yaml` 为准，沉淀网络管理页面可直接用于 Core 对齐的主维护材料
- 明确网络管理页面里哪些能力已经在 `Core` 冻结，哪些仍是后续待补能力
- 让 VM、容器等实例页面中的“网络关联入口”能够回指到清晰的网络资源定义
- 保持 `Console` 租户视角，不混入 `BOSS` 平台网络资源池运营语义

## 3. User Stories

### US-001: 查看网络管理总览
**Description:** 作为租户用户，我希望在一个网络管理页面里快速理解 VPC、子网、安全组、负载均衡之间的关系，以便找到要操作的网络资源。

**Acceptance Criteria:**
- [ ] 页面明确区分 `VPC / 子网 / 安全组 / 负载均衡`
- [ ] 页面明确标出哪些能力已经在 `Core v1.yaml` 冻结
- [ ] 页面明确标出“路由”等仍为规划项或待补能力
- [ ] 页面不混入平台侧全局网络池口径

### US-002: 查看和创建 VPC
**Description:** 作为租户用户，我希望查看和创建 VPC，以便建立自己的网络隔离边界。

**Acceptance Criteria:**
- [ ] 页面提供 VPC 列表和详情入口
- [ ] 页面提供创建 VPC 入口
- [ ] 查询接口对齐 `GET /api/v1/networks/vpcs`
- [ ] 创建接口对齐 `POST /api/v1/networks/vpcs`
- [ ] 创建请求体包含 `idempotency_key`
- [ ] 删除接口对齐 `DELETE /api/v1/networks/vpcs/{vpc_id}`

### US-003: 查看和创建子网
**Description:** 作为租户用户，我希望在指定 VPC 下查看和创建子网，以便进行地址划分和实例落点管理。

**Acceptance Criteria:**
- [ ] 页面提供子网列表和详情入口
- [ ] 页面提供创建子网入口
- [ ] 查询接口对齐 `GET /api/v1/networks/subnets`
- [ ] 创建接口对齐 `POST /api/v1/networks/subnets`
- [ ] 创建请求体包含 `idempotency_key`
- [ ] 删除接口对齐 `DELETE /api/v1/networks/subnets/{subnet_id}`

### US-004: 查看和创建安全组
**Description:** 作为租户用户，我希望查看和创建安全组，以便管理工作负载的基础网络访问策略。

**Acceptance Criteria:**
- [ ] 页面提供安全组列表和详情入口
- [ ] 页面提供创建安全组入口
- [ ] 查询接口对齐 `GET /api/v1/networks/security-groups`
- [ ] 创建接口对齐 `POST /api/v1/networks/security-groups`
- [ ] 创建请求体包含 `idempotency_key`
- [ ] 删除接口对齐 `DELETE /api/v1/networks/security-groups/{security_group_id}`
- [ ] 页面只承接当前已冻结的整体 `rules` 结构，不凭空扩展独立规则 CRUD 子路径

### US-005: 查看和创建负载均衡
**Description:** 作为租户用户，我希望查看和创建负载均衡入口，以便把服务稳定暴露给内部或外部访问者。

**Acceptance Criteria:**
- [ ] 页面提供负载均衡列表和详情入口
- [ ] 页面提供创建负载均衡入口
- [ ] 查询接口对齐 `GET /api/v1/networks/load-balancers`
- [ ] 创建接口对齐 `POST /api/v1/networks/load-balancers`
- [ ] 创建请求体包含 `idempotency_key`
- [ ] 删除接口对齐 `DELETE /api/v1/networks/load-balancers/{load_balancer_id}`
- [ ] 页面展示 `scheme / vip / listeners` 等已冻结字段

### US-006: 保持网络模块的 Core 边界清晰
**Description:** 作为 Service 团队成员，我希望网络管理正文只声明当前 Core 已冻结的网络资源和待补能力边界，以便后续与 Core 团队对齐时没有歧义。

**Acceptance Criteria:**
- [ ] 页面正文明确网络资源归 `Core`
- [ ] 页面正文不继续使用旧 `/api/v1/network/*` 或 `/api/v1/console/*` 路径
- [ ] 页面正文不要求前端显式传 `tenant_id / X-Tenant-Id`
- [ ] 页面正文明确“路由、IP 分配、安全组规则子资源、绑定解绑”当前不是已冻结 Core 路径
- [ ] 页面正文不把网络资源写入 `Services /api/v1/svc/*`

## 4. Functional Requirements

- FR-1: 系统必须提供网络管理总览，并明确各子资源关系
- FR-2: 系统必须提供 VPC 的查询、创建、详情和删除入口
- FR-3: 系统必须提供子网的查询、创建、详情和删除入口
- FR-4: 系统必须提供安全组的查询、创建、详情和删除入口
- FR-5: 系统必须提供负载均衡的查询、创建、详情和删除入口
- FR-6: 系统必须从认证上下文获取租户边界，不依赖前端传入 `tenant_id`
- FR-7: 系统必须对各创建接口应用 `idempotency_key`
- FR-8: 系统必须统一使用标准错误结构
- FR-9: 系统必须明确哪些网络能力已冻结、哪些仍待补
- FR-10: 系统必须让最终正文可直接作为网络模块对齐 Core 的主维护材料

## 5. Non-Goals (Out of Scope)

- 不在本轮把独立路由表资源写成已冻结 Core 契约
- 不在本轮把子网 IP 分配列表写成已冻结 Core 契约
- 不在本轮把安全组规则独立 CRUD 写成已冻结 Core 契约
- 不在本轮把安全组绑定/解绑实例或网卡写成已冻结 Core 契约
- 不在本轮实现 BOSS 侧平台网络资源池运营总览
- 不在本轮改变现有 `Core v1.yaml` 已冻结的返回码和主路径

## 6. Design Considerations

- 页面口径必须是“当前租户 / 当前权限范围内的网络资源”，不能写成全平台网络池视角
- 总览页可保留“路由”等产品规划语义，但必须标记为待补能力，不能伪装成已对齐 API
- 实例详情中的“网络”入口应当回指到本模块已冻结的 VPC、子网、安全组、负载均衡
- 风险操作应强调删除前依赖校验和资源引用关系提示

## 7. Technical Considerations

- 核心资源路径以 `Core v1.yaml` 现有 `networks/*` 路径组为准
- VPC 创建使用 `POST /api/v1/networks/vpcs`，成功返回 `201`
- 子网创建使用 `POST /api/v1/networks/subnets`，成功返回 `201`
- 安全组创建使用 `POST /api/v1/networks/security-groups`，成功返回 `201`
- 负载均衡创建使用 `POST /api/v1/networks/load-balancers`，成功返回 `201`
- 删除接口分别对齐对应 `DELETE /api/v1/networks/.../{id}`，成功返回 `200`
- 页面不得要求前端显式传 `tenant_id / X-Tenant-Id`
- 错误响应格式统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

## 8. Success Metrics

- 用户能在 30 秒内识别当前租户下有哪些网络资源
- 用户能在同一模块内完成已冻结网络资源的基础创建和查询
- 文档中关于网络模块的路径、返回码、字段与现有 `Core v1.yaml` 不再冲突
- 文档不再把“待补能力”误写成已经冻结的 Core 契约

## 9. Open Questions

- 路由表是否后续作为独立 `Core` 资源组补充，还是继续附着在 VPC/子网之下
- 子网 IP 分配列表是否需要在下一轮进入 `Core v1.yaml`
- 安全组规则是否保持内嵌在安全组资源中，还是未来拆成独立规则子资源

## 10. 回填前置依赖

- 后续若要把“路由”写成正式接口，需先补充到 `Core v1.yaml`，不能只留在 HTML
- 后续若要补子网 IP 分配列表、安全组规则 CRUD、绑定解绑等能力，应优先走 `Core v1.yaml` 扩充
- 若页面需要展示更细粒度的引用关系，应先确认 `NetworkVPC / NetworkSubnet / NetworkSecurityGroup / NetworkLoadBalancer` 当前 schema 是否已足够支撑前端展示
