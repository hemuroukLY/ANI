# CORE-CLI-B — Core CLI expanded resource coverage

完成日期：2026-06-04
对应 Sprint：Sprint 8（Core-only）
验证结果：`make validate-core-cli`、`make build-cli`、完整 Sprint 8 门禁见 `sprint8-closure-a-contract.md`

## 背景

Sprint 7 的 Core CLI 只覆盖最小资源。Sprint 8 扩展只读命令覆盖主要 Core 资源，继续拒绝 Services 业务资源，保持 CLI 只作为 ANI Core 控制面入口。

## 关键变更

| 文件 | 说明 |
|---|---|
| `cli/ani/main.go` | 扩展 `network-*`、`volumes`、`filesystems`、`objects`、`vector-stores`、`encryption-keys`、`observability-*` 命令 |
| `cli/ani/main_test.go` | 新增 CLI-B 解析测试和 observability query 参数测试 |

## 新增命令

```bash
ani network-vpcs list
ani network-subnets list
ani network-security-groups list
ani network-load-balancers list
ani volumes list
ani filesystems list
ani objects list
ani vector-stores list
ani encryption-keys list
ani observability-alert-rules list
ani observability-query get --query up --tenant-id tenant-a
```

## 边界

- 只新增 Core REST 只读请求和 query helper。
- 不实现 `model`、`kb`、`inference` 等 Services 命令。
- 不代表 CLI 全资源覆盖、交互式配置、发布包或生产支持完成。
