# ARCH-ADAPTER-B: Ports And Adapters Skeleton

Date: 2026-05-11

## Product Plan Mapping

- Cross-cutting architecture implementation following `ANI-13-开源组件松耦合适配器架构.md`.
- Supports `ANI-06 / 模块 1：基础设施底座` and all service modules that integrate open-source components.

## Scope

Implemented the first shared Go boundary for open-source component loose coupling:

- Added `repo/pkg/ports` capability interfaces:
  - `MetadataStore`
  - `ObjectStore`
  - `VectorStore`
  - `MessageBus`
  - `CacheStore`
  - `ImageRegistry`
  - `IdentityProvider`
- Added default adapter packages under `repo/pkg/adapters`:
  - `postgres` metadata adapter wrapping `pgxpool`.
  - `nats` message bus adapter wrapping JetStream.
  - `redis` cache adapter wrapping go-redis.
  - `objectstore`, `vectorstore`, `registry`, and `identity` not-configured adapters for future concrete providers.
- Added `bootstrap.Capabilities` and `Deps.Ports` so new services can use ANI capability ports while existing raw clients remain available during migration.
- Synchronized `AGENTS.md` with `CLAUDE.md` for the new mandatory `ANI-13` read order and adapter principle.

## Out Of Scope

- No service-layer migration from raw `pgx`, `nats`, or `redis` clients yet.
- No MinIO, Milvus, Harbor, Dex, or Keycloak SDK dependency added.
- No OpenAPI, Proto, DB schema, Helm, or CRD changes.

## Validation

Passed:

```bash
GOCACHE=/Users/zhangfan/ANI/repo/.cache/go-build GOMODCACHE=/Users/zhangfan/ANI/repo/.cache/gomod go test ./pkg/...
make validate-infra
make test
make build
```

## Version Impact

Expected release impact: `MINOR`.

Reason:
- Adds shared public Go interfaces and adapter packages.
- Does not change external OpenAPI/Proto contracts or database schema.

## Next Recommended Slice

- `ARCH-ADAPTER-C`: migrate direct component dependencies behind `pkg/ports` incrementally, starting with low-risk gateway/auth cache usage or task outbox publish path.
- Or resume product work with `M1-INFRA-D` / `M2.2-AUTH-D` while requiring new code to use `Deps.Ports`.
