# CORE-OFFLINE-CHECKSUM-A — Core offline checksum contract

完成日期：2026-06-04
对应 Sprint：Sprint 9（Core-only）
验证结果：`make validate-core-offline-pack`；完整 Sprint 9 门禁见 `sprint9-closure-a-contract.md`

## 背景

Sprint 8 的 offline package lock 已固定 artifact 与 verification commands，但 checksum 仍是 contract 占位值。Sprint 9 将它收敛为可复算的本地 source manifest checksum，并让校验器拒绝占位或不一致 checksum。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/offline/core-package-lock.yaml` | 新增 `source_manifest_sha256`，并将 `artifact.sha256` 对齐为 `deploy/offline/core-package.yaml` 当前内容 SHA256 |
| `scripts/validate_core_offline_pack.py` | 校验 source manifest checksum、拒绝占位 checksum，并要求 Sprint 9 contract package checksum 与 source manifest checksum 一致 |
| `scripts/validate_core_offline_pack_test.py` | 覆盖默认 lock checksum 和 mismatch 拒绝 |

## 边界

- 当前 checksum 是 Sprint 9 contract package 的 source manifest checksum，不是已签名真实离线 tarball 的交付 checksum。
- 真实离线包制作、签名、客户现场验证和供应链扫描仍归后续 release evidence/发布流程。
