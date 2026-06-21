# SPRINT13-OBJECTSTORE-MINIO-LIVE-A - object-store MinIO live gate result

> 记录类型：Sprint 13 B-track production-shaped live result
> 完成日期：2026-06-20
> 范围：ANI Core S05 object-store bucket/upload/download provider
> 状态：**production-shaped gate passed**；不代表 full platform production ready

## §153 五项实测结果

| 项 | 实测结果 |
|---|---|
| 当前状态 | S05 已执行 `--live --production-shaped --cleanup` live gate 并通过。Gateway 使用 `OBJECT_STORE_PROVIDER=minio` 接入 MinIO / S3-compatible object store，Core `/buckets`、`/objects/upload`、`/objects/{object_id}/download` 与 pre-signed URL 实际 PUT/GET 均通过。 |
| 真实组件 + 版本 | MinIO / S3-compatible object store；本次为 Sprint 13 S05 live-gate 临时验证部署，使用独立 namespace 与临时数据目录，不声明长期生产 MinIO 存储方案。 |
| live gate 命令 | `python scripts/validate_object_store_live_gate.py --live --production-shaped --cleanup --gateway-url <redacted>/api/v1 --ani-bearer-token <redacted> --minio-url <redacted> --evidence-output development-records/live-evidence/sprint13-objectstore-minio-live-evidence.json` |
| evidence 输出路径 | `repo/development-records/live-evidence/sprint13-objectstore-minio-live-evidence.json` |
| 边界 | Production-shaped gate passed 只证明 S05 Gateway runtime、MinIO object-store adapter、bucket create/list、upload/download pre-signed URL 与 cleanup 门禁通过；不代表 production ready / full platform release，不代表长期 MinIO HA、Ceph/RGW/PVC 存储规划、备份策略或对象生命周期策略全部完成。 |

## Evidence 摘要

```json
{
  "minio_health_status": 200,
  "bucket_create_status": 201,
  "bucket_list_status": 200,
  "bucket_list_count": 1,
  "upload_presign_status": 200,
  "actual_upload_status": 200,
  "download_presign_status": 200,
  "actual_download_status": 200,
  "cleanup_api_key_status": 201,
  "cleanup_status": 200,
  "cleanup_api_key_revoke_status": 200,
  "production_shape": {
    "status": "passed",
    "transport_profile": "production_gateway_and_object_store_service",
    "missing_items": [],
    "proof_items": [
      "production_gateway",
      "production_object_store_credentials",
      "production_presigned_url_endpoint"
    ]
  }
}
```

## 代码与部署闭环

- `validate_object_store_live_gate.py --production-shaped` 拒绝 local Gateway / local MinIO endpoint，并写入脱敏 evidence。
- `--cleanup` 使用短期 scoped API key 删除本次 live gate 创建的 object，并撤销该临时 key；evidence 不记录 bearer token、API key、access key、secret key、服务器 IP、MinIO endpoint 或完整 pre-signed URL。
- Gateway production-shaped manifest 使用 Secret 注入 object-store endpoint、public endpoint 与凭据；仓库不提交任何真实凭据。
