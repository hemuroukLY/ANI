# M1-INSTANCE-C Provider Renderer Boundary

This slice adds provider manifest dry-run rendering for the instance layer.

Implemented adapter:

- `pkg/adapters/runtime.KubernetesDryRunRenderer`

Provider mapping:

- VM -> KubeVirt `VirtualMachine`
- Container / GPU container / notebook / sandbox / inference -> Kubernetes `Deployment`
- Batch job -> Kubernetes `Job`

The renderer uses `PlanningRuntime` first, then emits JSON-formatted Kubernetes
manifests that are valid YAML documents and can be reviewed or passed to
server-side dry-run tooling later. It does not create cluster resources.

Validation:

```bash
cd repo
make validate-instance-renderer
make test
```
