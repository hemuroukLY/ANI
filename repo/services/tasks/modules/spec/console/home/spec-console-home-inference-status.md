# SPEC: Console 首页推理服务状态

> Revised: 2026-06-17

## 1. Summary

Services 聚合视角；`GET /api/v1/svc/inference-services`。

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/svc/inference-services` | `listInferenceServices` | `200` + items | `401` |

状态 enum: pending, deploying, running, stopping, stopped, failed

## 3. TODO-YAML

- 最近活跃服务数口径（metering/access log）

## 4. References

- `docs/console-modules/home/home-inference-status.md`
