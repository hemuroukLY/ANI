# M1-INSTANCE-S: VM Console and Remote Ops Sessions

## Status

Completed.

## Purpose

This slice closes the missing VM interactive operations boundary. VNC is only
one protocol. The platform must model VM access as a provider-owned remote ops
session so KubeVirt, OpenStack, VMware, and public cloud console mechanisms can
be integrated without leaking provider APIs into business services.

## Scope

- Adds VM console ops actions:
  - `vm_console`
  - `vm_vnc`
  - `vm_serial_console`
- Extends `WorkloadInstanceOpsRequest` with `Protocol`.
- Extends `WorkloadInstanceOpsResult` with:
  - `Protocol`
  - `ConnectURL`
  - `ExpiresAt`
- Adds KubeVirt subresource paths for VM console/VNC in the Kubernetes adapter.
- Adds staged Demo console session API:
  - `POST /api/v1/demo/instances/{instance_id}/console`
  - `POST /api/v1/demo/instances/{instance_id}/console/exec`
- Adds richer Demo lifecycle and ops controls:
  - start
  - stop
  - restart
  - resize
  - delete
  - logs/events/metrics/terminal/exec for container-like instances
  - console session for VM instances

## Guardrails

- VM console operations require a running VM.
- Container terminal/exec remains container-only.
- Business/API/UI layers receive session metadata only; provider-specific
  console/VNC APIs must stay inside adapters.
- Real provider execution remains controlled by provider execution switches.
- Demo console URLs are staged presentation URLs, not production noVNC proxies.
- Demo shell commands execute in a real backend shell workspace scoped to the
  demo VM instance. The workspace is not a real KubeVirt guest OS yet; it is the
  presentation bridge until the noVNC/WebSocket or cloud console proxy is added.
- Demo shell execution uses timeout and basic dangerous-command guardrails.

## Provider Mapping

- KubeVirt: `virtualmachineinstances/{name}/vnc` and console subresources.
- OpenStack: adapter may map to noVNC/spice/serial console URL APIs.
- VMware: adapter may map to WebMKS or remote console ticket APIs.
- Public cloud: adapter may map to cloud console URL/session APIs.

## Validation

```bash
make validate-demo-instances
go test ./pkg/adapters/runtime -run 'TestLocalInstanceServiceVMConsoleOpsCreatesSession|TestKubernetesInstanceOpsVMVNCUsesKubeVirtSubresource'
npm run build
```

## Release Impact

`MINOR`: expands instance ops contract and staged demo behavior without changing
core create/lifecycle provider contracts.
