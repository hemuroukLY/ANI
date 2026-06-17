# SPEC: Console VM 快照与恢复

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | 成功 | 错误 |
|---|---|---|---|
| POST | `/api/v1/instances/{instance_id}/lifecycle` (snapshot/rollback) | `200 + InstanceLifecycleResponse` | `400`,`401`,`403`,`404`,`409`,`422` |
| GET | `/api/v1/volumes/{volume_id}/snapshots` | `200 + VolumeSnapshotListResponse` | `401`,`403`,`404` |
| POST | `/api/v1/volumes/{volume_id}/snapshots` | `202 + VolumeSnapshot` | `400`,`401`,`403`,`404`,`409` |

Lifecycle required: action, idempotency_key

## 3. References

- `docs/console-modules/compute/vm-snapshot-restore.md`
