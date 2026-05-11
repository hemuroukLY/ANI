# M2.2-AUTH-C Development Record

Date: 2026-05-11

Product plan mapping:
- `ANI-06 / 模块 2：ANI Gateway（统一 Web Server 层） / 2.2 认证授权`
- Slice: `M2.2-AUTH-C`
- Scope: RLS-safe API Key lifecycle and validation.

Version impact:
- Current line: `v0.x` development.
- Release impact: `MINOR`.
- Reason: implements previously declared API Key RPC behavior and changes auth-service validation semantics for `X-API-Key`.
- No formal version tag was created.

## Implemented

- Implemented API Key creation.
- Implemented API Key list.
- Implemented API Key revocation.
- Implemented API Key validation through `ValidateToken`.
- Raw API Key format embeds tenant UUID:
  - `ani_dev_{tenant_uuid}_{random_secret}`
- Auth-service parses tenant UUID from raw key, sets PostgreSQL RLS tenant context, then queries `api_keys` by SHA256 hash.
- No `BYPASSRLS` role is introduced.
- API Key scopes are returned as roles alongside `service-account`.
- `CheckPermission` now recognizes scope-style permissions:
  - `scope:{resource}:{action}`
  - `{resource}:{action}`
  - `scope:{resource}:*`
  - `{resource}:*`
  - `scope:*:*`
- Added focused API Key format and scope tests.

## Files Added

- `services/auth-service/internal/service/api_keys.go`
- `services/auth-service/internal/service/api_keys_test.go`

## Files Updated

- `services/auth-service/main.go`
- `services/auth-service/internal/service/auth_service.go`
- `development-records/README.md`

## Validation

Commands run:

```bash
cd repo
GOCACHE=/private/tmp/ani-go-build GOMODCACHE="$PWD/.cache/gomod" go test ./services/auth-service/...
make gen-proto
make validate-infra
make test
make build
```

Result:
- Auth-service package tests: passed.
- `make gen-proto`: passed.
- `make validate-infra`: passed.
- `make test`: passed.
- `make build`: passed.

## Remaining Risks

- API Key creation RPC currently trusts caller-supplied `tenant_id`; caller authorization must be enforced by Gateway/RBAC and later internal auth policy.
- Key format uses `dev` as the environment segment until release packaging defines environment naming.
- API Key integration tests need a PostgreSQL test database to verify RLS behavior end to end.
- Dex/OIDC login and refresh token flows remain deferred.

## Next Boundary

Recommended next M2 slice:
- `M2.2-AUTH-D`: Dex/OIDC login integration or PostgreSQL-backed integration tests for API Key RLS behavior.
