# M1-INSTANCE-A Instance Object, Lifecycle, Network, and Storage Contract

This directory closes the missing platform instance layer before module 3.

It defines every user-visible runtime object as a first-class ANI instance:

- Traditional VM / cloud host.
- Traditional container or Pod instance.
- GPU container or Pod instance.
- Inference instance.
- Notebook instance.
- AI Agent sandbox instance.
- Batch job instance.

Network planning:

- `tenant_vpc` is used for tenant business traffic and VM/container mutual access.
- `foundation_mesh` is a platform-controlled service plane outside tenant VPCs.
- `storage` is isolated storage backend access.
- `management` is control-plane, observability, SSH/VNC proxy, and health access.
- `public_ingress` is explicit gateway or ingress exposure.

This means VM and Pod communication can use the same tenant VPC when they are in
the same business system, while platform, storage, and management dependencies
do not have to be nested inside that VPC.

Offline validation:

```bash
cd repo
make validate-instance-contracts
```
