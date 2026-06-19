# SPRINT13-NETROUTE-KUBEOVN-A-TRACK — 网络路由 Kube-OVN A 轨

> 批次类型：Sprint 13 A-track / code+contract ready
> 日期：2026-06-19
> 范围：仅 ANI Core；本 A 轨不跑真实集群写操作；后续 S01 B 轨 live gate 已通过，见 `sprint13-netroute-kubeovn-live-result.md`

## 背景

Sprint 12 已完成 `/api/v1/networks/routes` 的 Core handler、`ports.NetworkService.CreateRoute/ListRoutes` 与 Tier1 local profile。Sprint 13 S01 的 A 轨目标是在不改 port 签名、不改 Gateway handler、不触碰 `/api/v1/svc` 的前提下，把网络路由推进到 Kube-OVN adapter code + live-gate contract ready。

## 完成内容

- 在 `KubeOVNNetworkRenderer` 增加 adapter-local `RenderRoute`，将 gateway route 渲染为 Kube-OVN `Vpc.spec.staticRoutes` patch manifest。后续 B 轨 prep 已把 `RenderRoute` 纳入 `ports.NetworkProviderRenderer`，以支持 Core route provider 接线。
- route renderer 当前只支持 `next_hop_type=gateway`，把 `next_hop_id` 作为 `nextHopIP`；`instance` / `nat` 的真实映射未确认前不冒充支持。
- 增加 renderer 单测，覆盖 route manifest 的 `staticRoutes`、`cidr`、`nextHopIP`、`policyDst` 与 route annotation。
- 增加 provider dry-run fake 单测，确认 route manifest 可进入既有 `KubeOVNNetworkProviderAdapter.DryRun` 管线并生成 `kubeovn/Vpc/...` resource ref。
- 扩展 `validate_kubeovn_network_live_gate.py`、测试和 `deploy/real-k8s-lab/kubeovn-network-live-gate.yaml`，加入 `kubeovn-route-created` contract，fake live runner 覆盖 route manifest apply 与 Vpc static route observe。
- 未新增组件 SDK 依赖，因此不需要更新 component import allowlist。

## 边界

- 本批没有执行 `--live`，没有对真实 Kube-OVN / Kubernetes 集群做写操作，没有产出真实 evidence JSON。
- A 轨完成时 `NetworkService` route 方法仍保持 local profile；2026-06-20 的 B 轨 prep 已补显式 `NETWORK_PROVIDER=kubeovn_rest` forwarding，并要求 provider execution user/proof 由配置提供，不从 request 伪造 proof。
- 真实 B 轨执行前必须在 lab 确认 Kube-OVN `v1.15.8` 的静态路由真实 API 形态。如果 `Vpc.spec.staticRoutes` 与部署版本不一致，必须先修正 adapter 与 live-gate contract。
- 本 A 轨自身不证明真实 route apply/observe；后续 S01 B 轨已补真实 evidence。即便 S01 route provider evidence passed，也不得标 production ready。

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

后续状态：2026-06-20 的 B 轨 prep 已补 Core route provider 与 Gateway runtime 接线；同日真实 B 轨已运行 Kube-OVN route live gate，并写入 `sprint13-netroute-kubeovn-live-result.md` 与 `live-evidence/sprint13-netroute-kubeovn-live-evidence.json`。
