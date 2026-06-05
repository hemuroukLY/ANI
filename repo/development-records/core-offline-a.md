# CORE-OFFLINE-A — Core offline package manifest contract

完成日期：2026-06-04
对应 Sprint：Sprint 7（Core-only）
验证结果：`make validate-core-offline`、完整 Sprint 7 门禁见 `sprint7-closure-a-contract.md`

## 实现了什么

建立 Core 离线包 manifest 的最小可验证闭环：新增 `deploy/offline/core-package.yaml`，固定 Core 镜像、Helm chart 和脚本清单，并用 `scripts/validate_core_offline.py` 校验镜像 digest pin、Core-only scope 和路径存在性。

当前只证明离线包 manifest contract，不代表镜像已经实际拉取、打包、签名或客户现场可交付。

## 关键文件改动

| 文件 | 修改说明 |
|---|---|
| `deploy/offline/core-package.yaml` | 新增 Core offline package manifest |
| `scripts/validate_core_offline.py` | 新增离线包 manifest 校验器 |
| `scripts/validate_core_offline_test.py` | 覆盖 digest pin 和 Services 镜像拒绝 |
| `Makefile` | 新增 `validate-core-offline` |

## 边界

- 离线包只包含 Core 镜像和 Core chart/script 清单。
- 不纳入 model-service、kb-service、RAG、operator 或 frontend 镜像。
- 不宣称 offline package production ready。
