# Services Handler 实现指南（文档交付）

> **生成日期**：2026-06-17  
> **用途**：Services 团队按 `services/v1.yaml` 实现 Phase 2 handler 时的技术契约与验收 curl。  
> **权威源**：`ANI-main/repo/api/openapi/services/v1.yaml`  
> **验收记录**：`tasks/phase2/PHASE2-HANDLER-ACCEPTANCE-RECORD.md`  
> **任务索引**：`tasks/execution/SERVICES-TEAM-TASKS.md`

## 文档管理规则

1. 路径前缀 **`/api/v1/svc/*`**；不得写入 Core 资源路径。
2. 租户边界从认证上下文提取；handler 不得要求 body 传 `tenant_id`。
3. 写操作须校验 `idempotency_key`（YAML `required` 处）。
4. Services `422`：优先使用 YAML 已举例 code（如 `MODEL_NOT_READY`）；见 `../governance/module-delivery-workflow.md` §2.10。
5. **不实现** deprecated：`/gpu-containers*`、`/sandboxes*`。
6. OpenAPI 已声明 ≠ handler 已实现。

---

## TASK-SVC-013 — 推理 logs / test / policies

| 方法 | 路径 | operationId | 成功 | 模块详文 |
|---|---|---|---|---|
| GET | `/inference-services/{service_id}/logs` | `getInferenceServiceLogs` | 200 | `inference-observability.md` |
| POST | `/inference-services/{service_id}/test` | `testInferenceService` | 200 | `inference-call-test.md` |
| PUT | `/inference-services/{service_id}/policies` | `updateInferenceServicePolicies` | 200 | `inference-rate-limit-policy.md` |

### logs

- Query：`limit`（默认 100，max 500）、`cursor`
- 响应：`InferenceServiceLogListResponse`
- 服务不存在：`404`

### test

- Request：`InferenceServiceTestRequest`（含 `idempotency_key`、测试 payload，以 YAML 为准）
- 响应：`InferenceServiceTestResponse`
- 服务非 `running`：建议 **422**（具体 code 待 Services 冻结举例）
- 服务不存在：`404`

### policies

- Request：`InferenceServicePoliciesUpdateRequest`（`idempotency_key` + 策略字段）
- 响应：更新后的 `InferenceService` 或 YAML 声明体
- 成功：`200`

### 验收

```bash
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/inference-services/$SVC_ID/logs?limit=100"

curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'"}' \
  "$BASE/api/v1/svc/inference-services/$SVC_ID/test"

curl -sS -X PUT -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","rate_limit_rpm":1000}' \
  "$BASE/api/v1/svc/inference-services/$SVC_ID/policies"
```

---

## TASK-SVC-014 — 知识库 citations / sessions / permissions

| 方法 | 路径 | operationId | 模块详文 |
|---|---|---|---|
| GET | `/knowledge-bases/{kb_id}/citations` | `listKnowledgeBaseCitations` | `kb-source-citation.md` |
| GET | `/knowledge-bases/{kb_id}/sessions` | `listKnowledgeBaseSessions` | `kb-chat-history.md` |
| PUT | `/knowledge-bases/{kb_id}/permissions` | `updateKnowledgeBasePermissions` | `kb-permissions.md` |

### 要点

- KB 不存在：`404`；无读权限：`403`
- permissions PUT 必填 `idempotency_key`、`readers[]`、`editors[]`
- sessions/citations 空列表 → `200` + `items: []`

### 验收

```bash
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/knowledge-bases/$KB_ID/citations"

curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/knowledge-bases/$KB_ID/sessions"

curl -sS -X PUT -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","readers":["u1"],"editors":["u2"]}' \
  "$BASE/api/v1/svc/knowledge-bases/$KB_ID/permissions"
```

---

## TASK-SVC-015 — 租户 member / role / webhook deliveries

| 方法 | 路径 | operationId | 模块详文 |
|---|---|---|---|
| GET | `/tenant/members/{member_id}` | `getTenantMember` | `tenant-member-detail.md` |
| PUT | `/tenant/roles/{role_id}` | `updateTenantRole` | `role-permission-edit.md` |
| GET | `/tenant/webhooks/{webhook_id}/deliveries` | `listWebhookDeliveries` | `tenant-webhook-ops.md` |

### 要点

- `POST /tenant/webhooks` 422 验通（Phase 2 YAML 修补）— URL 无效等前置失败
- deliveries 支持 `limit`/`cursor`、可选 `status` filter

### 验收

```bash
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/tenant/members/$MEMBER_ID"

curl -sS -X PUT -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","permissions":["scope:instances:read"]}' \
  "$BASE/api/v1/svc/tenant/roles/$ROLE_ID"

curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/svc/tenant/webhooks/$WEBHOOK_ID/deliveries?limit=50"
```

---

## TASK-SVC-016 — integrations + bots

| 方法 | 路径 | operationId | 模块详文 |
|---|---|---|---|
| GET | `/integrations` | `listIntegrations` | `integration-third-party.md` |
| POST | `/integrations` | `createIntegration` | `integration-third-party.md` |
| POST | `/integrations/bots` | `createIntegrationBot` | `integration-bot.md` |

### 要点

- POST 必填 `idempotency_key`
- createIntegrationBot：`platform` ∈ wecom|dingtalk（以 YAML enum 为准）
- **未声明**：GET/PATCH/DELETE `/integrations/{id}`

### 验收

```bash
curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/integrations"

curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","name":"crm","provider":"custom"}' \
  "$BASE/api/v1/svc/integrations"

curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","platform":"wecom","webhook_url":"https://example.com/hook"}' \
  "$BASE/api/v1/svc/integrations/bots"
```

---

## 刻意不实现（Services Phase 2）

| 项 | 原因 |
|---|---|
| `/gpu-containers*`、`/sandboxes*` | deprecated |
| inference events | P0 草案 **待定**，见 `openapi-p0-yaml-draft.md` |
| integrations `{id}` CRUD | YAML 未声明 |

---

## 相关文件

- `tasks/phase2/PHASE2-HANDLER-ACCEPTANCE-RECORD.md`
- `tasks/execution/SERVICES-TEAM-TASKS.md`
- `docs/console-modules/governance/schema-completion-tracker.md`
