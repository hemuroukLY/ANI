# ANI Development Records — 批次归档索引

> 本文件是所有已完成开发批次的**唯一归档索引**。
> 进度追踪三层结构：
> - **全局状态快照** → `ANI-06-开发计划.md` Section 零（30秒定位）
> - **当前冲刺任务** → `repo/CURRENT-SPRINT.md`（每冲刺更新）
> - **已完成批次详情** → 本文件（每批次完成后追加）

> 当前执行处于 **Sprint 5 收敛中**：本地真实代码已完成 M1-K8S-A/B 的 CRUD+kubeconfig+proxy local profile、M1-ENCRYPT-A/B 的 keys+seal/unseal-token+rotate+revoke local profile、M1-SECRETS-A 的 Secret CRUD+bindings local profile、M1-RECONCILE-A 的 background controller adapter/capability 与默认关闭的 bootstrap opt-in 运行剖面、REAL-K8S-LAB-A 真实底座 contract gate；完整 vCluster 生命周期、真实 proxy 转发、controller leader election/指标/退避、真实 KMS/SM4 provider 和真实 K8s Secret 注入尚未完成。本文只做已完成批次归档，不作为当前任务清单使用。
> 2026-05-20 提交前闭环审查：Sprint 2 代码实现、OpenAPI 契约、冻结矩阵、校验脚本和批次记录已对齐；Sprint 3 当前优先项已切换为 `CORE-DEV-PROFILE-A`（原 `MOCK-DEV-A`，已收窄为 Core dev/local profile，不包含 Services 业务 mock）。
> 2026-05-21 Sprint 3 闭环门禁已通过，当前执行切换到 **Sprint 4**；`SPEC-SPLIT-A` 已完成，`SPEC-CORE-BETA` 已完成 Beta 准备矩阵、Core API v1 兼容性基线、SDK/Mock/API 文档加固、四语言 SDK-Mock 联动烟测和提交前审查。当前状态：开发与验收完成，待提交 GitHub；提交完成后再切换下一 Sprint。

---

## 已完成批次（按完成时间排列）

### Sprint 5 Delivery（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| M1-K8S-A | K8s 集群 create/get/list/delete + kubeconfig API + local dev profile + idempotency + tenant isolation；不含 proxy/真实 vCluster provider | m1-k8s-a-core-api-dev-profile.md |
| M1-K8S-B | K8s 集群 proxy Core API 契约 + local dev profile；method/path/query/body 请求边界、幂等 key、路径 allowlist 和 SDK/docs 生成；不含真实 vCluster API 转发 | m1-k8s-b-api-proxy-dev-profile.md |
| REAL-K8S-LAB-A | 真实底座验证门禁：定义三台云 VM K8s/Kube-OVN/KubeVirt/vCluster lab profile、`make validate-real-k8s-profile` 和 live kubectl 检查入口；不代表真实环境已经部署完成 | real-k8s-lab-a-validation-gate.md |
| M1-ENCRYPT-A | Encryption keys create/get/list/delete + seal + unseal-token API + local dev profile + idempotency + tenant isolation；不含真实 KMS/SM4 provider | m1-encrypt-a-core-api-dev-profile.md |
| M1-ENCRYPT-B | Encryption key rotate/revoke API + local dev profile + idempotency + state guard；不含真实 KMS/SM4 provider 生命周期操作 | m1-encrypt-b-key-rotation-revoke-local-profile.md |
| M1-SECRETS-A | Secret create/get/list/delete + bindings API + local dev profile + idempotency + tenant isolation；响应不返回明文，不含真实 K8s Secret 注入 | m1-secrets-a-core-api-dev-profile.md |
| M1-RECONCILE-A | WorkloadReconcileController adapter + bootstrap capability + opt-in 后台运行；扫描 reconcile target、观察 provider 状态、回写 instance 状态；不含 leader election/指标/退避 | m1-reconcile-a-background-controller.md |

### Sprint 5 Kickoff（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| SPRINT5-KICKOFF-A | Sprint 5 启动：执行入口切换与三份主文档状态对齐 | sprint5-kickoff-a.md |

### Sprint 4 API Beta Preparation（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| SPEC-SPLIT-A | Core/Services API 分层收口：Services 业务路径迁移到 Services API，Gateway Services stub 改挂 `/api/v1/svc`，SDK metadata 自然分层 | spec-split-a-core-services-api-boundary.md |
| SPEC-CORE-BETA-A | Core API Beta 准备矩阵：P0 path/schema、分页、幂等、状态机、dev_profile、RBAC scope 和 Core/Services 边界守卫 | spec-core-beta-a-readiness-matrix.md |
| SPEC-COMPAT-A | Core API v1 兼容性基线：保护 path/method/operationId/参数/响应/schema 字段，允许新增可选能力但阻止破坏性变更 | spec-compat-a-core-api-v1-baseline.md |
| SDK-BETA-A | 四语言 SDK 幂等 helper：生成 idempotency key、注入请求体、metadata 标出 Core 幂等操作 | sdk-beta-a-idempotency-helper.md |
| SDK-BETA-B | 四语言 SDK cursor 分页 helper：构造 limit/cursor 参数、metadata 标出 Core 分页操作 | sdk-beta-b-cursor-pagination-helper.md |
| SDK-BETA-C | 四语言 SDK 统一 API error helper：错误对象、错误码清单、错误码判断 | sdk-beta-c-api-error-helper.md |
| SDK-BETA-D | 四语言 SDK basic example：client 初始化、幂等、cursor 分页和 API error helper 组合用法 | sdk-beta-d-basic-examples.md |
| SDK-MOCK-SMOKE-A | Core Python SDK 调用 Mock Server 烟测：标准库 HTTP request 能力、分页响应和标准错误响应校验 | sdk-mock-smoke-a-python-sdk-mock-server.md |
| SDK-MOCK-SMOKE-B | Core TypeScript SDK 调用 Mock Server 烟测：fetch request 能力、分页响应和标准错误响应校验 | sdk-mock-smoke-b-typescript-sdk-mock-server.md |
| SDK-MOCK-SMOKE-C | Core Go SDK 调用 Mock Server 烟测：net/http Request 能力、分页响应和标准错误响应校验 | sdk-mock-smoke-c-go-sdk-mock-server.md |
| SDK-MOCK-SMOKE-D | Core Java SDK 调用 Mock Server 烟测：HttpClient request 能力、分页响应和标准错误响应校验 | sdk-mock-smoke-d-java-sdk-mock-server.md |
| MOCK-A | Core Mock Server：由 `api/openapi/v1.yaml` 驱动，覆盖 Core API 成功响应和统一错误结构 | mock-a-core-openapi-mock-server.md |
| DOC-API-A | 静态 API 文档生成：Core/Services API 契约生成 docs/api，并校验 operation/schema 覆盖 | doc-api-a-static-api-docs.md |
| SPRINT4-CLOSURE-A | Sprint 4 关联性闭环门禁：统一校验 API/SDK/Mock/Docs/Records 与 Makefile 入口 | sprint4-closure-a-contract.md |

### Sprint 3 Network / Storage / SDK（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| M1-NETWORK-A | VPC/Subnet/SecurityGroup/LoadBalancer Core API 契约、Gateway dev profile、持久化边界和网络合同守卫 | m1-network-a-core-api-dev-profile.md |
| M1-NETWORK-A | KubeOVN/Kubernetes provider 渲染边界：Vpc/Subnet、NetworkPolicy、Service 清单与 bootstrap capability | m1-network-a-kubeovn-renderer.md |
| M1-NETWORK-A | 网络 provider server-side dry-run、默认关闭 apply gate、KubeOVN/Kubernetes REST path 映射 | m1-network-a-provider-dry-run-apply-gate.md |
| M1-NETWORK-A | 网络 provider 状态读取边界：KubeOVN/Kubernetes 资源状态归一化为 ANI 网络状态与失败原因 | m1-network-a-provider-status-reader.md |
| M1-NETWORK-A | 网络状态 reconcile：provider observation 校验后回写网络资源 state/reason/updated_at | m1-network-a-status-reconcile.md |
| M1-STORAGE-A | volumes/filesystems/objects Core API 契约、Gateway dev profile、租户隔离和存储合同守卫 | m1-storage-a-core-api-dev-profile.md |
| M1-STORAGE-A | storage metadata 持久化边界、RLS 迁移、bootstrap capability 和持久化单元测试 | m1-storage-a-persistence-boundary.md |
| M1-STORAGE-A | 存储 provider 渲染边界：PVC manifest、objectstore metadata intent 和 bootstrap capability | m1-storage-a-provider-renderer.md |
| M1-STORAGE-A | 存储 provider server-side dry-run、默认关闭 apply gate、objectstore 执行边界保留 | m1-storage-a-provider-dry-run-apply-gate.md |
| M1-STORAGE-A | 存储 provider 状态读取和 metadata state/reason 回写闭环 | m1-storage-a-status-reconcile.md |
| M1-VSTORE-A | vector-stores Core API 契约、Gateway dev profile、搜索响应结构和合同守卫 | m1-vstore-a-core-api-dev-profile.md |
| SDK-ALPHA-A | Core/Services 四语言 SDK Alpha 生成、分层隔离和 smoke 门禁 | sdk-alpha-a-generation-smoke.md |
| M1-WKID-A | Workload Identity P0：实例 lifecycle-bound scoped API key、Secret 引用注入和删除 revoke | m1-wkid-a-workload-identity-p0.md |
| CORE-DEV-PROFILE-A | Core P0 API dev/local profile 显式标记、Core/Services mock 边界和合同守卫 | core-dev-profile-a-boundary-contract.md |
| SPRINT3-CLOSURE-A | Sprint 3 闭环审查门禁：批次记录、API/SDK 分层和各批次合同守卫统一校验 | sprint3-closure-a-contract.md |

### Sprint 2 Core API Alpha（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| SPEC-CORE-ALPHA-A | `/api/v1/instances` Core Alpha path/schema/RBAC scope + Gateway 主路径 + 合同守卫 | spec-core-alpha-a-instance-contract-guard.md |
| SPEC-CORE-ALPHA-B | Core API Alpha 机器可读冻结矩阵，校验 path/schema/error/state/RBAC scope 与 Gateway/runtime 对齐 | spec-core-alpha-b-freeze-matrix.md |
| M1-INSTANCE-U-A | VM `termination_protection` 危险操作 precheck、failed operation timeline 和 lifecycle policy 持久化 | m1-instance-u-a-termination-protection.md |
| M1-INSTANCE-U-B | VM SSH 连接元数据 schema、Gateway dev profile 响应和 `ssh_connection` 持久化 | m1-instance-u-b-vm-ssh-info.md |
| M1-INSTANCE-U-C | VM console/VNC/serial session 返回 `operation_id/url/expires_at` 并写入 operation timeline | m1-instance-u-c-console-session-timeline.md |
| M1-INSTANCE-U-D | VM `snapshot` local profile、`snapshots[]` 响应、operation timeline 和 JSONB 持久化 | m1-instance-u-d-vm-snapshot-local-profile.md |
| M1-INSTANCE-U-E | VM `attach_volume/detach_volume` local profile、`volumes[]` 响应和 operation timeline | m1-instance-u-e-vm-volume-binding-local-profile.md |
| M1-INSTANCE-V-A | Container/GPU Container `replicas/revision/rollout_status/history` 响应和 `container_status` 持久化 | m1-instance-v-a-container-rollout-status.md |
| M1-INSTANCE-V-B | Container/GPU Container `rollback` local profile、revision 回退和 `rollback_revision` operation timeline | m1-instance-v-b-container-rollback-local-profile.md |
| M1-INSTANCE-V-C | GPU Container `vendor/model/count/scheduling_reason/utilization_percent` 响应和 `gpu_status` 持久化 | m1-instance-v-c-gpu-status-local-profile.md |

### V8 架构重规划（2026-05-14~15）

| 批次 | 内容摘要 |
|---|---|
| V8-ARCH | Core/Services 分层、ANI-02/06 重写、CLAUDE.md 强制约定 |
| AWS-HARDENING | /healthz /readyz、idempotency_key port、ReconcileController port、operations DB 表、permissions schema |

### Sprint 1 Foundation（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| M1-HEALTH-A | Gateway/Auth/Model/Task 标准 /healthz 与 /readyz 探针 | m1-health-a-health-endpoints.md |
| M1-IDEM-A | 实例 create/lifecycle 幂等锁、DB 原子冲突回放和 bootstrap 接线 | m1-idem-a-idempotency-wire-up.md |

### M1 基础设施底座（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| M1-INFRA-A | ani-system 命名空间、NetworkPolicy、ServiceAccount 基线 | m1-infra-a-baseline.md |
| M1-INFRA-B | PostgreSQL/NATS/Redis/MinIO/Milvus/Harbor 组件安装 profile | m1-infra-b-component-profiles.md |
| M1-INFRA-C | KubeOVN VPC/Subnet 模板、沙箱出口限制 | m1-infra-c-network-isolation.md |
| M1-INFRA-D | cluster preflight validation profile | m1-infra-d-cluster-preflight.md |
| M1-INFRA-E | GPU scheduling baseline（Volcano/HAMi/DCGM）| m1-infra-e-gpu-scheduling-baseline.md |
| M1-INFRA-F | GPU preflight/e2e hardening | m1-infra-f-gpu-preflight-e2e.md |
| M1-GPU-A | 异构 GPU 发现调度契约（NVIDIA/昇腾/海光/GPUInventory port）| m1-gpu-a-heterogeneous-gpu-contract.md |
| M1-RUNTIME-A | WorkloadRuntime port（全实例类型抽象）| m1-runtime-a-workload-runtime.md |

### M1 Instance Fabric（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| M1-INSTANCE-A | 核心实例对象、生命周期、网络平面、存储附件契约 | m1-instance-a-instance-fabric.md |
| M1-INSTANCE-B | PlanningRuntime 实例规划器 | m1-instance-b-planning-runtime.md |
| M1-INSTANCE-C | K8s/KubeVirt provider dry-run renderer | m1-instance-c-provider-renderer.md |
| M1-INSTANCE-D | 本地 admission guardrail | m1-instance-d-admission-guardrail.md |
| M1-INSTANCE-E | 实例计划/渲染/准入审计持久化 | m1-instance-e-plan-audit.md |
| M1-INSTANCE-F | WorkloadProviderDryRun executor boundary | m1-instance-f-provider-dry-run.md |
| M1-INSTANCE-G | WorkloadProviderApply 执行门控 | m1-instance-g-provider-apply-gate.md |
| M1-INSTANCE-H | WorkloadStatusReconciler 状态回写 | m1-instance-h-status-reconcile.md |
| M1-INSTANCE-I | WorkloadProviderStatusReader + Orchestrator | m1-instance-i-orchestrator.md |
| M1-INSTANCE-J | WorkloadInstanceStore + workload_instances RLS 表 | m1-instance-j-instance-store.md |
| M1-INSTANCE-K | KubernetesProviderAdapter + Client | m1-instance-k-provider-adapter.md |
| M1-INSTANCE-L | WorkloadInstanceService API 层 | m1-instance-l-instance-service.md |
| M1-INSTANCE-M | 生命周期 + 可视化运维 API | m1-instance-m-lifecycle-ops.md |
| M1-INSTANCE-N | Kubernetes provider 执行剖面 | m1-instance-n-kubernetes-provider-execution.md |
| M1-INSTANCE-O | adapter-owned KubernetesRESTClient | m1-instance-o-kubernetes-rest-client.md |
| M1-INSTANCE-P | bootstrap/config provider wiring | m1-instance-p-kubernetes-bootstrap-wiring.md |
| M1-INSTANCE-Q | KubernetesLifecycleExecutor | m1-instance-q-kubernetes-lifecycle-execution.md |
| M1-INSTANCE-R | KubernetesInstanceOps | m1-instance-r-kubernetes-ops-execution.md |
| M1-INSTANCE-S | VM console/VNC/serial remote ops session 边界 | — |
| M1-INSTANCE-T | 操作语义横切基础：operation_id、timeline、幂等回放和操作查询 | m1-instance-t-operation-semantics.md |
| M1-E2E-A | M1 端到端集成剖面 | m1-e2e-a-instance-profile.md |
| M1-E2E-B | M1 real provider integration regression profile | m1-e2e-b-real-provider-profile.md |

### ARCH-ADAPTER 系列（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| ARCH-ADAPTER-A / M1-ARCH-A | 开源组件松耦合适配器架构设计 | m1-arch-a-component-adapter-design.md |
| ARCH-ADAPTER-B | pkg/ports + pkg/adapters + bootstrap.Capabilities 骨架 | arch-adapter-b-ports-adapters-skeleton.md |
| ARCH-ADAPTER-GUARD-A | 组件 SDK 直接导入扫描与 allowlist 护栏 | arch-adapter-guard-a-component-imports.md |
| ARCH-ADAPTER-C | 第一批迁移（CacheStore + MessageBus）| arch-adapter-c-first-migration.md |
| ARCH-ADAPTER-C-2 | pgx/metadata 依赖 bounded_direct 分类 | arch-adapter-c-2-metadata-boundaries.md |

### M2 Gateway / Auth（2026-05）

| 批次 | 内容摘要 | 文件 |
|---|---|---|
| M2.1-TASK-A/B | task-service + transactional outbox | m2-1-task-a-b-task-service-outbox.md |
| M2.1-TASK-C | worker mutation RPCs | m2-1-task-c-worker-mutations.md |
| M2.2-AUTH-A~K | auth-service 完整实现（JWT/OIDC/JWKS/RBAC/API Key）| m2-2-auth-*.md |
| M2.2-AUTH-FINAL | Auth 生产收尾：OIDC/Dex 护栏、Gateway Auth REST、API Key 管理、合同守卫与 Docker Dex smoke | m2-2-auth-final-production-closeout.md |

---

## 批次完工的更新流程

> 完整规约在 `CLAUDE.md` → "📋 开发进度更新规约"，以下是速查版本。

**批次完成时（必须按顺序）：**

```
① make test                              → 全通（零失败）
② 新建 {批次名}.md（用 TEMPLATE.md）    → 填入完成日期/验证结果/关键文件
③ 本文件 README.md                       → 在对应分组表格追加一行
④ repo/CURRENT-SPRINT.md                 → 该批次 🔄→✅，下一批次 ⏳→🔄
⑤ ANI-06-开发计划.md Section 零         → 更新批次/Sprint 状态行
⑥ git commit -m "feat: {批次名} {一句话}"
```

**Sprint 全部完成时，额外：**
```
⑦ ANI-06 Section 零 Sprint 行：🔄→✅（填完成日期）/ 下一Sprint：⏳→🔄
⑧ repo/CURRENT-SPRINT.md 整体重写为下一 Sprint 内容
⑨ git commit -m "sprint: Sprint N completed"
```
