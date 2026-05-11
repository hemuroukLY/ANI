# M1-INSTANCE-L Instance Service API

This slice defines the business-facing instance service API layer for VM,
container, and GPU container instances.

It covers:

- `ports.WorkloadInstanceService`
- `runtime.LocalInstanceService`
- `Create` through `WorkloadInstanceOrchestrator`
- `Get/List` through `WorkloadInstanceStore`

The service API does not expose provider manifests, Kubernetes/KubeVirt SDK
objects, or provider-specific status details.

Validate with:

```bash
make validate-instance-service
```
