# Console 辅助维护文档：概览、租户与接入模块

## 1. 适用范围

- `docs/console-modules/home/platform-overview.md`
- `docs/console-modules/tenant/api-key-management.md`
- `docs/console-modules/tenant/security-identity-overview.md`
- `docs/console-modules/tenant/usage-report.md`
- `docs/console-modules/tenant/tenant-management.md`
- `docs/console-modules/integration/open-integration-overview.md`

## 2. 日志规范

| 模块 | 必记事件 | 必记字段 | 脱敏要求 |
|---|---|---|---|
| 平台概览 | `overview_loaded`、`overview_widget_loaded`、`overview_widget_drilldown` | `module`、`widget`、`request_id?`、`result`、`duration_ms?` | 不记录整页聚合原始明细，只记区块摘要 |
| API Key 管理 | `api_key_list_loaded`、`api_key_created`、`api_key_secret_acknowledged`、`api_key_revoked` | `module`、`action`、`request_id`、`key_id?`、`result` | 不记录完整 `key_value`、token、用户敏感凭据 |
| 安全与身份概览 | `security_entry_loaded`、`security_submodule_opened`、`session_flow_redirected` | `module`、`entry`、`request_id?`、`result` | 不记录完整 token、会话密钥或认证断言 |
| 租户用量报表 | `usage_report_loaded`、`usage_filter_changed`、`usage_view_preset_changed` | `module`、`start_time`、`end_time`、`resource_type?`、`group_by?`、`request_id` | 不记录无关租户上下文和非当前筛选明细 |
| 租户管理 | `tenant_member_invited`、`tenant_member_removed`、`tenant_sso_updated`、`tenant_webhook_created` | `module`、`action`、`request_id`、`member_id?`、`result` | 不记录 Webhook secret、邀请 token 明文 |
| 开放与集成总入口 | `integration_entry_loaded`、`integration_path_selected`、`integration_jump_performed` | `module`、`entry_type`、`target`、`result` | 不记录访问凭据、未脱敏文档私有地址 |

## 3. 异常处理方案

| 模块 | 异常场景 | 前端处理 | 运维检查 |
|---|---|---|---|
| 平台概览 | 单个区块加载失败 | 仅区块报错并保留刷新入口，其余区块继续可用 | 检查对应来源模块和聚合口径是否一致 |
| API Key 管理 | 创建失败或吊销失败 | 保留列表上下文，提示 `request_id`，不把失败误报成已创建/已吊销 | 检查 `Auth / api-keys*` 路径返回码与权限 |
| 安全与身份概览 | 会话异常或入口无权限 | 引导重新登录或提示无权限，不展示无效子入口 | 检查 `OIDC / token / refresh / logout` 链路与权限配置 |
| 租户用量报表 | 查询失败、无数据、无权限 | 图表区和表格区分别提示，不伪造 0 值结果 | 检查 `metering/usage` 参数、时间范围与授权 |
| 开放与集成总入口 | 文档不可访问或接入面未准备 | 展示“文档暂不可用”或“先准备凭证”，不伪造替代入口 | 检查文档工件可用性、API Key 子页与 OIDC 入口状态 |

## 4. 运维操作指南

| 模块 | 日常检查 | 常见处置 | 变更注意 |
|---|---|---|---|
| 平台概览 | 检查 5 大主题区块是否都能独立返回 | 某一区块异常时先定位来源模块，再核对聚合映射 | 新增首页卡片前先确认其来源路径和主维护文档 |
| API Key 管理 | 检查列表、创建、吊销三条主链路 | 如出现投诉“看不到明文”，先确认是否是创建后关闭弹窗场景 | 任何新增明文展示能力都必须先过 Core 边界审查 |
| 安全与身份概览 | 检查总入口与 API Key 子页的导航一致性 | 会话问题与权限问题分开处理，不混用文案 | SSO、LDAP、审计等能力扩展前先冻结契约 |
| 租户用量报表 | 检查默认时间窗、预设视角与图表/表格联动 | 出现口径争议时先核对查询参数，不直接改页面文案 | 账单、发票、对账等能力不能直接叠加到本页 |
| 开放与集成总入口 | 检查推荐路径、文档入口、CLI 入口是否可跳转 | 接入准备不足优先回跳安全域，不伪造新入口 | OpenAI、SDK、Webhook 等能力必须等正式契约后再扩页 |

## 5. 版本变更记录

| 模块 | 版本 | 日期 | 变更内容 |
|---|---|---|---|
| 平台概览 | `v1.3` | `2026-06-16` | 复核聚合页无独立冻结接口口径，确认 5 大主题区块摘要、首页回跳与当前可用模块入口一致 |
| API Key 管理 | `v1.4` | `2026-06-16` | 复核与安全总入口、开放与集成及首页回跳的一致性，确认 `api-keys*`、一次性明文与 HTML 命名口径统一 |
| 安全与身份概览 | `v1.4` | `2026-06-16` | 复核 `OIDC begin / token / refresh / logout`、`API Key 管理` 总入口与首页待处理事项回跳口径，并同步 HTML 摘要 |
| 租户用量报表 | `v1.3` | `2026-06-16` | 复核 `GET /api/v1/metering/usage` 单一查询边界，确认预设视角、内联响应与 HTML 摘要一致 |
| 开放与集成总入口 | `v1.5` | `2026-06-16` | 复核 `REST API` 首屏摘要与安全域链路，统一 OIDC、API Key 管理、文档工件、CLI 与首页接入回跳口径，并将 OpenAI / SDK / Webhook / Bot / 第三方集成压回导航占位 |
