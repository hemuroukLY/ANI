# SPRINT8-CLOSURE-A — Sprint 8 Core-only closure contract

完成日期：2026-06-04
对应 Sprint：Sprint 8（Core-only）
验证结果：`make test`、`make validate-architecture`、`make validate-core-api-compatibility`、`make validate-doc-entrypoints`、`make validate-sdk-beta`、`make validate-sdk-mock-smoke`、`make validate-core-installer`、`make validate-core-offline`、`make validate-core-cli`、`make validate-sprint7-core-regression`、`make validate-core-release-hardening`、`make validate-core-installer-live`、`make validate-core-offline-pack`、`make validate-core-doc-consistency`、`make validate-sprint8-core-release`、`make build-cli`、`git diff --check` 均为 EXIT 0

## 背景

完成 Sprint 8 Core-only 收尾/发布前加固闭环：

1. `CORE-HARDEN-A`：Core release hardening profile。
2. `CORE-INSTALLER-LIVE-A`：Core installer live-readiness profile。
3. `CORE-OFFLINE-PACK-A`：Core offline package lock。
4. `CORE-CLI-B`：Core CLI 主要只读资源扩展。
5. `CORE-DOC-CONSISTENCY-A`：入口文档和 Sprint 8 gate 一致性校验。

## 边界

本 Sprint 不开发 RAG、Console、BOSS、model-service、kb-service、ai、operators 或 frontends，不新增 `M1-REAL-LAB-*` guard，不把 installer/offline/CLI 标记为 production ready。

## 后续

下一步进入 Sprint 9 RC 加固，只修 ANI Core bug、补真实 release evidence、替换真实离线包 checksum，并保持 Core API v1 兼容。
