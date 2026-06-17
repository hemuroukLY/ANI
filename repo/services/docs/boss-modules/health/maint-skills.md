# 运维 Skills

## 页面定位

`运维 Skills` 是 BOSS **运维与可观测** 域下的 **平台 SRE 标准化操作执行** 入口：将 gpu.drain、model.rollback、kb.reindex、inference.scale 等常见运维动作封装为 **可审计、可重试、可关联 AsyncTask** 的工作流。

本页属于 **Phase 2 产品能力**；Skill 编排 **整域 OpenAPI 未定义** — list/execute/definition CRUD 全部 **TODO-YAML**。一级权威源为 `ANI-main/repo/api/openapi/v1.yaml`（当前仅可关联 `AsyncTask` 与 ANI-14 写操作口径）。

Console **无对等页**；租户不得自助触发平台 Skill。触发入口来自 [`alert-rules.md`](alert-rules.md) 活动告警 Tab、[`incident-handling.md`](incident-handling.md) 工单页与本页目录。

## 文档管理规则

- 本文是 `运维 Skills` 的主维护源；页面定位、参数口径、字段映射和验收标准以本文为准
- `prototypes/ani-services-prototype-boss.html` 只保留摘要与详文入口
- [`prd-boss-maint-skills.md`](../../tasks/modules/prd/boss/health/prd-boss-maint-skills.md) 与 [`spec-boss-maint-skills.md`](../../tasks/modules/spec/boss/health/spec-boss-maint-skills.md) 为辅助材料，不替代本文
- 如本文与辅助材料冲突，先对照 `v1.yaml`，再统一回写 PRD/SPEC
- 流程：ANI-14 Phase 1 文档 + Phase 2 能力规划 + [`module-delivery-workflow.md`](../governance/module-delivery-workflow.md) + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- Skill 编排 API — **整域 TODO-YAML P2**；**不得** 自造 `/api/v1/observability/skills*`、`/api/v1/boss/skills*` 或 `/api/v1/svc/boss/skills*` 为已冻结
- 执行结果 **应** 关联 Core `AsyncTask`：`GET /api/v1/tasks/{task_id}` **路径已声明**（**无 operationId**、无 `x-ani-rbac-scope`）— 轮询只读参考
- `GET /api/v1/tasks` list — **TODO-YAML**；执行历史正式 list 见 [`job-history.md`](job-history.md)
- 所有 Skill **写操作** POST 须 `idempotency_key` + 平台 SRE RBAC（合入 YAML 后冻结）
- 平台 RBAC 鉴权；**不得** 信任 body/query 中未授权的 `tenant_id` 越权作用于其他租户资源
- 统一错误结构：`{"code":"UPPER_SNAKE","message":"...","request_id":"..."}`
- `422 PreconditionFailed` 仅 YAML 已声明 operation 可写冻结（§2.10）；Skill execute **待 YAML**
- 禁止自造 `MaintenanceSkill`、`SkillExecution` schema / operationId / 路径 / 返回码
- OpenAPI 已声明 ≠ handler 已实现

## 页面职责

- 展示 **Skills 目录**（产品枚举 · 非 REST 资源）：说明、风险级别、目标对象类型、所需 RBAC
- 提供 **执行弹窗**：选择目标对象、参数、风险确认、`idempotency_key`、提交 execute（API 待 YAML）
- 展示 **最近执行摘要**：关联 `task_id` → 跳转 [`job-history.md`](job-history.md) 或轮询 `GET /tasks/{task_id}`
- 接收来自 [`alert-rules.md`](alert-rules.md)、[`incident-handling.md`](incident-handling.md) 的 **预填上下文**（告警 ID / incident ID / 目标资源）
- **不承担** Console 租户自助运维；**不承担** Skill 定义 DSL 编辑器（Phase 3）

## 页面结构

- 首屏至少包含：`Skills 目录表`、`最近执行区`、`边界说明（Phase 2 / TODO-YAML）`
- 执行弹窗为独立 overlay；高风险 Skill 须二次确认（**UI 行为**，非 OpenAPI 字段）
- API 未就绪时：目录只读展示 + 执行按钮 disabled + 明确 TODO-YAML 文案
- 无权限态、执行失败态、AsyncTask 轮询失败态须可区分

```text
运维 Skills
├── Skills 目录表（skill_id / 说明 / 风险 / 目标类型 / RBAC）
├── 执行弹窗
│   ├── 目标对象（tenant / node / service / kb / …）
│   ├── 参数 JSON（待 YAML schema）
│   ├── 风险摘要 + 二次确认（high）
│   ├── idempotency_key（客户端生成）
│   └── 提交 → execute API **TODO-YAML** → 202 + task_id
├── 最近执行（task_id / skill_id / status / 时间）
└── 跳转 → job-history / incident-handling / audit
```

## 数据来源与分层约束

### 数据来源划分

| 层 | 路径 / 来源 | 本页用法 |
|---|---|---|
| Core | Skill execute/list API **TODO-YAML** | BOSS 正式编排数据源 |
| Core | `GET /api/v1/tasks/{task_id}` | 执行结果轮询只读参考 |
| Core | `GET /api/v1/tasks` list **TODO-YAML** | 执行历史 list → job-history |
| 产品 | Skills 目录枚举 | 页面静态/配置目录 · **非** REST |
| Services | 目标资源只读参考 | 校验对象存在（execute 前置 · 待 YAML） |

### 关键边界

- Skills 目录表是 **产品目标枚举**，**不得** 写成已实现 `GET /skills` REST list
- Skill 执行 **不得** 绕过 `AsyncTask` 审计链；成功 execute 应返回 `task_id`（目标契约 · 待 YAML）
- 现有 `AsyncTask.task_type` YAML 枚举为 `model.import` / `kb.parse` / `kb.index` / `inference.deploy` — **不包含** `gpu.drain` 等 Skill 类型；合入 YAML 前 **不得** 把 Skill task_type 写成已冻结
- `AsyncTask.status` 须用 `completed`，**禁止** UI 使用 `succeeded`
- 高风险 Skill（如 gpu.drain、tenant.quota）须 UI 二次确认；**不等于** OpenAPI 422

## 页面区块与数据来源映射

| 区块 | 主要来源 | 数据口径说明 | 跳转目标 |
|---|---|---|---|
| Skills 目录表 | 产品枚举 + 未来 list API **待 YAML** | `skill_id`、风险、目标类型 | 打开执行弹窗 |
| 执行弹窗 | POST execute **待 YAML** | `idempotency_key` 必填 | 提交后轮询 task |
| 最近执行 | task list **待 YAML** 或会话缓存 | `task_id` + `AsyncTask.status` | job-history |
| 执行详情 | `GET /tasks/{task_id}` | YAML `AsyncTask` 全字段 | job-history |
| 边界说明 | 规划项 | Phase 2 · 整域 TODO-YAML | — |
| 审计 | 平台 audit **待 YAML** | execute 须写 audit 事件 | platform-audit-log |

## BOSS 与 Console 分工

| 维度 | BOSS 运维 Skills | Console |
|---|---|---|
| 范围 | 平台 SRE 编排 | 无对等页 |
| 触发 | 告警/工单/本页 | — |
| 租户自助 | ❌ | — |
| 执行 API | **TODO-YAML** P2 | — |
| 任务追踪 | AsyncTask + job-history | async-task-center（单租户） |
| RBAC | 平台 SRE write | — |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/tasks/{task_id}` | **未声明** | 单任务查询；`AsyncTask` |

| `AsyncTask.status` | `pending` / `running` / `completed` / `failed` / `cancelled` / `dead_letter` |
| `AsyncTask.task_type`（YAML 现有） | `model.import` / `kb.parse` / `kb.index` / `inference.deploy` |

| 能力 | 状态 |
|---|---|
| Skill list / execute / definition CRUD | **TODO-YAML** P2 |
| `GET /api/v1/tasks` list | **ADDED-TO-YAML** (`listTasks`)（job-history 依赖） |
| Skill 专用 `task_type` 扩展 | **TODO-YAML** P2 |

## 字段级定义

### Skills 目录字段（产品目标 · **非冻结 schema**）

| 字段 | 说明 |
|---|---|
| `skill_id` | 稳定标识，如 `gpu.drain` |
| `display_name` | 展示名 |
| `description` | 说明 |
| `risk_level` | `low` / `medium` / `high`（**产品 enum · 待 YAML**） |
| `target_kind` | 目标类型：node / gpu / inference_service / kb / tenant / … |
| `required_scopes[]` | 执行所需 RBAC scope |
| `estimated_duration` | 预估耗时（展示） |
| `supports_dry_run` | 是否支持 dry-run（**待 YAML**） |

### 执行请求字段（目标契约 · **TODO-YAML**）

| 字段 | 必填 | 说明 |
|---|---|---|
| `idempotency_key` | ✅ | 1–128 字符；幂等 |
| `skill_id` | ✅ | 目录中的 skill_id |
| `target` | ✅ | `{ kind, id, tenant_id? }`；tenant_id 由 RBAC 校验 |
| `parameters` | ❌ | Skill 特定 JSON |
| `dry_run` | ❌ | 默认 false |
| `incident_id` | ❌ | 关联 [`incident-handling.md`](incident-handling.md) |
| `alert_incident_id` | ❌ | 关联 alert-rules 活动 Tab |

### 执行响应字段（目标契约 · **TODO-YAML**）

| 字段 | 说明 |
|---|---|
| `task_id` | UUID；轮询 `GET /tasks/{task_id}` |
| `status` | 初始 `pending` |
| `audit_id` | 审计事件 ID（**待 YAML**） |

### 返回字段 — `AsyncTask`（YAML 已冻结 · 轮询参考）

| 字段 | 说明 |
|---|---|
| `id` | task UUID |
| `idempotency_key` | 与 execute 请求一致 |
| `task_type` | 合入后可能扩展 Skill 类型 |
| `status` | 六态；**禁止** succeeded |
| `progress_pct` | 0–100 |
| `result` / `error_message` | 完成/失败详情 |
| `created_at` / `completed_at` | 时间戳 |

### 展示字段（UI 计算 · 非 API 字段）

| 字段 | 说明 |
|---|---|
| `success_rate_pct` | 最近 N 次成功率（aggregate **待 YAML**） |
| `last_run_at` | 最近执行时间 |
| `risk_badge` | 由 risk_level 映射颜色 |
| `execution_duration` | completed_at − created_at |

## 字段展示规则

| 场景 | 展示规则 | 说明 |
|---|---|---|
| API 未就绪 | 目录可读 + 执行 disabled + TODO-YAML 横幅 | **不** 伪装可执行 |
| 正常态（P2 后） | 目录 + 最近执行 + 可打开弹窗 | — |
| 执行提交成功 | 展示 task_id + 跳转 job-history | 202/200 待 YAML |
| 轮询 running | progress_pct 进度条 | GET task |
| 轮询 completed | 绿色完成 + result 摘要 | **禁止** succeeded 标签 |
| 轮询 failed / dead_letter | 红色 + error_message | 可重试（待 YAML） |
| 无权限态 | 403；隐藏执行入口 | 平台 SRE RBAC |
| 高风险 Skill | 二次确认 modal | UI 强制 |
| 目标对象不存在 | 404 或表单校验错误 | execute 前置 |
| 幂等冲突 | 返回同一 task_id（待 YAML 语义） | 复用 idempotency_key |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `AsyncTask.status` | YAML 六态 | string |
| `progress_pct` | 0–100 整数 | % |
| `risk_level` | 产品三档 | low/medium/high |
| `idempotency_key` | UUID 或客户端串 | 1–128 字符 |
| `execution_duration` | UI 计算 | 秒/分 humanize |

## 状态与能力口径

### AsyncTask.status（执行结果 · 已冻结）

| 状态 | 含义 | 本页展示 |
|---|---|---|
| `pending` | 已提交 | 排队中 |
| `running` | 执行中 | 进度条 |
| `completed` | 成功 | 绿色；**禁止** succeeded |
| `failed` | 失败 | 红色 + error |
| `cancelled` | 已取消 | 灰色 |
| `dead_letter` | 死信 | 红色 + 跳转 job-history |

### Skill 执行生命周期（目标 · **待 YAML**）

提交 execute → 返回 `task_id` → 轮询 AsyncTask → 写入 audit → 可选通知 incident/alert。

| 能力 | Phase |
|---|---|
| 查看目录 | P2 文档 / P2 API |
| 执行 Skill | P2 |
| 编辑 Skill 定义 | P3 |
| 取消执行中 task | **待 YAML** |

## 创建前置条件

| 依赖项 | 要求状态 | 未满足时的 HTTP 响应 |
|---|---|---|
| 平台 SRE 读 RBAC | 已授权 | `403 FORBIDDEN` |
| 平台 SRE 写 RBAC（执行） | 已授权 | `403 FORBIDDEN` |
| 用户已认证 | 平台登录态 | `401 UNAUTHORIZED` |
| `idempotency_key` | 非空 | `400 BAD_REQUEST`（execute 待 YAML） |
| `skill_id` | 目录中存在 | `400` / `404`（待 YAML） |
| 目标对象存在 | 资源可读 | `404 NOT_FOUND`（待 YAML） |
| 目标租户 suspended | 产品规则 | `422`（**待 YAML 合入后冻结**） |
| 高风险 Skill | UI 二次确认 | —（非 HTTP） |

## 操作可用性矩阵

| 操作 | 平台只读 | SRE | 平台管理员 | 说明 |
|---|---|---|---|---|
| 查看 Skills 目录 | ✅ | ✅ | ✅ | 产品枚举 |
| 执行 Skill | ❌ | ✅（待 YAML） | ✅（待 YAML） | POST execute |
| dry-run | ❌ | ✅（待 YAML） | ✅ | 待 YAML |
| 取消执行 | ❌ | Phase 2 | Phase 2 | task cancel 待 YAML |
| 编辑 Skill 定义 | ❌ | Phase 3 | Phase 3 | — |
| 查看执行历史 | ✅ → job-history | ✅ | ✅ | task list 待 YAML |
| 从告警/工单预填执行 | ❌ | ✅ | ✅ | 深链参数 |

## 删除前置校验与当前契约边界

本页 **无 Skill 资源 DELETE**（定义 CRUD 待 YAML）— Phase 2 文档口径 **N/A**。

合入 YAML 后若支持「撤销 Skill 定义」或「取消执行」，须单独写 DELETE/PATCH 冻结规则，**不得** 提前写入本文「已冻结」表。

## 接口冻结规则

### `GET /api/v1/tasks/{task_id}`（Core · **路径已声明** · **非 Skill execute**）

- `summary`：`查询异步任务状态`
- `tags`：`["Tasks"]`
- `path`：`task_id`（uuid）
- `success`：`200 + AsyncTask`
- `errors`：`404`（YAML **已声明**）
- 认证/授权错误为网关通用推断，**未** 在本 operation responses 中冻结
- **无 operationId** — 正文写「未声明」
- 用途：execute 提交后轮询；**不是** Skill list API

### Skill 编排 API（待补 · **整域未冻结**）

<!-- TODO-YAML: POST /api/v1/observability/skills/{skill_id}/execute 或等价 -->

| 目标能力 | 状态 |
|---|---|
| list skills | **TODO-YAML** |
| execute skill | **TODO-YAML**；须 `idempotency_key`；建议 `202 + task_id` |
| get execution | 复用 `GET /tasks/{task_id}` 或专用 path **待 YAML** |

- 合入前 **不得** 写入「已冻结路径」表
- 路径前缀须符合 Core `/api/v1/*`；**不得** 使用 `/api/v1/boss/*`

## 使用规则

- **不得** 把 Skills 目录表实现为伪造 REST list 并标注已上线
- **不得** 在无 execute API 时启用生产环境执行按钮
- 执行 **必须** 带 `idempotency_key`；UI 重试须复用同一 key
- 轮询 **必须** 使用 YAML `AsyncTask.status` 六态
- Skill 触达租户资源时 **必须** 经平台 RBAC 校验，不得信任前端 `tenant_id`
- execute 成功 **应** 产生 audit 事件（audit API 待 YAML）
- 从 alert-rules / incident-handling 跳转时须携带上下文 query，**不得** 重复定义 firing/incident API

## 待补能力边界

- Skill list / execute / cancel API — **TODO-YAML** P2
- Skill 专用 `AsyncTask.task_type` 扩展 — **TODO-YAML** P2
- 与 [`job-history.md`](job-history.md) list 联动 — 依赖 `GET /tasks` **ADDED-TO-YAML**
- 与 [`alert-rules.md`](alert-rules.md)「触发 Skill」— 依赖 `/observability/alert-events` **ADDED-TO-YAML**
- 与 [`../audit/platform-audit-log.md`](../audit/platform-audit-log.md) execute 审计 — Phase 2
- Skill 定义 DSL / 自定义 Skill — Phase 3

## 响应示例

### AsyncTask 轮询成功（Core · **路径已声明** · execute 后参考）

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "idempotency_key": "skill-exec-20260616-001",
  "task_type": "inference.deploy",
  "resource_type": "inference_service",
  "resource_id": "660e8400-e29b-41d4-a716-446655440001",
  "status": "completed",
  "attempt_count": 1,
  "max_attempts": 3,
  "progress_pct": 100,
  "result": { "message": "scale completed", "replicas": 4 },
  "error_message": null,
  "created_at": "2026-06-16T10:00:00Z",
  "completed_at": "2026-06-16T10:05:00Z"
}
```

### Skill execute 目标响应（**待 YAML · 非已冻结**）

```json
{
  "task_id": "550e8400-e29b-41d4-a716-446655440000",
  "skill_id": "gpu.drain",
  "status": "pending",
  "audit_id": "aud-20260616-001"
}
```

## 错误示例

### 无 Skill execute 写权限（**TODO-YAML**）

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-skill-403-001"
}
```

> **注**：适用于 **Skill execute API（TODO-YAML）** 场景。`message` 为示例文案，**非** YAML 已冻结 scope 名。

### 缺少 idempotency_key（execute 目标契约）

```json
{
  "code": "BAD_REQUEST",
  "message": "idempotency_key is required",
  "request_id": "req-boss-skill-400-001"
}
```

### 关联 task 不存在

```json
{
  "code": "NOT_FOUND",
  "message": "task not found",
  "request_id": "req-boss-skill-404-001"
}
```

## Skills 目录（产品目标 · 非冻结 API）

| skill_id | 说明 | 风险 | 关联域 |
|---|---|---|---|
| `gpu.drain` | GPU 节点降载 | high | [`../ops/gpu-pool-management.md`](../ops/gpu-pool-management.md) |
| `model.rollback` | 模型版本回滚 | medium | Services models |
| `kb.reindex` | 知识库重新索引 | medium | Services KB |
| `inference.scale` | 推理扩缩容 | medium | Services inference |
| `audit.export` | 审计导出 | low | [`../audit/platform-audit-log.md`](../audit/platform-audit-log.md) |
| `tenant.quota` | 租户配额调整 | high | [`../tenant/tenant-quota-policy.md`](../tenant/tenant-quota-policy.md) |

## 相关模块

- [alert-rules.md](alert-rules.md)、[incident-handling.md](incident-handling.md)、[job-history.md](job-history.md)
- [platform-health.md](platform-health.md)、[gpu-monitoring.md](gpu-monitoring.md)

## 回填验收标准

- [x] 满配章节齐全（对照 [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)）
- [x] Skill 整域 TODO-YAML；无伪造 REST 为已冻结
- [x] `AsyncTask.status` 六态与 `completed` 口径一致
- [x] 含响应示例与错误示例（400 + 403 + 404）
- [x] 独立字段定义（目录/execute/AsyncTask/展示）
- [ ] Skill execute YAML 合入后回写「已冻结路径」表
- [x] PRD/SPEC/HTML 与本文同步
