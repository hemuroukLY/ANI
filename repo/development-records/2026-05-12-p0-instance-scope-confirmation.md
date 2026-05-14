# P0 Instance Scope Confirmation

Date: 2026-05-12

## Decision

The v1.0.0 P0 scope for ANI instance lifecycle management is confirmed.

P0 is not intended to match every AWS feature. P0 must make ANI credible as a
production POC platform for AI private cloud scenarios while avoiding an
unbounded first-release scope.

## Confirmed P0 Instance Types

P0 includes:

- VM / cloud host.
- Container instance.
- GPU container instance.
- Inference instance, basic endpoint lifecycle only.

P1/P2, not P0:

- Notebook.
- Batch / training job.
- Agent Sandbox.

## Confirmed P0 Lifecycle Actions

P0 common actions:

- Create.
- List and detail query.
- Start.
- Stop.
- Restart.
- Resize.
- Delete.

P1/P2 actions, not P0:

- Clone.
- Snapshot.
- Backup and restore.
- Migration and recovery.
- Hibernate/suspend unless a provider makes it trivial behind the existing
  lifecycle boundary.

## Confirmed P0 Operations

P0 operations by type:

| Capability | VM | Container | GPU container | Inference |
|---|---|---|---|---|
| Logs | P0 | P0 | P0 | P0 |
| Events | P0 | P0 | P0 | P0 |
| Metrics | P0 | P0 | P0 + GPU metrics | P0 + QPS/latency/errors |
| Terminal / exec | P1 or limited P0 | P0 | P0 | P1 |
| Web console / VNC | P0 | N/A | N/A | N/A |
| Serial console | P1 | N/A | N/A | N/A |
| SSH/RDP management | P1; P0 shows connection metadata only | N/A | N/A | N/A |
| Audit | P0 | P0 | P0 | P0 |

## Confirmed P0 Inference Scope

P0 inference includes:

- Select model version.
- Select GPU specification.
- Set replica count.
- Create endpoint.
- View endpoint status.
- View logs, events, and metrics.
- Start, stop, restart, resize, delete.
- Basic endpoint test invocation.
- Token usage as P0-lite if it does not block the lifecycle work.

Not P0:

- Canary release.
- Multi-version traffic split.
- Automatic rollback.
- Advanced autoscaling.
- A/B testing.

## Confirmed P0 Operation Feedback

Every P0 operation must have:

- Precheck.
- Disabled-action reason.
- Confirmation for high-risk actions.
- Operation ID after backend acceptance.
- Operation timeline.
- Failure reason and recommended fix.
- Retry eligibility when applicable.
- Audit record.

This confirms `M1-INSTANCE-T` as the next implementation starting point before
VM/container/GPU/inference feature-depth work.

## Follow-up Documentation Updates

This decision is now reflected in:

- `ANI-02-产品功能设计.md`.
- `ANI-06-开发计划.md`.
- `repo/development-records/README.md`.

## Release Impact

This record is documentation only.

Expected release impact: `no-release-impact`.

