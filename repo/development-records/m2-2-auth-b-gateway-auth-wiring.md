# M2.2-AUTH-B Development Record

Date: 2026-05-11

Product plan mapping:
- `ANI-06 / 模块 2：ANI Gateway（统一 Web Server 层） / 2.2 认证授权`
- Slice: `M2.2-AUTH-B`
- Scope: wire ANI Gateway Auth/RBAC middleware to auth-service over gRPC.

Version impact:
- Current line: `v0.x` development.
- Release impact: `MINOR`.
- Reason: changes gateway runtime behavior to call auth-service and adds additive configuration via `AUTH_SERVICE_ADDR`.
- No formal version tag was created.

## Implemented

- Added gateway `AuthClient` abstraction and gRPC implementation.
- `AUTH_SERVICE_ADDR` config is supported, defaulting to `127.0.0.1:9101`.
- Updated Auth middleware:
  - Preserves public path bypass.
  - Preserves `ANI_AUTH_MODE=dev`.
  - Calls `auth.v1.AuthService.ValidateToken` for Bearer tokens and API key headers.
  - Fails closed when auth-service is unavailable.
- Updated RBAC middleware:
  - Calls `auth.v1.AuthService.CheckPermission`.
  - Infers conservative resource/action from HTTP method and `/api/v1/{resource}` path.
  - Fails closed when auth-service is unavailable or denies permission.
- Reuses one auth-service client for Auth and RBAC middleware registration.
- Added focused permission inference tests.

## Files Added

- `services/ani-gateway/internal/middleware/auth_client.go`
- `services/ani-gateway/internal/middleware/rbac_test.go`

## Files Updated

- `services/ani-gateway/go.mod`
- `services/ani-gateway/internal/middleware/auth.go`
- `services/ani-gateway/internal/middleware/rbac.go`
- `services/ani-gateway/internal/middleware/chain.go`
- `development-records/README.md`

## Validation

Commands run:

```bash
cd repo
GOCACHE=/private/tmp/ani-go-build GOMODCACHE="$PWD/.cache/gomod" go test ./services/ani-gateway/...
make gen-proto
make validate-infra
make test
make build
```

Result:
- Gateway package tests: passed.
- `make gen-proto`: passed.
- `make validate-infra`: passed.
- `make test`: passed.
- `make build`: passed.

## Remaining Risks

- Gateway-to-auth-service transport still lacks mTLS; this must be closed in deployment/security slices.
- Route-to-resource/action mapping is conservative and not final OPA policy metadata.
- API key validation is routed through `ValidateToken`, but auth-service currently only implements JWT validation; API key lookup remains deferred.
- gRPC client currently dials with insecure transport for local/internal baseline only.

## Next Boundary

Recommended next M2 slice:
- `M2.2-AUTH-C`: API key storage/lookup design with RLS-safe global key hash lookup, or Dex/OIDC login integration after identity provider deployment strategy is fixed.
