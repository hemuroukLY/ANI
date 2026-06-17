# 租户管理 — Webhook 投递运维

## 页面定位

租户 Webhook 的**投递日志**查询与运维视图，扩展 `tenant-management.md` 的 Webhook CRUD。

本页属于 **Services / Tenant** 运维子能力。

## 文档管理规则

- 本文是 Webhook 投递运维子页的主维护源
- 一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- Webhook CRUD 见 `tenant-management.md`
- TASK：`TASK-SVC-015`

## Services 层要求

<!-- ADDED-TO-YAML: Phase 2 2026-06-17 -->

| 方法 | 路径 | operationId |
|---|---|---|
| GET | `/api/v1/svc/tenant/webhooks/{webhook_id}/deliveries` | `listWebhookDeliveries` |

Query：`limit`、`cursor`、可选 `status` filter（以 OpenAPI 为准）。

响应：`WebhookDeliveryListResponse`（`status`: pending/success/failed/retrying）。

基础 CRUD：`tenant-management.md`（含 DELETE webhook，Phase 2 前已声明）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 租户 Webhook 读权限 | 租户管理员或等价 scope | `403 FORBIDDEN` |
| Webhook 存在 | 当前租户内 | `404 NOT_FOUND` |

本页 list 无写 body。**当前 YAML 未声明 deliveries 相关 `422`**。

## 页面职责

- Webhook 详情下「投递记录」Tab：分页、状态筛选
- 展示投递时间、HTTP 状态、重试次数、错误摘要
- 重试 — **YAML 未声明** POST retry

## 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 查看投递日志 | 不可用 | 可用 |
| 手动重试投递 | 不可用 | 不可用（待补） |
| 删除 Webhook | 跳转 tenant-management | 可用 |

## 接口冻结规则

### `GET /api/v1/svc/tenant/webhooks/{webhook_id}/deliveries`

- 成功：`200 + WebhookDeliveryListResponse`
- 错误：`401`、`403`、`404`
- 分页：遵循 YAML `limit`/`cursor`
- 空列表：`200` + `items: []`

## 待补边界

- `POST .../deliveries/{id}/retry` — **YAML 未声明**
- 投递 payload 明文重看 — 安全策略待产品定稿
- 批量导出投递日志 — 待补

## 相关模块

- `tenant-management.md`
- `integration/integration-webhook-overview.md`（接入域入口说明）

## 验收标准

- [ ] deliveries 路径与 Phase 2 services/v1.yaml 一致
- [ ] 未把 retry 写成冻结操作
- [ ] handler stub ≠ 已实现
- [ ] 与 tenant webhooks CRUD 分层清晰
