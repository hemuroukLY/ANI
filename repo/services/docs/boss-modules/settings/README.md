# BOSS 交付与安装域（settings）

12 个模块，权威源以 [`ANI-main/ANI-07-部署工程设计.md`](../../../ANI-main/ANI-07-部署工程设计.md) 为主；OpenAPI 安装工作流均为 **TODO-YAML: N/A**（除品牌 GET 与 k8s-clusters 只读参考）。

## 模块索引

| 模块 | 详文 | HTML ID | OpenAPI |
|---|---|---|---|
| ani-installer | [ani-installer.md](ani-installer.md) | `boss.settings.installer` | N/A |
| 裸机部署 | [baremetal-deploy.md](baremetal-deploy.md) | `boss.settings.baremetal` | N/A |
| VM 部署 | [vm-deploy.md](vm-deploy.md) | `boss.settings.vm_deploy` | N/A |
| 已有 K8s 接入 | [k8s-attach.md](k8s-attach.md) | `boss.settings.k8s_attach` | attach N/A；`k8s-clusters` 只读参考 |
| 硬件检测 | [hardware-check.md](hardware-check.md) | `boss.settings.hw_check` | N/A |
| GPU 驱动安装 | [gpu-driver-install.md](gpu-driver-install.md) | `boss.settings.driver` | N/A |
| 内部 CA | [internal-ca.md](internal-ca.md) | `boss.settings.ca` | N/A |
| 纯 IP HTTPS | [ip-https.md](ip-https.md) | `boss.settings.ip_https` | N/A |
| 离线安装包 | [offline-package.md](offline-package.md) | `boss.settings.offline` | N/A |
| 交付验收手册 | [acceptance-manual.md](acceptance-manual.md) | `boss.settings.acceptance` | N/A |
| 品牌配置 | [brand-config.md](brand-config.md) | `boss.settings.brand` | GET `/branding` 已声明；PATCH **TODO-YAML P2** |
| 升级与备份 | [upgrade-backup.md](upgrade-backup.md) | `boss.settings.upgrade` | 平台 patch **TODO-YAML P2** |

## 同步与门禁

```bash
python3 scripts/sync-boss-settings-domain.py
python3 scripts/validate-boss-settings-gate.py
```

Phase 0 GAP 摘要：[`../governance/boss-phase0-gap-settings.md`](../governance/boss-phase0-gap-settings.md)
