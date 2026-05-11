# M1-INSTANCE-H Status Reconcile

This slice defines the status write-back and lifecycle reconciliation boundary
for instances after provider apply/create.

It covers:

- `ports.WorkloadStatusReconciler`
- `runtime.LocalStatusReconciler`
- normalized provider observations
- provider phase to ANI `WorkloadState` mapping
- audit/apply/resource reference correlation before status write-back

Provider-specific status readers belong in adapters. Business services must not
poll Kubernetes, KubeVirt, or customer cloud status APIs directly.

Validate with:

```bash
make validate-instance-status-reconcile
```
