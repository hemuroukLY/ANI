# M1-K8S-B · K8s API Proxy Dev Profile

- 完成日期：2026-05-22
- 状态：✅ local dev profile

实现了 K8s 集群原生 API proxy 的 Core 管控面契约与 local dev profile。当前端点为 `POST /api/v1/k8s-clusters/{cluster_id}/proxy`，请求体表达 method/path/query/body，响应返回模拟的 Kubernetes 风格 JSON。该切片用于稳定 Core API、SDK 和 Services 调用边界，不会真实转发到 vCluster API Server。

## 关键实现

- `api/openapi/v1.yaml`：新增 `K8sClusterProxyRequest`、`K8sClusterProxyResponse` 和 `proxyK8sClusterAPI`。
- `pkg/ports/k8s_clusters.go`：在 `K8sClusterService` 中新增 `Proxy` 能力，保持 Gateway 不直接调用 K8s SDK。
- `pkg/adapters/runtime/local_k8s_cluster_service.go`：新增 local proxy 实现，校验 running cluster、`idempotency_key`、K8s API path allowlist，并返回模拟响应。
- `services/ani-gateway/internal/router/k8s_cluster_resources.go`：新增 `POST /k8s-clusters/{id}/proxy` 路由。
- `services/ani-gateway/internal/router/k8s_cluster_resources_test.go`：覆盖 proxy local profile、dev_profile 标识和非法 path 拒绝。
- `docs/api/` 与 `sdks/core/`：由 OpenAPI 重新生成，保持 API docs 和四语言 SDK metadata 同步。

## 真实边界

本批次完成 proxy 的 Core API 契约和 local profile，不包含：

- 真实 vCluster 创建、升级、删除。
- 将 proxy 请求真实转发到 vCluster API Server。
- kubectl/Helm 透明代理兼容性。
- proxy 审计日志、限流、长连接、watch/exec/logs 流式转发。
