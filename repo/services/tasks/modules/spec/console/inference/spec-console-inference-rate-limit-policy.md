# SPEC: Console 推理限流策略

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| PUT | `/api/v1/svc/inference-services/{service_id}/policies` | `updateInferenceServicePolicies` | `200 + InferenceService` | `400`,`401`,`403`,`404`,`422` |

InferenceServicePolicyUpdateRequest: idempotency_key required

## 3. References

- `docs/console-modules/inference/inference-rate-limit-policy.md`
