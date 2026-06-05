# CORE-ROOK-CEPH-REBOOT-RESILIENCE-A — Sprint 11 Rook-Ceph reboot resilience result

完成日期：2026-06-05
对应 Sprint：Sprint 11（Core Real Deployment Validation / reboot resilience）
验证结果：`make validate-sprint11-rook-ceph-reboot-resilience`

## 背景

本批次在 Rook-Ceph 正式部署、RBD PVC/Pod smoke 和 KubeVirt VM RBD storage smoke 之后，执行生产化 reboot resilience 验证。范围仍限定 ANI Core，只验证真实底座在逐台节点重启后的恢复能力，不开发 ANI Services、RAG、Console、BOSS、model-service、kb-service、ai、operators 或 frontends。

## 完成内容

1. 前置确认三节点 Ready，`rook-ceph` CephCluster 为 `Ready/HEALTH_OK`，`ceph-rbd-ssd` pool 为 `Ready`，KubeVirt 为 `Deployed`。
2. 按 worker-first、control-plane-last 顺序逐台重启，任一节点恢复前不进入下一台；没有执行并发重启。
3. 两个 worker 节点分别在重启前创建临时 KubeVirt VM + RBD Block PVC 基线；重启后同一 VM/PVC 恢复到 `Running/Bound`，guest 再次看到 RBD block device 并完成写入尝试。
4. control-plane 节点最后重启；恢复后 API `/readyz` 通过，三节点 Ready，control-plane 上 mon/mgr/OSD 均恢复，Ceph 回到 `HEALTH_OK`。
5. control-plane 重启期间保留的 worker VM/PVC 在 API 恢复后仍可观测为 `Running/Bound`。
6. 重启后 Ceph 曾短暂出现 `MON_CLOCK_SKEW`，通过重启系统时间同步服务并对高 offset 节点重新协商 NTP 后清除，最终恢复 `HEALTH_OK`。
7. 验证结束后删除全部临时 VM、VMI、PVC、PV、StorageClass 和 namespace；正式 `ani-rbd-ssd` StorageClass 仍保持 `Retain`、`WaitForFirstConsumer`、非默认。

## 关键文件

| 文件 | 说明 |
|---|---|
| `deploy/real-k8s-lab/sprint11-rook-ceph-reboot-resilience-result.yaml` | 机器可读 reboot resilience 结果 |
| `scripts/validate_sprint11_rook_ceph_reboot_resilience.py` | 校验逐台重启顺序、节点恢复、Ceph 健康、VM/PVC 恢复、清理状态和安全边界 |
| `scripts/validate_sprint11_rook_ceph_reboot_resilience_test.py` | 覆盖未执行重启、并发重启、顺序错误、VM 恢复缺失和清理缺失负例 |

## 安全边界

本批次执行了逐台服务器重启，但没有执行手工 `mount`、`/etc/fstab` 修改、系统盘变更、默认 StorageClass 切换、已有 PVC 迁移或 HDD class 引入。control-plane 恢复等待时间长于 worker，期间未对其他节点执行任何操作；恢复后 API readyz、三节点 Ready、Ceph `HEALTH_OK` 和 VM/PVC 观测均通过。重启后暴露的时间同步偏移已作为生产化问题记录并完成最小修复。当前结果证明逐节点 reboot resilience，不代表备份/恢复、长期 soak、故障注入、升级/回滚演练、业务迁移或完整 production ready 已完成。

## 验收

```bash
make validate-sprint11-rook-ceph-reboot-resilience
make validate-sprint11-real-deployment
make validate-doc-entrypoints
git diff --check
```
