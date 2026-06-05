# CORE-CLI-RELEASE-METADATA-A — ANI Core CLI release metadata

完成日期：2026-06-04
对应 Sprint：Sprint 10（Core-only）
验证结果：`make validate-core-cli`、`make build-cli`；完整 Sprint 10 门禁见 `sprint10-closure-a-contract.md`

## 背景

Sprint 9 已提供 `ani --version` 文本输出。Sprint 10 release-prep 需要机器可读版本信息，便于 release evidence、自动化检查和现场排障读取。

## 关键变更

| 文件 | 说明 |
|---|---|
| `cli/ani/main.go` | 新增 `--version-format text/json`；JSON 输出包含 `name`、`scope`、`version`、`build_time` |
| `cli/ani/main_test.go` | 新增 JSON version metadata 测试 |

## 边界

- 只扩展 Core CLI 本地元数据输出，不新增 Services 命令。
- CLI 正式发布、签名和渠道分发仍不在本批次内完成。
