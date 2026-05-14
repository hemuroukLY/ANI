# M2.2-AUTH-F: Refresh Token Foundation

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 2：ANI Gateway` / `2.2 认证授权`.

## Scope

Implemented the refresh-token foundation needed before Dex/OIDC callback integration:

- Added `refresh_tokens` table to the initial schema.
- Added RLS policy for `refresh_tokens`.
- Added JWT private-key config:
  - `AUTH_JWT_PRIVATE_KEY_PEM`
  - `AUTH_JWT_PRIVATE_KEY_FILE`
- Added RS256 `JWTIssuer`.
- Added persistent refresh token validation.
- Implemented `RefreshToken` RPC to validate a refresh token and issue a new AccessToken.
- Added focused tests proving issued AccessToken is valid and carries tenant/user/roles claims.

`Login` remains unimplemented because the current Proto has username/password semantics while the product design requires Dex/OIDC callback. The next slice should define that API boundary explicitly.

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
- Adds auth behavior and schema for a planned token lifecycle capability.
- Does not break existing Proto fields.

## Next Recommended Slice

- `M2.2-AUTH-G`: define and implement Dex/OIDC callback boundary that creates users/session records and inserts refresh token rows.
