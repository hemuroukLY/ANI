# M1-INSTANCE-G Provider Apply Gate

This slice defines the controlled execution gate for provider apply/create.

It covers:

- `ports.WorkloadProviderApply`
- `runtime.LocalProviderApply`
- disabled-by-default execution switch
- required admission, audit, and provider dry-run evidence
- create-only operation allowlist for the first implementation

The local gate validates execution evidence and returns provider resource refs
for tests. It is not a Kubernetes/KubeVirt client and must not bypass future
server-side dry-run, permission, tenant, audit, or idempotency checks.

Validate with:

```bash
make validate-instance-provider-apply
```
