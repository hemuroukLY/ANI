# SPEC: Console api-key-audit

> Phase 3 规划详化 · Source: `tasks/modules/prd/console/tenant/prd-console-api-key-audit.md`  
> 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md` · Revised: 2026-06-17

## 1. Summary

Core Auth Key 调用审计 list（规划）；只读。

## 2. Planned Facts（非冻结）

| Method | Path | operationId | 成功 |
|---|---|---|---|
| GET | `/api/v1/auth/api-keys/{key_id}/audit-events` | `listApiKeyAuditEvents` | 200 |

Schemas: `ApiKeyAuditEvent`, `ApiKeyAuditEventListResponse`.

## 3. 主维护源

- `docs/console-modules/tenant/api-key-audit.md`

OpenAPI **未合入** ≠ handler 已实现。
