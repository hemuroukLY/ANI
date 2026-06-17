# 告警与待处理事项

## 页面定位

`告警与待处理事项` 是 `Console` 首页「平台概览」第五主题区块的明细页，也是租户侧统一查看**高优先级风险摘要**与**待人工处理事项**的入口。

本页属于 **Console UI 聚合能力**，数据可同时来自 `Core` 与 `Services`，但不得自造独立告警资源域契约。

## 文档管理规则

- 本文是 `告警与待处理事项` 的主维护文档
- `tasks/modules/prd/console/alerts/prd-console-alerts-pending-items.md` 与 `tasks/modules/spec/console/alerts/spec-console-alerts-pending-items.md` 为辅助材料，冲突时以本文为准
- 一级权威源：`ANI-main/repo/api/openapi/v1.yaml`、`ANI-main/repo/api/openapi/services/v1.yaml`
- 当前 **无** 独立的 `GET /api/v1/alerts*` 或 `GET /api/v1/pending-items*` 冻结路径
- 正文不得把客户端聚合策略写成已存在的后端只读契约
- Handler 实现状态以 `CORE-HANDLER-IMPLEMENTATION-GUIDE` 为准；OpenAPI 已声明 ≠ 已实现

## Core 层要求

- Core 资源类接口必须使用 `/api/v1/*`
- 当前可引用的 Core 只读能力（**不等于告警实例 API**）：
  - `GET /api/v1/tasks/{task_id}` — 单任务状态查询（`AsyncTask`）
  - `GET /api/v1/observability/alert-rules*` — 告警**规则** CRUD，非 firing 告警事件
  - `GET /api/v1/observability/query` — PromQL 代理，可用于指标异常判断
- 不允许继续使用旧的 `/api/v1/console/home/alerts` 路径
- 页面不要求前端显式传 `tenant_id`；租户边界由认证上下文获取
- 错误结构统一为 `{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`

<!-- TODO-YAML: 待 Core/Services 冻结只读「告警事件 / 待办聚合」契约；当前无 listAlerts / listPendingItems -->

## Services 层要求

- Services 业务类接口必须使用 `/api/v1/svc/*`
- 当前可引用的 Services 事实来源（用于**推断**待处理事项，非专用告警 API）：
  - `GET /api/v1/svc/inference-services` — 筛选 `status=failed` 等服务异常
  - `GET /api/v1/svc/knowledge-bases/{kb_id}/documents` — 文档解析失败态（若 `KBDocument` 暴露 error）
  - Webhook 投递失败摘要见租户管理 `WebhookDelivery`（若页面有权限）
- 聚合仅发生在 Console 页面层，不写入 `services/v1.yaml` 新资源

## 页面职责

- 汇总展示高优先级告警数、处理中任务数、失败任务数、待确认事项数
- 提供按来源模块跳转的明细入口
- 支持区块级刷新与局部失败态
- **不承担**完整告警规则编辑、批量处置、审计导出（后者归属 observability / 各业务模块）

## 页面结构

```text
告警与待处理事项
├── 顶部摘要
│   ├── 高优先级告警数
│   ├── 处理中任务数
│   ├── 失败任务数
│   └── 待确认事项数
├── 分类列表（按来源）
│   ├── Core 异步任务（已知 task_id 时）
│   ├── 推理服务异常
│   ├── 知识库 / 文档异常
│   └── 安全 / Webhook 待确认（有权限时）
└── 跳转入口
    ├── 任务中心（async-task-center）
    ├── 推理服务
    ├── 知识库管理
    └── 安全与身份概览
```

## 数据来源与分层约束

### 当前可行策略（Phase P0 文档口径）

| 摘要字段 | 建议来源 | 说明 |
|---|---|---|
| 高优先级告警数 | **待补** 或 `observability/query` + 规则 | 无 firing 告警 list API；不得伪造数值 |
| 处理中任务数 | 已知 `task_id` 轮询 `GET /tasks/{task_id}` | **无** list tasks；首页只能展示最近触达任务 |
| 失败任务数 | 同上 + 各模块 `202` 返回的 `AsyncTask.status=failed` | 依赖用户会话内任务上下文 |
| 待确认事项数 | 各业务模块筛选 + 租户 Webhook 失败 | 口径待产品冻结 |

### 关键边界

- `observability/alert-rules` 管理的是**规则**，不是**告警事件**
- 不得把 `alert-rules` 列表长度误写成「当前告警数」
- 聚合仅发生在 Console 页面层，不写入 `services/v1.yaml` 新资源

## 字段级定义

| 字段 | 说明 | 展示建议 |
|---|---|---|
| 高优先级告警数 | 当前需立即处理的告警摘要 | 整数；无 API 时显示「待接入」 |
| 处理中任务数 | `AsyncTask.status` 为 `pending` 或 `running` 的任务 | 整数；需说明统计范围 |
| 失败任务数 | `AsyncTask.status` 为 `failed` 或 `dead_letter` | 大于 0 高亮 |
| 待确认事项数 | 待审批 / 待确认 / Webhook 失败等 | 整数；分项可展开 |
| 最近更新时间 | 本页聚合刷新时间 | 绝对时间或相对时间 |

## 与平台概览的关系

- 平台概览 (`platform-overview.md`) 第五区块为本页的**摘要视图**
- 本页为同一主题的**明细与跳转页**，不重复定义首页专属 schema
- 从首页「查看详情」应进入本页或任务中心

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 用户登录 | 已认证 | `401 UNAUTHORIZED` |
| 租户读权限 | 具备各来源模块读 scope | `403 FORBIDDEN` |
| 单来源接口 | 可用 | 本页局部降级，不整页失败 |

本页无 POST/PUT 写操作，不涉及 `idempotency_key`。

说明：各来源模块的前置条件（如推理服务 `404`）在跳转目标页校验；本页聚合层只读展示。

## 操作可用性矩阵

| 操作 | 只读用户 | 租户管理员 |
|---|---|---|
| 查看摘要与列表 | ✅ | ✅ |
| 刷新单来源区块 | ✅ | ✅ |
| 跳转来源模块处置 | ✅（有模块权限时） | ✅ |
| 在本页直接关闭告警 | ❌ | ❌（待专用 API 冻结） |
| 编辑告警规则 | ❌ | ❌（跳转 observability） |

## 接口冻结规则

本页**无独立冻结 API**。以下为聚合时可引用的已冻结 operation：

### `GET /api/v1/tasks/{task_id}`

- 成功：`200 + AsyncTask`
- 错误：`404 NOT_FOUND`
- **当前 YAML 未声明** `401`/`403`（仅 404）；租户边界由 handler 实现

### `GET /api/v1/observability/alert-rules`

- 成功：`200 + ObservabilityAlertRuleListResponse`
- 错误：`401`、`403`
- RBAC：`scope:observability:read`
- **用途限制**：规则管理，非 firing 告警事件列表

### `GET /api/v1/observability/query`

- 成功：`200 + ObservabilityQueryResponse`
- 错误：`400`、`401`、`403`
- Query 必填：`query`（PromQL）

### `GET /api/v1/svc/inference-services`

- 成功：`200` + items 数组（`InferenceService`）
- 错误：`401`（YAML 已声明）
- 页面层按 `status=failed` 计数

## 待补边界

- `listAlerts` / `listPendingItems` 或等价只读聚合 API — **TODO-YAML**
- firing 告警事件与规则分离的产品口径 — 待 observability 团队冻结
- 待确认事项（审批 / Webhook 失败）统一枚举 — 待产品定稿
- 全量租户任务视图 — 依赖 `GET /api/v1/tasks` list（见 async-task-center）

## 相关模块

- `platform-overview.md` — 首页第五区块摘要
- `async-task-center.md` — 异步任务明细
- `inference-service.md` — 推理服务异常来源
- `knowledge-base.md` — 知识库文档异常来源

## 验收标准

- [ ] 文档未自造 `Alert` / `PendingItem` schema 名称
- [ ] 无 API 处明确标注 `TODO-YAML`，不展示伪造 0 值
- [ ] 与 `platform-overview.md` 第五区块字段口径一致
- [ ] 接口冻结规则逐 operation 列出成功码与错误码
- [ ] HTML 摘要指向本文
