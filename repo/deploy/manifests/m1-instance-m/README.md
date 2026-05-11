# M1-INSTANCE-M Lifecycle and Visual Ops API

This slice completes the business-facing lifecycle and visual operations API for
VM, container, and GPU container instances.

It covers:

- `WorkloadInstanceService.Start`
- `WorkloadInstanceService.Stop`
- `WorkloadInstanceService.Restart`
- `WorkloadInstanceService.Resize`
- `WorkloadInstanceService.Delete`
- `WorkloadInstanceService.Ops`
- `WorkloadInstanceOps`
- `LocalInstanceOpsGuard`

Visual operations include logs, events, metrics, terminal, and exec. Provider
SDK calls for these operations must stay inside adapter-owned packages.

Validate with:

```bash
make validate-instance-lifecycle-ops
```
