# 账单导出

## 页面定位

租户 **用量报表异步导出**（CSV/PDF）；在线查询见 `usage-report.md`；结算/发票归 **BOSS**。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-billing-export-draft.md`。

## 文档管理规则

- 本文是 **账单导出** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-billing-export-draft.md`
- Core 只读：`GET /api/v1/metering/usage`
- 合入后 Services 权威源：`ANI-main/repo/api/openapi/services/v1.yaml`
- PRD/SPEC：`tasks/modules/prd/console/tenant/prd-console-billing-export.md`、`tasks/modules/spec/console/tenant/spec-console-billing-export.md`
- TASK：`TASK-SVC-018` 子项 / Phase 3 §4

## Core 层要求（只读引用）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/metering/usage` | 导出数据源 — `usage-report.md` |
| GET | `/api/v1/tasks/{task_id}` | 大 export 202 后轮询 |

## Services 层要求（规划 · 待合入 YAML）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| POST | `/api/v1/svc/billing/exports` | `createBillingExport` | 3a |
| GET | `/api/v1/svc/billing/exports/{export_id}` | `getBillingExport` | 3a |
| GET | `/api/v1/svc/billing/exports` | `listBillingExports` | 3b |

Schema（草案）：`CreateBillingExportRequest`、`BillingExportJob`。

RBAC（草案）：`scope:billing:read`、`scope:billing:write`。

异步：`202 + AsyncTask`，`task_type` 建议 `billing.export`。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 写权限 | billing:write | `403` |
| POST 幂等键 | 必填 | `400` |
| 时间窗 | start ≤ end，跨度上限 | `422`（建议语义：`EXPORT_WINDOW_TOO_LARGE`；**待 YAML 声明后定稿**） |
| 并发 export | 未超租户限额 | `422`（建议语义：`EXPORT_RATE_LIMITED`；**待 YAML 声明后定稿**） |

## 页面职责

- 从用量页发起 export：类型、时间窗、resource_type/model_id filter
- 进行中跳转 **异步任务中心**；ready 后 download_url
- export **不含** BOSS 结算金额（用量明细 only）
- 跳转 **usage-report**

## 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 | 说明 |
|---|---|---|---|
| 在线用量 | 可用 | 可用 | Core metering |
| 创建 export | 不可用 | 可用 | POST |
| 下载文件 | 不可用 | 可用 | download_url |
| export 历史 | 不可用 | 可用 | list（3b） |
| 发票/对账 | 不可用 | 不可用 | BOSS |

## 接口冻结规则

> **规划草案**；见 `../openapi-drafts/phase3/openapi-phase3-billing-export-draft.md` §4。

### `POST /api/v1/svc/billing/exports`（规划）

- 成功：`202 + AsyncTask`（默认）
- 错误：`400`、`422`

### `GET .../exports/{export_id}`（规划）

- 成功：`200 + BillingExportJob`

## 待补边界

- export_type 与 PDF 模板 — BOSS/Services
- download_url 存储（对象存储）与 TTL
- 与 model-usage-stats model_id filter 对齐

## 相关模块

- `tenant/usage-report.md`
- `inference/model-usage-stats.md`
- `alerts/async-task-center.md`

## 验收标准

- [ ] 大 export 202 + AsyncTask
- [ ] 与 BOSS 结算边界清晰
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
