# SPRINT13-NETROUTE-KUBEOVN-LIVE-A — 网络路由 Kube-OVN live gate 结果

> 记录类型：Sprint 13 B-track live result
> 完成日期：2026-06-20
> 范围：仅 ANI Core S01 网络路由 Kube-OVN real-provider evidence；不代表 production ready
> 状态：**real-provider evidence passed for S01 route gate, production-shaped gate pending**

## 目标

在人工确认真实写操作后，对 Sprint 13 S01 网络路由 Kube-OVN 执行真实 live gate，证明临时 Vpc/Subnet/route/NetworkPolicy/Service 能在真实 Kubernetes + Kube-OVN lab 中 create/observe，并在成功后清理临时资源。

## §153 五项实测结果

| 项 | 实测结果 |
|---|---|
| 当前状态 | S01 网络路由 Kube-OVN real-provider evidence passed。Core route provider 已通过 `NETWORK_PROVIDER=kubeovn_rest` 显式路径接入 `RenderRoute -> DryRun -> Apply -> Observe`；Gateway 可注入 provider-backed `ports.NetworkService`，handler 不绕过 port。 |
| 真实组件 + 版本 | Kubernetes `v1.36.1`，Kube-OVN `v1.15.8`。执行前已确认三节点 Ready、Kube-OVN 组件 Running，`Vpc.spec.staticRoutes` API 字段存在。 |
| live gate 命令 | `python scripts/validate_kubeovn_network_live_gate.py --live --cleanup --tenant-id sprint13-s01-b --vpc-name sprint13-s01-route --subnet-name sprint13-s01-subnet --route-name sprint13-s01-route --security-group-name sprint13-s01-sg --load-balancer-name sprint13-s01-lb --cidr 10.244.180.0/24 --gateway 10.244.180.1 --route-destination 10.250.0.0/16 --route-next-hop 10.244.180.1 --evidence-output development-records/live-evidence/sprint13-netroute-kubeovn-live-evidence.json` |
| evidence 输出路径 | `repo/development-records/live-evidence/sprint13-netroute-kubeovn-live-evidence.json` |
| 失败边界 | 本次已通过，因此 S01 route provider 可标 real-provider evidence passed；但不代表 production ready，不证明生产 RBAC/凭据管理、持久 route 元数据、`instance`/`nat` next hop 映射、外部负载均衡可达性或其他 Sprint 13 切片已完成。Production-shaped gate: **PENDING**，`production_shape.status=pending`。 |

## 关键输出

```text
M1-NETWORK-LIVE-A live checks valid; evidence written to development-records/live-evidence/sprint13-netroute-kubeovn-live-evidence.json
```

Evidence 摘要：

```json
{
  "id": "kubeovn-network-live-gate",
  "profile": "M1-NETWORK-LIVE-A",
  "status": "passed",
  "namespace": "ani-tenant-sprint13-s01-b",
  "vpc": "vpc-sprint13-s01-route",
  "subnet": "subnet-sprint13-s01-subnet",
  "route": "route-sprint13-s01-route",
  "security_group": "sg-sprint13-s01-sg",
  "load_balancer": "lb-sprint13-s01-lb",
  "production_shape": {
    "missing_items": [
      "production_rbac_and_credential_management",
      "persistent_route_metadata_reconciliation"
    ],
    "status": "pending",
    "transport_profile": "lab_kubeconfig"
  },
  "cleanup": {
    "status": "deleted",
    "resources": [
      "service/lb-sprint13-s01-lb",
      "networkpolicy/sg-sprint13-s01-sg",
      "subnet/subnet-sprint13-s01-subnet",
      "vpc/vpc-sprint13-s01-route",
      "namespace/ani-tenant-sprint13-s01-b"
    ]
  }
}
```

## 清理核验

已使用真实 Kubernetes API 查询以下临时资源，均为空输出：

```bash
kubectl get ns ani-tenant-sprint13-s01-b --ignore-not-found
kubectl get vpc vpc-sprint13-s01-route --ignore-not-found
kubectl get subnet subnet-sprint13-s01-subnet --ignore-not-found
kubectl get networkpolicy sg-sprint13-s01-sg -n ani-tenant-sprint13-s01-b --ignore-not-found
kubectl get service lb-sprint13-s01-lb -n ani-tenant-sprint13-s01-b --ignore-not-found
```

## 代码与契约边界

- `ports.NetworkProviderRenderer.RenderRoute` 已纳入 port，route 渲染经 `KubeOVNNetworkRenderer` 输出 `Vpc.spec.staticRoutes` manifest。
- `LocalNetworkService.CreateRoute` 在显式 provider 配置下执行 `RenderRoute -> DryRun -> Apply -> Observe`；默认未配置时仍为 Tier1 local profile。
- Gateway `RegisterWithOptions` 支持注入 `ports.NetworkService`，`services/ani-gateway/main.go` 可由 `NETWORK_PROVIDER=kubeovn_rest` 构造 provider-backed network service。
- `scripts/validate_kubeovn_network_live_gate.py --cleanup` 成功后清理临时 Service、NetworkPolicy、Subnet、Vpc 与 Namespace。

## 非目标

- 不声明 network route production ready。
- Production-shaped gate: **PENDING**；`production_shape.status=pending`。进入生产形态前必须补齐 `production_rbac_and_credential_management` 与 `persistent_route_metadata_reconciliation`，并用正式控制面凭据/RBAC 与持久化 route reconcile 路径重新产出 evidence。
- 不声明 S02-S07 已完成真实 live gate。
- 不把 Kube-OVN external LB 可达性、生产镜像供应链、生产凭据管理或多租户长期资源生命周期纳入本次结论。
