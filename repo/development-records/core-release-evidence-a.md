# CORE-RELEASE-EVIDENCE-A — Core release evidence manifest

完成日期：2026-06-04
对应 Sprint：Sprint 9（Core-only）
验证结果：`make validate-core-release-evidence`；完整 Sprint 9 门禁见 `sprint9-closure-a-contract.md`

## 背景

Sprint 9 需要让 RC readiness 的证据对人类和 AI agent 都可读。本批次新增 release evidence manifest，记录每项 Core release readiness 证据对应的可复跑命令和 artifact 引用。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/release/core-release-evidence.yaml` | 新增 Core release evidence manifest |
| `scripts/validate_core_release_evidence.py` | 校验 scope、Sprint 9 版本、RC readiness 标记、必需 evidence、Services 越界和敏感信息关键词 |
| `scripts/validate_core_release_evidence_test.py` | 覆盖默认 manifest 与敏感 artifact 拒绝 |
| `Makefile` | 新增 `validate-core-release-evidence` |

## 边界

- evidence manifest 只记录命令和非敏感 artifact 路径，不记录 token、password、credential、private key 或客户凭据。
- 该批次证明 RC readiness 证据可追溯，不代表实际 RC 发布或生产审批完成。
- 不触碰 Services 冻结目录。
