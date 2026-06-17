# 合规能力

## 页面定位

展示租户 **合规状态摘要**（控制项、待办）与 **历史报告** 只读列表；完整合规包生成与取证归 **BOSS**。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-compliance-draft.md`。

## 文档管理规则

- 本文是 **合规能力** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-compliance-draft.md`
- 合入后一级权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- PRD/SPEC：`tasks/modules/prd/console/security/prd-console-compliance.md`、`tasks/modules/spec/console/security/spec-console-compliance.md`
- TASK：`TASK-SVC-018` 子项 / Phase 3 §4

## Services 层要求（规划 · 待合入 YAML · Console 只读）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/compliance/summary` | `getComplianceSummary` | 3a |
| GET | `/api/v1/svc/compliance/reports` | `listComplianceReports` | 3a |
| GET | `/api/v1/svc/compliance/reports/{report_id}` | `getComplianceReport` | 3a |

Schema（草案）：`ComplianceSummary`、`ComplianceControlItem`、`ComplianceReport`。

RBAC（草案）：`scope:compliance:read`（**无 write**）。

## Core 层要求

无 Core Compliance 路径；平台审计见 `audit-log.md`（规划）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读权限 | compliance:read | `403` |
| report 存在 | 有效 id | `404` |

## 页面职责

- Dashboard：overall_status、控制项 checklist、open_actions
- 报告列表：类型、周期、status、download_url（ready 时）
- **无** Console 侧创建报告 POST
- 跳转 **审计日志**；BOSS 能力标注「非 Console API」

## 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 | 说明 |
|---|---|---|---|
| 合规摘要 | 可用 | 可用 | GET summary |
| 报告 list/detail | 可用 | 可用 | GET reports |
| 下载报告 | 可用 | 可用 | download_url |
| 请求新报告 | 不可用 | 不可用 | BOSS |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-compliance-draft.md` §4。

### `GET /api/v1/svc/compliance/summary`（规划）

- 成功：`200 + ComplianceSummary`

### `GET .../reports/{report_id}`（规划）

- 成功：`200 + ComplianceReport`

## 待补边界

- framework 枚举（iso27001、等保2.0）— BOSS 配置
- download_url 签发 TTL — BOSS
- 与 audit-log 导出联动 — 只读跳转

## 相关模块

- `security/audit-log.md`
- `tenant/tenant-management.md`
- `tenant/billing-export.md`

## 验收标准

- [ ] Console 无 write path
- [ ] BOSS 边界在 UI 明示
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
