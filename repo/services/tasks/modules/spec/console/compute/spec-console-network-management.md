# SPEC: Console 网络管理

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-network-management.md`
> Generated: 2026-06-13 | Scope: `Console / 算力与云资源 / 网络管理`

## 1. Summary

### 1.1 What This SPEC Covers

本 SPEC 说明 `Console / 算力与云资源 / 网络管理` 的技术边界、数据模型、接口冻结规则和错误处理方式。目标是把网络管理模块收口成一份可直接用于对齐 `Core v1.yaml` 的主维护材料，并明确区分“当前 Core 已冻结能力”与“仍待补充到 Core 的规划能力”。

### 1.2 PRD Reference

- Source: `tasks/modules/prd/console/compute/prd-console-network-management.md`
- User Stories covered: `US-001` ~ `US-006`
- Functional Requirements covered: `FR-1` ~ `FR-10`

### 1.3 Design Decisions Summary

| Decision | Choice | Rationale |
|---|---|---|
| 模块归属 | `Console` 租户侧资源管理页 | 与 `BOSS` 平台网络运营页分离 |
| 资源归属 | `VPC / 子网 / 安全组 / 负载均衡` 归 `Core` | 现有 `Core v1.yaml` 已冻结对应 schema 和路径 |
| 路由能力处理 | 标记为待补 `Core` 能力 | 当前 `Core v1.yaml` 未冻结独立路由资源 |
| 子网 IP 分配处理 | 标记为待补 `Core` 能力 | 当前 `Core v1.yaml` 未冻结相关路径 |
| 安全组细粒度规则 CRUD | 暂不扩展独立子路径 | 当前仅冻结整体 `rules` 结构 |
| 安全组绑定解绑 | 标记为待补 `Core` 能力 | 当前未冻结对应子路径 |

## 2. Architecture

### 2.1 System Context

- 页面属于 `Console / 算力与云资源`
- 页面直接消费 `Core /api/v1/networks/*`
- 页面不定义新的 `Services /api/v1/svc/*` 网络资源契约
- VM、容器、GPU 容器等页面中的“网络关联”只消费本模块已冻结的网络资源摘要和跳转入口

### 2.2 Component Design

- `网络管理总览`
  - 展示 VPC、子网、安全组、负载均衡的关系说明
  - 展示哪些能力已经 Core 冻结，哪些仍待补
- `VPC 资源区`
  - 承接列表、详情、创建、删除
- `子网资源区`
  - 承接列表、详情、创建、删除
- `安全组资源区`
  - 承接列表、详情、创建、删除
  - 当前仅消费整体 `rules`
- `负载均衡资源区`
  - 承接列表、详情、创建、删除
- `待补能力说明区`
  - 路由
  - 子网 IP 分配列表
  - 安全组规则子资源 CRUD
  - 绑定/解绑安全组到实例或网卡

### 2.3 Module Interactions

1. 用户打开网络管理页
2. 页面展示资源关系与能力边界
3. 用户在不同资源区执行列表查询
4. 用户通过创建抽屉或表单发起创建
5. 用户在详情页查看基础信息与关联摘要
6. 删除动作在确认依赖后调用对应 `DELETE` 接口

### 2.4 File Structure

```text
docs/
└── console-modules/
    └── compute/
        └── network-management.md        [NEW]

tasks/
├── prd-console-network-management.md   [NEW]
└── spec-console-network-management.md  [NEW]

root/
├── prototypes/ani-services-prototype-console.html [MODIFY]
└── prototypes/ani-services-prototype.html         [MODIFY]
```

## 3. Data Model

### 3.1 Core Frozen Schemas

已冻结并可直接引用 `Core v1.yaml` 的 schema：

- `NetworkVPC`
- `NetworkSubnet`
- `NetworkSecurityGroup`
- `NetworkSecurityGroupRule`
- `NetworkLoadBalancer`
- `NetworkLoadBalancerListener`
- `NetworkVPCListResponse`
- `NetworkSubnetListResponse`
- `NetworkSecurityGroupListResponse`
- `NetworkLoadBalancerListResponse`
- `CreateNetworkVPCRequest`
- `CreateNetworkSubnetRequest`
- `CreateNetworkSecurityGroupRequest`
- `CreateNetworkLoadBalancerRequest`

### 3.2 Entity Definitions

#### NetworkVPC

- 关键字段：`id`、`name`、`cidr`、`state`、`reason`、`created_at`、`updated_at`
- 注意：`tenant_id` 在响应中存在，但页面不要求前端主动传入

#### NetworkSubnet

- 关键字段：`id`、`vpc_id`、`name`、`cidr`、`gateway`、`state`、`reason`、`created_at`、`updated_at`

#### NetworkSecurityGroup

- 关键字段：`id`、`name`、`description`、`rules[]`、`state`、`reason`、`created_at`、`updated_at`
- 当前规则仅以 `rules[]` 整体消费

#### NetworkLoadBalancer

- 关键字段：`id`、`name`、`vpc_id`、`subnet_id`、`scheme`、`vip`、`listeners[]`、`state`、`reason`、`created_at`、`updated_at`

### 3.3 Non-Frozen Data Areas

以下内容在当前 `Core v1.yaml` 中未冻结，不应写成正式 Core 契约：

- `RouteTable`
- `RouteEntry`
- `SubnetIpAllocation`
- `SecurityGroupRuleResource`
- `SecurityGroupBinding`

### 3.4 Naming Convention

- 页面展示字段允许使用 camelCase 作为前端 view model
- API 路径参数、query 参数、请求体和响应 JSON schema 统一以现有 `Core v1.yaml` 命名为准
- 不能把旧 HTML 中的 `vpcId / subnetId / securityGroupId / pageNo / pageSize` 直接写成已冻结接口参数

## 4. API Design

### 4.1 Frozen Endpoints

| Method | Path | Description | Auth | Request | Response |
|---|---|---|---|---|---|
| GET | `/api/v1/networks/vpcs` | 查询 VPC 列表 | Bearer / ApiKey | `limit?`、`cursor?` | `200 + NetworkVPCListResponse` |
| POST | `/api/v1/networks/vpcs` | 创建 VPC | Bearer / ApiKey | `CreateNetworkVPCRequest` | `201 + NetworkVPC` |
| GET | `/api/v1/networks/vpcs/{vpc_id}` | 查询 VPC | Bearer / ApiKey | path: `vpc_id` | `200 + NetworkVPC` |
| DELETE | `/api/v1/networks/vpcs/{vpc_id}` | 删除 VPC | Bearer / ApiKey | path: `vpc_id` | `200 + NetworkVPC` |
| GET | `/api/v1/networks/subnets` | 查询子网列表 | Bearer / ApiKey | `limit?`、`cursor?` | `200 + NetworkSubnetListResponse` |
| POST | `/api/v1/networks/subnets` | 创建子网 | Bearer / ApiKey | `CreateNetworkSubnetRequest` | `201 + NetworkSubnet` |
| GET | `/api/v1/networks/subnets/{subnet_id}` | 查询子网 | Bearer / ApiKey | path: `subnet_id` | `200 + NetworkSubnet` |
| DELETE | `/api/v1/networks/subnets/{subnet_id}` | 删除子网 | Bearer / ApiKey | path: `subnet_id` | `200 + NetworkSubnet` |
| GET | `/api/v1/networks/security-groups` | 查询安全组列表 | Bearer / ApiKey | `limit?`、`cursor?` | `200 + NetworkSecurityGroupListResponse` |
| POST | `/api/v1/networks/security-groups` | 创建安全组 | Bearer / ApiKey | `CreateNetworkSecurityGroupRequest` | `201 + NetworkSecurityGroup` |
| GET | `/api/v1/networks/security-groups/{security_group_id}` | 查询安全组 | Bearer / ApiKey | path: `security_group_id` | `200 + NetworkSecurityGroup` |
| DELETE | `/api/v1/networks/security-groups/{security_group_id}` | 删除安全组 | Bearer / ApiKey | path: `security_group_id` | `200 + NetworkSecurityGroup` |
| GET | `/api/v1/networks/load-balancers` | 查询负载入口列表 | Bearer / ApiKey | `limit?`、`cursor?` | `200 + NetworkLoadBalancerListResponse` |
| POST | `/api/v1/networks/load-balancers` | 创建负载入口 | Bearer / ApiKey | `CreateNetworkLoadBalancerRequest` | `201 + NetworkLoadBalancer` |
| GET | `/api/v1/networks/load-balancers/{load_balancer_id}` | 查询负载入口 | Bearer / ApiKey | path: `load_balancer_id` | `200 + NetworkLoadBalancer` |
| DELETE | `/api/v1/networks/load-balancers/{load_balancer_id}` | 删除负载入口 | Bearer / ApiKey | path: `load_balancer_id` | `200 + NetworkLoadBalancer` |

### 4.2 Per-Endpoint Frozen Rules

#### `POST /api/v1/networks/vpcs`

| 项 | 值 |
|---|---|
| operationId | `createNetworkVPC` |
| summary | `创建 VPC` |
| tags | `["Networks"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| preconditions | 已通过标准鉴权；`name` 与 `idempotency_key` 非空；`cidr` 满足 schema 约束 |
| requestBody.required | `name`、`idempotency_key` |
| success | `201 + NetworkVPC` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN` |

#### `POST /api/v1/networks/subnets`

| 项 | 值 |
|---|---|
| operationId | `createNetworkSubnet` |
| summary | `创建子网` |
| tags | `["Networks"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| preconditions | 已通过标准鉴权；`vpc_id`、`name`、`idempotency_key` 非空；`vpc_id` 当前租户可见 |
| requestBody.required | `vpc_id`、`name`、`idempotency_key` |
| success | `201 + NetworkSubnet` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

#### `POST /api/v1/networks/security-groups`

| 项 | 值 |
|---|---|
| operationId | `createNetworkSecurityGroup` |
| summary | `创建安全组` |
| tags | `["Networks"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| preconditions | 已通过标准鉴权；`name` 与 `idempotency_key` 非空；若提供 `rules[]` 必须满足已冻结 schema |
| requestBody.required | `name`、`idempotency_key` |
| success | `201 + NetworkSecurityGroup` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN` |

#### `POST /api/v1/networks/load-balancers`

| 项 | 值 |
|---|---|
| operationId | `createNetworkLoadBalancer` |
| summary | `创建负载入口` |
| tags | `["Networks"]` |
| security | `[{BearerAuth: []}, {ApiKeyAuth: []}]` |
| preconditions | 已通过标准鉴权；`name`、`vpc_id`、`idempotency_key` 非空；若提供 `subnet_id` 必须属于目标 VPC |
| requestBody.required | `name`、`vpc_id`、`idempotency_key` |
| success | `201 + NetworkLoadBalancer` |
| error responses | `400 BAD_REQUEST`、`401 UNAUTHORIZED`、`403 FORBIDDEN`、`404 NOT_FOUND` |

### 4.3 Non-Frozen Endpoints

以下路径当前**不能**写成已对齐 Core 的正式接口：

- `/api/v1/networks/routes`
- `/api/v1/networks/subnets/{subnet_id}/ip-allocations`
- `/api/v1/networks/security-groups/{security_group_id}/rules`
- `/api/v1/networks/security-groups/{security_group_id}/bindings`

如果后续需要补这些路径，应先新增到 `Core v1.yaml`，再回写模块正文。

### 4.4 Operation Availability Matrix

| 资源 | list/get | create | delete | 删除前置校验 |
|---|---|---|---|---|
| VPC | 可用 | 可用 | 可用 | 不允许删除仍存在子网或被负载均衡引用的 VPC；应返回 `409 CONFLICT` |
| 子网 | 可用 | 可用 | 可用 | 不允许删除仍有实例落点或关键绑定的子网；应返回 `409 CONFLICT` |
| 安全组 | 可用 | 可用 | 可用 | 不允许删除仍被资源引用的安全组；应返回 `409 CONFLICT` |
| 负载均衡 | 可用 | 可用 | 可用 | 删除前需校验无不可删除的活动依赖；冲突时返回 `409 CONFLICT` |

### 4.5 Delete Conflict Rule

- `DELETE /api/v1/networks/vpcs/{vpc_id}`：若仍存在子网或被负载均衡引用，应返回 `409 CONFLICT`
- `DELETE /api/v1/networks/subnets/{subnet_id}`：若仍有实例落点或关键绑定，应返回 `409 CONFLICT`
- `DELETE /api/v1/networks/security-groups/{security_group_id}`：若仍被资源引用，应返回 `409 CONFLICT`
- `DELETE /api/v1/networks/load-balancers/{load_balancer_id}`：若存在不可删除的活动依赖，应返回 `409 CONFLICT`

### 4.6 Request/Response Examples

#### CreateNetworkVPCRequest

```json
{
  "idempotency_key": "idem-network-vpc-001",
  "name": "vpc-prod",
  "cidr": "10.20.0.0/16"
}
```

#### NetworkLoadBalancer

```json
{
  "id": "lb-3a1f",
  "tenant_id": "t-demo",
  "name": "ingress-public",
  "vpc_id": "vpc-8f3a",
  "subnet_id": "subnet-12ab",
  "scheme": "public",
  "vip": "203.0.113.10",
  "listeners": [
    {
      "protocol": "https",
      "port": 443,
      "target_port": 8443
    }
  ],
  "state": "available",
  "reason": null,
  "created_at": "2026-06-13T10:00:00Z",
  "updated_at": "2026-06-13T10:00:00Z"
}
```

## 5. Business Logic

### 5.1 Core Page Logic

- 网络管理页默认展示当前租户可见资源
- 总览页负责说明对象关系与能力边界
- 各资源区以 `Core v1.yaml` 现有接口作为唯一事实来源
- 对已冻结资源，只沉淀当前已存在的列表、详情、创建、删除能力
- 对未冻结能力，只记录为“待补 Core 契约”，不能写成正式接口表

### 5.2 Validation Rules

- 创建 VPC 必须提供 `name` 和 `idempotency_key`
- 创建子网必须提供 `vpc_id`、`name` 和 `idempotency_key`
- 创建安全组必须提供 `name` 和 `idempotency_key`
- 创建负载均衡必须提供 `name`、`vpc_id` 和 `idempotency_key`
- 前端不可要求用户主动填写 `tenant_id`
- 页面不得继续使用旧 `/api/v1/network/*` 或 `/api/v1/console/*` 路径

### 5.3 State Handling

- 资源状态使用 `NetworkResourceState`
- 页面只按已有 `state` 和 `reason` 展示标签与异常提示
- 不在本轮引入新的网络状态机

### 5.4 Edge Cases

- VPC 删除前仍存在子网或被负载均衡引用
- 子网删除前仍有实例落点
- 安全组删除前仍有资源引用
- 负载均衡创建时引用不存在的 `vpc_id`
- 列表无数据时展示空态而不是伪造统计值

## 6. Error Handling

### 6.1 Error Taxonomy

| Error Code | HTTP Status | Condition | User Message |
|---|---|---|---|
| `BAD_REQUEST` | 400 | 参数格式错误、CIDR 不合法、监听器配置错误 | 请求参数不合法，请检查后重试 |
| `UNAUTHORIZED` | 401 | 未登录或凭证失效 | 登录状态已失效，请重新登录 |
| `FORBIDDEN` | 403 | 无访问或操作权限 | 当前账号无权访问该网络资源 |
| `NOT_FOUND` | 404 | `vpc_id` / `subnet_id` / `security_group_id` / `load_balancer_id` 不存在 | 目标资源不存在或已删除 |
| `CONFLICT` | 409 | 删除前存在依赖引用，或资源当前状态不允许删除 | 资源仍被引用，无法执行删除 |

### 6.2 Shared Error Format

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-xxx"
}
```

### 6.3 Failure Modes

- 删除类动作若后端返回依赖冲突，页面展示错误码、错误消息和 `request_id`
- 查询失败时仅影响当前资源区，不阻断整个页面框架
- 待补能力若无后端契约，页面只能展示占位说明或跳转禁用，不得伪造接口

## 7. Security

### 7.1 Authentication & Authorization

- 全部网络接口使用 `security: [{BearerAuth: []}, {ApiKeyAuth: []}]`
- 后端必须从认证上下文获取租户边界
- 前端不可信任也不显式传 `tenant_id / X-Tenant-Id`

### 7.2 Input Validation

- `idempotency_key` 必须非空
- `name` 必须符合长度要求
- `cidr`、`gateway`、`listeners.port`、`listeners.target_port` 必须符合 schema 约束

### 7.3 Data Protection

- 不展示平台侧全局网络池信息
- 只展示当前租户可见资源
- 页面只消费网络资源自身和必要关联摘要，不向前端暴露额外平台内部实现细节

## 8. Performance

### 8.1 Expected Load

- 网络资源列表以租户级分页查询为主
- 首版不要求跨资源全局聚合接口

### 8.2 Optimization Strategy

- 列表统一沿用 `limit + cursor`
- 详情按需加载
- 不在首版为待补能力引入额外查询链路

### 8.3 Database Considerations

- 本文档不新增数据库设计
- 数据模型以现有 `Core` schema 为准

## 9. Testing Strategy

### 9.1 Unit Tests

- 校验创建请求字段映射
- 校验待补能力不会误显示为已冻结接口

### 9.2 Integration Tests

- VPC / 子网 / 安全组 / 负载均衡列表与详情接口映射正确
- 创建请求带有 `idempotency_key`
- 删除接口路径参数命名与 `Core v1.yaml` 一致

### 9.3 Edge Case Tests

- 删除不存在资源时展示 `NOT_FOUND`
- 创建子网引用不存在 `vpc_id` 时展示 `NOT_FOUND`
- 负载均衡监听器配置错误时展示 `BAD_REQUEST`

### 9.4 Acceptance Criteria Mapping

| US/FR | Test | Type | Description |
|---|---|---|---|
| `US-001` | 网络总览边界校验 | integration | 验证已冻结能力与待补能力被正确区分 |
| `US-002` | VPC 契约校验 | integration | 验证 VPC 路径、返回码、字段对齐 `Core v1.yaml` |
| `US-003` | 子网契约校验 | integration | 验证子网路径、返回码、字段对齐 `Core v1.yaml` |
| `US-004` | 安全组契约校验 | integration | 验证安全组只消费整体 `rules` 结构 |
| `US-005` | 负载均衡契约校验 | integration | 验证 `scheme / listeners / vip` 等字段对齐 |
| `US-006` | Core 边界校验 | unit | 验证未冻结能力不会写入正式接口表 |

## 10. Implementation Plan

### 10.1 Phases

1. 明确网络模块已冻结与未冻结边界
2. 生成 PRD 和 SPEC
3. 收口主维护源
4. 回填 HTML 摘要
5. 执行最终 Core 合规复审

### 10.2 Issue Mapping

| Issue | SPEC Sections | Priority | Depends On |
|---|---|---|---|
| 网络边界收口 | 1, 2, 4.3 | high | — |
| Core 已冻结接口整理 | 3, 4.1, 4.2 | high | 网络边界收口 |
| 主维护源落盘 | 5, 6, 7 | high | Core 已冻结接口整理 |
| HTML 摘要回填 | 2.4, 10.1 | medium | 主维护源落盘 |

### 10.3 Incremental Delivery

- 本轮先收口主维护材料，不直接扩写新的 `Core v1.yaml`
- 后续若要补“路由”“IP 分配”“规则子资源”等能力，再进入 `YAML` 扩充阶段

## 11. Open Questions & Risks

### 11.1 Unresolved Questions

- 路由表后续是否作为独立路径组进入 `Core v1.yaml`
- 子网 IP 分配列表是否需要成为下一个网络侧补充点
- 安全组规则是否需要拆成独立资源与动作

### 11.2 Technical Risks

| Risk | Impact | Mitigation |
|---|---|---|
| 旧 HTML 路径仍沿用 `/api/v1/network/*` | 会造成前后端继续对错路径 | 统一改为 `/api/v1/networks/*` |
| 旧文档继续要求 `X-Tenant-Id` | 会与 Core 网关边界冲突 | 文档统一改为认证上下文获取 |
| 把未冻结能力写进正式契约 | 后续 Core 对齐会失真 | 在正文和 SPEC 中单列“待补能力” |

### 11.3 Frozen Baseline

- 当前 `Core v1.yaml` 就是网络资源的唯一权威来源
- 页面不新增独立 `Services` 网络资源聚合接口
- 本轮只收口 `Console / 网络管理` 顶层模块，不同步细化每个子页面到独立主维护文档
