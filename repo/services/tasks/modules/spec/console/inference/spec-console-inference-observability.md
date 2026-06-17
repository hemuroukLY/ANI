# SPEC: Console 推理服务可观测性

> Revised: 2026-06-17

## 1. Summary

跨 Services（日志）+ Core（指标）两层观测页。

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/svc/inference-services/{service_id}/logs` | `getInferenceServiceLogs` | `200 + InferenceServiceLogListResponse` | `401`,`403`,`404` |
| GET | `/api/v1/observability/query` | `queryObservability` | `200 + ObservabilityQueryResponse` | `400`,`401`,`403` |

## 3. TODO-YAML

- `GET .../events`

## 4. Page Scope / Non-Goals

见详文。

## 5. References

- `docs/console-modules/inference/inference-observability.md`
