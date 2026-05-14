# M2.2-AUTH-E: Durable JWT Blocklist

Date: 2026-05-11

## Product Plan Mapping

- `ANI-06` / `模块 2：ANI Gateway` / `2.2 认证授权`.

## Scope

Improved JWT revocation stability:

- Added `persistentTokenBlocklist`.
- `RevokeToken` now writes to the existing `jwt_blocklist` PostgreSQL table and `CacheStore`.
- JWT validation checks `CacheStore` first, then PostgreSQL as a durability fallback, then backfills the cache.
- Preserves the fast path while keeping revocation effective after cache restart.

## Validation

Passed:

```bash
make validate-architecture
make test
make build
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Implements stronger behavior for an existing auth RPC without changing Proto or DB schema.

## Remaining Risks

- `Login` and `RefreshToken` remain unimplemented pending Dex/OIDC callback and refresh-token storage design.
