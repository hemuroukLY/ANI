# 接入 — Webhook 说明

## 页面定位

澄清 **接入域 Webhook** 与 **租户 Webhook** 的分工，避免重复契约。

本页为 **接入总入口说明页**（非独立 CRUD 资源域），归属 `open-integration-overview.md` 子导航。

## 文档管理规则

- 本文是接入域 Webhook 导航说明的主维护源
- 租户 Webhook CRUD 契约以 `tenant-management.md` 为准
- 投递运维以 `tenant-webhook-ops.md` 为准
- **无**独立的 `/api/v1/svc/integrations/webhooks` 冻结路径

## Services 层要求（引用，非本页新增）

| 能力 | 路径 | 维护模块 |
|---|---|---|
| 租户事件 Webhook CRUD | `/api/v1/svc/tenant/webhooks*` | `tenant-management.md` |
| Webhook 投递日志 | `GET .../webhooks/{id}/deliveries` | `tenant-webhook-ops.md` |
| 第三方集成 | `/api/v1/svc/integrations*` | `integration-third-party.md` |
| Bot 创建 | `POST .../integrations/bots` | `integration-bot.md` |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 接入域读权限 | 租户成员 | `403 FORBIDDEN` |

本页无写接口。**当前 YAML 未声明 `422`**。

## 页面职责

- 导航说明页：配置租户 Webhook → 跳转租户管理；Bot → Bot 页；第三方 → integrations 页
- 展示 Webhook 与 Bot 的能力边界对比表
- 不自造 `IntegrationWebhook` schema 或路径

## 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 阅读说明 / 导航 | 可用 | 可用 |
| 配置租户 Webhook | 跳转 | 可用（tenant-management） |
| 创建 Bot | 跳转 | 可用（integration-bot） |

## 接口冻结规则

本页不定义新 operation。引用 operation 的冻结规则见各子模块详文。

## 待补边界

- 接入域统一 Webhook 注册 API — **YAML 未声明**
- Webhook 签名算法轮换 UI — 待 tenant webhook 扩展
- 平台级出站 Webhook 模板 — 待产品规划

## 相关模块

- `open-integration-overview.md`
- `tenant-management.md`、`tenant-webhook-ops.md`
- `integration-bot.md`、`integration-third-party.md`
- TASK：`TASK-SVC-016`

## 验收标准

- [ ] 文档不重复 tenant webhooks 契约
- [ ] 未自造 integrations/webhooks 路径
- [ ] HTML 摘要指向正确子模块
