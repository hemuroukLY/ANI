# M2.2-AUTH-K: Auth Integration Profile

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 2：ANI Gateway` / `2.2 认证授权`.

## Scope

Added auth-service integration coverage for the OIDC and refresh-token path:

- `services/auth-service/internal/service/auth_integration_test.go`

The test exercises:

- `BeginOIDCLogin`
- `CompleteOIDCLogin`
- `RefreshToken`
- `ValidateToken`

The test uses fake OIDC exchange/verifier/session stores and the real JWT issuer/validator/blocklist composition, so it validates the service boundary without requiring an external Dex deployment.

## Validation

Passed:

```bash
make validate-infra
make test
make build
git diff --check
```

## Version Impact

Expected release impact: `PATCH`.

Reason:
- Adds test coverage only; no API, schema, deployment, or runtime behavior contract changes.

## Next Recommended Slice

- `M1-INFRA-F`: GPU scheduling preflight/e2e hardening.
- Or move to `模块 3` model metadata and object-storage boundary after M1/M2 baseline acceptance.
