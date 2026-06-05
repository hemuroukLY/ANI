# CORE-OFFLINE-PACK-A — Core offline package lock

完成日期：2026-06-04
对应 Sprint：Sprint 8（Core-only）
验证结果：`make validate-core-offline-pack`、完整 Sprint 8 门禁见 `sprint8-closure-a-contract.md`

## 背景

Sprint 7 已完成 Core offline package manifest。Sprint 8 增加 offline package lock，用于固定离线包 artifact 名称、checksum 占位、source manifest 和 verification commands，防止离线交付清单与 Core-only 边界漂移。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/offline/core-package-lock.yaml` | 新增 Core offline package lock |
| `scripts/validate_core_offline_pack.py` | 校验 source manifest、artifact sha256、verification commands 和 Services 业务越界 |
| `scripts/validate_core_offline_pack_test.py` | 覆盖默认 lock 与 checksum 缺失拒绝 |
| `Makefile` | 新增 `validate-core-offline-pack` |

## 边界

- 当前结果是离线包 lock/verification contract，不代表真实离线包已制作、签名或客户现场可交付。
- artifact checksum 为本地 contract 占位，真实发布时必须替换为实际构建产物 checksum。
- 不纳入 Services 业务镜像。
