# SPEC: Console billing-export

> Phase 3 规划详化 · Source: `tasks/modules/prd/console/tenant/prd-console-billing-export.md`  
> 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-billing-export-draft.md` · Revised: 2026-06-17

## 1. Summary

Services 用量 export（规划）；202 AsyncTask；Core metering 只读。

## 2. Planned Facts（非冻结）

| Method | Path | operationId | 成功 |
|---|---|---|---|
| POST | `/api/v1/svc/billing/exports` | `createBillingExport` | 202 |
| GET | `.../exports/{export_id}` | `getBillingExport` | 200 |
| GET | `/api/v1/svc/billing/exports` | `listBillingExports` | 200（3b） |

Schemas: `CreateBillingExportRequest`, `BillingExportJob`.  
task_type: `billing.export`.

## 3. 主维护源

- `docs/console-modules/tenant/billing-export.md`

OpenAPI **未合入** ≠ handler 已实现。
