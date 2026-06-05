# CORE-REAL-DEPLOY-DOC-CONSISTENCY-A — Sprint 11 Core documentation consistency gate

完成日期：2026-06-04
对应 Sprint：Sprint 11（Core Real Deployment Validation）
验证结果：`make validate-sprint11-core-doc-consistency`

## 背景

Sprint 11 进入真实服务器验证阶段后，入口文档必须同时让人类和 AI agent 清楚当前只完成了只读验证和风险建模，不能误读为已经执行存储部署或生产化交付。

## 关键变更

| 文件 | 说明 |
|---|---|
| `scripts/validate_sprint11_core_doc_consistency.py` | 校验 `ANI-DOCS-INDEX.md`、`ANI-06-开发计划.md`、`repo/CURRENT-SPRINT.md`、`repo/README.md`、Makefile targets 和 records 索引 |
| `scripts/validate_sprint11_core_doc_consistency_test.py` | 覆盖完整文档工作区和缺失 Sprint 11 marker 的失败场景 |
| `Makefile` | 新增 `validate-sprint11-core-doc-consistency` |
| `ANI-DOCS-INDEX.md`、`ANI-06-开发计划.md`、`repo/CURRENT-SPRINT.md`、`repo/README.md`、`repo/development-records/README.md` | 同步 Sprint 11 当前状态 |

## 边界

本批次完成时，真实服务器只读验证已完成，尚未执行磁盘挂载、格式化、Rook-Ceph 安装或 OSD 认领，因此当时文档不得把 Sprint 11 kickoff/read-only validation 写成 production ready、actual release 或真实存储上线。后续 `CORE-ROOK-CEPH-LIVE-DEPLOYMENT-A` 已完成正式 live 部署；当前状态以 `deploy/real-k8s-lab/sprint11-rook-ceph-live-deployment-result.yaml` 和 `repo/CURRENT-SPRINT.md` 为准。
