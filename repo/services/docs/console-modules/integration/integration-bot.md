# 接入 — 企微 / 钉钉 Bot

## 页面定位

创建与管理**企业微信 / 钉钉** Bot 集成。

本页属于 **Services / Integrations** 子类型（Bot）。

## 文档管理规则

- 本文是 Bot 集成子页的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- 通用 integrations 列表见 `integration-third-party.md`
- TASK：`TASK-SVC-016`

## Services 层要求

| 方法 | 路径 | operationId |
|---|---|---|
| POST | `/api/v1/svc/integrations/bots` | `createIntegrationBot` |

请求必填：`idempotency_key`、`platform`（wecom/dingtalk）、`webhook_url`；可选 `name`（以 OpenAPI 为准）。

列表：复用 `GET /api/v1/svc/integrations`，客户端按 `provider`/`platform` 筛选 Bot 项。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 集成写权限 | 租户管理员或等价 scope | `403 FORBIDDEN` |
| 幂等键 | POST 携带 `idempotency_key` | `400`（若必填缺失） |
| webhook_url | 合法 URL 格式 | `400`（校验失败） |

**`createIntegrationBot` YAML 已声明 `422 PreconditionFailed`**（Services v1.yaml Phase 2 2026-06-17）；具体 `code` 值待 Services description 补充后再写成冻结契约。

## 页面职责

- Bot 创建表单（平台选择、Webhook URL、名称）
- Bot 列表（复用 integrations list + 筛选）
- Bot 删除 — **YAML 未声明** DELETE

## 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 查看 Bot 列表 | 可用 | 可用 |
| 创建 Bot | 不可用 | 可用 |
| 删除/更新 Bot | 不可用 | 不可用（待补） |

## 接口冻结规则

### `POST /api/v1/svc/integrations/bots`

- 成功：`201 + Integration`（或 YAML 声明的响应 schema）
- 错误：`400`、`401`、`403`、`422`
- 幂等：相同 `idempotency_key` 行为以 YAML 为准

### `GET /api/v1/svc/integrations`（列表筛选）

- 成功：`200 + IntegrationListResponse`
- 错误：`401`、`403`

## 待补边界

- `DELETE /integrations/bots/{id}` — **YAML 未声明**
- `PATCH` 更新 Bot — **YAML 未声明**
- Bot 消息投递日志 — 待补（或复用 tenant webhook deliveries 口径）

## 相关模块

- `integration-third-party.md`
- `integration-webhook-overview.md`
- `open-integration-overview.md`

## 验收标准

- [ ] POST bots 路径与 services/v1.yaml 一致
- [ ] 创建表单校验 idempotency_key
- [ ] 未把 DELETE 写成可用操作
- [ ] handler stub ≠ 已实现
