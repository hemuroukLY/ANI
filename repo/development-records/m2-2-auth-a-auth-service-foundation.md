# M2.2-AUTH-A Development Record

Date: 2026-05-11

Product plan mapping:
- `ANI-06 / 模块 2：ANI Gateway（统一 Web Server 层） / 2.2 认证授权`
- Slice: `M2.2-AUTH-A`
- Scope: minimal auth-service foundation for internal token validation and RBAC checks.

Version impact:
- Current line: `v0.x` development.
- Release impact: `MINOR`.
- Reason: adds a new Go service and build target; no protobuf field changes were required.
- No formal version tag was created.

## Implemented

- Added `services/auth-service`.
- Registered generated `auth.v1.AuthService`.
- Implemented `ValidateToken` for RS256 JWT validation.
- Extracted ANI claims:
  - `tid` -> tenant id.
  - `uid` -> user id.
  - `roles` -> role list.
  - `jti` -> blocklist lookup key.
- Added issuer, expiry, not-before, algorithm, signature, tenant id, and user id validation.
- Added Redis JWT blocklist lookup using `jwt:blocklist:{jti}`.
- Implemented minimal `CheckPermission` role rules:
  - `platform-admin` and `tenant-admin`: allowed.
  - `auditor`: read-only.
  - `user`: read/use/create baseline actions only.
- Kept login, refresh, token revocation, and API key management explicitly unimplemented until Dex/OIDC and tenant-safe API key lookup are designed.
- Added focused JWT tests.
- Added auth-service to `go.work`, `Makefile` test packages, and `make build`.

## Files Added

- `services/auth-service/go.mod`
- `services/auth-service/main.go`
- `services/auth-service/internal/config/config.go`
- `services/auth-service/internal/service/auth_service.go`
- `services/auth-service/internal/service/jwt.go`
- `services/auth-service/internal/service/jwt_test.go`

## Files Updated

- `go.work`
- `Makefile`
- `development-records/README.md`

## Validation

Commands run:

```bash
cd repo
make gen-proto
make validate-infra
make test
make build
```

Result:
- `make gen-proto`: passed.
- `make validate-infra`: passed.
- `make test`: passed.
- `make build`: passed.

## Remaining Risks

- `Login`, `RefreshToken`, and `RevokeToken` still require Dex/OIDC integration and persistent token blocklist design.
- API key validation and management remain unimplemented because safe global lookup by key hash needs an explicit RLS-compatible design.
- Gateway middleware still uses local JWT/API key stubs; wiring gateway to auth-service should be a separate `M2.2-AUTH-B` slice.
- Current `CheckPermission` is a minimal role rule engine, not final OPA-backed policy.

## Next Boundary

Recommended next M2 slice:
- `M2.2-AUTH-B`: wire ANI Gateway Auth/RBAC middleware to auth-service over gRPC and preserve local dev mode.
