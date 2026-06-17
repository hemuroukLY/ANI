# Phase 3 — 合规能力 OpenAPI 详化草案

> **状态**：Phase 3 **规划详化**（非 ANI-main 冻结事实）  
> **模块详文**：`security/compliance.md`  
> **TASK**：`TASK-SVC-018` 子项 / Phase 3 §4  
> **原则**：Console **只读**；完整合规包导出与取证归 **BOSS**。

---

## 1. 与审计日志 / BOSS 合规的边界

| 维度 | 平台审计日志 | 合规摘要（本草案） | BOSS |
|---|---|---|---|
| 路径 | 规划 `/audit/events` | `/svc/compliance/*` | 平台运营 |
| 内容 | 操作事件流 | 控制项状态、报告元数据 | 全量导出/取证 |
| 写操作 | — | **无**（Console） | 报告生成 |
| 模块 | `audit-log.md` | 本文 | 非 Console API |

租户可见：**摘要 + 历史报告 list**；`download_url` 由 BOSS 签发只读链接（短期有效）。

---

## 2. 建议 tag 与 RBAC

```yaml
tags: [Compliance]
x-ani-rbac-scope:
  read: scope:compliance:read
```

无 write scope（Console 租户侧）。

---

## 3. Schemas（草案）

### ComplianceControlItem

```yaml
ComplianceControlItem:
  type: object
  required: [control_id, title, status]
  properties:
    control_id:      { type: string, description: 如 SOC2-CC6.1 }
    title:           { type: string }
    status:          { type: string, enum: [pass, fail, unknown, not_applicable] }
    last_checked_at: { type: string, format: date-time, nullable: true }
    remediation_hint: { type: string, nullable: true }
```

### ComplianceSummary

```yaml
ComplianceSummary:
  type: object
  required: [overall_status, items, evaluated_at]
  properties:
    overall_status: { type: string, enum: [compliant, non_compliant, partial, unknown] }
    framework:      { type: string, nullable: true, description: 如 iso27001, 等保2.0 }
    items:
      type: array
      items: { $ref: '#/components/schemas/ComplianceControlItem' }
    open_actions:   { type: integer, minimum: 0 }
    evaluated_at:   { type: string, format: date-time }
```

### ComplianceReport

```yaml
ComplianceReport:
  type: object
  required: [id, report_type, period_start, period_end, status, created_at]
  properties:
    id:            { type: string, format: uuid }
    report_type:   { type: string, enum: [periodic_summary, audit_pack, retention_attestation] }
    period_start:  { type: string, format: date-time }
    period_end:    { type: string, format: date-time }
    status:        { type: string, enum: [generating, ready, expired, failed] }
    download_url:  { type: string, format: uri, nullable: true, description: BOSS 签发；ready 时有效 }
    expires_at:    { type: string, format: date-time, nullable: true }
    created_at:    { type: string, format: date-time }
```

### ComplianceReportListResponse

```yaml
ComplianceReportListResponse:
  type: object
  required: [items]
  properties:
    items:       { type: array, items: { $ref: '#/components/schemas/ComplianceReport' } }
    next_cursor: { type: string, nullable: true }
```

---

## 4. Operations（草案）

### 4.1 `GET /api/v1/svc/compliance/summary`

- operationId: `getComplianceSummary`
- Query: `framework`（optional）
- 成功：`200 + ComplianceSummary`
- 错误：`401`、`403`

### 4.2 `GET /api/v1/svc/compliance/reports`

- operationId: `listComplianceReports`
- Query: `report_type`、`status`、`limit`、`cursor`
- 成功：`200 + ComplianceReportListResponse`

### 4.3 `GET /api/v1/svc/compliance/reports/{report_id}`

- operationId: `getComplianceReport`
- 成功：`200 + ComplianceReport`
- 错误：`404`

> **无** `POST /compliance/reports` — 报告生成请求走 BOSS；Console 仅展示 BOSS 同步至租户可见 catalog。

---

## 5. Console 操作可用性矩阵（草案）

| 操作 | 只读用户 | 租户管理员 | 说明 |
|---|---|---|---|
| 合规摘要 | 可用 | 可用 | GET summary |
| 报告列表/详情 | 可用 | 可用 | GET reports |
| 下载 ready 报告 | 可用 | 可用 | download_url 跳转 |
| 请求新报告 | 不可用 | 不可用 | BOSS / 工单 |

---

## 6. Handler 验收 curl（合入 YAML 后）

```bash
curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/compliance/summary"
curl -sS -H "X-API-Key: $TEST_KEY" "$BASE/api/v1/svc/compliance/reports?limit=20"
```

---

## 7. 评审检查清单

- [ ] Console 无 write path
- [ ] download_url 不暴露 BOSS 内部 admin API
- [ ] 与 audit-log 导出边界清晰
- [ ] 合入后更新 `compliance.md`

---

## 相关文件

- `docs/console-modules/security/compliance.md`
- `docs/console-modules/security/audit-log.md`
- `docs/console-modules/tenant/tenant-management.md`
