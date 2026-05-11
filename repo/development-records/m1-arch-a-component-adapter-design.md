# ARCH-ADAPTER-A / M1-ARCH-A Development Record

Date: 2026-05-11

Product plan mapping:
- Cross-cutting architecture governance for `ANI-06 / 模块 1：基础设施底座` and all upper modules.
- Slice: `ARCH-ADAPTER-A / M1-ARCH-A`
- Scope: open-source component loose-coupling adapter architecture design.

Version impact:
- Current line: `v0.x` development.
- Release impact: `no-release-impact`.
- Reason: documentation and architecture governance only; no runtime code, API, Proto, DB, CRD, or Helm behavior changed in this slice.
- No formal version tag was created.

## Implemented

- Added a dedicated architecture document:
  - `ANI-13-开源组件松耦合适配器架构.md`
- Defined mandatory loose-coupling principle:
  - Kubernetes API is the only stable bottom-layer platform API.
  - MinIO, Milvus, NATS JetStream, Redis, Harbor, CloudNativePG, Dex, and similar components are default adapters only.
  - Business services must depend on ANI-defined ports/capabilities instead of component SDKs.
- Defined initial capability ports:
  - `MetadataStore`
  - `ObjectStore`
  - `VectorStore`
  - `MessageBus`
  - `CacheStore`
  - `ImageRegistry`
  - `IdentityProvider`
- Defined allowed and forbidden locations for component names.
- Defined review checklist and merge blockers for future PRs.
- Added migration plan:
  - `ARCH-ADAPTER-B`: introduce `pkg/ports` and `pkg/adapters` skeleton.
  - `ARCH-ADAPTER-C`: migrate existing direct dependencies behind adapters.

## Files Added

- `../ANI-13-开源组件松耦合适配器架构.md`

## Files Updated

- `../ANI-00-产品战略与开发哲学.md`
- `../ANI-04-技术栈设计.md`
- `../ANI-05-系统架构设计.md`
- `../ANI-11-代码实现规范.md`
- `../CLAUDE.md`
- `development-records/README.md`

## Validation

This is a documentation-only architecture governance slice.

Validation performed:
- Cross-checked existing architecture and implementation references for MinIO, Milvus, NATS, Redis, Harbor, PostgreSQL, and CloudNativePG coupling points.
- No code validation required for this slice.

## Remaining Risks

- Existing code still has direct dependencies on NATS, Redis, and PostgreSQL through `bootstrap` and service implementations.
- Existing proto comments and design docs still mention default components in some flow examples; future cleanup should convert those to capability-first language where they define business semantics.
- Helm values currently expose component keys; `ARCH-ADAPTER-B` should introduce capability-first values while preserving compatibility.

## Next Boundary

Recommended next slice:
- `ARCH-ADAPTER-B`: add `pkg/ports` and `pkg/adapters` skeleton without changing runtime behavior.

Do not continue feature expansion that directly imports or exposes third-party component SDKs before this adapter skeleton exists.
