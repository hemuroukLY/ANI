# K8s 集群

## 页面定位

`K8s 集群` 是 `Console / 算力与云资源` 下的租户侧集群管理页面，用于帮助用户查看、创建、升级、删除和维护自己权限范围内的 K8s 集群及节点池，并获取 kubeconfig 或通过受控代理访问集群原生 API。

本页是 `Console` 页面，不是 `BOSS` 的平台集群资源池运营页。

## 文档管理规则

- 本文件是 `K8s 集群` 的主维护源
- `prototypes/ani-services-prototype-console.html` 只保留页面摘要、模块边界和详细材料入口
- `tasks/modules/prd/console/compute/prd-console-k8s-cluster-management.md` 与 `tasks/modules/spec/console/compute/spec-console-k8s-cluster-management.md` 作为阶段性产物保留，不替代本文件

## Core 层要求

- K8s 集群能力属于 `Core / K8sClusters`
- 查询与写操作接口必须使用 `/api/v1/k8s-clusters*`
- 当前页面不继续使用旧 `/api/v1/console/*`
- 页面不要求前端显式传 `tenant_id`
- 创建、升级、节点池写操作必须按现有 schema 承接 `idempotency_key`
- 当前页面不把工作负载、事件、Helm、应用市场写成已冻结独立路径
- 受控访问原生 API 的正式入口是 `POST /api/v1/k8s-clusters/{cluster_id}/proxy`

## 页面职责

- 展示当前租户下的 K8s 集群列表与详情
- 提供创建、删除、升级集群入口
- 提供 kubeconfig 获取入口
- 提供节点池增删改查入口
- 提供受控的原生 API 代理入口

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 集群列表 | Core | 对齐 `K8sClusterListResponse` | 集群详情 |
| 集群详情 | Core | 对齐 `K8sCluster` | 节点池 / kubeconfig |
| 创建与升级 | Core | 对齐 `K8sClusterCreateRequest` / `K8sClusterUpgradeRequest` | 集群详情 |
| 节点池管理 | Core | 对齐 `K8sClusterNodePool*` | 节点池详情 |
| 受控代理 | Core | 对齐 `proxyK8sClusterAPI` | 代理响应 |
| 边界说明 | 规划项 | 标明工作负载与事件不在当前正式路径内 | 后续模块 |

## 字段级定义

### 集群字段

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 集群 ID | 文本 |
| name | 集群名称 | 文本 |
| version | K8s 版本 | 文本 |
| state | 集群状态 | 标签 |
| reason | 状态原因 | 文本 |
| created_at | 创建时间 | 时间 |
| updated_at | 更新时间 | 时间 |

### 节点池字段

| 字段 | 说明 | 展示建议 |
|---|---|---|
| id | 节点池 ID | 文本 |
| name | 节点池名称 | 文本 |
| node_count | 节点数 | 数字 |
| instance_type | 规格 | 文本 |
| gpu.vendor/model/count | GPU 意图 | 文本 |
| state | 节点池状态 | 标签 |

## 已冻结的 Core 能力

- `GET /api/v1/k8s-clusters`
- `POST /api/v1/k8s-clusters`
- `GET /api/v1/k8s-clusters/{cluster_id}`
- `DELETE /api/v1/k8s-clusters/{cluster_id}`
- `GET /api/v1/k8s-clusters/{cluster_id}/kubeconfig`
- `POST /api/v1/k8s-clusters/{cluster_id}/upgrade`
- `GET /api/v1/k8s-clusters/{cluster_id}/node-pools`
- `POST /api/v1/k8s-clusters/{cluster_id}/node-pools`
- `GET /api/v1/k8s-clusters/{cluster_id}/node-pools/{node_pool_id}`
- `PATCH /api/v1/k8s-clusters/{cluster_id}/node-pools/{node_pool_id}`
- `DELETE /api/v1/k8s-clusters/{cluster_id}/node-pools/{node_pool_id}`
- `POST /api/v1/k8s-clusters/{cluster_id}/proxy`

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录与租户上下文 | 已认证 | `401` / `403` |
| `name` | 租户内唯一 | `409 CONFLICT` |
| `version` | 支持的 K8s 版本 | `422 PRECONDITION_FAILED`（**YAML 已声明**，Core v1.yaml Phase 2 2026-06-17）；具体 `code` 待 Core description 补充；建议语义：版本不支持<!-- ADDED-TO-YAML: POST /api/v1/k8s-clusters 增 422 PreconditionFailed (Core v1.yaml, Phase 2 2026-06-17) --> |
| 租户集群配额（若启用） | 未超限 | `422 PRECONDITION_FAILED`（**YAML 已声明**，同上）；具体 `code` 待 Core description 补充；建议语义：配额超限<!-- ADDED-TO-YAML: 同上 POST /api/v1/k8s-clusters (Phase 2 2026-06-17) --> |

## 操作可用性矩阵

基于 `K8sCluster.state`：`provisioning` / `running` / `deleting`。

| 操作 | provisioning | running | deleting |
|---|---|---|---|
| 查看详情 | ✅ | ✅ | ⚠️ 只读 |
| 获取 kubeconfig | ❌ | ✅ | ❌ |
| 升级集群 | ❌ | ✅ | ❌ |
| 管理节点池 | ❌ | ✅ | ❌ |
| 原生 API 代理 | ❌ | ✅ | ❌ |
| 删除集群 | ⚠️ 二次确认 | ⚠️ 二次确认 | ❌ |

`provisioning` 期间写操作置灰；冲突返回 `409 CONFLICT`。`POST /k8s-clusters` 已在 YAML 声明 `422`（Phase 2 2026-06-17），前置条件类失败产品语义见上表；具体 code 待 Core description 补充。

## 工作负载（Phase 2 已声明）

- 详文：`k8s-workloads.md`
- `GET /api/v1/k8s-clusters/{cluster_id}/workloads` <!-- ADDED-TO-YAML: Core v1.yaml, Phase 2 2026-06-17 -->

## 待补 Core 契约能力

- 独立事件页路径
- Helm / 应用市场
- 平台级集群运营总览

## 接口冻结规则

### 集群主资源

| 项 | 路径 |
|---|---|
| 列表 | `GET /api/v1/k8s-clusters` |
| 创建 | `POST /api/v1/k8s-clusters` |
| 详情 | `GET /api/v1/k8s-clusters/{cluster_id}` |
| 删除 | `DELETE /api/v1/k8s-clusters/{cluster_id}` |

### 集群扩展能力

| 项 | 路径 |
|---|---|
| kubeconfig | `GET /api/v1/k8s-clusters/{cluster_id}/kubeconfig` |
| 升级 | `POST /api/v1/k8s-clusters/{cluster_id}/upgrade` |
| 节点池列表/创建 | `GET/POST /api/v1/k8s-clusters/{cluster_id}/node-pools` |
| 节点池详情/更新/删除 | `GET/PATCH/DELETE /api/v1/k8s-clusters/{cluster_id}/node-pools/{node_pool_id}` |
| 原生 API 代理 | `POST /api/v1/k8s-clusters/{cluster_id}/proxy` |

## 使用规则

- 创建集群时至少承接 `name`、`idempotency_key`
- 升级集群时至少承接 `version` 与 `idempotency_key`
- 节点池写操作必须沿用对应 schema 的 `idempotency_key`
- 工作负载如需展示，优先引用 `k8s-workloads.md` 中 Phase 2 已声明路径；事件仍只能作为 `proxy` 或后续模块能力

## 验收标准

- 正文路径、返回码、字段与现有 `Core v1.yaml` 一致
- 正文不把工作负载、事件和 Helm 写成当前模块正式接口
- HTML 摘要、PRD、SPEC 与本文件一致
- 本文件可以独立作为 `Console / K8s 集群` 的主维护材料

## 产品经理补充定义

### 目标用户与权限视角

- 租户管理员：查看集群、节点池、kubeconfig、升级与代理入口
- 平台使用成员：查看集群状态、节点池状态和访问入口，但是否可执行变更以后端权限为准
- 只读用户：仅查看摘要、详情与最近操作，不显示高风险变更动作

### 默认视图与页面状态

- 首屏默认展示集群列表与状态摘要，无集群时提供创建入口
- 长任务场景必须把“进行中”与“失败待处理”独立展示，不能只靠列表状态猜测
- kubeconfig 或代理能力不可用时，应保留说明与重试入口，而不是静默隐藏
- 集群、节点池、升级三类长任务失败时必须提供回看入口和错误上下文

### 核心任务流

1. 用户从列表发起集群创建，随后进入详情或操作反馈查看创建进度
2. 用户在详情中管理节点池、申请 kubeconfig 或查看代理入口
3. 用户在升级或变更节点池后，通过任务状态和详情页复核最终结果

### 跨模块协同

- 创建集群依赖网络、存储、GPU 等基础资源时，只通过跳转引用，不在本页复制其契约
- 集群异常可从首页或告警页回跳到本页并带入状态筛选
- 工作负载、Helm、事件等未冻结能力只保留边界说明，不形成子页承诺

### 产品验收补充

- 用户必须能分辨“列表状态”“长任务状态”“升级状态”这三类不同反馈
- 每个长任务都必须存在明确回流点，避免用户刷新后丢失上下文
- `kubeconfig`、`proxy` 等入口不可用时必须有可读说明
- 本页不得把工作负载管理写成当前 K8s 集群模块的已冻结能力
