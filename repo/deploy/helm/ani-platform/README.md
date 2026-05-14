# ANI Platform Helm Chart

`ani-platform` is the umbrella chart entrypoint for the ANI control plane.

Current scope:

- `M1-INFRA-A` chart metadata and values contract.
- Shared infrastructure dependency contract for PostgreSQL, NATS JetStream, Redis, MinIO, Milvus, and Harbor.
- Service image and port values for the initial Go services.
- `M1-INFRA-B` profile contract:
  - `profiles/dev.yaml`
  - `profiles/attach-k8s.yaml`
  - `profiles/offline.yaml`
  - `component-contracts/*.yaml`

Rendering templates will be added after the raw `deploy/manifests/m1-infra-a` baseline is accepted. The raw manifests remain the first validation target because they can be checked with `kubectl --dry-run=client` in environments where Helm is not installed.

The component contracts intentionally do not download public charts. Chart versions are recorded as compatibility targets and must be pinned into a lockfile before any offline installer package is produced.
