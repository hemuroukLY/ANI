# SPEC: Console netsec-policy

> Phase 3 规划详化 · Source: `tasks/modules/prd/console/security/prd-console-netsec-policy.md`  
> 草案: `docs/console-modules/openapi-drafts/phase3/openapi-phase3-netsec-policy-draft.md` · Revised: 2026-06-17

## 1. Summary

Core `/networks/policies` CRUD（规划）；与 security-groups 分资源。

## 2. Planned Facts（非冻结）

| Method | Path | operationId |
|---|---|---|
| GET/POST | `/api/v1/networks/policies` | list/create |
| GET/PUT/DELETE | `.../policies/{policy_id}` | get/update/delete |

Schemas: `NetworkSecurityPolicy`, `CreateNetworkSecurityPolicyRequest`.

## 3. 主维护源

- `docs/console-modules/security/netsec-policy.md`

OpenAPI **未合入** ≠ handler 已实现。
