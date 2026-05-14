# M1-INFRA-E GPU Scheduling Baseline

This directory defines the first GPU scheduling baseline for ANI.

It does not install GPU Operator, HAMi, Volcano, or DCGM charts. Those remain
deploy-profile controlled components so ANI can support offline sites and
customer-managed clusters.

This baseline provides:

- GPU system namespaces.
- GPU node label contract.
- Volcano queue templates for inference and training workloads.
- HAMi resource contract.
- DCGM observability contract.
- Offline-verifiable preflight script contract.

Validation:

```bash
cd repo
make validate-infra
```
