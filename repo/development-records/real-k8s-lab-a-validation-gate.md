# REAL-K8S-LAB-A · Real K8s Lab Validation Gate

- 完成日期：2026-05-23
- 状态：✅ contract gate

建立 Sprint 5 真实底座验证门禁。该批次不部署真实环境，而是把三台公有云 VM 上的 K8s/Kube-OVN/KubeVirt/vCluster 验证要求固化为 profile 和 Makefile 入口，避免后续把 local profile 误写成真实 provider 完成。

## 关键实现

- `deploy/real-k8s-lab/profile.yaml`：定义 REAL-K8S-LAB-A 最小真实底座组件、三节点要求和 live check。
- `scripts/validate_real_k8s_profile.py`：默认校验门禁定义和文档引用；`--live` 模式通过 `kubectl` 检查真实 lab。
- `Makefile`：新增 `make validate-real-k8s-profile`。
- `CLAUDE.md`、`ANI-DOCS-INDEX.md`、`ANI-06-开发计划.md`、`repo/CURRENT-SPRINT.md`：同步真实底座强制门禁。

## 使用方式

默认合同门禁：

```bash
make validate-real-k8s-profile
```

三台云 VM K8s 环境就绪后执行 live 门禁：

```bash
KUBECONFIG=/path/to/real-lab.kubeconfig python scripts/validate_real_k8s_profile.py --live
```

## 真实边界

本批次只完成验证门禁，不代表以下能力已经完成：

- 三台云 VM 已部署完成。
- K8s/Kube-OVN/KubeVirt/vCluster 已安装。
- vCluster 生命周期、真实 proxy 转发、KubeVirt VM 或 Kube-OVN 网络真实跑通。
- 生产级 HA、安全加固、故障恢复和性能基线。
