# 镜像仓库 / 项目配额

## 页面定位

`镜像仓库 / 项目配额` 是 BOSS **资源池与基础设施 / 镜像仓库** 域下的 **全平台 Harbor 项目配额** 专页：查看与调整各租户 registry 项目的存储配额、已用量与超额风险。

Console **无** 直接对等页（平台专属）。关联：[`tenant-quota-policy.md`](../tenant/tenant-quota-policy.md)（租户算力配额，**非** registry 配额）、[`registry-vulnerability-scan.md`](registry-vulnerability-scan.md)。

## 文档管理规则

- 本文是 `镜像仓库 / 项目配额` 的主维护源
- PRD/SPEC 为辅助；冲突以本文 + `v1.yaml` 为准
- 流程：ANI-14 + [`boss-full-depth-checklist.md`](../governance/boss-full-depth-checklist.md)

## Core 层要求

- Core `/api/v1/registry/*`；**当前 `RegistryProject` schema 不含 quota 字段**
- 租户 **只读参考**：`GET /api/v1/registry/projects`（`listRegistryProjects`）
- 平台 registry 项目 quota list/PATCH — **TODO-YAML**
- RBAC 租户路径：`scope:registry:read` / `scope:registry:write`
- 平台 RBAC — **TODO-YAML**；禁止信任未授权 `tenant_id`
- POST/PATCH 须 `idempotency_key`（YAML 已声明 operation）
- 错误结构：`code` / `message` / `request_id`
- **不得** 把 `tenant-quota-policy` 的 GPU/CPU quota 与 registry 存储 quota 混为一表

## 页面职责

- 列表展示全平台 registry 项目：租户、项目名、配额上限、已用、使用率
- 支持按租户/超额筛选；调整配额（待 YAML 写 API）
- 跳转漏洞扫描、垃圾回收、租户列表
- 标明 `RegistryProject` 当前 **无** quota 字段 — 合入前不得伪造列值

## 页面结构

```text
镜像仓库 / 项目配额
├── 筛选（租户 / 超额 / 项目名）
├── 项目配额表格
├── 调整配额抽屉（PATCH — 待 YAML）
└── 跳转
    ├── 漏洞扫描
    └── 垃圾回收
```

## 数据来源与分层约束

| 层 | 路径 | 用法 |
|---|---|---|
| Core | `GET /api/v1/registry/projects` | 租户 **只读参考**；`RegistryProject` |
| Core | 平台 registry quota list **TODO-YAML** | BOSS 正式数据源 |
| Core | PATCH project quota **TODO-YAML** | 写操作 |

### 关键边界

- `listRegistryProjects` 返回 `tenant_id` — BOSS 跨租户 list 须 **platform API**，禁止 JWT 轮询
- `CreateRegistryProjectRequest` 不含 quota — 创建时配额策略 **待 YAML**
- Harbor 存储 quota **≠** tenant GPU quota

## 页面区块与数据来源映射

| 区块 | 来源 | 跳转 |
|---|---|---|
| 项目列表 | aggregate **待 YAML** + projects 参考 | — |
| 配额列 | quota API **待 YAML** | — |
| 超额告警 | 产品计算 | registry-vuln（若扫描阻断） |

## BOSS 与 Console 分工

| 维度 | BOSS 本页 | Console |
|---|---|---|
| 范围 | 全平台 registry 项目配额 | 无对等页 |
| 操作 | 平台调整 Harbor quota | 租户可能 push 镜像（registry write） |
| API | platform quota **TODO-YAML** | `registry/projects` 租户上下文 |

## 当前冻结事实

| 方法 | 路径 | operationId | RBAC | BOSS 用法 |
|---|---|---|---|---|
| GET | `/api/v1/registry/projects` | `listRegistryProjects` | `scope:registry:read` | 租户 **只读参考** |
| POST | `/api/v1/registry/projects` | `createRegistryProject` | `scope:registry:write` | 租户创建；BOSS 代建 **待 YAML** |

| Schema | 说明 |
|---|---|
| `RegistryProject` | `id`, `tenant_id`, `name`, `public`, `created_at` — **无 quota** |

| 能力 | 状态 |
|---|---|
| 平台 registry quota list/PATCH | **TODO-YAML** P2 |
| BOSS 跨租户 projects list | **TODO-YAML** P2 |

## 字段级定义

| 字段 | 来源 | 说明 |
|---|---|---|
| `project_id` | `RegistryProject.id` | 项目 ID |
| `tenant_id` | `RegistryProject.tenant_id` | 所属租户 |
| `project_name` | `RegistryProject.name` | 项目名 |
| `public` | `RegistryProject.public` | 是否公开 |
| `storage_quota_bytes` | **待 YAML** | 配额上限 |
| `storage_used_bytes` | **待 YAML** | 已用 |
| `usage_pct` | UI 计算 | 使用率 % |
| `over_quota` | UI 计算 | 是否超额 |

## 字段展示规则

| 场景 | 规则 |
|---|---|
| quota API 未就绪 | 配额列「待 Core 合入」；可展示 projects 基础列 |
| `over_quota` | 红色高亮 |
| 403 | 无权限 |
| PATCH 422 | 仅 YAML 声明 operation 可写 |

## 字段口径与单位

| 字段 | 口径 | 单位/格式 |
|---|---|---|
| `quota_bytes` | 项目存储配额目标值 | bytes；UI 可换算 GiB |
| `used_bytes` | 项目已用空间 | bytes；UI 可换算 GiB |
| `usage_pct` | `used_bytes / quota_bytes * 100` | 0–100%，保留 1 位 |
| `project_count` | 项目数量 | integer |

## 状态与能力口径

本页写操作（调整 quota）待 YAML。`RegistryProject` 无状态 enum。

## 创建前置条件

| 依赖 | 响应 |
|---|---|
| 平台 registry 读（**TODO-YAML**） | `403` |
| PATCH quota 写（**TODO-YAML**） | `403` / `422` |
| POST create 引用时 | `idempotency_key` + `scope:registry:write` |
| 未认证 | `401` |

## 操作可用性矩阵

| 操作 | 只读 | SRE | 平台管理员 |
|---|---|---|---|
| 查看项目配额 | ✅ | ✅ | ✅ |
| 调整项目 quota | 待 YAML | 待 YAML | 待 YAML |
| 创建 registry 项目 | 待 YAML | 待 YAML | 待 YAML |
| 跳转扫描/GC | ✅ | ✅ | ✅ |

## 删除前置校验

**N/A**（本页 Phase 2 文档不写 DELETE project；若未来 YAML 声明须单列）

## 接口冻结规则

### `GET /api/v1/registry/projects`（租户 · **只读参考**）

- Query：`limit`、`cursor`
- `success`：`200` + `RegistryProjectListResponse`
- `errors`：`401`、`403`
- **非** BOSS 跨租户正式契约

### `POST /api/v1/registry/projects`（租户 · 创建）

- `CreateRegistryProjectRequest`：`idempotency_key`、`name`（`public` 可选，默认 `false`）
- `success`：`201` + `RegistryProject`
- `errors`：`400`、`401`、`403`
- BOSS 代建其他租户 — **platform API TODO-YAML**

### 平台 registry quota（待补）

<!-- TODO-YAML: GET/PATCH /api/v1/registry/projects/{project}/quota 或 platform list -->

## 使用规则

- 不得伪造 `storage_quota_bytes`
- registry quota 与 tenant GPU quota 分表展示
- OpenAPI 已声明 ≠ handler 已实现

## 待补能力边界

- 平台 registry quota list/PATCH — **TODO-YAML** P2
- `RegistryProject` 内嵌 quota 或子资源设计 — 待 Core
- 超额自动阻断 push — Phase 2 产品策略

## 响应示例

### listRegistryProjects 引用（租户 · **非 BOSS quota 正式响应**）

```json
{
  "items": [
    {
      "id": "reg-proj-001",
      "tenant_id": "t-001",
      "name": "acme-models",
      "public": false,
      "created_at": "2026-06-01T08:00:00Z"
    }
  ],
  "total": 1,
  "next_cursor": null
}
```

### 平台 quota 目标（**待 YAML**）

```json
{
  "items": [
    {
      "project_id": "reg-proj-001",
      "tenant_id": "t-001",
      "project_name": "acme-models",
      "storage_quota_bytes": 1099511627776,
      "storage_used_bytes": 912680550400,
      "usage_pct": 83.1,
      "over_quota": false
    }
  ]
}
```

## 错误示例

### POST create 参数无效（Core · 已声明）

```json
{
  "code": "BAD_REQUEST",
  "message": "name is required",
  "request_id": "req-boss-rpq-400-001"
}
```

### 无 registry 读权限

```json
{
  "code": "FORBIDDEN",
  "message": "permission denied: scope:registry:read required",
  "request_id": "req-boss-rpq-403-001"
}
```

### PATCH quota 前置失败（**待 YAML · 422 示例**）

```json
{
  "code": "PRECONDITION_FAILED",
  "message": "quota below current usage",
  "request_id": "req-boss-rpq-422-001"
}
```

> 注：422 仅当对应 operation YAML 已声明时可标为已冻结。

## 相关模块

- [`registry-vulnerability-scan.md`](registry-vulnerability-scan.md)、[`registry-garbage-collection.md`](registry-garbage-collection.md)
- [`tenant-quota-policy.md`](../tenant/tenant-quota-policy.md)（算力配额，非 registry）

## 回填验收标准

- [x] 满配章节齐全
- [x] `RegistryProject` 无 quota 字段已标注；租户 paths 只读参考
- [x] 400 + 403 错误示例
- [ ] registry quota YAML 合入后回写
- [x] PRD/SPEC/HTML 与本文同步
