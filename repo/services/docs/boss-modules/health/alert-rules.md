# 告警规则

## 页面定位

`告警规则` 是 BOSS **运维与可观测** 域下的 **平台告警治理** 页面，包含 **两类能力**（同一导航下分 Tab，不拆独立 REST 域）：

1. **告警规则管理** — 基于 Core `ObservabilityAlertRule` CRUD（PromQL 规则）
2. **活动告警处置** — P0/P1/P2 firing 事件确认/静默/分派/关闭（**无**冻结 API，**TODO-YAML**）

本页属于 **Core / Observability** 视角下的 **平台 RBAC** 页面，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`。

Console [`alerts-pending-items.md`](../../console-modules/alerts/alerts-pending-items.md) 只做租户侧摘要聚合，**不承担** 规则 CRUD 与平台 firing 处置。本页 **不得** 把 `alert-rules` list 长度误写成「当前告警数」。

## 文档管理规则

- 本文是 `告警规则` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-alert-rules.md`](../../tasks/modules/prd/boss/health/prd-boss-alert-rules.md) 与 [`spec-boss-alert-rules.md`](../../tasks/modules/spec/boss/health/spec-boss-alert-rules.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

### 已冻结 — 告警规则 CRUD

| 方法 | 路径 | operationId | 成功 |
|---|---|---|---|
| GET | `/api/v1/observability/alert-rules` | `listObservabilityAlertRules` | `200` |
| POST | `/api/v1/observability/alert-rules` | `createObservabilityAlertRule` | `201` |
| GET | `/api/v1/observability/alert-rules/{rule_id}` | `getObservabilityAlertRule` | `200` |
| PATCH | `/api/v1/observability/alert-rules/{rule_id}` | `updateObservabilityAlertRule` | `200` |
| DELETE | `/api/v1/observability/alert-rules/{rule_id}` | `deleteObservabilityAlertRule` | `200` |

- RBAC：`scope:observability:read` / `scope:observability:write`
- POST/PATCH 必填 `idempotency_key`（`CreateObservabilityAlertRuleRequest` / `UpdateObservabilityAlertRuleRequest`）
- Schema：`ObservabilityAlertRule` 含 `tenant_id` — 规则在 YAML 中为 **租户归属**；BOSS 平台视图须 **平台 RBAC 跨租户 list**，**不得** 信任前端伪造 `tenant_id`

### 未冻结 — 活动告警事件

<!-- TODO-YAML: listAlertIncidents / acknowledge / silence / assign / resolve -->

| 能力 | 状态 |
|---|---|
| firing 告警 list | **TODO-YAML** |
| 确认/静默/分派/关闭 | **TODO-YAML** |

**禁止** 把 `alert-rules` CRUD 写成 firing 事件 API；**禁止** 把规则 `state=active` 等同于「当前 firing」。

## 页面职责

- **Tab 规则管理**：平台视角 list/create/edit/delete `ObservabilityAlertRule`；启用/禁用规则
- **Tab 活动告警**：展示 firing 事件 list、筛选、详情抽屉与处置操作（API 待 YAML）
- 支持从活动告警跳转日志（[`platform-logs.md`](platform-logs.md)）、指标（PromQL）、运维 Skill（[`maint-skills.md`](maint-skills.md)）
- 为 [`platform-health.md`](platform-health.md)、[`gpu-monitoring.md`](gpu-monitoring.md) 等监控页提供告警创建入口
- **不承担** Console 租户侧待办摘要（见 alerts-pending-items）

## 页面结构

- 首屏 Tab 切换：`规则管理` / `活动告警`
- 规则 Tab：list + 新建/编辑表单 + 删除确认
- 活动 Tab：筛选 + 表格 + 详情抽屉 + 处置按钮（待 YAML）
- 无数据态、无权限态、写操作失败态、活动 API 未就绪态须可区分

```text
告警规则
├── Tab：规则管理
│   ├── 规则列表（ObservabilityAlertRule）
│   ├── 新建/编辑（promql / severity / duration / labels）
│   ├── 启用/禁用（enabled + state）
│   └── 删除（DELETE → state=deleted）
└── Tab：活动告警
    ├── 筛选（severity / tenant / status）
    ├── 告警表格（TODO-YAML）
    ├── 详情抽屉（指标值 / 规则 / 影响范围）
    └── 操作（确认/静默/分派/关闭/触发 Skill — TODO-YAML）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | `/api/v1/observability/alert-rules*` | 规则 CRUD — **已冻结** |
| Core | `GET /api/v1/observability/query` | 活动告警详情中的指标值；**非** firing list |
| Core | 活动告警 events API **TODO-YAML** | firing list / ack / silence |
| Console | alerts-pending 客户端聚合 | **非** BOSS 正式数据源 |

### 关键边界

- `observability/alert-rules` 管理 **规则定义**，不是 **告警实例/firing 事件**
- `ObservabilityAlertRule.state`（`active`/`disabled`/`deleted`）是 **规则生命周期**，不是 firing `status`
- `severity` API 枚举为 `info`/`warning`/`critical`；UI 可映射 P0/P1/P2 **展示**，**不得** 改 API 枚举
- BOSS 跨租户 rules list 筛选 query — **待 RBAC/YAML 设计**；不得逐租户 JWT 轮询作为正式方案
- **`CreateObservabilityAlertRuleRequest` 不含 `tenant_id`**：创建时 `ObservabilityAlertRule.tenant_id` 由 **JWT claims** 写入；BOSS 为 **其他租户** 创建/编辑规则须 **platform API 或 TODO-YAML**，不得切换租户 JWT 轮询
- PromQL 验证可在创建规则前通过 `observability/query` 探测，**不等于** 规则已触发

## 页面区块与数据来源映射

| 区块 | Tab | 主要来源 | 说明 |
|---|---|---|---|
| 规则 list | 规则管理 | `listObservabilityAlertRules` | cursor 分页 |
| 新建/编辑表单 | 规则管理 | POST/PATCH alert-rules | `idempotency_key` 必填 |
| 删除确认 | 规则管理 | DELETE alert-rules | 返回 `state=deleted` |
| 活动告警表格 | 活动告警 | events API **待 YAML** | 不得用 rules list 替代 |
| 详情抽屉 | 活动告警 | events + PromQL | 指标值来自 query |
| 处置按钮 | 活动告警 | ack/silence API **待 YAML** | Phase 1 可 disabled |
| 跳转 Skill | 活动告警 | maint-skills **待 YAML** | Phase 2 |

## BOSS 与 Console 分工

| 能力 | BOSS 告警规则 | Console |
|---|---|---|
| 规则 CRUD | ✅ 平台视角 | ❌ |
| 活动告警处置 | ✅ 平台 SRE（API 待 YAML） | 摘要 only |
| 租户告警摘要 | 可选排行 | ✅ alerts-pending |
| firing list API | **TODO-YAML** | 客户端聚合 |
| RBAC | 平台 observability read/write | 租户 JWT |

## 当前冻结事实

| 方法 | 路径 | operationId | RBAC |
|---|---|---|---|
| GET | `/api/v1/observability/alert-rules` | `listObservabilityAlertRules` | read |
| POST | `/api/v1/observability/alert-rules` | `createObservabilityAlertRule` | write |
| GET | `/api/v1/observability/alert-rules/{rule_id}` | `getObservabilityAlertRule` | read |
| PATCH | `/api/v1/observability/alert-rules/{rule_id}` | `updateObservabilityAlertRule` | write |
| DELETE | `/api/v1/observability/alert-rules/{rule_id}` | `deleteObservabilityAlertRule` | write |

| 字段 | 冻结值 |
|---|---|
| `ObservabilityAlertRule.severity` | `info` / `warning` / `critical` |
| `ObservabilityAlertRule.state` | `active` / `disabled` / `deleted` |
| `CreateObservabilityAlertRuleRequest.required` | `idempotency_key`, `name`, `promql`, `severity` |

| 能力 | 状态 |
|---|---|
| firing events list/ack/silence | **ADDED-TO-YAML** `/observability/alert-events` |
| BOSS 跨租户 rules 筛选 query | 待 RBAC 设计 |

## 字段级定义

### 查询字段 — `listObservabilityAlertRules`（YAML 已冻结）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `limit` | query | ❌ | 1–100 |
| `cursor` | query | ❌ | 分页游标 |

<!-- TODO-YAML: tenant_id / severity / enabled 等平台筛选 -->

### 请求字段 — `CreateObservabilityAlertRuleRequest`（YAML 已冻结）

| 字段 | 必填 | 说明 |
|---|---|---|
| `idempotency_key` | ✅ | 1–128 字符 |
| `name` | ✅ | 1–128 字符 |
| `promql` | ✅ | minLength 1 |
| `severity` | ✅ | info / warning / critical |
| `duration` | ❌ | 默认 `5m` |
| `labels` | ❌ | string map |
| `annotations` | ❌ | string map |
| `enabled` | ❌ | 默认 true |

### 请求字段 — `UpdateObservabilityAlertRuleRequest`（YAML 已冻结）

| 字段 | 必填 | 说明 |
|---|---|---|
| `idempotency_key` | ✅ | 幂等键 |
| `name` / `promql` / `severity` / `duration` / `labels` / `annotations` / `enabled` | ❌ | 部分更新 |

### 返回字段 — `ObservabilityAlertRule`（YAML 已冻结）

| 字段 | 类型/枚举 | 说明 |
|---|---|---|
| `id` | string | 规则 ID |
| `tenant_id` | string | 租户归属 |
| `name` | string | 规则名 |
| `promql` | string | PromQL 表达式 |
| `duration` | string | 持续阈值；默认 5m |
| `severity` | info / warning / critical | 严重级别 |
| `labels` / `annotations` | object | 可选 |
| `enabled` | boolean | 是否启用 |
| `state` | active / disabled / deleted | 规则状态 |
| `created_at` / `updated_at` | date-time | 时间戳 |

### 活动告警字段（页面目标 · **TODO-YAML** · 非冻结 schema）

| 字段 | 说明 |
|---|---|
| `incident_id` | 事件 ID |
| `severity` | 展示可映射 P0/P1/P2；API 待 YAML |
| `alert_name` | 告警名称 |
| `target` | 对象 |
| `tenant_id` | 租户 |
| `status` | firing / acknowledged / silenced / resolved（**待 YAML**） |
| `fired_at` | 触发时间 |
| `duration` | 持续时间 |
| `assignee` | 负责人 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `severity_display` | P0↔critical、P1↔warning、P2↔info（**仅 UI**） |
| `rule_count_by_severity` | 规则 Tab 统计 |
| `firing_count` | 活动 Tab；**须** 来自 events API，不可用 rules total |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 规则 list 正常 | 表格 + 分页 | cursor 加载 |
| 规则无数据 | 「暂无告警规则」 | 引导新建 |
| 活动告警 API 未就绪 | Tab 内说明 TODO-YAML；按钮 disabled | **不** 用 rules list 冒充 |
| 创建/编辑成功 | 201/200 + 刷新 list | 展示返回体 |
| 删除成功 | 200；行标记 deleted 或移除 | 遵循 YAML 返回 |
| 写操作失败 | 表单内错误 + `request_id` | 保留输入 |
| 无读权限 | 403；两 Tab 均不可见数据 | — |
| 无写权限 | 规则 Tab 只读；隐藏新建/编辑/删除 | SRE 可写 |
| `state=deleted` | 默认不展示或灰显 | 软删语义 |
| `enabled=false` | 禁用标签 | 与 state=disabled 区分展示 |
| critical 规则/告警 | 红色高亮 | 优先排序 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `severity` | YAML 三态 | API 必须用 info/warning/critical |
| `duration` | PromQL 风格 duration | 如 `5m`、`1h` |
| `state` | 规则生命周期三态 | 非 firing status |
| P0/P1/P2 | **仅 UI 映射** | 不得写入 API enum |
| `fired_at` | ISO 8601 | date-time（events 待 YAML） |

## 状态与能力口径

### ObservabilityAlertRule.state（规则生命周期）

| 状态 | 含义 | 本页展示 |
|---|---|---|
| `active` | 规则生效中 | 正常行 |
| `disabled` | 规则停用 | 灰显/禁用标签 |
| `deleted` | 已删除（DELETE 返回） | 默认隐藏 |

### enabled vs state

- `enabled=false` 可在 PATCH 中设置；与 `state=disabled` 可能并存 — UI 以 `state` 为主、`enabled` 为辅展示
- DELETE 后 `state=deleted`；**不得** 物理删除后仍展示为 active

### 活动告警 status（**待 YAML · 页面目标**）

| 状态 | 含义 |
|---|---|
| `firing` | 触发中 |
| `acknowledged` | 已确认 |
| `silenced` | 已静默 |
| `resolved` | 已恢复 |

**禁止** 把上述 firing 状态机写成已冻结；合入 YAML 前仅作 UI 目标。

## 创建前置条件

| 场景 | 依赖 | 未满足响应 |
|---|---|---|
| 列出规则 | observability read RBAC | `403` |
| 创建规则 | name/promql/severity + `idempotency_key` + write RBAC | `400` / `403` |
| 更新规则 | rule 存在 + `idempotency_key` | `404` / `400` |
| 删除规则 | rule 存在 + write RBAC | `404` / `403` |
| 查看活动告警 | events API **待 YAML** + read RBAC | — |
| 处置活动告警 | events write API **待 YAML** | — |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 列出规则 | ✅ | ✅ | ✅ | GET alert-rules |
| 查看规则详情 | ✅ | ✅ | ✅ | GET by id |
| 创建规则 | ❌ | ✅ | ✅ | POST + idempotency_key |
| 编辑规则 | ❌ | ✅ | ✅ | PATCH |
| 删除规则 | ❌ | ✅ | ✅ | DELETE |
| 查看活动告警 | ✅ | ✅ | ✅ | TODO-YAML |
| 确认/静默/分派/关闭 | ❌ | ✅ | ✅ | TODO-YAML |
| 触发运维 Skill | ❌ | ✅ | ✅ | maint-skills Phase 2 |
| PromQL 探测 | ✅ | ✅ | ✅ | observability/query |

## 删除前置校验与当前契约边界

### DELETE `/api/v1/observability/alert-rules/{rule_id}`

| 校验项 | 要求 | 失败响应 |
|---|---|---|
| `rule_id` 存在 | 路径参数合法 UUID/ID | `404 NOT_FOUND` |
| write RBAC | `scope:observability:write` | `403 FORBIDDEN` |
| 幂等 | 重复 DELETE 已删规则 | 遵循 YAML（通常仍 200 或 404） |

- DELETE **成功**返回 `200 + ObservabilityAlertRule`（`state=deleted`）
- UI 删除前须二次确认；**不得** 假设物理删除
- 活动告警 Tab **无 DELETE** 规则操作 — 处置 API 待 YAML

## 接口冻结规则

### `GET /api/v1/observability/alert-rules`

- `operationId`：`listObservabilityAlertRules`
- `x-ani-rbac-scope`：`scope:observability:read`
- `query`：`limit`（1–100）、`cursor`
- `success`：`200 + ObservabilityAlertRuleListResponse`（`items`, `total`, `next_cursor`）
- `errors`：`401`、`403`

### `POST /api/v1/observability/alert-rules`

- `operationId`：`createObservabilityAlertRule`
- `x-ani-rbac-scope`：`scope:observability:write`
- `body.required`：`idempotency_key`, `name`, `promql`, `severity`
- `success`：`201 + ObservabilityAlertRule`
- `errors`：`400`、`401`、`403`

### `PATCH /api/v1/observability/alert-rules/{rule_id}`

- `operationId`：`updateObservabilityAlertRule`
- `body.required`：`idempotency_key`
- `success`：`200 + ObservabilityAlertRule`
- `errors`：`400`、`401`、`403`、`404`

### `DELETE /api/v1/observability/alert-rules/{rule_id}`

- `operationId`：`deleteObservabilityAlertRule`
- `success`：`200 + ObservabilityAlertRule`（`state=deleted`）
- `errors`：`401`、`403`、`404`

### 活动告警 events（待补）

<!-- TODO-YAML: GET /api/v1/observability/alert-incidents 等 -->

- **不得** 写入「已冻结」正文
- **不得** 用 rules CRUD path 冒充 incidents API

## 使用规则

- **不得** 把 `alert-rules` items.length 显示为「当前告警数」
- **不得** 把 `ObservabilityAlertRule.state=active` 等同于 firing
- 创建/更新 **必须** 带 `idempotency_key`；UI 重试须复用同一 key
- severity 提交 **必须** 使用 YAML 枚举；P0/P1/P2 仅作 UI 标签
- BOSS 跨租户 list 未冻结前，**禁止** 生产环境逐租户 JWT 轮询
- BOSS **代建/代改** 其他租户规则须待 platform API；当前 CRUD 为 **租户 JWT 上下文** 契约
- 活动 Tab API 未就绪时，处置按钮 **必须** disabled 并标注 TODO-YAML
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 活动告警 events list/ack/silence — **ADDED-TO-YAML**（assign/resolve — P2）
- BOSS 跨租户 rules 筛选（tenant_id/severity/enabled）— 待 RBAC/query
- 与 [`maint-skills.md`](maint-skills.md) 联动触发 remediation — Phase 2
- 告警审计导出 — Phase 2
- 批量导入规则 — Phase 3

## 响应示例

### 规则 list 成功（Core · 已冻结）

```json
{
  "items": [
    {
      "id": "rule-001",
      "tenant_id": "t-platform",
      "name": "gpu-utilization-critical",
      "promql": "avg(gpu_utilization_percent) > 95",
      "duration": "5m",
      "severity": "critical",
      "labels": { "team": "sre" },
      "annotations": { "summary": "GPU utilization above 95%" },
      "enabled": true,
      "state": "active",
      "created_at": "2026-06-10T08:00:00Z",
      "updated_at": "2026-06-15T10:00:00Z"
    }
  ],
  "total": 1,
  "next_cursor": null
}
```

### 创建规则成功（Core · 已冻结）

```json
{
  "id": "rule-002",
  "tenant_id": "t-platform",
  "name": "inference-error-rate-warning",
  "promql": "sum(rate(inference_errors_total[5m])) / sum(rate(inference_requests_total[5m])) > 0.05",
  "duration": "10m",
  "severity": "warning",
  "enabled": true,
  "state": "active",
  "created_at": "2026-06-16T12:00:00Z",
  "updated_at": "2026-06-16T12:00:00Z"
}
```

### 删除规则成功（Core · 已冻结）

```json
{
  "id": "rule-001",
  "tenant_id": "t-platform",
  "name": "gpu-utilization-critical",
  "promql": "avg(gpu_utilization_percent) > 95",
  "duration": "5m",
  "severity": "critical",
  "enabled": false,
  "state": "deleted",
  "created_at": "2026-06-10T08:00:00Z",
  "updated_at": "2026-06-16T14:00:00Z"
}
```

## 错误示例

### 创建规则缺少 idempotency_key

```json
{
  "code": "BAD_REQUEST",
  "message": "idempotency_key is required",
  "request_id": "req-boss-alert-400-001"
}
```

### 无 observability 写权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: scope:observability:write required",
  "request_id": "req-boss-alert-403-001"
}
```

### 规则不存在

```json
{
  "code": "NOT_FOUND",
  "message": "alert rule not found",
  "request_id": "req-boss-alert-404-001"
}
```

## 相关模块

- [platform-health.md](platform-health.md)、[gpu-monitoring.md](gpu-monitoring.md)、[platform-logs.md](platform-logs.md)
- [maint-skills.md](maint-skills.md)
- Console：[alerts-pending-items.md](../../console-modules/alerts/alerts-pending-items.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] 规则 CRUD 五 path 与 `v1.yaml` operationId 一致
- [x] 规则 CRUD 与 firing 事件严格分离；无伪造 incidents API
- [x] `severity`/`state` 枚举与 YAML 一致；P0/P1/P2 仅 UI
- [x] DELETE 前置校验已写；含响应示例与错误示例（400 + 403 + 404）
- [x] POST/PATCH `idempotency_key` 要求已冻结
- [ ] 活动告警 events YAML 合入后回写 Tab 2 正式契约
- [x] PRD/SPEC/HTML 与本文同步
