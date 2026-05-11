# M1-INSTANCE-I Provider Status Reader and Instance Orchestrator

This slice defines the provider status reader and business-facing instance
orchestration API.

It covers:

- `ports.WorkloadProviderStatusReader`
- `ports.WorkloadInstanceOrchestrator`
- `runtime.LocalProviderStatusReader`
- `runtime.LocalInstanceOrchestrator`
- create pipeline sequencing through ANI ports only

Business services should use `WorkloadInstanceOrchestrator` to create instances
instead of manually sequencing provider renderer, dry-run, apply, status reader,
or reconcile calls.

Validate with:

```bash
make validate-instance-orchestrator
```
