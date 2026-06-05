# CORE-REGRESSION-A — Sprint 7 Core regression profile

完成日期：2026-06-04
对应 Sprint：Sprint 7（Core-only）
验证结果：`make validate-sprint7-core-regression`、完整 Sprint 7 门禁见 `sprint7-closure-a-contract.md`

## 实现了什么

新增 Sprint 7 Core regression profile，将新增的 installer/offline/CLI contract gate 与历史 real path / SDK smoke gate 组合成固定回归入口。

当前 profile 只定义 contract regression composition，不执行真实 live 模式，不新增 `M1-REAL-LAB-*` guard。

## 关键文件改动

| 文件 | 修改说明 |
|---|---|
| `deploy/real-k8s-lab/sprint7-core-regression.yaml` | 新增 Sprint 7 Core regression target profile |
| `scripts/validate_sprint7_core_regression.py` | 新增 regression profile 校验器 |
| `scripts/validate_sprint7_core_regression_test.py` | 覆盖必需 targets 和 guard 冻结约束 |
| `Makefile` | 新增 `validate-sprint7-core-regression` |

## 边界

- 不新增 REAL-K8S-LAB guard。
- 不把 contract target 写成 live evidence。
- 不把 Services/RAG/Console 纳入本仓库回归目标。
