# SPEC: Console LDAP 配置

> Revised: 2026-06-17

## 2. Frozen Facts（OIDC 引用）

| Method | Path | operationId | 成功 | 错误 |
|---|---|---|---|---|
| GET | `/api/v1/svc/tenant/sso` | `getSsoConfig` | `200 + SsoConfig` | `401`,`403` |
| PUT | `/api/v1/svc/tenant/sso` | `updateSsoConfig` | `200 + SsoConfig` | `400`,`403` |

## 3. TODO-YAML

- `/tenant/ldap` GET/PUT
- LDAP sync / group mapping

## 4. References

- `docs/console-modules/tenant/ldap-config.md`
