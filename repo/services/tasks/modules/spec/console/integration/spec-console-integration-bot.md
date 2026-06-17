# SPEC: Console integration-bot

> Source: `tasks/modules/prd/console/integration/prd-console-integration-bot.md`  
> Revised: 2026-06-17

## 1. Summary

创建与管理**企业微信 / 钉钉** Bot 集成。属于 **Services / Integrations** 子类型（Bot）。通用 integrations 列表见 `integration-third-party.md`。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | 权限 |
|---|---|---|---|---|
| POST | `/api/v1/svc/integrations/bots` | `createIntegrationBot` | `201 + Integration` | 租户管理员或等价 scope |
| GET | `/api/v1/svc/integrations` | `listIntegrations` | `200 + IntegrationListResponse` | 对应 scope（客户端按 `provider`/`platform` 筛选 Bot 项） |

### 2.3 Verified Schemas

- 创建请求必填：`idempotency_key`、`platform`（wecom/dingtalk）、`webhook_url`；可选 `name`（以 OpenAPI 为准）
- `Integration` / `IntegrationListResponse`（以 YAML 为准）

## 3. Page Scope

- Bot 创建表单（平台选择、Webhook URL、名称）
- Bot 列表（复用 integrations list + 筛选）
- 创建表单校验 `idempotency_key`

## 4. Non-Goals

- `DELETE /integrations/bots/{id}` — **YAML 未声明**
- `PATCH` 更新 Bot — **YAML 未声明**
- Bot 消息投递日志 — 待补（或复用 tenant webhook deliveries 口径）
- **当前 YAML 未声明 `422`**（除非 OpenAPI 已为 createIntegrationBot 声明 PreconditionFailed）

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 集成写权限 | 租户管理员或等价 scope | `403 FORBIDDEN` |
| 幂等键 | POST 携带 `idempotency_key` | `400`（若必填缺失） |
| webhook_url | 合法 URL 格式 | `400`（校验失败） |

## 6. 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 查看 Bot 列表 | ✅ | ✅ |
| 创建 Bot | ❌ | ✅ |
| 删除/更新 Bot | ❌ | ❌（待补） |

## 7. 主维护源

- `docs/console-modules/integration/integration-bot.md`
- 相关：`integration-third-party.md`、`integration-webhook-overview.md`、`open-integration-overview.md`
- TASK：`TASK-SVC-016`

## 8. Handler 验收（Services 团队）

```bash
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"idempotency_key":"...","platform":"wecom","webhook_url":"https://..."}' \
  "$BASE/api/v1/svc/integrations/bots"
curl -H "Authorization: Bearer $TOKEN" "$BASE/api/v1/svc/integrations"
```

- 幂等：相同 `idempotency_key` 行为以 YAML 为准
- 未把 DELETE 写成可用操作
- OpenAPI 已声明 ≠ handler 已实现
