# SPEC: Console 批任务实例

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/instances?kind=batch_job` | `listInstances` | `200 + InstanceListResponse` | `401`,`403` — **kind 枚举未含 batch_job** |
| POST | `/api/v1/instances` | `createInstance` | `201 + CreateInstanceResponse` | `400`,`401`,`403`,`409`,`422` |
| GET | `/api/v1/instances/{instance_id}` | `getInstance` | `200 + InstanceRecord` | `401`,`403`,`404` |
| POST | `/api/v1/instances/{instance_id}/lifecycle` | `applyInstanceLifecycle` | `200 + InstanceLifecycleResponse` | `400`,`401`,`403`,`404`,`409`,`422` |

## 3. TODO-YAML

- listInstances kind 枚举扩展 batch_job

## 4. References

- `docs/console-modules/compute/batch-job-instances.md`
