# Phase 3 — API Key 审计 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`tenant/api-key-audit.md`  
> **TASK**：`TASK-CORE-015` 子项 / Phase 3 §4  
> **原则**：Core Auth 子资源；与 `api-key-management.md` CRUD 分路径；只读 list。

---

## 1. 与 API Key CRUD / 平台审计日志的边界

| 维度 | API Key 管理 | API Key 审计（本草案） | 平台审计日志 |
|---|---|---|---|
| 路径 | `/auth/api-keys*` | `/auth/api-keys/{id}/audit-events` | 规划 `/audit/events` |
| 事件 | 创建/吊销 | **使用该 Key 的 API 调用** | 租户内全平台操作 |
| 粒度 | Key 对象 | per-key 调用记录 | actor + resource |
| 模块 | `api-key-management.md` | 本文 | `audit-log.md` |

Console：从 Key 列表「查看审计」进入本模块；不替代 BOSS 取证导出。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [Auth]
x-ani-rbac-scope:
  read: scope:auth:api-keys:audit:read
```

与 `api-keys:read` 分离，便于只授审计只读角色（**评审可合并为 admin-only**）。

---

## 3. Schemas（草案）

### ApiKeyAuditEvent

```yaml
ApiKeyAuditEvent:
  type: object
  required: [id, api_key_id, occurred_at, method, path_template, status_code]
  properties:
    id:            { type: string, format: uuid }
    api_key_id:    { type: string, format: uuid }
    occurred_at:   { type: string, format: date-time }
    source_ip:     { type: string, nullable: true }
    method:        { type: string, enum: [GET, POST, PUT, PATCH, DELETE] }
    path_template: { type: string, description: 归一化路径，如 /api/v1/instances/{id} }
    status_code:   { type: integer, minimum: 100, maximum: 599 }
    latency_ms:    { type: number, nullable: true }
    user_agent:    { type: string, nullable: true }
    region:        { type: string, nullable: true }
    error_code:    { type: string, nullable: true, description: 4xx/5xx 业务 code }
```

### ApiKeyAuditEventListResponse

```yaml
ApiKeyAuditEventListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/ApiKeyAuditEvent' } }
    next_cursor: { type: string, nullable: true }
```

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/auth/api-keys/{key_id}/audit-events`

- operationId: `listApiKeyAuditEvents`
- Path params: `key_id`
- Query: `start_time`、`end_time`（ISO8601）、`status_code`（int）、`path_prefix`、`limit`（default 50, max 200）、`cursor`
- 成功：`200 + ApiKeyAuditEventListResponse`
- 错误：`401`、`403`、`404`（key 不存在或不可见）

### 4.2 `GET /api/v1/auth/api-keys/{key_id}/audit-events/{event_id}`（3b）

- operationId: `getApiKeyAuditEvent`
- 成功：`200 + ApiKeyAuditEvent`（含可选 request_id 关联）
- 3a 可不实现，列表足够

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | Key 管理员 | 说明 |
|---|---|---|---|
| 查看 Key 审计列表 | 不可用 | 可用 | GET list |
| 从 Key 管理跳转 | — | 可用 | key_id 上下文 |
| 平台级导出 | 不可用 | 不可用 | BOSS |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -H "X-API-Key: $ADMIN_KEY" \
  "$BASE/api/v1/auth/api-keys/$KEY_ID/audit-events?limit=20&start_time=2026-06-01T00:00:00Z"
```

---

## 7. 评审检查清单

- [ ] 与 `DELETE /auth/api-keys/{id}` 无 path 冲突
- [ ] 不记录 `key_value` 明文
- [ ] 与 `audit-log.md` 事件字段可对齐映射
- [ ] 合入后更新 `api-key-audit.md`

---

## 相关文件

- `docs/console-modules/tenant/api-key-audit.md`
- `docs/console-modules/tenant/api-key-management.md`
- `docs/console-modules/security/audit-log.md`
