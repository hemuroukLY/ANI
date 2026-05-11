# M1-RUNTIME-A Workload Runtime / Instance Abstraction

This directory defines the platform-level runtime contract for user-visible
instances.

Instance kinds:

- `vm`: traditional virtual machine or cloud host.
- `container`: traditional Kubernetes/container instance.
- `gpu_container`: container instance that requires GPU scheduling.
- `inference`: model serving instance built on runtime and GPU contracts.
- `notebook`: interactive development instance.
- `agent_sandbox`: restricted AI agent execution sandbox.
- `batch_job`: finite training, evaluation, conversion, or import job.

ANI services must depend on `ports.WorkloadRuntime`. VM-specific logic belongs
in a VM adapter such as KubeVirt or customer cloud provider integration. Pod,
Deployment, Job, RuntimeClass, and Volcano details belong in Kubernetes runtime
adapters.

Offline validation:

```bash
cd repo
make validate-runtime-contracts
```
