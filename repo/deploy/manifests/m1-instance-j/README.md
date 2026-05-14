# M1-INSTANCE-J Instance Persistence and Query

This slice defines the persistent instance state and query boundary.

It covers:

- `ports.WorkloadInstanceStore`
- `runtime.MetadataInstanceStore`
- `workload_instances` table with tenant RLS
- instance status upsert from `WorkloadInstanceOrchestrator`
- query recovery through `Get` and `List`

Business-facing instance query APIs should use `WorkloadInstanceStore` or a
higher-level service wrapping it instead of relying on `PlanningRuntime` memory.

Validate with:

```bash
make validate-instance-store
```
