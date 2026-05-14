# AWS Instance Lifecycle Reference for ANI

Date: 2026-05-12

## Purpose

This record captures the product-depth gap found from the Demo console:
the current ANI instance implementation validates the minimum technical chain,
but it is not yet the full instance lifecycle experience expected from the
product goal.

The goal of this reference is to turn AWS public-cloud instance behavior into
an ANI development guide for:

- VM / cloud host instances.
- Container and GPU container instances.
- Model inference instances.
- Notebook / interactive development instances.
- Batch / training / large-compute jobs.
- Foundation-model invocation capacity and cost tracking.

This guide does not require ANI to copy AWS concepts one-to-one. It uses AWS as
a functional benchmark while preserving ANI's existing architecture rule:
business code must continue to use `WorkloadRuntime`, `WorkloadInstanceService`,
`WorkloadInstanceOps`, `GPUInventory`, `ObjectStore`, `MessageBus`, and other
ANI ports rather than calling provider APIs directly.

## Local Documentation Consistency Check

### Consistent Points

- `ANI-06` defines all VM, container, GPU container, inference, notebook,
  sandbox, and batch job workloads as first-class ANI instances.
- `ANI-13` correctly requires provider details to stay behind ports/adapters.
- M1 implementation records from `M1-RUNTIME-A` through `M1-INSTANCE-S` are
  internally consistent: creation goes through planning, rendering, admission,
  audit, dry-run, gated apply, status observation, reconcile, persistence, and
  service APIs.
- `DEMO-INSTANCE-WORKSPACE-UI-A` explicitly states that non-instance resource
  pages are mock surfaces and that real VM/container/GPU container flows must
  remain behind `WorkloadInstanceService` and `WorkloadInstanceOps`.

### Gap / Risk

The docs are technically consistent, but the wording around "completed" can be
misread as "full product feature completed". In reality, most completed M1
records prove architectural and lifecycle control-plane boundaries, not the
full user-facing depth of a production cloud console.

Specific gaps:

- Lifecycle is currently modeled as start/stop/restart/resize/delete, but AWS
  class products expose richer lifecycle states, safety checks, scheduled
  events, alarms, recovery, hibernation/suspend variants, forced termination
  warnings, and deletion consequences.
- Ops currently covers logs/events/metrics/terminal/exec/VM console, but AWS
  separates normal connection, serial/break-glass console, Session Manager-like
  browser shell, command auditing, session duration, and connection
  prerequisites.
- Resource tabs exist in the Demo, but production behavior for images, disks,
  snapshots, security groups, network interfaces, tags, quotas, cost, alarms,
  jobs, deployments, rollbacks, and model autoscaling is still not specified at
  AWS-level depth.
- The Demo creates a strong information-architecture preview, but mock data
  must not become the feature target. The feature target should be the concrete
  lifecycle and operation matrix below.

## AWS Functional Benchmarks

The following AWS references are used as benchmark material:

- EC2 lifecycle states, billing, reboot/stop/hibernate/terminate differences:
  https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
- EC2 connection methods: SSH, RDP, EC2 Instance Connect, Session Manager,
  Fleet Manager, and Instance Connect Endpoint:
  https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/connect.html
- EC2 status checks and scheduled events:
  https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/monitoring-instances-status-check.html
- EC2 serial console:
  https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/connect-to-serial-console.html
- ECS task lifecycle:
  https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-lifecycle-explanation.html
- ECS service deployments and deployment history:
  https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-deployment.html
- ECS deployment strategies:
  https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs_service-options.html
- ECS Exec:
  https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-exec.html
- Container Insights for ECS/EKS metrics and logs:
  https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/ContainerInsights.html
- EKS/Kubernetes GPU, pod, node, service metrics:
  https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-metrics-EKS.html
- SageMaker real-time inference endpoints:
  https://docs.aws.amazon.com/sagemaker/latest/dg/realtime-endpoints.html
- SageMaker endpoint autoscaling:
  https://docs.aws.amazon.com/sagemaker/latest/dg/endpoint-auto-scaling.html
- SageMaker Studio running applications and spaces:
  https://docs.aws.amazon.com/sagemaker/latest/dg/studio-updated-running-stop.html
- AWS Batch job states:
  https://docs.aws.amazon.com/batch/latest/userguide/job_states.html
- Bedrock Provisioned Throughput:
  https://docs.aws.amazon.com/bedrock/latest/userguide/prov-throughput.html
- Bedrock inference profiles:
  https://docs.aws.amazon.com/bedrock/latest/userguide/inference-profiles.html
- Bedrock model invocation logging:
  https://docs.aws.amazon.com/bedrock/latest/userguide/model-invocation-logging.html

## Target ANI Instance Domain Model

ANI should keep one top-level "Instance" product concept, but internally use
typed sub-kinds:

| ANI instance kind | AWS benchmark | ANI provider target |
|---|---|---|
| `vm` | EC2 instance | KubeVirt / OpenStack / VMware / existing cloud adapter |
| `container` | ECS task/service or EKS workload | Kubernetes Deployment/Pod/Job adapter |
| `gpu_container` | EKS GPU workload / ECS GPU task | Kubernetes + GPUInventory + HAMi/Volcano |
| `inference` | SageMaker endpoint / Bedrock provisioned capacity | WorkloadRuntime + model-service + inference-operator |
| `notebook` | SageMaker Studio app/space | WorkloadRuntime + PVC/object storage + browser session |
| `agent_sandbox` | Managed isolated execution environment | WorkloadRuntime + restricted network/storage policy |
| `batch_job` | AWS Batch job | Kubernetes Job/Volcano/Kueue-style queued workload |

## Lifecycle State Targets

Current ANI states are sufficient as a minimal baseline, but production should
add state reasons, substatus, timestamps, and user-visible explanations.

### Common States

- `draft`: configuration saved, not submitted.
- `validating`: quota, image, network, storage, GPU, policy checks running.
- `pending`: accepted but not scheduled.
- `provisioning`: provider resources being created.
- `starting`: runtime process is starting.
- `running`: usable.
- `updating`: spec, version, image, model, or capacity update in progress.
- `scaling`: replica, GPU, CPU, memory, or token capacity change in progress.
- `stopping`: graceful shutdown in progress.
- `stopped`: stopped and restartable.
- `hibernated`: memory/root-volume preserved where supported.
- `degraded`: running but health/status checks are impaired.
- `recovering`: automatic or manual recovery running.
- `failed`: failed terminal or recoverable state with reason.
- `deleting`: cleanup in progress.
- `deleted`: final deleted state.

### Kind-Specific State Notes

- VM should distinguish `rebooting`, `stopping`, `stopped`, `hibernated`,
  `shutting_down`, and `terminated/deleted` semantics. Deletion must explain
  which disks, IPs, snapshots, and backups are retained or removed.
- Container should expose task/workload states similar to `provisioning`,
  `pending`, `activating`, `running`, `deactivating`, `stopping`,
  `deprovisioning`, and `stopped`.
- Inference should expose endpoint states such as creating, updating,
  in-service/running, rolling back, scaling, failed, and deleting.
- Batch/training should expose submitted, pending, runnable, starting,
  running, succeeded, failed, cancelled, and retrying.
- Notebook should separate application runtime state from persistent workspace
  storage state. Stopping the app should preserve files; deleting the workspace
  should warn that data is removed unless backed up.

## Create Flow Requirements

All instance creation flows should follow the same UX and API shape:

1. Choose instance kind.
2. Choose blueprint/template.
3. Configure compute.
4. Configure image/model/runtime.
5. Configure network.
6. Configure storage.
7. Configure security/access.
8. Configure observability.
9. Configure cost/quota/labels.
10. Review plan.
11. Submit.
12. Show timeline with per-step status.

### VM Create

Required fields:

- Name, description, tenant/project, tags.
- Region/AZ/placement target.
- Image/template/boot source.
- CPU, memory, GPU/NPU if applicable.
- Root disk size, type, retention policy.
- Data disks, attach/detach policy, mount hints.
- Network plane, subnet/VPC, private IP mode, public ingress mode.
- Security group/firewall rules.
- Login method: SSH key, password, no-login, SSM-like managed session.
- Cloud-init/user data.
- Backup/snapshot policy.
- Termination protection.

Submit feedback:

- Disable submit while validation runs.
- Show estimated quota impact and known blocking issues.
- On submit, create a plan ID immediately.
- Timeline must show validation, audit, dry-run, provider apply, status
  observation, and final ready/failed state.

### Container / GPU Container Create

Required fields:

- Image, tag/digest, pull secret, command, args, working directory.
- CPU/memory request and limit.
- GPU vendor/model/count/slice type, runtime class, queue.
- Replica count or job completion policy.
- Environment variables, secrets, config maps.
- Ports, service exposure, ingress, internal-only flag.
- Volumes: PVC, object mount, ephemeral, host-path denied by default.
- Health checks: startup, readiness, liveness.
- Restart policy and rollout strategy.
- Logs and metrics collection options.
- Exec/terminal enabled flag with audit policy.

Submit feedback:

- Image pull secret and registry access should be prechecked.
- GPU scheduling failures must be visible before provider apply when possible.
- If GPU inventory is missing, reject with a scheduling reason rather than
  creating a stuck workload.

### Inference Create

Required fields:

- Model and model version.
- Runtime engine: vLLM/Triton/custom adapter.
- Instance shape: CPU/memory/GPU/NPU, tensor parallel, replicas.
- Endpoint mode: internal, tenant VPC, public ingress, OpenAI-compatible.
- Autoscaling: min/max replicas, target metrics, cooldown.
- Traffic policy: canary, weighted variants, rollback threshold.
- Model loading: object-store path, cache policy, encrypted model key ref.
- Request limits: max tokens, concurrency, timeout, rate limit.
- Observability: QPS, latency, TTFT, token counts, error rate, GPU metrics.
- Invocation logging and data retention.

Submit feedback:

- Return endpoint ID and task ID immediately.
- Show model download, cache, decrypt, warmup, health check, route registration,
  and readiness steps.
- If update is requested, show traffic-shift and rollback progress.

### Notebook / Workspace Create

Required fields:

- Image/runtime template.
- CPU/memory/GPU.
- Persistent space/PVC size.
- Git repo or object-store mount.
- Idle shutdown policy.
- Network access policy.
- Browser access method.
- Collaboration/share policy.

Submit feedback:

- Starting an app creates or attaches runtime compute.
- Stopping the app preserves workspace files.
- Deleting the workspace requires explicit confirmation because persistent data
  can be removed.

### Batch / Training Create

Required fields:

- Job definition/template.
- Image/command/args.
- Input datasets and output location.
- CPU/memory/GPU/NPU.
- Queue, priority, retry strategy, timeout.
- Multi-node or distributed training topology.
- Checkpoint path.
- Metrics parser and log collection.

Submit feedback:

- Timeline must show submitted, pending dependencies, runnable, starting,
  running attempts, retrying, succeeded/failed.
- Stop must be graceful first and then forceful after timeout.
- Failed job detail must include exit code, attempt count, retry decision, and
  log shortcut.

## Detail Page / Tab Requirements

Each instance detail page should have stable tabs. Tabs can be hidden or disabled
when a kind does not support them, but the mental model should remain consistent.

### Overview

- State, state reason, provider state, health.
- Create/update/delete timestamps.
- Owner, tenant/project, tags.
- Provider references and audit IDs.
- Endpoint/IP/URL where applicable.
- Cost/quota impact.
- Recent timeline events.

### Configuration

- Runtime kind, image/model/template.
- CPU, memory, GPU/NPU, replicas.
- Placement: region, AZ, node class, queue.
- Restart, rollout, autoscaling, timeout policies.

### Metrics

- VM: CPU, memory, disk IO, network IO, status checks.
- Container: CPU, memory, restarts, network, pod/task health.
- GPU: GPU allocation, GPU utilization, memory utilization, temperature,
  scheduler queue and reserved capacity.
- Inference: QPS, latency percentiles, TTFT, input/output tokens, errors,
  concurrency, replica saturation, model load time.
- Batch/training: elapsed time, resource utilization, progress, custom metrics.

### Logs

- Runtime logs.
- System/provider events.
- Application stdout/stderr.
- Previous attempt logs.
- Export/copy log links.

### Events

- Lifecycle transitions.
- Provider events.
- Scheduler events.
- Health check changes.
- Scaling events.
- Alarm state changes.

### Network

- Network plane.
- VPC/subnet/security group/firewall policy.
- Private/public address.
- Ports and ingress routes.
- DNS/service discovery.
- Egress restrictions.

### Storage

- Root disk.
- Data disks.
- PVC/object mounts.
- Snapshot/backup policy.
- Retention on stop/delete.
- Capacity and IO metrics.

### Security

- Access method.
- Keys/secrets.
- Service account / IAM-like role.
- Exec/terminal permission.
- Audit policy.
- Privileged and host network guardrail result.

### Operations

- Lifecycle actions.
- Console/terminal/exec.
- Snapshot/backup/restore.
- Resize/update.
- Clone/create image.
- Redeploy/rollback for inference and service-like instances.

### Audit

- Who created/updated/operated/deleted.
- Permission proof.
- Request ID.
- Plan ID / audit ID.
- Provider action ID.
- Before/after spec diffs.

## Action Matrix

### Common Actions

| Action | Applies to | Preconditions | User feedback |
|---|---|---|---|
| Create | all | valid spec, quota, policy | create timeline, plan ID, task ID |
| Start | stopped VM/notebook/container-like service | stopped and restartable | state moves to starting, then running or failed |
| Stop | running restartable instances | running/degraded | confirm retention and billing impact, show graceful timeout |
| Restart | running VM/container/inference/notebook | running/degraded | state moves to restarting, preserve identity where supported |
| Resize | VM/container/GPU/inference/notebook | supported state, quota | show whether restart/recreate is required |
| Delete | all | permission, optional protection disabled | destructive confirmation, retention summary |
| Force delete | stuck/deleting/failed | elevated permission | high-risk confirmation, best-effort cleanup result |
| Tag/edit metadata | all | write permission | inline update, audit event |
| Clone | VM/container/notebook/inference | snapshot/template available | new draft with copied spec |
| Export spec | all | read permission | downloadable manifest/spec without secrets |

### VM-Specific Actions

- Connect by SSH/RDP where network policy allows.
- Browser managed session where agent/role is installed.
- Serial/VNC console for break-glass troubleshooting.
- Get console output / screenshot equivalent where provider supports it.
- Create snapshot.
- Create image/template from instance.
- Attach/detach data disk.
- Attach/detach network interface.
- Change security group/firewall rule.
- Enable/disable termination protection.
- Hibernate/suspend where provider supports it.
- Recover/migrate when host/system status is impaired.

### Container / GPU Container Actions

- Scale replicas.
- Rollout new image.
- Rollback deployment.
- Pause/resume rollout.
- View deployment history.
- Exec shell.
- Run one-off command.
- View logs, previous logs, events, metrics.
- Restart one replica/pod.
- Drain/evict from node.
- Change resource requests/limits.
- Change GPU allocation or queue, with recreate warning when required.

### Inference Actions

- Invoke test request.
- Open OpenAI-compatible endpoint details.
- Scale replicas.
- Configure autoscaling.
- Update model version.
- Canary/weighted traffic shift.
- Rollback.
- Warm up model.
- Drain old replicas.
- Enable/disable invocation logging.
- Set rate limits and token limits.
- View per-model and per-variant metrics.
- Track usage and cost by tenant/project/application.

### Notebook Actions

- Open browser session.
- Stop app but preserve workspace.
- Restart app.
- Change instance size.
- Attach Git repo or object mount.
- Manage kernels/terminals.
- Delete workspace with data-loss confirmation.
- Idle shutdown policy.

### Batch / Training Actions

- Submit.
- Cancel before running.
- Stop running job.
- Retry failed job.
- Clone and rerun.
- View attempts.
- View dependency graph.
- View queue position.
- Update priority if allowed.
- Export artifacts.
- Open checkpoint/output location.

## Click Feedback Rules

Every action button must produce deterministic feedback:

- `Precheck`: validate state, permission, quota, provider support, and
  destructive consequences before opening final confirmation.
- `Confirmation`: required for delete, force delete, stop with data-loss risk,
  resize requiring recreate, workspace deletion, and public ingress exposure.
- `Immediate response`: return operation ID, task ID, or audit ID immediately.
- `Optimistic state`: show transitional state only after backend accepts the
  operation.
- `Timeline`: append step-by-step progress with timestamps.
- `Failure`: show reason, provider message, recommended fix, retry eligibility,
  and related logs/events.
- `Audit`: every operation writes user, tenant, request ID, permission proof,
  resource ID, before/after diff, and provider action references.
- `Idempotency`: repeated create/update/operation submissions must not duplicate
  resources when an idempotency key is provided.
- `Long-running operations`: show non-blocking task progress, allow user to
  leave page, and surface final notification.

## ANI Development Slices Suggested From This Guide

### M1-INSTANCE-T: Instance Product Spec Enrichment

Goal:

- Extend the instance domain contract with state reason, operation ID, action
  precheck, destructive-action metadata, provider support flags, and operation
  timeline.

Success criteria:

- No provider SDK leaks into business services.
- Existing M1 tests still pass.
- Instance detail APIs can explain why an action is allowed, disabled, or
  dangerous.

### M1-INSTANCE-U: VM Production Ops Depth

Goal:

- Add VM connection modes, console output/screenshot metadata, serial/VNC
  session lifecycle, termination protection, snapshot/backup placeholders, and
  disk/network attachment contracts.

Success criteria:

- KubeVirt/OpenStack/VMware/cloud adapters can map their own console/session
  mechanisms behind `WorkloadInstanceOps`.

### M1-INSTANCE-V: Container/GPU Deployment Operations

Goal:

- Add rollout, rollback, replica scaling, deployment history, per-replica
  restart, scheduler reason, GPU queue/state, and previous logs.

Success criteria:

- GPU scheduling failures are explainable before or during admission/provider
  dry-run.

### M3-INFERENCE-A: Inference Instance Lifecycle

Goal:

- Model inference becomes a first-class instance kind with model version,
  runtime engine, endpoint, autoscaling, traffic policy, token limits, logs, and
  metrics.

Success criteria:

- Inference create/update/delete uses `WorkloadRuntime` and never directly
  creates Deployment/KubeVirt resources from model-service.

### M4-CONSOLE-INSTANCE-A: Production Instance Console

Goal:

- Replace demo-only surfaces with production APIs, using the AWS-derived tab and
  action matrix above.

Success criteria:

- Every disabled button has a reason.
- Every accepted operation has an operation ID and timeline.
- Every destructive operation has an explicit retention/cost/security summary.

## Product Priority

P0 for production credibility:

- State reason and operation timeline.
- Action precheck and disabled-action reasons.
- Logs/events/metrics linked from every instance.
- VM console/session boundary.
- Container exec with audit.
- Delete/stop/resize confirmations with retention impact.
- GPU scheduling reason and GPU utilization visibility.
- Inference endpoint lifecycle and autoscaling basics.

P1 for public-cloud parity:

- Snapshots/backups/restore.
- Clone/create image.
- Deployment history and rollback.
- Serial console / break-glass access.
- Alarms and scheduled events.
- Notebook workspace/app separation.
- Batch retries, timeout, attempts, queue position.

P2 for advanced cloud maturity:

- Hibernation/suspend where supported.
- Automated recovery/migration.
- Cross-region or multi-AZ placement policies.
- Cost allocation profiles and per-application model invocation tracking.
- Canary/blue-green traffic policies for inference.
- Advanced invocation logging controls.

## Release Impact

This record is documentation and planning only.

Expected release impact: `no-release-impact`.

