# 接入 — Go SDK

## 页面定位

**非 Console API 交付页**：指向仓内 Go SDK 源码与示例，帮助租户自助接入 Core/Services HTTP API。

## 文档管理规则

- 本文是 **接入 — Go SDK** 的主维护源
- 一级权威源：`仓内 SDK 源码（只读引用）`
- 不得把规划路径或 handler stub 写成已实现
- PRD/SPEC 同步规则见 `docs/console-modules/governance/module-delivery-workflow.md` §3.5
- **TODO-YAML: N/A** — 非 Console API 交付页

## Core 层要求

N/A — 不新增 Console 冻结路径。

## 仓内源码路径（只读引用）

- `ANI-main/repo/sdks/core/go/anisdk/client.go`
- `ANI-main/repo/sdks/core/go/examples/basic/main.go`
- `ANI-main/repo/sdks/services/go/anisdk/client.go`
- `ANI-main/repo/sdks/services/go/examples/basic/main.go`

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证（若页面需登录） | 页面层拦截 |
| OpenAPI | N/A | 本页不调用专有 API |

无 POST 创建接口。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员/运维 | 说明 |
|---|---|---|---|
| 查看 SDK 文档 | 可用 | 可用 | Console 静态页 |
| 下载/安装指引 | 可用 | 可用 | go get / module path |

## 页面职责

- 占位 UI + 明确 YAML/OpenAPI 缺口（若适用）
- 跳转关联模块（见「相关模块」）
- 不把 BOSS/平台运营能力写入 Console 冻结契约

## 接口冻结规则

**当前无本模块冻结 API（TODO-YAML: N/A）。**

本页为 **非 Console API 交付页**，不新增 OpenAPI 路径；接入契约以 Core/Services YAML 为准。

## 待补边界

- SDK 版本与 OpenAPI 生成流水线 — 见 repo Makefile
- <!-- TODO-YAML: N/A -->

## 相关模块

- `integration/open-integration-overview.md`
- `integration/integration-ts-client.md`

## 验收标准

- [ ] 路径与 OpenAPI 权威源一致（或明确 TODO-YAML / N/A）
- [ ] 正文不把 handler stub 写成已实现
- [ ] 含创建前置条件与操作可用性矩阵
- [ ] PRD/SPEC/HTML 摘要已同步
