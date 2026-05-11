# M1-INSTANCE-E Plan Audit Trail

This slice persists the instance planning trail before any provider-side dry-run
or real create/apply action is allowed.

Implemented:

- `ports.WorkloadPlanAuditStore`
- `runtime.MetadataPlanAuditStore`
- `instance_plan_audits` table with tenant RLS

The audit row records:

- Tenant and user.
- Instance ID/name/kind.
- Provider.
- Rendered manifests.
- Admission result and warnings.
- Creation time.

Validation:

```bash
cd repo
make validate-instance-audit
make test
```
