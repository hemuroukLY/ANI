# Next Development Plan

Date: 2026-05-11

Planning scope:
- Continue module 1 and module 2 in parallel.
- Do not start implementation from this document alone; each slice should start with a clean branch or an explicit continuation decision.

Current completed baseline:
- `M1-INFRA-A`: infrastructure baseline, manifests, Helm values contract, offline YAML validation.
- `M2.1-TASK-A/B/C`: task query path, transactional outbox, tenant-safe worker mutations.
- `M2.2-AUTH-A`: auth-service foundation with RS256 JWT validation and minimal role checks.
- `M1-INFRA-B/C`: component install profiles and tenant network isolation templates.
- `M2.2-AUTH-B/C`: gateway auth-service wiring and RLS-safe API key lifecycle.
- `ARCH-ADAPTER-A/B`: open-source component loose-coupling design plus `pkg/ports` and `pkg/adapters` skeleton.

Required validation before starting the next implementation slice:

```bash
cd repo
make gen-proto
make validate-infra
make test
make build
```

## Parallel Track A: M1-INFRA-B

Product plan mapping:
- `ANI-06 / 模块 1：基础设施底座`
- Focus: component install profiles for storage, metadata, queue, cache, vector DB, and registry dependencies.

Goal:
- Turn the `M1-INFRA-A` baseline into a concrete install profile contract for the first infrastructure dependencies ANI needs to run locally and in an attached K8s cluster.

In scope:
- Extend `deploy/helm/ani-platform/values.yaml` with explicit profiles:
  - `dev`
  - `attach-k8s`
  - `offline`
- Add chart dependency placeholders or subchart contract files for:
  - PostgreSQL / CloudNativePG
  - NATS JetStream
  - Redis
  - MinIO
  - Milvus
  - Harbor
- Add manifest or values documentation for:
  - required namespaces
  - required secrets
  - persistent storage expectations
  - MinIO bucket bootstrap contract
  - PostgreSQL RLS migration application contract
- Keep NetworkPolicy deny-by-default behavior.
- Add validation that remains useful without a live cluster.

Out of scope:
- RKE2 node installation.
- KubeOVN VPC/Subnet CR implementation.
- GPU Operator, HAMi, Volcano, DCGM, Prometheus, Grafana.
- Actual Helm dependency downloads from public chart repositories unless explicitly approved.
- Installer TUI implementation.

Success criteria:
- `make validate-infra` validates all new YAML/values files offline.
- `make gen-proto`, `make test`, and `make build` still pass.
- A new record exists: `repo/development-records/m1-infra-b-component-profiles.md`.
- Version impact is classified according to `ANI-12`; expected impact is `MINOR`.

Primary risks:
- Public Helm chart versions can change, so pinned versions must be documented before any real dependency fetch.
- Harbor is intentionally independent from ANI, so the chart contract must avoid making ANI hard-dependent on Harbor availability.
- RLS migration execution must not require an app role with `BYPASSRLS`.

## Parallel Track B: M2.2-AUTH-B

Product plan mapping:
- `ANI-06 / 模块 2：ANI Gateway / 2.2 认证授权`
- Focus: wire ANI Gateway Auth/RBAC middleware to auth-service.

Goal:
- Replace gateway local verifier stubs with an internal gRPC auth-service client while preserving `ANI_AUTH_MODE=dev`.

In scope:
- Add gateway auth client configuration:
  - auth-service address
  - request timeout
  - fail-closed default
- Update gateway `Auth()` middleware to call `auth.v1.AuthService.ValidateToken`.
- Update gateway `RBAC()` middleware to call `auth.v1.AuthService.CheckPermission`.
- Preserve public path bypass behavior.
- Preserve local dev mode for routes that need manual testing before Dex/OIDC is installed.
- Add focused unit tests for middleware behavior where practical.
- Avoid storing raw tokens in logs.

Out of scope:
- Dex/OIDC login flow.
- Refresh token issuance.
- API key creation/list/revocation.
- OPA policy engine integration.
- Frontend auth screens.

Success criteria:
- Gateway rejects unauthenticated requests when not in dev mode.
- Gateway accepts a request when auth-service returns a valid `TenantContext`.
- Gateway blocks requests when RBAC returns `allowed=false`.
- `make gen-proto`, `make validate-infra`, `make test`, and `make build` pass.
- A new record exists: `repo/development-records/m2-2-auth-b-gateway-auth-wiring.md`.
- Version impact is classified according to `ANI-12`; expected impact is `MINOR` if configs/build targets are additive.

Primary risks:
- Gateway currently constructs middleware without dependency injection, so the minimal change may need a small middleware config object.
- Route-to-resource/action mapping is not yet explicit; `M2.2-AUTH-B` should choose a conservative placeholder mapping rather than inventing full policy metadata.
- Without mTLS, gateway-to-auth-service transport trust is incomplete; record this as a deferred security control tied to deployment work.

## Recommended Execution Order

1. Start `ARCH-ADAPTER-C` first if the next goal is enforcing the new loose-coupling principle in existing services.
2. Start `M1-INFRA-D` first if the next goal is deeper cluster-side infrastructure validation.
3. Start `M2.2-AUTH-D` first if the next goal is completing login/refresh/OIDC-facing auth behavior.
4. Keep each implementation slice independently reviewable and validated.

## Branching

Recommended branches:
- `codex/arch-adapter-c-migrate-components`
- `codex/m1-infra-d-cluster-validation`
- `codex/m2-2-auth-d-login-refresh`

If continuing from the current dirty workspace, first decide whether to commit completed slices together or split them into reviewable commits.
