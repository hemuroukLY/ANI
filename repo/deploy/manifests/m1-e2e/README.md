# M1 End-to-End Integration Profile

This profile validates the M1 instance contract chain for:

- VM
- container
- GPU container

It covers:

- create
- lifecycle, including delete
- query
- visual ops: logs, events, metrics, terminal, and exec

The default profile is offline/local and does not require a real Kubernetes
cluster. Production e2e may swap in `KubernetesProviderAdapter` and a real
`KubernetesProviderClient`.

Validate with:

```bash
make validate-m1-e2e
```
