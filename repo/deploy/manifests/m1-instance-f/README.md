# M1-INSTANCE-F Provider Dry-Run Executor

This slice adds the provider dry-run execution boundary for instances.

Implemented:

- `ports.WorkloadProviderDryRun`
- `runtime.LocalProviderDryRun`

The local executor validates provider/kind/apiVersion mappings after local
admission and audit have completed. It does not create cluster resources.

Provider mapping:

- `kubevirt` -> `kubevirt.io/v1` `VirtualMachine`
- `kubernetes` -> `apps/v1` `Deployment`
- `kubernetes` -> `batch/v1` `Job`

Validation:

```bash
cd repo
make validate-instance-provider-dry-run
make test
```
