# SPEC: Console K8s 集群

> Technical specification derived from: `tasks/modules/prd/console/compute/prd-console-k8s-cluster-management.md`
> Generated: 2026-06-13 | Scope: `Console / 算力与云资源 / K8s 集群`

## 1. Summary

本 SPEC 定义 `Console / 算力与云资源 / K8s 集群` 的技术边界、字段映射、接口约束和回填规则。该页面直接承接 `Core / K8sClusters` 已冻结能力，不扩写为工作负载平台或集群运营平台。

## 2. Architecture

### 2.1 System Context

```text
浏览器 Console
    -> Core K8sClusters API (/api/v1/k8s-clusters*)
        -> 集群列表 / 详情 / 创建 / 删除
        -> kubeconfig
        -> upgrade
        -> node-pools
        -> proxy
```

约束说明：

1. 本页直接读取 `Core` 的租户集群能力，不要求 `Services` 再包装
2. 工作负载、事件等原生对象不在当前页面形成独立正式路径
3. 受控原生 API 访问通过 `proxy` 实现，不直接暴露 provider 访问方式

## 3. Data Model

```ts
interface K8sClusterListItem {
  id?: string;
  name?: string;
  version?: string;
  state?: "provisioning" | "running" | "deleting";
  reason?: string;
  createdAt?: string;
  updatedAt?: string;
}

interface K8sClusterNodePoolView {
  id?: string;
  name?: string;
  nodeCount?: number;
  instanceType?: string;
  gpu?: {
    vendor?: string;
    model?: string;
    count?: number;
    resourceName?: string;
  };
  state?: "running" | "deleting";
}
```

## 4. API Design

### 4.1 Endpoints

| Method | Path | operationId | summary | Success |
|---|---|---|---|---|
| GET | `/api/v1/k8s-clusters` | `listK8sClusters` | 列出 K8s 集群 | `200 + K8sClusterListResponse` |
| POST | `/api/v1/k8s-clusters` | `createK8sCluster` | 创建 K8s 集群 | `201 + K8sCluster` |
| GET | `/api/v1/k8s-clusters/{cluster_id}` | `getK8sCluster` | 查询 K8s 集群 | `200 + K8sCluster` |
| DELETE | `/api/v1/k8s-clusters/{cluster_id}` | `deleteK8sCluster` | 删除 K8s 集群 | `200 + K8sCluster` |
| GET | `/api/v1/k8s-clusters/{cluster_id}/kubeconfig` | `getK8sClusterKubeconfig` | 获取 kubeconfig | `200 + K8sClusterKubeconfig` |
| POST | `/api/v1/k8s-clusters/{cluster_id}/upgrade` | `upgradeK8sCluster` | 升级集群版本 | `200 + K8sCluster` |
| GET | `/api/v1/k8s-clusters/{cluster_id}/node-pools` | `listK8sClusterNodePools` | 列出节点池 | `200 + K8sClusterNodePoolListResponse` |
| POST | `/api/v1/k8s-clusters/{cluster_id}/node-pools` | `createK8sClusterNodePool` | 创建节点池 | `201 + K8sClusterNodePool` |
| GET | `/api/v1/k8s-clusters/{cluster_id}/node-pools/{node_pool_id}` | `getK8sClusterNodePool` | 获取节点池 | `200 + K8sClusterNodePool` |
| PATCH | `/api/v1/k8s-clusters/{cluster_id}/node-pools/{node_pool_id}` | `updateK8sClusterNodePool` | 更新节点池 | `200 + K8sClusterNodePool` |
| DELETE | `/api/v1/k8s-clusters/{cluster_id}/node-pools/{node_pool_id}` | `deleteK8sClusterNodePool` | 删除节点池 | `200 + K8sClusterNodePool` |
| POST | `/api/v1/k8s-clusters/{cluster_id}/proxy` | `proxyK8sClusterAPI` | 代理访问原生 API | `200 + K8sClusterProxyResponse` |

### 4.2 Schema Rules

- 创建集群请求体对齐 `K8sClusterCreateRequest`
- 升级请求体对齐 `K8sClusterUpgradeRequest`
- 节点池创建请求体对齐 `K8sClusterNodePoolCreateRequest`
- 节点池更新请求体对齐 `K8sClusterNodePoolUpdateRequest`
- 写操作中出现的 `idempotency_key` 必须沿用现有 schema

### 4.2.1 创建前置条件

产品建议语义见主维护源 `docs/console-modules/compute/k8s-cluster-management.md`。**当前 `POST /k8s-clusters` YAML 未声明 `422`**。

### 4.3 Non-Frozen Capabilities

以下能力当前**不能**写成已对齐 Core 的正式接口：

- 独立工作负载管理路径
- 独立事件页路径
- Helm / 应用市场
- 平台级资源池运营视图

## 5. Business Logic

### 5.1 Core Flows

#### 集群列表

1. 页面打开时请求 `GET /api/v1/k8s-clusters`
2. 用户点击某个集群后加载详情和节点池

#### 创建集群

1. 用户填写名称和版本
2. 提交 `POST /api/v1/k8s-clusters`
3. 成功后返回 `201 + K8sCluster`

#### 节点池维护

1. 用户进入节点池页
2. 读取 `GET /node-pools`
3. 按需执行创建、更新、删除

#### 受控代理

1. 用户发起受控 API 代理请求
2. 页面调用 `POST /proxy`
3. 返回代理响应，不展示 provider 内部细节

## 6. Error Handling

- 集群详情、kubeconfig、节点池读取支持 `404`
- 升级、节点池更新、代理访问应正确处理 `409 CONFLICT`
- 错误结构统一为：

```json
{
  "code": "UPPER_SNAKE",
  "message": "error message",
  "request_id": "req-20260613-001"
}
```

## 7. Testing Strategy

- 校验所有 `k8s-clusters*` 路径、返回码、schema 与 `Core v1.yaml` 对齐
- 校验工作负载、事件没有被写成独立冻结路径
- 校验节点池写操作正确沿用 `idempotency_key`

## 8. Implementation Plan

1. 生成 PRD 和 SPEC
2. 收口主维护文档
3. 回填 K8s HTML 摘要
4. 执行最终 Core 合规复审
