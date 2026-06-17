# SPEC: Console prompt-injection-guard

> Phase 3 规划详化  
> Source: `tasks/modules/prd/console/ai-native/prd-console-prompt-injection-guard.md`  
> OpenAPI 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-prompt-injection-guard-draft.md`  
> Revised: 2026-06-17

## 1. Summary

租户级 Prompt 注入/越狱扫描策略；与 inference RPM 策略分 path；运行时写 Audit `policy_hit`。

## 2. Planned Facts（草案 · 非冻结）

### 2.1 Authority Source（合入后）

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Planned Paths

| Method | Path | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/api/v1/svc/agent/guard-policies` | `listPromptGuardPolicies` | 200 | `scope:agent-guard:read` |
| GET | `/api/v1/svc/agent/guard-policies/{policy_id}` | `getPromptGuardPolicy` | 200 | read |
| PUT | `/api/v1/svc/agent/guard-policies/{policy_id}` | `updatePromptGuardPolicy` | 200 | write |
| PUT | `/api/v1/svc/inference-services/{service_id}/guard-binding` | `updateInferenceServiceGuardBinding` | 200 | write |
| POST | `/api/v1/svc/agent/guard-policies/evaluate` | `evaluatePromptGuard` | 200 | write |

### 2.3 Planned Schemas

- `PromptGuardPolicy`（`scan_target`, `action`, `rules[]`, `is_default`）
- `PromptGuardRule`（builtin/custom, `pattern`, `severity`）
- `PromptGuardPolicyUpdateRequest`
- `PromptGuardEvaluateRequest` / `PromptGuardEvaluateResponse`

### 2.4 处置语义

| action | 运行时 | Audit |
|---|---|---|
| block | 中断请求/回合 | policy_hit（critical） |
| warn | 放行 + 标记 | policy_hit（warning） |
| log_only | 放行 | policy_hit（info） |

## 3. Page Scope

- 默认策略编辑 + 规则矩阵 + 演练（3b）+ 推理绑定（3b）
- **Non-Goals**：RPM PUT；BOSS 全局模板写 API

## 4. 创建前置条件与操作矩阵

见主维护源。

## 5. Handler 验收（合入 YAML 后）

```bash
curl -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/agent/guard-policies"
curl -X PUT .../guard-policies/$ID -d '{"idempotency_key":"...","enabled":true,"action":"block",...}'
curl -X POST .../guard-policies/evaluate -d '{"text":"...","scan_target":"input"}'
```

## 6. 主维护源

- `docs/console-modules/ai-native/prompt-injection-guard.md`

OpenAPI **未合入** ≠ handler 已实现。
