# M1-INSTANCE-D Admission Guardrail

This slice adds the local admission boundary for rendered instance manifests.

Implemented adapter:

- `pkg/adapters/runtime.LocalAdmissionGuard`

The guardrail validates provider manifests before future provider adapters can
perform server-side dry-run or real create/apply.

Checks:

- Allowed kinds are KubeVirt `VirtualMachine`, Kubernetes `Deployment`, and
  Kubernetes `Job`.
- Tenant and instance labels are required.
- `ani.kubercloud.io/render-mode=dry-run` is required.
- Network plane annotation is required.
- `hostNetwork=true` is denied.
- Privileged containers are denied.

Validation:

```bash
cd repo
make validate-instance-admission
make test
```
