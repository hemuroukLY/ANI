# Core Handler 实现指南（文档交付）

> **生成日期**：2026-06-17  
> **用途**：Core 团队按 OpenAPI 实现 handler 时的**唯一技术契约文档**（本阶段仅文档，不含代码补丁）。  
> **权威源**：`ANI-main/repo/api/openapi/v1.yaml`  
> **Console 模块详文**：`docs/console-modules/`（页面边界与待补说明）  
> **任务索引**：`tasks/execution/CORE-TEAM-TASKS.md`  
> **依赖关系**：`tasks/execution/TASK-DEPENDENCY-MAP.md`

## 文档管理规则（Core 层）

1. 路径、schema、成功/错误返回码以 `v1.yaml` 为准；不得自造 `/api/v1/console/*` 或 Services 路径。
2. GPU 容器 / Sandbox 统一实例：`/api/v1/instances*` + `kind=gpu_container|sandbox`；**不实现** deprecated Services `/gpu-containers*`、`/sandboxes*`。
3. 租户 ID 从 JWT / API Key 上下文提取；handler **不得**要求请求体传 `tenant_id`。
4. 写操作（POST 创建类）须校验 `idempotency_key`（YAML `required` 处）。
5. `422 PreconditionFailed`：仅使用 YAML 已举例或 description 中声明的 `code`（如 `INSTANCE_STATE_INVALID`、`IMAGE_NOT_FOUND`）；其余前置失败写「建议语义」，见 `../governance/module-delivery-workflow.md` §2.10。
6. OpenAPI 已声明 ≠ handler 已实现；验收以本文 curl 与集成测试为准。

## 错误响应结构（全 Core API）

```json
{
  "code": "UPPER_SNAKE",
  "message": "human readable",
  "request_id": "..."
}
```

标准错误码：`BAD_REQUEST`、`UNAUTHORIZED`、`FORBIDDEN`、`NOT_FOUND`、`CONFLICT`、`PRECONDITION_FAILED`（HTTP 422 时 body.code 可为 YAML 举例值）。

---

## TASK-CORE-001 — 实例列表 kind 过滤

| 项 | 值 |
|---|---|
| 状态 | 文档 + handler 验通（Sprint 1） |
| 接口 | `GET /api/v1/instances?kind={vm\|container\|gpu_container\|sandbox}` |
| operationId | `listInstances` |
| RBAC | `scope:instances:read` |
| 模块详文 | `docs/console-modules/compute/sandbox-instance-management.md` 等 |

### Handler 契约

- Query `kind` 可选；若提供则只返回对应 `InstanceRecord.kind`。
- 成功：`200` + `InstanceListResponse`（`items`、`next_cursor`）。
- 错误：`401`、`403`。

### 验收

```bash
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/instances?kind=sandbox&limit=20"
# 期望 200；items[].kind 均为 sandbox
```

---

## TASK-CORE-002 — 实例创建 / 生命周期 422

| 项 | 值 |
|---|---|
| 状态 | 文档 + handler 验通（Sprint 1） |
| 接口 | `POST /api/v1/instances`、`POST /api/v1/instances/{instance_id}/lifecycle` |
| operationId | `createInstance`、`applyInstanceLifecycle` |
| 模块详文 | `vm-management.md`、`container-instance-management.md`、`gpu-container-instance-management.md`、`sandbox-instance-management.md` |

### Handler 契约

**创建 `POST /instances`**

- Request：`CreateInstanceRequest`（含 `idempotency_key`、`kind`）。
- 成功：`201` + `CreateInstanceResponse`。
- 镜像/引导镜像不存在：HTTP `422`，body.code **`IMAGE_NOT_FOUND`**（YAML 已举例）。
- 幂等冲突：`409`。

**生命周期 `POST .../lifecycle`**

- Request：`InstanceLifecycleRequest`（`action`、`idempotency_key`）。
- 成功：`200` + `InstanceLifecycleResponse`。
- 当前状态不允许该 action：HTTP `422`，body.code **`INSTANCE_STATE_INVALID`**（YAML 已举例）。
- 典型验通：`running` 实例再次 `start` → 422。

### 状态 × action 矩阵（Console 置灰依据）

| action | provisioning | running | stopped | failed | deleted |
|---|---|---|---|---|---|
| start | ❌ | ❌ | ✅ | ✅ | ❌ |
| stop | ❌ | ✅ | ❌ | ❌ | ❌ |
| restart | ❌ | ✅ | ❌ | ❌ | ❌ |
| delete | 视保护 | 视保护 | ✅ | ✅ | ❌ |

（VM `termination_protection` 等细化规则见各实例模块详文。）

### 验收

```bash
# 422 INSTANCE_STATE_INVALID
curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"action":"start","idempotency_key":"'"$(uuidgen)"'"}' \
  "$BASE/api/v1/instances/$INSTANCE_ID/lifecycle"
```

---

## TASK-CORE-003 — GPU 库存只读

| 项 | 值 |
|---|---|
| 接口 | `GET /api/v1/gpu-inventory`、`GET /api/v1/gpu-inventory/occupancy` |
| operationId | `listGPUInventory`、`getGPUOccupancy` |
| RBAC | `scope:gpu-inventory:read` |
| 模块详文 | `compute/gpu-inventory-ui.md`、`home/home-gpu-utilization.md`、`compute/gpu-management.md` |
| SPEC | `tasks/modules/spec/console/compute/spec-console-gpu-inventory-ui.md` |

### Handler 契约

**listGPUInventory**

- Query：`vendor`、`state`（`available|allocated|unavailable|maintenance`）、`limit`、`cursor`。
- 成功：`200` + `GPUInventoryListResponse`。
- **只读**；无 GPU CRUD、无 POST/PATCH/DELETE。
- 租户边界：仅返回当前租户可见设备（与 BOSS 全平台视图区分）。

**getGPUOccupancy**

- 成功：`200` + `GPUOccupancySummary`（`total`、`allocated`、`available`、`avg_utilization_pct`）。

### 创建前置条件

- 无创建接口；页面聚合读能力。

### 验收

```bash
curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/gpu-inventory"
curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/gpu-inventory/occupancy"
```

---

## TASK-CORE-004 — AsyncTask 单任务查询

| 项 | 值 |
|---|---|
| 接口 | `GET /api/v1/tasks/{task_id}` |
| operationId | `getTask`（以 YAML 为准） |
| 模块详文 | `alerts/async-task-center.md` |
| 被依赖 | Services `202` 响应、`Location: /api/v1/tasks/{id}` |

### Handler 契约

- 成功：`200` + `AsyncTask`（schema 见 `components/schemas/AsyncTask`）。
- 不存在或无权：`404`。
- **未声明**：`GET /api/v1/tasks` list — 不得自造 list schema。

### 验收

```bash
curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/tasks/$TASK_ID"
```

---

## TASK-CORE-005 — 实例观测（logs / events / metrics / exec）

| 项 | 值 |
|---|---|
| 依赖 | TASK-CORE-001（实例存在） |
| 模块详文 | `compute/container-observability.md` |
| SPEC | `tasks/modules/spec/console/compute/spec-console-container-observability.md` |

| 方法 | 路径 | operationId | 成功 | RBAC |
|---|---|---|---|---|
| GET | `/instances/{instance_id}/logs` | `listInstanceLogs` | 200 + `InstanceLogListResponse` | `scope:instances:read` |
| GET | `/instances/{instance_id}/events` | `listInstanceEvents` | 200 + `InstanceEventListResponse` | `scope:instances:read` |
| GET | `/instances/{instance_id}/metrics` | `getInstanceMetrics` | 200 + `InstanceMetricsResponse` | `scope:instances:read` |
| POST | `/instances/{instance_id}/exec` | `createInstanceExecSession` | 200 + `InstanceExecSession` | `scope:instances:exec` |

### Handler 契约要点

- 实例不存在：`404`。
- exec 时实例非 running：HTTP `422`（YAML 已挂 PreconditionFailed；具体 code 待 Core 统一举例时可沿用 `INSTANCE_STATE_INVALID` 建议语义）。
- exec Request 必填：`idempotency_key`、`command[]`。
- logs/events 支持 `limit`、`cursor` 分页。

### 验收

```bash
curl -sS -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/instances/$INSTANCE_ID/logs?limit=100"
```

---

## TASK-CORE-006 — Sandbox 模板与安全事件

| 项 | 值 |
|---|---|
| 依赖 | TASK-CORE-001 |
| 模块详文 | `compute/sandbox-templates.md`、`compute/sandbox-instance-management.md` |

| 方法 | 路径 | operationId | 成功 |
|---|---|---|---|
| GET | `/api/v1/sandbox-templates` | `listSandboxTemplates` | 200 + `SandboxTemplateListResponse` |
| GET | `/instances/{instance_id}/security-events` | `listInstanceSecurityEvents` | 200 + `InstanceSecurityEventListResponse` |

### 待补（不得写成已冻结）

- Sandbox 模板 CRUD、延长/暂停专属 action — 见 sandbox 主模块待补边界。

### 验收

```bash
curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/sandbox-templates"
```

---

## TASK-CORE-007 — 网络路由

| 项 | 值 |
|---|---|
| 模块详文 | `compute/network/route.md` |

| 方法 | 路径 | operationId | 成功 |
|---|---|---|---|
| GET | `/api/v1/networks/routes` | `listNetworkRoutes` | 200 + 列表 schema |
| POST | `/api/v1/networks/routes` | `createNetworkRoute` | 201 + `NetworkRoute` |

### 待补

- `DELETE /networks/routes/{id}` — YAML 未声明。

---

## TASK-CORE-008 — 卷快照

| 项 | 值 |
|---|---|
| 依赖 | volumes CRUD 基础 handler |
| 模块详文 | `compute/storage/block-storage-snapshot.md` |

| 方法 | 路径 | operationId | 成功 |
|---|---|---|---|
| GET | `/volumes/{volume_id}/snapshots` | `listVolumeSnapshots` | 200 |
| POST | `/volumes/{volume_id}/snapshots` | `createVolumeSnapshot` | **202** + `VolumeSnapshot` |

### 前置条件（建议语义）

- 卷不存在：`404`。
- 卷状态不允许快照：HTTP `422`（具体 code 待 Core 举例）。

---

## TASK-CORE-009 — 对象存储桶 / 上传 / 下载

| 项 | 值 |
|---|---|
| 依赖 | objects 元数据 CRUD |
| 模块详文 | `compute/storage/object-storage-upload.md`、`compute/storage/object-storage.md` |

| 方法 | 路径 | operationId | 成功 |
|---|---|---|---|
| GET | `/api/v1/buckets` | `listStorageBuckets` | 200 |
| POST | `/api/v1/buckets` | `createStorageBucket` | 201 |
| POST | `/api/v1/objects/upload` | `uploadStorageObject` | **202** + `AsyncTask` |
| GET | `/api/v1/objects/{object_id}/download` | `downloadStorageObject` | 200（流或重定向，以 YAML 为准） |

### 要点

- upload 桶不存在：`422`（YAML 已声明 upload 路径）。
- 上传进度轮询：`GET /tasks/{task_id}`。

---

## TASK-CORE-010 — 文件系统挂载目标

| 项 | 值 |
|---|---|
| 接口 | `GET /api/v1/filesystems/{filesystem_id}/mount-targets` |
| operationId | `listFilesystemMountTargets` |
| 模块详文 | `compute/storage/filesystem-mount-targets.md`、`compute/storage/file-storage.md` |

### Handler 契约

- 文件系统不存在：`404`。
- 成功：`200` + mount targets 列表（schema 见 YAML）。

---

## TASK-CORE-011 — 向量写入 + search 422

| 项 | 值 |
|---|---|
| 依赖 | vector-stores CRUD |
| 模块详文 | `compute/storage/vector-store-write.md`、`compute/storage/vector-storage.md` |

| 项 | 路径 | 成功 |
|---|---|---|
| POST documents | `/vector-stores/{vector_store_id}/documents` | **202** + AsyncTask 或 YAML 定义体 |
| search 422 | `POST .../search` | 向量库非 ready → **422** PreconditionFailed |

### 验收（search 422）

- 对 `state != ready` 的 VectorStore 调用 search 返回 422。

---

## TASK-CORE-012 — K8s 工作负载 + 创建集群 422

| 项 | 值 |
|---|---|
| 依赖 | k8s-clusters CRUD |
| 模块详文 | `compute/k8s-workloads.md`、`compute/k8s-cluster-management.md` |

| 项 | 路径 | operationId |
|---|---|---|
| GET workloads | `/k8s-clusters/{cluster_id}/workloads` | `listK8sClusterWorkloads` |
| 422 验通 | `POST /k8s-clusters` | `createK8sCluster` |

### Handler 契约

- workloads：集群不存在 `404`；成功 `200` + 工作负载摘要列表。
- create 集群前置失败（配额/版本/网络未就绪等）：HTTP `422`（Phase 2 已挂 PreconditionFailed）。

---

## TASK-CORE-013 — 收尾（Branding / Metering / 遗留 schema）

| 项 | 说明 |
|---|---|
| 参考 | `docs/console-modules/governance/schema-completion-tracker.md` |
| 范围 | Branding stub 验通、`GET /metering/usage` handler、volume 挂载冲突 409/422（**409/422 须先补 YAML** 再写冻结文档） |

---

## 刻意不实现

| 路径/能力 | 原因 |
|---|---|
| Services `/gpu-containers*`、`/sandboxes*` | deprecated |
| `GET /api/v1/alerts` | YAML 未声明 |
| `GET /api/v1/tasks` list | YAML 未声明 |
| `listInstances?kind=batch_job` | query 枚举待扩展 |

---

## 文档交付清单（Core 对齐）

| 交付物 | 路径 |
|---|---|
| 本指南 | `tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` |
| Phase 2 验收记录 | `tasks/phase2/PHASE2-HANDLER-ACCEPTANCE-RECORD.md` |
| Services 指南 | `tasks/execution/SERVICES-HANDLER-IMPLEMENTATION-GUIDE.md` |
| 任务索引 | `tasks/execution/CORE-TEAM-TASKS.md` |
| Console 模块详文 | `docs/console-modules/**` |
| Console SPEC（辅助） | `tasks/modules/spec/spec-console-*.md` |
| OpenAPI | `ANI-main/repo/api/openapi/v1.yaml` |
| Schema 进度 | `docs/console-modules/governance/schema-completion-tracker.md` |
| Core-ready 复审 | `docs/console-modules/governance/core-ready-review-checklist.md` |

## 相关文件

- `docs/console-modules/governance/YAML-EXPANSION-SUMMARY-2026-06-17.md`
- `docs/console-modules/governance/module-delivery-workflow.md`
- `tasks/support/support-console-compute-modules.md`
