# P0 OpenAPI 扩充草案（文档层 · 非冻结）

> **状态**：**待定** — 暂缓评审合入 ANI-main；优先 Phase 2 handler 验收与 Phase 3 Agent Session 详化。  
> **TASK**：TASK-CORE-014 / TASK-SVC-017 同步标记为 **待定**。  
> **用途**：供 Core/Services 后续评审；**不得**写入 ANI-main 直至评审通过。

---

## 1. 告警事件只读列表（P0-01）

**模块详文**：`alerts/alerts-pending-items.md`  
**TASK**：`TASK-CORE-014`（见 `tasks/execution/CORE-TEAM-TASKS.md`）

### 建议归属

| 方案 | 路径前缀 | 说明 |
|---|---|---|
| **A（推荐）** | Core `/api/v1/observability/alerts` | 与 `alert-rules` 同 tag，区分规则 vs 事件 |
| B | Core `/api/v1/alerts` | 更短，但与 branding 占位易混淆 |
| C | 仅 Console 聚合 | 不扩 YAML；维持现状 |

### 草案 operation（方案 A）

```yaml
GET /api/v1/observability/alerts:
  operationId: listObservabilityAlerts
  summary: 列出当前租户可见的告警事件
  tags: [Observability]
  security: [{ ApiKeyAuth: [] }]
  parameters:
    - name: severity
      in: query
      schema: { type: string, enum: [critical, warning, info] }
    - name: state
      in: query
      schema: { type: string, enum: [firing, resolved, acknowledged] }
    - name: limit
      in: query
      schema: { type: integer, default: 50, maximum: 200 }
    - name: cursor
      in: query
      schema: { type: string }
  responses:
    '200':
      content:
        application/json:
          schema: { $ref: '#/components/schemas/ObservabilityAlertListResponse' }
    '401': { $ref: '#/components/responses/Unauthorized' }
    '403': { $ref: '#/components/responses/Forbidden' }
```

### 建议 schema（草案）

| Schema | 关键字段 |
|---|---|
| `ObservabilityAlert` | `id`, `severity`, `state`, `title`, `source_module`, `resource_type`, `resource_id`, `fired_at`, `resolved_at` |
| `ObservabilityAlertListResponse` | `items[]`, `next_cursor` |

### 待办聚合（可选 Phase 2）

若产品需要独立「待处理事项」资源，建议 **Services** 侧：

`GET /api/v1/svc/pending-items` — 与告警事件分 schema，避免混写。

---

## 2. 异步任务列表（P0-02）

**模块详文**：`alerts/async-task-center.md`  
**TASK**：`TASK-CORE-014`

### 现状

- 已冻结：`GET /api/v1/tasks/{task_id}` → `AsyncTask`
- 缺失：`GET /api/v1/tasks`（cursor 分页）

### 草案 operation

```yaml
GET /api/v1/tasks:
  operationId: listTasks
  summary: 分页列出当前租户异步任务
  tags: [Tasks]
  parameters:
    - name: status
      in: query
      schema:
        type: string
        enum: [pending, running, completed, failed, cancelled, dead_letter]
    - name: task_type
      in: query
      schema: { type: string }
    - name: resource_type
      in: query
      schema: { type: string }
    - name: limit
      in: query
      schema: { type: integer, default: 50, maximum: 200 }
    - name: cursor
      in: query
      schema: { type: string }
  responses:
    '200':
      content:
        application/json:
          schema: { $ref: '#/components/schemas/AsyncTaskListResponse' }
```

### 建议 schema

| Schema | 说明 |
|---|---|
| `AsyncTaskListResponse` | `items: AsyncTask[]`, `next_cursor` |

### 可选扩展（非 P0 阻塞）

- `POST /api/v1/tasks/{task_id}/cancel` — 取消任务；须定义幂等与终态冲突 `409`

---

## 3. 推理服务事件（P0-04）

**模块详文**：`inference/inference-observability.md`  
**TASK**：`TASK-SVC-017`

### 现状

- 已冻结：`GET /api/v1/svc/inference-services/{service_id}/logs`
- 缺失：`GET .../events`

### 草案 operation（Services）

```yaml
GET /api/v1/svc/inference-services/{service_id}/events:
  operationId: listInferenceServiceEvents
  summary: 列出推理服务运维事件
  tags: [InferenceServices]
  parameters:
    - $ref: '#/components/parameters/ServiceId'
    - name: limit
      in: query
      schema: { type: integer, default: 100, maximum: 500 }
    - name: cursor
      in: query
      schema: { type: string }
  responses:
    '200':
      content:
        application/json:
          schema: { $ref: '#/components/schemas/InferenceServiceEventListResponse' }
    '404': { $ref: '#/components/responses/NotFound' }
```

### 建议 schema

| Schema | 关键字段 |
|---|---|
| `InferenceServiceEvent` | `id`, `event_type`, `message`, `occurred_at`, `metadata` |
| `InferenceServiceEventListResponse` | `items[]`, `next_cursor` |

---

## 4. 评审检查清单

评审通过前，每项须确认：

- [ ] `operationId` 全局唯一
- [ ] 中文 `summary` + `tags` + `security`
- [ ] 2xx 响应挂完整 schema（非空 object）
- [ ] 401/403/404 引用统一 responses
- [ ] 列表 operation 含 `limit`/`cursor`
- [ ] Console 模块详文 `TODO-YAML` 改为「YAML 已声明」或保留待补说明
- [ ] `schema-completion-tracker.md` 与 TASK 清单同步

---

## 相关文件

- `docs/console-modules/governance/console-undefined-features-backlog.md` §建议执行顺序 #1
- `tasks/execution/CORE-TEAM-TASKS.md` — TASK-CORE-014
- `tasks/execution/SERVICES-TEAM-TASKS.md` — TASK-SVC-017
- `docs/console-modules/governance/schema-completion-tracker.md`
