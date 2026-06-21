# Sprint 13 切片 05 — object-store bucket/upload/download MinIO real provider 就绪声明

> 记录类型：Per-slice readiness（ANI-06「真实底座组件引入强制门禁」§153 的执行前声明）
> 工件归属：Sprint 13 / Core real provider 与 live gate 收敛
> 执行地图：[`sprint13-real-provider-readiness-plan.md`](sprint13-real-provider-readiness-plan.md)
> 状态：**production-shaped gate passed**（B 轨已完成；evidence：`live-evidence/sprint13-objectstore-minio-live-evidence.json`）。该结论只覆盖 S05 object-store provider/live gate，不代表 full platform production ready。

---

## 0. 已核对的真实事实（禁止臆测）

1. Sprint 12 已落地 storage bucket、object upload/download 的 Core API contract：`listStorageBuckets`、`createStorageBucket`、`uploadStorageObject`、`downloadStorageObject` 均在 `storage_resources.go` 经 `ports.StorageService` 暴露。
2. `ports.ObjectStore` 已存在 `EnsureBucket`、`SignedUploadURL`、`SignedDownloadURL` 边界；`LocalStorageService` 在显式注入 object store 时会复用这些 port 方法，upload/download 返回 pre-signed URL，不走 multipart。
3. OpenAPI 已定义 `StorageBucketRecord` / `StorageBucketListResponse` / `StorageObjectUploadRequest` / `StorageObjectUploadResponse` / `StorageObjectDownloadInfo`，POST 请求保留 `idempotency_key`，RBAC scope 保持 `scope:objects:*`。
4. S05 A 轨只允许新增 MinIO/S3 兼容 adapter、fake/mock 单测、契约级 live-gate 和文档闭环；B 轨已在人工授权后部署临时 MinIO 验证组件并执行真实上传/下载 live gate。

## 1. §153 五项声明

| 项 | 内容 |
|---|---|
| **当前状态** | `OBJECT_STORE_PROVIDER=minio` 下 production-shaped Gateway 已接入真实 MinIO/S3-compatible backend；bucket create/list、upload/download pre-signed URL、实际 PUT/GET 与 cleanup 已通过。 |
| **真实组件 + 版本** | MinIO / S3 compatible object store；本次为 Sprint 13 S05 live-gate 临时验证部署，独立 namespace + 临时数据目录，不声明长期生产 MinIO 存储方案。 |
| **live gate 命令** | 本地契约：`make validate-storage-alpha validate-object-store-live-gate`；真实 B 轨：`python scripts/validate_object_store_live_gate.py --live --production-shaped --cleanup --evidence-output development-records/live-evidence/sprint13-objectstore-minio-live-evidence.json` 已通过。 |
| **evidence 输出路径** | result：`repo/development-records/sprint13-objectstore-minio-live-result.md`；evidence：`repo/development-records/live-evidence/sprint13-objectstore-minio-live-evidence.json`；未归档 access key、secret key、session token、服务器 IP 或完整 pre-signed URL。 |
| **失败边界（不得声称）** | S05 production-shaped gate passed 不等于 full platform production ready；不得据此声明长期 MinIO HA、Ceph/RGW/PVC 存储规划、备份策略、对象生命周期策略或全平台 release gate 已完成。 |

## 2. 代码边界

- A 轨新增 `pkg/adapters/objectstore/MinIOObjectStore`，实现既有 `ports.ObjectStore`，不改 port 接口签名，不改 Gateway handler，不新增 `/api/v1/svc`。
- adapter 使用标准 HTTP + AWS SigV4 生成 S3-compatible pre-signed URL；bucket 创建走 `HEAD /bucket` + `PUT /bucket`，已有 bucket 幂等成功。
- adapter 支持内部 `OBJECT_STORE_ENDPOINT` 与外部 `OBJECT_STORE_PUBLIC_ENDPOINT` 分离：Gateway 用集群内 Service 访问 MinIO API，pre-signed URL 使用 live gate 可访问的 public endpoint。
- `pkg/bootstrap` 与 Gateway runtime 仅在 `OBJECT_STORE_PROVIDER=minio` 显式配置时构造并注入 MinIO object store；默认 dev/local profile 不变，避免把未配置 adapter 误标为 runtime ready。

## 3. 真实服务器安全

- B 轨只部署 Sprint 13 S05 live-gate 临时 MinIO 验证组件；使用 `emptyDir` 是为了避免触碰 Rook-Ceph、默认 StorageClass、节点裸盘或 fstab。
- 若后续把 MinIO 定位为长期生产组件，必须另行评审 Rook-Ceph PVC/Ceph RGW/分布式 MinIO/专用磁盘方案；不得在未授权时格式化或挂载 ANI3 未使用 HDD。
- 凭据、endpoint、服务器 IP 与 pre-signed URL 不得写入可提交文件、evidence 或回复。

## 4. 完成判定（S05 B 轨）

```bash
cd repo && make test && make validate-storage-alpha && make validate-object-store-live-gate && make validate-sprint13-b-track-production-shape && python scripts/validate_yaml.py api/openapi/v1.yaml && make validate-doc-entrypoints && git diff --check
```

## 5. 关联文档

- Sprint 13 执行地图：[`sprint13-real-provider-readiness-plan.md`](sprint13-real-provider-readiness-plan.md)
- 当前冲刺入口：[`../CURRENT-SPRINT.md`](../CURRENT-SPRINT.md)
- S05 A 轨记录：[`sprint13-objectstore-minio-a-track.md`](sprint13-objectstore-minio-a-track.md)
- S05 B 轨结果：[`sprint13-objectstore-minio-live-result.md`](sprint13-objectstore-minio-live-result.md)
- 代码：`pkg/ports/object_store.go`、`pkg/adapters/objectstore/minio_store.go`、`pkg/bootstrap/deps.go`、`services/ani-gateway/internal/router/storage_resources.go`
