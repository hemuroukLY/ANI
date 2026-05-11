# M1-INSTANCE-N Kubernetes Provider Execution Profile

This profile validates the real-cluster execution boundary before adding a
concrete client-go or KubeVirt client implementation.

It covers:

- `KubernetesProviderClient.ServerSideDryRun` with `dryRun=All`
- controlled `KubernetesProviderClient.Apply`
- `KubernetesProviderClient.Observe`
- audit id, permission proof, admission, dry-run, and resource refs
- `WorkloadInstanceOrchestrator` integration with `KubernetesProviderAdapter`

Validate with:

```bash
make validate-kubernetes-provider-execution
```
