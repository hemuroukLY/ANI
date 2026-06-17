# SPEC: Console compliance

> Phase 3 规划详化 · Source: `tasks/modules/prd/console/security/prd-console-compliance.md`  
> 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-compliance-draft.md` · Revised: 2026-06-17

## 1. Summary

Services 合规只读 GET（规划）；无 Console write。

## 2. Planned Facts（非冻结）

| Method | Path | operationId |
|---|---|---|
| GET | `/api/v1/svc/compliance/summary` | `getComplianceSummary` |
| GET | `/api/v1/svc/compliance/reports` | `listComplianceReports` |
| GET | `.../reports/{report_id}` | `getComplianceReport` |

Schemas: `ComplianceSummary`, `ComplianceReport`.

## 3. 主维护源

- `docs/console-modules/security/compliance.md`

OpenAPI **未合入** ≠ handler 已实现。
