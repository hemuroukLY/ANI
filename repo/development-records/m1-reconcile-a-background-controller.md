# M1-RECONCILE-A · Workload Reconcile Controller

- 完成日期：2026-05-22
- 状态：✅

实现了 `WorkloadReconcileController` 的 local adapter 和默认关闭的 bootstrap opt-in 运行剖面：后台 controller 通过 `ReconcileTargetLister` 扫描需要对齐的实例，调用 `WorkloadProviderStatusReader.Observe` 获取 provider observation，再交给 `WorkloadStatusReconciler.Reconcile` 统一映射并通过 `WorkloadInstanceStore.UpsertStatus` 回写状态。使用 `WORKLOAD_RECONCILE_CONTROLLER_ENABLED=true` 时，`bootstrap.RunGRPC` 会随服务生命周期启动/停止 controller。

## 关键实现

- `pkg/ports/reconcile_controller.go`：补齐 `ReconcileTargetListRequest` 和 `ReconcileTargetLister`。
- `pkg/adapters/runtime/reconcile_controller.go`：新增 `LocalWorkloadReconcileController`，支持 `Start` 循环和 `ReconcileNow`。
- `pkg/adapters/runtime/instance_store.go`：`MetadataInstanceStore` 支持 `ListReconcileTargets`，通过 platform tx 扫描非终态或 stale 实例。
- `pkg/bootstrap/deps.go`：将 `WorkloadController` 暴露为 bootstrap capability，并注入 reconcile 轮询配置。
- `pkg/bootstrap/server.go`：新增 controller 运行开关与环境变量覆盖；`RunGRPC` 按 opt-in 配置启动后台循环，默认关闭。
- `pkg/bootstrap/server_test.go`：覆盖环境变量覆盖、默认不启动和显式启动。
- `pkg/adapters/runtime/reconcile_controller_test.go`：覆盖状态回写、provider missing 转 failed、target lister 调用。

## 真实边界

本批次完成 controller adapter、capability 与 bootstrap opt-in 后台运行，不包含：

- 独立 worker 镜像/部署形态，或 ani-gateway 主入口接入。
- 真实 vCluster/Kubernetes provider 的生产状态观察。
- controller 指标、leader election、分布式锁和重试退避。
- 对 Secret binding、K8s cluster provider 的真实状态回写。
