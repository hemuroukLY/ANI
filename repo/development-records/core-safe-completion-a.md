# CORE-SAFE-COMPLETION-A — Sprint 11 Core safe completion profile

完成日期：2026-06-05
对应 Sprint：Sprint 11（Core Real Deployment Validation）
验证结果：`make validate-sprint11-safe-completion`

## 背景

用户明确要求 Sprint 11 必须“安全：不会造成数据丢失和无法重启的情况出现”。因此 Sprint 11 的完成条件不能是直接安装 Rook-Ceph 或挂载数据盘，而是先把安全完成条件写成机器可验证 profile：只读验证、稳定设备 ID、禁止 destructive 操作、人工审批前不改变服务器状态。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/real-k8s-lab/sprint11-core-safe-completion.yaml` | 新增 Sprint 11 safe completion profile |
| `scripts/validate_sprint11_safe_completion.py` | 校验 Core scope、安全完成状态、false destructive flags、开源最佳实践、只读证据、必需 profile、completion gates 和人工审批动作 |
| `scripts/validate_sprint11_safe_completion_test.py` | 覆盖默认 profile、mutation、Rook-Ceph 安装误标、缺失持久设备 ID 原则和缺失重启审批的拒绝场景 |
| `Makefile` | 新增 `validate-sprint11-safe-completion`，并纳入 `validate-sprint11-real-deployment` 聚合门禁 |

## 安全边界

- 未执行 `wipefs`、`sgdisk`、`mkfs`、`mount`、`fstab` 修改、Rook-Ceph cluster install、OSD claim 或服务器重启。
- 不接受任何数据丢失风险作为 Sprint 11 完成条件。
- 不把 `/dev/sdX` 作为自动化设备身份；后续必须使用 `/dev/disk/by-id`、WWN、序列号、UUID/PARTUUID。
- 遵循上游 Kubernetes/Rook-Ceph 最佳实践：Rook-Ceph OSD 使用 raw unmounted devices，HDD 初期不混入 VM 优先 SSD pool。

## 边界

本批次只证明 Sprint 11 安全完成门禁，不代表 Rook-Ceph 已安装、StorageClass 已上线、PV/PVC 已迁移、VM 持久化存储已生产可用或实际 v1.0.0 发布。
