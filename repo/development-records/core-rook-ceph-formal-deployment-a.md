# CORE-ROOK-CEPH-FORMAL-DEPLOYMENT-A — Sprint 11 Rook-Ceph formal deployment code

完成日期：2026-06-05
对应 Sprint：Sprint 11（Core Real Deployment Validation / formal deployment code）
验证结果：`make validate-sprint11-rook-ceph-formal-deployment`

## 背景

Sprint 11 不能只停留在环境验证和只读盘点，需要进入正式部署所需的 Core 代码与部署资产开发。但三台物理服务器的数据盘仍属于高风险存储操作边界，任何 wipe、format、mount、Rook-Ceph install 或 OSD claim 都必须在人工审批后执行。

## 完成内容

新增 Rook-Ceph 正式部署代码包，面向 VM 优先块存储：

1. `CephCluster`：只使用三台物理服务器上的 SSD OSD 候选盘，全部以 `/dev/disk/by-id` 持久设备身份引用。
2. `CephBlockPool`：定义 `ceph-rbd-ssd`，`failureDomain=host`，副本数为 3，并要求安全副本大小。
3. `StorageClass`：定义 `ani-rbd-ssd`，使用 RBD CSI，`reclaimPolicy=Retain`，`volumeBindingMode=WaitForFirstConsumer`，且不自动设为默认 StorageClass。
4. 审批门禁：正式部署包要求人工确认设备 ID、无挂载/无文件系统/无分区/LVM、回滚窗口和 Rook operator 版本后才允许 apply。

## 关键文件

| 文件 | 说明 |
|---|---|
| `deploy/real-k8s-lab/sprint11-rook-ceph-formal-deployment.yaml` | 新增 Rook-Ceph 正式部署代码包 |
| `scripts/validate_sprint11_rook_ceph_formal_deployment.py` | 校验部署包只使用稳定设备 ID、排除 HDD、禁止默认 StorageClass 切换和无审批 live mutation |
| `scripts/validate_sprint11_rook_ceph_formal_deployment_test.py` | 覆盖默认包、live claim、`/dev/sdX`、HDD/计划外设备和默认 StorageClass 风险 |
| `Makefile` | 新增 `validate-sprint11-rook-ceph-formal-deployment`，并接入 Sprint 11 聚合门禁 |

## 安全边界

本批次完成时的范围是正式部署代码资产，不是 live Rook-Ceph 安装；当时仍未执行 `kubectl apply`、Rook-Ceph operator install、CephCluster 创建、OSD 认领、StorageClass 默认切换、磁盘 wipe/format/mount 或服务器重启。后续 `CORE-ROOK-CEPH-LIVE-DEPLOYMENT-A` 已在同一 Sprint 内完成正式 live 部署，当前状态以 `deploy/real-k8s-lab/sprint11-rook-ceph-live-deployment-result.yaml` 为准。

## 验收

```bash
make validate-sprint11-rook-ceph-formal-deployment
make validate-sprint11-real-deployment
make validate-doc-entrypoints
make validate-architecture
make test
git diff --check
```
