# M1-INSTANCE-K Kubernetes/KubeVirt Provider Adapter

This slice defines the provider adapter boundary for Kubernetes Deployment,
Kubernetes Job, and KubeVirt VirtualMachine execution.

It covers:

- `runtime.KubernetesProviderAdapter`
- `runtime.KubernetesProviderClient`
- server-side dry-run through `dryRun=All`
- disabled-by-default apply execution
- provider status observation through ANI normalized observations

Production Kubernetes/KubeVirt SDK usage must remain inside adapter-owned
packages. Business services continue to use ANI workload ports.

Validate with:

```bash
make validate-instance-provider-adapter
```
