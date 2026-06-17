# Phase 3 — 模型加密与密钥绑定 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`inference/model-encryption.md`  
> **TASK**：`TASK-SVC-018` 子项（模型增强）  
> **原则**：Services 模型资源绑定 Core Encryption Key；密钥 CRUD 不重复声明。

---

## 1. 与 Core Encryption / ModelVersion 的边界

| 维度 | Core Encryption | 模型加密绑定（本草案） |
|---|---|---|
| 路径 | `/api/v1/encryption/keys*` | `/api/v1/svc/models/{model_id}/encryption` |
| 职责 | 密钥生命周期、seal/unseal | 模型 artifact ↔ `encryption_key_id` 绑定 |
| Schema | `EncryptionKey` | `ModelEncryptionBinding` |
| 写操作 | create/rotate/revoke key | PUT binding（引用已有 key_id） |

`ModelVersion.is_encrypted`（已冻结字段）由 binding 成功后 Services 回写；Console 以 binding API 为配置入口。

详见 `security/crypto-sm.md`。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [Models]
x-ani-rbac-scope:
  read:  scope:models:read
  write: scope:models:write
```

---

## 3. Schemas（草案）

### ModelEncryptionBinding

```yaml
ModelEncryptionBinding:
  type: object
  required: [model_id, status, updated_at]
  properties:
    model_id:           { type: string, format: uuid }
    version_id:         { type: string, format: uuid, nullable: true, description: 空表示默认/全部版本策略 }
    encryption_key_id:  { type: string, format: uuid, nullable: true }
    algorithm:          { type: string, enum: [aes256_gcm, sm4_gcm], nullable: true }
    status:             { type: string, enum: [unbound, binding, bound, failed, rotating] }
    bound_at:           { type: string, format: date-time, nullable: true }
    updated_at:         { type: string, format: date-time }
    error_message:      { type: string, nullable: true }
```

### UpdateModelEncryptionRequest

```yaml
UpdateModelEncryptionRequest:
  type: object
  required: [idempotency_key]
  properties:
    idempotency_key:    { type: string, format: uuid }
    encryption_key_id:  { type: string, format: uuid, nullable: true, description: null 表示解绑 }
    version_id:         { type: string, format: uuid, nullable: true }
    algorithm:          { type: string, enum: [aes256_gcm, sm4_gcm], nullable: true }
```

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/svc/models/{model_id}/encryption`

- operationId: `getModelEncryptionBinding`
- 成功：`200 + ModelEncryptionBinding`
- 错误：`401`、`403`、`404`（model）

### 4.2 `PUT /api/v1/svc/models/{model_id}/encryption`

- operationId: `updateModelEncryptionBinding`
- Body: `UpdateModelEncryptionRequest`
- 成功：`200 + ModelEncryptionBinding`（同步完成）或 `202 + AsyncTask`（re-encrypt 大 artifact 时）
- 错误：
  - `400`：缺 idempotency_key
  - `404`：model / key 不存在
  - `422`：model 非 ready — `MODEL_NOT_READY`；key 已 revoke — `ENCRYPTION_KEY_REVOKED`；版本不存在 — `MODEL_VERSION_NOT_FOUND`

### 4.3 Core 只读引用（绑定前选 key）

- `GET /api/v1/encryption/keys` — `listEncryptionKeys`
- `GET /api/v1/encryption/keys/{key_id}` — `getEncryptionKey`

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 模型管理员 | 说明 |
|---|---|---|---|
| 查看绑定状态 | 可用 | 可用 | GET |
| 绑定/解绑/轮换 | 不可用 | 可用 | PUT |
| 跳转密钥管理 | 跳转 | 跳转 | crypto-sm |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/models/$MODEL_ID/encryption"

curl -sS -X PUT -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","encryption_key_id":"'"$KEY_ID"'","algorithm":"sm4_gcm"}' \
  "$BASE/api/v1/svc/models/$MODEL_ID/encryption"
```

---

## 7. 评审检查清单

- [ ] 不重复声明 Core encryption keys CRUD
- [ ] PUT 必填 `idempotency_key`
- [ ] 422 与 inference `MODEL_NOT_READY` 口径一致
- [ ] 合入后更新 `model-encryption.md`

---

## 相关文件

- `docs/console-modules/inference/model-encryption.md`
- `docs/console-modules/inference/model-center.md`
- `docs/console-modules/security/crypto-sm.md`
