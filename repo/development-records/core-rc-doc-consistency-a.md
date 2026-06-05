# CORE-RC-DOC-CONSISTENCY-A — Sprint 9 Core documentation consistency gate

完成日期：2026-06-04
对应 Sprint：Sprint 9（Core-only）
验证结果：`make validate-sprint9-core-doc-consistency`；完整 Sprint 9 门禁见 `sprint9-closure-a-contract.md`

## 背景

Sprint 9 新增 RC readiness profile、release evidence、offline checksum 和 CLI version 后，需要让入口文档、Makefile target 和 development records 保持一致，避免人类和 AI agent 读取到不同阶段判断。

## 关键变更

| 文件 | 说明 |
|---|---|
| `scripts/validate_sprint9_core_doc_consistency.py` | 校验 `ANI-DOCS-INDEX.md`、`ANI-06-开发计划.md`、`repo/CURRENT-SPRINT.md`、Makefile 和 records 索引的 Sprint 9 一致性 |
| `scripts/validate_sprint9_core_doc_consistency_test.py` | 覆盖默认文档一致性和 stale Sprint 8 current marker 拒绝 |
| `Makefile` | 新增 `validate-sprint9-core-doc-consistency` |

## 边界

- 本 gate 只校验文档和入口一致性，不替代真实 release、安装、离线包签名或客户现场验收。
- 不更新 `CLAUDE.md` 的动态进度，遵守入口文档职责边界。
