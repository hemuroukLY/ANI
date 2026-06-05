# CORE-ROOK-CEPH-VM-STORAGE-SMOKE-A — Sprint 11 KubeVirt VM RBD storage smoke result

完成日期：2026-06-05
对应 Sprint：Sprint 11（Core Real Deployment Validation / VM storage smoke）
验证结果：`make validate-sprint11-rook-ceph-vm-storage-smoke`

## 背景

本批次在 Rook-Ceph 正式部署完成后，补充验证 ANI Core 的 VM 优先存储路径：KubeVirt VM 能使用 Rook-Ceph RBD CSI 动态创建的 Block PVC。范围仍限定 ANI Core，不开发 ANI Services、RAG、Console、BOSS、model-service、kb-service、ai、operators 或 frontends。

## 完成内容

1. 创建临时 namespace、临时 `Delete` 策略 StorageClass `ani-rbd-ssd-vm-smoke-delete`、1Gi Block PVC 和 KubeVirt VM。
2. PVC 通过 `rook-ceph.rbd.csi.ceph.com` 动态创建 RBD PV 并达到 `Bound`。
3. KubeVirt VMI 达到 `Running/Ready`，virt-launcher pod 运行，VM 调度到真实物理节点。
4. VM guest probe 看到 RBD block device `/dev/vdb`，并对临时块设备完成写入尝试。
5. 验证后删除临时 VM、VMI、PVC、PV、StorageClass 和 namespace；CephCluster 仍为 `Ready/HEALTH_OK`，`ceph-rbd-ssd` pool 仍为 `Ready`。

## 关键文件

| 文件 | 说明 |
|---|---|
| `deploy/real-k8s-lab/sprint11-rook-ceph-vm-storage-smoke-result.yaml` | 机器可读 VM storage smoke 结果 |
| `scripts/validate_sprint11_rook_ceph_vm_storage_smoke.py` | 校验临时 StorageClass、Block PVC、VM guest probe、清理状态和安全边界 |
| `scripts/validate_sprint11_rook_ceph_vm_storage_smoke_test.py` | 覆盖 live flag、默认 StorageClass、guest probe 和 PV 清理负例 |
| `deploy/real-k8s-lab/sprint11-rook-ceph-live-deployment-result.yaml` | Rook-Ceph 正式部署结果源 |

## 安全边界

本批次没有执行手工 `mount`、`/etc/fstab` 修改、系统盘变更、服务器重启、默认 StorageClass 切换或已有 PVC 迁移。临时 PVC 使用 `Delete` 回收策略，仅用于 smoke test；正式 `ani-rbd-ssd` StorageClass 仍保持 `Retain`、`WaitForFirstConsumer`、非默认。当前结果证明 VM attach 与 guest block device 可见性，不代表备份/恢复、长期 soak、故障注入、升级/回滚演练、业务迁移或完整 production ready 已完成。

## 验收

```bash
make validate-sprint11-rook-ceph-vm-storage-smoke
make validate-sprint11-real-deployment
make validate-doc-entrypoints
git diff --check
```
