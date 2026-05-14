# M2.2-AUTH-I: OIDC JWKS Discovery

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 2：ANI Gateway` / `2.2 认证授权`.

## Scope

Improved OIDC ID token verification:

- Added `AUTH_OIDC_JWKS_URL`.
- OIDC verifier now prefers JWKS-based verification when configured.
- Selects RSA public keys by JWT header `kid`.
- Validates RS256 signature, issuer, audience, expiry, subject, and email.
- Keeps static RS256 public-key verification as an offline fallback.
- Added a no-listener unit test using a custom `http.RoundTripper` for JWKS responses.

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
- Adds OIDC verifier configuration and behavior without breaking existing API fields.

## Remaining Work

- `M2.2-AUTH-J`: enterprise group mapping hardening and configurable role mapping.
