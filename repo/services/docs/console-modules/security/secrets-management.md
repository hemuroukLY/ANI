# 密钥管理

## 页面定位

租户 **Secret 元数据** CRUD 与工作负载绑定页；API **不返回明文值**（见 YAML description）。

## 文档管理规则

- 本文是 **密钥管理** 的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`
- 不得把规划路径或 handler stub 写成已实现
- PRD/SPEC 同步规则见 `docs/console-modules/governance/module-delivery-workflow.md` §3.5

## Core 层要求

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| GET | `/api/v1/secrets` | `listSecrets` | `scope:secrets:read` |
| POST | `/api/v1/secrets` | `createSecret` | `scope:secrets:create` |
| GET | `/api/v1/secrets/{secret_id}` | `getSecret` | `scope:secrets:read` |
| DELETE | `/api/v1/secrets/{secret_id}` | `deleteSecret` | `scope:secrets:delete` |
| POST | `/api/v1/secrets/{secret_id}/bindings` | `bindSecret` | `scope:secrets:bind` |

Schema：`Secret`、`SecretListResponse`、`SecretCreateRequest`、`SecretBinding`。

<!-- ADDED-TO-YAML: Core Secrets 组 -->

## Services 层要求

无 Services 路径。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 模块读/写权限 | 对应 RBAC scope | `403 FORBIDDEN` |

写操作 POST/PUT 的 `idempotency_key` 与 422 口径 — **待 YAML 冻结后按 ../governance/module-delivery-workflow.md §2.10 补充**。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员/运维 | 说明 |
|---|---|---|---|
| 列表/详情 | 可用 | 可用 | 无明文 |
| 创建/删除 | 不可用 | 管理员 | create/delete |
| 绑定工作负载 | 不可用 | 运维 | bindSecret |

## 页面职责

- 占位 UI + 明确 YAML/OpenAPI 缺口（若适用）
- 跳转关联模块（见「相关模块」）
- 不把 BOSS/平台运营能力写入 Console 冻结契约

## 接口冻结规则

### `GET /api/v1/secrets`

- operationId: `listSecrets`
- 成功：`200 + SecretListResponse`
- RBAC：`scope:secrets:read`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

### `POST /api/v1/secrets`

- operationId: `createSecret`
- 成功：`201 + Secret`
- RBAC：`scope:secrets:create`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

### `GET /api/v1/secrets/{secret_id}`

- operationId: `getSecret`
- 成功：`200 + Secret`
- RBAC：`scope:secrets:read`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

### `DELETE /api/v1/secrets/{secret_id}`

- operationId: `deleteSecret`
- 成功：`200 + Secret`
- RBAC：`scope:secrets:delete`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文

### `POST /api/v1/secrets/{secret_id}/bindings`

- operationId: `bindSecret`
- 成功：`201 + SecretBinding`
- RBAC：`scope:secrets:bind`
- 错误：`401`、`403`（及 YAML 已声明的 404/409 等）
- 租户边界：认证上下文


## 待补边界

- Secret 值轮换 PUT — **YAML 未声明**
- K8s provider 真实注入 — handler 边界

## 相关模块

- `tenant/security-identity-overview.md`
- `security/crypto-sm.md`

## 验收标准

- [ ] 路径与 OpenAPI 权威源一致（或明确 TODO-YAML / N/A）
- [ ] 正文不把 handler stub 写成已实现
- [ ] 含创建前置条件与操作可用性矩阵
- [ ] PRD/SPEC/HTML 摘要已同步
