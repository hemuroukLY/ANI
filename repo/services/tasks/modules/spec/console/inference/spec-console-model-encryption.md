# SPEC: Console model-encryption

> Technical specification — Phase 3 规划详化  
> Source: `tasks/modules/prd/console/inference/prd-console-model-encryption.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md`  
> Revised: 2026-06-17

## 1. Summary

模型 ↔ Encryption Key 绑定（规划）；Services PUT/GET；Core keys 只读引用。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- Services: `ANI-main/repo/api/openapi/services/v1.yaml`
- Core 只读: `ANI-main/repo/api/openapi/v1.yaml`

### 2.2 Planned Paths（Services）

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/svc/models/{model_id}/encryption` | `getModelEncryptionBinding` | 200 | `scope:models:read` |
| PUT | `/api/v1/svc/models/{model_id}/encryption` | `updateModelEncryptionBinding` | 200/202 | write |

### 2.3 Planned Schemas

- `ModelEncryptionBinding`（status: unbound|binding|bound|failed|rotating）
- `UpdateModelEncryptionRequest`

## 3. Page Scope

- 绑定状态、算法、跳转 crypto-sm
- **Non-Goals**：Core keys CRUD

## 4. 创建前置条件

`MODEL_NOT_READY`、`ENCRYPTION_KEY_REVOKED` — 见主维护源。

## 5. Handler 验收（合入 YAML 后）

```bash
curl .../svc/models/$ID/encryption
curl -X PUT .../svc/models/$ID/encryption -d '{"idempotency_key":"...","encryption_key_id":"..."}'
```

## 6. 主维护源

- `docs/console-modules/inference/model-encryption.md`
- 草案：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-encryption-draft.md`

OpenAPI **未合入** ≠ handler 已实现。
