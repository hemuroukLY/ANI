# Sprint 13 — S05/S06/S07 B 轨执行提示词（codex goal，达到 S01–S04 同级 production-shaped passed）

> 工件归属：Sprint 13 执行驱动配套。地图见 [`sprint13-real-provider-readiness-plan.md`](sprint13-real-provider-readiness-plan.md)，队列见 [`sprint13-loop-execution-prompts.md`](sprint13-loop-execution-prompts.md)。
> 用法：把 **§2 编排提示** 整段粘进 codex goal，逐切片 S05→S06→S07 把 B 轨做到 `production_shape.status=passed`，每切片提交分支。
> 验收标准与 S01–S04 一致：每切片 evidence 必须 `production_shape.status=passed` 且通过 `make validate-sprint13-b-track-production-shape`。

---

## 0. 已核对的现状（2026-06-20，禁止臆测）

- S01–S04 B 轨已 passed（四份 `production_shape.status=passed` evidence，`validate-sprint13-b-track-production-shape` 绿）；Auth/Dex production gate 已通过。
- **S05 MinIO real adapter 未写**、**S06 Milvus real adapter 未写**；**S07 Prometheus real adapter 已存在**（`pkg/adapters/runtime/prometheus_instance_observability.go`）。
- provider 切换开关已存在：`OBJECT_STORE_PROVIDER` / `VECTOR_STORE_PROVIDER` / `INSTANCE_OBSERVABILITY_PROVIDER`（real 分支需补实现/接线）。
- 架构 allowlist 已登记 `github.com/minio/`、`github.com/milvus-io/`（`scripts/validate_component_imports.py`）；Prometheus 走 HTTP query，无需新 SDK。
- production-shape 标准已定义：`deploy/real-k8s-lab/sprint13-production-shaped-gateway-profile.yaml` 已含 S05/S06/S07 `slice_proof_items`。
- 三个 live-gate 校验器 `validate-object-store-live-gate` / `validate-vector-store-live-gate` / `validate-instance-observability-live-gate` 当前为 contract（LIVE PENDING），**需补 `--live` + `--production-shaped` + evidence 输出**（对照已 passed 的 `validate_gpu_inventory_live_gate.py` 实现样板）。

## 1. 硬约束 & 安全协议（CLAUDE.md / ANI-06 / Sprint 11）

- 分支 `feature/sprint12-core-support`，不 push、不动 main、不改 `ports.*` 接口签名、不改 Gateway handler、不动 `/api/v1/svc`。
- 能力经 `pkg/ports` + `pkg/adapters`；组件 SDK 只在 adapter 内，import 已在 allowlist。
- 真实服务器安全：任何写操作前**只读盘点 + 列出预期影响**；组件部署到独立 namespace 并在 gate 用 `--cleanup` 清理临时资源；**凭据/IP 绝不写进可提交文件、evidence、日志**；遇到破坏性/不可逆操作（磁盘/系统盘/fstab/默认 StorageClass）**停下来等人工确认**。
- 契约差异先改 `api/openapi/v1.yaml`（只增可选字段）再实现，保持 v1 兼容。
- **代码-文档一致是硬门禁**：每个 slice 的 live-gate 校验器要求 4 份文档（ANI-DOCS-INDEX.md、ANI-06-开发计划.md、repo/CURRENT-SPRINT.md、repo/development-records/README.md）都包含该 slice 的全部 `REQUIRED_DOC_TOKENS`（见 §3）；任一缺失即门禁失败。

## 2. 编排提示（粘进 codex goal 持续运行）

```text
角色：ANI Core 平台工程师（生产级，严谨可落地）。你在 codex goal 推进 Sprint 13 S05/S06/S07 B 轨，目标与 S01–S04 同级：production_shape.status=passed。

每切片按序加载：CLAUDE.md → ANI-DOCS-INDEX.md → repo/CURRENT-SPRINT.md →
ANI-06-开发计划.md（§0 + §真实底座组件引入强制门禁）→
repo/development-records/sprint13-real-provider-readiness-plan.md →
repo/development-records/sprint13-s05-s07-b-track-prompts.md（§0 现状 + §1 约束 + §3 切片事实）→
对应 per-slice readiness + a-track 记录 → repo/api/openapi/v1.yaml 对应段。

分支：feature/sprint12-core-support。不 push、不动 main、不改 port 签名/handler、不动 /api/v1/svc。

逐切片 S05→S06→S07，每个做完整 B 轨：
1. 只读盘点三台物理服务器（凭据见本机 local-secrets/，绝不入库/回显），确认组件未部署与真实 API（MinIO S3/pre-signed、Milvus collection schema、Prometheus query）。按 §153 更新该 slice readiness。
2. 写/补 real adapter（S05 新建 MinIO ObjectStore adapter；S06 新建 Milvus VectorStore adapter；S07 复用并校验现有 prometheus_instance_observability.go），实现对应 ports 接口；组件 SDK 只在 adapter（allowlist 已含 minio/milvus）。接线 *_PROVIDER 开关 real 分支 + bootstrap/Gateway in-cluster provider 装配（对照 S04 GPU_INVENTORY_PROVIDER=kubernetes_rest / S01 NETWORK_PROVIDER=kubeovn_rest）。
3. 给该 slice live-gate 校验器补 --live + --production-shaped + evidence 输出与 proof_items（样板：scripts/validate_gpu_inventory_live_gate.py）；proof_items 对齐 deploy/real-k8s-lab/sprint13-production-shaped-gateway-profile.yaml 中该 slice 条目。
4. 写 adapter 单测（fake/mock，不依赖真实后端）。
5. 在独立 namespace 部署组件（idempotent），经 production-shaped Gateway 跑真实业务证据：S05 bucket create/list + upload/download pre-signed；S06 vector store create + documents insert + search readiness；S07 instance logs/events/metrics/security-events/exec。用 --cleanup 清理临时资源，输出 production_shape.status=passed 的非敏感 evidence JSON 到 development-records/live-evidence/。
6. 更新关联文档（代码-文档一致，硬门禁）：把该 slice 的全部 REQUIRED_DOC_TOKENS（见加载的 §3）写进 ANI-DOCS-INDEX.md、ANI-06-开发计划.md、repo/CURRENT-SPRINT.md、repo/development-records/README.md；新增/更新 sprint13-<slice>-live-result.md、a-track 记录、loop 驱动 §2 队列 status、readiness-plan 矩阵。
7. 全套门禁全绿并贴输出：
   cd repo && make test && make validate-<该slice domain> && make validate-<该slice live-gate> && make validate-sprint13-b-track-production-shape && python scripts/validate_yaml.py api/openapi/v1.yaml && make validate-doc-entrypoints && git diff --check
8. git commit 到 feature/sprint12-core-support（不 push）。进入下一个 slice。

绝对禁止：改 main/push、改 port 签名或 handler、动 /api/v1/svc、把凭据或 IP 写进任何可提交文件/evidence/日志、对磁盘/系统盘/fstab/默认 StorageClass 等破坏性不可逆操作自作主张（必须停下来等人工确认）。

停止条件（命中即停并报告）：某 slice 需破坏性/不可逆写操作；该 slice 真实组件无法在 lab 就绪；任一门禁失败且无法在不越界前提下修复；S05–S07 全部 passed。

每切片结束输出：slice / 写了哪些 adapter+gate+manifest / 门禁结果 / evidence 路径 / 文档 token 已补 / 下一步。
```

## 3. 每切片钉死事实（防幻觉）

| Slice | 真实 adapter（写/复用） | provider 开关 | ports（不改签名） | live-gate 校验器 | REQUIRED live_checks | REQUIRED_DOC_TOKENS（4 份文档都要有） |
|---|---|---|---|---|---|---|
| **S05 MinIO** | 新建 MinIO ObjectStore adapter（allowlist `github.com/minio/`） | `OBJECT_STORE_PROVIDER=minio` | `ports.ObjectStore` + `ports.StorageService` | `validate-object-store-live-gate` | minio-health-ready, core-bucket-create, core-buckets-list, core-object-upload-presign, core-object-download-presign | `SPRINT13-OBJECTSTORE-MINIO-A-TRACK`, `validate-object-store-live-gate`, `MinIO`, `pre-signed URL`, `LIVE PENDING` |
| **S06 Milvus** | 新建 Milvus VectorStore adapter（allowlist `github.com/milvus-io/`） | `VECTOR_STORE_PROVIDER=milvus` | `ports.VectorStore` + `ports.VectorStoreService` | `validate-vector-store-live-gate` | milvus-health-ready, core-vector-store-create, core-vector-documents-insert, core-vector-search-readiness | `SPRINT13-VECTOR-MILVUS-A-TRACK`, `validate-vector-store-live-gate`, `Milvus`, `LIVE PENDING` |
| **S07 Prometheus** | 复用并校验 `prometheus_instance_observability.go`（HTTP query，无新 SDK） | `INSTANCE_OBSERVABILITY_PROVIDER=prometheus` | `ports.InstanceObservability` | `validate-instance-observability-live-gate` | prometheus-health-ready, core-instance-logs-list, core-instance-events-list, core-instance-metrics-get, core-instance-security-events-list, core-instance-exec-session-create | `SPRINT13-INSTANCE-OBSERVABILITY-PROMETHEUS-A-TRACK`, `validate-instance-observability-live-gate`, `Prometheus`, `kubelet`, `LIVE PENDING` |

通用：每切片对应 domain 校验 = S05/S07 `validate-storage-alpha` 或 `validate-demo-instances`、S06 `validate-vector-alpha`；handler/port 签名不改；deploy manifest 放 `deploy/real-k8s-lab/`；evidence JSON 放 `development-records/live-evidence/sprint13-<slice>-live-evidence.json`，结构对齐 S01–S04。

## 4. 代码-文档一致检查清单（每切片收尾必须做）

1. 把该 slice §3 的 `REQUIRED_DOC_TOKENS` 全部写进 4 份文档（ANI-DOCS-INDEX.md / ANI-06 / CURRENT-SPRINT.md / development-records/README.md）。保留 `LIVE PENDING` 字样（可在边界/历史语境）。
2. 跑该 slice 的 `validate-<slice>-live-gate` 与 `validate-sprint13-b-track-production-shape`、`validate-doc-entrypoints` 确认 token 与 evidence 一致（**本 sprint 已因 ANI-DOCS-INDEX 漏 token 触发过 gate 失败，必须验证**）。
3. evidence `production_shape.status=passed`，且 readiness/a-track/live-result/loop 队列/readiness-plan 矩阵状态同步更新为 passed。

## 5. 关联文档

- 执行地图：[`sprint13-real-provider-readiness-plan.md`](sprint13-real-provider-readiness-plan.md)
- 队列与两轨道：[`sprint13-loop-execution-prompts.md`](sprint13-loop-execution-prompts.md)
- S01 样板就绪声明：[`sprint13-netroute-kubeovn-readiness.md`](sprint13-netroute-kubeovn-readiness.md)
- production-shape 标准：`deploy/real-k8s-lab/sprint13-production-shaped-gateway-profile.yaml`、`validate-sprint13-b-track-production-shape`
- 强制规则 / 真实底座门禁：[`../../CLAUDE.md`](../../CLAUDE.md)、[`../../ANI-06-开发计划.md`](../../ANI-06-开发计划.md)
