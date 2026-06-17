# Phase 2 Handler 文档验收记录

> **生成日期**：2026-06-17  
> **用途**：Core/Services 团队实现 Phase 2 handler 时的**文档层验收清单**（curl + 断言 + 签核）。  
> **约束**：本仓库当前阶段**不提交** ANI-main handler 代码；验收在目标环境执行后回填「Handler」列。  
> **契约指南**：Core → `tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md` · Services → `tasks/execution/SERVICES-HANDLER-IMPLEMENTATION-GUIDE.md`

---

## 验收环境（通用）


| 变量              | 说明                                        |
| --------------- | ----------------------------------------- |
| `$BASE`         | Gateway 根 URL，如 `https://api.example.com` |
| `$TEST_KEY`     | 具备对应 RBAC scope 的 API Key                 |
| `$READONLY_KEY` | 只读 Key（403 负例）                            |
| `$NOAUTH`       | 无 Key（401 负例）                             |


**通用断言**

- 错误体含 `code`、`message`、`request_id`（字符串非空）
- 2xx 响应体 schema 与 OpenAPI `components/schemas` 一致（字段类型、required）
- 列表响应含 `items` 数组；有 cursor 时 `next_cursor` 可为 null
- 租户隔离：Key A 不可见 Key B 租户资源 → `404` 或空列表（以产品口径为准，须在 PR 描述中固定）

---

## 总表（Phase 2 范围）


| TASK     | 层        | Ops | 模块详文                           | 文档验收 | YAML | Handler |
| -------- | -------- | --- | ------------------------------ | ---- | ---- | ------- |
| CORE-003 | Core     | 2   | `gpu-inventory-ui.md`          | ✅    | ✅    | ☐       |
| CORE-005 | Core     | 4   | `container-observability.md`   | ✅    | ✅    | ☐       |
| CORE-006 | Core     | 2   | `sandbox-templates.md`         | ✅    | ✅    | ☐       |
| CORE-007 | Core     | 2   | `network/route.md`             | ✅    | ✅    | ☐       |
| CORE-008 | Core     | 2   | `block-storage-snapshot.md`    | ✅    | ✅    | ☐       |
| CORE-009 | Core     | 4   | `object-storage-upload.md`     | ✅    | ✅    | ☐       |
| CORE-010 | Core     | 1   | `filesystem-mount-targets.md`  | ✅    | ✅    | ☐       |
| CORE-011 | Core     | 2   | `vector-store-write.md`        | ✅    | ✅    | ☐       |
| CORE-012 | Core     | 2   | `k8s-workloads.md`             | ✅    | ✅    | ☐       |
| SVC-013  | Services | 3   | `inference-observability.md` 等 | ✅    | ✅    | ☐       |
| SVC-014  | Services | 3   | `kb-`* 子模块                     | ✅    | ✅    | ☐       |
| SVC-015  | Services | 3   | `tenant-*` 子模块                 | ✅    | ✅    | ☐       |
| SVC-016  | Services | 3   | `integration-*.md`             | ✅    | ✅    | ☐       |


> CORE-001/002/004 为 Sprint 1 基线，见 `CORE-HANDLER-IMPLEMENTATION-GUIDE`；本表聚焦 Phase 2 **+19 Core / +12 Services** 新增 ops。

**文档验收 ✅ 定义**：模块详文 + SPEC + 本记录含 curl 与负例；`schema-completion-tracker.md` 已链到 TASK。

---

## CORE-003 — GPU 库存

**RBAC**：`scope:gpu-inventory:read`

### 正向

```bash
# T1 list
curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/gpu-inventory?limit=20"
# 期望 HTTP 200；body.items 为数组；items[].state ∈ available|allocated|unavailable|maintenance

# T2 occupancy
curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/gpu-inventory/occupancy"
# 期望 HTTP 200；含 total/allocated/available（以 GPUOccupancySummary 为准）
```

### 负向


| 用例                   | 期望  |
| -------------------- | --- |
| 无 Key                | 401 |
| 无 gpu-inventory:read | 403 |
| `state=invalid_enum` | 400 |


### 签核

- [ ] T1–T2 通过
- [ ] 无 POST/PATCH/DELETE GPU 写路径被误实现

---

## CORE-005 — 实例观测

**前置**：存在 `$INSTANCE_ID`（任意 kind，当前租户可见）

### 正向

```bash
curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/instances/$INSTANCE_ID/logs?limit=100"

curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/instances/$INSTANCE_ID/events?limit=50"

curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/instances/$INSTANCE_ID/metrics"
```

### exec（写操作）

```bash
curl -sS -w "\nHTTP:%{http_code}\n" -X POST \
  -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","command":["echo","ok"]}' \
  "$BASE/api/v1/instances/$INSTANCE_ID/exec"
# running 实例 → 200 + InstanceExecSession
# stopped 实例 → 422（PreconditionFailed；code 以 YAML 为准）
```

### 负向


| 用例                     | 期望  |
| ---------------------- | --- |
| 伪造 instance_id         | 404 |
| exec 缺 idempotency_key | 400 |


### 签核

- [ ] logs/events/metrics/exec 四类路径验通
- [ ] 422 非 running exec 已验通

---

## CORE-006 — Sandbox 模板与安全事件

```bash
curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/sandbox-templates"

curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/instances/$SANDBOX_INSTANCE_ID/security-events"
```

### 签核

- [ ] 模板 list 200
- [ ] 安全事件 200（sandbox 实例）
- [ ] 未实现模板 CRUD

---

## CORE-007 — 网络路由

```bash
curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/networks/routes"

curl -sS -w "\nHTTP:%{http_code}\n" -X POST \
  -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","name":"rt-test","destination_cidr":"0.0.0.0/0","next_hop":"10.0.0.1"}' \
  "$BASE/api/v1/networks/routes"
# 字段以 CreateNetworkRouteRequest 为准 → 201
```

### 签核

- [ ] list + create 验通
- [ ] DELETE route 未自造（YAML 未声明）

---

## CORE-008 — 卷快照

**前置**：`$VOLUME_ID` 存在

```bash
curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/volumes/$VOLUME_ID/snapshots"

curl -sS -w "\nHTTP:%{http_code}\n" -X POST \
  -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","name":"snap-accept"}' \
  "$BASE/api/v1/volumes/$VOLUME_ID/snapshots"
# 期望 202 + VolumeSnapshot（或 YAML 声明体）
```

### 签核

- [ ] list 200；create 202
- [ ] 卷不存在 → 404

---

## CORE-009 — 桶 / 上传 / 下载

```bash
curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/buckets"

# upload → 202 AsyncTask；轮询 GET /tasks/{id}
curl -sS -w "\nHTTP:%{http_code}\n" -X POST \
  -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","bucket":"b1","key":"k1"}' \
  "$BASE/api/v1/objects/upload"

curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/objects/$OBJECT_ID/download"
```

### 签核

- [ ] upload 202 + task 可查询
- [ ] 桶不存在 upload → 422（YAML 已声明 upload 路径）

---

## CORE-010 — 文件系统挂载目标

```bash
curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/filesystems/$FILESYSTEM_ID/mount-targets"
```

### 签核

- [ ] 200 + 列表 schema
- [ ] filesystem 不存在 → 404

---

## CORE-011 — 向量写入 + search 422

```bash
curl -sS -w "\nHTTP:%{http_code}\n" -X POST \
  -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","documents":[{"id":"d1","text":"hello"}]}' \
  "$BASE/api/v1/vector-stores/$VS_ID/documents"

curl -sS -w "\nHTTP:%{http_code}\n" -X POST \
  -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"query":"test","limit":5}' \
  "$BASE/api/v1/vector-stores/$VS_NOT_READY_ID/search"
# state != ready → 422
```

### 签核

- [ ] documents POST 202（或 YAML 体）
- [ ] search 422 已验通

---

## CORE-012 — K8s 工作负载 + 创建集群 422

```bash
curl -sS -w "\nHTTP:%{http_code}\n" -H "X-API-Key: $TEST_KEY" \
  "$BASE/api/v1/k8s-clusters/$CLUSTER_ID/workloads"

# 422：故意触发前置失败（配额/无效 version 等，以环境 fixture 为准）
curl -sS -w "\nHTTP:%{http_code}\n" -X POST \
  -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","name":"bad-cluster","version":"invalid"}' \
  "$BASE/api/v1/k8s-clusters"
```

### 签核

- [ ] workloads 200
- [ ] create 集群 422 至少 1 条用例

---

## SVC-013 — 推理 logs / test / policies

详见 `tasks/execution/SERVICES-HANDLER-IMPLEMENTATION-GUIDE.md` §SVC-013。

### 签核

- [ ] logs 200 + cursor
- [ ] test 对 running 服务 200
- [ ] policies PUT 200 + idempotency_key

---

## SVC-014 — 知识库 citations / sessions / permissions

### 签核

- [ ] citations list 200
- [ ] sessions list 200
- [ ] permissions PUT 200；缺 idempotency_key → 400

---

## SVC-015 — 租户 member / role / webhook deliveries

### 签核

- [ ] getTenantMember 200
- [ ] updateTenantRole 200
- [ ] listWebhookDeliveries 200 + 分页

---

## SVC-016 — integrations + bots

### 签核

- [ ] listIntegrations 200
- [ ] createIntegration 201
- [ ] createIntegrationBot 201

---

## 文档层收口签核（Phase 2）


| 项                                    | 状态                    |
| ------------------------------------ | --------------------- |
| Core 指南含 CORE-003～012                | ✅                     |
| Services 指南含 SVC-013～016             | ✅                     |
| 本记录含全部 Phase 2 TASK curl             | ✅                     |
| `schema-completion-tracker.md` 链到本记录 | ✅                     |
| Handler 运行时验通                        | ☐（待 Core/Services PR） |


---

## 相关文件

- `tasks/execution/CORE-HANDLER-IMPLEMENTATION-GUIDE.md`
- `tasks/execution/SERVICES-HANDLER-IMPLEMENTATION-GUIDE.md`
- `docs/console-modules/governance/schema-completion-tracker.md`
- `tasks/execution/TASK-DEPENDENCY-MAP.md`

