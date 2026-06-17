# SPEC: Console 对象存储上传下载

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/buckets` | `listStorageBuckets` | `200 + StorageBucketListResponse` | `401`,`403` |
| POST | `/api/v1/buckets` | `createStorageBucket` | `201 + StorageBucket` | `400`,`401`,`403`,`409` |
| POST | `/api/v1/objects/upload` | `uploadStorageObject` | `202 + AsyncTask` | `400`,`401`,`403`,`422` |
| GET | `/api/v1/objects/{object_id}/download` | `downloadStorageObject` | `200 + ObjectDownloadResponse` | `401`,`403`,`404` |

## 3. TODO-YAML

- 桶策略 API

## 4. References

- `docs/console-modules/compute/storage/object-storage-upload.md`
