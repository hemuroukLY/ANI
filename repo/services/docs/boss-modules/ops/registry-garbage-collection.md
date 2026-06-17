# 镜像仓库 / 垃圾回收

## 页面定位

`镜像仓库 / 垃圾回收` 是 BOSS **资源池与基础设施 / 镜像仓库** 域下的 **全平台 Harbor GC** 专页：查看 GC 任务历史、释放空间预估、调度策略与手动触发入口。

Console **无** 直接对等页。**当前 v1.yaml 无 registry GC 相关 path** — 全部 **TODO-YAML**。

## 文档管理规则

- 本文是 `镜像仓库 / 垃圾回收` 的主维护源
- PRD/SPEC 为辅助
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- **无** 已冻结 `GET/POST /api/v1/registry/gc*` 或等价 operation
- 可 **只读参考** 租户 registry 体量间接指标：
  - `GET /api/v1/registry/projects` — 项目 list
  - `GET /api/v1/registry/projects/{project}/repositories` — `artifact_count`
  - `GET /api/v1/registry/projects/{project}/repositories/{repository}/artifacts` — `size_bytes` 汇总（页面计算 · **非** GC API）
- GC 触发、调度、历史 list — **TODO-YAML** P2
- 平台 RBAC — **TODO-YAML**；GC 写操作须审计 + `idempotency_key`（合入时须 YAML 声明）
- **不得** 自造 `RegistryGCJob` schema 为已冻结

## 页面职责

- 展示最近 GC 任务：状态、开始/结束时间、释放空间、触发人
- 手动触发 GC（dry-run + 正式运行 — 待 YAML）
- 调度策略摘要（cron/保留天数 — 待 YAML）
- 跳转项目配额（回收后配额使用率）、漏洞扫描
- **不承担** 删除单个 artifact（归属 registry 权限/Console 流程 · 待产品）

## 页面结构

```text
镜像仓库 / 垃圾回收
├── 顶部 KPI（可回收空间 / 上次 GC 释放量）
├── GC 任务历史表格
├── 手动触发（dry-run → confirm → run — 待 YAML）
├── 调度策略（待 YAML）
└── 跳转
    ├── 项目配额
    └── 漏洞扫描
```

## 数据来源与分层约束

| 层 | 路径 | 用法 |
|---|---|---|
| Core | registry GC list/run **TODO-YAML** | BOSS 正式数据源 |
| Core | `listRegistryProjects` / repositories / artifacts | 体量 **只读参考**；**非** GC 状态 |
| Core | `GET /api/v1/tasks/{task_id}` | 若 GC 异步化 — **路径已声明**；list 待 YAML |

### 关键边界

- artifact `size_bytes` 之和 **≠** GC 可回收空间（未引用 blob 口径由 Harbor GC 定义）
- 不得把 `listRegistryArtifacts` 长度写成「待 GC 数量」
- 异步 GC 须用 `AsyncTask.status`=`completed`（禁 `succeeded`）

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| KPI | GC aggregate **待 YAML** | — |
| 任务历史 | GC jobs list **待 YAML** | 任务详情 |
| 手动触发 | POST GC **待 YAML** | job-history |
| 调度 | GC schedule **待 YAML** | — |

## BOSS 与 Console 分工

| 维度 | BOSS | Console |
|---|---|---|
| 范围 | 全平台 Harbor GC | 无 |
| 操作 | 平台触发 GC | 租户可能 delete artifact（permissions） |
| API | **全部 TODO-YAML** | 租户 registry CRUD |

## 当前冻结事实

| 方法 | 路径 | operationId | 说明 |
|---|---|---|---|
| GET | `/api/v1/registry/projects` | `listRegistryProjects` | 体量参考 |
| GET | `/api/v1/registry/projects/{project}/repositories` | `listRegistryRepositories` | `artifact_count` 参考 |
| GET | `/api/v1/registry/projects/{project}/repositories/{repository}/artifacts` | `listRegistryArtifacts` | `size_bytes` 汇总参考 |
| GET | `/api/v1/tasks/{task_id}` | — | 异步 GC **路径已声明** |

| 能力 | 状态 |
|---|---|
| GC jobs list | **TODO-YAML** P2 |
| POST trigger GC / dry-run | **TODO-YAML** P2 |
| GC schedule CRUD | **TODO-YAML** P2 |
| `GET /api/v1/tasks` list | **TODO-YAML** |

## 字段级定义

| 字段 | 说明 | 来源 |
|---|---|---|
| `gc_job_id` | GC 任务 ID | **待 YAML** |
| `status` | `pending`/`running`/`completed`/`failed` | **待 YAML** / AsyncTask |
| `triggered_by` | 操作者 | **待 YAML** |
| `dry_run` | 是否演练 | **待 YAML** |
| `space_reclaimed_bytes` | 释放字节 | **待 YAML** |
| `started_at` / `finished_at` | 时间 | **待 YAML** |
| `schedule_cron` | 调度 | **待 YAML** |
| `retention_days` | 保留策略 | **待 YAML** |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| 无 GC API | 全页「GC API 待 Core 合入」 |
| `failed` | 红标 + error 展开 |
| dry-run | 明确标注「未实际删除」 |
| 403 | 无权限 |
| 运行中 | 禁止重复触发（产品规则 · 待 YAML 409/422） |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `space_reclaimed_bytes` | 实际 GC 或 dry-run 预估回收空间 | bytes；UI 可换算 GiB/TiB |
| `artifact_count` | 命中 artifact 数 | integer |
| `dry_run` | 是否仅预估 | boolean；必须显式标注 |
| `started_at` / `completed_at` | 任务时间 | ISO 8601 date-time |

## 状态与能力口径

GC job 状态对齐 `AsyncTask` 或专用 enum（**待 YAML**）。`completed` 禁 UI 写 `succeeded`。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台 registry GC 写 RBAC（**TODO-YAML**） | `403` |
| 上次 GC 未完成 | `409` 或 `422`（**待 YAML**） |
| POST 须 `idempotency_key` | `400` |
| 未认证 | `401` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 平台管理员 |
|---|---|---|---|
| 查看 GC 历史 | 待 YAML | 待 YAML | 待 YAML |
| dry-run GC | 待 YAML | 待 YAML | 待 YAML |
| 正式触发 GC | 待 YAML | 待 YAML | 待 YAML |
| 修改调度 | 待 YAML | 待 YAML | 待 YAML |

## 删除前置校验

GC 本身即删除未引用 blob — 须 YAML 声明 dry-run 与 confirm 流程；本页文档不写 artifact DELETE 细节。

## 接口冻结规则

### 当前无 GC 冻结 API

<!-- TODO-YAML: GET /api/v1/registry/gc/jobs -->
<!-- TODO-YAML: POST /api/v1/registry/gc/run -->
<!-- TODO-YAML: GET/PATCH /api/v1/registry/gc/schedule -->

- 合入前 **不得** 写入「已冻结路径」表
- 引用 `listRegistryProjects` 仅作体量参考 — **非** GC 契约

### `GET /api/v1/tasks/{task_id}`（若 GC 异步 · **路径已声明**）

- 无 operationId / 无 x-ani-rbac-scope
- **无** list — 历史表格须 GC jobs API 或 tasks list **TODO-YAML**

## 使用规则

- 禁止伪造 GC 历史表格
- dry-run 与正式 run 须 UI 区分
- GC 为平台级高危操作 — 须审计日志（audit 域）

## 待补能力边界

- 全部 GC REST — **TODO-YAML** P2
- 与 Harbor 原生 GC 映射 — 待 Core/交付文档
- 租户级 artifact delete — Console/registry permissions · 非本页

## 响应示例

### GC run 目标响应（**待 YAML · 非已冻结**）

```json
{
  "gc_job_id": "gc-20260617-001",
  "status": "pending",
  "dry_run": false,
  "triggered_by": "platform-admin",
  "task_id": "task-gc-001",
  "started_at": "2026-06-17T11:00:00Z"
}
```

### GC job 完成（**待 YAML**）

```json
{
  "gc_job_id": "gc-20260617-001",
  "status": "completed",
  "space_reclaimed_bytes": 85899345920,
  "finished_at": "2026-06-17T11:45:00Z"
}
```

## 错误示例

### 无平台 GC 权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied",
  "request_id": "req-boss-rgc-403-001"
}
```

### GC 已在运行（**待 YAML**）

```json
{
  "code": "CONFLICT",
  "message": "garbage collection already in progress",
  "request_id": "req-boss-rgc-409-001"
}
```

> 注：`409`/`CONFLICT` 仅当 YAML 声明后可标为已冻结；当前为产品目标示例。

## 相关模块

- [`registry-project-quota.md`](registry-project-quota.md)、[`registry-vulnerability-scan.md`](registry-vulnerability-scan.md)
- [`job-history.md`](../health/job-history.md)（异步 task 路径已声明）

## 回填验收标准

- [x] 满配章节齐全
- [x] 明确 **无** GC 冻结 API；artifact 体量参考边界清晰
- [x] 403 + 409 错误示例（409 标注待 YAML）
- [ ] GC YAML 合入后回写
- [x] PRD/SPEC/HTML 与本文同步
