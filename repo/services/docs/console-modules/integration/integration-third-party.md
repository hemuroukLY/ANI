# 接入 — 第三方集成

## 页面定位

第三方业务系统（通用 provider）集成列表与创建。

本页属于 **Services / Integrations** 资源域。

## 文档管理规则

- 本文是第三方集成主维护源
- 一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- Bot 子类型见 `integration-bot.md`
- TASK：`TASK-SVC-016`

## Services 层要求

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/svc/integrations` | `listIntegrations` |
| POST | `/api/v1/svc/integrations` | `createIntegration` |

列表项字段（以 YAML 为准）：`id`、`name`、`provider`、`status`（active/inactive/error）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 集成读/写权限 | 对应 scope | `403 FORBIDDEN` |
| 幂等键 | POST 携带 `idempotency_key` | `400` |
| provider | 枚举内合法值 | `400` |

**当前 YAML 未声明 `422`**（除非 createIntegration 已声明 PreconditionFailed）。

## 页面职责

- 集成列表、创建、状态展示（active/inactive/error）
- 详情抽屉：展示 provider 配置摘要（不暴露 secret 明文）
- 更新/删除集成 — **YAML 未声明** `{integration_id}` 路径

## 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 列表/详情 | 可用 | 可用 |
| 创建 | 不可用 | 可用 |
| 更新/删除 | 不可用 | 不可用（待补） |

## 接口冻结规则

### `GET /api/v1/svc/integrations`

- 成功：`200 + IntegrationListResponse`
- 错误：`401`、`403`
- 分页：遵循 YAML query

### `POST /api/v1/svc/integrations`

- 成功：`201 + Integration`
- 错误：`401`、`403`、`400`
- 幂等：`idempotency_key`

## 待补边界

- `GET/PATCH/DELETE /integrations/{integration_id}` — **YAML 未声明**
- 集成连通性测试 API — 待补
- OAuth 第三方授权回调 — 待 Gateway/Services 规划

## 相关模块

- `integration-bot.md`（Bot 为 integrations 子类型）
- `integration-webhook-overview.md`
- `open-integration-overview.md`

## 验收标准

- [ ] list/create 路径与 services/v1.yaml 一致
- [ ] 未自造 integration_id 写路径
- [ ] secret 字段不在 UI 明文展示
- [ ] handler stub ≠ 已实现
