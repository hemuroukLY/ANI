# CORE-STORAGE-DISK-RISK-A — Sprint 11 Core storage disk risk plan

完成日期：2026-06-04
对应 Sprint：Sprint 11（Core Real Deployment Validation）
验证结果：`make validate-sprint11-storage-disk-plan`

## 背景

三台物理服务器的数据盘尚未挂载。为了后续 VM-first 块存储和 Rook-Ceph 部署不误碰系统盘，本批次先用只读盘点结果建立稳定设备清单和风险策略。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/real-k8s-lab/sprint11-storage-disk-plan.yaml` | 记录三台服务器系统盘、数据盘、稳定 `/dev/disk/by-id` 路径、SSD/HDD 分类、Rook-Ceph 初始策略和禁止操作 |
| `scripts/validate_sprint11_storage_disk_plan.py` | 校验 Core scope、Sprint 11 版本、只读状态、禁止 `/dev/sdX` 自动化、Rook-Ceph raw unmounted device 策略和风险控制 |
| `scripts/validate_sprint11_storage_disk_plan_test.py` | 覆盖默认计划、破坏性模式、`/dev/sdX` 身份和 HDD 混入初始 VM SSD pool 的拒绝场景 |
| `Makefile` | 新增 `validate-sprint11-storage-disk-plan` |

## 只读盘点结论

- ANI1 系统盘观测为 `sdb`，数据 SSD 观测为 `sda`。
- ANI2 系统盘观测为 `sdc`，数据 SSD 观测为 `sda`、`sdb`。
- ANI3 系统盘观测为 `sdd`，数据 SSD 观测为 `sda`、`sdb`，HDD 观测为 `sdc`。
- 当前策略不把这些 `/dev/sdX` 作为自动化依据，只作为人工阅读时的观测值；真实选择必须使用 `/dev/disk/by-id`。

## 边界

真实服务器只读验证已完成；未执行磁盘挂载、格式化、Rook-Ceph 安装或 OSD 认领。后续如需继续部署，必须先人工审批具体设备 ID、预期影响和回滚方案。
