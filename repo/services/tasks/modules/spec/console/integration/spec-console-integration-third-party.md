# SPEC: Console integration-third-party

> Source: `tasks/modules/prd/console/integration/prd-console-integration-third-party.md`  
> Revised: 2026-06-17

## 1. Summary

第三方业务系统（通用 provider）集成列表与创建。属于 **Services / Integrations** 资源域。Bot 子类型见 `integration-bot.md`。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功响应 | 权限 |
|---|---|---|---|---|
| GET | `/api/v1/svc/integrations` | `listIntegrations` | `200 + IntegrationListResponse` | 对应 scope |
| POST | `/api/v1/svc/integrations` | `createIntegration` | `201 + Integration` | 租户管理员或等价 scope |

### 2.3 Verified Schemas

- `Integration` 列表项字段（以 YAML 为准）：`id`、`name`、`provider`、`status`（active/inactive/error）
- `IntegrationListResponse`（以 YAML 为准）
- 创建请求含 `idempotency_key`、`provider`（枚举内合法值）

## 3. Page Scope

- 集成列表、创建、状态展示（active/inactive/error）
- 详情抽屉：展示 provider 配置摘要（不暴露 secret 明文）
- 分页：遵循 YAML query

## 4. Non-Goals

- `GET/PATCH/DELETE /integrations/{integration_id}` — **YAML 未声明**
- 集成连通性测试 API — 待补
- OAuth 第三方授权回调 — 待 Gateway/Services 规划
- **当前 YAML 未声明 `422`**（除非 createIntegration 已声明 PreconditionFailed）

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 集成读/写权限 | 对应 scope | `403 FORBIDDEN` |
| 幂等键 | POST 携带 `idempotency_key` | `400` |
| provider | 枚举内合法值 | `400` |

## 6. 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 列表/详情 | ✅ | ✅ |
| 创建 | ❌ | ✅ |
| 更新/删除 | ❌ | ❌（待补） |

## 7. 主维护源

- `docs/console-modules/integration/integration-third-party.md`
- 相关：`integration-bot.md`、`integration-webhook-overview.md`、`open-integration-overview.md`
- TASK：`TASK-SVC-016`

## 8. Handler 验收（Services 团队）

```bash
curl -H "Authorization: Bearer $TOKEN" "$BASE/api/v1/svc/integrations"
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"idempotency_key":"...","provider":"...","name":"..."}' \
  "$BASE/api/v1/svc/integrations"
```

- secret 字段不在 UI 明文展示
- 未自造 integration_id 写路径
- OpenAPI 已声明 ≠ handler 已实现
