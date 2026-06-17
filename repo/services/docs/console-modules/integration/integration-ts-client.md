# 接入 — TypeScript Client

## 页面定位

**非 Console API 交付页**：TypeScript / openapi-typescript 生成客户端与 Console 前端共用类型策略说明。

## 文档管理规则

- 本文是 **接入 — TypeScript Client** 的主维护源
- 一级权威源：`仓内 SDK 源码`
- 不得把规划路径或 handler stub 写成已实现
- PRD/SPEC 同步规则见 `docs/console-modules/governance/module-delivery-workflow.md` §3.5
- **TODO-YAML: N/A** — 非 Console API 交付页

## Core 层要求

N/A

## 仓内源码路径（只读引用）

- `ANI-main/repo/sdks/core/typescript/src/index.ts`
- `ANI-main/repo/sdks/services/typescript/src/index.ts`
- `ANI-main/repo/frontends/console/src/api/client.ts`
- `ANI-main/repo/Makefile（openapi-typescript 生成命令）`

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证（若页面需登录） | 页面层拦截 |
| OpenAPI | N/A | 本页不调用专有 API |

无 POST 创建接口。

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员/运维 | 说明 |
|---|---|---|---|
| 查看 Client 文档 | 可用 | 可用 | Console 静态页 |
| npm 安装 | 可用 | 可用 | @kubercloud/ani-client 规划名 |

## 页面职责

- 占位 UI + 明确 YAML/OpenAPI 缺口（若适用）
- 跳转关联模块（见「相关模块」）
- 不把 BOSS/平台运营能力写入 Console 冻结契约

## 接口冻结规则

**当前无本模块冻结 API（TODO-YAML: N/A）。**

本页为 **非 Console API 交付页**，不新增 OpenAPI 路径；接入契约以 Core/Services YAML 为准。

## 待补边界

- <!-- TODO-YAML: N/A -->

## 相关模块

- `integration/open-integration-overview.md`
- `integration/openai-compatible-api.md`

## 验收标准

- [ ] 路径与 OpenAPI 权威源一致（或明确 TODO-YAML / N/A）
- [ ] 正文不把 handler stub 写成已实现
- [ ] 含创建前置条件与操作可用性矩阵
- [ ] PRD/SPEC/HTML 摘要已同步
