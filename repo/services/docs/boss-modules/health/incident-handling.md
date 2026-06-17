# 故障处理

## 页面定位

`故障处理` 是 BOSS **运维与可观测** 域下的 **平台 Incident 工单** 聚合页：关联活动告警、AsyncTask 失败、平台健康降级，支撑分派、时间线、 remediation 与复盘。

本页属于 **Phase 2 产品能力**；Incident CRUD / timeline / 分派 API **整域 TODO-YAML**。Phase 2 文档阶段以 **只读跳转** 已有满配模块（[`alert-rules.md`](alert-rules.md)、[`job-history.md`](job-history.md)、[`platform-health.md`](platform-health.md)）为主，**不得** 自造 REST 契约为已实现。

一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`（当前无 `/incidents*` path；可引用 `GET /api/v1/tasks/{task_id}`）。

Console [`alerts-pending-items.md`](../../console-modules/alerts/alerts-pending-items.md) 仅为租户侧摘要；**不承担** 平台 Incident 生命周期。

## 文档管理规则

- 本文是 `故障处理` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-incident-handling.md`](../../tasks/modules/prd/boss/health/prd-boss-incident-handling.md) 与 [`spec-boss-incident-handling.md`](../../tasks/modules/spec/boss/health/spec-boss-incident-handling.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 文档 + Phase 2 能力规划 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- Incident API — **整域 TODO-YAML P2**；**不得** 自造 `GET/POST/PATCH /api/v1/incidents*` 为已冻结
- 只读数据源跳转（**不重复定义** 其 API）：
  - 活动告警：[`alert-rules.md`](alert-rules.md) Tab 2 · events **TODO-YAML**
  - 失败任务：[`job-history.md`](job-history.md) · `GET /tasks/{task_id}` + list **TODO-YAML**
  - 健康降级：[`platform-health.md`](platform-health.md) · aggregate **TODO-YAML**
- 写操作（创建/分派/关闭/评论）须 `idempotency_key` — 待 Incident API 冻结
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）
- 禁止自造 `Incident`、`IncidentTimelineEntry` schema / operationId
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 展示 **进行中 / 已解决** Incident list（API 待 YAML）；Phase 2 前可 UI 占位 + 跳转入口
- 支持 **新建 Incident** 或 **从活动告警升级**（POST 待 YAML）
- 展示 **时间线**（分派、确认、Skill 触发、状态变更）与 **负责人**
- 支持 **触发 Skill** → [`maint-skills.md`](maint-skills.md)（预填 incident 上下文）
- 支持 **复盘与关闭**（mitigated → resolved；POST 待 YAML）
- 可选通知：跳转 [`../integration/enterprise-notification.md`](../integration/enterprise-notification.md)（集成 API 待 YAML）
- **不承担** Console 租户告警摘要；**不承担** 告警规则 CRUD（见 alert-rules）

## 页面结构

- 首屏至少包含：`筛选区`、`Incident 列表`、`新建/升级入口`、`边界说明`
- 详情为 drawer 或独立路由：时间线 + 关联对象 + 操作栏
- API 未就绪时：列表占位 + 新建/分派 disabled + 深链已有模块
- 无权限态、写失败态、关联对象 404 态须可区分

```text
故障处理
├── 筛选（status / severity / assignee / tenant — 待 YAML）
├── Incident 列表（incident_id / title / severity / status / assignee）
├── 新建 Incident / 从告警升级
├── 详情 drawer
│   ├── 摘要（title / severity / status / assignee）
│   ├── 关联（related_alerts[] / related_tasks[] / health_context）
│   ├── 时间线 timeline[]
│   └── 操作（分派 / 触发 Skill / mitigated / resolve — 待 YAML）
└── 复盘区（resolved 后 · Phase 2）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 / 模块 | 本页用法 |
|---|---|---|
| Core | Incident CRUD **TODO-YAML** | BOSS 正式工单数据源 |
| Core | alert events **TODO-YAML** | 关联告警；跳转 alert-rules |
| Core | `GET /api/v1/tasks/{task_id}` | 关联失败任务只读 |
| Core | `GET /api/v1/tasks` list **TODO-YAML** | job-history 正式 list |
| Core | platform-health aggregate **TODO-YAML** | 降级上下文 |
| 产品 | 深链 query 参数 | `?from=alert&incident_id=` 等 |

### 关键边界

- **不得** 在本页重复定义 firing 告警 API（见 alert-rules）
- **不得** 把 `ObservabilityAlertRule` 或 rules list 当作 Incident list
- Incident `status`（open/mitigated/resolved）是 **工单状态**，与 `AsyncTask.status`、`alert firing status` **不同域**
- 从告警升级时须引用 **alert incident id**（events API 待 YAML），不得用 `rule_id` 代替
- 企业 IM 通知走 integration 域，**不得** 在 Incident API 未冻结时写死 webhook 路径

## 页面区块与数据来源映射

| 区块 | 主要来源 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 筛选区 | Incident list API **待 YAML** | status/severity/assignee | 刷新列表 |
| Incident 列表 | Incident CRUD **待 YAML** | 工单摘要字段 | 打开详情 |
| 新建/升级 | POST incidents **待 YAML** | `idempotency_key` | 详情 drawer |
| 关联告警 | alert events **待 YAML** | related_alerts[] | alert-rules |
| 关联任务 | job-history / GET task | related_tasks[] | job-history |
| 健康上下文 | platform-health **待 YAML** | 组件 degraded/error | platform-health |
| 触发 Skill | maint-skills execute **待 YAML** | 预填 incident_id | maint-skills |
| 时间线 | Incident timeline **待 YAML** | 审计化操作记录 | platform-audit-log |
| 通知 | enterprise-notification **待 YAML** | 分派/升级通知 | integration |

## BOSS 与 Console 分工

| 维度 | BOSS 故障处理 | Console |
|---|---|---|
| 范围 | 平台 Incident 全生命周期 | 租户告警摘要 |
| 工单 API | **TODO-YAML** P2 | 无 |
| 活动告警 | alert-rules + events 待 YAML | alerts-pending 客户端聚合 |
| 任务失败 | job-history | async-task-center |
| Skill remediation | maint-skills | 无 |
| RBAC | 平台 incident read/write | 租户 JWT |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/tasks/{task_id}` | **未声明** | 关联任务只读 |

| 数据源 | 模块 | 状态 |
|---|---|---|
| 活动告警 events | alert-rules | **ADDED-TO-YAML** `/observability/alert-events` |
| 任务 list | job-history | **ADDED-TO-YAML** `listTasks` |
| 平台健康 aggregate | platform-health | **ADDED-TO-YAML** `getPlatformHealth` |
| Incident CRUD / timeline | 本页 | **TODO-YAML** P2 |

## 字段级定义

### 查询字段（Incident list 目标 · **TODO-YAML**）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `status` | query | ❌ | open / mitigated / resolved |
| `severity` | query | ❌ | critical / warning / info 或 P0/P1/P2 映射 |
| `assignee` | query | ❌ | 负责人 user id |
| `tenant_id` | query | ❌ | 平台 RBAC 筛选 |
| `limit` / `cursor` | query | ❌ | 分页 |

### 请求字段 — 创建 Incident（目标 · **TODO-YAML**）

| 字段 | 必填 | 说明 |
|---|---|---|
| `idempotency_key` | ✅ | 幂等 |
| `title` | ✅ | 标题 |
| `severity` | ✅ | 严重级别 |
| `description` | ❌ | 描述 |
| `related_alert_ids[]` | ❌ | 来自 alert events |
| `related_task_ids[]` | ❌ | UUID |
| `assignee` | ❌ | 初始负责人 |

### 返回字段 — Incident（页面目标 · **TODO-YAML**）

| 字段 | 说明 |
|---|---|
| `incident_id` | 工单 ID |
| `title` | 标题 |
| `severity` | critical / warning / info |
| `status` | open / mitigated / resolved |
| `assignee` | 负责人 |
| `created_at` / `updated_at` / `resolved_at` | 时间戳 |
| `related_alerts[]` | `{ alert_incident_id, rule_id?, summary }` |
| `related_tasks[]` | `{ task_id, task_type, status }` |
| `health_context` | 可选；组件名 + HealthCheck.status 摘要 |

### 时间线字段 — timeline[]（目标 · **TODO-YAML**）

| 字段 | 说明 |
|---|---|
| `occurred_at` | date-time |
| `actor` | 操作人 |
| `action` | created / assigned / acknowledged / skill_triggered / mitigated / resolved / comment |
| `message` | 人类可读摘要 |
| `metadata` | 可选 JSON（task_id / skill_id 等） |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `duration_open` | now − created_at（未 resolve） |
| `severity_badge` | P0/P1/P2 映射（**仅 UI**） |
| `related_count` | 关联告警+任务数 |
| `sla_breach` | 是否超 SLA（**产品规则 · 待 YAML**） |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| API 未就绪 | 列表占位 + 深链入口 + TODO-YAML 横幅 | **不** 伪造工单数据 |
| 正常态（P2 后） | 列表 + 详情 + 时间线 | — |
| 无 Incident | 「暂无进行中故障」 | 引导新建或跳转 alert-rules |
| open | 默认排序靠前；红色/橙色 severity | — |
| mitigated | 黄色；仍显示在「进行中」Tab 可选 | 待产品 Tab 规则 |
| resolved | 移入已解决；显示 resolved_at | 复盘区可展开 |
| 关联告警不可用 | 显示「events API 待接入」 | 不得用 rules list 替代 |
| 写操作失败 | 表单错误 + request_id | 保留输入 |
| 无权限态 | 403 | 不渲染假工单 |
| 关联 task 404 | 时间线标注「任务不存在」 | GET task |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `Incident.status` | open / mitigated / resolved（**待 YAML**） | string |
| `severity` | 与 alert-rules API 对齐或独立 enum | 待 YAML |
| P0/P1/P2 | **仅 UI** 映射 | 不得改 API enum |
| `duration_open` | UI 计算 | humanize |
| `AsyncTask.status` | 关联任务用 YAML 六态 | 引用 job-history |

## 状态与能力口径

### Incident.status（目标 · **待 YAML**）

| 状态 | 含义 | 本页展示 |
|---|---|---|
| `open` | 新建/升级，未缓解 | 进行中 |
| `mitigated` | 已缓解，待复盘关闭 | 进行中或独立 Tab |
| `resolved` | 已关闭 | 已解决 |

状态迁移（目标）：open → mitigated → resolved；可分派/触发 Skill 不改变 status 除非 PATCH（待 YAML）。

### 与关联域状态区分

| 域 | 状态字段 | 禁止混用 |
|---|---|---|
| Incident | open/mitigated/resolved | — |
| Alert firing | firing/acknowledged/… | 不得写入 Incident API 为已冻结 |
| AsyncTask | pending/running/completed/… | 不得用 completed 表示 Incident resolved |
| HealthCheck | ok/degraded/error | 仅 health_context 摘要 |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 incident 读 RBAC | 已授权 | `403 FORBIDDEN` |
| 平台 incident 写 RBAC | 已授权（创建/分派/关闭） | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `idempotency_key` | 非空（创建/写操作） | `400 BAD_REQUEST`（待 YAML） |
| `title` / `severity` | 创建时必填 | `400`（待 YAML） |
| 关联 alert id 合法 | events 存在 | `404`（待 YAML） |
| 关联 task id 合法 | GET task 200 | `404` |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 查看 Incident list | ✅（待 YAML） | ✅ | ✅ | API P2 |
| 查看详情/时间线 | ✅（待 YAML） | ✅ | ✅ | — |
| 新建 Incident | ❌ | ✅（待 YAML） | ✅ | POST |
| 从告警升级 | ❌ | ✅（待 YAML） | ✅ | 依赖 events API |
| 分派/改派 | ❌ | ✅（待 YAML） | ✅ | PATCH |
| 触发 Skill | ❌ | ✅ → maint-skills | ✅ | execute 待 YAML |
| mitigated / resolve | ❌ | ✅（待 YAML） | ✅ | PATCH |
| 添加评论 | ❌ | ✅（待 YAML） | ✅ | timeline |
| 跳转 alert/job/health | ✅ | ✅ | ✅ | 只读深链 |
| 导出/复盘报告 | ❌ | Phase 2 | Phase 2 | — |

## 删除前置校验与当前契约边界

Phase 2 目标以 **关闭（resolved）** 为主，**不建议物理 DELETE** Incident。

若 YAML 合入 DELETE：

| 校验项 | 要求 | 失败响应 |
|---|---|---|
| `incident_id` 存在 | 路径合法 | `404` |
| write RBAC | 平台 incident write | `403` |
| status | 仅 `resolved` 可删（产品规则 · 待 YAML） | `422`（**待 YAML 合入后冻结**） |

合入前本文 **N/A**，不得写 DELETE 为已冻结。

## 接口冻结规则

### `GET /api/v1/tasks/{task_id}`（Core · 关联只读 · **非 Incident API**）

- `summary`：`查询异步任务状态`
- `success`：`200 + AsyncTask`
- `errors`：`404`
- 用途：校验 `related_task_ids[]`；**不是** Incident CRUD

### Incident API（待补 · **整域未冻结**）

<!-- TODO-YAML: GET/POST /api/v1/incidents, GET/PATCH /api/v1/incidents/{incident_id}, timeline sub-resource -->

| 目标能力 | 状态 |
|---|---|
| list incidents | **TODO-YAML** P2 |
| create incident | **TODO-YAML**；`idempotency_key` |
| get / patch incident | **TODO-YAML** |
| timeline append | **TODO-YAML** |
| link alert / task | **TODO-YAML** |

- 路径前缀须 Core `/api/v1/*`
- **不得** 使用 `/api/v1/boss/incidents`
- 合入前不得写入「已冻结路径」表

### 跳转数据源（已在他模块定义 · **不得重复冻结**）

| 模块 | 能力 | 状态 |
|---|---|---|
| alert-rules | alert events | **ADDED-TO-YAML** |
| job-history | task list | **ADDED-TO-YAML** |
| platform-health | health aggregate | **ADDED-TO-YAML** |
| maint-skills | skill execute | **TODO-YAML** P2 |

## 使用规则

- Phase 2 API 未就绪时，**禁止** 生产环境展示伪造 Incident 列表为真实数据
- **不得** 把 alert-rules CRUD 或 rules total 当作 Incident 数据源
- 创建/分派/关闭 **必须** 带 `idempotency_key`（YAML 合入后）
- 从告警升级须用 **alert event id**，不得仅用 `rule_id`
- 触发 Skill 须跳转 maint-skills 并携带 `incident_id` query，**不得** 在本页重复 execute API 定义
- 时间线条目 **应** 与 platform-audit-log 审计口径可对齐（audit API 待 YAML）
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- Incident CRUD + timeline API — **TODO-YAML** P2
- 与 alert events API 联动（升级/关联）— events **ADDED-TO-YAML**；Incident 写 API — P2
- 与 [`maint-skills.md`](maint-skills.md) execute 联动 — P2
- 与 [`../integration/enterprise-notification.md`](../integration/enterprise-notification.md) 分派通知 — integration 待 YAML
- SLA 计时 / 值班表 — Phase 3
- 复盘模板 / postmortem 导出 — Phase 3

## 响应示例

### 关联 AsyncTask 只读（Core · **路径已声明**）

```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "idempotency_key": "kb-index-20260616",
  "task_type": "kb.index",
  "status": "failed",
  "error_message": "vector store timeout",
  "created_at": "2026-06-16T09:00:00Z",
  "completed_at": "2026-06-16T09:12:00Z"
}
```

### Incident 详情目标响应（**待 YAML · 非已冻结**）

```json
{
  "incident_id": "inc-20260616-001",
  "title": "Model service degraded — vector store timeout",
  "severity": "critical",
  "status": "open",
  "assignee": "sre-oncall",
  "created_at": "2026-06-16T09:15:00Z",
  "updated_at": "2026-06-16T09:30:00Z",
  "resolved_at": null,
  "related_alerts": [
    { "alert_incident_id": "alert-ev-001", "summary": "KB index failure rate > 5%" }
  ],
  "related_tasks": [
    { "task_id": "770e8400-e29b-41d4-a716-446655440002", "task_type": "kb.index", "status": "failed" }
  ],
  "health_context": {
    "service_name": "model-service",
    "status": "degraded"
  },
  "timeline": [
    {
      "occurred_at": "2026-06-16T09:15:00Z",
      "actor": "sre-alice",
      "action": "created",
      "message": "Incident created from alert alert-ev-001"
    },
    {
      "occurred_at": "2026-06-16T09:20:00Z",
      "actor": "sre-alice",
      "action": "assigned",
      "message": "Assigned to sre-oncall"
    }
  ]
}
```

## 错误示例

### 无 platform incident 读权限（**TODO-YAML**）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-inc-403-001"
}
```

> **注**：适用于 **Incident list API（TODO-YAML）** 场景。`message` 为示例文案，**非** YAML 已冻结 scope 名。

### 创建缺少 idempotency_key（目标契约）

```json
{
  "code": "BAD_REQUEST",
  "message": "idempotency_key is required",
  "request_id": "req-boss-inc-400-001"
}
```

### Incident 不存在

```json
{
  "code": "NOT_FOUND",
  "message": "incident not found",
  "request_id": "req-boss-inc-404-001"
}
```

## 相关模块

- [alert-rules.md](alert-rules.md)、[maint-skills.md](maint-skills.md)、[job-history.md](job-history.md)、[platform-health.md](platform-health.md)
- [`../integration/enterprise-notification.md`](../integration/enterprise-notification.md)
- Console：[alerts-pending-items.md](../../console-modules/alerts/alerts-pending-items.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] Incident 整域 TODO-YAML；无伪造 REST 为已冻结
- [x] 与 alert/job/health/maint-skills 边界清晰；无重复定义 firing API
- [x] 含响应示例与错误示例（400 + 403 + 404）
- [x] 独立字段定义（查询/Incident/timeline/展示）
- [ ] Incident YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
