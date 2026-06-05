# CORE-HARDEN-A — Core release hardening contract

完成日期：2026-06-04
对应 Sprint：Sprint 8（Core-only）
验证结果：`make validate-core-release-hardening`、完整 Sprint 8 门禁见 `sprint8-closure-a-contract.md`

## 背景

Sprint 8 进入 Core 收尾/发布前加固阶段。本批次建立 Core release hardening profile，用固定 gate 覆盖 Core API 兼容性、ports/adapters 架构边界、文档入口、SDK Beta、SDK Mock smoke、Core CLI、Core installer 和 Core offline manifest。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/release/core-hardening.yaml` | 新增 Sprint 8 Core release hardening gate profile |
| `scripts/validate_core_release_hardening.py` | 校验 hardening profile scope、必需 gates、重复 ID 和 Services 业务越界 |
| `scripts/validate_core_release_hardening_test.py` | 覆盖默认 profile 与 Services gate 拒绝 |
| `Makefile` | 新增 `validate-core-release-hardening` |

## 边界

- 只覆盖 ANI Core 发布前 contract/local gates。
- 不纳入 RAG、Console、BOSS、model-service、kb-service、ai、operators 或 frontends。
- 不替代后续真实安全扫描、镜像签名或生产发布审批。
