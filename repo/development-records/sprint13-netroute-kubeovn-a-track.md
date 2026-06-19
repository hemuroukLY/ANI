# SPRINT13-NETROUTE-KUBEOVN-A-TRACK — 网络路由 Kube-OVN A 轨

> 批次类型：Sprint 13 A-track / code+contract ready
> 日期：2026-06-19
> 范围：仅 ANI Core；不跑真实集群写操作；状态为 LIVE PENDING

## 背景

Sprint 12 已完成 `/api/v1/networks/routes` 的 Core handler、`ports.NetworkService.CreateRoute/ListRoutes` 与 Tier1 local profile。Sprint 13 S01 的 A 轨目标是在不改 port 签名、不改 Gateway handler、不触碰 `/api/v1/svc` 的前提下，把网络路由推进到 Kube-OVN adapter code + live-gate contract ready。

## 完成内容

- 在 `KubeOVNNetworkRenderer` 增加 adapter-local `RenderRoute`，将 gateway route 渲染为 Kube-OVN `Vpc.spec.staticRoutes` patch manifest；不扩展 `ports.NetworkProviderRenderer` 接口签名。
- route renderer 当前只支持 `next_hop_type=gateway`，把 `next_hop_id` 作为 `nextHopIP`；`instance` / `nat` 的真实映射未确认前不冒充支持。
- 增加 renderer 单测，覆盖 route manifest 的 `staticRoutes`、`cidr`、`nextHopIP`、`policyDst` 与 route annotation。
- 增加 provider dry-run fake 单测，确认 route manifest 可进入既有 `KubeOVNNetworkProviderAdapter.DryRun` 管线并生成 `kubeovn/Vpc/...` resource ref。
- 扩展 `validate_kubeovn_network_live_gate.py`、测试和 `deploy/real-k8s-lab/kubeovn-network-live-gate.yaml`，加入 `kubeovn-route-created` contract，fake live runner 覆盖 route manifest apply 与 Vpc static route observe。
- 未新增组件 SDK 依赖，因此不需要更新 component import allowlist。

## 边界

- 本批没有执行 `--live`，没有对真实 Kube-OVN / Kubernetes 集群做写操作，没有产出真实 evidence JSON。
- `NetworkService` route 方法仍保持 local profile；未开启 real-mode forwarding。原因是当前 `NetworkRouteCreateRequest` 不携带 provider execution 所需的 user / permission proof，本轮又禁止改 port 签名和 Gateway handler；不使用 fake proof 伪造真实 provider 接线。
- 真实 B 轨执行前必须在 lab 确认 Kube-OVN `v1.15.8` 的静态路由真实 API 形态。如果 `Vpc.spec.staticRoutes` 与部署版本不一致，必须先修正 adapter 与 live-gate contract。
- 未跑通真实 route apply/observe 前，网络路由只可标 Tier1 local profile，不得标 real-provider、runtime ready 或 production ready。

## 验证命令

本批最终提交前按 Sprint 13 loop-safe gate 执行：

```bash
cd repo && make test && make validate-network-alpha validate-kubeovn-network-live-gate && python scripts/validate_yaml.py api/openapi/v1.yaml && make validate-doc-entrypoints && git diff --check
```

## 关键文件

- `pkg/adapters/runtime/kubeovn_network_renderer.go`
- `pkg/adapters/runtime/kubeovn_network_renderer_test.go`
- `pkg/adapters/runtime/kubeovn_network_provider_test.go`
- `scripts/validate_kubeovn_network_live_gate.py`
- `scripts/validate_kubeovn_network_live_gate_test.py`
- `deploy/real-k8s-lab/kubeovn-network-live-gate.yaml`
- `development-records/sprint13-netroute-kubeovn-readiness.md`
- `development-records/sprint13-loop-execution-prompts.md`

## 下一步

进入 S02 前，S01 保持 `code+contract ready, LIVE PENDING`。真实 B 轨由人工先做只读盘点与审批，再运行 Kube-OVN route live gate 并写入 `sprint13-netroute-kubeovn-live-result.md`。
