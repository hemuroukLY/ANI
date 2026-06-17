# SPEC: Console 角色权限编辑

> Revised: 2026-06-17

## 2. Frozen Facts

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/svc/tenant/roles` | `listTenantRoles` | `200` + TenantRole items | `401` |
| PUT | `/api/v1/svc/tenant/roles/{role_id}` | `updateTenantRole` | `200 + TenantRole` | `400`,`401`,`403`,`404` |

UpdateTenantRoleRequest: idempotency_key, permissions[] (minItems 1)

## 3. TODO-YAML

- 角色 POST/DELETE

## 4. References

- `docs/console-modules/tenant/role-permission-edit.md`
