# 模型 — 调用与使用统计

## 页面定位

按 **模型/推理服务** 维度展示调用量、Token 与延迟聚合；**优先扩展 Core Metering** 查询 filter，可选 Services summary 端点。

OpenAPI **详化草案**：`../openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md`。

## 文档管理规则

- 本文是 **模型 — 调用与使用统计** 的主维护源
- **规划权威源**：`docs/console-modules/openapi-drafts/phase3/openapi-phase3-model-usage-stats-draft.md`
- Core 扩展：`ANI-main/repo/api/openapi/v1.yaml`（metering/usage query）
- Services 可选：`ANI-main/repo/api/openapi/services/v1.yaml`（usage-stats summary）
- PRD/SPEC：`tasks/modules/prd/console/inference/prd-console-model-usage-stats.md`、`tasks/modules/spec/console/inference/spec-console-model-usage-stats.md`
- TASK：`TASK-SVC-018` 子项（模型增强）

## Core 层要求（已冻结 + 规划扩展）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/metering/usage` | 租户用量（**规划扩展** `model_id`、`inference_service_id` filter） |
| POST | `/api/v1/metering/token-usage` | Services 上报（非 Console 查询） |

规划 `group_by` 扩展（3b）：`model_id`、`inference_service_id`。

## Services 层要求（规划 · 可选）

| 方法 | 路径 | operationId | 阶段 |
|---|---|---|---|
| GET | `/api/v1/svc/models/{model_id}/usage-stats` | `getModelUsageStatsSummary` | 3a 可选 |
| GET | `/api/v1/svc/models/usage-stats` | `listModelsUsageStatsSummaries` | 3b |

Schema（草案）：`ModelUsageStatsSummary`。

RBAC：Metering 读 scope（与 `usage-report.md` 对齐）。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401` |
| 读权限 | metering/models read | `403` |
| 时间窗 | start_time、end_time 合法 | `400` |
| model 存在 | 单模型查询时 | `404` |

## 页面职责

- 模型中心 Dashboard：请求量、Token 入/出、延迟 p50/p95
- 按日趋势：`group_by=day` + `model_id`
- 按推理服务拆分：`by_inference_service` 或 filter `inference_service_id`
- 跳转 **租户用量报表**（全租户视角）、**推理服务**详情
- 不把 `token-usage` POST 写成查询 API

## 操作可用性矩阵

| 操作 | 只读用户 | 管理员 | API |
|---|---|---|---|
| 租户级 metering | 可用 | 可用 | Core usage |
| 按 model 筛选 | 可用 | 可用 | usage + model_id（扩展） |
| 模型 summary | 可用 | 可用 | Services usage-stats（可选） |
| 导出账单 | 跳转 | 跳转 | billing-export |

## 接口冻结规则

### `GET /api/v1/metering/usage`（Core · 已冻结子集）

- 成功：`200 + items[]`
- 规划扩展 query：`model_id`、`inference_service_id` — **评审前不得写为已冻结**

### `GET /api/v1/svc/models/{model_id}/usage-stats`（规划 · 可选）

- 成功：`200 + ModelUsageStatsSummary`
- 见草案 §4

## 待补边界

- Core vs Services 双路径取舍 — 评审优先 Core filter
- latency 指标来源（access log vs APM）— Services 实现
- 与 `tenant/usage-report.md` preset 视角不重复

## 相关模块

- `tenant/usage-report.md`
- `inference/inference-service.md`
- `inference/model-center.md`

## 验收标准

- [ ] 优先 Core metering 扩展，避免重复存储
- [ ] 与 usage-report group_by 不冲突
- [ ] 合入 YAML 后切换冻结口径
- [ ] PRD/SPEC 已同步
