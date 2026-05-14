# M2.2-AUTH-G: OIDC Login Boundary

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 2：ANI Gateway` / `2.2 认证授权`.

## Scope

Added the first Dex/OIDC login boundary without pretending code exchange is complete:

- Added `BeginOIDCLogin` RPC.
- Added `CompleteOIDCLogin` RPC.
- Added OIDC config:
  - `AUTH_OIDC_ISSUER_URL`
  - `AUTH_OIDC_CLIENT_ID`
  - `AUTH_OIDC_AUTH_URL`
- `BeginOIDCLogin` creates a random state, stores it in `CacheStore`, and returns a Dex/OIDC authorization URL.
- `CompleteOIDCLogin` validates state and redirect URI, then returns `UNIMPLEMENTED` until Dex code exchange and ID token verification are added.
- Added focused tests for state storage, URL construction, and redirect URI mismatch rejection.

## Validation

Passed:

```bash
make gen-proto
make validate-architecture
make test
make build
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds new auth-service RPCs and generated protobuf code in a backward-compatible way.

## Next Recommended Slice

- `M2.2-AUTH-H`: implement Dex code exchange, ID token verification, user/role mapping, and refresh-token row creation after successful OIDC login.
