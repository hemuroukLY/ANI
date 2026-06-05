# CORE-CLI-VERSION-A — ANI Core CLI version output

完成日期：2026-06-04
对应 Sprint：Sprint 9（Core-only）
验证结果：`make validate-core-cli`、`make build-cli`；完整 Sprint 9 门禁见 `sprint9-closure-a-contract.md`

## 背景

Sprint 9 RC readiness 需要 CLI 二进制可自报版本和构建时间，便于 release evidence、客户现场排障和 AI agent 判断构建来源。

## 关键变更

| 文件 | 说明 |
|---|---|
| `cli/ani/main.go` | 新增 `Version` / `BuildTime` 变量和 `ani --version` 输出；继续支持 Makefile `-ldflags` 注入 |
| `cli/ani/main_test.go` | 新增 version 输出测试，确保 `--version` 不触发 Core API 请求 |

## 边界

- 只新增 Core CLI 本地版本输出，不新增 Services 业务命令。
- `--version` 不代表 CLI 已正式发布；发布包、签名和安装渠道仍在后续 Sprint 10 收敛。
