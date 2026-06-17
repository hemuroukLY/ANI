# SPEC: Console 块存储卷快照

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/volumes/{volume_id}/snapshots` | `listVolumeSnapshots` | `200 + VolumeSnapshotListResponse` | `401`,`403`,`404` |
| POST | `/api/v1/volumes/{volume_id}/snapshots` | `createVolumeSnapshot` | `202 + VolumeSnapshot` | `400`,`401`,`403`,`404`,`409` |

## 3. TODO-YAML

- DELETE snapshot

## 4. References

- `docs/console-modules/compute/storage/block-storage-snapshot.md`
