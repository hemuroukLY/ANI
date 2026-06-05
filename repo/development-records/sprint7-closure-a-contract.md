# SPRINT7-CLOSURE-A — Sprint 7 Core-only closure contract

完成日期：2026-06-04
对应 Sprint：Sprint 7（Core-only）
验证结果：`make test`、`make validate-architecture`、`make validate-core-api-compatibility`、`make validate-doc-entrypoints`、`make validate-sdk-beta`、`make validate-sdk-mock-smoke`、`make validate-core-installer`、`make validate-core-offline`、`make validate-core-cli`、`make validate-sprint7-core-regression`、`make build-cli`、`git diff --check` 均为 EXIT 0

## 实现了什么

完成 Sprint 7 Core-only 代码开发闭环：

1. `CORE-INSTALLER-A`：Core installer 三种 profile 与 contract validator。
2. `CORE-OFFLINE-A`：Core offline package manifest 与 contract validator。
3. `CORE-CLI-A`：`ani` Core CLI 最小只读资源访问和 Services 命令拒绝。
4. `CORE-REGRESSION-A`：Sprint 7 Core regression profile，串联新增 Core gates 和历史回归门禁。

本 Sprint 不开发 RAG、Console、model-service、kb-service、ai、operators 或 frontends，不新增 `M1-REAL-LAB-*` guard，不把 installer/offline/CLI 标记为 production ready。

## 关键文件改动

| 文件 | 修改说明 |
|---|---|
| `installer/ani-installer/profiles/*.yaml` | Core installer profile |
| `deploy/offline/core-package.yaml` | Core offline package manifest |
| `cli/ani/*` | Core CLI module |
| `deploy/real-k8s-lab/sprint7-core-regression.yaml` | Sprint 7 regression profile |
| `scripts/validate_core_installer.py` / `scripts/validate_core_offline.py` / `scripts/validate_sprint7_core_regression.py` | 新增 Sprint 7 validators |
| `Makefile` / `go.work` | 接入 CLI 构建和 Sprint 7 validation targets |

## 边界

- 当前完成的是 contract/local validation 和 CLI 最小行为。
- 真实安装、离线包实际制作/签名、CLI 发布和生产化部署仍需后续 live/offline evidence。
- Services 业务能力由外部团队负责，本仓库只维护 ANI Core。
