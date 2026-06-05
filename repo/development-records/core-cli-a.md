# CORE-CLI-A — ANI Core CLI minimal contract

完成日期：2026-06-04
对应 Sprint：Sprint 7（Core-only）
验证结果：`make validate-core-cli`、`make build-cli`、完整 Sprint 7 门禁见 `sprint7-closure-a-contract.md`

## 实现了什么

新增 `cli/ani` Go 模块，提供 `ani` Core CLI 的最小资源访问闭环。CLI 通过 Core REST API base URL 和 bearer token 调用现有 Core API，当前覆盖实例、K8s 集群、Secrets、Registry projects 和 metering usage 的只读入口，并显式拒绝 Services 业务资源命令。

当前只证明 Core CLI contract/minimal behavior，不代表 CLI 已覆盖所有 Core 资源、交互体验或发布包。

## 关键文件改动

| 文件 | 修改说明 |
|---|---|
| `cli/ani/go.mod` | 新增 Core CLI Go 模块 |
| `cli/ani/main.go` | 新增 CLI 解析、Core-only 资源限制和 HTTP request 执行 |
| `cli/ani/main_test.go` | 覆盖 Services 资源拒绝、实例 list 请求构造和 bearer token 注入 |
| `go.work` | 将 `./cli/ani` 纳入 workspace |
| `Makefile` | 新增 `build-cli`、`validate-core-cli` |

## 边界

- CLI 不新增 Services 命令。
- CLI 不重复定义 Core API schema，只调用现有 Core REST path。
- 当前不宣称 CLI 覆盖 Core 全资源或 production release。
