# M2.2-AUTH-D: JWT Token Revocation

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 2：ANI Gateway` / `2.2 认证授权`.

## Scope

Implemented the first stable part of token lifecycle management:

- `RevokeToken` validates `jti`.
- Writes `jwt:blocklist:{jti}` into `ports.CacheStore`.
- JWT validation checks the same blocklist key through `CacheStore`.
- Added focused tests for revocation and blocked JWT rejection.

`Login` and `RefreshToken` remain explicitly unimplemented until Dex/OIDC callback and refresh-token storage are designed.

## Validation

Passed:

```bash
make test
make build
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Implements previously declared auth-service RPC behavior without breaking existing API contracts.

## Remaining Risks

- Revocation TTL currently follows the default access-token TTL of 1 hour.
- Refresh-token revocation and OIDC session handling are deferred to `M2.2-AUTH-E`.
