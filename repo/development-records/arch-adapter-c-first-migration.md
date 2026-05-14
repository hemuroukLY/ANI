# ARCH-ADAPTER-C: First Ports Migration

Date: 2026-05-11

## Product Plan Mapping

- Cross-cutting architecture implementation for `ANI-13-开源组件松耦合适配器架构.md`.

## Scope

Migrated two low-risk direct component dependencies behind existing ports:

- `auth-service` JWT blocklist now uses `ports.CacheStore`.
- `task-service` outbox publisher now uses `ports.MessageBus`.
- `auth-service/main.go` wires `deps.Ports.Cache`.
- `task-service/main.go` wires `deps.Ports.MessageBus`.
- Removed corresponding Redis and NATS direct import exceptions from `repo/architecture/component-import-allowlist.yaml`.

## Out Of Scope

- pgx/metadata direct dependencies remain as temporary exceptions.
- No Proto or DB schema changes.
- No MinIO/Milvus/Harbor concrete adapter implementation.

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
- Internal dependency boundary refactor with no external API, Proto, DB, or Helm contract change.

## Next Recommended Slice

- `ARCH-ADAPTER-C-2`: decide whether pgx repository usage remains bounded direct or migrates behind `MetadataStore`.
