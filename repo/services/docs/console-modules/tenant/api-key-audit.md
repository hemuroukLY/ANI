# API Key 审计

## 页面定位

检索 **API Key 使用记录**（调用时间、来源 IP、路径摘要、结果码），扩展 `api-key-management.md` 的凭证 CRUD。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md`。

## 文档管理规则

- 本文是 **API Key 审计** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/v1.yaml`（Core Auth）
- 不得把草案路径写成**已实现**
- PRD/SPEC：`tasks/modules/prd/console/tenant/prd-console-api-key-audit.md`、`tasks/modules/spec/console/tenant/spec-console-api-key-audit.md`
- TASK：`TASK-CORE-015` 子项 / Phase 3 §4

## Core 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/auth/api-keys/{key_id}/audit-events` | `listApiKeyAuditEvents` | 3a |
| GET | `.../audit-events/{event_id}` | `getApiKeyAuditEvent` | 3b |

Schema（草案）：`ApiKeyAuditEvent`、`ApiKeyAuditEventListResponse`。

RBAC（草案）：`scope:auth:api-keys:audit:read`。

Auth CRUD 只读引用：`/api/v1/auth/api-keys*` — `api-key-management.md`。

## Services 层要求

无 Services 路径。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 审计读权限 | auth api-keys audit read | `403` |
| key 存在 | 有效 `key_id`、租户可见 | `404` |
| 时间窗 | start_time ≤ end_time | `400` |

## 页面职责

- 从 Key 列表进入：按 key_id 过滤的调用事件时间线
- 筛选：时间窗、status_code、path_prefix
- **不**展示 key_value 明文；不替代 BOSS 取证导出
- 跳转 **API Key 管理**、**平台审计日志**

## 操作可用性矩阵

| 操作 | 只读用户 | Key 管理员 | YAML 合入后 |
|---|---|---|---|
| 审计列表 | 不可用 | 可用 | GET audit-events |
| 单条详情 | 不可用 | 可用 | GET event（3b） |
| BOSS 导出 | 不可用 | 不可用 | — |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-api-key-audit-draft.md` §4。

### `GET .../audit-events`（规划）

- 成功：`200 + ApiKeyAuditEventListResponse`
- 错误：`401`、`403`、`404`

## 待补边界

- 与 `audit-log.md` 统一事件 schema 映射
- 保留周期与 GDPR — BOSS 策略
- audit read scope 是否与 api-keys:read 合并 — 评审

## 相关模块

- `tenant/api-key-management.md`
- `security/audit-log.md`

## 验收标准

- [ ] 与 auth api-keys CRUD 无 path 冲突
- [ ] 不记录 key 明文
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
