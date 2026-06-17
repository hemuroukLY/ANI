# BOSS integration 域 Phase 0 GAP 摘要

> **日期**：2026-06-17 · **详文**：`docs/boss-modules/integration/*.md`

## Services 租户 path（只读参考 · 非 BOSS 正式）

| 路径 | operationId | 模块 |
|---|---|---|
| `GET/POST /api/v1/svc/tenant/webhooks` | listWebhooks / createWebhook | ops-webhook 对照 |
| `GET .../deliveries` | listWebhookDeliveries | ops-webhook 对照 |
| `DELETE .../{webhook_id}` | deleteWebhook | ops-webhook 对照 |
| `GET/POST /api/v1/svc/integrations` | listIntegrations / createIntegration | ops-system-integration 对照 |
| `POST /api/v1/svc/integrations/bots` | createIntegrationBot（422 已声明） | enterprise-notification 对照 |

## Core 缺口

| 项 | 结论 |
|---|---|
| platform ops webhook | **不存在** |
| platform bot 治理 list | **不存在** |
| platform integrations aggregate | **不存在** |
| `{integration_id}` GET/PATCH/DELETE | **不存在** |
| bot DELETE/PATCH | **不存在** |

## Phase 2 建议批次

| 优先级 | 能力 |
|---|---|
| **P1** | Core platform ops webhook CRUD + deliveries |
| **P2** | 平台 Bot 渠道 list/CRUD |
| **P2** | 平台第三方集成治理 aggregate + 连通性探测 |

## 架构待决

- 平台 RBAC scope（integration / notifications）
- platform ops webhook schema 是否独立于 tenant `Webhook`
- Bot 平台实例 vs 租户 Bot 是否分表

门禁：`python3 scripts/validate-boss-integration-gate.py`
