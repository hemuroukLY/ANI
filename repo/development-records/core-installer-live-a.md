# CORE-INSTALLER-LIVE-A — Core installer live readiness contract

完成日期：2026-06-04
对应 Sprint：Sprint 8（Core-only）
验证结果：`make validate-core-installer-live`、完整 Sprint 8 门禁见 `sprint8-closure-a-contract.md`

## 背景

Sprint 7 已完成 Core installer 三种 profile 的 contract/local validation。Sprint 8 在此基础上新增 installer live readiness profile，固定 baremetal、VM、existing-k8s 三种模式的 evidence 输出入口和校验规则。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/real-k8s-lab/sprint8-core-installer-live.yaml` | 新增 Sprint 8 Core installer live-readiness profile |
| `scripts/validate_core_installer_live.py` | 校验 Core scope、三种安装模式覆盖、profile 路径存在和 production ready 禁止声明 |
| `scripts/validate_core_installer_live_test.py` | 覆盖默认 profile 与 production ready 误声明拒绝 |
| `Makefile` | 新增 `validate-core-installer-live` |

## 边界

- 当前结果是 `contract-local` / `live-ready` 入口，不宣称真实安装已经完成。
- 不把 15 分钟安装、生产化部署或客户现场交付标记为完成。
- 不纳入 Services 组件或前端组件。
