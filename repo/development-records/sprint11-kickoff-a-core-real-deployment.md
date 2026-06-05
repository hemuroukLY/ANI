# SPRINT11-KICKOFF-A — Sprint 11 Core real deployment validation kickoff

完成日期：2026-06-04
对应 Sprint：Sprint 11（Core Real Deployment Validation）
验证结果：真实服务器只读验证已完成；未执行磁盘挂载、格式化、Rook-Ceph 安装或 OSD 认领。代码门禁见 `CORE-STORAGE-DISK-RISK-A`、`CORE-REAL-DEPLOY-A` 和 `CORE-REAL-DEPLOY-DOC-CONSISTENCY-A`。

## 背景

Sprint 6-10 已形成 Core contract/local/release-prep 成果，但这些成果还需要放到三台物理服务器和现有 Kubernetes/KubeVirt 环境上做真实部署前校验。Sprint 11 的第一步不是直接安装存储或改磁盘，而是建立只读盘点、风险建模和可复跑门禁。

## 边界

- 只推进 ANI Core。
- 不开发 RAG、Console、BOSS、model-service、kb-service、ai、operators 或 frontends。
- 不新增 `M1-REAL-LAB-*` guard。
- 不执行 `wipefs`、`sgdisk`、`mkfs`、`mount`、`fstab` 修改、Rook-Ceph cluster install 或 OSD claim。

## 结果

入口文档已切换到 Sprint 11 / Core Real Deployment Validation 已启动状态。当前可继续执行的安全动作是只读验证、profile 校验和人工审批前的存储部署方案审查。
