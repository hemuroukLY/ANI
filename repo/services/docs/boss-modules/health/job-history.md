# 任务历史

## 页面定位

`任务历史` 是 BOSS **平台可观测与健康** 域下的 **跨租户 AsyncTask 运维视图**，用于排查平台级异步作业（模型导入、KB 解析/索引、推理部署等）的失败、取消与 dead letter 队列。

本页属于 **Core / Tasks** 视角下的 **平台 RBAC** 页面，一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`（`AsyncTask` schema 与 `GET /api/v1/tasks/{task_id}`）。

Console [`async-task-center.md`](../../console-modules/alerts/async-task-center.md) 仅 **单租户 + 无 list API**（客户端聚合已知 `task_id`）；本页须 **平台 RBAC + 跨租户 list**，**不得** 逐租户轮询单查冒充平台 list。

## 文档管理规则

- 本文是 `任务历史` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-job-history.md`](../../tasks/modules/prd/boss/health/prd-boss-job-history.md) 与 [`spec-boss-job-history.md`](../../tasks/modules/spec/boss/health/spec-boss-job-history.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- AsyncTask 归属 **Core**；单查路径 `GET /api/v1/tasks/{task_id}` 已声明
- 平台 list `GET /api/v1/tasks` — **TODO-YAML**；不得把单查直接写成 BOSS list 正式契约
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id` 越权
- 写操作（重试/取消/标记处理）— **TODO-YAML**；合入后须 `idempotency_key`
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）；重试 API 待 YAML
- UI **不得** 使用 `succeeded`，须用 YAML 的 `completed`
- 禁止自造 task list schema / operationId / 路径为已冻结
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **全平台** 异步任务 list 与筛选（list API 待 YAML）
- 展示任务详情：`AsyncTask` 全字段只读（`GET /tasks/{task_id}` **路径已声明**）
- 支持按 `task_type`、`status`、租户、时间范围筛选（list **ADDED-TO-YAML**）
- 支持 dead letter 队列排查与高亮
- 为 SRE 提供重试/标记处理入口边界（写 API 待 YAML + `idempotency_key`）
- 与 Console 任务中心对照，明确 BOSS 跨租户运维差异

## 页面结构

- 首屏至少包含：`筛选区`、`任务表格`、`详情抽屉`、`操作区（待 YAML）`、`边界说明`
- 无数据态、无权限态、list API 未就绪态、查询失败态须可区分

```text
任务历史
├── 筛选（tenant_id / task_type / status / 时间 — list **ADDED-TO-YAML**）
├── 任务表格（id / type / status / progress / tenant / 创建时间）
├── 详情抽屉（AsyncTask 全字段）
│   └── 跳转关联资源（resource_type / resource_id）
└── 操作：重试 / 标记处理 / 查看 dead letter（写 API 待 YAML）
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | `GET /api/v1/tasks/{task_id}` | 单任务详情；**路径已声明**（无 operationId / 无 `x-ani-rbac-scope`） |
| Core | `GET /api/v1/tasks` list **TODO-YAML** | BOSS 正式跨租户 list |
| Services | `202 + AsyncTask` 触发源 | 任务来源参考；查询仍走 Core |

### 关键边界

- Services 层 `202` 响应引用同一 `AsyncTask` schema，但 **任务查询统一走 Core**
- 单查 `{task_id}` **不能** 直接作为 BOSS 跨租户 list 正式方案（禁止逐租户轮询）
- 重试须待 YAML 声明 + **`idempotency_key`**（与 ANI-14 写操作约定一致）
- **`succeeded` 不是合法 status** — 必须使用 `completed`

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 筛选区 | UI + list API **待 YAML** | task_type / status / tenant / 时间 | 刷新表格 |
| 任务表格 | list API **待 YAML** | 跨租户 AsyncTask 摘要 | 详情抽屉 |
| 详情抽屉 | `GET /tasks/{task_id}` | AsyncTask 全字段 | 关联资源 |
| dead letter 区 | list API filter | `status=dead_letter` | — |
| 操作区 | 写 API **待 YAML** | 重试/标记处理 | — |
| 边界说明 | 规划项 | list **ADDED-TO-YAML** | async-task-center |

## BOSS 与 Console 分工

| 维度 | BOSS 任务历史 | Console 任务中心 |
|---|---|---|
| 范围 | 全平台多租户 | 当前租户 |
| List | `GET /tasks` **待 YAML** | 会话内 task_id 客户端聚合 |
| 单查 | `GET /tasks/{task_id}` | 同 Core 路径 |
| 重试 | 平台运维（待 YAML + idempotency_key） | 只读为主 |
| dead letter | 平台排查 | 租户可见范围内 |
| RBAC | 平台 tasks 读/写（待 YAML） | 租户 JWT |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/tasks/{task_id}` | **未声明** | 返回 `AsyncTask`；`200` / `404` |

| 字段 | 冻结值 |
|---|---|
| `AsyncTask.status` | `pending` / `running` / `completed` / `failed` / `cancelled` / `dead_letter` |
| `AsyncTask.task_type` | `model.import` / `kb.parse` / `kb.index` / `inference.deploy` |
| `AsyncTask.resource_type` | `inference_service` / `kb_document` / `model_version`（可空） |

| 能力 | 状态 |
|---|---|
| GET | `/api/v1/tasks` | `listTasks` | 跨租户 list · **ADDED-TO-YAML** |
| 平台重试 task | **TODO-YAML** P2 P2 + `idempotency_key` |
| 平台取消 task | **TODO-YAML** |
| 标记处理/审计备注 | **TODO-YAML** |

## 字段级定义

### 查询字段（list API · **ADDED-TO-YAML**）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `tenant_id` | query | 可选 | 平台 RBAC 下筛选；**不得** 越权 |
| `task_type` | query | 可选 | YAML enum 四值 |
| `status` | query | 可选 | YAML enum 六值 |
| `start_time` | query | 可选 | 创建时间下限 |
| `end_time` | query | 可选 | 创建时间上限 |
| `limit` | query | 可选 | 分页条数 |
| `cursor` | query | 可选 | 游标分页 |

### 路径参数（单查 · **路径已声明**）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `task_id` | path | ✅ | UUID |

### 返回字段（AsyncTask · Core 已冻结）

| 字段 | 来源 | 说明 |
|---|---|---|
| `id` | `AsyncTask` | 任务 UUID |
| `idempotency_key` | `AsyncTask` | 幂等键 |
| `task_type` | `AsyncTask` | model.import / kb.parse / kb.index / inference.deploy |
| `resource_type` | `AsyncTask` | inference_service / kb_document / model_version；可空 |
| `resource_id` | `AsyncTask` | 关联资源 UUID；可空 |
| `status` | `AsyncTask` | 六态 enum；**禁止** succeeded |
| `attempt_count` | `AsyncTask` | 当前尝试次数 |
| `max_attempts` | `AsyncTask` | 最大尝试次数 |
| `progress_pct` | `AsyncTask` | 0–100 |
| `result` | `AsyncTask` | 结果对象；可空 |
| `error_message` | `AsyncTask` | 失败原因；可空 |
| `dead_letter_at` | `AsyncTask` | 进入 dead letter 时间；可空 |
| `created_at` | `AsyncTask` | 创建时间 |
| `completed_at` | `AsyncTask` | 完成时间；可空 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `tenant_id_display` | list API 返回；单查时从 task 上下文 — **ADDED-TO-YAML listTasks** |
| `duration_display` | completed_at − created_at |
| `is_terminal` | status ∈ {completed, failed, cancelled, dead_letter} |
| `can_retry_ui` | failed/dead_letter 且写 API 可用 — **待 YAML** |
| `task_type_label` | 中文展示名映射 |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 筛选 + 表格 + 可开详情 | list **ADDED-TO-YAML** |
| 无数据态 | 「当前条件下暂无任务」 | **不** 伪造行 |
| list API 未就绪 | 说明「跨租户 list 待 Core 合入」+ 可选单 task_id 查询 | 不得伪装 list |
| 单查 404 | 「任务不存在或无权查看」 | 不泄露跨租户存在性 |
| 无权限态 | 403 提示 | 平台 RBAC |
| failed / dead_letter | 高亮 + 展示 error_message | — |
| completed | 展示 completed_at；**禁止** 标签写 succeeded | YAML 枚举 |
| running | 展示 progress_pct 进度条 | 轮询 ≥ 2s |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `status` | YAML 六态；**never** succeeded | string |
| `task_type` | YAML 四值 enum | string |
| `progress_pct` | 0–100 整数 | integer |
| `attempt_count` / `max_attempts` | 后端原值 | integer |
| `created_at` / `completed_at` / `dead_letter_at` | ISO 8601 | date-time |
| `idempotency_key` | UUID 字符串 | string |

## 状态与能力口径

### AsyncTask.status 状态机（只读展示 + 待 YAML 写操作）

| 状态 | 含义 | 终态 |
|---|---|---|
| `pending` | 已创建，等待执行 | 否 |
| `running` | 执行中 | 否 |
| `completed` | 成功完成 | **是** |
| `failed` | 失败 | **是** |
| `cancelled` | 已取消 | **是** |
| `dead_letter` | 超过最大重试 | **是** |

**禁止** 在 UI/文档中使用 `succeeded` — 须用 `completed`。

### AsyncTask.task_type（冻结 enum）

| task_type | 典型触发 |
|---|---|
| `model.import` | 模型导入 |
| `kb.parse` | KB 文档解析 |
| `kb.index` | KB 索引构建 |
| `inference.deploy` | 推理服务部署 |

### 状态 × 操作可用性矩阵

| 操作 \ status | pending | running | completed | failed | cancelled | dead_letter |
|---|---|---|---|---|---|---|
| 查看详情 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| 列表筛选 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| 取消 | ⏳ 待 YAML | ⏳ 待 YAML | ❌ | ❌ | ❌ | ❌ |
| 重试 | ❌ | ❌ | ❌ | ⏳ 待 YAML | ❌ | ⏳ 待 YAML |
| 标记处理 | ❌ | ❌ | ❌ | ⏳ 待 YAML | ⏳ 待 YAML | ⏳ 待 YAML |
| 跳转关联资源 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

⏳ = 写 API 待 YAML 合入；合入后重试/标记处理须带 **`idempotency_key`**。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 tasks list 读 RBAC | 已授权（**TODO-YAML** `GET /tasks` scope 待 Core 定义） | `403 FORBIDDEN` |
| 单查 `GET /tasks/{task_id}` | 路径已声明；YAML **未** 声明 security / 401 / 403 | `404`（已声明）；401/403 为组件级推断 |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `task_id`（单查） | 合法 UUID | `404 NOT_FOUND` |
| list 筛选时间 | start < end（list **ADDED-TO-YAML**） | `400 BAD_REQUEST` |
| 重试（待 YAML） | status=failed/dead_letter + 写 RBAC + idempotency_key | `422` 待声明 |
| 取消（待 YAML） | status=pending/running | `422` 待声明 |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 列表（待 YAML） | ✅ | ✅ | ✅ | GET /tasks |
| 查看详情 | ✅ | ✅ | ✅ | GET /tasks/{id} |
| 筛选 dead letter | ✅ **待 YAML** | ✅ | ✅ | status 筛选 |
| 重试 | ❌ | ⏳ 待 YAML | ⏳ 待 YAML | 须 idempotency_key |
| 取消 | ❌ | ⏳ 待 YAML | ⏳ 待 YAML | — |
| 标记处理 | ❌ | ⏳ 待 YAML | ⏳ 待 YAML | 审计备注 |
| 导出 CSV | ❌ | Phase 2 | Phase 2 | **待 YAML** |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 操作 — **N/A**。任务取消为状态变更（待 YAML），非物理删除。

## 接口冻结规则

### `GET /api/v1/tasks/{task_id}`（Core · **路径已声明** · RBAC/平台 list **待 YAML**）

- `operationId`：权威源 **未声明**
- `summary`：`查询异步任务状态`
- `tags`：`["Tasks"]`
- `path.required`：`task_id`（UUID）
- `success`：`200 + AsyncTask`
- `errors`：`404`（YAML **已声明**）
- 认证/授权错误为网关通用推断，**未** 在本 operation responses 中冻结
- **不得** 自造 operationId
- **不得** 将本 path 等同于 BOSS 跨租户 list 或已定义平台 RBAC

### `GET /api/v1/tasks` list（待补）

<!-- TODO-YAML: GET /api/v1/tasks with cursor + tenant_id + status + task_type filter -->

- 须支持跨租户 list 与 cursor 分页
- 须平台 RBAC 鉴权
- query 预期：`tenant_id`、`task_type`、`status`、`start_time`、`end_time`、`limit`、`cursor`
- **合入前不得写入「已冻结」正文**

### 重试/取消（待补）

<!-- TODO-YAML: POST /api/v1/tasks/{task_id}/retry 或等价 -->

- 须声明 `idempotency_key`（request body）
- 须声明 `422 PreconditionFailed` 方可冻结前置条件错误
- 合入前不得写入已实现

## 使用规则

- **不得** 写 `GET /api/v1/tasks` list 为已实现（当前仅 `{task_id}` 单查已声明）
- UI status **必须** 使用 YAML 的 `completed`，**禁止** `succeeded`
- 重试/取消须待 YAML + `idempotency_key`
- list API 未上线时，**禁止** 生产环境逐租户轮询单查冒充平台 list
- 轮询刷新间隔 ≥ 2s，避免风暴
- `dead_letter_at` 非空时须高亮 dead letter 态

## 待补能力边界

- `GET /api/v1/tasks` — **ADDED-TO-YAML** P1 (`listTasks`)（cursor、tenant_id、status、task_type）
- 重试/取消 task — Phase 2 + `idempotency_key`
- 标记处理/审计备注 — Phase 2
- 平台 CSV 导出 — Phase 2
- Phase 3 扩展 task_type（见 Console async-task-center 规划项）— 须 YAML 合入后回写

## 响应示例

### 单任务查询成功（Core · **路径已声明**）

```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "idempotency_key": "880e8400-e29b-41d4-a716-446655440003",
  "task_type": "inference.deploy",
  "resource_type": "inference_service",
  "resource_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "failed",
  "attempt_count": 3,
  "max_attempts": 3,
  "progress_pct": 65,
  "result": null,
  "error_message": "GPU allocation failed: InsufficientGPU",
  "dead_letter_at": "2026-06-16T10:30:00Z",
  "created_at": "2026-06-16T10:00:00Z",
  "completed_at": "2026-06-16T10:30:00Z"
}
```

### 任务 list 成功（页面目标 · **待 YAML 合入后对齐**）

```json
{
  "items": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "task_type": "kb.parse",
      "status": "running",
      "progress_pct": 40,
      "tenant_id": "t-001",
      "created_at": "2026-06-16T09:00:00Z"
    }
  ],
  "next_cursor": "cursor-tasks-page-2",
  "total": 128
}
```

## 错误示例

### task_id 格式非法

```json
{
  "code": "BAD_REQUEST",
  "message": "task_id must be a valid UUID",
  "request_id": "req-boss-task-400-001"
}
```

### 无平台 tasks list 读权限（**TODO-YAML** · 非单查 path）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-task-403-001"
}
```

> **注**：适用于 **`GET /api/v1/tasks` list（TODO-YAML）** 场景。`message` 为示例文案，**非** YAML 已冻结 scope 名。单查 `GET /tasks/{task_id}` 的 YAML 仅声明 `404`，**不** 应用本示例替代单查鉴权语义。

## 相关模块

- Console：[async-task-center.md](../../console-modules/alerts/async-task-center.md)
- [alert-rules.md](alert-rules.md)、[maint-skills.md](maint-skills.md)
- [`inference-monitoring.md`](inference-monitoring.md)、[`kb-monitoring.md`](kb-monitoring.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] `AsyncTask.status` 六态与 `v1.yaml` 一致（completed 非 succeeded）
- [x] `task_type` 四值 enum 与 YAML 一致
- [x] 含状态 × 操作矩阵
- [x] 重试须 idempotency_key 已声明
- [x] 含响应示例与错误示例（400 + 403）
- [x] 独立字段定义（查询/返回/展示）
- [ ] list/retry YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
