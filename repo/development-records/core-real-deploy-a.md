# CORE-REAL-DEPLOY-A — Sprint 11 Core real deployment validation profile

完成日期：2026-06-04
对应 Sprint：Sprint 11（Core Real Deployment Validation）
验证结果：`make validate-sprint11-core-real-deployment`

## 背景

Sprint 11 需要把 Sprint 6-10 的 Core 成果放到真实物理服务器上验一遍，但真实部署验证必须先区分只读检查、可重复门禁和需要人工审批的写操作。本批次新增聚合 profile，固定 Sprint 11 的安全边界。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/real-k8s-lab/sprint11-core-real-deployment.yaml` | 新增 Sprint 11 Core real deployment validation profile |
| `scripts/validate_sprint11_core_real_deployment.py` | 校验 profile scope、只读状态、K8s/KubeVirt/storage 观测、必需 gates、禁止 Services 越界和禁止新增 REAL-K8S-LAB guard |
| `scripts/validate_sprint11_core_real_deployment_test.py` | 覆盖默认 profile、mutation 开关、Rook-Ceph 未审批安装状态和缺失 storage gate 的拒绝 |
| `Makefile` | 新增 `validate-sprint11-core-real-deployment` 和 `validate-sprint11-real-deployment` |

## 真实环境只读结论

- 三节点 Kubernetes Ready，KubeVirt phase 为 `Deployed`。
- 当前无 StorageClass，存在 Pending PVC，说明 VM-first 持久化存储还没有生产化底座。
- Rook-Ceph namespace 未安装；Sprint 11 当前不直接安装 Rook-Ceph。

## 边界

本批次只证明 real deployment validation profile 和只读盘点结果，不代表生产部署完成、Rook-Ceph 可用、块存储可用或实际 `v1.0.0` 发布。
