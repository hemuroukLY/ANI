# M1-INSTANCE-O Kubernetes REST Client

This profile adds the first adapter-owned `KubernetesProviderClient`
implementation.

It uses standard library HTTP to avoid coupling business services to
client-go/KubeVirt SDKs while still preserving real Kubernetes API semantics:

- server-side dry-run with `dryRun=All`
- server-side apply with `fieldManager` and `force=true`
- observe for Deployment, Job, and KubeVirt VirtualMachine resources

Validate with:

```bash
make validate-kubernetes-rest-client
```
