# CORE-DOC-CONSISTENCY-A — Core documentation consistency gate

完成日期：2026-06-04
对应 Sprint：Sprint 8（Core-only）
验证结果：`make validate-core-doc-consistency`、完整 Sprint 8 门禁见 `sprint8-closure-a-contract.md`

## 背景

Sprint 8 要求人类和 AI agent 都能从入口文档准确判断当前阶段、代码范围、验收命令和历史边界。本批次新增文档一致性 validator，防止当前入口停留在 Sprint 7 completed / Sprint 8 prep 或遗漏 Sprint 8 Makefile targets 与批次记录。

## 关键变更

| 文件 | 说明 |
|---|---|
| `scripts/validate_core_doc_consistency.py` | 校验三份入口文档 Sprint 8 completed marker、Makefile targets 和 development records 索引 |
| `scripts/validate_core_doc_consistency_test.py` | 覆盖默认文档一致性与 stale Sprint 7 current marker 拒绝 |
| `Makefile` | 新增 `validate-core-doc-consistency` |

## 边界

- 该 validator 只校验当前阶段和 Sprint 8 Core gates 的一致性。
- 历史归档文档允许保留历史语境，不要求反向改写。
- 不修改 `CLAUDE.md` 的稳定规则职责。
