# SPRINT13-NETROUTE-KUBEOVN-B-TRACK-PREP — 网络路由 Kube-OVN B 轨接线准备

> 记录类型：Sprint 13 B-track preparation record
> 完成日期：2026-06-20
> 范围：仅 ANI Core 代码接线与本地契约门禁；真实集群写操作见后续 live result
> 状态：Core route provider wiring ready；后续 S01 live gate 已通过，见 `sprint13-netroute-kubeovn-live-result.md`

## 目标

把 Sprint 12 已落地的 `/api/v1/networks/routes` 从纯 Tier1 local profile 推进到可显式切换的 Kube-OVN route provider execution path，为 S01 真实 live gate 做代码准备。

本记录不是 `sprint13-netroute-kubeovn-live-result.md`。真实 Kube-OVN apply/observe evidence 已由后续 live result 归档；即便 S01 route provider evidence passed，仍不得标 production ready。

## 完成内容

- `ports.NetworkProviderRenderer` 增加 `RenderRoute`，把 route 渲染纳入正式 provider renderer port。
- `LocalNetworkService.CreateRoute` 支持显式配置 route provider，执行 `RenderRoute -> DryRun -> Apply -> Observe`，默认未配置时仍保持 local profile。
- bootstrap 与 Gateway runtime 增加 `NETWORK_PROVIDER=kubeovn_rest` 显式开关，配套 `NETWORK_PROVIDER_APPLY_ENABLED`、`NETWORK_PROVIDER_USER_ID`、`NETWORK_PROVIDER_PERMISSION_PROOF`。
- provider execution identity 必须由配置提供；代码不伪造 user/proof，不把凭据写入仓库。
- `validate_kubeovn_network_live_gate.py --live` 增加 `--cleanup`，成功 observe 后按 Service、NetworkPolicy、Subnet、Vpc、Namespace 顺序删除临时资源，并把 cleanup 结果写入 evidence。
- Gateway handler 与 Core OpenAPI 未修改。

## 历史边界与后续状态

- 本 prep 记录自身未执行 `validate_kubeovn_network_live_gate.py --live`，未创建真实 Namespace/Vpc/Subnet/staticRoute/NetworkPolicy/Service。
- 后续 `SPRINT13-NETROUTE-KUBEOVN-LIVE-A` 已执行真实 live gate，产出 `sprint13-netroute-kubeovn-live-result.md` 与 `live-evidence/sprint13-netroute-kubeovn-live-evidence.json`。
- route 状态仍只记录在 `LocalNetworkService` 内存路径；本批不扩展 metadata schema。
- S01 route provider 可标 real-provider evidence passed；不代表 production ready。

## 验证

已执行：

```bash
go test ./pkg/adapters/runtime ./pkg/bootstrap ./services/ani-gateway/internal/router
make validate-network-alpha validate-kubeovn-network-live-gate
python scripts/validate_kubeovn_network_live_gate_test.py
```

关键覆盖：

- `TestLocalNetworkServiceRouteCanUseKubeOVNProviderPipeline`
- `TestNewCapabilitiesCanWireKubeOVNNetworkRouteProvider`
- `TestNewCapabilitiesRejectsKubeOVNNetworkProviderWithoutExecutionProof`
- `TestConfigEnvironmentOverridesNetworkProvider`
- `test_live_gate_cleanup_deletes_temporary_resources_after_observe`

## 后续 B 轨结果

已在人工确认后执行带 `--cleanup` 的真实 live gate 写操作，产出非敏感 evidence JSON，清理临时资源，并新增 `sprint13-netroute-kubeovn-live-result.md`。S01 网络路由能力可标 real-provider evidence passed；不代表 runtime/production ready。
