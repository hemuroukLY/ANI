# ANI · 当前冲刺上手指南

> 新开发者（人类或 AI 工具）的第一个入口文件。本文只描述当前真实执行状态；历史完成批次查 `repo/development-records/README.md`。

> 历史门禁保留：Sprint 4 的 `SPEC-SPLIT-A`、`SPEC-CORE-BETA`、`SPEC-COMPAT-A`、`SDK-BETA-A`、`SDK-BETA-B`、`SDK-BETA-C`、`SDK-BETA-D`、`SDK-MOCK-SMOKE-A`、`SDK-MOCK-SMOKE-B`、`SDK-MOCK-SMOKE-C`、`SDK-MOCK-SMOKE-D`、`MOCK-A`、`DOC-API-A`、`SPRINT4-CLOSURE-A` 和 `validate-sprint4-closure` 仍是提交前回归门禁，不作为当前任务清单。

## 当前冲刺

| 字段 | 值 |
|---|---|
| **冲刺编号** | Sprint 5（收敛中） |
| **主题** | K8s 集群管理 + 加解密主链路 |
| **当前状态** | 🔄 local profile 主链路与 controller 可配置运行剖面已收敛，真实 provider/HA 生产能力仍未完成 |
| **已由代码证明完成** | `M1-K8S-A` local profile（CRUD + kubeconfig + proxy）、`M1-ENCRYPT-A/B` local profile（keys + seal/unseal-token + rotate + revoke）、`M1-SECRETS-A` local profile（Secret CRUD + bindings）、`M1-RECONCILE-A` controller adapter/capability + bootstrap opt-in 运行、`REAL-K8S-LAB-A` contract gate |
| **不可标记为完成** | vCluster 真实生命周期、真实 vCluster API 转发、controller leader election/指标/重试退避、真实 KMS/SM4 provider、真实 K8s Secret 注入 |
| **关联历史门禁** | Sprint 4 `SPEC-CORE-BETA` / `api/core-beta-readiness.yaml`、`DOC-API-A` 仍需保持 API docs 与 OpenAPI 同步；`SDK-BETA-A`、`SDK-BETA-B`、`SDK-BETA-C`、`SDK-BETA-D` 仍需保持 SDK helper 与新增 Core API 同步 |
| **最后校准日期** | 2026-05-23 |

## 已完成切片

1. `M1-K8S-A` local profile：`/api/v1/k8s-clusters` create/get/list/delete、`GET /api/v1/k8s-clusters/{cluster_id}/kubeconfig`、`POST /api/v1/k8s-clusters/{cluster_id}/proxy` 已有路由、OpenAPI schema、ports 接口、local runtime service、租户隔离、幂等请求、SDK 生成和 router 单元测试。
2. `M1-ENCRYPT-A/B` local profile：`/api/v1/encryption/keys` create/get/list/delete、`POST /api/v1/encryption/seal`、`POST /api/v1/encryption/unseal-token`、`POST /api/v1/encryption/keys/{key_id}/rotate`、`POST /api/v1/encryption/keys/{key_id}/revoke` 已有路由、OpenAPI schema、ports 接口、local runtime service、租户隔离、幂等 create/seal/rotate/revoke、SDK 生成和 router 单元测试。
3. `M1-SECRETS-A` local profile：`/api/v1/secrets` create/get/list/delete 与 `POST /api/v1/secrets/{secret_id}/bindings` 已有路由、OpenAPI schema、ports 接口、local runtime service、租户隔离、幂等创建、SDK 生成和 router 单元测试；响应只返回元数据和 key 名称，不返回明文。
4. `M1-RECONCILE-A` controller：`LocalWorkloadReconcileController` 已能扫描 reconcile target、调用 provider status reader、复用 `WorkloadStatusReconciler` 回写 instance 状态，并在 provider missing 时标记 failed；`bootstrap.RunGRPC` 已支持通过 `WORKLOAD_RECONCILE_CONTROLLER_ENABLED=true` 显式启动后台循环，默认关闭。
5. `REAL-K8S-LAB-A` contract gate：`deploy/real-k8s-lab/profile.yaml`、`scripts/validate_real_k8s_profile.py` 和 `make validate-real-k8s-profile` 已建立；默认校验门禁定义和文档闭环，三台云 VM 就绪后可用 live 模式执行真实 kubectl 检查。

## 当前事实边界

- K8s 集群当前是 local dev profile 模拟服务，不是真实 vCluster provider。
- Kubeconfig 当前指向模拟 vCluster endpoint，不能当作真实 kubectl/Helm 生产凭据。
- K8s proxy 当前是 Core API 契约和 local dev profile 响应，不会真实转发到 vCluster API Server。
- Encryption 当前提供 local profile seal/unseal-token/rotate/revoke，不是真实 KMS/SM4 加解密或真实 KMS 生命周期管理。
- Secret binding 当前只记录绑定意图，不执行真实 K8s Secret 写入、实例环境变量注入或文件挂载。
- Reconcile controller 当前完成 adapter、bootstrap capability 与默认关闭的 opt-in 后台运行剖面；不代表已具备 leader election、指标、重试退避或独立 worker 部署形态。
- REAL-K8S-LAB-A 当前只完成验证门禁定义，不代表三台云 VM 已部署，也不代表 K8s/Kube-OVN/KubeVirt/vCluster 已真实跑通。
- Sprint 6 不能作为当前执行入口，直到 Sprint 5 未完成项被代码、API 契约和测试共同证明完成或明确重新排期。

## 真实底座门禁

从 Sprint 5 起，涉及 K8s、Kube-OVN、KubeVirt、vCluster、KMS/SM4、K8s Secret 注入等真实组件的能力，不能再只靠 local profile 宣称完成。local profile 只能证明 API、SDK、状态机和调用边界正确；真实运行能力必须由真实组件环境、固定验证命令或批次记录证明。当前固定入口是 `REAL-K8S-LAB-A` 和 `make validate-real-k8s-profile`。

当前必须并行准备的真实验证环境：

| 组件 | 进入时机 | 验证目标 |
|---|---|---|
| K8s 测试集群 | Sprint 5 当前起 | API Server、Namespace、RBAC、ServiceAccount、StorageClass 基础可用 |
| Kube-OVN | Sprint 5 当前起 | VPC、Subnet、NetworkPolicy、Service/LB 可创建、可观察 |
| KubeVirt | Sprint 5 当前起 | VM 创建、启动、停止、删除、console/VNC 可运行 |
| vCluster | Sprint 5 当前起 | K8s 集群创建、kubeconfig、proxy 能真实访问租户集群 |
| KMS/SM4 + K8s Secret | Sprint 5~6 | 加解密、密钥轮换、Secret 写入和实例注入真实跑通 |

## 下一步

1. 继续 Sprint 5，不切换 Sprint 6。
2. 并行启动 `REAL-K8S-LAB-A` 真实底座验证线，优先补 K8s、Kube-OVN、KubeVirt、vCluster 的部署与固定验证记录。
3. 下一个代码闭环：K8s 真实 vCluster provider 或将 proxy port 接到真实 vCluster API Server；controller 后续只补 leader election、指标、退避和独立 worker 部署形态。
4. Encryption/Secrets 后续优先级：真实 KMS/SM4 provider、真实 provider 下的密钥生命周期实现、真实 K8s Secret 注入。

## 文档入口边界

- `CLAUDE.md` 只维护稳定强制规则、读取顺序、架构边界、提交门禁和 Karpathy 四条开发原则。
- 当前 Sprint 的详细完成项、未完成项、验收命令、下一步和真实底座边界以本文为准。
- 批次实现细节只写入 `repo/development-records/*.md`，不得把每日开发流水账或 API path 长列表写回 `CLAUDE.md`。
- 修改入口文档后必须运行 `make validate-doc-entrypoints`。

## 验收命令

```bash
make validate-doc-entrypoints
python scripts/validate_yaml.py api/openapi/v1.yaml api/openapi/services/v1.yaml
make validate-mock-a
make validate-doc-api
make validate-sdk-beta
make validate-sdk-mock-smoke
make validate-real-k8s-profile
go test ./services/ani-gateway/internal/router ./pkg/adapters/runtime
go test ./pkg/adapters/runtime ./pkg/bootstrap -run 'TestLocalWorkloadReconcileController|TestNewCapabilitiesDefaults|TestConfigEnvironmentOverridesWorkloadReconcileController|TestStartWorkloadReconcileControllerRequiresOptIn' -v
git diff --check
```

> 在没有联网依赖缓存时，`go test` 可能需要下载 Go module；本地可复用 `Makefile` 中的 `GOCACHE`/`GOMODCACHE` 设置。
