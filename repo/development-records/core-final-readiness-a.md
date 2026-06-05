# CORE-FINAL-READINESS-A — Sprint 10 Core final readiness profile

完成日期：2026-06-04
对应 Sprint：Sprint 10（Core-only）
验证结果：`make validate-sprint10-final-readiness`；完整 Sprint 10 门禁见 `sprint10-closure-a-contract.md`

## 背景

Sprint 10 需要把 Sprint 9 RC readiness、Core artifact manifest、version policy、CLI build、SDK/API/doc gates 串成统一 release-prep 入口。本批次新增 final readiness profile 作为 Sprint 10 的聚合配置。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/release/sprint10-core-final-readiness.yaml` | 新增 Sprint 10 Core release-prep readiness profile |
| `scripts/validate_sprint10_final_readiness.py` | 校验 profile scope、target release、false release flags、必需 gates、重复 ID 和 Services 越界 |
| `scripts/validate_sprint10_final_readiness_test.py` | 覆盖默认 profile 和 actual_release 越界拒绝 |
| `Makefile` | 新增 `validate-sprint10-final-readiness` 和 `validate-sprint10-release-prep` |

## 边界

- `release_preparation_complete=true` 只表示发布前准备门禁完成。
- `actual_release=false`、`release_candidate=false`、`production_release=false`，不是实际 v1.0.0 发布。
- 不新增 `M1-REAL-LAB-*` guard。
