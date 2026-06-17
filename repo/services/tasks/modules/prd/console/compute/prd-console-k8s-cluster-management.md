# PRD: Console K8s 集群

## 1. Introduction/Overview

`Console / 算力与云资源 / K8s 集群` 用于帮助租户用户查看、创建、升级、删除和维护自己权限范围内的 K8s 集群及节点池，并获取 kubeconfig 或通过受控代理访问集群原生 API。

本轮 PRD 严格按现有 `Core v1.yaml` 约束生成，当前已核实的页面边界如下：

- 当前 K8s 集群能力拥有独立 `Core` 路径组 `/api/v1/k8s-clusters*`
- 当前已冻结的集群能力包括：列表、详情、创建、删除、获取 kubeconfig、升级、节点池增删改查、代理访问原生 API
- 当前 `Core v1.yaml` 没有冻结独立的“工作负载页”“事件页”正式路径，相关能力仅能通过 `proxy` 或后续模块承接
- 当前页面属于 `Console` 租户侧 K8s 集群管理页，不是 `BOSS` 的平台集群资源池运营页

## 2. Goals

- 让租户管理自己的 K8s 集群与节点池
- 让页面可以获取 kubeconfig 并进行基础集群操作
- 严格以现有 `Core v1.yaml` 为准，沉淀可直接用于 Core 对齐的主维护材料
- 明确集群页与工作负载、事件观测、平台运营视角的边界

## 3. User Stories

### US-001: 查看 K8s 集群列表
**Description:** 作为租户用户，我希望查看当前租户下的 K8s 集群列表，以便快速找到需要维护的集群。

**Acceptance Criteria:**
- [ ] 页面展示集群名称、版本、状态、原因、创建时间、更新时间
- [ ] 查询接口对齐 `GET /api/v1/k8s-clusters`
- [ ] 支持 `limit` 和 `cursor`

### US-002: 创建和删除 K8s 集群
**Description:** 作为租户用户，我希望创建或删除 K8s 集群，以便交付或回收集群资源。

**Acceptance Criteria:**
- [ ] 创建接口对齐 `POST /api/v1/k8s-clusters`
- [ ] 创建请求体对齐 `K8sClusterCreateRequest`
- [ ] 删除接口对齐 `DELETE /api/v1/k8s-clusters/{cluster_id}`
- [ ] 创建成功返回 `201 + K8sCluster`

### US-003: 获取 kubeconfig 和升级集群
**Description:** 作为租户用户，我希望获取 kubeconfig 或升级 K8s 版本，以便进行日常维护。

**Acceptance Criteria:**
- [ ] kubeconfig 接口对齐 `GET /api/v1/k8s-clusters/{cluster_id}/kubeconfig`
- [ ] 升级接口对齐 `POST /api/v1/k8s-clusters/{cluster_id}/upgrade`
- [ ] 升级请求体对齐 `K8sClusterUpgradeRequest`
- [ ] 页面正确处理 `409 CONFLICT`

### US-004: 管理节点池
**Description:** 作为租户用户，我希望查看、创建、修改和删除节点池，以便调整集群容量和规格。

**Acceptance Criteria:**
- [ ] 节点池接口对齐 `/api/v1/k8s-clusters/{cluster_id}/node-pools*`
- [ ] 写操作按现有 schema 使用 `idempotency_key`
- [ ] 页面可承接 CPU/GPU 意图，但不扩写平台级资源池调度逻辑

### US-005: 通过代理访问原生 API
**Description:** 作为租户用户，我希望在受控前提下代理访问集群原生 API，以便实现受限运维操作。

**Acceptance Criteria:**
- [ ] 接口对齐 `POST /api/v1/k8s-clusters/{cluster_id}/proxy`
- [ ] 页面明确该能力是“受控代理”，不是直接暴露 provider SDK
- [ ] 页面不把原生工作负载 CRUD 写成当前模块独立冻结路径

## 4. Functional Requirements

- FR-1: 系统必须提供 K8s 集群列表与详情视图
- FR-2: 系统必须支持创建与删除 K8s 集群
- FR-3: 系统必须支持获取 kubeconfig 和升级集群
- FR-4: 系统必须支持节点池增删改查
- FR-5: 系统必须支持代理访问集群原生 API
- FR-6: 系统必须保持租户侧集群管理而非平台资源池运营口径

## 5. Non-Goals (Out of Scope)

- 不在本轮把工作负载、事件、应用商店写成独立冻结路径
- 不在本轮实现平台级集群运营总览
- 不在本轮直接暴露 provider SDK 或长期凭据

## 6. Technical Considerations

- 正式路径为 `/api/v1/k8s-clusters*`
- 创建、升级、节点池写操作需要按 schema 承接 `idempotency_key`
- 代理访问仅通过 `POST /api/v1/k8s-clusters/{cluster_id}/proxy`
- 错误响应需对齐 `400/404/409`

## 7. Success Metrics

- 用户能完成集群查看、创建、删除、升级和节点池维护
- 文档中的路径、字段与返回结构与 `Core v1.yaml` 一致
- 页面不再把工作负载/事件写成超出当前冻结范围的正式接口

## 8. Open Questions

- 后续是否需要把工作负载管理拆为独立模块
- 事件与可观测是否进入单独的 Console 集群观察页
