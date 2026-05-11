# M1-E2E-B Real Provider Regression Profile

This profile validates the unified Kubernetes REST provider path for create,
observe, lifecycle, and visual ops.

It uses a fake HTTP transport, not a live cluster, but keeps the real adapter
chain:

- `KubernetesRESTClient`
- `KubernetesProviderAdapter`
- `KubernetesLifecycleExecutor`
- `KubernetesInstanceOps`

Validate with:

```bash
make validate-m1-real-provider-profile
```
