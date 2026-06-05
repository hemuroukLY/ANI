# CORE-INSTALLER-A — Core installer profile contract

完成日期：2026-06-04
对应 Sprint：Sprint 7（Core-only）
验证结果：`make validate-core-installer`、`make build-cli`、完整 Sprint 7 门禁见 `sprint7-closure-a-contract.md`

## 实现了什么

建立 ani-installer 的 Core-only 最小可验证闭环：新增 baremetal、VM、existing-k8s 三种安装 profile，并用 `scripts/validate_core_installer.py` 固定校验范围、组件清单、preflight 命令和安装计划阶段。

当前只证明 installer contract/local validation，不代表 15 分钟安装、真实裸机安装或 production ready。

## 关键文件改动

| 文件 | 修改说明 |
|---|---|
| `installer/ani-installer/profiles/*.yaml` | 新增三种 Core installer profile |
| `scripts/validate_core_installer.py` | 新增 Core installer profile 校验器 |
| `scripts/validate_core_installer_test.py` | 覆盖模式完整性、Services 组件拒绝、preflight 阶段要求 |
| `Makefile` | 新增 `validate-core-installer` |

## 边界

- 不包含 Services/RAG/Console/frontends。
- 不安装或声明真实生产集群。
- 不新增 REAL-K8S-LAB guard。
