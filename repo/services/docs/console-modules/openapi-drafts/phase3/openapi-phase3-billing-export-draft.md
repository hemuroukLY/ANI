# Phase 3 — 账单/用量导出 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`tenant/billing-export.md`  
> **TASK**：`TASK-SVC-018` 子项 / Phase 3 §4  
> **原则**：扩展 **用量报表**；大 export **202 + AsyncTask**；结算/发票归 BOSS。

---

## 1. 与 usage-report / BOSS 结算的边界

| 维度 | 租户用量报表 | 账单导出（本草案） | BOSS 结算 |
|---|---|---|---|
| 路径 | Core `GET /metering/usage` | `POST /svc/billing/exports` | 平台 |
| 职责 | 在线查询/图表 | **异步文件导出**（CSV/PDF） | 发票、对账、金额 |
| 金额 | 不展示 | export 不含结算金额（**评审确认**） | BOSS |
| 模块 | `usage-report.md` | 本文 | 非 Console |

Console：从用量页「导出」进入；export 仅含用量明细，非正式发票。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [Billing]
x-ani-rbac-scope:
  read:  scope:billing:read
  write: scope:billing:write
```

与 metering 读权限关系 — **评审与 tenant admin 对齐**（可复用 `scope:metering:read` + export write）。

---

## 3. Schemas（草案）

### CreateBillingExportRequest

```yaml
CreateBillingExportRequest:
  type: object
  required: [idempotency_key, export_type, start_time, end_time]
  properties:
    idempotency_key: { type: string, format: uuid }
    export_type:
      type: string
      enum: [usage_csv, usage_pdf, usage_detail_csv]
    start_time:      { type: string, format: date-time }
    end_time:        { type: string, format: date-time }
    resource_type:   { type: string, nullable: true, description: 与 metering filter 一致 }
    model_id:        { type: string, format: uuid, nullable: true }
    group_by:        { type: string, nullable: true }
    locale:          { type: string, nullable: true, default: zh-CN }
```

### BillingExportJob

```yaml
BillingExportJob:
  type: object
  required: [id, export_type, status, created_at]
  properties:
    id:            { type: string, format: uuid }
    export_type:   { type: string, enum: [usage_csv, usage_pdf, usage_detail_csv] }
    status:        { type: string, enum: [pending, processing, ready, failed, expired] }
    start_time:    { type: string, format: date-time }
    end_time:      { type: string, format: date-time }
    task_id:       { type: string, format: uuid, nullable: true }
    download_url:  { type: string, format: uri, nullable: true }
    file_size_bytes: { type: integer, nullable: true }
    expires_at:    { type: string, format: date-time, nullable: true }
    error_message: { type: string, nullable: true }
    created_at:    { type: string, format: date-time }
    completed_at:  { type: string, format: date-time, nullable: true }
```

### BillingExportListResponse

```yaml
BillingExportListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/BillingExportJob' } }
    next_cursor: { type: string, nullable: true }
```

---

## 4. Operations（草案）

### 4.1 `POST /api/v1/svc/billing/exports`

- operationId: `createBillingExport`
- Body: `CreateBillingExportRequest`
- 成功：`202 + AsyncTask`（`task_type` 建议 `billing.export`）+ body 可选嵌套 `BillingExportJob`
- 小文件同步：`201 + BillingExportJob`（status=ready）— **评审二选一默认 202**
- 错误：
  - `400`：时间窗非法
  - `401`、`403`
  - `422`：时间窗过大 — `EXPORT_WINDOW_TOO_LARGE`；并发 export 上限 — `EXPORT_RATE_LIMITED`

### 4.2 `GET /api/v1/svc/billing/exports/{export_id}`

- operationId: `getBillingExport`
- 成功：`200 + BillingExportJob`
- 错误：`404`

### 4.3 `GET /api/v1/svc/billing/exports`（3b）

- operationId: `listBillingExports`
- Query: `status`、`limit`、`cursor`
- 成功：`200 + BillingExportListResponse`

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 租户管理员 | 说明 |
|---|---|---|---|
| 在线用量 | 可用 | 可用 | 跳转 usage-report |
| 创建 export | 不可用 | 可用 | POST |
| 下载 ready 文件 | 不可用 | 可用 | download_url |
| 查看 export 历史 | 不可用 | 可用 | list（3b） |
| 发票/对账 | 不可用 | 不可用 | BOSS |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -X POST -H "Content-Type: application/json" -H "X-API-Key: $TEST_KEY" \
  -d '{"idempotency_key":"'"$(uuidgen)"'","export_type":"usage_csv","start_time":"2026-06-01T00:00:00Z","end_time":"2026-06-15T23:59:59Z"}' \
  "$BASE/api/v1/svc/billing/exports"

curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/billing/exports/$EXPORT_ID"
curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/tasks/$TASK_ID"
```

---

## 7. 评审检查清单

- [ ] export 不含 BOSS 结算金额（或明确 export_type 边界）
- [ ] 大 export 必须 202 + AsyncTask
- [ ] 与 `metering/usage` query 参数对齐
- [ ] 合入后更新 `billing-export.md`

---

## 相关文件

- `docs/console-modules/tenant/billing-export.md`
- `docs/console-modules/tenant/usage-report.md`
- `docs/console-modules/alerts/async-task-center.md`
