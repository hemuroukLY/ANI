# ANI Development Records

This directory records implementation progress that maps AI-generated code stages back to the product development plan.

## Current Position

As of 2026-05-11, development is in:

- Product plan file: `ANI-06-开发计划.md`
- Product plan section: `模块 2：ANI Gateway（统一 Web Server 层）（M1-M2）`
- Product plan item: `2.1 Gateway 骨架 / NATS JetStream 异步任务框架`
- Implementation status: `Stage 3A + 3B` completed; `Stage 3C` is next.

## Stage Mapping

| Implementation Stage | Product Plan Mapping | Status | Record |
|---|---|---|---|
| Stage 3A | `ANI-06` / `2.1 Gateway 骨架` / `NATS JetStream 异步任务框架` / task query path | Completed | `stage-3A-3B-task-service-outbox.md` |
| Stage 3B | `ANI-06` / `2.1 Gateway 骨架` / `NATS JetStream 异步任务框架` / transactional outbox + NATS publisher | Completed | `stage-3A-3B-task-service-outbox.md` |
| Stage 3C | `ANI-06` / `2.1 Gateway 骨架` / `NATS JetStream 异步任务框架` / worker mutation RPCs with tenant-safe writes | Pending | To be created: `stage-3C-task-worker-mutations.md` |

## Read Order For Continuation

Before continuing code generation, read these files in order:

1. `ANI-06-开发计划.md`
2. `ANI-08-安全架构设计.md`
3. `ANI-09-数据模型设计.md`
4. `ANI-11-代码实现规范.md`
5. `repo/development-records/stage-3A-3B-task-service-outbox.md`
6. `repo/development-records/2026-05-11-handoff-codex-cloud.md`

## Next Work Boundary

Next implementation should be Stage 3C only:

- Fix `repo/api/proto/task/v1/task_service.proto` so worker mutation RPCs have tenant/security context.
- Regenerate protobuf code.
- Implement tenant-safe task worker mutation RPCs.
- Preserve PostgreSQL RLS and transactional outbox guarantees.
- Run `make gen-proto`, `make test`, and `make build`.
- Add `repo/development-records/stage-3C-task-worker-mutations.md`.
