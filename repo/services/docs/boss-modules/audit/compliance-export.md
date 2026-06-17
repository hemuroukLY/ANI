# 合规导出与取证

## 页面定位

`合规导出与取证` 是 BOSS **安全审计与合规** 域下的 **跨租户合规数据包生成与异步导出** 专页：发起全平台或指定租户的审计/计量/配置取证导出，跟踪异步作业状态，并提供下载与留存边界说明。

本页属于 **平台 RBAC** 视角。合规包 **创建/导出写 API** — **TODO-YAML** P2；异步作业状态跟踪可引用 Core `GET /api/v1/tasks/{task_id}`（**路径已声明**，`operationId` **未声明**）。Console [`compliance.md`](../../console-modules/security/compliance.md) 为 **单租户只读** 合规摘要与历史报告（Services 草案路径 — **待 YAML，非冻结**）；完整包生成与平台取证 **归 BOSS**。

## 文档管理规则

- 本文是 `合规导出与取证` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-compliance-export.md`](../../tasks/modules/prd/boss/audit/prd-boss-compliance-export.md) 与 [`spec-boss-compliance-export.md`](../../tasks/modules/spec/boss/audit/spec-boss-compliance-export.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)
- Console 对照：[`compliance.md`](../../console-modules/security/compliance.md)（草案见 `openapi-phase3-compliance-draft.md`，**非冻结**）

## Core 层要求

- 合规导出写 API — **TODO-YAML** P2；合入后 POST 须 `idempotency_key`
- 异步任务查询：`GET /api/v1/tasks/{task_id}` — **路径已声明**；返回 `AsyncTask`
- 平台跨租户导出 job list — **TODO-YAML** P2（不可仅依赖单查轮询）
- Console Services 草案：`GET /api/v1/svc/compliance/*` — **草案/待 YAML**；Console 只读，**非** BOSS 写入口
- 平台 RBAC 鉴权；**不得** 信任未授权 `tenant_id` 导出他租户数据
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422` 仅 YAML 已声明 operation 可写冻结（§2.10）
- `AsyncTask.status` 须用 `completed`，**禁止** `succeeded`
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 提供 **合规导出任务创建** 表单（framework、租户范围、时间窗、包含数据集 — 写 API 待 YAML）
- 提供 **全平台导出作业 list** 与状态跟踪（list 待 YAML；单查 `GET /tasks/{task_id}` 已声明）
- 展示 `AsyncTask` 进度：`status`、`progress_pct`、`error_message`、`result` 下载链接
- 支持跳转 **平台审计日志**、**任务历史**（[`job-history.md`](../health/job-history.md)）
- 与 Console 合规摘要对照，明确 BOSS 独占写/导出能力

## 页面结构

- 首屏至少包含：`创建导出区`、`导出作业表格`、`任务详情抽屉`、`下载区`、`边界说明`
- 无数据态、无权限态、写 API 未就绪态、任务失败/dead_letter 态须可区分

```text
合规导出与取证
├── 创建导出（framework / tenant_scope / 时间窗 — POST 待 YAML）
├── 导出作业表格（task_id / status / progress / 创建时间）
├── 详情抽屉（AsyncTask 全字段）
│   ├── 下载 result（completed 时）
│   └── 跳转 job-history / platform-audit-log
└── 边界：Console compliance 只读；BOSS 独占创建
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 | 本页用法 |
|---|---|---|
| Core | POST compliance export **TODO-YAML** | 创建导出；返回 `202 + AsyncTask` 或 task_id |
| Core | `GET /api/v1/tasks/{task_id}` | 单任务状态；**路径已声明** |
| Core | `GET /api/v1/tasks` list **TODO-YAML** | BOSS 跨租户作业 list（见 job-history） |
| Services | `GET /api/v1/svc/compliance/*`（草案） | Console **只读参考**；**非** BOSS 创建入口 |

### 关键边界

- 合规包 **创建** 无 Core/Services 已冻结 path — 全部 **TODO-YAML** P2
- Console 草案 compliance reports **无 write** — BOSS 不得复用为平台导出 API
- 导出作业预期 `task_type` — **待 YAML** 扩展（当前 enum 四值不含 compliance；合入前标规划）
- 单查 `{task_id}` **不能** 替代 BOSS 跨租户 job list（禁止逐 task_id 轮询冒充 list）
- 下载 URL 签发/TTL — **待 YAML** 产品策略

## 页面区块与数据来源映射

| 区块 | 主要来源层 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| 创建导出 | POST **待 YAML** | idempotency_key + 范围 | 作业表 |
| 作业表格 | list **待 YAML** | AsyncTask 摘要 | 详情 |
| 详情抽屉 | `GET /tasks/{task_id}` | 全字段 | 下载 |
| 下载区 | `AsyncTask.result` | completed 时 | — |
| 边界说明 | 规划 | Console 只读 | compliance（Console） |

## BOSS 与 Console 分工

| 维度 | BOSS 合规导出与取证 | Console 合规能力 |
|---|---|---|
| 范围 | 全平台 / 指定租户 | 当前租户 |
| 创建导出 | BOSS 独占（**待 YAML**） | **无** write path |
| 报告 list | platform export jobs **待 YAML** | 草案 `GET .../compliance/reports` **待 YAML** |
| 摘要 dashboard | 可选平台汇总 **待 YAML** | 草案 `GET .../compliance/summary` **待 YAML** |
| 下载 | BOSS 签发 + AsyncTask result | 草案 download_url（只读） |
| RBAC | 平台 compliance 写/读（待 YAML） | `scope:compliance:read`（草案） |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/tasks/{task_id}` | **未声明** | 返回 `AsyncTask`；`200` / `404` |

| 字段 | 冻结值 |
|---|---|
| `AsyncTask.status` | `pending` / `running` / `completed` / `failed` / `cancelled` / `dead_letter` |
| `AsyncTask.task_type`（当前） | `model.import` / `kb.parse` / `kb.index` / `inference.deploy` |
| UI status 标签 | **必须用** `completed`；**禁止** `succeeded` |

| 能力 | 状态 |
|---|---|
| POST 平台合规导出 | **TODO-YAML** P2 |
| `GET /api/v1/tasks` 跨租户 list | **ADDED-TO-YAML** `listTasks`（共用 job-history） |
| compliance `task_type` enum 扩展 | **TODO-YAML** |
| Console `svc/compliance/*` | **草案/待 YAML** |

## 字段级定义

### 创建请求字段（POST 目标 · 待 YAML）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `idempotency_key` | body | ✅ | UUID |
| `framework` | body | ✅ | 如 iso27001、等保2.0（enum 待 YAML） |
| `tenant_scope` | body | 可选 | 全平台 / 单租户 / 租户列表 |
| `tenant_id` | body | 条件 | tenant_scope=single 时 |
| `start_time` / `end_time` | body | ✅ | 取证时间窗 |
| `include_datasets` | body | 可选 | audit_log / apikey / inference / config |
| `format` | body | 可选 | zip / tar.gz（待 YAML） |

### 路径参数（单查 · **路径已声明**）

| 字段 | 来源 | 必填 | 说明 |
|---|---|---|---|
| `task_id` | path | ✅ | UUID |

### 返回字段（AsyncTask · Core 已冻结）

| 字段 | 来源 | 说明 |
|---|---|---|
| `id` | `AsyncTask` | 任务 UUID |
| `idempotency_key` | `AsyncTask` | 幂等键 |
| `task_type` | `AsyncTask` | 合规导出待扩展 enum |
| `resource_type` | `AsyncTask` | 可空 |
| `resource_id` | `AsyncTask` | 可空 |
| `status` | `AsyncTask` | 六态；**禁止** succeeded |
| `attempt_count` | `AsyncTask` | 尝试次数 |
| `max_attempts` | `AsyncTask` | 最大尝试 |
| `progress_pct` | `AsyncTask` | 0–100 |
| `result` | `AsyncTask` | 含 download_url 等；completed 时 |
| `error_message` | `AsyncTask` | 失败原因 |
| `dead_letter_at` | `AsyncTask` | dead letter 时间 |
| `created_at` | `AsyncTask` | 创建时间 |
| `completed_at` | `AsyncTask` | 完成时间 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `tenant_scope_display` | 全平台 / 单租户摘要 |
| `framework_label` | framework 中文名 |
| `duration_display` | completed_at − created_at |
| `download_ready` | status=completed 且 result 含 URL |
| `is_terminal` | completed / failed / cancelled / dead_letter |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| 正常态 | 创建区 + 作业表 | 写 API 待 YAML 时创建区 disabled |
| 无数据态 | 「暂无导出作业」 | **不** 伪造 |
| 写 API 未就绪 | 创建按钮禁用 + 说明待 Core | — |
| `status=running` | progress_pct 进度条；轮询 ≥ 2s | 单查 task |
| `status=completed` | 展示 completed_at；启用下载 | **禁止** succeeded 标签 |
| `status=failed` / `dead_letter` | 高亮 error_message | 跳转 job-history |
| 单查 404 | 「任务不存在或无权查看」 | — |
| 403 | 无平台 compliance 写/读权限 | — |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `status` | YAML 六态；**never** succeeded | string |
| `progress_pct` | 0–100 整数 | integer |
| `start_time` / `end_time` | ISO 8601 | date-time |
| `framework` | 待 YAML enum | string |
| `result.download_url` | HTTPS URL；TTL 待 YAML | string |
| `include_datasets` | 字符串数组 enum | string[] |

## 状态与能力口径

### AsyncTask.status 状态机（导出作业）

| 状态 | 含义 | 终态 |
|---|---|---|
| `pending` | 已创建，等待执行 | 否 |
| `running` | 打包/导出中 | 否 |
| `completed` | 成功完成，可下载 | **是** |
| `failed` | 失败 | **是** |
| `cancelled` | 已取消 | **是** |
| `dead_letter` | 超过最大重试 | **是** |

**禁止** UI/文档使用 `succeeded` — 须用 `completed`。

### 状态 × 操作可用性矩阵

| 操作 \ status | pending | running | completed | failed | cancelled | dead_letter |
|---|---|---|---|---|---|---|
| 查看详情 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| 轮询进度 | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| 下载包 | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| 取消导出 | ⏳ 待 YAML | ⏳ 待 YAML | ❌ | ❌ | ❌ | ❌ |
| 重试 | ❌ | ❌ | ❌ | ⏳ 待 YAML | ❌ | ⏳ 待 YAML |
| 跳转 job-history | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

⏳ = 写 API 待 YAML。

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 compliance 导出写 RBAC | **TODO-YAML** | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| 时间窗 | start < end | `400 BAD_REQUEST` |
| `idempotency_key`（POST） | 合法 UUID | `400 BAD_REQUEST` |
| `tenant_scope=single` | 有效 tenant_id + 平台授权 | `403` / `404` 待 YAML |
| 单查 `GET /tasks/{task_id}` | 合法 UUID | `404 NOT_FOUND` |
| 下载 | status=completed + result | `422` 待 YAML |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 创建合规导出 | ❌ **待 YAML** | ⏳ | ⏳ | POST + idempotency_key |
| 作业 list | ✅ **待 YAML** | ✅ | ✅ | 共用 tasks list |
| 查看 task 详情 | ✅ | ✅ | ✅ | GET /tasks/{id} |
| 下载取证包 | ✅ completed 时 | ✅ | ✅ | result URL |
| 取消/重试 | ⏳ 待 YAML | ⏳ | ⏳ | job-history 写 API |
| 查看 Console 报告 | ✅ 跳转 | ✅ | ✅ | 只读对照 |

## 删除前置校验与当前契约边界

本页 **无 DELETE** 导出记录操作 — **N/A**。取证包对象存储生命周期与合规留存 — **待 YAML** 策略；物理删除非本页 Phase 2 范围。

取消导出为状态变更（`cancelled`）— **TODO-YAML**，非 DELETE。

## 接口冻结规则

### `GET /api/v1/tasks/{task_id}`（Core · **路径已声明**）

- `operationId`：权威源 **未声明**
- `summary`：`查询异步任务状态`
- `tags`：`["Tasks"]`
- `path.required`：`task_id`（UUID）
- `success`：`200 + AsyncTask`
- `errors`：`404`（YAML **已声明**）
- 认证/授权错误为网关通用推断，**未** 在本 operation responses 中冻结
- 本页用于 **轮询导出作业**；**不得** 自造 operationId

### POST 平台合规导出（待补）

<!-- TODO-YAML: POST /api/v1/compliance/exports 或 /platform/compliance/packages -->

- 须 `idempotency_key`
- 预期 `202 + AsyncTask`（或 task_id）
- 须扩展 `task_type` 含 compliance 类（如 `compliance.export`）
- **合入前不得写入「已冻结」正文**

### Console 草案 Services compliance（**草案/待 YAML · 非 BOSS 写入口**）

- 规划见 `docs/console-modules/openapi-drafts/phase3/openapi-phase3-compliance-draft.md`
- `GET /api/v1/svc/compliance/summary`、`.../reports` — Console 只读
- **合入前不得写入 BOSS 正式契约**

### `GET /api/v1/tasks` list（待补 · 见 job-history）

<!-- TODO-YAML: GET /api/v1/tasks with compliance task_type filter -->

- BOSS 跨租户作业 list — P1
- **合入前不得写入已实现**

## 使用规则

- **不得** 写合规导出 POST 为已实现
- UI status **必须** 使用 `completed`，**禁止** `succeeded`
- **不得** 把 Console 草案 compliance path 写成 BOSS 创建入口
- 轮询 `GET /tasks/{task_id}` 间隔 ≥ 2s
- list 未就绪时 **禁止** 伪造作业列表
- 下载链接须 HTTPS；TTL 过期须提示重新导出
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- POST 平台合规导出 — **TODO-YAML** P2
- `task_type` compliance 扩展 — P2
- `GET /api/v1/tasks` 跨租户 list — **ADDED-TO-YAML**（job-history）
- framework enum、download_url TTL — P2
- 与 audit-log / apikey / inference 数据集打包映射 — P2
- 取消/重试导出 — Phase 2 + `idempotency_key`

## 响应示例

### 创建导出成功（页面目标 · **待 YAML**）

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440010",
  "idempotency_key": "990e8400-e29b-41d4-a716-446655440011",
  "task_type": "compliance.export",
  "resource_type": null,
  "resource_id": null,
  "status": "pending",
  "attempt_count": 0,
  "max_attempts": 3,
  "progress_pct": 0,
  "result": null,
  "error_message": null,
  "dead_letter_at": null,
  "created_at": "2026-06-16T11:00:00Z",
  "completed_at": null
}
```

> `task_type=compliance.export` 为 **页面目标**；当前 YAML enum 四值 — 合入扩展前标规划。

### 单任务查询 completed（Core · **路径已声明**）

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440010",
  "idempotency_key": "990e8400-e29b-41d4-a716-446655440011",
  "task_type": "compliance.export",
  "status": "completed",
  "attempt_count": 1,
  "max_attempts": 3,
  "progress_pct": 100,
  "result": {
    "download_url": "https://storage.example.com/compliance/pkg-001.zip?sig=...",
    "expires_at": "2026-06-17T11:00:00Z",
    "sha256": "abc123..."
  },
  "error_message": null,
  "dead_letter_at": null,
  "created_at": "2026-06-16T11:00:00Z",
  "completed_at": "2026-06-16T11:05:00Z"
}
```

## 错误示例

### 时间窗非法（POST · **待 YAML**）

```json
{
  "code": "BAD_REQUEST",
  "message": "start_time must be before end_time",
  "request_id": "req-boss-compliance-400-001"
}
```

### 无平台 compliance 导出权限（**TODO-YAML**）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-compliance-403-001"
}
```

## 相关模块

- Console：[`compliance.md`](../../console-modules/security/compliance.md)
- [`job-history.md`](../health/job-history.md)（AsyncTask 跨租户 list / 重试边界）
- [`platform-audit-log.md`](platform-audit-log.md)
- [`platform-apikey-audit.md`](platform-apikey-audit.md)
- [`platform-inference-audit.md`](platform-inference-audit.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] `AsyncTask.status` 六态；**completed** 非 succeeded
- [x] `GET /tasks/{task_id}` 路径已声明；POST 导出 **TODO-YAML** P2
- [x] Console compliance 草案标 **待 YAML**；BOSS 独占写已说明
- [x] 含字段展示规则、口径、状态矩阵、400+403 错误示例
- [x] 删除前置校验 N/A
- [ ] compliance export + task_type YAML 合入后回写冻结表
- [x] PRD/SPEC/HTML 与本文同步
