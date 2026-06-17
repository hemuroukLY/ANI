# SPEC: Console tenant-webhook-ops

> Source: `tasks/modules/prd/console/tenant/prd-console-tenant-webhook-ops.md`  
> Revised: 2026-06-17

## 1. Summary

租户 Webhook 的**投递日志**查询与运维视图，扩展 `tenant-management.md` 的 Webhook CRUD。属于 **Services / Tenant** 运维子能力。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | 权限 |
|---|---|---|---|---|
| GET | `/api/v1/svc/tenant/webhooks/{webhook_id}/deliveries` | `listWebhookDeliveries` | `200 + WebhookDeliveryListResponse` | 租户管理员或等价 scope |

Query：`limit`、`cursor`、可选 `status` filter（以 OpenAPI 为准）。

基础 CRUD：`tenant-management.md`（含 DELETE webhook，Phase 2 前已声明）。

### 2.3 Verified Schemas

- `WebhookDeliveryListResponse`（`status`: pending/success/failed/retrying，以 YAML 为准）

## 3. Page Scope

- Webhook 详情下「投递记录」Tab：分页、状态筛选
- 展示投递时间、HTTP 状态、重试次数、错误摘要
- 空列表：`200` + `items: []`

## 4. Non-Goals

- `POST .../deliveries/{id}/retry` — **YAML 未声明**
- 投递 payload 明文重看 — 安全策略待产品定稿
- 批量导出投递日志 — 待补
- **当前 YAML 未声明 deliveries 相关 `422`**

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 租户 Webhook 读权限 | 租户管理员或等价 scope | `403 FORBIDDEN` |
| Webhook 存在 | 当前租户内 | `404 NOT_FOUND` |

## 6. 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 查看投递日志 | ❌ | ✅ |
| 手动重试投递 | ❌ | ❌（待补） |
| 删除 Webhook | 跳转 tenant-management | ✅ |

## 7. 主维护源

- `docs/console-modules/tenant/tenant-webhook-ops.md`
- 相关：`tenant-management.md`、`integration/integration-webhook-overview.md`
- TASK：`TASK-SVC-015`

## 8. Handler 验收（Services 团队）

```bash
curl -H "Authorization: Bearer $TOKEN" "$BASE/api/v1/svc/tenant/webhooks/{webhook_id}/deliveries"
```

- 分页：遵循 YAML `limit`/`cursor`
- 未把 retry 写成冻结操作
- 与 tenant webhooks CRUD 分层清晰
- OpenAPI 已声明 ≠ handler 已实现
