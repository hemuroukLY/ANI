# SPEC: Console 容器实例可观测性

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/instances/{instance_id}/logs` | `listInstanceLogs` | `200 + InstanceLogListResponse` | `401`,`403`,`404` |
| GET | `/api/v1/instances/{instance_id}/events` | `listInstanceEvents` | `200 + InstanceEventListResponse` | `401`,`403`,`404` |
| GET | `/api/v1/instances/{instance_id}/metrics` | `getInstanceMetrics` | `200 + InstanceMetricsResponse` | `401`,`403`,`404` |
| POST | `/api/v1/instances/{instance_id}/exec` | `createInstanceExecSession` | `200 + InstanceExecSession` | `400`,`401`,`403`,`404`,`422` |

## 3. References

- `docs/console-modules/compute/container-observability.md`
