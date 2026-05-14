# M1-INSTANCE-P Kubernetes Provider Bootstrap Wiring

This profile wires the adapter-owned `KubernetesRESTClient` into bootstrap
configuration.

Default behavior remains local/offline. Set `WORKLOAD_PROVIDER=kubernetes_rest`
and `KUBERNETES_API_HOST` to use the Kubernetes REST client.

Apply remains disabled unless `WORKLOAD_PROVIDER_APPLY_ENABLED=true`.

Validate with:

```bash
make validate-kubernetes-bootstrap-wiring
```
