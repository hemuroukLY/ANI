# M1-INSTANCE-B Instance Planner Adapter

This slice adds the first executable instance-layer logic without binding ANI to
KubeVirt, Kubernetes Pod/Deployment/Job APIs, or any customer cloud provider.

`PlanningRuntime` validates and normalizes:

- VM specs, boot image, and root disk.
- Container and GPU container image requirements.
- Default network planes by instance kind.
- Required `tenant_vpc` business connectivity.
- Storage attachment requirements.
- GPU scheduling dependency on `GPUInventory`.
- Lifecycle transitions for start, stop, restart, resize, and delete.

It is intentionally a planner, not a provider executor. Real KubeVirt,
Kubernetes, Volcano, or cloud adapters must consume the same `WorkloadRuntime`
port after this layer.

Validation:

```bash
cd repo
make validate-instance-planner
make test
```
