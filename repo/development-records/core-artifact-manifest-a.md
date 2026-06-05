# CORE-ARTIFACT-MANIFEST-A — Core artifact manifest

完成日期：2026-06-04
对应 Sprint：Sprint 10（Core-only）
验证结果：`make validate-core-artifact-manifest`；完整 Sprint 10 门禁见 `sprint10-closure-a-contract.md`

## 背景

Sprint 10 进入 Core 发布前收敛阶段。本批次为核心 release-prep artifact 建立机器可读 SHA256 清单，让人类和 AI agent 可以复核 Core OpenAPI、Core SDK metadata、CLI source、offline lock 和 release evidence 是否与当前工作区一致。

## 关键变更

| 文件 | 说明 |
|---|---|
| `deploy/release/core-artifacts.yaml` | 新增 Sprint 10 Core artifact manifest |
| `scripts/validate_core_artifact_manifest.py` | 校验 artifact scope、必需条目、SHA256、验证命令和 Services 越界 |
| `scripts/validate_core_artifact_manifest_test.py` | 覆盖默认 manifest 与 checksum mismatch 拒绝 |
| `Makefile` | 新增 `validate-core-artifact-manifest` |

## 边界

- 这是 release-prep artifact 清单，不是真实 release tarball、镜像签名或客户现场交付物。
- 不纳入 Services、RAG、Console、BOSS、ai、operators 或 frontends。
