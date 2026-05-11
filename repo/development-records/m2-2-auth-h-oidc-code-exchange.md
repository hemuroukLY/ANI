# M2.2-AUTH-H: OIDC Code Exchange

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 2：ANI Gateway` / `2.2 认证授权`.

## Scope

Completed the first runtime OIDC callback path:

- Added OIDC token endpoint config:
  - `AUTH_OIDC_TOKEN_URL`
  - `AUTH_OIDC_CLIENT_SECRET`
  - `AUTH_OIDC_PUBLIC_KEY_PEM`
  - `AUTH_OIDC_PUBLIC_KEY_FILE`
- `CompleteOIDCLogin` now validates state and redirect URI, exchanges the authorization code, verifies an RS256 ID token with a configured static public key, creates/updates the OIDC user, grants mapped roles, creates a refresh token, and returns a `TokenPair`.
- Added bounded PostgreSQL session persistence for OIDC user/session creation.
- Added tests for successful OIDC callback token-pair issuance.

## Design Boundary

This slice intentionally avoids adding a third-party OIDC library or network-dependent JWKS discovery.

Current verifier:
- Static RS256 public key.
- Validates issuer, audience, expiry, subject, and email.

Deferred:
- Dex JWKS discovery.
- JWKS key rotation.
- Customer-specific group mapping rules.
- SAML-specific mapping.

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
- Adds runtime behavior for previously declared OIDC login RPCs and additional configuration.

## Next Recommended Slice

- `M2.2-AUTH-I`: implement JWKS discovery and hardened enterprise group mapping.
