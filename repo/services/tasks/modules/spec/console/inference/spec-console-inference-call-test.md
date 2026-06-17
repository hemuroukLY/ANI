# SPEC: Console 推理服务调用测试

> Source: `tasks/modules/prd/console/inference/prd-console-inference-call-test.md`  
> Revised: 2026-06-17

## 1. Summary

Services 子路径调试页；`POST /api/v1/svc/inference-services/{service_id}/test`。

## 2. Frozen Facts

### 2.1 Authority Source

- `ANI-main/repo/api/openapi/services/v1.yaml`

### 2.2 Verified Paths

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| POST | `/api/v1/svc/inference-services/{service_id}/test` | `testInferenceService` | `200 + InferenceServiceTestResponse` | `400`,`401`,`403`,`404`,`422` |

### 2.3 Request

- required: `idempotency_key`, `prompt`
- optional: `max_tokens`

## 3. Page Scope

- 服务上下文只读 + 测试表单 + 结果区

## 4. Non-Goals

- Gateway OpenAI 路径

## 5. References

- 详文：`docs/console-modules/inference/inference-call-test.md`
