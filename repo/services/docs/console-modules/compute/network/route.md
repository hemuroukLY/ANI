# 网络管理 — 路由

## 页面定位

`路由` 是网络管理下路由表/路由条目的管理页。

父级：`network-management.md`。

## 文档管理规则

- 本文是路由子模块主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`

## Core 层要求

<!-- ADDED-TO-YAML: GET/POST /api/v1/networks/routes (Core v1.yaml, Phase 2 2026-06-17) -->

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| GET | `/api/v1/networks/routes` | `listNetworkRoutes` | `scope:networks:read` |
| POST | `/api/v1/networks/routes` | `createNetworkRoute` | `scope:networks:create` |

Query（list）：`vpc_id`、`limit`、`cursor`。

Schema：`NetworkRoute`、`NetworkRouteListResponse`、`CreateNetworkRouteRequest`。

**当前 YAML 未声明** 单条路由 GET/DELETE 路径；文档不得自造 `{route_id}` 路径。

POST 创建必须带 `idempotency_key`（以 `CreateNetworkRouteRequest` 为准）。

## 页面职责

- 按 VPC 展示路由列表；创建路由
- 展示 destination、next_hop、优先级等（以 `NetworkRoute` 为准）

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 读/写权限 | 对应 networks scope | `403 FORBIDDEN` |
| POST 请求体 | 满足 `CreateNetworkRouteRequest` | `400 BAD_REQUEST` |
| 关联 `vpc_id` | 有效 | 具体 `404`/`422` **当前 create 未声明** |
| 路由冲突 | 无重复 destination | `409 CONFLICT`（create 已声明） |

## 操作可用性矩阵

| 操作 | 只读用户 | 网络管理员 |
|---|---|---|
| 列表 | ✅ | ✅ |
| 创建 | ❌ | ✅ |
| 查看单条详情 | ❌ | ❌（无 GET by id） |
| 删除单条 | ❌ | ❌（待 YAML） |

## 接口冻结规则

### `GET /api/v1/networks/routes`

- 成功：`200 + NetworkRouteListResponse`
- 错误：`401`、`403`
- Query 可选：`vpc_id`、`limit`、`cursor`

### `POST /api/v1/networks/routes`

- 成功：`201 + NetworkRoute`
- 错误：`400`、`401`、`403`、`409`

## 待补边界

- `GET /networks/routes/{route_id}` — **TODO-YAML**
- `DELETE /networks/routes/{route_id}` — **TODO-YAML**
- 路由更新 PATCH — **TODO-YAML**

## 验收标准

- [ ] 与 Phase 2 路由路径一致
- [ ] 未声明的 DELETE/GET by id 不写进冻结事实
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
