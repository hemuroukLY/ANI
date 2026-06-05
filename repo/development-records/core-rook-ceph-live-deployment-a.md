# CORE-ROOK-CEPH-LIVE-DEPLOYMENT-A — Sprint 11 Rook-Ceph live deployment result

完成日期：2026-06-05
对应 Sprint：Sprint 11（Core Real Deployment Validation / formal live deployment）
验证结果：`make validate-sprint11-rook-ceph-live-deployment-result`

## 背景

本批次在完成只读盘点、存储风险计划、正式部署代码包和人工授权后，执行 Sprint 11 的正式 Rook-Ceph 部署。范围仍限定 ANI Core：只补真实 VM 优先块存储底座，不开发 ANI Services、RAG、Console、BOSS、model-service、kb-service、ai、operators 或 frontends。

## 完成内容

1. 安装 Rook `v1.20.0` CRD/common/operator，并补齐 CSI operator 与 CSI-Addons `v0.14.0` CRD。
2. 应用 Sprint 11 正式部署 manifest：`CephCluster`、`CephBlockPool` 和 `ani-rbd-ssd` StorageClass。
3. Rook-Ceph 集群达到 `Ready/HEALTH_OK`；3 个 mon、1 个 mgr、5 个 SSD OSD 运行。
4. `ceph-rbd-ssd` pool 达到 `Ready`；`ani-rbd-ssd` StorageClass 使用 RBD CSI，`Retain`、`WaitForFirstConsumer`，且不设为默认。
5. 受控 RBD smoke test 通过：临时 PVC 绑定、Pod 挂载、写读 marker 成功；临时 StorageClass/PVC/Pod/PV 已删除。
6. ANI3 盘点结论校准：ANI3 只有 2 块 SSD 候选盘；另 1 块 3.6T 盘为 HDD，已排除出 VM 优先 SSD pool。

## 关键文件

| 文件 | 说明 |
|---|---|
| `deploy/real-k8s-lab/sprint11-rook-ceph-live-deployment-result.yaml` | 机器可读 live 部署结果 |
| `scripts/validate_sprint11_rook_ceph_live_deployment_result.py` | 校验 live 结果、OSD 数量、StorageClass 策略、smoke test 清理和安全边界 |
| `scripts/validate_sprint11_rook_ceph_live_deployment_result_test.py` | 覆盖 live flag、默认 StorageClass、OSD 数量和 smoke 清理负例 |
| `deploy/real-k8s-lab/sprint11-rook-ceph-formal-deployment.yaml` | 正式部署 manifest 源 |

## 安全边界

本批次未执行手工 `wipefs`、`sgdisk`、`mkfs`、`mount`、`/etc/fstab` 修改、系统盘变更或服务器重启。Rook-Ceph 按审批后的 manifest 对已验证的 SSD raw devices 执行 OSD prepare 和 OSD 认领。当前结果不是实际 v1.0.0 发布，也不是完整 production ready；备份/恢复演练、容量告警、故障注入、长期 soak、升级/回滚演练、默认 StorageClass 切换、已有 PVC 迁移和 HDD class 引入仍需后续单独审批与验证。

## 验收

```bash
make validate-sprint11-rook-ceph-live-deployment-result
make validate-sprint11-real-deployment
make validate-doc-entrypoints
git diff --check
```
