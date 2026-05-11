# M1-INSTANCE-Q Kubernetes Lifecycle Execution

This profile adds provider-side lifecycle execution for already-created
instances.

Default behavior remains local/offline. Set
`WORKLOAD_LIFECYCLE_PROVIDER=kubernetes_rest` and
`WORKLOAD_LIFECYCLE_APPLY_ENABLED=true` to execute lifecycle operations through
`KubernetesRESTClient`.

Validate with:

```bash
make validate-kubernetes-lifecycle-execution
```
