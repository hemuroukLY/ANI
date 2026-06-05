# CORE-FINAL-DOC-CONSISTENCY-A — Sprint 10 Core documentation consistency gate

完成日期：2026-06-04
对应 Sprint：Sprint 10（Core-only）
验证结果：`make validate-sprint10-core-doc-consistency`；完整 Sprint 10 门禁见 `sprint10-closure-a-contract.md`

## 背景

Sprint 10 新增 artifact manifest、version policy、final readiness profile 和 CLI release metadata 后，需要确保入口文档、Makefile target 和 development records 对同一状态给出一致结论。

## 关键变更

| 文件 | 说明 |
|---|---|
| `scripts/validate_sprint10_core_doc_consistency.py` | 校验 `ANI-DOCS-INDEX.md`、`ANI-06-开发计划.md`、`repo/CURRENT-SPRINT.md`、Makefile 和 records 索引的 Sprint 10 一致性 |
| `scripts/validate_sprint10_core_doc_consistency_test.py` | 覆盖默认文档一致性和 stale Sprint 9 current marker 拒绝 |
| `Makefile` | 新增 `validate-sprint10-core-doc-consistency` |

## 边界

- 本 gate 不替代正式发布审批或真实 release artifact 签名。
- 必须明确 Sprint 10 completion 不是实际 v1.0.0 发布。
