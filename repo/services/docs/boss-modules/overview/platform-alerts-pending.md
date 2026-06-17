# 平台告警与待处理

## 页面定位

`平台告警与待处理` 是 BOSS **平台运营总览** 域下的 **全平台风险与待办** 明细页，汇总 P0/P1 当前告警、失败任务、推理/KB 异常与待确认事项，供 SRE 与平台运营统一处置入口。

Console 对照：[`alerts-pending-items.md`](../../console-modules/alerts/alerts-pending-items.md)（租户侧）。深度专页：[`alert-rules.md`](../health/alert-rules.md)、[`job-history.md`](../health/job-history.md)、[`platform-alerts-pending` 上级摘要](operations-overview.md)。

## 文档管理规则

- 本文是 `平台告警与待处理` 的主维护源
- PRD/SPEC 为辅助
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- **无** `GET /api/v1/alerts*` 或 `listPendingItems` 冻结路径
- `GET /api/v1/observability/alert-rules` — 告警**规则** CRUD；**非** firing 事件
- `CreateObservabilityAlertRuleRequest` **不含 `tenant_id`**
- `GET /api/v1/tasks/{task_id}` — **路径已声明**；**无** list（`GET /api/v1/tasks` — **TODO-YAML**）
- `GET /api/v1/observability/query` — PromQL 辅助判断
- Services `inference-services`、`knowledge-bases` — 租户 **只读参考**，用于异常推断
- 平台 firing 告警 list / 待办 aggregate — **TODO-YAML**
- `AsyncTask.status`：`completed`（禁 `succeeded`）

## 页面职责

- 汇总 P0/P1 告警数、失败任务数、跨租户异常服务数
- 按来源分类列表（告警事件、任务、推理、KB、交付）
- 提供跳转：告警规则、任务历史、推理/KB 监控、运维 Skills
- 区块级刷新与局部失败态
- **不承担** 完整规则 CRUD（跳转 alert-rules）

## 页面结构

```text
平台告警与待处理
├── 顶部摘要（P0 / P1 / 失败任务 / 待确认）
├── 分类列表
│   ├── 活动告警（firing — 待 YAML）
│   ├── 异步任务（list 待 YAML；单查已声明）
│   ├── 推理服务异常
│   ├── 知识库/文档异常
│   └── 交付/验收待办（settings 待收口）
└── 跳转
    ├── 告警规则
    ├── 任务历史
    └── 运维 Skills
```

## 数据来源与分层约束

| 摘要字段 | 建议来源 | 说明 |
|---|---|---|
| P0/P1 告警数 | firing events **TODO-YAML** 或 PromQL+规则 | **不得** 用 rules 列表长度冒充 |
| 失败任务数 | `GET /tasks` list **TODO-YAML** | 单查 `{task_id}` 已声明 |
| 推理异常数 | 平台 inference aggregate **TODO-YAML** | 租户 list 仅参考 |
| KB 异常数 | 平台 KB aggregate **TODO-YAML** | 同上 |
| 待确认事项 | 各模块 + Webhook | 口径待产品冻结 |

### 关键边界

- `alert-rules` = **规则**，不是 **告警事件**
- BOSS 代建其他租户规则须 platform API / **TODO-YAML**
- 403 示例限定 list **TODO-YAML** 场景；单查 YAML 仅声明 404

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| P0/P1 摘要 | firing **TODO-YAML** | alert-rules |
| 活动告警列表 | **TODO-YAML** | alert-rules |
| 失败任务 | tasks list **TODO-YAML** | job-history |
| 推理异常 | aggregate / Services 参考 | inference-monitoring |
| KB 异常 | aggregate / Services 参考 | kb-monitoring |

## BOSS 与 Console 分工

| 维度 | BOSS 本页 | Console alerts-pending |
|---|---|---|
| 范围 | 全平台跨租户 | 当前租户 |
| 告警 | 平台 firing 待 YAML | 租户推断 + 局部任务 |
| 任务 list | 平台 **TODO-YAML** | 会话内已知 task_id |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/observability/alert-rules` | `listObservabilityAlertRules` | **规则**；`scope:observability:read` |
| GET | `/api/v1/tasks/{task_id}` | — | **路径已声明** |
| GET | `/api/v1/observability/query` | `queryObservability` | PromQL |
| GET | `/api/v1/svc/inference-services` | `listInferenceServices` | 租户参考 |
| GET | `/api/v1/svc/knowledge-bases` | `listKnowledgeBases` | 租户参考 |

| 能力 | 状态 |
|---|---|
| firing 告警 events list/ack/silence | **ADDED-TO-YAML** P1 |
| `GET /api/v1/tasks` list | **ADDED-TO-YAML** P1 |
| 平台 pending aggregate | **ADDED-TO-YAML** `getPlatformPendingSummary` |

## 字段级定义

| 字段 | 说明 |
|---|---|
| `alert_p0_count` | P0 当前告警；无 API 时「待接入」 |
| `alert_p1_count` | P1 当前告警 |
| `failed_task_count` | `failed` / `dead_letter` |
| `running_task_count` | `pending` / `running` |
| `inference_anomaly_count` | 失败推理服务数 |
| `kb_anomaly_count` | 解析失败等 |
| `pending_confirm_count` | 待确认事项 |
| `last_refreshed_at` | 聚合刷新时间 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 正常 | 摘要 + 分类列表 |
| alert-events 已 ADDED-TO-YAML | 告警数「待接入」；**不** 用 rules 数冒充 |
| listTasks 已 ADDED-TO-YAML | 任务区说明 TODO-YAML |
| P0 > 0 | 红色高亮 |
| 单分类失败 | 局部失败 + 重试 |
| 403 | 无权限 |

## 字段口径与单位

| 字段 | 口径 |
|---|---|
| `AsyncTask.status` | YAML enum；禁 `succeeded` |
| 告警 severity | `info` / `warning` / `critical`（规则 schema） |

## 状态与能力口径

### AsyncTask.status（任务）

`pending` / `running` / `completed` / `failed` / `cancelled` / `dead_letter`

本页可跳转处置；写操作（ack/silence/retry）在专页或 **TODO-YAML** API。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台 observability 读（**TODO-YAML** 扩展） | `403` |
| 未认证 | `401` |
| 跳转 alert-rules 写 | `scope:observability:write` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 管理员 |
|---|---|---|---|
| 查看待处理 | ✅ | ✅ | ✅ |
| 确认/静默告警 | 待 YAML | ✅ | ✅ |
| 跳转规则/任务/Skills | ✅ | ✅ | ✅ |
| 在本页 CRUD 规则 | ❌ | → alert-rules | → alert-rules |

## 删除前置校验

**N/A**（本页无 DELETE；规则删除在 alert-rules）

## 接口冻结规则

### `GET /api/v1/observability/alert-rules`（**规则** · 非事件）

- CRUD 五 path 已冻结；`scope:observability:read|write`
- **不得** 把 list rules 当作 firing 告警 list

### `GET /api/v1/tasks/{task_id}`（**路径已声明**）

- 无 operationId / 无 x-ani-rbac-scope
- `success`：`200` + `AsyncTask`；`errors`：`404`（YAML）
- **无** list path

### `GET /api/v1/tasks` list（待补）

<!-- TODO-YAML: GET /api/v1/tasks cursor + tenant filter -->

- 403 示例适用于 **list** 场景

### firing events（待补）

<!-- TODO-YAML: observability alert events list/ack/silence -->

## 使用规则

- 聚合在页面层；不写入新 services schema
- 禁止用 rules 数量冒充告警数
- 禁止跨租户 JWT 轮询
- BOSS 跨租户创建规则须 platform API

## 待补能力边界

- firing events API — **ADDED-TO-YAML**
- tasks list — **ADDED-TO-YAML**
- 平台 pending aggregate — **ADDED-TO-YAML**
- 批量 ack/silence — Phase 2

## 响应示例

### 单任务查询成功（**路径已声明 · 非 list**）

```json
{
  "task_id": "task-abc123",
  "task_type": "kb.reindex",
  "status": "failed",
  "progress": 0.45,
  "error": { "code": "INDEX_TIMEOUT", "message": "vector index timeout" },
  "created_at": "2026-06-17T08:00:00Z",
  "updated_at": "2026-06-17T08:05:00Z"
}
```

### 平台 pending aggregate 目标（**待 YAML**）

```json
{
  "alert_p0_count": 2,
  "alert_p1_count": 5,
  "failed_task_count": 3,
  "inference_anomaly_count": 8,
  "kb_anomaly_count": 4,
  "last_refreshed_at": "2026-06-17T10:00:00Z"
}
```

## 错误示例

### 平台 tasks list 无权限（**TODO-YAML list 场景**）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-pap-403-001"
}
```

> 注：403 示例限定 **list TODO-YAML**；单查 YAML 仅声明 404。

### 任务不存在（单查）

```json
{
  "code": "NOT_FOUND",
  "message": "task not found",
  "request_id": "req-boss-pap-404-001"
}
```

### 非法筛选参数（平台 list · **TODO-YAML** 目标）

```json
{
  "code": "BAD_REQUEST",
  "message": "status filter must be one of: pending, running, completed, failed, cancelled",
  "request_id": "req-boss-pap-400-001"
}
```

## 相关模块

- [`operations-overview.md`](operations-overview.md)、[`alert-rules.md`](../health/alert-rules.md)、[`job-history.md`](../health/job-history.md)
- Console：[`alerts-pending-items.md`](../../console-modules/alerts/alerts-pending-items.md)

## 回填验收标准

- [x] 满配章节齐全
- [x] rules vs firing、tasks 单查 vs list 分层清晰
- [x] `CreateObservabilityAlertRuleRequest` 不含 tenant_id 边界
- [x] 403 + 404 + 400 错误示例
- [ ] firing / tasks list YAML 合入后回写
- [x] PRD/SPEC/HTML 与本文同步
