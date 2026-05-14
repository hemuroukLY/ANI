# ANI Development Records

This directory records implementation progress that maps AI-generated code stages back to the product development plan.

## Current Position

As of 2026-05-11, development is in:

- Product plan file: `ANI-06-开发计划.md`
- Product plan sections:
  - `模块 1：基础设施底座`
  - `模块 2：ANI Gateway（统一 Web Server 层）（M1-M2）`
- Product plan items:
  - `M1-INFRA-A` infrastructure baseline completed and validated.
  - `2.1 Gateway 骨架 / NATS JetStream 异步任务框架` completed and validated.
  - `M2.2-AUTH-A/B/C` auth-service and gateway auth path completed and validated.
  - `ARCH-ADAPTER-A/B` component loose-coupling architecture and code skeleton completed.
  - `ARCH-ADAPTER-GUARD-A` component import guardrail completed.
  - `ARCH-ADAPTER-C`, `M2.2-AUTH-D`, and `M1-INFRA-D` first slices completed and validated.
  - `M1-INFRA-E` GPU scheduling baseline completed and validated.
  - `M2.2-AUTH-K` auth integration profile completed and validated.
  - `M1-INFRA-F` GPU scheduling preflight/e2e hardening completed and validated.
  - `M1-GPU-A` heterogeneous GPU discovery and scheduling contract completed and validated.
  - `M1-RUNTIME-A` workload runtime / instance abstraction completed and validated.
  - `M1-INSTANCE-A` instance object, lifecycle, network plane, and storage attachment contracts completed and validated.
  - `M1-INSTANCE-B` planning runtime adapter completed and validated.
  - `M1-INSTANCE-C` provider dry-run renderer completed and validated.
  - `M1-INSTANCE-D` local admission guardrail completed and validated.
  - `M1-INSTANCE-E` plan/render/admission audit persistence completed and validated.
  - `M1-INSTANCE-F` provider dry-run executor boundary completed and validated.
  - `M1-INSTANCE-G` provider apply/create execution gate completed and validated.
  - `M1-INSTANCE-H` status write-back and lifecycle reconcile contract completed and validated.
  - `M1-INSTANCE-I` provider status reader and instance orchestration API completed and validated.
  - `M1-INSTANCE-J` instance persistence and query API completed and validated.
  - `M1-INSTANCE-K` Kubernetes/KubeVirt provider adapter boundary completed and validated.
  - `M1-INSTANCE-L` instance service API layer completed and validated.
  - `M1-INSTANCE-M` lifecycle and visual ops API completed and validated.
  - `M1-E2E-A` M1 end-to-end integration profile completed and validated.
  - `M1-INSTANCE-N` Kubernetes provider execution profile completed and validated.
  - `M1-INSTANCE-O` adapter-owned Kubernetes REST client completed and validated.
  - `M1-INSTANCE-P` Kubernetes provider bootstrap wiring completed and validated.
  - `M1-INSTANCE-Q` Kubernetes lifecycle execution completed and validated.
  - `M1-INSTANCE-R` Kubernetes visual ops execution completed and validated.
  - `M1-E2E-B` M1 real provider integration regression profile completed and validated.
  - `DEMO-INSTANCE-CONSOLE-A` staged VM/container/GPU container demo console completed and validated.
  - `M1-INSTANCE-S` VM console and remote ops session boundary completed and validated.
  - `DEMO-INSTANCE-WORKSPACE-UI-A` production-oriented instance workspace UI completed and validated.
  - `2026-05-12-demo-handoff` records the paused Demo environment, mock boundaries, startup commands, and next development options.
  - `2026-05-12-aws-instance-lifecycle-reference` records the AWS-derived instance lifecycle feature-depth guide before the next production implementation slice.
  - `2026-05-12-instance-lifecycle-implementation-plan` decomposes the AWS-derived instance lifecycle guide into ANI implementation slices.
  - `2026-05-12-p0-instance-scope-confirmation` confirms the v1.0.0 P0 instance scope: VM, container, GPU container, and basic inference instance.

Important naming note:
- This is not `模块 3：模型管理平台`.
- Earlier records used `Stage 3A/3B/3C` as internal code-generation slice names.
- To avoid confusion with the product plan modules, use `M2.1-TASK-A/B/C` going forward.

## Stage Mapping

| Implementation Slice | Legacy Name | Product Plan Mapping | Status | Record |
|---|---|---|---|---|
| M1-INFRA-A | - | `ANI-06` / `模块 1：基础设施底座` / infrastructure-as-code baseline | Completed | `m1-infra-a-baseline.md` |
| M2.1-TASK-A | Stage 3A | `ANI-06` / `2.1 Gateway 骨架` / `NATS JetStream 异步任务框架` / task query path | Completed | `m2-1-task-a-b-task-service-outbox.md` |
| M2.1-TASK-B | Stage 3B | `ANI-06` / `2.1 Gateway 骨架` / `NATS JetStream 异步任务框架` / transactional outbox + NATS publisher | Completed | `m2-1-task-a-b-task-service-outbox.md` |
| M2.1-TASK-C | Stage 3C | `ANI-06` / `2.1 Gateway 骨架` / `NATS JetStream 异步任务框架` / worker mutation RPCs with tenant-safe writes | Completed | `m2-1-task-c-worker-mutations.md` |
| M2.2-AUTH-A | - | `ANI-06` / `2.2 认证授权` / auth service foundation | Completed | `m2-2-auth-a-auth-service-foundation.md` |
| M1-INFRA-B | - | `ANI-06` / `模块 1：基础设施底座` / component install profiles | Completed | `m1-infra-b-component-profiles.md` |
| M2.2-AUTH-B | - | `ANI-06` / `2.2 认证授权` / gateway auth-service wiring | Completed | `m2-2-auth-b-gateway-auth-wiring.md` |
| M1-INFRA-C | - | `ANI-06` / `模块 1：基础设施底座` / KubeOVN tenant network isolation templates | Completed | `m1-infra-c-network-isolation.md` |
| M2.2-AUTH-C | - | `ANI-06` / `2.2 认证授权` / RLS-safe API Key lifecycle and validation | Completed | `m2-2-auth-c-api-keys.md` |
| ARCH-ADAPTER-A / M1-ARCH-A | - | Cross-cutting architecture / open-source component loose-coupling adapter design | Completed | `m1-arch-a-component-adapter-design.md` |
| ARCH-ADAPTER-B | - | Cross-cutting architecture / `pkg/ports` and `pkg/adapters` capability skeleton | Completed | `arch-adapter-b-ports-adapters-skeleton.md` |
| ARCH-ADAPTER-GUARD-A | - | Cross-cutting architecture / component SDK import guardrail, coupling levels, and allowlist | Completed | `arch-adapter-guard-a-component-imports.md` |
| ARCH-ADAPTER-C | - | Cross-cutting architecture / first direct SDK dependency migration to ports | Completed | `arch-adapter-c-first-migration.md` |
| M2.2-AUTH-D | - | `ANI-06` / `2.2 认证授权` / JWT revocation via CacheStore blocklist | Completed | `m2-2-auth-d-token-revocation.md` |
| M1-INFRA-D | - | `ANI-06` / `模块 1：基础设施底座` / cluster preflight validation profile | Completed | `m1-infra-d-cluster-preflight.md` |
| ARCH-ADAPTER-C-2 | - | Cross-cutting architecture / PostgreSQL metadata bounded direct classification | Completed | `arch-adapter-c-2-metadata-boundaries.md` |
| M2.2-AUTH-E | - | `ANI-06` / `2.2 认证授权` / durable JWT blocklist with CacheStore fast path | Completed | `m2-2-auth-e-durable-token-blocklist.md` |
| M2.2-AUTH-F | - | `ANI-06` / `2.2 认证授权` / refresh token foundation and RS256 access token issuance | Completed | `m2-2-auth-f-refresh-token-foundation.md` |
| M2.2-AUTH-G | - | `ANI-06` / `2.2 认证授权` / OIDC begin/callback boundary and state handling | Completed | `m2-2-auth-g-oidc-boundary.md` |
| M2.2-AUTH-H | - | `ANI-06` / `2.2 认证授权` / OIDC code exchange, ID token verification, and session issuance | Completed | `m2-2-auth-h-oidc-code-exchange.md` |
| M2.2-AUTH-I | - | `ANI-06` / `2.2 认证授权` / OIDC JWKS discovery and key selection | Completed | `m2-2-auth-i-oidc-jwks.md` |
| M2.2-AUTH-J | - | `ANI-06` / `2.2 认证授权` / explicit OIDC group-to-role mapping | Completed | `m2-2-auth-j-oidc-group-mapping.md` |
| M1-INFRA-E | - | `ANI-06` / `模块 1：基础设施底座` / GPU scheduling baseline | Completed | `m1-infra-e-gpu-scheduling-baseline.md` |
| M2.2-AUTH-K | - | `ANI-06` / `2.2 认证授权` / OIDC-refresh-ValidateToken integration profile | Completed | `m2-2-auth-k-auth-integration-profile.md` |
| M1-INFRA-F | - | `ANI-06` / `模块 1：基础设施底座` / GPU scheduling preflight/e2e hardening | Completed | `m1-infra-f-gpu-preflight-e2e.md` |
| M1-GPU-A | - | `ANI-06` / `模块 1：基础设施底座` / heterogeneous GPU discovery and scheduling contract | Completed | `m1-gpu-a-heterogeneous-gpu-contract.md` |
| M1-RUNTIME-A | - | `ANI-06` / `模块 1：基础设施底座` / workload runtime and instance abstraction | Completed | `m1-runtime-a-workload-runtime.md` |
| M1-INSTANCE-A | - | `ANI-06` / `模块 1：基础设施底座` / instance object, lifecycle, network plane, and storage attachment contract | Completed | `m1-instance-a-instance-fabric.md` |
| M1-INSTANCE-B | - | `ANI-06` / `模块 1：基础设施底座` / planning runtime adapter and lifecycle validation | Completed | `m1-instance-b-planning-runtime.md` |
| M1-INSTANCE-C | - | `ANI-06` / `模块 1：基础设施底座` / Kubernetes/KubeVirt provider dry-run renderer | Completed | `m1-instance-c-provider-renderer.md` |
| M1-INSTANCE-D | - | `ANI-06` / `模块 1：基础设施底座` / local admission guardrail for provider manifests | Completed | `m1-instance-d-admission-guardrail.md` |
| M1-INSTANCE-E | - | `ANI-06` / `模块 1：基础设施底座` / instance plan/render/admission audit persistence | Completed | `m1-instance-e-plan-audit.md` |
| M1-INSTANCE-F | - | `ANI-06` / `模块 1：基础设施底座` / provider dry-run executor boundary | Completed | `m1-instance-f-provider-dry-run.md` |
| M1-INSTANCE-G | - | `ANI-06` / `模块 1：基础设施底座` / provider apply/create execution gate | Completed | `m1-instance-g-provider-apply-gate.md` |
| M1-INSTANCE-H | - | `ANI-06` / `模块 1：基础设施底座` / status write-back and lifecycle reconcile contract | Completed | `m1-instance-h-status-reconcile.md` |
| M1-INSTANCE-I | - | `ANI-06` / `模块 1：基础设施底座` / provider status reader and instance orchestration API | Completed | `m1-instance-i-orchestrator.md` |
| M1-INSTANCE-J | - | `ANI-06` / `模块 1：基础设施底座` / instance persistence and query API | Completed | `m1-instance-j-instance-store.md` |
| M1-INSTANCE-K | - | `ANI-06` / `模块 1：基础设施底座` / Kubernetes/KubeVirt provider adapter boundary | Completed | `m1-instance-k-provider-adapter.md` |
| M1-INSTANCE-L | - | `ANI-06` / `模块 1：基础设施底座` / instance service API layer | Completed | `m1-instance-l-instance-service.md` |
| M1-INSTANCE-M | - | `ANI-06` / `模块 1：基础设施底座` / lifecycle and visual ops API | Completed | `m1-instance-m-lifecycle-ops.md` |
| M1-E2E-A | - | `ANI-06` / `模块 1：基础设施底座` / M1 end-to-end integration profile | Completed | `m1-e2e-a-instance-profile.md` |
| M1-INSTANCE-N | - | `ANI-06` / `模块 1：基础设施底座` / Kubernetes provider execution profile | Completed | `m1-instance-n-kubernetes-provider-execution.md` |
| M1-INSTANCE-O | - | `ANI-06` / `模块 1：基础设施底座` / adapter-owned Kubernetes REST client | Completed | `m1-instance-o-kubernetes-rest-client.md` |
| M1-INSTANCE-P | - | `ANI-06` / `模块 1：基础设施底座` / Kubernetes provider bootstrap wiring | Completed | `m1-instance-p-kubernetes-bootstrap-wiring.md` |
| M1-INSTANCE-Q | - | `ANI-06` / `模块 1：基础设施底座` / Kubernetes lifecycle execution | Completed | `m1-instance-q-kubernetes-lifecycle-execution.md` |
| M1-INSTANCE-R | - | `ANI-06` / `模块 1：基础设施底座` / Kubernetes visual ops execution | Completed | `m1-instance-r-kubernetes-ops-execution.md` |
| M1-E2E-B | - | `ANI-06` / `模块 1：基础设施底座` / M1 real provider integration regression profile | Completed | `m1-e2e-b-real-provider-profile.md` |

## Read Order For Continuation

Before continuing code generation, read these files in order:

1. `ANI-06-开发计划.md`
2. `ANI-08-安全架构设计.md`
3. `ANI-09-数据模型设计.md`
4. `ANI-11-代码实现规范.md`
5. `ANI-12-版本管理策略.md`
6. `repo/development-records/m2-1-task-a-b-task-service-outbox.md`
7. `repo/development-records/2026-05-11-handoff-codex-cloud.md`
8. `repo/development-records/2026-05-11-next-development-plan.md`
9. `ANI-13-开源组件松耦合适配器架构.md`
10. `repo/development-records/2026-05-12-demo-handoff.md`
11. `repo/development-records/2026-05-12-aws-instance-lifecycle-reference.md`
12. `repo/development-records/2026-05-12-instance-lifecycle-implementation-plan.md`
13. `repo/development-records/2026-05-12-p0-instance-scope-confirmation.md`

## Version Position

- Current release line: `v0.x` development.
- First formal release target: `v1.0.0` on 2026-09-30.
- Versioning policy: `ANI-12-版本管理策略.md`.
- `M2.1-TASK-C` changed `task_service.proto` and the initial schema; its stage record classifies the release impact as `MINOR`.
- `M2.1-TASK-C` validation passed: `make gen-proto`, `make test`, and `make build`.
- `M1-INFRA-A` adds infrastructure manifests and Helm chart contract; release impact is `MINOR`.
- `M2.2-AUTH-A` adds auth-service; release impact is `MINOR`.
- `M1-INFRA-B` adds install profiles and component contracts; release impact is `MINOR`.
- `M2.2-AUTH-B` wires gateway auth/RBAC to auth-service; release impact is `MINOR`.
- `M1-INFRA-C` adds KubeOVN/network isolation templates; release impact is `MINOR`.
- `M2.2-AUTH-C` implements API Key lifecycle and validation; release impact is `MINOR`.
- `ARCH-ADAPTER-A / M1-ARCH-A` adds architecture governance docs only; release impact is `no-release-impact`.
- `ARCH-ADAPTER-B` adds shared Go ports/adapters skeleton and bootstrap capability wiring; release impact is `MINOR`.
- `ARCH-ADAPTER-GUARD-A` adds validation tooling and Makefile wiring; release impact is `PATCH`.
- `ARCH-ADAPTER-C` migrates internal direct SDK usage to existing ports; release impact is `PATCH`.
- `M2.2-AUTH-D` implements JWT revocation blocklist behavior; release impact is `MINOR`.
- `M1-INFRA-D` adds cluster preflight manifests/profile; release impact is `MINOR`.
- `ARCH-ADAPTER-C-2` classifies metadata pgx usage as bounded direct and removes an unnecessary pgx error dependency from auth-service; release impact is `PATCH`.
- `M2.2-AUTH-E` adds durable JWT blocklist fallback using the existing schema; release impact is `MINOR`.
- `M2.2-AUTH-F` adds refresh token table and AccessToken issuance; release impact is `MINOR`.
- `M2.2-AUTH-G` adds OIDC login RPCs and generated protobuf code; release impact is `MINOR`.
- `M2.2-AUTH-H` adds OIDC callback runtime behavior and config; release impact is `MINOR`.
- `M2.2-AUTH-I` adds OIDC JWKS URL configuration and verifier behavior; release impact is `MINOR`.
- `M2.2-AUTH-J` hardens OIDC role mapping behavior; release impact is `MINOR`.
- `M1-INFRA-E` adds GPU scheduling manifests, Helm profile, and preflight contracts; release impact is `MINOR`.
- `M2.2-AUTH-K` adds integration test coverage only; release impact is `PATCH`.
- `M1-INFRA-F` adds executable GPU preflight/e2e manifests and validation tooling; release impact is `MINOR`.
- `M1-GPU-A` adds GPUInventory port, heterogeneous GPU contracts, and validation tooling; release impact is `MINOR`.
- `M1-RUNTIME-A` adds WorkloadRuntime port, runtime instance contracts, and validation tooling; release impact is `MINOR`.
- `M1-INSTANCE-A` extends WorkloadRuntime instance model and adds Instance Fabric manifests/profiles; release impact is `MINOR`.
- `M1-INSTANCE-B` adds PlanningRuntime adapter, lifecycle validation tests, and planner contract validation; release impact is `MINOR`.
- `M1-INSTANCE-C` adds WorkloadRenderer port, Kubernetes/KubeVirt dry-run renderer, tests, and renderer contract validation; release impact is `MINOR`.
- `M1-INSTANCE-D` adds WorkloadAdmission port, local admission guardrail, tests, and admission contract validation; release impact is `MINOR`.
- `M1-INSTANCE-E` adds WorkloadPlanAuditStore, instance_plan_audits schema/RLS, tests, and audit contract validation; release impact is `MINOR`.
- `M1-INSTANCE-F` adds WorkloadProviderDryRun port, local dry-run executor, tests, and provider dry-run contract validation; release impact is `MINOR`.
- `M1-INSTANCE-G` adds WorkloadProviderApply port, disabled-by-default local apply gate, tests, and provider apply gate contract validation; release impact is `MINOR`.
- `M1-INSTANCE-H` adds WorkloadStatusReconciler port, local phase mapping, tests, and status reconcile contract validation; release impact is `MINOR`.
- `M1-INSTANCE-I` adds WorkloadProviderStatusReader, WorkloadInstanceOrchestrator, local orchestration tests, and orchestration contract validation; release impact is `MINOR`.
- `M1-INSTANCE-J` adds WorkloadInstanceStore, workload_instances schema/RLS, persistence tests, and instance store contract validation; release impact is `MINOR`.
- `M1-INSTANCE-K` adds KubernetesProviderAdapter, KubernetesProviderClient boundary, tests, and provider adapter validation; release impact is `MINOR`.
- `M1-INSTANCE-L` adds WorkloadInstanceService, local service tests, and instance service validation; release impact is `MINOR`.
- `M1-INSTANCE-M` adds lifecycle service methods, WorkloadInstanceOps, local ops guard, tests, and lifecycle/ops validation; release impact is `MINOR`.
- `M1-E2E-A` adds offline end-to-end profile coverage for VM, container, and GPU container create/lifecycle/query/ops; release impact is `PATCH`.
- `M1-INSTANCE-N` adds Kubernetes provider execution profile tests, manifests, and validation tooling without a concrete provider SDK dependency; release impact is `PATCH`.
- `M1-INSTANCE-O` adds adapter-owned KubernetesRESTClient, tests, manifests, and validation tooling without public API/schema changes; release impact is `MINOR`.
- `M1-INSTANCE-P` wires KubernetesRESTClient into bootstrap/config with default-local and apply-disabled guardrails; release impact is `MINOR`.
- `M1-INSTANCE-Q` adds WorkloadInstanceLifecycleExecutor, KubernetesLifecycleExecutor, tests, manifests, and wiring; release impact is `MINOR`.
- `M1-INSTANCE-R` adds KubernetesInstanceOps, tests, manifests, and ops provider wiring; release impact is `MINOR`.
- `M1-E2E-B` adds real-provider integration regression test coverage, manifests, and validation tooling; release impact is `PATCH`.
- `DEMO-INSTANCE-CONSOLE-A` adds a staged presentation API and Console page for VM, container, and GPU container instance flows; release impact is `MINOR`.
- `M1-INSTANCE-S` adds VM console/VNC/serial ops session boundaries and richer staged lifecycle/ops presentation controls; release impact is `MINOR`.
- `DEMO-INSTANCE-WORKSPACE-UI-A` redesigns the staged instance demo as a production-oriented instance workspace reference; release impact is `MINOR`.
- `2026-05-12-demo-handoff` is documentation only and records the paused local Demo environment; release impact is `no-release-impact`.
- `2026-05-12-aws-instance-lifecycle-reference` is documentation only and records AWS-derived instance lifecycle feature-depth requirements; release impact is `no-release-impact`.
- `2026-05-12-instance-lifecycle-implementation-plan` is documentation only and decomposes instance lifecycle feature-depth work into future slices; release impact is `no-release-impact`.
- `2026-05-12-p0-instance-scope-confirmation` is documentation only and records the confirmed v1.0.0 P0 instance feature boundary; release impact is `no-release-impact`.

## Next Work Boundary

Recommended next implementation options:

- Optional `M1-DEMO-SMOKE-A`: live KubeVirt/Kubernetes/GPU smoke profile for presentation environments.
- `M1-INSTANCE-T`: AWS-derived instance product spec enrichment for state reasons, action precheck, operation timeline, and production-grade click feedback.
- P0 instance scope reference: `2026-05-12-p0-instance-scope-confirmation.md`.
- Instance implementation planning reference: `2026-05-12-instance-lifecycle-implementation-plan.md`.
- `M3-MODEL-A`: module 3 model metadata and object-storage boundary design/implementation, after accepting the M1/M2/GPU/Runtime/Instance baseline.
- Demo resume reference: `2026-05-12-demo-handoff.md`.
- Instance lifecycle feature-depth reference: `2026-05-12-aws-instance-lifecycle-reference.md`.

Continue module 1 and module 2 in parallel, but keep each implementation slice independently reviewable and validated.

Planning status:
- Detailed plan is ready in `2026-05-11-next-development-plan.md`.
- `M1-INFRA-B`, `M2.2-AUTH-B`, `M1-INFRA-C`, `M2.2-AUTH-C`, `ARCH-ADAPTER-A / M1-ARCH-A`, `ARCH-ADAPTER-B`, `ARCH-ADAPTER-GUARD-A`, `ARCH-ADAPTER-C`, `M2.2-AUTH-D`, `M1-INFRA-D`, `ARCH-ADAPTER-C-2`, `M2.2-AUTH-E`, `M2.2-AUTH-F`, `M2.2-AUTH-G`, `M2.2-AUTH-H`, `M2.2-AUTH-I`, `M2.2-AUTH-J`, `M1-INFRA-E`, `M2.2-AUTH-K`, `M1-INFRA-F`, `M1-GPU-A`, `M1-RUNTIME-A`, `M1-INSTANCE-A`, `M1-INSTANCE-B`, `M1-INSTANCE-C`, `M1-INSTANCE-D`, `M1-INSTANCE-E`, `M1-INSTANCE-F`, `M1-INSTANCE-G`, `M1-INSTANCE-H`, `M1-INSTANCE-I`, `M1-INSTANCE-J`, `M1-INSTANCE-K`, `M1-INSTANCE-L`, `M1-INSTANCE-M`, `M1-E2E-A`, `M1-INSTANCE-N`, `M1-INSTANCE-O`, `M1-INSTANCE-P`, `M1-INSTANCE-Q`, `M1-INSTANCE-R`, `M1-E2E-B`, `DEMO-INSTANCE-CONSOLE-A`, `M1-INSTANCE-S`, and `DEMO-INSTANCE-WORKSPACE-UI-A` have been implemented.
- Before starting implementation, rerun `make gen-proto`, `make validate-infra`, `make test`, and `make build`.
- Before a Demo presentation, rerun `make validate-demo-instances`, `git diff --check`, and `npm run build` in `frontends/console`.
