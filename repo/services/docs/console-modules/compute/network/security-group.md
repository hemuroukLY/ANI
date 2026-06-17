# 网络管理 — 安全组

## 页面定位

`安全组` 是网络管理下安全组资源的独立明细页。

父级：`network-management.md`。

## 文档管理规则

- 本文是安全组子模块主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`

## Core 层要求

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| GET | `/api/v1/networks/security-groups` | `listNetworkSecurityGroups` | `scope:networks:read` |
| POST | `/api/v1/networks/security-groups` | `createNetworkSecurityGroup` | `scope:networks:create` |
| GET | `/api/v1/networks/security-groups/{security_group_id}` | `getNetworkSecurityGroup` | `scope:networks:read` |
| DELETE | `/api/v1/networks/security-groups/{security_group_id}` | `deleteNetworkSecurityGroup` | `scope:networks:delete` |

Schema：`NetworkSecurityGroup`（含内嵌 `rules[]`）。

POST 创建必须带 `idempotency_key`。

## 页面职责

- 安全组 CRUD；规则以摘要展示（数量 + 协议/端口摘要）
- 不提供独立 rules 子资源 CRUD（待补）

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 读/写权限 | 对应 networks scope | `403 FORBIDDEN` |
| POST 必填 | `name`、`idempotency_key` | `400 BAD_REQUEST` |
| `rules[]` 若提供 | 满足 YAML schema | `400 BAD_REQUEST` |

**当前 YAML 未在 create 声明 `404`/`422`**。

## 操作可用性矩阵

| 操作 | 只读用户 | 网络管理员 |
|---|---|---|
| 列表/详情 | ✅ | ✅ |
| 创建/删除 | ❌ | ✅ |
| 编辑 rules 子资源 | ❌ | ❌（待 YAML） |

## 接口冻结规则

### `GET /api/v1/networks/security-groups`

- 成功：`200 + NetworkSecurityGroupListResponse`
- 错误：`401`、`403`

### `POST /api/v1/networks/security-groups`

- 成功：`201 + NetworkSecurityGroup`
- 错误：`400`、`401`、`403`

### `GET /api/v1/networks/security-groups/{security_group_id}`

- 成功：`200 + NetworkSecurityGroup`
- 错误：`401`、`403`、`404`

### `DELETE /api/v1/networks/security-groups/{security_group_id}`

- 成功：`200 + NetworkSecurityGroup`
- 错误：`401`、`403`、`404`

## 待补边界

- 安全组规则子资源 CRUD — **TODO-YAML**
- 安全组绑定/解绑实例或网卡 — **TODO-YAML**
- 安全组更新 PATCH — **当前 YAML 未声明**

## 验收标准

- [ ] 不把 rules 子路径写成已冻结
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
