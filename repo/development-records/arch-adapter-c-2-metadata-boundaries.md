# ARCH-ADAPTER-C-2: Metadata Boundaries

Date: 2026-05-11

## Product Plan Mapping

- Cross-cutting architecture implementation for `ANI-13-开源组件松耦合适配器架构.md`.

## Scope

Reviewed remaining PostgreSQL/pgx direct dependencies after the first ports migration.

Decision:
- PostgreSQL metadata access is a `bounded_direct` integration where direct pgx usage is allowed inside persistence modules because RLS, transaction behavior, row locking, and outbox atomicity are stability-critical.
- Direct pgx usage remains blocked by default outside allowlisted bounded modules.

Changes:
- Updated `repo/architecture/component-import-allowlist.yaml` from `temporary_exception` to `bounded_direct` for pgx metadata modules.
- Removed an unnecessary `pgx.ErrNoRows` dependency from `auth_service.go` by returning `types.ErrNotFound` from API key revocation.

## Validation

Passed:

```bash
make validate-architecture
make test
make build
```

## Version Impact

Expected release impact: `PATCH`.

Reason:
- Internal architecture classification and minor error-boundary cleanup.
