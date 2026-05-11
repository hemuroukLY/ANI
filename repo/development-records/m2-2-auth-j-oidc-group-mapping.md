# M2.2-AUTH-J: OIDC Group Mapping

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 2：ANI Gateway` / `2.2 认证授权`.

## Scope

Hardened OIDC group-to-role mapping:

- Added `AUTH_OIDC_GROUP_ROLE_MAP_JSON`.
- Default behavior grants only `user`, even if external groups contain names like `tenant-admin` or `platform-admin`.
- High-privilege ANI roles require explicit allowlist mapping.
- Invalid mapped roles are ignored.
- Added tests for default least-privilege behavior and explicit enterprise mappings.

Example:

```json
{
  "/corp/ani-admins": ["tenant-admin"],
  "CN=ANI-Auditors": ["auditor"]
}
```

## Validation

Passed:

```bash
make validate-architecture
make test
make build
make validate-infra
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds security-relevant OIDC mapping configuration and hardens default role behavior.

## Next Recommended Slice

- `M1-INFRA-E`: GPU scheduling baseline.
- Or `M2.2-AUTH-K`: integration/e2e profile for OIDC, refresh token, and gateway auth wiring.
