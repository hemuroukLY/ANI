# ARCH-ADAPTER-GUARD-A: Component Import Guardrail

Date: 2026-05-11

## Product Plan Mapping

- Cross-cutting architecture guardrail for `ANI-13-开源组件松耦合适配器架构.md`.
- Protects future work across `ANI-06 / 模块 1` and service modules from bypassing `pkg/ports`.

## Scope

Implemented a lightweight engineering guardrail:

- Added `repo/scripts/validate_component_imports.py`.
- Added `repo/architecture/component-import-allowlist.yaml`.
- Added `make validate-architecture`.
- Wired `make test` to run `validate-architecture` before Go/Python tests.
- Documented the guardrail in `ANI-11`, `ANI-13`, `CLAUDE.md`, and `AGENTS.md`.
- Added coupling levels so direct binding is a controlled architecture decision rather than a blanket ban:
  - `port_required`
  - `adapter_with_extensions`
  - `bounded_direct`
  - `temporary_exception`

The checker blocks direct component SDK imports outside approved areas:

- `pkg/adapters/`
- `pkg/bootstrap/`

Existing direct imports in service/repository code are explicitly listed as `temporary_exception` migration debt and must be removed in `ARCH-ADAPTER-C`.

## Out Of Scope

- No service code migration was performed in this slice.
- No runtime behavior changed.
- No new third-party SDK dependency was added.

## Validation

Passed:

```bash
make validate-architecture
```

## Version Impact

Expected release impact: `PATCH`.

Reason:
- Adds validation tooling and build workflow guardrails.
- Does not change runtime API, Proto, DB, Helm, or service behavior.

## Next Recommended Slice

- `ARCH-ADAPTER-C`: start removing allowlist entries by migrating direct Redis/NATS/pgx usage behind `Deps.Ports`.
