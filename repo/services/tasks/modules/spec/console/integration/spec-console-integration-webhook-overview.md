# SPEC: Console integration-webhook-overview

> Source: `tasks/modules/prd/console/integration/prd-console-integration-webhook-overview.md`  
> Revised: 2026-06-17

## 1. Summary

澄清 **接入域 Webhook** 与 **租户 Webhook** 的分工，避免重复契约。本页为 **接入总入口说明页**（非独立 CRUD 资源域），归属 `open-integration-overview.md` 子导航。

## 2. Frozen Facts

### 2.1 Authority Source

- 引用：`ANI-main/repo/api/openapi/services/v1.yaml`（tenant webhooks、integrations）
- **无**独立的 `/api/v1/svc/integrations/webhooks` 冻结路径

### 2.2 Verified Paths

本页不定义新 operation。引用路径如下：

| 能力 | Method | Path | 维护模块 |
|---|---|---|---|
| 租户事件 Webhook CRUD | （见 tenant-management） | `/api/v1/svc/tenant/webhooks*` | `tenant-management.md` |
| Webhook 投递日志 | GET | `/api/v1/svc/tenant/webhooks/{webhook_id}/deliveries` | `tenant-webhook-ops.md` |
| 第三方集成 | GET/POST | `/api/v1/svc/integrations*` | `integration-third-party.md` |
| Bot 创建 | POST | `/api/v1/svc/integrations/bots` | `integration-bot.md` |

### 2.3 Verified Schemas

- 不自造 `IntegrationWebhook` schema 或路径
- 各子模块 schema 以对应 OpenAPI 声明为准

## 3. Page Scope

- 导航说明页：配置租户 Webhook → 跳转租户管理；Bot → Bot 页；第三方 → integrations 页
- 展示 Webhook 与 Bot 的能力边界对比表
- 引用 operation 的冻结规则见各子模块详文

## 4. Non-Goals

- 接入域统一 Webhook 注册 API — **YAML 未声明**
- Webhook 签名算法轮换 UI — 待 tenant webhook 扩展
- 平台级出站 Webhook 模板 — 待产品规划
- 本页无写接口；**当前 YAML 未声明 `422`**

## 5. 创建前置条件

| 依赖项 | 要求 | 未满足时 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 接入域读权限 | 租户成员 | `403 FORBIDDEN` |

## 6. 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 阅读说明 / 导航 | ✅ | ✅ |
| 配置租户 Webhook | 跳转 | ✅（tenant-management） |
| 创建 Bot | 跳转 | ✅（integration-bot） |

## 7. 主维护源

- `docs/console-modules/integration/integration-webhook-overview.md`
- 父级：`open-integration-overview.md`
- 相关：`tenant-management.md`、`tenant-webhook-ops.md`、`integration-bot.md`、`integration-third-party.md`
- TASK：`TASK-SVC-016`

## 8. 验收要点

- 文档不重复 tenant webhooks 契约
- 未自造 integrations/webhooks 路径
- HTML 摘要指向正确子模块
