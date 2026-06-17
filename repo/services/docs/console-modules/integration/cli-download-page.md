# 接入 — ani CLI

## 页面定位

**非 Console API 交付页**：`ani` CLI 源码入口、构建与常用子命令说明；不等于独立 Console 下载 CDN 页（待运维交付）。

## 文档管理规则

- 本文是 **接入 — ani CLI** 的主维护源
- 一级权威源：`CLI 源码`
- 不得把规划路径或 handler stub 写成已实现
- PRD/SPEC 同步规则见 `docs/console-modules/governance/module-delivery-workflow.md` §3.5
- **TODO-YAML: N/A** — 非 Console API 交付页

## Core 层要求

N/A

## 仓内源码路径（只读引用）

- `ANI-main/repo/cli/ani/main.go`

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证（若页面需登录） | 页面层拦截 |
| OpenAPI | N/A | 本页不调用专有 API |

无 POST 创建接口。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员/运维 | 说明 |
|---|---|---|---|
| 查看 CLI 说明 | 可用 | 可用 | 静态文档 |
| 二进制下载 | 不可用 | 不可用 | 无冻结 download API |

## 页面职责

- 占位 UI + 明确 YAML/OpenAPI 缺口（若适用）
- 跳转关联模块（见「相关模块」）
- 不把 BOSS/平台运营能力写入 Console 冻结契约

## 接口冻结规则

**当前无本模块冻结 API（TODO-YAML: N/A）。**

本页为 **非 Console API 交付页**，不新增 OpenAPI 路径；接入契约以 Core/Services YAML 为准。

## 待补边界

- Release 产物与 checksum — 运维待补
- <!-- TODO-YAML: N/A -->

## 相关模块

- `integration/open-integration-overview.md`

## 验收标准

- [ ] 路径与 OpenAPI 权威源一致（或明确 TODO-YAML / N/A）
- [ ] 正文不把 handler stub 写成已实现
- [ ] 含创建前置条件与操作可用性矩阵
- [ ] PRD/SPEC/HTML 摘要已同步
