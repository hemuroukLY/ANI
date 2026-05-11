# M1-INSTANCE-R Kubernetes Visual Ops Execution

This profile adds provider-side visual operations for container instances.

Default behavior remains local/offline. Set `WORKLOAD_OPS_PROVIDER=kubernetes_rest`
and `WORKLOAD_OPS_ENABLED=true` to execute operations through
`KubernetesRESTClient`.

Validate with:

```bash
make validate-kubernetes-ops-execution
```
