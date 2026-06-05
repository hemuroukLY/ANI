# CORE-VERSION-POLICY-A — Core version policy

完成日期：2026-06-04
对应 Sprint：Sprint 10（Core-only）
验证结果：`make validate-core-version-policy`；完整 Sprint 10 门禁见 `sprint10-closure-a-contract.md`

## 背景

`CLAUDE.md` 明确当前不得标记实际 `v1.0.0` 或 RC。Sprint 10 因此只能完成 Core release-prep readiness，不能伪造正式发布。本批次新增版本策略 manifest 和 validator，确保 release YAML 与入口文档都保留这个边界。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/release/core-version-policy.yaml` | 新增 Sprint 10 Core version policy manifest |
| `scripts/validate_core_version_policy.py` | 校验 `actual_release`、`release_candidate`、`production_release`、`rc_cut` 均为 false，并要求入口文档声明不是实际 v1.0.0 发布 |
| `scripts/validate_core_version_policy_test.py` | 覆盖默认策略和 release flag 越界拒绝 |
| `Makefile` | 新增 `validate-core-version-policy` |

## 边界

- Sprint 10 Core-only 代码开发已完成不是实际 v1.0.0 发布。
- 不创建 tag、不切 RC、不声明 production ready。
