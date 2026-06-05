# CORE-RC-GATE-A — Sprint 9 Core RC readiness profile

完成日期：2026-06-04
对应 Sprint：Sprint 9（Core-only）
验证结果：`make validate-sprint9-core-rc`；完整 Sprint 9 门禁见 `sprint9-closure-a-contract.md`

## 背景

Sprint 9 进入 Core RC readiness 加固阶段。本批次建立 Sprint 9 聚合 gate profile，用固定命令清单串联 Core API 兼容性、架构边界、文档入口、SDK、Sprint 8 release gates、release evidence、offline package lock、CLI version 和 Sprint 9 文档一致性。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/release/sprint9-core-rc.yaml` | 新增 Sprint 9 Core RC readiness gate profile |
| `scripts/validate_sprint9_core_rc.py` | 校验 profile scope、版本、RC readiness 标记、必需 gates、重复 ID 和 Services 越界 |
| `scripts/validate_sprint9_core_rc_test.py` | 覆盖默认 profile 与 Services gate 拒绝 |
| `Makefile` | 新增 `validate-sprint9-core-rc` 和 `validate-sprint9-rc` |

## 边界

- 这是 RC readiness gate，不是实际 `v1.0.0-rc.*` cut。
- `release_candidate` 和 `production_release` 均保持 `false`。
- 不纳入 RAG、Console、BOSS、model-service、kb-service、ai、operators 或 frontends。
- 不新增 `M1-REAL-LAB-*` guard。
